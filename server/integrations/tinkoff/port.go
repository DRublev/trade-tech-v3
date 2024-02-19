package tinkoff

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/db"
	sdk "main/integrations/tinkoff/sdk"
	"main/types"
	"main/utils"
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

var dbInstance = db.DB{}

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
		return "", errors.New("cannot decrypt token")
	}

	return token, nil
}
