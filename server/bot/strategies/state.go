package strategies

type IStrategyState[T any] interface {
	Get() *T
	Set(state T) error
	Persist() error
	Restore() error
}

type StrategyState[T any] struct {
	IStrategyState[T]
	value T
}

func (s *StrategyState[T]) String() string {
	// TODO: Дергать Marshall
	return "s"
}

func (s *StrategyState[T]) Get() *T {
	return &s.value
}
