package indicators

import (
	"errors"
)

type MacdResult struct {
	Signal float64
	Value  float64
}

// MacdIndicator Дивергенция скользящих средних
type MacdIndicator struct {
	Indicator[MacdResult, []float64]
	signalPeriod int
	prevPrices   []float64
	values       []MacdResult
	emaSlow      EmaIndicator
	emaFast      EmaIndicator
	emaSignal    EmaIndicator
}

// NewMacd Конструктор
func NewMacd(periodSlow int, periodFast int, signalPeriod int) *MacdIndicator {
	inst := &MacdIndicator{}
	inst.signalPeriod = signalPeriod
	inst.prevPrices = []float64{}
	inst.values = []MacdResult{}
	inst.emaSlow = *NewEma(periodSlow)
	inst.emaFast = *NewEma(periodFast)
	inst.emaSignal = *NewEma(signalPeriod)
	return inst
}

// Latest Получить последнее значение
func (i *MacdIndicator) Latest() (MacdResult, error) {
	if len(i.values) == 0 {
		return MacdResult{}, errors.New("No latest value")
	}

	_, err := i.emaSignal.Latest()
	if err != nil {
		return MacdResult{}, err
	}
	return i.values[len(i.values)-1], nil
}

// Get Получить все значения
func (i *MacdIndicator) Get() []MacdResult {
	return i.values
}

// Update Уточнить значение. Юзать при поступлени новых данных
// MACD = fast EMA - slow EMA
func (i *MacdIndicator) Update(prices []float64) {
	i.emaFast.Update(prices)
	i.emaSlow.Update(prices)

	slowEma, err := i.emaSlow.Latest()
	fastEma, err := i.emaFast.Latest()
	if err != nil {
		return
	}

	macd := fastEma - slowEma
	i.emaSignal.Update([]float64{macd})
	signal, err := i.emaSignal.Latest()
	if err != nil {
		return
	}

	i.values = append(i.values, MacdResult{
		Signal: signal,
		Value:  macd,
	})
}
