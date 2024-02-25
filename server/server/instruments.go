package server

import (
	"context"
	"fmt"
	"main/bot"
	shares "main/grpcGW/grpcGW.shares"
	"main/types"
)

func (s *Server) GetShares(ctx context.Context, in *shares.GetInstrumentsRequest) (*shares.GetSharesResponse, error) {
	fmt.Println("blya")
	err := bot.Init(ctx, types.Tinkoff)
	if err != nil {
		fmt.Println("shares GetShares request err", err)
		return &shares.GetSharesResponse{Instruments: []*shares.Share{}}, err
	}

	var res []*shares.Share

	// Вызываем созданный ранее сервис
	_, err = bot.Broker.GetShares(types.InstrumentStatus(in.InstrumentStatus),)

	if err != nil {
		return &shares.GetSharesResponse{Instruments: res}, err
	}

	fmt.Println("shares GetShares request")
	return &shares.GetSharesResponse{Instruments: res}, nil
}
