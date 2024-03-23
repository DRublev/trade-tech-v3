package server

import (
	"context"
	"main/bot/broker"
	shares "main/grpcGW/grpcGW.shares"
	"main/types"
	"strings"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var instL = log.WithFields(log.Fields{
	"controller": "instruments",
})

func (s *Server) GetShares(ctx context.Context, in *shares.GetInstrumentsRequest) (*shares.GetSharesResponse, error) {
	instL.Info("GetShares requested")

	err := broker.Init(ctx, types.Tinkoff)
	if err != nil {
		instL.Errorf("Cannot init broker: %v", err)
		return &shares.GetSharesResponse{Instruments: []*shares.Share{}}, err
	}

	var res []*shares.Share

	instL.Trace("Requesting broker to get shares")
	sharesArr, err := broker.Broker.GetShares(types.InstrumentStatus(in.InstrumentStatus))

	if err != nil {
		instL.Errorf("Failed getting shares list: %v", err)
		return &shares.GetSharesResponse{Instruments: res}, err
	}

	instL.Tracef("Got %v shares", len(sharesArr))
	for _, share := range sharesArr {
		minPrice := shares.Quatation{
			Units: int32(share.MinPriceIncrement.Units),
			Nano:  int32(share.MinPriceIncrement.Nano),
		}
		res = append(res, &shares.Share{
			Name:                share.Name,
			Figi:                share.Figi,
			Exchange:            share.Exchange,
			Ticker:              share.Ticker,
			Lot:                 share.Lot,
			IpoDate:             timestamppb.New(share.IpoDate),
			TradingStatus:       int32(share.TradingStatus),
			MinPriceIncrement:   &minPrice,
			Uid:                 share.Uid,
			First1MinCandleDate: timestamppb.New(share.First1minCandleDate),
			First1DayCandleDate: timestamppb.New(share.First1dayCandleDate),
		})
	}

	instL.Info("GetShares responding")
	return &shares.GetSharesResponse{Instruments: res}, nil
}

func (s *Server) GetTradingSchedules(ctx context.Context, in *shares.GetTradingSchedulesRequest) (*shares.GetTradingSchedulesResponse, error) {
	instL.Info("GetTradingSchedules requested")
	err := broker.Init(ctx, types.Tinkoff)
	if err != nil {
		instL.Errorf("Cannot init broker: %v", err)
		return &shares.GetTradingSchedulesResponse{Exchanges: []*shares.TradingSchedule{}}, err
	}

	var res []*shares.TradingSchedule

	exchangesArr, err := broker.Broker.GetTradingSchedules(in.Exchange, in.From.AsTime(), in.To.AsTime())

	if err != nil {
		return &shares.GetTradingSchedulesResponse{Exchanges: res}, err
	}

	for _, exchange := range exchangesArr {
		if strings.Contains(exchange.Exchange, "MOEX") && !strings.Contains(exchange.Exchange, "WEEKEND") {
			var days []*shares.TradingDay
			for _, day := range exchange.Days {
				days = append(days, &shares.TradingDay{
					Date:                    timestamppb.New(day.Date),
					IsTradingDay:            day.IsTradingDay,
					StartTime:               timestamppb.New(day.StartTime),
					EndTime:                 timestamppb.New(day.EndTime),
					OpeningAuctionStartTime: timestamppb.New(day.OpeningAuctionEndTime),
				})
			}

			res = append(res, &shares.TradingSchedule{
				Exchange: exchange.Exchange,
				Days:     days,
			})
		}

	}

	instL.Info("GetTradingSchedules responding")
	return &shares.GetTradingSchedulesResponse{Exchanges: res}, nil
}
