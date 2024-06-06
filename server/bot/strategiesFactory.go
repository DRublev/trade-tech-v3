package bot

import (
	"context"
	"errors"
	"main/bot/broker"
	"main/bot/strategies"
	"main/bot/strategies/macd"
	"main/bot/strategies/spread"
	"main/types"
	"os"
	"os/signal"
)

// Assemble Фабрика для сборки стратегий
func Assemble(key strategies.StrategyKey, config *strategies.Config) (strategies.IStrategy, error) {
	backCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	err := broker.Init(backCtx, types.Tinkoff)
	if err != nil {
		return nil, err
	}

	switch key {
	case strategies.Spread:
		return spread.New(), nil
	case strategies.Macd:
		return macd.New(), nil
	}

	// TODO: Инициализировать стратегию в зависимости от ключа
	return nil, errors.New("Strategy not implemented")
}
