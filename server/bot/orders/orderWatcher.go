package orders

import (
	"errors"
	"main/bot/broker"
	"main/types"
	"sync"

	log "github.com/sirupsen/logrus"
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
		log.Info("Creating new order watcher")
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
	log.Infof("Starting to watch order with idempodent id %v", idempodentID)

	for _, candidate := range ow.idempodentsToWatch {
		if candidate == idempodentID {
			log.Warnf("IdempodentID %v is already watching", idempodentID)
			return nil
		}
	}

	ow.RWMutex.Lock()
	ow.idempodentsToWatch = append(ow.idempodentsToWatch, idempodentID)
	ow.RWMutex.Unlock()

	log.Infof("Added %v to watch list", idempodentID)
	return nil
}

// PairWithOrderID Матчит ордер с idempodentID на orderID
func (ow *OrderWatcher) PairWithOrderID(idempodentID types.IdempodentID, orderID types.OrderID) error {
	l := log.WithFields(log.Fields{
		"idempodentID": idempodentID,
		"orderID":      orderID,
	})
	l.Infof("Pairing idempodent id %v with id %v", idempodentID, orderID)

	id, ok := ow.idempodentsToOrdersMap[orderID]
	if ok {
		l.Warnf("IDs %v and %v are already matched", id, orderID)
		return errors.New("already matched with this idempodent")
	}

	ow.RWMutex.Lock()
	ow.idempodentsToOrdersMap[orderID] = idempodentID
	ow.RWMutex.Unlock()

	go func() {
		l.Infof("Getting initial state of order")
		s, err := broker.Broker.GetOrderState(orderID)
		if err != nil {
			l.Error("Failed getting initial state of order: %v", err)
			return
		}

		if s.Status != types.New && s.Status != types.Unspecified {
			l.Info("State of order", s)
			ow.notify(s)
		} else {
			l.Info("Initial state is NEW, no need to notify")
		}
	}()

	l.Info("Paired IDs")
	return nil
}

func (ow *OrderWatcher) notify(state types.OrderExecutionState) {
	idempodentID, ok := ow.idempodentsToOrdersMap[state.ID]
	l := log.WithFields(log.Fields{
		"idempodentID": idempodentID,
		"orderID":      state.ID,
		"instrumentID": state.InstrumentID,
		"direction":    state.Direction,
	})
	if !ok {
		l.Error("Order is not watched, or no one subscribed for it")
		return
	}

	if state.Status == types.Fill {
		l.Info("Order is fullfilled, unsubscribing")
		delete(ow.idempodentsToOrdersMap, state.ID)
	}

	l.Info("Notifying about new order state")
	*ow.notifyCh <- state
}

func (ow *OrderWatcher) ErrNotify(order types.PlaceOrder) {
	orderErr := &types.OrderExecutionState{
		LotsExecuted:       int(order.Quantity),
		ExecutedOrderPrice: float64(order.Price),
		InstrumentID:       order.InstrumentID,
		Direction:          order.Direction,
		Status:             types.ErrorPlacing,
	}

	*ow.notifyCh <- *orderErr
}
