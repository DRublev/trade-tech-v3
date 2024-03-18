package spread

import (
	"fmt"
	"main/bot/orderbook"
	"main/bot/strategies"
	"main/types"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	strategies.Config

	// Минимальная разница bid-ask, при которой выставлять ордер
	minProfit float32

	// Сколько мс ждать после исполнения итерации покупка-продажа перед следующей
	nextOrderCooldownMs int32

	// Каким количчеством акций торговать? Макс
	maxSharesToHold int32

	// Лотность инструмента
	lotSize int32

	// Если цена пошла ниже чем цена покупки - stopLossAfter, продать по лучшей цене
	// Нужно чтобы  выходить из позиции, когда акция пошла вниз
	stopLossAfter float32
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
	strategies.Strategy
	config Config
	// Канал для стакана
	obCh              *chan *types.Orderbook
	state             strategies.StrategyState[State]
	nextOrderCooldown *time.Timer
	isBuying          isWorking
	isSelling         isWorking

	toPlaceOrders chan *types.PlaceOrder
}

func New() *SpreadStrategy {
	inst := &SpreadStrategy{}
	inst.toPlaceOrders = make(chan *types.PlaceOrder)
	return inst
}

var l *log.Entry

func (s *SpreadStrategy) Start(config *strategies.Config, ordersToPlaceCh *chan *types.PlaceOrder, orderStateChangeCh *chan types.OrderExecutionState) (bool, error) {
	// TODO: Нужен метод ConvertSerialsableToType[T](candidate) T, который конвертирует типы через json.Marshall
	debugCfg := Config{

		Config: strategies.Config{
			// InstrumentId: "BBG004730N88", // SBER
			InstrumentId: "4c466956-d2ce-4a95-abb4-17947a65f18a", // TGLD
			// InstrumentId: "BBG004730RP0", // GAZP
			// InstrumentId: "BBG004PYF2N3", // POLY
			// InstrumentId: "ba64a3c7-dd1d-4f19-8758-94aac17d971b", // FIXP
			// InstrumentId: "BBG004730ZJ9", // VTBR
			Balance: 400,
		},
		maxSharesToHold:     1,
		nextOrderCooldownMs: 0,
		lotSize:             1,
		minProfit:           0,
		stopLossAfter:       0.02,
		// VTBR
		// lotSize: 10_000,
		// minProfit: 0.00002,
		// stopLossAfter: 0.00002,
	}

	l = log.WithFields(log.Fields{
		"strategy":   "spread",
		"instrument": debugCfg.InstrumentId,
	})

	s.config = debugCfg //((any)(*config)).(Config)
	l.Infof("Starting strategy with config: %v", s.config)

	// Создаем или получаем канал, в который будет постаупать инфа о стакане
	l.Tracef("Getting orderbook channel")
	obProvider := orderbook.NewProvider()
	ch, err := obProvider.GetOrCreate(s.config.InstrumentId)
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

	s.obCh = ch

	go func(ch *chan *types.Orderbook) {
		l.Info("Start listening changes in orderbook")
		for {
			select {
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
	return false, nil
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

	isHoldingMaxShares := s.state.Get().holdingShares+s.state.Get().pendingBuyShares >= s.config.maxSharesToHold
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
	if leftBalance < (minBuyPrice * float32(s.config.lotSize)) {
		l.WithField("state", s.state.Get()).Tracef("Not enough money")
		return
	}

	canBuySharesAmount := leftBalance / (minBuyPrice * float32(s.config.lotSize))
	l.Tracef("First bid price: %v; Left money: %v; Can buy %v shares\n", minBuyPrice, leftBalance, canBuySharesAmount)
	if canBuySharesAmount == 0 {
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
	if canBuySharesAmount > float32(s.config.maxSharesToHold) {
		l.Tracef("Can buy more shares, than config allows")
		canBuySharesAmount = float32(s.config.maxSharesToHold)
	}

	order := &types.PlaceOrder{
		InstrumentID: s.config.InstrumentId,
		Quantity:     int64(canBuySharesAmount),
		Price:        types.Price(minBuyPrice),
		Direction:    types.Buy,
	}
	l.Infof("Order to place: %v", order)

	s.toPlaceOrders <- order

	l.Trace("Updating state")
	newState := *s.state.Get()
	newState.pendingBuyShares += int32(canBuySharesAmount)
	newState.notConfirmedBlockedMoney += canBuySharesAmount * minBuyPrice
	newState.lastBuyPrice = minBuyPrice
	s.state.Set(newState)
	l.WithField("state", s.state.Get()).Trace("State updated after place buy order")

	s.isBuying.value = false
	l.Trace("Is buy released")
}

func (s *SpreadStrategy) sell(wg *sync.WaitGroup, ob *types.Orderbook) {
	defer wg.Done()
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
	if state.holdingShares < 0 {
		l.WithField("state", state).Warn("Holding less than 0 shares")
		return
	}

	minAskPrice := ob.Asks[0].Price
	l.Tracef("Min ask price %v", minAskPrice)

	isGoodPrice := minAskPrice-state.lastBuyPrice >= s.config.minProfit
	hasStopLossBroken := s.config.stopLossAfter != 0 && ob.Bids[0].Price <= s.state.Get().lastBuyPrice-s.config.stopLossAfter

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
			"stopLoss":      s.config.stopLossAfter,
			"stopLossPrice": state.lastBuyPrice - s.config.stopLossAfter,
		}).Info("Stop loss broken")
	}

	order := &types.PlaceOrder{
		InstrumentID: s.config.InstrumentId,
		Quantity:     int64(state.holdingShares),
		Direction:    types.Sell,
		Price:        types.Price(price),
	}
	l.Infof("Order to place: %v", order)

	s.toPlaceOrders <- order

	l.Trace("Updating state")
	state = *s.state.Get()
	state.pendingSellShares += state.holdingShares
	s.state.Set(state)
	l.WithField("state", s.state.Get()).Trace("State updated after place sell order")

	s.isSelling.value = false
	l.Trace("Is sell released")
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
	// TODO: Обновлять последнюю цену покупки
	// TODO: Обновлять Оставшийся баланс и остальной стейт

	l.Info("Order state changed %v", state)

	newState := s.state.Get()

	if state.Direction == types.Buy {
		l.Trace("Updating state after buy order executed")
		newState.holdingShares += int32(state.LotsExecuted / int(s.config.lotSize))
		newState.pendingBuyShares -= int32(state.LotsExecuted / int(s.config.lotSize))
		newState.notConfirmedBlockedMoney -= float32(state.ExecutedOrderPrice)
		newState.leftBalance -= float32(state.ExecutedOrderPrice)
		l.Tracef(
			"Lots executed %v of %v; Executed buy price %v",
			state.ExecutedOrderPrice,
			state.LotsExecuted,
			state.LotsRequested,
		)
	}
	if state.Direction == types.Sell {
		l.Trace("Updating state after sell order executed")

		newState.pendingSellShares -= int32(state.LotsExecuted / int(s.config.lotSize))
		newState.leftBalance += float32(state.ExecutedOrderPrice)
		newState.holdingShares += int32(state.LotsExecuted / int(s.config.lotSize))
		l.Tracef(
			"Lots executed %v of %v; Executed sell price %v",
			state.ExecutedOrderPrice,
			state.LotsExecuted,
			state.LotsRequested,
		)
	}
	l.WithField("state", s.state.Get()).Info("State updated")

	s.state.Set(*newState)
}
