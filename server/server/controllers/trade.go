package controllers

import (
	"context"
	"encoding/json"
	"errors"
	bot "main/bot"
	config "main/bot/config"
	"main/bot/strategies"
	trade "main/server/contracts/contracts.trade"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	log "github.com/sirupsen/logrus"
)

var tradeL = log.WithFields(log.Fields{
	"controller": "trade",
})

var configRepository = config.New()

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

// TODO: Допилить, чтобы возвращал инфу о том  какие заявки закрыл и что вообще сделал для остановки, потом это отображать на фронте
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

func (s *Server) IsStarted(ctx context.Context, in *trade.StartRequest) (*trade.StartResponse, error) {
	tradeL.Info("IsStarted requested")
	strategyPool := bot.NewPool()
	isStarted, err := strategyPool.IsStarted(strategies.StrategyKey(in.Strategy), in.InstrumentId)

	tradeL.Tracef("IsStarted responding: %v (%v for %v)", isStarted, in.Strategy, in.InstrumentId)
	return &trade.StartResponse{
		Ok: isStarted,
	}, err
}

func (s *Server) ChangeConfig(ctx context.Context, in *trade.ChangeConfigRequest) (*trade.ChangeConfigResponse, error) {
	tradeL.Info("ChangeConfig requested")

	configKey := in.Strategy + "_" + in.InstrumentId
	config := make(strategies.Config)

	for key, value := range in.Config.AsMap() {
		config[key] = value
	}
	config["InstrumentID"] = in.InstrumentId

	err := configRepository.Set(configKey, config)
	errMsg := ""
	if err != nil {
		log.Errorf("Error setting config: %v", err)
		errMsg = err.Error()
	}

	tradeL.Info("ChangeConfig responding")
	return &trade.ChangeConfigResponse{
		Ok:    true,
		Error: errMsg,
	}, err
}

func (s *Server) GetConfig(ctx context.Context, in *trade.GetConfigRequest) (*trade.GetConfigResponse, error) {
	tradeL.Info("GetConfig requested")

	configKey := in.Strategy + "_" + in.InstrumentId

	c, err := configRepository.Get(configKey)
	if err != nil {
		return nil, err
	}

	res := &trade.GetConfigResponse{}

	b, err := json.Marshal(c)
	configStruct := &structpb.Struct{}

	err = protojson.Unmarshal(b, configStruct)
	if err != nil {
		return nil, err
	}

	res.Config = configStruct

	tradeL.Info("GetConfig responding")
	return res, nil

}

func (s *Server) SubscribeStrategiesEvents(in *trade.SubscribeStrategiesEventsRequest, stream trade.Trade_SubscribeStrategiesEventsServer) error {
	tradeL.WithField("strategy", in.Strategy).Info("SubscribeStrategiesEvents requested")

	activitiesPubSub := strategies.NewActivityPubSub()

	activitiesChannel := activitiesPubSub.Subscribe(in.Strategy)

	if activitiesChannel == nil {
		return errors.New("No container for " + in.Strategy)
	}

	streamCtx := stream.Context()
	for {
		select {
		case <-streamCtx.Done():
			activitiesChannel = nil
			return nil
		case activity, ok := <-*activitiesChannel:
			if !ok {
				tradeL.Info("Subscription closed")
				continue
			}

			b, err := json.Marshal(activity.Value)
			activityStruct := &structpb.Struct{}

			err = protojson.Unmarshal(b, activityStruct)
			if err != nil {
				tradeL.Warnf("Activity conversion error %v", err)
				continue
			}

			err = stream.Send(&trade.StrategyEvent{
				ID:    activity.ID,
				Kind:  string(activity.Kind),
				Value: activityStruct,
			})
			if err != nil {
				tradeL.Warnf("Activity send error %v", err)
			}
		}
	}

	return nil
}
