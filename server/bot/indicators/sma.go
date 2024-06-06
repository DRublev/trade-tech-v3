package indicators

import (
	"errors"

	"github.com/shopspring/decimal"
)

// SmaIndicator Простая скользящая средняя
// Обычное вычисление средней - сумма всех значений за период / количество значений
type SmaIndicator struct {
	Indicator[decimal.Decimal, []float64]
	period     int
	prevPrices []float64
	values     []decimal.Decimal
}

// NewSma Конструктор
func NewSma(period int) *SmaIndicator {
	inst := &SmaIndicator{}
	inst.period = period
	inst.prevPrices = []float64{}
	inst.values = []decimal.Decimal{}
	return inst
}

// Latest Получить последнее значение
func (i *SmaIndicator) Latest() (decimal.Decimal, error) {
	if len(i.values) == 0 {
		return decimal.NewFromInt(0), errors.New("No latest value")
	}
	return i.values[len(i.values)-1], nil
}

// Get Получить все значения
func (i *SmaIndicator) Get() []decimal.Decimal {
	return i.values
}

// Update Уточнить значение. Юзать при поступлени новых данных
func (i *SmaIndicator) Update(prices []float64) {
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
		for k := startIdx; k <= len(prices); k++ {
			value += prices[k]
		}
		value /= float64(i.period)
		value = roundFloat(value/float64(i.period), roundPrecision)

		i.values = append(i.values, decimal.NewFromFloat(value))
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
	for k := len(source) - i.period; k <= len(prices); k++ {
		value += prices[k]
	}

	value /= float64(i.period)
	value = roundFloat(value/float64(i.period), roundPrecision)

	i.values = append(i.values, decimal.NewFromFloat(value))
	i.prevPrices = append(i.prevPrices, prices...)
	return
}
