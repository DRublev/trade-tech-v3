package spread

import (
	"context"
	"encoding/json"
	"fmt"
	"main/bot/orderbook"
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

	MinProfit float32

	// Сколько мс ждать после исполнения итерации покупка-продажа перед следующей
	NextOrderCooldownMs int32

	// Каким количчеством акций торговать? Макс
	MaxSharesToHold int32

	// Лотность инструмента
	LotSize int32

	// Если цена пошла ниже чем цена покупки - StopLossAfter, продать по лучшей цене
	// Нужно чтобы  выходить из позиции, когда акция пошла вниз
	StopLossAfter float32
}

type State struct {
	// Оставшееся количество денег
	leftBalance float32

	// Сумма, которая должна списаться при выставлении ордера на покупку
	// Инкрементим когда хотим выставить бай ордер
	// Декрементим когда закрываем бай ордер
	notConfirmedBlockedMoney float32

	// Количество акций, купленных на данный момент
	holdingShares int32

	// Количество акций, на которое выставлены ордера на покупку
	pendingBuyShares int32

	// Количество акций, на которое выставлены ордера на продажу
	pendingSellShares int32

	lastBuyPrice float32

	placedOrders []types.OrderExecutionState
}

func (s *State) String() string {
	return fmt.Sprintf(
		"Holding %v\nLeft balance %v; Blocked money %v\nPending buy %v, sell %v\nLast buy price %v",
		s.holdingShares,
		s.leftBalance,
		s.notConfirmedBlockedMoney,
		s.pendingBuyShares,
		s.pendingSellShares,
		s.lastBuyPrice,
	)
}

type isWorking struct {
	sync.RWMutex
	value bool
}

type SpreadStrategy struct {
	strategies.IStrategy
	strategies.Strategy[Config]
	config Config
	// Канал для стакана
	obCh              *chan *types.Orderbook
	state             strategies.StrategyState[State]
	nextOrderCooldown *time.Timer
	isBuying          isWorking
	isSelling         isWorking

	toPlaceOrders chan *types.PlaceOrder

	stopCtx context.Context
}

var cancelSwitch context.CancelFunc

func New() *SpreadStrategy {
	inst := &SpreadStrategy{}
	inst.toPlaceOrders = make(chan *types.PlaceOrder)
	inst.stopCtx, cancelSwitch = context.WithCancel(context.Background())
	return inst
}

var l *log.Entry

func (s *SpreadStrategy) Start(
	config *strategies.Config,
	ordersToPlaceCh *chan *types.PlaceOrder,
	orderStateChangeCh *chan types.OrderExecutionState,
) (bool, error) {
	l = log.WithFields(log.Fields{
		"strategy":   "spread",
		"instrument": (*config)["InstrumentID"],
	})

	var res Config

	bts, err := json.Marshal(config)
	if err != nil {
		l.Errorf("Error parsing config %v", err)
		return false, err
	}

	err = json.Unmarshal(bts, &res)
	if err != nil {
		l.Errorf("Error parsing config %v", err)
		return false, err
	}
	s.config = res

	l.Infof("Starting strategy with config: %v", s.config)

	// Создаем или получаем канал, в который будет постаупать инфа о стакане
	l.Tracef("Getting orderbook channel")
	obProvider := orderbook.NewProvider()
	ch, err := obProvider.GetOrCreate(s.config.InstrumentID)
	if err != nil {
		l.Errorf("Failed to get orderbook channel: %v", err)
		return false, err
	}

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

	// стакан!
	s.obCh = ch

	go func(ch *chan *types.Orderbook) {
		l.Info("Start listening changes in orderbook")
		for {
			select {
			case <-s.stopCtx.Done():
				l.Info("Strategy stopped")
				return
			case ob, ok := <-*ch:
				l.Trace("New orderbook change")
				if !ok {
					l.Trace("Orderbook channel closed")
					return
				}

				go s.onOrderbook(ob)
			}
		}
	}(ch)

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

	return true, nil
}

func (s *SpreadStrategy) Stop() (bool, error) {
	l.Info("Stopping strategy")
	close(s.toPlaceOrders)
	s.isBuying.value = true
	s.isSelling.value = true
	cancelSwitch()
	return true, nil
}

func (s *SpreadStrategy) onOrderbook(ob *types.Orderbook) {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go s.checkForRottenBuys(wg, ob)
	wg.Add(1)
	go s.checkForRottenSells(wg, ob)

	wg.Add(1)
	go s.buy(wg, ob)
	wg.Add(1)
	go s.sell(wg, ob)

	wg.Wait()
}

// TODO: Перенести в отдельный файл
func (s *SpreadStrategy) buy(wg *sync.WaitGroup, ob *types.Orderbook) {
	defer wg.Done()
	l.Trace("Checking for buy")

	if s.isBuying.value {
		l.Warnf("Not processed prev orderbook item for buy")
		return
	}

	isHoldingMaxShares := s.state.Get().holdingShares+s.state.Get().pendingBuyShares >= s.config.MaxSharesToHold
	if isHoldingMaxShares {
		l.WithField("state", s.state.Get()).Tracef("Cannot buy, holding max shares")
		return
	}

	// Аукцион закрытия, только заявки на продажу
	if len(ob.Bids) == 0 {
		l.Trace("No bids")
		return
	}

	minBuyPrice := ob.Bids[0].Price
	l.Tracef("Min buy price: %v", minBuyPrice)
	leftBalance := s.state.Get().leftBalance - s.state.Get().notConfirmedBlockedMoney
	if leftBalance < (minBuyPrice * float32(s.config.LotSize)) {
		l.WithField("state", s.state.Get()).Tracef("Not enough money")
		return
	}

	canBuySharesAmount := int32(math.Abs(float64(leftBalance / (minBuyPrice * float32(s.config.LotSize)))))
	l.Tracef("First bid price: %v; Left money: %v; Can buy %v shares\n", minBuyPrice, leftBalance, canBuySharesAmount)
	if canBuySharesAmount <= 0 {
		l.WithField("state", s.state.Get()).Trace("Can buy 0 shares")
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
	if canBuySharesAmount > s.config.MaxSharesToHold {
		l.Tracef("Can buy more shares, than config allows")
		canBuySharesAmount = s.config.MaxSharesToHold
	}

	order := &types.PlaceOrder{
		InstrumentID: s.config.InstrumentID,
		Quantity:     int64(canBuySharesAmount),
		Price:        types.Price(minBuyPrice),
		Direction:    types.Buy,
	}
	l.Infof("Order to place: %v", order)

	l.Trace("Updating state")
	newState := *s.state.Get()
	newState.pendingBuyShares += canBuySharesAmount
	newState.notConfirmedBlockedMoney += float32(canBuySharesAmount) * minBuyPrice
	newState.lastBuyPrice = minBuyPrice
	s.state.Set(newState)
	l.WithField("state", s.state.Get()).Trace("State updated after place buy order")

	s.isBuying.value = false
	l.Trace("Is buy released")

	s.toPlaceOrders <- order
}

func (s *SpreadStrategy) sell(wg *sync.WaitGroup, ob *types.Orderbook) {
	defer wg.Done()
	l.Trace("Checking for sell")

	if s.isSelling.value {
		l.Warnf("Not processed prev orderbook item for sell")
		return
	}

	state := *s.state.Get()

	if state.holdingShares < 0 {
		l.WithField("state", state).Warn("Holding less than 0 shares")
		return
	}

	minAskPrice := ob.Asks[0].Price
	l.Tracef("Min ask price %v", minAskPrice)

	isGoodPrice := minAskPrice-state.lastBuyPrice >= s.config.MinProfit
	hasStopLossBroken := state.holdingShares-state.pendingSellShares > 0 && s.config.StopLossAfter != float32(0) && ob.Bids[0].Price <= state.lastBuyPrice-s.config.StopLossAfter

	if state.holdingShares-state.pendingSellShares == 0 && !hasStopLossBroken {
		l.WithField("state", state).Trace("Nothing to sell")
		return
	}

	fmt.Printf("320 strategy %v <= %v - %v (%vis %v\n", ob.Bids[0].Price, state.lastBuyPrice, s.config.StopLossAfter, state.lastBuyPrice-s.config.StopLossAfter, hasStopLossBroken)
	shouldMakeSell := isGoodPrice || hasStopLossBroken
	if !shouldMakeSell {
		l.WithField("lastBuyPrice", state.lastBuyPrice).Tracef("Not a good deal")
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

	price := minAskPrice
	l.Tracef("Selling price: %v", price)
	if hasStopLossBroken {
		price = ob.Bids[0].Price
		l.WithFields(log.Fields{
			"lastBuyPrice":  state.lastBuyPrice,
			"stopLoss":      s.config.StopLossAfter,
			"stopLossPrice": state.lastBuyPrice - s.config.StopLossAfter,
		}).Info("Stop loss broken")
	}

	order := &types.PlaceOrder{
		InstrumentID: s.config.InstrumentID,
		Quantity:     int64(state.holdingShares),
		Direction:    types.Sell,
		Price:        types.Price(price),
	}

	if hasStopLossBroken {
		for _, o := range state.placedOrders {
			if (o.Direction == types.Buy || o.Direction == types.Sell) && o.Status != types.New {
				order.CancelOrder = o.ID
				break
			}
		}
	}

	l.Infof("Order to place: %v", order)

	l.Trace("Updating state")
	state = *s.state.Get()
	state.pendingSellShares += state.holdingShares
	s.state.Set(state)
	l.WithField("state", s.state.Get()).Trace("State updated after place sell order")

	s.isSelling.value = false
	l.Trace("Is sell released")

	s.toPlaceOrders <- order

}

func (s *SpreadStrategy) checkForRottenBuys(wg *sync.WaitGroup, ob *types.Orderbook) {
	defer wg.Done()
	// TODO: Чекать неаткуальные выставленные ордера и отменять их
	// TODO: Сбрасывать lastBuyPrice на предыдущий, если закрываем какой то бай ордер
}

func (s *SpreadStrategy) checkForRottenSells(wg *sync.WaitGroup, ob *types.Orderbook) {
	defer wg.Done()
	// TODO: Чекать неаткуальные выставленные ордера и отменять их
}

func (s *SpreadStrategy) onOrderSateChange(state types.OrderExecutionState) {
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
		newState.leftBalance += float32(state.ExecutedOrderPrice)
		newState.pendingBuyShares -= int32(state.LotsExecuted / int(s.config.LotSize))
		newState.notConfirmedBlockedMoney -= float32(state.ExecutedOrderPrice)
	} else if isSellPlaceError {
		newState.pendingSellShares -= int32(state.LotsExecuted / int(s.config.LotSize))
	}

	if isSellOk || isBuyCancel {
		l.Trace("Updating state after sell order executed")
		newState.pendingSellShares -= int32(state.LotsExecuted / int(s.config.LotSize))
		newState.leftBalance += float32(state.ExecutedOrderPrice)
		newState.holdingShares -= int32(state.LotsExecuted / int(s.config.LotSize))
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
		newState.holdingShares += int32(state.LotsExecuted / int(s.config.LotSize))
		newState.pendingBuyShares -= int32(state.LotsExecuted / int(s.config.LotSize))
		newState.notConfirmedBlockedMoney -= float32(state.ExecutedOrderPrice)
		newState.leftBalance -= float32(state.ExecutedOrderPrice)
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
