package rosshook

import (
	"context"
	"fmt"
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
	Balance float64

	// Акция для торговли
	InstrumentID string

	// Каким количчеством акций торговать? Макс
	MaxSharesToHold int64

	// Лотность инструмента
	LotSize int64

	// При падении ниже 2 точки минус этот парамер выставим продажу
	StopLoss float64

	// Нужен для Trailing take profit
	// При какой просадке от максимума выставить продажу
	SaveProfit float64
}

type isWorking struct {
	sync.RWMutex
	value bool
}

type RossHookStrategy struct {
	strategies.IStrategy
	strategies.Strategy[Config]

	provider candles.BaseCandlesProvider
	// Канал для стакана
	obCh              *chan *types.Orderbook
	nextOrderCooldown *time.Timer
	isBuying          isWorking
	isSelling         isWorking

	macd indicators.MacdIndicator

	toPlaceOrders chan *types.PlaceOrder

	stopCtx context.Context

	vault strategies.Vault
}

var cancelSwitch context.CancelFunc

func New(provider candles.BaseCandlesProvider) *RossHookStrategy {
	inst := &RossHookStrategy{}
	inst.provider = provider
	inst.toPlaceOrders = make(chan *types.PlaceOrder)
	inst.stopCtx, cancelSwitch = context.WithCancel(context.Background())
	return inst
}

var l *log.Entry

func (s *RossHookStrategy) Start(
	config *strategies.Config,
	ordersToPlaceCh *chan *types.PlaceOrder,
	orderStateChangeCh *chan types.OrderExecutionState,
) (bool, error) {
	l = log.WithFields(log.Fields{
		"strategy":   "rosshook",
		"instrument": (*config)["InstrumentID"],
	})

	err := s.SetConfig(*config)

	if err != nil {
		l.Error("Error parsing config %v", err)
		return false, err
	}

	s.vault = *strategies.NewVault(s.Config.LotSize, s.Config.Balance)

	l.Infof("Starting strategy with config: %v", s.Config)

	// Создаем или получаем канал, в который будет постаупать инфа о стакане
	l.Tracef("Getting candles channel")
	now := time.Now()

	ch, err := s.provider.GetOrCreate(s.Config.InstrumentID, now, now, false)
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

				go s.OnCandle(candle)
			}
		}
	}()

	go s.OnOrderSateChangeSubscribe(s.stopCtx, orderStateChangeCh, s.vault.OnOrderSateChange)

	return true, nil
}

func (s *RossHookStrategy) Stop() (bool, error) {
	l.Info("Stopping strategy")
	close(s.toPlaceOrders)
	s.isBuying.value = true
	s.isSelling.value = true
	cancelSwitch()
	return true, nil
}

var candlesHistory = []types.OHLC{}

var high *types.OHLC
var low *types.OHLC
var targetGrow *types.OHLC
var less *types.OHLC
var buy float64
var takeProfit *types.OHLC

func (s *RossHookStrategy) OnCandle(c types.OHLC) {
	candlesHistory = append(candlesHistory, c)

	if s.isBuying.value || s.isSelling.value {
		return
	}

	s.watchBuySignal(c)

	s.watchSellSignal(c)
}

var lastBuyPendingCandle *types.OHLC

func isSameTF(candidate types.OHLC, toCompare types.OHLC) bool {
	cH, cM, _ := candidate.Time.Clock()
	nH, nM, _ := toCompare.Time.Clock()
	return cH == nH && cM == nM
}

func (s *RossHookStrategy) watchBuySignal(c types.OHLC) {
	// Закрываем висящие на заявку покупки при поступлении новой свечи - мы проебали момент
	// Однако, если текущая цена равна цене, по которой выставляли заявку, есть шанс что еще исполнится
	if lastBuyPendingCandle != nil && !isSameTF(c, *lastBuyPendingCandle) && c.High.Float() > lastBuyPendingCandle.High.Float() {
		s.closePendingBuys()
	}

	isLowDifTF := high != nil && !isSameTF(c, *high)
	isTargetGrowDiffTF := isLowDifTF && low != nil && !isSameTF(c, *low)
	isLessDiffTF := isTargetGrowDiffTF && targetGrow != nil && !isSameTF(c, *targetGrow)

	if high == nil || high.High.Float() <= c.High.Float() {
		high = &c
		low = nil
		targetGrow = nil
		less = nil
		l.Infof("Set point 1. high: %v;", high.High.Float())
		return
	} else if low == nil || (low.Low.Float() >= c.Low.Float() && isLowDifTF) {
		low = &c
		targetGrow = nil
		less = nil
		l.Infof("Set point 2. high: %v; low: %v;", high.High.Float(), low.Low.Float())
		return
	} else if targetGrow == nil || (targetGrow.High.Float() < c.High.Float() && less == nil && isTargetGrowDiffTF) {
		targetGrow = &c
		less = nil
		takeProfit = &c
		l.Infof("Set point 3. high: %v; low: %v; targetGrow: %v;", high.High.Float(), low.Low.Float(), targetGrow.High.Float())
		return
	} else if less == nil || (less.Low.Float() >= c.Low.Float() && less.Low.Float() < targetGrow.High.Float() && isLessDiffTF) {
		less = &c
		l.Infof("Set point 4. high: %v; low: %v; targetGrow: %v; less: %v;", high.High.Float(), low.Low.Float(), targetGrow.High.Float(), less.Low.Float())
		return
	}

	if high != nil && low != nil && targetGrow != nil && less != nil {
		if targetGrow.High.Float() <= c.High.Float() {
			go s.buy(*targetGrow)
		}
	}
}

func (s *RossHookStrategy) watchSellSignal(c types.OHLC) {
	// Stop-loss
	if less != nil && less.Close.Float()-s.Config.StopLoss >= c.Close.Float() {
		l.Infof("Placing stop-loss (less: %v; loss: %v; current: %v)", less.Close.Float(), s.Config.StopLoss, c.Close.Float())
		go s.sell(c)
		return
	}

	if takeProfit == nil || less == nil {
		return
	}
	if takeProfit.High.Float() <= c.High.Float() {
		// Копируем свечу. Подозрение на баг, что свеча перезаписывается следующей, поэтомуне прокидываем просто &c
		takeProfit = &types.OHLC{
			Open:  c.Open,
			High:  c.High,
			Low:   c.Low,
			Close: c.Close,
			Time:  c.Time,
		}
	} else if takeProfit.High.Float()-float64(s.Config.SaveProfit) >= c.Close.Float() {
		l.Infof("Placing take-profit (take: %v; save: %v; current: %v)", takeProfit.High.Float(), s.Config.SaveProfit, c.Close.Float())
		go s.sell(*takeProfit)
	}
}

func (s *RossHookStrategy) sell(c types.OHLC) {
	l.Trace("Checking for sell")

	if s.isSelling.value {
		l.Warnf("Not processed prev orderbook item for sell")
		return
	}

	if s.vault.HoldingShares-s.vault.PendingSellShares == 0 {
		l.WithField("state", s.vault).Info("Nothing to sell")
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
	l.WithField("state", s.vault).Trace("State updated after place sell order")

	s.isSelling.value = false
	l.Trace("Is sell released")

	high = nil
	low = nil
	targetGrow = nil
	less = nil
	lastBuyPendingCandle = nil

	s.toPlaceOrders <- order
}

func (s *RossHookStrategy) buy(c types.OHLC) {
	l.Trace("Checking for buy")

	if s.isBuying.value {
		l.Warnf("Not processed prev orderbook item for buy")
		return
	}

	leftBalance := s.vault.LeftBalance - s.vault.NotConfirmedBlockedMoney

	canBuySharesAmount := int64(math.Abs(leftBalance / (c.Close.Float() * float64(s.Config.LotSize))))
	fmt.Printf("266 strategy lotSize %v; left balance %v; can buy %v \n", s.Config.LotSize, leftBalance, canBuySharesAmount)
	if canBuySharesAmount <= 0 {
		l.WithField("state", s.vault).Trace("Can buy 0 shares")
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

	s.vault.PendingBuyShares += int64(canBuySharesAmount)
	s.vault.NotConfirmedBlockedMoney += float64(canBuySharesAmount) * c.Close.Float()
	s.vault.LastBuyPrice = c.Close.Float()
	lastBuyPendingCandle = &c
	l.WithField("state", s.vault).Trace("State updated after place buy order")

	s.isBuying.value = false
	l.Trace("Is buy released")

	takeProfit = nil

	s.toPlaceOrders <- order
}

func (s *RossHookStrategy) closePendingBuys() {
	l.Infof("Pending buys: %v", len(s.vault.PlacedBuyOrders))
	for _, order := range s.vault.PlacedBuyOrders {
		o := &types.PlaceOrder{
			InstrumentID: s.Config.InstrumentID,
			CancelOrder:  order.ID,
		}
		s.toPlaceOrders <- o
	}
}
