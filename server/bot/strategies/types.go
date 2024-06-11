package strategies

import (
	"main/types"
)

// StrategyKey Идентификатор стратегии
type StrategyKey string

// Коллекция доступных стратегий
const (
	Spread   StrategyKey = "spread_v0"
	Macd     StrategyKey = "macd"
	RossHook StrategyKey = "ross_hook"
)

// IsValid Валидность идентификатора стратегии
func (s StrategyKey) IsValid() bool {
	switch s {
	case Spread:
		return true
	case Macd:
		return true
	case RossHook:
		return true
	}
	return false
}

// Config Общий конфиг для стратегии
type Config map[string]any

// IStrategy Интерфейс стратегии
type IStrategy interface {
	Start(config *Config, ordersToPlaceCh *chan *types.PlaceOrder, ordersStateCh *chan types.OrderExecutionState) (bool, error)
	Stop() (bool, error)
}

// Strategy Контракт для стратегии
type Strategy struct {
	IStrategy
	Key StrategyKey
}
