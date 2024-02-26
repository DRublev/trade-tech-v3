package bot

import (
	"errors"
	"main/bot/strategies"
	"sync"
)

type IStrategyPool interface {
	Start(key strategies.StrategyKey, instrumentId string) (bool, error)
	Stop(key strategies.StrategyKey, instrumentId string) (bool, error)
}

type StrategyPool struct {
	IStrategyPool
}

var once sync.Once
var pool *StrategyPool

func NewPool() *StrategyPool {
	if pool != nil {
		return pool
	}

	once.Do(func() {
		pool = &StrategyPool{}
	})

	// TODO: Подписаться на os.Exit и вызвать Stop для каждой стратегии

	return pool
}

func (sp *StrategyPool) Start(key strategies.StrategyKey, instrumentId string) (bool, error) {
	_, err := sp.getConfig(key, instrumentId)
	if err != nil {
		return false, errors.New("no config found for " + string(key) + " " + instrumentId)
	}

	// go func() { strategy, err := factory.Assemble(key, config) }

	return false, errors.New("method is not implemented")
}

func (sp *StrategyPool) Stop(key strategies.StrategyKey, instrumentId string) (bool, error) {
	return false, errors.New("method is not implemented")
}

func (sp *StrategyPool) getConfig(key strategies.StrategyKey, instrumentId string) (*strategies.Config, error) {
	//  TODO: Вызывать ConfigRepository.Get
	return nil, errors.New("method is not implemented")
}
