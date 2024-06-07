package indicators

import (
	"errors"
)

// SmaIndicator Простая скользящая средняя
// TODO: Допилить, не очень то и рабочая реализация
// Обычное вычисление средней - сумма всех значений за период / количество значений
type SmaIndicator struct {
	Indicator[float64, float64]
	period     int
	prevPrices []float64
	values     []float64
}

// NewSma Конструктор
func NewSma(period int) *SmaIndicator {
	inst := &SmaIndicator{}
	inst.period = period
	inst.prevPrices = []float64{}
	inst.values = []float64{}
	return inst
}

// Latest Получить последнее значение
func (i *SmaIndicator) Latest() (float64, error) {
	if len(i.values) == 0 {
		return 0, errors.New("No latest value")
	}
	return i.values[len(i.values)-1], nil
}

// Get Получить все значения
func (i *SmaIndicator) Get() []float64 {
	return i.values
}

// Update Уточнить значение. Юзать при поступлени новых данных
func (i *SmaIndicator) Update(price float64) {
	i.prevPrices = append(i.prevPrices, price)
	// Недостаточно данных для рассчета
	if len(i.prevPrices) < i.period {
		return
	}

	roundPrecision := detectPrecision(price)

	var sma []float64
	for j := 0; j+i.period <= len(i.prevPrices); j++ {
		var sum float64
		for _, p := range i.prevPrices[j : j+i.period] {
			sum += p
		}

		sma = append(sma, roundFloat(sum/float64(i.period), roundPrecision))
	}
	i.values = sma
}
