package tinkoff

import (
	"context"
	"errors"
	"main/configuration"
	"main/db"
	"main/utils"
	"os"
	"sync"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	log "github.com/sirupsen/logrus"
)

var instance *investgo.Client = nil
var once sync.Once
var dbInstance = db.DB{}

var sdkL = log.WithFields(log.Fields{
	"broker": "Tinkoff",
})

var endpoint string

func initialize(ctx context.Context, config investgo.Config, l investgo.Logger) *investgo.Client {
	once.Do(func() {
		sdkL.Trace("Initializing Tinkoff sdk")
		inst, err := investgo.NewClient(ctx, config, l)
		if err != nil {
			sdkL.Errorf("Cannot init Tinkoff sdk! %v", err)
			return
		}
		instance = inst
		sdkL.Trace("Tinkoff sdk created")
	})

	return instance
}
func (c *TinkoffBrokerPort) GetSdk() (*investgo.Client, error) {
	if instance == nil {
		sdkL.Error("Tinkoff sdk is not inited, but requested")
		return nil, errors.New("sdk is not inited")
	}

	return instance, nil
}

func (c *TinkoffBrokerPort) NewSdk(accountID *string) (*investgo.Client, error) {
	sdkL.WithField("is accountId empty", accountID == nil).Trace("NewSdk called")
	if instance != nil {
		return instance, nil
	}

	token, err := getToken()
	if err != nil {
		sdkL.Errorf("Cannot get token: %v", err)
		return nil, err
	}

	conf := configuration.Configuration{
		TinkoffEndpoint: "invest-public-api.tinkoff.ru:443",
	}
	endpoint = conf.Get().TinkoffEndpoint

	config := &investgo.Config{
		EndPoint: endpoint,
		Token:    token,
		// TODO: Для прод енвы кидать другое название
		AppName: "trade-tech-dev",
	}
	if accountID != nil {
		config.AccountId = *accountID
	}

	logger := log.New()
	logger.Out = os.Stdout

	if err != nil {
		log.Errorf("logger creating error %v", err)
	}

	ctx := context.Background()

	s := initialize(ctx, *config, logger)

	sdkL.Trace("NewSdk returning ok")
	return s, nil
}

func getToken() (string, error) {
	sdkL.Trace("Getting token")
	secret, exists := os.LookupEnv("SECRET")
	if !exists {
		sdkL.Error("Missing sekret key for token encryption")
		return "", errors.New("secret not provided")
	}

	sdkL.Trace("Getting token from storage")
	// TODO: Вынести бы в константу
	encrypted, err := dbInstance.Get([]string{"auth"})
	if err != nil {
		sdkL.Errorf("Cannot getting token from storage: %v", err)
		return "", errors.New("no or missing token")
	}

	sdkL.Trace("Decrypting token")
	token, err := utils.Decrypt(string(encrypted), secret)
	if err != nil {
		sdkL.Errorf("Cannot decrypt token. Clearing token in storage. Error: %v", err)
		dbInstance.Prune([]string{"auth"})
		return "", errors.New("cannot decrypt token")
	}

	sdkL.Trace("Returning token")
	return token, nil
}
