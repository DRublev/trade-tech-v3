package server

import (
	"context"
	"fmt"
	"main/bot/broker"
	accounts "main/grpcGW/grpcGW.accounts"
	"main/types"
)

func (s *Server) GetAccounts(ctx context.Context, in *accounts.GetAccountsRequest) (*accounts.GetAccountsResponse, error) {
	err := broker.Init(ctx, types.Tinkoff)
	if err != nil {
		fmt.Println("accounts GetAccounst request err", err)
		return &accounts.GetAccountsResponse{Accounts: []*accounts.Account{}}, nil
	}

	var res []*accounts.Account

	accs, err := broker.Broker.GetAccounts()
	if err != nil {
		return &accounts.GetAccountsResponse{Accounts: res}, err
	}

	for _, a := range accs {
		res = append(res, &accounts.Account{Id: a.Id, Name: a.Name})
	}

	fmt.Println("accounts GetAccounst request")
	return &accounts.GetAccountsResponse{Accounts: res}, nil
}

// accountID, err := dbInstance.Get([]string{"accounts"})

func (s *Server) SetAccount