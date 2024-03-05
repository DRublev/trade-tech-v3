package broker

import (
	"context"
	"errors"
	"main/db"
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
		dbInstance := db.DB{}
		tinkoffBroker := &tinkoff.TinkoffBrokerPort{}

		Broker = tinkoffBroker

		accountID, err := dbInstance.Get([]string{"accounts"})
		if err != nil {
			return err
		}

		_, err = tinkoffBroker.NewSdk(string(accountID))
		return err
	default:
		return errors.New("unknown broker type")
	}
}