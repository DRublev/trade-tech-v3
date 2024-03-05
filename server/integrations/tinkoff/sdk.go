package tinkoff

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/db"
	"main/utils"
	"os"
	"sync"
	"time"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var instance *investgo.Client = nil
var once sync.Once
var dbInstance = db.DB{}

func initialize(ctx context.Context, config investgo.Config, l investgo.Logger) *investgo.Client {
	once.Do(func() {
		inst, err := investgo.NewClient(ctx, config, l)
		if err != nil {
			log.Fatalln("Cannot init sdk!" + err.Error())
			return
		}
		instance = inst
		fmt.Println("Instance created", inst != nil)
	})

	return instance
}
func (c *TinkoffBrokerPort) GetSdk() (*investgo.Client, error) {
	if instance == nil {
		return nil, errors.New("sdk is not inited")
	}

	return instance, nil
}
<<<<<<< HEAD
func (c *TinkoffBrokerPort) NewSdk() (*investgo.Client, error) {
=======
func (c *TinkoffBrokerPort) NewSdk(accountId string) (*investgo.Client, error) {
>>>>>>> eb34c9e (feat: Новый метод в accounts.proto, контракт в go)
	if instance != nil {
		return instance, nil
	}

	if len(accountId) == 0 {
		return nil, errors.New("account id is empty")
	}

	token, err := getToken()
	if err != nil {
		return nil, err
	}
	config := &investgo.Config{
		EndPoint: ENDPOINT,
		Token:    token,
		// TODO: Для прод енвы кидать другое название
		AppName: "trade-tech-dev",
<<<<<<< HEAD
=======
		AccountId: accountId, //"2020306672"
>>>>>>> eb34c9e (feat: Новый метод в accounts.proto, контракт в go)
	}

	// TODO: Норм логгер надо бы
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"
	l, err := zapConfig.Build()
	logger := l.Sugar()

	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Printf(err.Error())
		}
	}()

	if err != nil {
		log.Fatalf("logger creating error %v", err)
	}

	ctx := context.Background()

	s := initialize(ctx, *config, logger)
	fmt.Println("Tinkoff sdk inited")

	return s, nil
}

func getToken() (string, error) {
	secret, exists := os.LookupEnv("SECRET")
	if !exists {
		return "", errors.New("secret not provided")
	}

	// TODO: Вынести бы в константу
	encrypted, err := dbInstance.Get([]string{"auth"})
	if err != nil {
		fmt.Println("No or missing token", err)
		return "", errors.New("no or missing token")
	}

	token, err := utils.Decrypt(string(encrypted), secret)
	if err != nil {
		fmt.Println("cannot decrypt token", err)
		dbInstance.Prune([]string{"auth"})
		return "", errors.New("cannot decrypt token")
	}

	return token, nil
}
