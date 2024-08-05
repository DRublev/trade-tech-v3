package macd

import (
	"context"
	"main/bot/candles"
	"main/bot/indicators"
	"main/bot/strategies"
	"main/types"
	"math"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	strategies.Config
	// Доступный для торговли баланс
	Balance float32

	// Акция для торговли
	InstrumentID string

	// Каким количчеством акций торговать? Макс
	MaxSharesToHold int64

	// Лотность инструмента
	LotSize int64

	// Если цена пошла ниже чем цена покупки - StopLossAfter, продать по лучшей цене
	// Нужно чтобы  выходить из позиции, когда акция пошла вниз
	StopLossAfter float64
}

// Храним пару последних значений индикаторов
// Чтобы отслеживать их пересечения
var latestSignals = []float64{}
var latestMacd = []float64{}

type isWorking struct {
	sync.RWMutex
	value bool
}

type MacdStrategy struct {
	strategies.IStrategy
	strategies.Strategy[Config]

	provider       candles.BaseCandlesProvider
	activityPubSub strategies.IStrategyActivityPubSub
	vault          strategies.Vault

	isBuying  isWorking
	isSelling isWorking

	macd indicators.MacdIndicator

	toPlaceOrders chan *types.PlaceOrder

	stopCtx context.Context
}

var cancelSwitch context.CancelFunc

func New(provider candles.BaseCandlesProvider, activityPubSub strategies.IStrategyActivityPubSub) *MacdStrategy {
	inst := &MacdStrategy{}
	inst.provider = provider
	inst.activityPubSub = activityPubSub
	inst.stopCtx, cancelSwitch = context.WithCancel(context.Background())
	inst.macd = *indicators.NewMacd(21, 16, 9)
	return inst
}

var l *log.Entry

func (s *MacdStrategy) Start(
	config *strategies.Config,
	ordersToPlaceCh *chan *types.PlaceOrder,
	orderStateChangeCh *chan types.OrderExecutionState,
) (bool, error) {
	l = log.WithFields(log.Fields{
		"strategy":   "macd",
		"instrument": (*config)["InstrumentID"],
	})

	err := s.SetConfig(*config)
	if err != nil {
		l.Errorf("Error parsing config %v", err)
		return false, err
	}

	l.WithField("lotSize", s.Config.LotSize).Infof("Starting strategy with config: %v", s.Config)

	// Создаем или получаем канал, в который будет постаупать инфа о стакане
	l.Tracef("Getting candles channel")
	now := time.Now()

	ch, err := s.provider.GetOrCreate(s.Config.InstrumentID, now.Add(-time.Duration(time.Minute)*5*21), now, true)
	if err != nil {
		l.Errorf("Failed to get candles channel: %v", err)
		return false, err
	}

	s.toPlaceOrders = *ordersToPlaceCh
	go func() {
		l.Info("Start listening latest candles")
		for {
			select {
			case <-s.stopCtx.Done():
				l.Info("Strategy stopped")
				return
			case candle, ok := <-*ch:
				l.Trace("New candle")
				if !ok {
					l.Trace("Candles channel closed")
					return
				}

				go s.onCandle(candle)
			}
		}
	}()

	go s.OnOrderSateChangeSubscribe(s.stopCtx, orderStateChangeCh, s.vault.OnOrderSateChange)

	return true, nil
}

func (s *MacdStrategy) Stop() (bool, error) {
	l.Info("Stopping strategy")
	close(s.toPlaceOrders)
	s.isBuying.value = true
	s.isSelling.value = true
	cancelSwitch()
	return true, nil
}

func (s *MacdStrategy) onCandle(c types.OHLC) {
	wg := &sync.WaitGroup{}

	close := c.Close.Float()
	// TODO: Добавить сюда время. Пересчитывать индикатор, если время в интервале не поменялось
	s.macd.Update(close)
	allMacd, allSignals := s.macd.Get()
	minPeriod := 5
	if len(allMacd) < minPeriod {
		l.Infof("Not enough data for macd")
		return
	}

	latestMacd = allMacd[len(allMacd)-minPeriod:]
	latestSignals = allSignals[len(allSignals)-minPeriod:]

	l.Tracef("Updating signal with new values: signal %v; macd: %v", latestSignals, latestMacd)

	wg.Add(1)
	go s.buy(wg, c)
	wg.Add(1)
	go s.sell(wg, c)

	wg.Wait()
}

func (s *MacdStrategy) buy(wg *sync.WaitGroup, c types.OHLC) {
	defer wg.Done()
	l.Trace("Checking for buy")

	if s.isBuying.value {
		l.Warnf("Not processed prev orderbook item for buy")
		return
	}

	lastIdx := len(latestMacd) - 1
	isNowOver := latestMacd[lastIdx] >= latestSignals[lastIdx]

	isPrevUnder := false
	for i := len(latestMacd) - 1; i >= 0; i-- {
		if latestMacd[i] < latestSignals[i] {
			isPrevUnder = true
			break
		}
	}
	// Если дивергенция растет, то можно войти в позу
	signalEntryPoint := isNowOver && isPrevUnder

	if !signalEntryPoint {
		l.Infof("Not a good entry: macd %v, signal %v", latestMacd, latestSignals)
		return
	}
	l.Infof("Good entry for buy: %v macd, %v signal, %v close price", latestMacd, latestSignals, c.Close.Float())
	leftBalance := s.vault.LeftBalance - s.vault.NotConfirmedBlockedMoney

	canBuySharesAmount := int64(math.Abs(leftBalance / (c.Close.Float() * float64(s.Config.LotSize))))
	if canBuySharesAmount <= 0 {
		l.WithField("state", s.vault.String()).Trace("Can buy 0 shares")
		return
	}

	ok := s.isBuying.TryLock()
	if !ok {
		l.Warn("IsBuiyng mutex cannot be locked")
		return
	}
	defer s.isBuying.Unlock()

	l.Trace("Set is buiyng")
	s.isBuying.value = true
	if canBuySharesAmount > s.Config.MaxSharesToHold {
		l.Tracef("Can buy more shares, than config allows")
		canBuySharesAmount = s.Config.MaxSharesToHold
	}

	order := &types.PlaceOrder{
		InstrumentID: s.Config.InstrumentID,
		Quantity:     int64(canBuySharesAmount),
		Direction:    types.Buy,
		Price:        types.Price(c.Close.Float()),
	}

	l.Infof("Order to place: %v", order)

	priceForAllShares := float64(s.Config.LotSize) * float64(order.Price)
	s.vault.PendingBuyShares += int64(canBuySharesAmount)
	s.vault.NotConfirmedBlockedMoney += float64(canBuySharesAmount) * priceForAllShares
	s.vault.LastBuyPrice = priceForAllShares
	l.WithField("state", s.vault.String()).Trace("State updated after place buy order")

	s.isBuying.value = false
	l.Trace("Is buy released")

	s.toPlaceOrders <- order
}

func (s *MacdStrategy) sell(wg *sync.WaitGroup, c types.OHLC) {
	defer wg.Done()
	l.Trace("Checking for sell")

	if s.isSelling.value {
		l.Warnf("Not processed prev orderbook item for sell")
		return
	}

	if s.vault.HoldingShares-s.vault.PendingSellShares == 0 {
		l.WithField("state", s.vault.String()).Trace("Nothing to sell")
		return
	}

	lastIdx := len(latestMacd) - 1

	isNowUnder := latestMacd[lastIdx] <= latestSignals[lastIdx]
	isPrevOver := false
	for i := len(latestMacd); i >= 0; i-- {
		if latestMacd[i] >= latestSignals[i] {
			isPrevOver = true
			break
		}
	}

	lastBuyPrice := s.vault.LastBuyPrice / float64(s.Config.LotSize)
	hasStopLossBroken := lastBuyPrice-s.Config.StopLossAfter >= c.Close.Float()

	shouldSell := (isNowUnder && isPrevOver) || hasStopLossBroken

	if !shouldSell {
		l.Infof("Not a good exit: macd %v, signal %v", latestMacd, latestSignals)
		return
	}

	ok := s.isSelling.TryLock()
	if !ok {
		l.Warn("isSelling mutex cannot be locked")

		return
	}
	defer s.isSelling.Unlock()

	l.Trace("Set is selling")
	s.isSelling.value = true

	price := c.Close.Float()
	order := &types.PlaceOrder{
		InstrumentID: s.Config.InstrumentID,
		Quantity:     int64(s.vault.HoldingShares),
		Direction:    types.Sell,
		Price:        types.Price(price),
	}
	l.Infof("Order to place: %v", order)

	l.Trace("Updating state")
	s.vault.PendingSellShares += s.vault.HoldingShares
	l.WithField("state", s.vault.String()).Trace("State updated after place sell order")

	s.isSelling.value = false
	l.Trace("Is sell released")

	s.toPlaceOrders <- order

}
