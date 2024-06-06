package indicators

import (
	"errors"
)

// SmaIndicator Простая скользящая средняя
// TODO: Допилить, не очень то и рабочая реализация
// Обычное вычисление средней - сумма всех значений за период / количество значений
type SmaIndicator struct {
	Indicator[float64, []float64]
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
func (i *SmaIndicator) Update(p []float64) {
	i.prevPrices = append(i.prevPrices, p...)
	prices := i.prevPrices

	if len(i.prevPrices) < i.period {
		return
	}

	// Рассчитываем значение первый раз
	if len(i.values) == 0 {
		// Запоминаем цены, которые уже учитывали
		// Для кейса, когда пытались инициализировать индикатор с недостатком данных, а потом пытаемся его обновить
		source := prices
		if len(source) < i.period {
			source = append(i.prevPrices, source...)
		}
		startIdx := len(source) - i.period
		// Данных недостаточно для рассчета, ниче не делаем
		if startIdx < 0 {
			i.prevPrices = source
			return
		}

		roundPrecision := detectPrecision(prices[0])

		var value float64 = 0
		for k := startIdx; k < len(prices); k++ {
			value += prices[k]
		}
		value /= float64(i.period)
		value = roundFloat(value/float64(i.period), roundPrecision)

		i.values = append(i.values, value)
		i.prevPrices = append(i.prevPrices, prices...)

		return
	}

	// TODO: Так как мы знаем предыдущие цены и значения индикатора, можем применить алгоритм slidingWindow
	// Но на маленьких периодах это не так страшно, поэтому пока не делаю

	// Нужно обновить индикатор, с учетом новых значений
	source := prices
	startIdx := len(source) - i.period
	if startIdx < 0 {
		// Докидываем старые значения, чтобы было достаточно данных для рассчета
		source = append(i.prevPrices[i.period-len(prices)-1:len(prices)-1], prices...)
		startIdx = 0
	}

	roundPrecision := detectPrecision(prices[0])
	var value float64 = 0
	for k := len(source) - i.period; k < len(prices); k++ {
		value += prices[k]
	}

	value /= float64(i.period)
	value = roundFloat(value/float64(i.period), roundPrecision)

	i.values = append(i.values, value)
	i.prevPrices = append(i.prevPrices, prices...)
	return
}
