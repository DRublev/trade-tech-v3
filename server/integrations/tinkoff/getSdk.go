package tinkoff

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	sdk "main/integrations/tinkoff/sdk"

	"github.com/tinkoff/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (c *TinkoffBrokerPort) getSdk() (*investgo.Client, error) {
	if sdk.IsInited() {
		return sdk.GetInstance(), nil
	}
	// TODO: Придумать как нормально брать токен
	token, ok := os.LookupEnv("TINKOFF_TOKEN_RO")
	if !ok {
		return nil, errors.New("no token provided")
	}
	config := &investgo.Config{
		EndPoint: ENDPOINT,
		Token:    token,
		// TODO: Для прод енвы кидать другое название
		AppName: "trade-tech-dev",
	}

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
