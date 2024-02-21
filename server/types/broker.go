package types

import (
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

type Broker interface {
	GetAccounts() ([]Account, error)
	SetAccount(string) error
	GetCandles(string, Interval, time.Time, time.Time) ([]OHLC, error)
}
