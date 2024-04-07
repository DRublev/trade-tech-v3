package bot

import (
	"main/bot/strategies"
	"sync"
)

// StrategiesMap Хранилище запущенных стратегий
type StrategiesMap struct {
	sync.RWMutex
	value map[string]strategies.IStrategy
}

func (sm *StrategiesMap) SetValue(key string, strategy strategies.IStrategy) {
	sm.Lock()
	sm.value[key] = strategy
	sm.Unlock()
}

func (sm *StrategiesMap) GetValue(key string) (strategies.IStrategy, bool) {
	defer sm.RUnlock()

	sm.RLock()
	strategy, exist := sm.value[key]
	return strategy, exist
}