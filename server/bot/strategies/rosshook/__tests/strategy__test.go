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

func TestBuyTakeProfit(t *testing.T) {
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
			t.Fatalf("Не выставили заявку")
		}
		// TODO: Дописать проверку на корректность ордера
		// TODO: Дописать проверку на наличие продажи
		if placedOrder.Quantity != 1 {
			t.Fatalf("Ордер выставлен неверно %v", placedOrder)
		}
		break
	}
}

// TODO: Написать тест на сценарий покупка-стоп лосс
// TODO: Написать тест на сценарий, когда стратегия не должна отработать
