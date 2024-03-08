package orders

import (
	"errors"
	"fmt"
	"main/bot/broker"
	"main/types"
	"sync"
)

type IOrderWatcher interface {
	Watch(string) error
	Register(*chan types.OrderExecutionState) error
	PairWithOrderId(types.IdempodentId, types.OrderId) error
}

type OrderWatcher struct {
	sync.RWMutex
	IOrderWatcher
	idempodentsToWatch     []types.IdempodentId
	idempodentsToOrdersMap map[types.OrderId]types.IdempodentId
	notifyCh               *chan types.OrderExecutionState
}

var onceOw sync.Once
var ow *OrderWatcher

func NewOrderWatcher(notifyCh *chan types.OrderExecutionState) *OrderWatcher {
	if ow != nil {
		return ow
	}

	onceOw.Do(func() {
		ow = &OrderWatcher{
			idempodentsToWatch:     []types.IdempodentId{},
			idempodentsToOrdersMap: make(map[types.OrderId]types.IdempodentId),
			notifyCh:               notifyCh,
		}
	})

	go broker.Broker.SubscribeOrders(ow.notify)

	return ow
}

func (ow *OrderWatcher) Watch(idempodentId types.IdempodentId) error {

	for _, candidate := range ow.idempodentsToWatch {
		if candidate == idempodentId {
			fmt.Println("50 orderWatcher ", "already watching this id")
			return nil
		}
	}

	ow.RWMutex.Lock()
	ow.idempodentsToWatch = append(ow.idempodentsToWatch, idempodentId)
	ow.RWMutex.Unlock()
	fmt.Printf("58 orderWatcher watching idempodent: %v\n", idempodentId)
	return nil
}

func (ow *OrderWatcher) PairWithOrderId(idempodentId types.IdempodentId, orderId types.OrderId) error {
	id, ok := ow.idempodentsToOrdersMap[orderId]
	if ok {
		fmt.Println("65 orderWatcher ", id)
		return errors.New("already matched with this idempodent")
	}

	ow.RWMutex.Lock()
	ow.idempodentsToOrdersMap[orderId] = idempodentId
	ow.RWMutex.Unlock()
	fmt.Println("72 orderWatcher ", "paired")
	return nil
}

func (ow *OrderWatcher) notify(state types.OrderExecutionState) {
	idempodentID, ok := ow.idempodentsToOrdersMap[state.Id]
	if !ok {
		fmt.Printf("75 orderWatcher %v\n", ow.idempodentsToOrdersMap)
		fmt.Printf("found no orders with this id or id not watching id %v; idempodent: %v\n", state.Id, state.IdempodentId)
		return
	}

	fmt.Printf("state for order %v changed: %v\n", idempodentID, state)

	*ow.notifyCh <- state
}
