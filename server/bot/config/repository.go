package bot

import (
	"errors"
	"main/bot/strategies"
)

type IConfigRepository interface {
}

type ConfigRepository struct {
	IConfigRepository
}

func (cr *ConfigRepository) Get(key string) (*strategies.Config, error) {
	return nil, errors.New("method is not implemented")
}
