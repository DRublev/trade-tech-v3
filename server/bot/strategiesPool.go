package bot

import (
	"errors"
	"fmt"
	config "main/bot/config"
	"main/bot/orders"
	"main/bot/strategies"
	"main/types"
	"sync"
)

type IStrategyPool interface {
	Start(key strategies.StrategyKey, instrumentId string) (bool, error)
	Stop(key strategies.StrategyKey, instrumentId string) (bool, error)
}

type StrategiesMap struct {
	sync.RWMutex
	value map[string]*strategies.Strategy
}

type StrategyPool struct {
	IStrategyPool
	configRepository *config.ConfigRepository
	strategies       StrategiesMap
}

var oncePool sync.Once
var pool *StrategyPool

func NewPool() *StrategyPool {
	if pool != nil {
		return pool
	}

	oncePool.Do(func() {
		pool = &StrategyPool{}
		pool.configRepository = &config.ConfigRepository{}
		pool.strategies = StrategiesMap{
			value: make(map[string]*strategies.Strategy),
		}
	})

	// TODO: Подписаться на os.Exit и вызвать Stop для каждой стратегии

	return pool
}

func (sp *StrategyPool) Start(key strategies.StrategyKey, instrumentId string) (bool, error) {
	if !key.IsValid() {
		return false, errors.New("unknown strategy key")
	}

	config, err := sp.getConfig(key, instrumentId)
	if err != nil {
		return false, errors.New("no config found for " + string(key) + " " + instrumentId)
	}

	// TODO: Перенести StrategyMap в отдельный файл с геттерами\сеттерами и контролить мьютекс там
	sp.strategies.RLock()
	_, exists := sp.strategies.value[sp.getMapKey(key, instrumentId)]
	sp.strategies.RUnlock()
	if exists {
		return false, errors.New("strategy already exists")
	}

	strategy, err := Assemble(key, config)
	if err != nil {
		return false, err
	}

	sp.strategies.Lock()
	sp.strategies.value[sp.getMapKey(key, instrumentId)] = strategy
	sp.strategies.Unlock()

	ordersToPlaceCh := make(chan *types.PlaceOrder)
	ordersStateCh := make(chan orders.OrderExecutionState)

	okCh := make(chan bool, 1)
	go func(ordersToPlaceCh chan *types.PlaceOrder, ordersStateCh *chan orders.OrderExecutionState) {
		ok, err := strategy.Start(config, &ordersToPlaceCh, ordersStateCh)
		if err != nil {
			fmt.Println("Error starting strategy ", err)
		}
		okCh <- ok
	}(ordersToPlaceCh, &ordersStateCh)

	ow := orders.NewOrderWatcher()
	go func(source chan *types.PlaceOrder, ordersStateCh *chan orders.OrderExecutionState) {
		err := ow.Register(ordersStateCh)
		if err != nil {
			fmt.Println("error registering notification channel!", err)
			return
		}
		// TODO: Тут будет WithIdempodentId
		for {
			select {
			case order, ok := <- ordersToPlaceCh:
				if !ok {
					fmt.Println("orders to place channel closed")
					return
				}
				orderId, err := Broker.PlaceOrder(order)
				if err != nil {
					fmt.Printf("error placing order: %v\n", err)
					continue
				}
				ow.Watch(order.IdempodentId)
				ow.PairWithOrderId(order.IdempodentId, orderId)
			}
		}
	}(ordersToPlaceCh, &ordersStateCh)

	return <-okCh, nil
}

func (sp *StrategyPool) Stop(key strategies.StrategyKey, instrumentId string) (bool, error) {
	if !key.IsValid() {
		return false, errors.New("unknown strategy key")
	}

	mapKey := sp.getMapKey(key, instrumentId)

	sp.strategies.RLock()
	strategy, exists := sp.strategies.value[mapKey]
	sp.strategies.RUnlock()
	if exists {
		return false, errors.New("strategy not exists")
	}

	ok, err := strategy.Stop()

	return ok, err
}

func (sp *StrategyPool) getConfig(key strategies.StrategyKey, instrumentId string) (*strategies.Config, error) {
	config, err := sp.configRepository.Get(sp.getMapKey(key, instrumentId))
	return config, err
}

func (sp *StrategyPool) getMapKey(key strategies.StrategyKey, instrumentId string) string {
	return string(key) + instrumentId
}
