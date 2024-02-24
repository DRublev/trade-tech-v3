package server

import (
	"context"
	"fmt"
	"main/bot"
	marketdata "main/grpcGW/grpcGW.marketdata"
	"main/types"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Обьявляем нвоый обработчик эндпоинта GetCandles
func (s *Server) GetCandles(ctx context.Context, in *marketdata.GetCandlesRequest) (*marketdata.GetCandlesResponse, error) {
	err := bot.Init(ctx, types.Tinkoff)
	if err != nil {
		fmt.Println("marketdata GetCandles request err", err)
		return &marketdata.GetCandlesResponse{Candles: []*marketdata.GetCandlesResponse_OHLC{}}, err
	}

	var res []*marketdata.GetCandlesResponse_OHLC

	// Вызываем созданный ранее сервис
	candles, err := bot.Broker.GetCandles(
		in.InstrumentId,
		types.Interval(in.Interval),
		in.Start.AsTime(),
		in.End.AsTime())

	if err != nil {
		return &marketdata.GetCandlesResponse{Candles: res}, err
	}

	// Мапим в нужный формат
	for _, candle := range candles {
		o := marketdata.GetCandlesResponse_Quant{
			Units: int32(candle.Open.Units),
			Nano:  int32(candle.Open.Nano),
		}
		h := marketdata.GetCandlesResponse_Quant{
			Units: int32(candle.High.Units),
			Nano:  int32(candle.High.Nano),
		}
		l := marketdata.GetCandlesResponse_Quant{
			Units: int32(candle.Low.Units),
			Nano:  int32(candle.Low.Nano),
		}
		c := marketdata.GetCandlesResponse_Quant{
			Units: int32(candle.Close.Units),
			Nano:  int32(candle.Close.Nano),
		}
		res = append(res, &marketdata.GetCandlesResponse_OHLC{
			Open:  &o,
			High:  &h,
			Low:   &l,
			Close: &c,
			Time:  timestamppb.New(candle.Time),
			Volume: candle.Volume,
		})
	}

	fmt.Println("marketdata GetCandles request")
	return &marketdata.GetCandlesResponse{Candles: res}, nil
}
