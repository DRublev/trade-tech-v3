package server

import (
	"context"
	"fmt"
	"main/bot"
	 instruments "main/grpcGW/grpcGW.instruments"
	"main/types"
)

func (s *Server) GetShares(ctx context.Context, in *instruments.GetInstrumentsRequest) (*instruments.GetSharesResponse, error) {
	err := bot.Init(ctx, types.Tinkoff)
	if err != nil {
		fmt.Println("instrument GetShares request err", err)
		return &instruments.GetSharesResponse{Instruments: []*instruments.Share{}}, err
	}

	var res []*instruments.Share

	// Вызываем созданный ранее сервис
	shares, err := bot.Broker.GetShares(types.InstrumentStatus{})

	if err != nil {
		return &instruments.GetSharesResponse{Instruments: res}, err
	}

	fmt.Println("instrumetns GetShares request")
	fmt.Println(shares)
	return &instruments.GetSharesResponse{Instruments: res}, nil
}