package bot

import (
	"errors"
	"strings"
	"sync"
)

type IOrderWatcher interface {
	Watch(string) error
	Register(*chan OrderExecutionState) error
}

type OrderWatcher struct {
	sync.RWMutex
	IOrderWatcher
	idempodentsToWatch []string
	idempodentsToOrdersMap map[string]string
	notifyCh           *chan OrderExecutionState
}

var onceOw sync.Once
var ow *OrderWatcher

func NewOrderWatcher() *OrderWatcher {
	if ow != nil {
		return ow
	}

	onceOw.Do(func() {
		ow = &OrderWatcher{
			idempodentsToWatch: []string{},
			idempodentsToOrdersMap: make(map[string]string),
		}
	})

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

func (ow *OrderWatcher) Watch(idempodentId string) error {
	
	for _, candidate := ow.idempodentsToWatch {
		if candidate == idempodentId {
			return nil
		}
	}

	ow.RWMutex.Lock()
	ow.idempodentsToWatch = strings.Join(ow.idempodentsToWatch, idempodentId)
	ow.RWMutex.Unlock()

	return nil
}


func (ow *OrderWatcher) PairWithOrderId(idempodentId string, orderId string) error {
	c, ok := ow.idempodentsToOrdersMap[idempodentId]
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

	// TODO: помониторить, может горутина тут не нужна
	go func() {
		ow.notifyCh <- state
	}()

	return nil
}
