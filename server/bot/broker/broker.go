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

// Broker Инстанс брокера
var Broker types.IBroker

// Init Конструктор для брокера. Создает инстанс брокера, исходя из ключа
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
				fmt.Printf("34 broker %v\n", err)
				return err
			}
			_, err = tinkoffBroker.NewSdk(&accountID)
			return nil
		}

		accountID = string(accountIDBytes)
		accountID = accountID[:len(accountID)-len("\n")]

		_, err = tinkoffBroker.NewSdk(&accountID)

		return err
	default:
		return errors.New("unknown broker type")
	}
}
