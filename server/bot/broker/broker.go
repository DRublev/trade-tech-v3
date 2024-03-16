package broker

import (
	"context"
	"errors"
	"main/db"
	"main/integrations/tinkoff"
	"main/types"
	"os"

	log "github.com/sirupsen/logrus"
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
		log.Info("Setting Tinkoff as a Broker")
		dbInstance := &db.DB{}
		tinkoffBroker := &tinkoff.TinkoffBrokerPort{}

		Broker = tinkoffBroker

		var accountID string
		var err error

		log.Trace("Getting account ID")
		accountIDBytes, err := dbInstance.Get([]string{"accounts"})
		if err != nil {
			if !os.IsNotExist(err) {
				log.Errorf("Failed getting accountID from database: %v", err)
				return err
			}
			log.Trace("No account ID stored, creating with empty")
			_, err = tinkoffBroker.NewSdk(&accountID)
			return nil
		}

		accountID = string(accountIDBytes)
		accountID = accountID[:len(accountID)-len("\n")]

		log.Trace("Account ID has been found in store, creating Tinkoff SDK with it")
		_, err = tinkoffBroker.NewSdk(&accountID)

		return err
	default:
		log.Errorf("Trying to instantiate unsupported broker: %v", key)
		return errors.New("unknown broker type")
	}
}
