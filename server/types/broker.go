package types

type Broker interface {
	GetAccounts() ([]string, error)
	SetAccount(string) error
}
