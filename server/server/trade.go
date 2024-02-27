package server

import (
	"context"
	bot "main/bot"
	"main/bot/strategies"
	trade "main/grpcGW/grpcGW.trade"
)

func (s *Server) Start(ctx context.Context, in *trade.StartRequest) (*trade.StartResponse, error) {
	strategyPool := bot.NewPool()

	ok, err := strategyPool.Start(strategies.StrategyKey(in.Strategy), in.InstrumentId)

	return &trade.StartResponse{
		Ok:    ok,
		Error: err.Error(),
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
