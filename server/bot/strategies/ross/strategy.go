package ross

import (
	"main/types"
)

type RossStrategy struct {

}

type Extremums struct {
	// Минимум локального нисходящего тренда
	minOfTrend types.OHLC
	
	// Максимум коррекции
	// Стоп-лосс ставим чуть под него
	correctionMax types.OHLC
	// Флаг для подтверждения correctionMax (локальной смены тренда)
	isPotentialCorrectionMax bool

	// Минимумальная точка, которая больше чем minOfTrend, но меньше correctionMax
	aboveMinOfTrend types.OHLC

	// Следующий экстремум, после correctionMax
	enterExtremum types.OHLC
	// Флаг для подтверждения enterExtremum (локальной смены тренда)
	isPotentialEnterExtremum bool
}

func (s *RossStrategy) processOrderbook(ob *types.Orderbook) {
	// find Extremums
	// 1 текущее закрытие сравниваем с minOfTrend
	// 1.1 если меньше, то обновить minOfTrend
	// 1.2 если больше, то ставим isPotentialCorrectionMax
	// 1.2 если больше и isPotentialCorrectionMax, то ставим как correctionMax

	// 2 текущее закрытие сравниваем с correctionMax
	// 2.1 если больше, то обновить correctionMax
	// 2.2 если меньше и , то ставим как aboveMinOfTrend
	
	// 3 текущее закрытие сравниваем с aboveMinOfTrend
	// 3.1 если меньше, то обновить aboveMinOfTrend
	// 3.2 если aboveMinOfTrend меньше minOfTrend, то сбрасываем correctionMax и aboveMinOfTrend и идем на 1.1
	// 3.3 если больше, то сетим enterExtremum

	// 4 текущее закрытие сравниваем с enterExtremum
	// 4.1 если больше, то обновляем enterExtremum
	// 4.2 если меньше, то сетим isPotentialEnterExtremum и ждем следующую свечу, возврааемся к 4
	// 4.3 если меньше и isPotentialEnterExtremum, то идем к 5

	// 5 сравниваем закрытие свечи экстремума и закрытие прошлой свечи
	// 5.1 если закрытие прошлой свечи меньше, то сигнал НЕ сработал
	// 5.2 если закрытие прошлой свечи больше, то сигнал сработал, ставим точку входа НАД enterExtremum

	// Учесть закрытие бай ордера при пробитии стопа
	// Тейк-профит фикс процент (2 к 1 от стопа)
}