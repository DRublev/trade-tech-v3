package tinkoff

import (
	"context"
	"errors"
	"fmt"
	"log"
	sdk "main/integrations/tinkoff/sdk"
	"main/types"
	"os"
	"time"

	"github.com/tinkoff/invest-api-go-sdk/investgo"
	investapi "github.com/tinkoff/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const ENDPOINT = "sandbox-invest-public-api.tinkoff.ru:443"

// TODO: Хорошо бы явно наследовать types.Broker (чтоб были подсказки при имплементации метода)
type TinkoffBrokerPort struct{}

func (c *TinkoffBrokerPort) GetAccounts(ctx context.Context) ([]types.Account, error) {
	s, err := c.getSdk()
	if err != nil {
		fmt.Println("Cannot init sdk! ", err)
		return []types.Account{}, nil
	}

	us := s.NewUsersServiceClient()
	accountsRes, err := us.GetAccounts()

	if err != nil {
		fmt.Println(err)
		fmt.Println("Cannot get accounts ", err)
		return []types.Account{}, nil
	}
	accounts := []types.Account{}

	for _, acc := range accountsRes.Accounts {
		isOpen := acc.Status == investapi.AccountStatus_ACCOUNT_STATUS_OPEN
		hasAccess := acc.AccessLevel == investapi.AccessLevel_ACCOUNT_ACCESS_LEVEL_FULL_ACCESS
		isValidType := acc.Type == investapi.AccountType_ACCOUNT_TYPE_TINKOFF

		if isOpen && hasAccess && isValidType {
			accounts = append(accounts, types.Account{Id: acc.Id, Name: acc.Name})
		}
	}

	return accounts, nil
}

func (c *TinkoffBrokerPort) SetAccount(ctx context.Context, accountId string) error {
	return nil
}

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
