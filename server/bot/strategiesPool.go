package bot

import (
	"errors"
	"fmt"
	"main/bot/broker"
	config "main/bot/config"
	"main/bot/orders"
	"main/bot/strategies"
	"main/types"
	"sync"

	"github.com/google/uuid"
)

// IStrategyPool Интерфейс пула стратегий
type IStrategyPool interface {
	Start(key strategies.StrategyKey, instrumentID string) (bool, error)
	Stop(key strategies.StrategyKey, instrumentID string) (bool, error)
}

// StrategiesMap Хранилище запущенных стратегий
type StrategiesMap struct {
	sync.RWMutex
	value map[string]strategies.IStrategy
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
		pool.configRepository = &config.ConfigRepository{}
		pool.strategies = StrategiesMap{
			value: make(map[string]strategies.IStrategy),
		}
	})

	// TODO: Подписаться на os.Exit и вызвать Stop для каждой стратегии

	return pool
}

// Start Запуск стратегии
func (sp *StrategyPool) Start(key strategies.StrategyKey, instrumentID string) (bool, error) {
	if !key.IsValid() {
		return false, errors.New("unknown strategy key")
	}

	config, err := sp.getConfig(key, instrumentID)
	if err != nil {
		return false, errors.New("no config found for " + string(key) + " " + instrumentID)
	}

	// TODO: Перенести StrategyMap в отдельный файл с геттерами\сеттерами и контролить мьютекс там
	sp.strategies.RLock()
	_, exists := sp.strategies.value[sp.getMapKey(key, instrumentID)]
	sp.strategies.RUnlock()
	if exists {
		return false, errors.New("strategy already exists")
	}

	strategy, err := Assemble(key, config)
	if err != nil {
		return false, err
	}

	sp.strategies.Lock()
	sp.strategies.value[sp.getMapKey(key, instrumentID)] = strategy
	sp.strategies.Unlock()

	ordersToPlaceCh := make(chan *types.PlaceOrder)
	ordersStateCh := make(chan types.OrderExecutionState)

	okCh := make(chan bool, 1)
	go func(s strategies.IStrategy, ordersToPlaceCh chan *types.PlaceOrder, ordersStateCh *chan types.OrderExecutionState) {
		fmt.Printf("84 strategiesPool %v\n", s)
		ok, err := s.Start(config, &ordersToPlaceCh, ordersStateCh)
		if err != nil {
			fmt.Println("Error starting strategy ", err)
		}
		okCh <- ok
	}(strategy, ordersToPlaceCh, &ordersStateCh)

	go func(source chan *types.PlaceOrder, ordersStateCh *chan types.OrderExecutionState) {
		ow := orders.NewOrderWatcher(ordersStateCh)

		if err != nil {
			fmt.Println("error registering notification channel!", err)
			return
		}
		for {
			select {
			case order, ok := <-ordersToPlaceCh:
				if !ok {
					fmt.Println("orders to place channel closed")
					return
				}

				// TODO: Тут сделать WithIdempodentId
				order.IdempodentID = types.IdempodentID(uuid.New().String())

				orderID, err := broker.Broker.PlaceOrder(order)
				if err != nil {
					fmt.Printf("error placing order: %v\n", err)
					continue
				}
				ow.Watch(order.IdempodentID)
				ow.PairWithOrderID(order.IdempodentID, orderID)
			}
		}
	}(ordersToPlaceCh, &ordersStateCh)

	return <-okCh, nil
}

// Stop Остановить работу стратегии
func (sp *StrategyPool) Stop(key strategies.StrategyKey, instrumentID string) (bool, error) {
	if !key.IsValid() {
		return false, errors.New("unknown strategy key")
	}

	mapKey := sp.getMapKey(key, instrumentID)

	sp.strategies.RLock()
	strategy, exists := sp.strategies.value[mapKey]
	sp.strategies.RUnlock()
	if exists {
		return false, errors.New("strategy not exists")
	}

	ok, err := strategy.Stop()

	return ok, err
}

func (sp *StrategyPool) getConfig(key strategies.StrategyKey, instrumentID string) (*strategies.Config, error) {
	config, err := sp.configRepository.Get(sp.getMapKey(key, instrumentID))
	return config, err
}

func (sp *StrategyPool) getMapKey(key strategies.StrategyKey, instrumentID string) string {
	return string(key) + instrumentID
}
