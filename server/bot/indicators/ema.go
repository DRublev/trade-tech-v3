package indicators

import (
	"errors"

	"github.com/shopspring/decimal"
)

// EmaIndicator Экспоненциальная скользящая средняя
type EmaIndicator struct {
	Indicator[decimal.Decimal, []float64]
	period     int
	prevPrices []float64
	values     []decimal.Decimal
	sma        SmaIndicator
}

// NewEma Конструктор
func NewEma(period int) *EmaIndicator {
	inst := &EmaIndicator{}
	inst.period = period
	inst.prevPrices = []float64{}
	inst.values = []decimal.Decimal{}
	inst.sma = *NewSma(period)
	return inst
}

// Latest Получить последнее значение
func (i *EmaIndicator) Latest() (decimal.Decimal, error) {
	if len(i.values) == 0 {
		return decimal.NewFromInt(0), errors.New("No latest value")
	}

	return i.values[len(i.values)-1], nil
}

// Get Получить все значения
func (i *EmaIndicator) Get() []decimal.Decimal {
	return i.values
}

// Update Уточнить значение. Юзать при поступлени новых данных
func (i *EmaIndicator) Update(prices []float64) {
	i.sma.Update(prices)
	roundPrecision := detectPrecision(prices[0])

	prevEma, err := i.sma.Latest()
	if err != nil {
		return
	}

	i.values = append(i.values, prevEma)

	k := 2 / float64(1+i.period)
	for _, p := range prices[i.period:] {
		prevEma = i.values[len(i.values)-1]
		floatEma, _ := prevEma.Float64()
		ema := (p * k) + (floatEma * (1 - k))
		i.values = append(i.values, decimal.NewFromFloat(roundFloat(ema, roundPrecision)))
	}

}
