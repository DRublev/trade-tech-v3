package tinkoff

import (
	"main/types"
	"math"

	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const nanoPrecision = 1_000_000_000
const BILLION int64 = 1_000_000_000

func quantToNumber(q types.Quant) float64 {
	return float64(q.Units) + (float64(q.Nano) / nanoPrecision)
}

func toQuant(iq *investapi.Quotation) types.Quant {
	if iq == nil {
		return types.Quant{Units: 0}
	}
	return types.Quant{
		Units: int(iq.Units),
		Nano:  int(iq.Nano),
	}
}

// FloatToQuotation - Перевод float в Quotation, step - шаг цены для инструмента (min_price_increment)
func FloatToQuotation(number float64, step *investapi.Quotation) investapi.Quotation {
	// делим дробь на дробь и округляем до ближайшего целого
	k := math.Round(number / step.ToFloat())
	// целое умножаем на дробный шаг и получаем готовое дробное значение
	roundedNumber := step.ToFloat() * k
	// разделяем дробную и целую части
	decNumber := decimal.NewFromFloat(roundedNumber)

	intPart := decNumber.IntPart()
	fracPart := decNumber.Sub(decimal.NewFromInt(intPart))

	nano := fracPart.Mul(decimal.NewFromInt(BILLION)).IntPart()
	return investapi.Quotation{
		Units: intPart,
		Nano:  int32(nano),
	}
}

func roundFloat(val float32, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(float64(val)*ratio) / ratio
}
func toQuotation(n float64) investapi.Quotation {
	units := math.Floor(float64(n))
	nano := roundFloat(float32(n-units), 9) * nanoPrecision

	return investapi.Quotation{
		Units: int64(units),
		Nano:  int32(nano),
	}
}

type IInvestCandleWithLastTrade interface {
	IInvestCandle
	GetLastTradeTs() *timestamppb.Timestamp
}

type IInvestCandle interface {
	GetTime() *timestamppb.Timestamp
	GetOpen() *investapi.Quotation
	GetHigh() *investapi.Quotation
	GetLow() *investapi.Quotation
	GetClose() *investapi.Quotation
	GetVolume() int64
}

func toOHLCWithTrade(c IInvestCandleWithLastTrade) types.OHLC {
	candle := types.OHLC{
		Time:        c.GetTime().AsTime(),
		Open:        toQuant(c.GetOpen()),
		High:        toQuant(c.GetHigh()),
		Low:         toQuant(c.GetLow()),
		Close:       toQuant(c.GetClose()),
		Volume:      c.GetVolume(),
		LastTradeTS: c.GetLastTradeTs().AsTime(),
	}

	return candle
}

func toOHLC(c IInvestCandle) types.OHLC {
	candle := types.OHLC{
		Time:        c.GetTime().AsTime(),
		Open:        toQuant(c.GetOpen()),
		High:        toQuant(c.GetHigh()),
		Low:         toQuant(c.GetLow()),
		Close:       toQuant(c.GetClose()),
		Volume:      c.GetVolume(),
		LastTradeTS: c.GetTime().AsTime(),
	}

	return candle
}

func toBidAsk(in []*investapi.Order) []*types.BidAsk {
	var items []*types.BidAsk
	for _, inItem := range in {
		item := &types.BidAsk{
			Price:    float32(quantToNumber(toQuant(inItem.Price))),
			Quantity: inItem.Quantity,
		}
		items = append(items, item)
	}

	return items
}
