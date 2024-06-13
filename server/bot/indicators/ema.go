package indicators

import (
	"errors"
)

// EmaIndicator Экспоненциальная скользящая средняя
type EmaIndicator struct {
	Indicator[float64, float64]
	// Период, за который будем считать среднюю. Условно диапазон массива с конца
	period int

	// Коэффицент усреднения
	k float64

	// Уже обработанные цены
	prevPrices []float64

	// Индекс последней обработанной цены
	// Нужен чтобы не вводить чанки
	latestCalcedPriceIdx int

	// Значния индикаторы
	values []float64

	sma SmaIndicator

	precision uint
}

// NewEma Конструктор
func NewEma(period int) *EmaIndicator {
	inst := &EmaIndicator{}
	inst.period = period
	inst.k = 2 / float64(1+period)
	inst.latestCalcedPriceIdx = -1 // Невалидный индекс, чтобы сетить его в коде в нужный момент
	inst.prevPrices = []float64{}
	inst.values = []float64{}
	inst.sma = *NewSma(period)
	return inst
}

// Latest Получить последнее значение
func (i *EmaIndicator) Latest() (float64, error) {
	if len(i.values) == 0 {
		return 0, errors.New("No latest value")
	}

	return i.values[len(i.values)-1], nil
}

// Get Получить все значения
func (i *EmaIndicator) Get() []float64 {
	return i.values
}

// Update Уточнить значение. Юзать при поступлени новых данных
func (i *EmaIndicator) Update(price float64) {
	i.prevPrices = append(i.prevPrices, price)
	i.sma.Update(price)
	// Количество данных должно быть минимум х2 от периода рассчета
	if len(i.prevPrices) < 2*i.period {
		return
	}
	// i.prevPrices = i.prevPrices[len(i.prevPrices)-2*i.period:]

	var emas []float64

	roundPrecision := detectPrecision(price)
	if roundPrecision > i.precision {
		i.precision = roundPrecision
	}
	// first ema value = sma value
	allSma := i.sma.Get()
	if len(allSma) < i.period {
		return
	}
	sma := allSma[:i.period]
	previousEma := sma[0]

	emas = append(emas, roundFloat(previousEma, i.precision))

	for _, p := range i.prevPrices[i.period:] {
		previousEma = emas[len(emas)-1]
		ema := (p * i.k) + (previousEma * (1 - i.k))
		emas = append(emas, roundFloat(ema, i.precision))
	}

	i.values = emas
}
