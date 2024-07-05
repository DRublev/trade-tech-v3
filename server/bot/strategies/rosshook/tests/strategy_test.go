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
	mock    []types.OHLC
	channel chan types.OHLC
}

func (p MockProvider) GetOrCreate(instrumentID string, initialFrom time.Time, initialTo time.Time, onlyCompletedCandles bool) (*chan types.OHLC, error) {
	p.channel = make(chan types.OHLC)
	go func() {
		for _, candle := range p.mock {
			// Таймаут, чтобы успела отработать логи обновления стейта при выставлении заявки
			<-time.After(0.2 * 1000 * 1000 * 1000)
			p.channel <- candle
		}
	}()
	return &p.channel, nil
}

// Выставляем бай
func TestBuyOrderPlaced(t *testing.T) {
	mockProvider := MockProvider{
		mock: getShouldBuyMock(),
	}

	strategy := rosshook.New(mockProvider, strategies.ActivityContainer{})

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

	select {
	case placedOrder, ok := <-placedOrders:
		if !ok {
			t.Fatalf("Не выставили заявку")
		}
		// TODO: Дописать проверку на корректность ордера
		if placedOrder.Quantity != 1 {
			t.Fatalf("Ордер выставлен неверно %v", placedOrder)
		}
		return
	case <-time.After(time.Second * 10):
		t.Fatalf("Таймаут выставления заявки")
		break
	}

	t.Fatal("Не выставили заявку")
}

func TestShouldCloseBuyIfNotExecuted(t *testing.T) {
	mockProvider := MockProvider{
		mock: getShouldCloseBuyWhenNotExecutedMock(),
	}

	strategy := rosshook.New(mockProvider, strategies.ActivityContainer{})

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

	var closedOrderID string
	go func() {
		for {
			select {
			case placedOrder := <-placedOrders:
				// Нам нужно проверить, что мы закрываем неисполнившиеся заявки на покупку
				if len(placedOrder.CancelOrder) > 0 {
					closedOrderID = string(placedOrder.CancelOrder)
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
	}()
	// Без этого, тест будет висеть дефолтный таймаут (30 секунд), пока не упадет сам
	<-time.After(time.Second * 15)
	if len(closedOrderID) <= 0 {
		t.Fatalf("Не закрыли заявку")
	}
}

func TestBuyAndStopLoss(t *testing.T) {
	mockProvider := MockProvider{
		mock: getBuyAndStopLossMock(),
	}

	strategy := rosshook.New(mockProvider, strategies.ActivityContainer{})

	var c rosshook.Config
	c.MaxSharesToHold = 1
	c.LotSize = 1
	c.Balance = 1000
	c.SaveProfit = 0.2
	c.StopLoss = 0.2
	var config strategies.Config
	b, _ := json.Marshal(c)
	json.Unmarshal(b, &config)

	placedOrders := make(chan *types.PlaceOrder)
	ordersStates := make(chan types.OrderExecutionState)

	go strategy.Start(&config, &placedOrders, &ordersStates)

	go func(t *testing.T) {
		// Без этого, тест будет висеть дефолтный таймаут (30 секунд), пока не упадет сам
		<-time.After(time.Second * 20)
		t.Fatal("Таймаут")
	}(t)

	shouldBeBuyOrder := <-placedOrders
	if shouldBeBuyOrder.Direction == types.Buy {
		// Эмулируем выставления заявки на покупку, но она будет висеть не исполнившаяся
		ordersStates <- types.OrderExecutionState{
			Status:             types.New,
			LotsExecuted:       0,
			LotsRequested:      int(shouldBeBuyOrder.Quantity),
			ExecutedOrderPrice: float64(shouldBeBuyOrder.Price),
			InstrumentID:       shouldBeBuyOrder.InstrumentID,
			Direction:          shouldBeBuyOrder.Direction,
			ID:                 "placedBuyOrderID",
		}
		ordersStates <- types.OrderExecutionState{
			Status:             types.Fill,
			LotsExecuted:       int(shouldBeBuyOrder.Quantity),
			LotsRequested:      int(shouldBeBuyOrder.Quantity),
			ExecutedOrderPrice: float64(shouldBeBuyOrder.Price),
			InstrumentID:       shouldBeBuyOrder.InstrumentID,
			Direction:          shouldBeBuyOrder.Direction,
			ID:                 "placedBuyOrderID",
		}
		// По цене targetGrow
		if shouldBeBuyOrder.Price != 527.2 {
			t.Fatalf("Заход по уебищной цене: %v", shouldBeBuyOrder.Price)
		}
	}

	shouldBeTakeOrder := <-placedOrders
	// 524.8 - цена первой свечи, которая пробила стоп
	if shouldBeTakeOrder.Price != 524.8 {
		t.Fatalf("Неверный стоп-лосс %v", shouldBeTakeOrder.Price)
	}
}

// TODO: Тест что корректно выставляется закрытие пендинг бай ордеров
// TODO: Тест, когда LotSize > 1
