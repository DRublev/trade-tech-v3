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

	mockedCandles := getShouldBuyMock()
	go func() {
		for _, candle := range mockedCandles {
			candlesCh <- candle
		}
	}()

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

func TestShouldNotBuy(t *testing.T) {
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

	mockedCandles := getShouldNotBuyMock()
	go func() {
		for _, candle := range mockedCandles {
			candlesCh <- candle
		}
	}()

	select {
	case placedOrder := <-placedOrders:
		if placedOrder.Direction == 1 {
			t.Fatalf("Не должны, но купили по %v", placedOrder.Price)
		}
		break
	}
}

func TestShouldCloseBuyIfNotExecuted(t *testing.T) {
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

	mockedCandles := getShouldCloseBuyWhenNotExecutedMock()
	go func() {
		for _, candle := range mockedCandles {
			// Таймаут, чтобы успела отработать логи обновления стейта при выставлении заявки
			<-time.After(0.1 * 1000 * 1000 * 1000)

			candlesCh <- candle
		}
	}()

	var closedOrderID string
	for {
		select {
		case placedOrder := <-placedOrders:
			// Нам нужно проверить, что мы закрываем неисполнившиеся заявки на покупку
			if len(placedOrder.CancelOrder) > 0 {
				closedOrderID = string(placedOrder.CancelOrder)
				return
			}
			if placedOrder.Direction == types.Buy {
				// Эмулируем выставления заявки на покупку, но она будет висеть не исполнившаяся
				ordersStates <- types.OrderExecutionState{
					Status:             types.New,
					LotsExecuted:       0,
					LotsRequested:      int(placedOrder.Quantity),
					ExecutedOrderPrice: float64(placedOrder.Price),
					InstrumentID:       placedOrder.InstrumentID,
					Direction:          placedOrder.Direction,
					ID:                 "placedBuyOrderID",
				}
			}
		}
	}

	// Падаем, если через 5 секунд мы не закрыли заявку
	<-time.After(5 * 1000 * 1000 * 1000)
	// Без этого, тест будет висеть дефолтный таймаут (30 секунд), пока не упадет сам
	if len(closedOrderID) <= 0 {
		t.Fatalf("Не закрыли заявку")
	}
}

// TODO: Написать тест на сценарий покупка-стоп лосс
