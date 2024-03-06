package bot

import (
	"main/bot/strategies"
)

type IConfigRepository interface {
}

type ConfigRepository struct {
	IConfigRepository
}

func (cr *ConfigRepository) Get(key string) (*strategies.Config, error) {
	// TODO: Возможно тут понадобится мьютекс
	return nil, nil
	// return nil, errors.New("method is not implemented")
}
