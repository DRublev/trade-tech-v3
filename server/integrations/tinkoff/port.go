package tinkoff

import (
	"context"
	"fmt"
	"main/types"

	investapi "github.com/tinkoff/invest-api-go-sdk/proto"
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

// func (c *TinkoffBrokerPort) GetCandles() {
// 	s, err := c.getSdk()
// 	if err != nil {
// 		fmt.Println("Cannot init sdk! ", err)
// 		return []types.Account{}, nil
// 	}
// }
