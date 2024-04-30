package controllers

import (
	"context"
	"main/bot/broker"
	"main/bot/orderbook"
	marketdata "main/server/contracts/contracts.marketdata"
	"main/types"
	"math"
	"os"
	"os/signal"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var mdL = log.WithFields(log.Fields{
	"controller": "marketdata",
})

func (s *Server) GetCandles(ctx context.Context, in *marketdata.GetCandlesRequest) (*marketdata.GetCandlesResponse, error) {
	mdL.WithField("instrument", in.InstrumentId).Info("GetCandles requested")

	err := broker.Init(ctx, types.Tinkoff)
	if err != nil {
		mdL.Errorf("Cannot init broker: %v", err)
		return &marketdata.GetCandlesResponse{Candles: []*marketdata.OHLC{}}, err
	}

	var res []*marketdata.OHLC

	mdL.Trace("Requesting broker for candles")
	// Вызываем созданный ранее сервис
	candles, err := broker.Broker.GetCandles(
		in.InstrumentId,
		types.Interval(in.Interval),
		in.Start.AsTime(),
		in.End.AsTime())

	if err != nil {
		mdL.Errorf("Failed getting candles from broker: %v", err)
		return &marketdata.GetCandlesResponse{Candles: res}, err
	}

	mdL.Tracef("Got %v candles, mapping", len(candles))
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

	mdL.Info("GetCandles responding")
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

// TODO: Изменить сигнатуру на () [units, nano] и вынести в utils
func toMDQuantFromNum(p float32) *marketdata.Quant {
	units := math.Floor(float64(p))
	nano := roundFloat(p-float32(units), 9)

	return &marketdata.Quant{
		Units: int32(units),
		Nano:  int32(nano),
	}
}
func (s *Server) SubscribeOrders(in *marketdata.SubscribeOrderRequest, stream marketdata.MarketData_SubscribeOrdersServer) error {
	mdL.Info("SubscribeOrders requested")

	bCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	err := broker.Init(bCtx, types.Tinkoff)
	if err != nil {
		mdL.Errorf("Cannot init broker: %v", err)
		return err
	}

	mdL.Trace("Requesting broker for subscribing to orderState")
	err = broker.Broker.SubscribeOrders(func(st types.OrderExecutionState) {
		err = stream.Send(&marketdata.OrderState{
			IdempodentID:    string(st.ID), // TODO: Поправить несоответствие
			ExecutionStatus: int32(st.Status),
			OperationType:   int32(st.Direction),
			LotsRequested:   int32(st.LotsRequested),
			LotsExecuted:    int32(st.LotsExecuted),
			InstrumentID:    st.InstrumentID,
			Strategy:        "", // TODO: Брать стратегию из единого мета
			PricePerLot:     st.ExecutedOrderPrice / float64(st.LotsExecuted),
			Time:            timestamppb.New(time.Now()),
		})
		if err != nil {
			mdL.Warnf("Failed sending orderState to stream: %v", err)
		}
	})
	if err != nil {
		mdL.Errorf("Failed subscribing orderState: %v", err)
		return err
	}

	mdL.Info("SubscribeOrders responding")
	return err
}

func (s *Server) SubscribeCandles(in *marketdata.SubscribeCandlesRequest, stream marketdata.MarketData_SubscribeCandlesServer) error {
	mdL.WithField("instrument", in.InstrumentId).Info("SubscribeCandles requested")

	ctx := stream.Context()
	bCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	err := broker.Init(bCtx, types.Tinkoff)
	if err != nil {
		mdL.Errorf("Cannot init broker: %v", err)
		return err
	}

	candlesCh := make(chan types.OHLC)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(ctx context.Context, ch *chan types.OHLC) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				mdL.Info("Subscribe candles context closed")
				return
			case c, ok := <-*ch:
				if !ok {
					mdL.Infof("Subscribe candles stream is done for %v", in.InstrumentId)
					return
				}

				mdL.Tracef("Sending new candle to stream. Candle time: %v", c.Time)
				err = stream.Send(&marketdata.OHLC{
					Open:   toMDQuant(&c.Open),
					High:   toMDQuant(&c.High),
					Low:    toMDQuant(&c.Low),
					Close:  toMDQuant(&c.Close),
					Time:   timestamppb.New(c.Time),
					Volume: c.Volume,
				})
				if err != nil {
					mdL.Warnf("Failed sending candle to stream: %v", err)
				}
			}
		}
	}(bCtx, &candlesCh)

	mdL.Trace("Requesting broker for subscribing to candles")
	err = broker.Broker.SubscribeCandles(ctx, &candlesCh, in.InstrumentId, types.Interval(in.Interval))
	if err != nil {
		mdL.Errorf("Failed subscribing candles: %v", err)
		return err
	}

	wg.Wait()

	mdL.Info("SubscribeCandles responding")
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
	mdL.WithField("instrument", in.InstrumentId).Info("SubscribeOrderbook requested")
	var err error

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	err = broker.Init(ctx, types.Tinkoff)
	if err != nil {
		mdL.Errorf("Cannot init broker: %v", err)
		return err
	}

	mdL.Trace("Creating orderbook provider and channel")
	orderbookProvider := orderbook.NewProvider()
	orderbookCh, err := orderbookProvider.GetOrCreate(in.InstrumentId)
	if err != nil {
		mdL.Errorf("Failed getting channel for orderbook: %v", err)
		return err
	}

	streamCtx := stream.Context()

	mdL.Trace("Requesting broker for subscribe to orderbook")
	err = broker.Broker.SubscribeOrderbook(streamCtx, orderbookCh, in.InstrumentId, in.Depth)
	select {
	case <-streamCtx.Done():
		mdL.Infof("Subscribe orderbook context closed for %v", in.InstrumentId)
		return err
	case o, ok := <-*orderbookCh:
		if !ok {
			mdL.Info("Subscribe orderbook")
			return nil
		}
		mdL.Tracef("Sending orderbook to stream. Orderbok time: %v", o.Time)
		err = stream.Send(&marketdata.Orderbook{
			InstrumentId: o.InstrumentId,
			Depth:        o.Depth,
			Time:         timestamppb.New(o.Time),
			LimitUp:      toMDQuant(&o.LimitUp),
			LimitDown:    toMDQuant(&o.LimitDown),
			Bids:         toMDBidAsk(o.Bids),
			Asks:         toMDBidAsk(o.Asks),
		})

		if err != nil {
			mdL.Warnf("Failed sending orderbook to stream: %v", err)
		}
	}

	mdL.Info("SubscribeOrderbook responding")
	return err
}
