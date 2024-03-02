package server

import (
	"context"
	"errors"
	"fmt"
	"main/bot"
	"main/bot/orderbook"
	marketdata "main/grpcGW/grpcGW.marketdata"
	"main/types"
	"math"
	"os"
	"os/signal"
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

func toMDQuant(q *types.Quant) *marketdata.Quant {
	return &marketdata.Quant{
		Units: int32(q.Units),
		Nano:  int32(q.Nano),
	}
}

func roundFloat(val float32, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(float64(val)*ratio) / ratio
}
func toMDQuantFromNum(p float32) *marketdata.Quant {
	units := math.Floor(float64(p))
	nano := roundFloat(p-float32(units), 9)

	return &marketdata.Quant{
		Units: int32(units),
		Nano:  int32(nano),
	}
}

func (s *Server) SubscribeCandles(in *marketdata.SubscribeCandlesRequest, stream marketdata.MarketData_SubscribeCandlesServer) error {
	var err error

	ctx := stream.Context()
	err = bot.Init(ctx, types.Tinkoff)
	if err != nil {
		fmt.Println("marketdata SubscribeCandles request err", err)
		return err
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(ctx context.Context, instrumentId string, interval int32) {
		defer wg.Done()
		candlesCh := make(chan types.OHLC)
		if ctx == nil {
			fmt.Println("78 marketdata", "ctx is nil")

			return
		}

		fmt.Println("83 marketdata", ctx, &candlesCh, instrumentId, interval)

		er := bot.Broker.SubscribeCandles(ctx, &candlesCh, instrumentId, types.Interval(interval))
		if er != nil {
			fmt.Println("80 marketdata", er)

			return
		}

		for c := range candlesCh {
			fmt.Println("New candle ", c.Time)
			err = stream.Send(&marketdata.OHLC{
				Open:   toMDQuant(&c.Open),
				High:   toMDQuant(&c.High),
				Low:    toMDQuant(&c.Low),
				Close:  toMDQuant(&c.Close),
				Time:   timestamppb.New(c.Time),
				Volume: c.Volume,
			})
			if err != nil {
				return
			}
		}
	}(ctx, in.InstrumentId, in.Interval)

	wg.Wait()

	return err
}

func toMDBidAsk(in []*types.BidAsk) []*marketdata.BidAsk {
	var items []*marketdata.BidAsk

	for _, inItem := range in {
		item := &marketdata.BidAsk{
			Price:    toMDQuantFromNum(inItem.Price),
			Quantity: inItem.Quantity,
		}
		items = append(items, item)
	}

	return items
}

func (s *Server) SubscribeOrderbook(in *marketdata.SubscribeOrderbookRequest, stream marketdata.MarketData_SubscribeOrderbookServer) error {
	var err error

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	err = bot.Init(ctx, types.Tinkoff)
	if err != nil {
		fmt.Println("marketdata SubscribeOrderbook request err", err)
		return err
	}

	orderbookProvider := orderbook.NewOrederbookProvider()
	orderbookCh, err := orderbookProvider.GetOrCreate(in.InstrumentId)
	if err != nil {
		return err
	}

	streamCtx := stream.Context()
	err = bot.Broker.SubscribeOrderbook(streamCtx, orderbookCh, in.InstrumentId, in.Depth)

	select {
	case <-streamCtx.Done():
		return err
	case o, ok := <-*orderbookCh:
		if !ok {
			return errors.New("stream is end")
		}
		err = stream.Send(&marketdata.Orderbook{
			InstrumentId: o.InstrumentId,
			Depth:        o.Depth,
			Time:         timestamppb.New(o.Time),
			LimitUp:      toMDQuant(&o.LimitUp),
			LimitDown:    toMDQuant(&o.LimitDown),
			Bids:         toMDBidAsk(o.Bids),
			Asks:         toMDBidAsk(o.Asks),
		})
	}

	return err
}
