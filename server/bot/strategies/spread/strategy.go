package spread

import (
	"fmt"
	"main/bot/orderbook"
	"main/bot/strategies"
	"main/types"
	"sync"
	"time"
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

func (s *SpreadStrategy) Start(config *strategies.Config, ordersToPlaceCh *chan *types.PlaceOrder, orderStateChangeCh *chan types.OrderExecutionState) (bool, error) {
	// TODO: Нужен метод ConvertSerialsableToType[T](candidate) T, который конвертирует типы через json.Marshall
	debugCfg := Config{
		
		Config: strategies.Config{
			// InstrumentId: "BBG004730N88", // SBER
			// InstrumentId: "4c466956-d2ce-4a95-abb4-17947a65f18a", // TGLD
			// InstrumentId: "BBG004730RP0", // GAZP
			// InstrumentId: "BBG004PYF2N3", // POLY
			InstrumentId: "ba64a3c7-dd1d-4f19-8758-94aac17d971b", // FIXP
			Balance: 400,
		},
		maxSharesToHold: 1,
		nextOrderCooldownMs: 0,
		lotSize: 1,
		minProfit: 0.1,
		stopLossAfter: 0,
	}
	s.config = debugCfg //((any)(*config)).(Config)
	fmt.Printf("78 strategy %v\n", s.config)

	// Создаем или получаем канал, в который будет постаупать инфа о стакане
	obProvider := orderbook.NewOrderbookProvider()
	fmt.Printf("95 strategy %v\n", obProvider)
	ch, err := obProvider.GetOrCreate(s.config.InstrumentId)
	if err != nil {
		fmt.Printf("98 strategy %v\n", err)
		return false, err
	}

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
		fmt.Printf("113 strategy %v\n", err)
		return false, err
	}

	s.nextOrderCooldown = time.NewTimer(time.Duration(0) * time.Millisecond)

	s.obCh = ch
	// Слушаем изменения в стакане
	go func(ch *chan *types.Orderbook) {
		for {
			select {
			case ob, ok := <-*ch:
				if !ok {
					fmt.Println("spread orderbook channel end")
					return
				}

				go s.onOrderbook(ob)
			}
		}
	}(ch)

	// Копируем выставляемые ордера в другой канал
	go func(source *chan *types.PlaceOrder, target *chan *types.PlaceOrder) {
		for {
			select {
			case orderToPlace, ok := <-*source:
				if !ok {
					fmt.Println("orders to place channel end, closing target channel")
					return
				}
				*target <- orderToPlace
			}
		}
	}(&s.toPlaceOrders, ordersToPlaceCh)

	// Подписка на изменения в ордерах
	go func(source *chan types.OrderExecutionState) {
		for {
			select {
			case state, ok := <-*source:
				if !ok {
					fmt.Println("order state channel end")
					return
				}
				go s.onOrderSateChange(state)
			}

		}
	}(orderStateChangeCh)

	return true, nil
}

func (s *SpreadStrategy) Stop() (bool, error) {
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

	if s.isBuying.value {
		fmt.Println("Already buying")
		return
	}

	isHoldingMaxShares := s.state.Get().holdingShares+s.state.Get().pendingBuyShares >= s.config.maxSharesToHold
	if isHoldingMaxShares {
		fmt.Printf("already holding max shares. holding %v, processing %v,  max %v\n", s.state.Get().holdingShares, s.state.Get().pendingBuyShares, s.config.maxSharesToHold)
		return
	}

	// Аукцион закрытия, только заявки на продажу
	if len(ob.Bids) == 0 {
		return
	}

	minBuyPrice := ob.Bids[0].Price
	leftBalance := s.state.Get().leftBalance - s.state.Get().notConfirmedBlockedMoney
	if leftBalance < minBuyPrice {
		fmt.Printf("Not enough money to enter position. First bid price: %v; Left money: %v\n", minBuyPrice, leftBalance)
		return
	}

	canBuySharesAmount := leftBalance / (minBuyPrice * float32(s.config.lotSize))
	fmt.Printf("First bid price: %v; Left money: %v; Can buy %v shares\n", minBuyPrice, leftBalance, canBuySharesAmount)
	if canBuySharesAmount == 0 {
		fmt.Println("Can buy 0 shares")
		return
	}

	ok := s.isBuying.TryLock()
	if !ok {
		fmt.Println("isBuy mutex cannot be locked")
		return
	}
	defer s.isBuying.Unlock()

	s.isBuying.value = true

	if canBuySharesAmount > float32(s.config.maxSharesToHold) {
		canBuySharesAmount = float32(s.config.maxSharesToHold)
	}

	order := &types.PlaceOrder{
		InstrumentID: s.config.InstrumentId,
		Quantity: int64(canBuySharesAmount),
		Price: types.Price(minBuyPrice),
		Direction: types.Buy,
	}
	fmt.Printf("Order to place: %v\n", order)

	s.toPlaceOrders <- order

	newState := *s.state.Get()
	newState.pendingBuyShares += int32(canBuySharesAmount)
	newState.notConfirmedBlockedMoney += canBuySharesAmount * minBuyPrice
	newState.lastBuyPrice = minBuyPrice
	s.state.Set(newState)
	s.isBuying.value = false
}

func (s *SpreadStrategy) sell(wg *sync.WaitGroup, ob *types.Orderbook) {
	defer wg.Done()

	if s.isSelling.value {
		fmt.Println("Already buying")
		return
	}

	state := *s.state.Get()
	if state.holdingShares-state.pendingSellShares == 0 {
		fmt.Println("Nothing to sell, hold 0 shares")
		return
	}
	if state.holdingShares < 0 {
		fmt.Printf("ERROR, holding less than 0 shares: %v\n", state.holdingShares)
		return
	}

	minAskPrice := ob.Asks[0].Price
	isGoodPrice := minAskPrice-s.state.Get().lastBuyPrice >= s.config.minProfit
	hasStopLossBroken := s.config.stopLossAfter != 0 && ob.Bids[0].Price <= s.state.Get().lastBuyPrice - s.config.stopLossAfter
	shouldMakeSell := isGoodPrice || hasStopLossBroken
	if !shouldMakeSell {
		fmt.Printf("Not a good deal. asks[0].price: %v; lastBuyPrice: %v; minProfit: %v\n", minAskPrice, s.state.Get().lastBuyPrice, s.config.minProfit)
		return
	}

	ok := s.isSelling.TryLock()
	if !ok {
		fmt.Println("isSelling mutex cannot be locked")
		return
	}
	defer s.isSelling.Unlock()

	price := minAskPrice
	if hasStopLossBroken {
		price = ob.Bids[0].Price
		fmt.Printf("Stop loss broken. stop loss: %v; current buy price: %v\n", s.state.Get().lastBuyPrice - s.config.stopLossAfter, price)
	}
	

	order := &types.PlaceOrder{
		InstrumentID: s.config.InstrumentId,
		Quantity: int64(state.holdingShares),
		Direction: types.Sell,
		Price: types.Price(price),
	}
	fmt.Printf("Order to place: %v\n", order)
	s.toPlaceOrders <- order

	state = *s.state.Get()
	state.pendingSellShares += state.holdingShares
	s.state.Set(state)
	s.isSelling.value = false
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
	fmt.Printf("291 strategy %v %v \n", state, state.ExecutedOrderPrice)
	newState := s.state.Get()
	
	if state.Direction == types.Buy {
		newState.holdingShares += int32(state.LotsExecuted)
		newState.pendingBuyShares -= int32(state.LotsExecuted)
		newState.notConfirmedBlockedMoney -= float32(state.ExecutedOrderPrice)
		newState.leftBalance -= float32(state.ExecutedOrderPrice)
	}
	if state.Direction == types.Sell {
		newState.pendingSellShares -= int32(state.LotsExecuted)
		newState.leftBalance += float32(state.ExecutedOrderPrice)
	}
fmt.Printf("331 strategy %v\n", newState)
	s.state.Set(*newState)
}
