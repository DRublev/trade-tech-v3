package strategies

import (
	"context"
	"encoding/json"
	"main/types"
)

// StrategyKey Идентификатор стратегии
type StrategyKey string

// Коллекция доступных стратегий
const (
	Spread   StrategyKey = "spread_v0"
	Macd     StrategyKey = "macd"
	RossHook StrategyKey = "rosshook"
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
type Strategy[T any] struct {
	IStrategy
	Key    StrategyKey
	Config T
}

func (s *Strategy[T]) SetConfig(config Config) error {
	bts, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bts, s.Config)
	if err != nil {
		return err
	}

	return nil
}

func (s *Strategy[T]) OnOrderSateChangeSubscribe(stopCtx context.Context, orderStateChangeCh *chan types.OrderExecutionState, onOrderSateChange func(state types.OrderExecutionState)) {
	l.Info("Start listening for orders")
	for {
		select {
		case <-stopCtx.Done():
			l.Info("Strategy stopped")
			return
		case state, ok := <-*orderStateChangeCh:
			if !ok {
				l.Warn("Orders state channel closed")
				return
			}
			onOrderSateChange(state)
		}
	}
}
