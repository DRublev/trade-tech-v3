package orders

import (
	"errors"
	"fmt"
	"main/bot/broker"
	"main/types"
	"sync"
)

// IOrderWatcher Интерфейс для pub-sub на ордеры
type IOrderWatcher interface {
	Watch(string) error
	Register(*chan types.OrderExecutionState) error
	PairWithOrderID(types.IdempodentID, types.OrderID) error
}

// OrderWatcher Провайдит функционал подписки на ордера и их состояние
type OrderWatcher struct {
	sync.RWMutex
	IOrderWatcher
	idempodentsToWatch     []types.IdempodentID
	idempodentsToOrdersMap map[types.OrderID]types.IdempodentID
	notifyCh               *chan types.OrderExecutionState
}

var onceOw sync.Once
var ow *OrderWatcher

// NewOrderWatcher Конструктор для OrderWatcher
func NewOrderWatcher(notifyCh *chan types.OrderExecutionState) *OrderWatcher {
	if ow != nil {
		return ow
	}

	onceOw.Do(func() {
		ow = &OrderWatcher{
			idempodentsToWatch:     []types.IdempodentID{},
			idempodentsToOrdersMap: make(map[types.OrderID]types.IdempodentID),
			notifyCh:               notifyCh,
		}
	})

	go broker.Broker.SubscribeOrders(ow.notify)

	return ow
}

// Watch Подписаться на ордер по idempodentID
// TODO: Переименовать в Subscribe
func (ow *OrderWatcher) Watch(idempodentID types.IdempodentID) error {

	for _, candidate := range ow.idempodentsToWatch {
		if candidate == idempodentID {
			fmt.Println("50 orderWatcher ", "already watching this id")
			return nil
		}
	}

	ow.RWMutex.Lock()
	ow.idempodentsToWatch = append(ow.idempodentsToWatch, idempodentID)
	ow.RWMutex.Unlock()
	fmt.Printf("58 orderWatcher watching idempodent: %v\n", idempodentID)
	return nil
}

// PairWithOrderID Матчит ордер с idempodentID на orderID
func (ow *OrderWatcher) PairWithOrderID(idempodentID types.IdempodentID, orderID types.OrderID) error {
	id, ok := ow.idempodentsToOrdersMap[orderID]
	if ok {
		fmt.Println("65 orderWatcher ", id)
		return errors.New("already matched with this idempodent")
	}

	ow.RWMutex.Lock()
	ow.idempodentsToOrdersMap[orderID] = idempodentID
	ow.RWMutex.Unlock()

	go func() {
		fmt.Println("74 orderWatcher ", "getting order state")
		s, err := broker.Broker.GetOrderState(orderID)
		if err!= nil {
			fmt.Printf("77 orderWatcher %v\n", err)
			return
		}
		fmt.Printf("80 orderWatcher %v\n", s)
		if s.Status != types.New && s.Status != types.Unspecified {
			ow.notify(s)
		}
	}()

	fmt.Println("72 orderWatcher ", "paired")
	return nil
}

func (ow *OrderWatcher) notify(state types.OrderExecutionState) {
	idempodentID, ok := ow.idempodentsToOrdersMap[state.ID]
	if !ok {
		fmt.Printf("75 orderWatcher %v\n", ow.idempodentsToOrdersMap)
		fmt.Printf("found no orders with this id or id not watching id %v; idempodent: %v\n", state.ID, state.IdempodentID)
		return
	}

	fmt.Printf("state for order %v changed: %v\n", idempodentID, state)
	if state.Status == types.Fill {
		delete(ow.idempodentsToOrdersMap, state.ID)
	}

	*ow.notifyCh <- state
}
