package strategies

type IStrategyState[T any] interface {
	Get() *T
	Set(state T) error
	Persist() error
	Restore() error
}

type StrategyState[T any] struct {
	IStrategyState[T]
	state T
}

func (s *StrategyState[T]) Get() *T {
	return &s.state
}
