package server

import (
	"context"
	"fmt"
	bot "main/bot"
	"main/bot/strategies"
	trade "main/grpcGW/grpcGW.trade"
)

func (s *Server) Start(ctx context.Context, in *trade.StartRequest) (*trade.StartResponse, error) {
	strategyPool := bot.NewPool()

	ok, err := strategyPool.Start(strategies.StrategyKey(in.Strategy), in.InstrumentId)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	fmt.Printf("18 trade %v\n", err)
	return &trade.StartResponse{
		Ok:    ok,
		Error: errMsg,
	}, err
}

func (s *Server) Stop(ctx context.Context, in *trade.StopRequest) (*trade.StopResponse, error) {
	strategyPool := bot.NewPool()

	ok, err := strategyPool.Stop(strategies.StrategyKey(in.Strategy), in.InstrumentId)

	return &trade.StopResponse{
		Ok:    ok,
		Error: err.Error(),
	}, err
}
