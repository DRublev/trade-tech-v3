package candles

import (
	"context"
	"main/bot/broker"
	"main/types"
	"os"
	"os/signal"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type candlesChannels map[string]*chan types.OHLC

// Provider Провайдер свечей
type Provider struct {
	channels candlesChannels
}

var once sync.Once
var instance *Provider

// NewProvider Конструктор провайдера свечей
func NewProvider() *Provider {
	if instance != nil {
		return instance
	}

	// once.Do(func() {
	log.Infof("Creating candles provider")

	instance := &Provider{}
	instance.channels = make(map[string]*chan types.OHLC)
	// })

	return instance
}

// GetOrCreate Создать провайдер или взять готовый
func (p *Provider) GetOrCreate(instrumentID string, initialFrom time.Time, initialTo time.Time) (*chan types.OHLC, error) {
	log.Infof("Getting candles channel for %v", instrumentID)

	ch, exists := p.channels[instrumentID]

	if !exists {
		log.Tracef("No candles channel found for %v, creating a new one", instrumentID)
		newCh := make(chan types.OHLC)
		p.channels[instrumentID] = &newCh
		ch = &newCh
	}

	backCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	log.Tracef("Getting candles from %v to %v", initialFrom, initialTo)
	initialCandles, err := broker.Broker.GetCandles(instrumentID, 1, initialFrom, initialTo)
	if err != nil {
		log.Errorf("Error getting candles %v", err)
	} else {
		log.Tracef("Sending initial candles from %v for %v", initialFrom, instrumentID)
		go func() {
			for _, candle := range initialCandles {
				*ch <- candle
			}
		}()
	}

	go broker.Broker.SubscribeCandles(backCtx, ch, instrumentID, 1)

	return ch, nil
}
