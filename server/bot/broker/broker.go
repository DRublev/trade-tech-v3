package broker

import (
	"context"
	"errors"
	"fmt"
	"main/db"
	"main/integrations/tinkoff"
	"main/types"
	"os"
)

var Broker types.IBroker

func Init(ctx context.Context, key types.BrokerKey) error {
	if Broker != nil {
		return nil
	}
	switch key {
	case types.Tinkoff:
		dbInstance := &db.DB{}
		tinkoffBroker := &tinkoff.TinkoffBrokerPort{}

		Broker = tinkoffBroker

		var accountID string
		var err error

		accountIDBytes, err := dbInstance.Get([]string{"accounts"})
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			_, err = tinkoffBroker.NewSdk(&accountID)
			return err
		}

		accountID = string(accountIDBytes)
		fmt.Printf("30 broker %v\n", accountID)
		accountID = accountID[:len(accountID)-len("\n")]

		_, err = tinkoffBroker.NewSdk(&accountID)
		return err
	default:
		return errors.New("unknown broker type")
	}
}
