package strategies

import "main/types"

type StrategyKey string

const (
	Spread StrategyKey = "spread_v0"
)

func (s StrategyKey) IsValid() bool {
	switch s {
	case Spread:
		return true
	}
	return false
}

type Config interface{}

type IStrategy interface {
	Start(config Config) (bool, error)
	Stop() (bool, error)
}

type Strategy struct {
	Key             StrategyKey
	OrdersToPlaceCh *chan types.Order
}
