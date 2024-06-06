package indicators

import (
	"errors"

	"github.com/shopspring/decimal"
)

type MacdResult struct {
	Signal decimal.Decimal
	Value  decimal.Decimal
}

// MacdIndicator Дивергенция скользящих средних
type MacdIndicator struct {
	Indicator[MacdResult, []float64]
	signalPeriod int
	prevPrices   []float64
	values       []decimal.Decimal
	emaSlow      EmaIndicator
	emaFast      EmaIndicator
	emaSignal    EmaIndicator
}

// NewMacd Конструктор
func NewMacd(periodSlow int, periodFast int, signalPeriod int) *MacdIndicator {
	inst := &MacdIndicator{}
	inst.signalPeriod = signalPeriod
	inst.prevPrices = []float64{}
	inst.values = []decimal.Decimal{}
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
	signal, err := i.emaSignal.Latest()
	if err != nil {
		return MacdResult{}, err
	}
	return MacdResult{
		Value:  i.values[len(i.values)-1],
		Signal: signal,
	}, nil
}

// Get Получить все значения
func (i *MacdIndicator) Get() []MacdResult {
	// TODO: Допилить
	return []MacdResult{}
}

// Update Уточнить значение. Юзать при поступлени новых данных
func (i *MacdIndicator) Update(prices []float64) {
	i.emaSlow.Update(prices)
	i.emaFast.Update(prices)

	emaSlow := i.emaSlow.Get()
	if len(emaSlow) == 0 {
		return
	}
	emaFast := i.emaFast.Get()
	if len(emaFast) == 0 {
		return
	}

	roundPrecision := emaFast[0].Exponent()

	emaFast = emaFast[len(emaFast)-len(emaSlow):]

	macd := make([]float64, 0)
	for i := 0; i < len(emaSlow); i++ {
		diff, _ := emaFast[i].Sub(emaSlow[i]).Float64()
		macd = append(macd, roundFloat(diff, uint(roundPrecision)))
	}

	i.emaSignal.Update(macd)
	_, err := i.emaSignal.Latest()
	if err != nil {
		return
	}
	i.values = append(i.values, decimal.NewFromFloat(macd[len(macd)-1]))
}
