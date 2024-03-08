package broker

import (
	"context"
	"errors"
	"fmt"
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
		dbInstance := &db.DB{}
		tinkoffBroker := &tinkoff.TinkoffBrokerPort{}

		Broker = tinkoffBroker

		accountIDBytes, err := dbInstance.Get([]string{"accounts"})
		if err != nil {
			return err
		}
		accountID := string(accountIDBytes)
		fmt.Printf("30 broker %v\n", accountID)
		accountID = accountID[:len(accountID) - len("\n")]
fmt.Printf("28 broker %v\n", accountID)
		_, err = tinkoffBroker.NewSdk(accountID)
		return err
	default:
		return errors.New("unknown broker type")
	}
}