package orderbook

import (
	"context"
	"fmt"
	"main/bot/broker"
	"main/types"
	"os"
	"os/signal"
	"sync"
)

type OrderbookChannels struct {
	sync.RWMutex
	value map[string]*chan *types.Orderbook
}

type OrderbookProvider struct {
	channels OrderbookChannels
}

var onceOp sync.Once
var op *OrderbookProvider

func NewOrderbookProvider() *OrderbookProvider {
	if op != nil {
		return op
	}

	onceOp.Do(func() {
	})
	op := &OrderbookProvider{}
	op.channels = OrderbookChannels{
		value: make(map[string]*chan *types.Orderbook),
	}
	fmt.Printf("32 orderbookProvider %v\n", op)
	return op
}

func (op *OrderbookProvider) GetOrCreate(instrumentId string) (*chan *types.Orderbook, error) {
	fmt.Printf("36 orderbookProvider %v\n", op)
	op.channels.RLock()
	ch, exists := op.channels.value[instrumentId]
	op.channels.RUnlock()

	if !exists {
		fmt.Printf("Creating orderbook channel for %v\n", instrumentId)
		op.channels.Lock()
		newCh := make(chan *types.Orderbook)
		op.channels.value[instrumentId] = &newCh
		ch = &newCh
		op.channels.Unlock()
	}
	backCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	go broker.Broker.SubscribeOrderbook(backCtx, ch, instrumentId, 30)

	return ch, nil
}
