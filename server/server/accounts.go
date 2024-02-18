package server

import (
	"context"
	"fmt"
	"main/bot"
	accounts "main/grpcGW/grpcGW.accounts"
	"main/types"
)

func (s *Server) GetAccounts(ctx context.Context, in *accounts.GetAccountsRequest) (*accounts.GetAccountsResponse, error) {
	err := bot.Init(ctx, types.Tinkoff)
	if err != nil {
		fmt.Println("accounts GetAccounst request err", err)
		return &accounts.GetAccountsResponse{Accounts: []*accounts.Account{}}, nil
	}

	var res []*accounts.Account = {}

	accs, err := bot.Broker.GetAccounts(ctx)
	if err != nil {
		return res, errors.New('error getting accounts' + err.Error())
	}

	for _, a := range accs {
		res = append(res, &accounts.Account{Id: a.Id, Name: a.Name})
	}

	fmt.Println("accounts GetAccounst request")
	return &accounts.GetAccountsResponse{Accounts: res}, nil
}