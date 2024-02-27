package bot

import (
	"errors"
	"main/bot/strategies"
)

func Assemble(key strategies.StrategyKey, config *strategies.Config) (*strategies.Strategy, error) {
	// TODO: Инициализировать стратегию в зависимости от ключа
	return nil, errors.New("method not implemented")
}
