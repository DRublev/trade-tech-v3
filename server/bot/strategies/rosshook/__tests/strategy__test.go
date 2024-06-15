package rosshook__tests

import (
	"encoding/json"
	"main/bot/candles"
	"main/bot/strategies"
	"main/bot/strategies/rosshook"
	"main/types"
	"testing"
	"time"
)

type MockProvider struct {
	candles.BaseCandlesProvider
}

var candlesCh = make(chan types.OHLC)

func (p MockProvider) GetOrCreate(instrumentID string, initialFrom time.Time, initialTo time.Time) (*chan types.OHLC, error) {
	return &candlesCh, nil
}

func TestTest(t *testing.T) {
	mockProvider := MockProvider{}

	strategy := rosshook.New(mockProvider)

	var c rosshook.Config
	c.MaxSharesToHold = 1
	c.LotSize = 1
	c.Balance = 1000
	var config strategies.Config
	b, _ := json.Marshal(c)
	json.Unmarshal(b, &config)

	placedOrders := make(chan *types.PlaceOrder)
	ordersStates := make(chan types.OrderExecutionState)

	strategy.Start(&config, &placedOrders, &ordersStates)

	mockedCandles := GetMock()
	for _, candle := range mockedCandles {
		candlesCh <- candle
	}

	select {
	case placedOrder, ok := <-placedOrders:
		if !ok {
			t.Fatalf("NO placed order")
		}
		if placedOrder.Quantity != 1 {
			t.Fatalf("Wrong placed order %v", placedOrder)
		}
		break

	}

}
