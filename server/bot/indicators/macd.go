package indicators

import (
	"errors"
)

// MacdIndicator Дивергенция скользящих средних
type MacdIndicator struct {
	signalPeriod int
	prevPrices   []float64
	values       []float64
	signals      []float64
	emaSlow      EmaIndicator
	emaFast      EmaIndicator
	emaSignal    EmaIndicator
}

// NewMacd Конструктор
func NewMacd(periodSlow int, periodFast int, signalPeriod int) *MacdIndicator {
	inst := &MacdIndicator{}
	inst.signalPeriod = signalPeriod
	inst.prevPrices = []float64{}
	inst.values = []float64{}
	inst.signals = []float64{}
	inst.emaSlow = *NewEma(periodSlow)
	inst.emaFast = *NewEma(periodFast)
	inst.emaSignal = *NewEma(signalPeriod)
	return inst
}

// Latest Получить последнее значение
func (i *MacdIndicator) Latest() (float64, float64, error) {
	if len(i.values) == 0 {
		return 0, 0, errors.New("No latest value")
	}

	_, err := i.emaSignal.Latest()
	if err != nil {
		return 0, 0, err
	}
	return i.values[len(i.values)-1], i.signals[len(i.signals)-1], nil
}

// Get Получить все значения
func (i *MacdIndicator) Get() ([]float64, []float64) {
	return i.values, i.signals
}

// Update Уточнить значение. Юзать при поступлени новых данных
// MACD = fast EMA - slow EMA
func (i *MacdIndicator) Update(price float64) {
	i.prevPrices = append(i.prevPrices, price)
	i.emaSlow.Update(price)
	i.emaFast.Update(price)
	// Количество данных должно быть минимум х2 от периода рассчета
	if len(i.prevPrices) < 2*i.emaSlow.period {
		return
	}

	ema26 := i.emaSlow.Get()
	ema12 := i.emaFast.Get()

	roundPrecision := detectPrecision(ema12[0])

	ema12 = ema12[len(ema12)-len(ema26):]

	macd := make([]float64, 0)
	for j := 0; j < len(ema26); j++ {
		macdCurrent := roundFloat(ema12[j]-ema26[j], roundPrecision)
		macd = append(macd, macdCurrent)
		i.emaSignal.Update(macdCurrent)
	}

	i.values = macd
	i.signals = i.emaSignal.Get()
}
