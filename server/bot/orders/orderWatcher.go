package orders

import (
	"errors"
	"main/types"
	"sync"
)

type IOrderWatcher interface {
	Watch(string) error
	Register(*chan OrderExecutionState) error
}

type OrderWatcher struct {
	sync.RWMutex
	IOrderWatcher
	idempodentsToWatch     []types.IdempodentId
	idempodentsToOrdersMap map[types.IdempodentId]types.OrderId
	notifyCh               *chan OrderExecutionState
}

var onceOw sync.Once
var ow *OrderWatcher

func NewOrderWatcher() *OrderWatcher {
	if ow != nil {
		return ow
	}

	onceOw.Do(func() {
		ow = &OrderWatcher{
			idempodentsToWatch:     []types.IdempodentId{},
			idempodentsToOrdersMap: make(map[types.IdempodentId]types.OrderId),
		}
	})

	// TODO: Вызывать Broker.SubscribeOrders в горутине, как коллбек дергать ow.Notify 

	return ow
}

func (ow *OrderWatcher) Register(notifyCh *chan OrderExecutionState) error {
	if ow.notifyCh != nil {
		return errors.New("notification channel already set")
	}

	ow.Lock()
	ow.notifyCh = notifyCh
	ow.Unlock()

	return nil
}

func (ow *OrderWatcher) Watch(idempodentId types.IdempodentId) error {

	for _, candidate := range ow.idempodentsToWatch {
		if candidate == idempodentId {
			return nil
		}
	}

	ow.RWMutex.Lock()
	ow.idempodentsToWatch = append(ow.idempodentsToWatch, idempodentId)
	ow.RWMutex.Unlock()

	return nil
}

func (ow *OrderWatcher) PairWithOrderId(idempodentId types.IdempodentId, orderId types.OrderId) error {
	_, ok := ow.idempodentsToOrdersMap[idempodentId]
	if ok {
		return errors.New("already matched with this idempodent")
	}

	ow.RWMutex.Lock()
	ow.idempodentsToOrdersMap[idempodentId] = orderId
	ow.RWMutex.Unlock()

	return nil
}

func (ow *OrderWatcher) Notify(idempodentId string, state OrderExecutionState) error {
	_, ok := ow.idempodentsToOrdersMap[state.IdempodentId]
	if ok {
		return errors.New("found no orders with this idempodent or idempodent not watching")
	}

	*ow.notifyCh <- state

	return nil
}
