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

	MinProfit float64

	// Каким количчеством акций торговать? Макс
	MaxSharesToHold int64

	// Лотность инструмента
	LotSize int64

	// Если цена пошла ниже чем цена покупки - StopLossAfter, продать по лучшей цене
	// Нужно чтобы  выходить из позиции, когда акция пошла вниз
	StopLossAfter float64

	SaveProfit float64
}

type State struct {
	// Оставшееся количество денег
	leftBalance float64

	// Сумма, которая должна списаться при выставлении ордера на покупку
	// Инкрементим когда хотим выставить бай ордер
	// Декрементим когда закрываем бай ордер
	notConfirmedBlockedMoney float64

	// Количество акций, купленных на данный момент
	holdingShares int64

	// Количество акций, на которое выставлены ордера на покупку
	pendingBuyShares int64

	// Количество акций, на которое выставлены ордера на продажу
	pendingSellShares int64

	lastBuyPrice float64

	placedOrders []types.OrderExecutionState
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
	state             strategies.StrategyState[State]
	nextOrderCooldown *time.Timer
	isBuying          isWorking
	isSelling         isWorking

	macd indicators.MacdIndicator

	toPlaceOrders chan *types.PlaceOrder

	stopCtx context.Context
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
		"strategy":   "ross_hook",
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

	l.Infof("Starting strategy with config: %v", s.config)

	// Создаем или получаем канал, в который будет постаупать инфа о стакане
	l.Tracef("Getting candles channel")
	now := time.Now()

	ch, err := s.provider.GetOrCreate(s.config.InstrumentID, now, now)
	if err != nil {
		l.Errorf("Failed to get candles channel: %v", err)
		return false, err
	}

	go func(c *chan types.OHLC) {
		l.Info("Start listening latest candles")
		for {
			select {
			case <-s.stopCtx.Done():
				l.Info("Strategy stopped")
				return
			case candle, ok := <-*c:
				l.Trace("New candle")
				if !ok {
					l.Trace("Candles channel closed")
					return
				}

				go s.OnCandle(candle)
			}
		}
	}(ch)

	go func(source *chan types.OrderExecutionState) {
		l.Info("Start listening for orders state changes")
		for {
			select {
			case <-s.stopCtx.Done():
				l.Info("Strategy stopped")
				return
			case state, ok := <-*source:
				if !ok {
					l.Warn("Orders state channel closed")
					return
				}
				go s.onOrderSateChange(state)
			}

		}
	}(orderStateChangeCh)

	go func(source *chan *types.PlaceOrder, target *chan *types.PlaceOrder) {
		l.Info("Start listening for new place order requests")
		for {
			select {
			case <-s.stopCtx.Done():
				l.Info("Strategy stopped")
				return
			case orderToPlace, ok := <-*source:
				if !ok {
					l.Warn("Place orders channel closed")
					return
				}
				*target <- orderToPlace
			}
		}
	}(&s.toPlaceOrders, ordersToPlaceCh)

	l.Trace("Setting state to empty")
	// Заполняем изначальное состояние
	s.state = strategies.StrategyState[State]{}
	err = s.state.Set(State{
		holdingShares:            0,
		pendingBuyShares:         0,
		pendingSellShares:        0,
		leftBalance:              s.config.Balance,
		notConfirmedBlockedMoney: 0,
		lastBuyPrice:             0,
	})
	if err != nil {
		l.Errorf("Failed to set strategy initial state: %v", err)
		return false, err
	}

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
	if high == nil || high.High.Float() < c.High.Float() {
		high = &c
		low = nil
		targetGrow = nil
		less = nil
	} else if low == nil || low.Low.Float() > c.Low.Float() {
		low = &c
		targetGrow = nil
		less = nil
	} else if targetGrow == nil || targetGrow.High.Float() < c.High.Float() {
		targetGrow = &c
		less = nil
		takeProfit = &c
	} else if less == nil || less.Low.Float() > c.Low.Float() {
		less = &c
	}

	if high != nil && low != nil && targetGrow != nil && less != nil {
		if targetGrow.High.Float() <= c.High.Float() {
			go s.buy(c)
		}
	}
}

func (s *RossHookStrategy) watchSellSignal(c types.OHLC) {
	if high != nil && low != nil && targetGrow != nil && less != nil {
		if less.Close.Float() >= c.Close.Float() {
			go s.sell(c)
		}
	}

	if takeProfit == nil {
		return
	}
	if takeProfit.High.Float() < c.High.Float() {
		takeProfit = &c
	} else if takeProfit.Close.Float()-float64(s.config.SaveProfit) >= c.Close.Float() {
		go s.sell(c)
	}
}

func (s *RossHookStrategy) sell(c types.OHLC) {
	l.Trace("Checking for sell")

	if s.isSelling.value {
		l.Warnf("Not processed prev orderbook item for sell")
		return
	}

	state := *s.state.Get()

	if state.holdingShares-state.pendingSellShares == 0 {
		l.WithField("state", state).Trace("Nothing to sell")
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
		Quantity:     int64(state.holdingShares),
		Direction:    types.Sell,
		Price:        types.Price(price),
	}
	l.Infof("Order to place: %v", order)

	l.Trace("Updating state")
	state.pendingSellShares += state.holdingShares
	s.state.Set(state)
	l.WithField("state", s.state.Get()).Trace("State updated after place sell order")

	s.isSelling.value = false
	l.Trace("Is sell released")

	s.toPlaceOrders <- order
}

func (s *RossHookStrategy) buy(c types.OHLC) {
	l.Trace("Checking for buy")

	if s.isBuying.value {
		l.Warnf("Not processed prev orderbook item for buy")
		return
	}

	// TODO: Выставить ордер на покупку

	state := s.state.Get()
	leftBalance := state.leftBalance - state.notConfirmedBlockedMoney

	canBuySharesAmount := math.Abs(leftBalance / (c.Close.Float() * float64(s.config.LotSize)))
	fmt.Printf("266 strategy lotSize %v; left balance %v; can buy %v \n", s.config.LotSize, leftBalance, canBuySharesAmount)
	if canBuySharesAmount == 0 {
		l.WithField("state", state).Trace("Can buy 0 shares")
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

	state.pendingBuyShares += int64(canBuySharesAmount)
	state.notConfirmedBlockedMoney += canBuySharesAmount * c.Close.Float()
	state.lastBuyPrice = c.Close.Float()
	s.state.Set(*state)
	l.WithField("state", s.state.Get()).Trace("State updated after place buy order")

	s.isBuying.value = false
	l.Trace("Is buy released")

	s.toPlaceOrders <- order
}

func (s *RossHookStrategy) onOrderSateChange(state types.OrderExecutionState) {
	l.Infof("Order state changed %v", state)

	if state.Status == types.ErrorPlacing {
		l.Error("Order placing error. State restored")
	}

	newState := *s.state.Get()
	defer l.WithField("state", s.state.Get()).Info("State updated")

	if state.Status == types.New {
		newState.placedOrders = append(newState.placedOrders, state)
		s.state.Set(newState)
		l.Infof("Adding new order to placed list")
		return
	}
	if state.Status == types.Fill {
		filteredOrders := []types.OrderExecutionState{}

		for _, order := range newState.placedOrders {
			if order.ID != state.ID {
				filteredOrders = append(filteredOrders, order)
			}
		}

		newState.placedOrders = filteredOrders
	}

	if state.Status != types.PartiallyFill &&
		state.Status != types.Fill &&
		state.Status != types.ErrorPlacing &&
		state.Status != types.Cancelled {
		l.Warnf("Not processed order state change: %v", state)
		return
	}

	isBuyPlaceError := state.Direction == types.Buy && state.Status == types.ErrorPlacing
	isSellPlaceError := state.Direction == types.Sell && state.Status == types.ErrorPlacing
	isBuyCancel := state.Direction == types.Buy && state.Status == types.Cancelled
	isSellCancel := state.Direction == types.Sell && state.Status == types.Cancelled
	isSellOk := state.Direction == types.Sell && !isSellPlaceError && !isSellCancel
	isBuyOk := state.Direction == types.Buy && !isBuyPlaceError && !isBuyCancel

	if isBuyPlaceError {
		l.Trace("Updating state after buy order place error")
		newState.leftBalance += state.ExecutedOrderPrice
		newState.pendingBuyShares -= int64(state.LotsExecuted / int(s.config.LotSize))
		newState.notConfirmedBlockedMoney -= state.ExecutedOrderPrice
	} else if isSellPlaceError {
		newState.pendingSellShares -= int64(state.LotsExecuted / int(s.config.LotSize))
	}

	if isSellOk || isBuyCancel {
		l.Trace("Updating state after sell order executed")
		newState.pendingSellShares -= int64(state.LotsExecuted / int(s.config.LotSize))
		newState.leftBalance += state.ExecutedOrderPrice
		newState.holdingShares -= int64(state.LotsExecuted / int(s.config.LotSize))
		l.WithField("orderId", state.ID).Infof(
			"Lots executed (cancelled %v, erroPlacing: %v) %v of %v; Executed sell price %v",
			isBuyCancel,
			isBuyPlaceError,
			state.LotsExecuted,
			state.LotsRequested,
			state.ExecutedOrderPrice,
		)
	} else if isBuyOk || isSellPlaceError || isSellCancel {
		l.Trace("Updating state after buy order executed")
		newState.holdingShares += int64(state.LotsExecuted / int(s.config.LotSize))
		newState.pendingBuyShares -= int64(state.LotsExecuted / int(s.config.LotSize))
		newState.notConfirmedBlockedMoney -= state.ExecutedOrderPrice
		newState.leftBalance -= state.ExecutedOrderPrice
		newState.lastBuyPrice = state.ExecutedOrderPrice / float64(state.LotsExecuted)
		l.WithField("orderId", state.ID).Infof(
			"Lots executed (cancelled %v, erroPlacing: %v) %v of %v; Executed buy price %v",
			isSellCancel,
			isSellPlaceError,
			state.LotsExecuted,
			state.LotsRequested,
			state.ExecutedOrderPrice,
		)
	} else {
		l.Warnf("Order state change not handled: %v", state)
	}
	s.state.Set(newState)
}
