package orderbook

import (
	"context"
	"main/bot/broker"
	"main/types"
	"os"
	"os/signal"
	"sync"

	log "github.com/sirupsen/logrus"
)

type orderbookChannels struct {
	sync.RWMutex
	value map[string]*chan *types.Orderbook
}

// Provider Провайдер для стакана
type Provider struct {
	channels orderbookChannels
}

var onceOp sync.Once
var op *Provider

// NewProvider Конструктор для Provider
func NewProvider() *Provider {
	if op != nil {
		return op
	}

	// onceOp.Do(func() {
	log.Info("Creating orderbook provider")
	op := &Provider{}
	op.channels = orderbookChannels{
		value: make(map[string]*chan *types.Orderbook),
	}
	// })

	return op
}

// GetOrCreate Подписывается на стакан для инструмента instrumentID
// Возвращает канал для стакана или создает новый
func (op *Provider) GetOrCreate(instrumentID string) (*chan *types.Orderbook, error) {
	log.Infof("Getting orderbook channel for %v", instrumentID)

	op.channels.RLock()
	ch, exists := op.channels.value[instrumentID]
	op.channels.RUnlock()

	if !exists {
		log.Tracef("No orderbook channel found for %v, creating a new one", instrumentID)
		op.channels.Lock()
		newCh := make(chan *types.Orderbook)
		op.channels.value[instrumentID] = &newCh
		ch = &newCh
		op.channels.Unlock()
	}
	backCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	go broker.Broker.SubscribeOrderbook(backCtx, ch, instrumentID, 30)

	return ch, nil
}
