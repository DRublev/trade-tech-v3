package bot

import (
	"encoding/json"
	"fmt"
	"main/bot/strategies"
)

// ConfigRepository Репозиторий для доступа к конфигам стратегии
type ConfigRepository struct {
	storage map[string]strategies.Config
}

var instance *ConfigRepository

func getDebugConfig() *strategies.Config {
	/*
		"BBG004730N88"; // SBER
		"BBG004730RP0"; // GAZP
		"BBG004730ZJ9"; // VTBR
		"BBG004PYF2N3"; // POLY
		"4c466956-d2ce-4a95-abb4-17947a65f18a"; // TGLD
		"ba64a3c7-dd1d-4f19-8758-94aac17d971b"; // FIXP
	*/
	var debugCfg strategies.Config = make(strategies.Config)
	debugCfg["InstrumentID"] = "BBG004PYF2N3" // POLY
	debugCfg["Balance"] = 450
	debugCfg["MaxSharesToHold"] = 1
	debugCfg["NextOrderCooldownMs"] = 0
	debugCfg["LotSize"] = 1
	debugCfg["MinProfit"] = 0.34
	debugCfg["StopLossAfter"] = 1
	// VTBR
	// lotSize: 10_000,
	// minProfit: 0.00002,
	// stopLossAfter: 0.00002,

	var res strategies.Config

	debugBytes, err := json.Marshal(debugCfg)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(debugBytes, &res)
	if err != nil {
		return nil
	}

	return &res
}

// New Конструктор
func New() *ConfigRepository {
	if instance == nil {
		instance = &ConfigRepository{
			storage: make(map[string]strategies.Config),
		}

		debugConfig := getDebugConfig()
		if debugConfig != nil {
			instance.storage["spread_v0_BBG004PYF2N3"] = *debugConfig
		}

	}
	return instance
}

// Get Получить конфиг стратегии по ее ключу
func (cr *ConfigRepository) Get(key string) (*strategies.Config, error) {
	// TODO: Нужен метод ConvertSerialsableToType[T](candidate) T, который конвертирует типы через json.Marshall
	// TODO: В будущем  тут понадобится мьютекс

	stored, exists := cr.storage[key]
	if !exists {
		return nil, fmt.Errorf("config with key %s not found", key)
	}

	var res strategies.Config

	b, err := json.Marshal(stored)
	if err != nil {
		fmt.Printf("42 repository %v\n", err)
		return nil, err
	}

	err = json.Unmarshal(b, &res)
	if err != nil {
		fmt.Printf("48 repository %v\n", err)
		return nil, err
	}
	fmt.Printf("52 repository %v\n", res)

	return &res, nil
}

// Set Сохранить конфиг
func (cr *ConfigRepository) Set(key string, config strategies.Config) error {
	cr.storage[key] = config
	return nil
}
