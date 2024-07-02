package bot

import (
	"errors"
	"main/bot/broker"
	config "main/bot/config"
	errs "main/bot/errors"
	"main/bot/orders"
	"main/bot/strategies"
	"main/types"
	"sync"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// IStrategyPool Интерфейс пула стратегий
type IStrategyPool interface {
	Start(key strategies.StrategyKey, instrumentID string) (bool, error)
	Stop(key strategies.StrategyKey, instrumentID string) (bool, error)
}

// StrategyPool Аггрегатор стратегий. Весь доступ к стратегии ведется через него
type StrategyPool struct {
	IStrategyPool
	configRepository *config.ConfigRepository
	strategies       StrategiesMap
}

var oncePool sync.Once
var pool *StrategyPool

// NewPool Хранилище для StrategyPool
func NewPool() *StrategyPool {
	if pool != nil {
		return pool
	}

	oncePool.Do(func() {
		pool = &StrategyPool{}
		pool.configRepository = config.New()
		pool.strategies = StrategiesMap{
			value: make(map[string]strategies.IStrategy),
		}
	})

	// TODO: Подписаться на os.Exit и вызвать Stop для каждой стратегии

	return pool
}

func (sp *StrategyPool) IsStarted(key strategies.StrategyKey, instrumentID string) (bool, error) {
	_, exists := sp.strategies.GetValue(sp.getMapKey(key, instrumentID))
	return exists, nil
}

// Start Запуск стратегии
func (sp *StrategyPool) Start(key strategies.StrategyKey, instrumentID string) (bool, error) {
	l := log.WithFields(log.Fields{
		"method":       "Start",
		"instrument": instrumentID,
		"strategy":     key,
	})

	l.Info("Starting strategy")

	if !key.IsValid() {
		l.Tracef("Unknown strategy key %v; %v", key, instrumentID)
		return false, errs.UnknownStrategy
	}

	config, err := sp.getConfig(key, instrumentID)
	if err != nil {
		l.Tracef("No config found: %v", err)
		return false, errors.New("no config found for " + string(key) + " " + instrumentID)
	}

	_, exists := sp.strategies.GetValue(sp.getMapKey(key, instrumentID))
	if exists {
		l.Trace("Strategy already exists")
		return false, errors.New("strategy already exists")
	}

	strategy, err := Assemble(key, config)
	if err != nil {
		l.Errorf("Error assembling strategy: %v", err)
		return false, err
	}

	sp.strategies.SetValue(sp.getMapKey(key, instrumentID), strategy)

	l.Trace("Creating channels")
	ordersToPlaceCh := make(chan *types.PlaceOrder)
	ordersStateCh := make(chan types.OrderExecutionState)

	okCh := make(chan bool, 1)
	go func() {
		l.Trace("Starting strategy")
		ok, err := strategy.Start(config, &ordersToPlaceCh, &ordersStateCh)
		if err != nil {
			l.Errorf("Error starting strategy in pool %v", err)
		}
		okCh <- ok
	}()

	go func() {
		l.Trace("Registering channel for orders to place")
		ow := orders.NewOrderWatcher(&ordersStateCh)

		for {
			select {
			case order, ok := <-ordersToPlaceCh:
				l.Tracef("New order to place")
				if !ok {
					l.Tracef("Orders to place channel closed; %v", instrumentID)
					return
				}

				// TODO: Тут сделать WithIdempodentId
				order.IdempodentID = types.IdempodentID(uuid.New().String())

				orderID, err := broker.Broker.PlaceOrder(order)
				if err != nil {
					l.Errorf("Error placing order: %v", err)

					// Если ошибка закрытия ордера, не надо уведомлять об этом стратегию, тк у нее неверно обновится стейт
					if err.Error() != "error closing order" || err.Error() == "order cancelled" {
						ow.ErrNotify(*order)
					}
					continue
				}
				l.Trace("Order place processed")

				ow.Watch(order.IdempodentID)
				ow.PairWithOrderID(order.IdempodentID, orderID)
			}
		}
	}()

	l.Info("Strategy started")
	return <-okCh, nil
}

// Stop Остановить работу стратегии
func (sp *StrategyPool) Stop(key strategies.StrategyKey, instrumentID string) (bool, error) {
	l := log.WithFields(log.Fields{
		"method":       "Stop",
		"instrument": instrumentID,
		"strategy":     key,
	})
	l.Info("Stopping strategy")

	if !key.IsValid() {
		l.Tracef("Unknown strategy key %v; %v", key, instrumentID)
		return false, errs.UnknownStrategy
	}

	mapKey := sp.getMapKey(key, instrumentID)

	sp.strategies.RLock()
	strategy, exists := sp.strategies.value[mapKey]
	sp.strategies.RUnlock()
	if !exists {
		l.Error("Strategy doesnt exists")
		return false, errors.New("strategy not exists")
	}

	l.Trace("Trying to call Stop of a strategy instance")
	ok, err := strategy.Stop()
	l.Trace("Called Stop of a strategy instance")

	sp.strategies.RLock()
	delete(sp.strategies.value, mapKey)
	sp.strategies.RUnlock()

	return ok, err
}

func (sp *StrategyPool) getConfig(key strategies.StrategyKey, instrumentID string) (*strategies.Config, error) {
	l := log.WithFields(log.Fields{
		"instrument": instrumentID,
		"strategy":     key,
	})

	configKey := sp.getMapKey(key, instrumentID)
	l.Tracef("Getting config for %v", configKey)
	config, err := sp.configRepository.Get(configKey)
	return config, err
}

func (sp *StrategyPool) getMapKey(key strategies.StrategyKey, instrumentID string) string {
	return string(key) + "_" + instrumentID
}
