package strategies

import "sync"

// IStrategyState Интерфейс для типа состояния стратегии
// T - тип конкретного состояния, специфичен для стратегии
type IStrategyState[T any] interface {
	Get() *T
	Set(state T) error
	Persist() error
	Restore() error
}

// StrategyState Контракт для состояния стратегии
// T - тип конкретного состояния, специфичен для стратегии
type StrategyState[T any] struct {
	IStrategyState[T]
	sync.RWMutex
	value T
}

func (s *StrategyState[T]) String() string {
	// TODO: Дергать Marshall или выводить читаемый лог стейта
	return "s"
}

func (s *StrategyState[T]) Get() *T {
	s.RLock()
	defer s.RUnlock()
	return &s.value
}

func (s *StrategyState[T]) Set(state T) error {
	s.Lock()
	defer s.Unlock()
	s.value = state
	return nil
}
