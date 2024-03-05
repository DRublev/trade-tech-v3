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

<<<<<<< HEAD
		_, err := tinkoffBroker.NewSdk()

=======
		accountID, err := dbInstance.Get([]string{"accounts"})
		if err != nil {
			return err
		}

		_, err = tinkoffBroker.NewSdk(string(accountID))
>>>>>>> eb34c9e (feat: Новый метод в accounts.proto, контракт в go)
		return err
	default:
		return errors.New("unknown broker type")
	}
}