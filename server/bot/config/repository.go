package bot

import (
	"main/bot/strategies"
)

// ConfigRepository Репозиторий для доступа к конфигам стратегии
type ConfigRepository struct {
}

// Get Получить конфиг стратегии по ее ключу
func (cr *ConfigRepository) Get(key string) (*strategies.Config, error) {
	// TODO: Возможно тут понадобится мьютекс
	return nil, nil
	// return nil, errors.New("method is not implemented")
}
