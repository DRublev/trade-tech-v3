package controllers

import (
	"context"
	bot "main/bot"
	"main/bot/strategies"
	trade "main/server/contracts/contracts.trade"

	log "github.com/sirupsen/logrus"
)

var tradeL = log.WithFields(log.Fields{
	"controller": "trade",
})

func (s *Server) Start(ctx context.Context, in *trade.StartRequest) (*trade.StartResponse, error) {
	tradeL.Info("Start requested")
	strategyPool := bot.NewPool()

	tradeL.Tracef("Calling poll to start a strategy %v for %v", in.Strategy, in.InstrumentId)
	ok, err := strategyPool.Start(strategies.StrategyKey(in.Strategy), in.InstrumentId)
	errMsg := ""
	if err != nil {
		log.Errorf("Error starting strategy: %v", err)
		errMsg = err.Error()
	}

	tradeL.Info("Start responding")
	return &trade.StartResponse{
		Ok:    ok,
		Error: errMsg,
	}, err
}

func (s *Server) Stop(ctx context.Context, in *trade.StopRequest) (*trade.StopResponse, error) {
	tradeL.Info("Stop requested")
	strategyPool := bot.NewPool()

	tradeL.Tracef("Calling poll to stop a strategy %v for %v", in.Strategy, in.InstrumentId)
	ok, err := strategyPool.Stop(strategies.StrategyKey(in.Strategy), in.InstrumentId)

	errMsg := ""
	if err != nil {
		log.Errorf("Error stopping strategy: %v", err)
		errMsg = err.Error()
	}

	tradeL.Info("Stop responding")
	return &trade.StopResponse{
		Ok:    ok,
		Error: errMsg,
	}, err
}
