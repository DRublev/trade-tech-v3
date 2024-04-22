package bot

import (
	"encoding/json"
	"fmt"
	"main/bot/strategies"
	"main/db"
)

var dbInstance = db.DB{}

// ConfigRepository Репозиторий для доступа к конфигам стратегии
type ConfigRepository struct {
	storage map[string]strategies.Config
}

var instance *ConfigRepository

func getDebugConfig() *strategies.Config {
	/*
		"BBG004730N88" // SBER
		"BBG004730RP0" // GAZP
		"BBG004730ZJ9" // VTBR
		"BBG004PYF2N3" // POLY
		"b71bd174-c72c-41b0-a66f-5f9073e0d1f5" // VKCO
		"4c466956-d2ce-4a95-abb4-17947a65f18a" // TGLD
		"ba64a3c7-dd1d-4f19-8758-94aac17d971b" // FIXP
		"4163e41d-55f4-4f93-82fc-6c44fe5d444e" // SPBE
		"BBG00F9XX7H4" // RNFT
		"9654c2dd-6993-427e-80fa-04e80a1cf4da" // TMOS
	*/
	var debugCfg strategies.Config = make(strategies.Config)
	debugCfg["InstrumentID"] = "4c466956-d2ce-4a95-abb4-17947a65f18a" // TGLD
	debugCfg["Balance"] = 1000
	debugCfg["MaxSharesToHold"] = 100
	debugCfg["NextOrderCooldownMs"] = 0
	debugCfg["LotSize"] = 1
	debugCfg["MinProfit"] = 0
	debugCfg["StopLossAfter"] = 0.01
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

		config, _ := dbInstance.Get([]string{"strategyConf"})
		if config != nil {
			var storage map[string]strategies.Config
			json.Unmarshal(config, &storage)
			instance.storage = storage
		}

		debugConfig := getDebugConfig()
		if debugConfig != nil {
			instance.storage[fmt.Sprintf("spread_v0_%v", (*debugConfig)["InstrumentID"])] = *debugConfig
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
		return nil, err
	}

	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Set Сохранить конфиг
func (cr *ConfigRepository) Set(key string, config strategies.Config) error {
	cr.storage[key] = config

	configJson, _ := json.Marshal(cr.storage)
	dbInstance.Prune([]string{"strategyConf"})
	dbInstance.Append([]string{"strategyConf"}, configJson)

	return nil
}
