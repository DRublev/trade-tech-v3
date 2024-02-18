package bot

import (
	"context"
	"errors"
	"main/integrations/tinkoff"
	"main/types"
)

var Broker types.Broker

func Init(ctx context.Context, key types.BrokerKey) error {
	if Broker != nil {
		return nil
	}
	switch key {
	case types.Tinkoff:
		tinkoffBroker := &tinkoff.TinkoffBrokerPort{}
		Broker = tinkoffBroker
		return nil
	default:
		return errors.New("unknown broker type")
	}
}
