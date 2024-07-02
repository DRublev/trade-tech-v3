package bot

import (
	"context"
	"errors"
	"main/bot/broker"
	"main/bot/candles"
	"main/bot/strategies"
	"main/bot/strategies/macd"
	"main/bot/strategies/rosshook"
	"main/bot/strategies/spread"
	"main/types"
	"os"
	"os/signal"
)

type StreamerCb func(event StrategyEvent)
type StrategyEvent struct {
	Type string
}

type Streamer struct {
	subscribers []StreamerCb
}

func NewStreamer() {
	
}

func (s *Streamer) Subscribe(cb StreamerCb) {
	s.subscribers = append(s.subscribers, cb)
}

func (s *Streamer) notify(e StrategyEvent) {
	for _, cb := range s.subscribers {
		cb(e)
	}
}

// Assemble Фабрика для сборки стратегий
func Assemble(key strategies.StrategyKey, config *strategies.Config) (strategies.IStrategy, error) {
	backCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	err := broker.Init(backCtx, types.Tinkoff)
	if err != nil {
		return nil, err
	}

	switch key {
	case strategies.Spread:
		return spread.New(), nil
	case strategies.Macd:
		return macd.New(), nil
	case strategies.RossHook:
		return rosshook.New(candles.NewProvider()), nil
	}

	// TODO: Инициализировать стратегию в зависимости от ключа
	return nil, errors.New("Strategy not implemented")
}
