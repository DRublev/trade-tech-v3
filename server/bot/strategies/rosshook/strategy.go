package rosshook

import (
	"context"
	"encoding/json"
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
	strategies.Strategy

	provider candles.BaseCandlesProvider
	config   Config
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

	var res Config

	// TODO: Вынести в сущность конфига стратегии
	bts, err := json.Marshal(config)
	if err != nil {
		l.Error("Error parsing config %v", err)
		return false, err
	}

	err = json.Unmarshal(bts, &res)
	if err != nil {
		l.Error("Error parsing config %v", err)
		return false, err
	}
	s.config = res

	s.vault = *strategies.NewVault(s.config.LotSize, s.config.Balance)

	l.Infof("Starting strategy with config: %v", s.config)

	// Создаем или получаем канал, в который будет постаупать инфа о стакане
	l.Tracef("Getting candles channel")
	now := time.Now()

	ch, err := s.provider.GetOrCreate(s.config.InstrumentID, now, now)
	if err != nil {
		l.Errorf("Failed to get candles channel: %v", err)
		return false, err
	}

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

	go func() {
		l.Info("Start listening for orders")
		for {
			select {
			case <-s.stopCtx.Done():
				l.Info("Strategy stopped")
				return
			case state, ok := <-*orderStateChangeCh:
				if !ok {
					l.Warn("Orders state channel closed")
					return
				}
				go s.vault.OnOrderSateChange(state)
			case orderToPlace, ok := <-s.toPlaceOrders:
				if !ok {
					l.Warn("Place orders channel closed")
					return
				}
				*ordersToPlaceCh <- orderToPlace
			}
		}
	}()

	s.nextOrderCooldown = time.NewTimer(time.Duration(0) * time.Millisecond)

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

func (s *RossHookStrategy) watchBuySignal(c types.OHLC) {
	if high == nil || high.High.Float() <= c.High.Float() {
		high = &c
		low = nil
		targetGrow = nil
		less = nil
		l.Infof("Set point 1. high: %v;", high.High.Float())
	} else if low == nil || low.Low.Float() >= c.Low.Float() {
		low = &c
		targetGrow = nil
		less = nil
		l.Infof("Set point 2. high: %v; low: %v;", high.High.Float(), low.Low.Float())
	} else if targetGrow == nil || (targetGrow.High.Float() < c.High.Float() && less == nil) {
		targetGrow = &c
		less = nil
		takeProfit = &c
		l.Infof("Set point 3. high: %v; low: %v; targetGrow: %v;", high.High.Float(), low.Low.Float(), targetGrow.High.Float())
	} else if less == nil || less.Low.Float() >= c.Low.Float() {
		less = &c
		l.Infof("Set point 4. high: %v; low: %v; targetGrow: %v; less: %v;", high.High.Float(), low.Low.Float(), targetGrow.High.Float(), less.Low.Float())
	}

	if high != nil && low != nil && targetGrow != nil && less != nil {
		if targetGrow.High.Float() <= c.High.Float() {
			go s.buy(*targetGrow)
		}
	}
}

func (s *RossHookStrategy) watchSellSignal(c types.OHLC) {
	// Stop-loss
	if high != nil && low != nil && targetGrow != nil && less != nil {
		if less.Close.Float()-s.config.StopLoss >= c.Close.Float() {
			go s.sell(c)
			return
		}
	}

	if takeProfit == nil {
		return
	}
	if takeProfit.High.Float() < c.High.Float() {
		takeProfit = &c
	} else if takeProfit.Close.Float()-float64(s.config.SaveProfit) >= c.Close.Float() {
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
		l.WithField("state", s.vault).Trace("Nothing to sell")
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
		InstrumentID: s.config.InstrumentID,
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
	s.toPlaceOrders <- order
}

func (s *RossHookStrategy) buy(c types.OHLC) {
	l.Trace("Checking for buy")

	if s.isBuying.value {
		l.Warnf("Not processed prev orderbook item for buy")
		return
	}

	leftBalance := s.vault.LeftBalance - s.vault.NotConfirmedBlockedMoney

	canBuySharesAmount := math.Round(math.Abs(leftBalance / (c.Close.Float() * float64(s.config.LotSize))))
	fmt.Printf("266 strategy lotSize %v; left balance %v; can buy %v \n", s.config.LotSize, leftBalance, canBuySharesAmount)
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
	if canBuySharesAmount > float64(s.config.MaxSharesToHold) {
		l.Tracef("Can buy more shares, than config allows")
		canBuySharesAmount = float64(s.config.MaxSharesToHold)
	}

	order := &types.PlaceOrder{
		InstrumentID: s.config.InstrumentID,
		Quantity:     int64(canBuySharesAmount),
		Direction:    types.Buy,
		Price:        types.Price(c.Close.Float()),
	}

	l.Infof("Order to place: %v", order)

	s.vault.PendingBuyShares += int64(canBuySharesAmount)
	s.vault.NotConfirmedBlockedMoney += canBuySharesAmount * c.Close.Float()
	s.vault.LastBuyPrice = c.Close.Float()
	l.WithField("state", s.vault).Trace("State updated after place buy order")

	s.isBuying.value = false
	l.Trace("Is buy released")

	s.toPlaceOrders <- order
}
