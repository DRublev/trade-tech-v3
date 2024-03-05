package tinkoff

import (
	"fmt"
	"main/types"
	"math"

	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const nanoPrecision = 1_000_000_000

// TODO: Вынести на уровень аппа
func quantToNumber(q types.Quant) float64 {
	return float64(q.Units) + (float64(q.Nano) / nanoPrecision)
}

func toQuant(iq *investapi.Quotation) types.Quant {
	return types.Quant{
		Units: int(iq.Units),
		Nano:  int(iq.Nano),
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
		Nano: int32(nano),
	}
}

type IInvestCandle interface {
	GetTime() *timestamppb.Timestamp
	GetOpen() *investapi.Quotation
	GetHigh() *investapi.Quotation
	GetLow() *investapi.Quotation
	GetClose() *investapi.Quotation
	GetVolume() int64
}

func toOHLC(c IInvestCandle) types.OHLC {
	candle := types.OHLC{
		Time:   c.GetTime().AsTime(),
		Open:   toQuant(c.GetOpen()),
		High:   toQuant(c.GetHigh()),
		Low:    toQuant(c.GetLow()),
		Close:  toQuant(c.GetClose()),
		Volume: c.GetVolume(),
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
