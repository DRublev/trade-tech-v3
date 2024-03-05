package broker

import (
	"context"
	"errors"
	"main/integrations/tinkoff"
	"main/types"
)

var Broker types.IBroker

func Init(ctx context.Context, key types.BrokerKey) error {
	if Broker != nil {
		return nil
	}
	switch key {
	case types.Tinkoff:
		tinkoffBroker := &tinkoff.TinkoffBrokerPort{}

		Broker = tinkoffBroker

		_, err := tinkoffBroker.NewSdk()

		return err
	default:
		return errors.New("unknown broker type")
	}
}