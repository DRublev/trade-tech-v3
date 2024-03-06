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

<<<<<<< HEAD
<<<<<<< HEAD
		_, err := tinkoffBroker.NewSdk()

=======
		accountID, err := dbInstance.Get([]string{"accounts"})
		if err != nil {
			return err
		}

		_, err = tinkoffBroker.NewSdk(string(accountID))
>>>>>>> eb34c9e (feat: Новый метод в accounts.proto, контракт в go)
=======
		_, err := tinkoffBroker.NewSdk()

>>>>>>> 8fc68c3 (fix: Получение айдишника аккаунта и инит сдк)
		return err
	default:
		return errors.New("unknown broker type")
	}
}