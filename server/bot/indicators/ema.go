package indicators

import (
	"errors"
)

// EmaIndicator Экспоненциальная скользящая средняя
type EmaIndicator struct {
	Indicator[float64, []float64]
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
func (i *EmaIndicator) Update(prices []float64) {
	i.prevPrices = append(i.prevPrices, prices...)

	if i.latestCalcedPriceIdx != 0 {
		i.values = append(i.values, i.prevPrices[0])
		i.latestCalcedPriceIdx = 0
	}

	j := i.latestCalcedPriceIdx + 1
	for j < len(i.prevPrices) {
		// EMA = (price(t) * k) + (EMA(t - 1) * (1 – k))
		ema := i.prevPrices[j]*i.k + i.values[j-1]*(1-i.k)

		i.values = append(i.values, ema)
		j++
	}
	i.latestCalcedPriceIdx = j
}
