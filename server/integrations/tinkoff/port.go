package tinkoff

import (
	"context"
	"fmt"
	"main/types"
)

const ENDPOINT = "sandbox-invest-public-api.tinkoff.ru:443"

// https://github.com/RussianInvestments/invest-api-go-sdk
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
		isOpen := acc.Status == 2                                 //pb.AccountStatus_ACCOUNT_STATUS_OPEN
		hasAccess := acc.AccessLevel == 1 || acc.AccessLevel == 2 //AccessLevel_ACCOUNT_ACCESS_LEVEL_FULL_ACCESS || AccessLevel_ACCOUNT_ACCESS_LEVEL_READ_ONLY
		isValidType := acc.Type == 1                
		fmt.Println(acc)              // pb.AccountType_ACCOUNT_TYPE_TINKOFF

		if isOpen && hasAccess && isValidType {
			accounts = append(accounts, types.Account{Id: acc.GetId(), Name: acc.GetName() })
		}
	}

	return accounts, nil
}

func (c *TinkoffBrokerPort) SetAccount(ctx context.Context, accountId string) error {
	return nil
}
