package rosshook__tests

import (
	"main/bot/candles"
	"main/bot/strategies"
	"main/bot/strategies/rosshook"
	"main/types"
	"testing"
	"time"
)

type MockProvider struct {
}

func (p *MockProvider) GetOrCreate(instrumentID string, initialFrom time.Time, initialTo time.Time) (*chan types.OHLC, error) {
	return nil, nil
}
func Test(t *testing.T) {
	mockProvider := candles.Provider{}
	strategy := rosshook.New(&mockProvider)
	config := strategies.Config{}

	placedOrders := make(chan *types.PlaceOrder)
	ordersStates := make(chan types.OrderExecutionState)

	strategy.Start(&config, &placedOrders, &ordersStates)

}
