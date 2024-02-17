package tinkoff

type TinkoffBrokerPort struct{}

func (c *TinkoffBrokerPort) GetAccounts() ([]string, error) {
	return []string{}, nil
}
