package bot

import (
	"errors"
	"main/bot/strategies"
	"main/bot/strategies/spread"
)

func Assemble(key strategies.StrategyKey, config *strategies.Config) (*strategies.Strategy, error) {

	switch key {
	case strategies.Spread:
		var s strategies.Strategy
		instance := spread.New()
		// TODO: Разобраться с типам
		s = ((any)(*instance)).(strategies.Strategy)
		return &s, nil
	}

	// TODO: Инициализировать стратегию в зависимости от ключа
	return nil, errors.New("method not implemented")
}
