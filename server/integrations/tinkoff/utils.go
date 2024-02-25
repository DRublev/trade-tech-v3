package tinkoff

import (
	"main/types"

	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const nanoPrecision = 1_000_000_000

func quantToNumber(q types.Quant) float64 {
	return float64(q.Units) + (float64(q.Nano) / nanoPrecision)
}

func toQuant(iq *investapi.Quotation) types.Quant {
	return types.Quant{
		Units: int(iq.Units),
		Nano:  int(iq.Nano),
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
