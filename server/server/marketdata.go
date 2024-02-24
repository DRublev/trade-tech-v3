package server

import (
	"context"
	"fmt"
	"main/bot"
	marketdata "main/grpcGW/grpcGW.marketdata"
	"main/types"
	"sync"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Обьявляем нвоый обработчик эндпоинта GetCandles
func (s *Server) GetCandles(ctx context.Context, in *marketdata.GetCandlesRequest) (*marketdata.GetCandlesResponse, error) {
	err := bot.Init(ctx, types.Tinkoff)
	if err != nil {
		fmt.Println("marketdata GetCandles request err", err)
		return &marketdata.GetCandlesResponse{Candles: []*marketdata.OHLC{}}, err
	}

	var res []*marketdata.OHLC

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
		o := marketdata.Quant{
			Units: int32(candle.Open.Units),
			Nano:  int32(candle.Open.Nano),
		}
		h := marketdata.Quant{
			Units: int32(candle.High.Units),
			Nano:  int32(candle.High.Nano),
		}
		l := marketdata.Quant{
			Units: int32(candle.Low.Units),
			Nano:  int32(candle.Low.Nano),
		}
		c := marketdata.Quant{
			Units: int32(candle.Close.Units),
			Nano:  int32(candle.Close.Nano),
		}
		res = append(res, &marketdata.OHLC{
			Open:   &o,
			High:   &h,
			Low:    &l,
			Close:  &c,
			Time:   timestamppb.New(candle.Time),
			Volume: candle.Volume,
		})
	}

	return &marketdata.GetCandlesResponse{Candles: res}, nil
}

func (s *Server) SubscribeCandles(in *marketdata.SubscribeCandlesRequest, stream marketdata.MarketData_SubscribeCandlesServer) error {
	var err error

	ctx := stream.Context()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(ctx context.Context, instrumentId string, interval int32) {
		defer wg.Done()
		candlesCh := make(chan types.OHLC)
		if ctx == nil {
			fmt.Println("78 marketdata", "ctx is nil")

			return
		}

		fmt.Println("83 marketdata", ctx, &candlesCh, instrumentId, interval);
		

		er := bot.Broker.SubscribeCandles(ctx, &candlesCh, instrumentId, types.Interval(interval))
		if er != nil {
			fmt.Println("80 marketdata", er)

			return
		}

		for candle := range candlesCh {
			fmt.Println("New candle ", candle.Time)
			o := marketdata.Quant{
				Units: int32(candle.Open.Units),
				Nano:  int32(candle.Open.Nano),
			}
			h := marketdata.Quant{
				Units: int32(candle.High.Units),
				Nano:  int32(candle.High.Nano),
			}
			l := marketdata.Quant{
				Units: int32(candle.Low.Units),
				Nano:  int32(candle.Low.Nano),
			}
			c := marketdata.Quant{
				Units: int32(candle.Close.Units),
				Nano:  int32(candle.Close.Nano),
			}
			err = stream.Send(&marketdata.OHLC{
				Open:   &o,
				High:   &h,
				Low:    &l,
				Close:  &c,
				Time:   timestamppb.New(candle.Time),
				Volume: candle.Volume,
			})
			if err != nil {
				return
			}
		}
	}(ctx, in.InstrumentId, in.Interval)

	wg.Wait()

	return err
}
