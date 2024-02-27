package bot

import (
	"main/types"
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

func NewOrederbookProvider() *OrderbookProvider {
	if op != nil {
		return op
	}

	onceOp.Do(func() {
		op := &OrderbookProvider{}
		op.channels = OrderbookChannels{
			value: make(map[string]*chan *types.Orderbook),
		}
	})

	return op
}

func (op *OrderbookProvider) GetOrCreate(instrumentId string) (*chan *types.Orderbook, error) {
	op.channels.RLock()
	ch, exists := op.channels.value[instrumentId]
	op.channels.RUnlock()

	if !exists {
		op.channels.Lock()
		newCh := make(chan *types.Orderbook)
		op.channels.value[instrumentId] = &newCh
		ch = op.channels.value[instrumentId]
		op.channels.Unlock()
	}

	return ch, nil
}
