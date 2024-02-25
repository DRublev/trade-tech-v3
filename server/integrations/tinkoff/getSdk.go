package tinkoff

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/db"
	sdk "main/integrations/tinkoff/sdk"
	"main/utils"
	"os"
	"time"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var dbInstance = db.DB{}

func (c *TinkoffBrokerPort) getSdk() (*investgo.Client, error) {
	if sdk.IsInited() {
		return sdk.GetInstance(), nil
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

	s := sdk.Init(ctx, *config, logger)
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
