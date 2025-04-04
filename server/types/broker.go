package types

import (
	"context"
	"time"
)

// BrokerKey Ключ для обозначения брокера
// Влияет на сборку инстанса Broker
type BrokerKey string

const (
	// Tinkoff Тинькофф брокер
	// https://github.com/RussianInvestments/invest-api-go-sdk
	Tinkoff BrokerKey = "tinkoff"
)

// Account Аккаунт пользователя на стороне брокера
type Account struct {
	Id   string
	Name string
}

// IBroker Интерфейс, который должен имплементировать каждая интеграция с брокером
type IBroker interface {
	GetAccounts() ([]Account, error)
	SetAccount(string) error
	GetCandles(instrumentID string, interval Interval, start time.Time, end time.Time) ([]OHLC, error)
	SubscribeCandles(context.Context, *chan OHLC, string, Interval, bool) error
	SubscribeOrderbook(context.Context, *chan *Orderbook, string, int32) error
	GetShares(InstrumentStatus) ([]Share, error)
	PlaceOrder(order *PlaceOrder) (OrderID, error)
	SubscribeOrders(func(OrderExecutionState)) error
	GetTradingSchedules(string, time.Time, time.Time) ([]TradingSchedule, error)
	GetOrderState(orderID OrderID) (OrderExecutionState, error)
}
