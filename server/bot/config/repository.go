package bot

import (
	"errors"
	"fmt"
	"main/bot/strategies"
)

type IConfigRepository interface {
}

type ConfigRepository struct {
	IConfigRepository
}

// 	// Минимальная разница bid-ask, при которой выставлять ордер
// 	minSpread float32

// 	// Сколько мс ждать после исполнения итерации покупка-продажа перед следующей
// 	nextOrderCooldownMs int32

// 	// Каким количчеством акций торговать? Макс
// 	maxSharesToHold int32

// 	// Лотность инструмента
// 	lotSize int32

// Доступный для торговли баланс
	// Balance float32

	// // Акция для торговли
	// InstrumentId string

func (cr *ConfigRepository) Get(key string) (*strategies.Config, error) {
	// TODO: Возможно тут понадобится мьютекс
	fmt.Printf("18 repository %v  \n", key)
	_ = map[string]any {
		// InstrumentId: "BBG004730N88", // SBER
		// "InstrumentId": "4c466956-d2ce-4a95-abb4-17947a65f18a", // TGLD
		"InstrumentId": "BBG004730RP0", // GAZP
		"Balance": 200,
		"maxSharesToHold": 1,
		"nextOrderCooldownMs": 0,
		"lotSize": 1,
		"minSpread": 0.2,
	}
	return nil, errors.New("method is not implemented")
}
