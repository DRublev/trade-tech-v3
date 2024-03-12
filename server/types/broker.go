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
	GetCandles(string, Interval, time.Time, time.Time) ([]OHLC, error)
	SubscribeCandles(context.Context, *chan OHLC, string, Interval) error
	SubscribeOrderbook(context.Context, *chan *Orderbook, string, int32) error
	GetShares(InstrumentStatus) ([]Share, error)
	PlaceOrder(order *PlaceOrder) (OrderID, error)
	SubscribeOrders(func(OrderExecutionState)) error
	GetOrderState(orderID OrderID) (OrderExecutionState, error)
}
