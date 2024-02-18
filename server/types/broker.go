package types

import "context"

type BrokerKey string

const (
	Tinkoff BrokerKey = "tinkoff"
)

type Account struct {
	Id   string
	Name string
}

type Broker interface {
	GetAccounts(context.Context) ([]Account, error)
	SetAccount(context.Context, string) error
}
