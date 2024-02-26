package types

import (
	"context"
	"time"
)

type BrokerKey string

const (
	Tinkoff BrokerKey = "tinkoff"
)

type Account struct {
	Id   string
	Name string
}

type IBroker interface {
	GetAccounts() ([]Account, error)
	SetAccount(string) error
	GetCandles(string, Interval, time.Time, time.Time) ([]OHLC, error)
	SubscribeCandles(context.Context, *chan OHLC, string, Interval) error
	PlaceOrder(order Order) (string, error)
}
