package server

import (
	"context"
	"fmt"
	"main/bot"
	shares "main/grpcGW/grpcGW.shares"
	"main/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) GetShares(ctx context.Context, in *shares.GetInstrumentsRequest) (*shares.GetSharesResponse, error) {
	err := bot.Init(ctx, types.Tinkoff)
	if err != nil {
		fmt.Println("shares GetShares request err", err)
		return &shares.GetSharesResponse{Instruments: []*shares.Share{}}, err
	}

	var res []*shares.Share

	sharesArr, err := bot.Broker.GetShares(types.InstrumentStatus(in.InstrumentStatus))

	if err != nil {
		return &shares.GetSharesResponse{Instruments: res}, err
	}

	for _, share := range sharesArr {
		minPrice := shares.Quatation{
			Units: int32(share.MinPriceIncrement.Units),
			Nano:  int32(share.MinPriceIncrement.Nano),
		}
		res = append(res, &shares.Share{
			Name: share.Name,
			Figi: share.Figi,
			Exchange: share.Exchange,
			Ticker: share.Ticker,
			Lot: share.Lot,
			IpoDate :  timestamppb.New(share.IpoDate),
			TradingStatus: int32(share.TradingStatus),
			MinPriceIncrement: &minPrice,
			Uid: share.Uid,
			First1MinCandleDate: timestamppb.New(share.First1minCandleDate),
			First1DayCandleDate: timestamppb.New(share.First1dayCandleDate),
		})
	}

	fmt.Println("shares GetShares request")
	return &shares.GetSharesResponse{Instruments: res}, nil
}
