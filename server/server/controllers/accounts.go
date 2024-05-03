package controllers

import (
	"context"
	"errors"
	"main/bot/broker"
	accounts "main/server/contracts/contracts.accounts"
	"main/types"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

var accL = log.WithFields(log.Fields{
	"controller": "accounts",
})

func (s *Server) GetAccounts(ctx context.Context, in *accounts.GetAccountsRequest) (*accounts.GetAccountsResponse, error) {
	accL.Info("GetAccounts requested")

	err := broker.Init(ctx, types.Tinkoff)
	if err != nil {
		accL.Errorf("Cannot init broker: %v", err)
		return &accounts.GetAccountsResponse{Accounts: []*accounts.Account{}}, nil
	}

	var res []*accounts.Account

	accs, err := broker.Broker.GetAccounts()
	if err != nil {
		accL.Errorf("Failed getting accounts from broker: %v", err)
		return &accounts.GetAccountsResponse{Accounts: res}, err
	}

	accL.Tracef("Found %v accounts", len(accs))
	for _, a := range accs {
		res = append(res, &accounts.Account{Id: a.Id, Name: a.Name})
	}

	accL.Info("GetAccounts responding")
	return &accounts.GetAccountsResponse{Accounts: res}, nil
}

func (s *Server) SetAccount(ctx context.Context, in *accounts.SetAccountRequest) (*accounts.SetAccountResponse, error) {
	accL.Info("SetAccount requested")

	if in.AccountId == "" {
		return &accounts.SetAccountResponse{}, errors.New("accountId is empty")
	}

	content := []byte(in.AccountId)

	accL.Trace("Setting account id to storage")
	dbInstance.Prune([]string{"accounts"})
	err := dbInstance.Append([]string{"accounts"}, content)

	accL.Info("SetAccount responding")
	return &accounts.SetAccountResponse{}, err
}

func (s *Server) GetAccount(ctx context.Context, in *accounts.GetAccountRequest) (*accounts.GetAccountResponse, error) {
	accL.Info("GetAccount requested")

	accL.Trace("Getting account id")
	accountIDRaw, err := dbInstance.Get([]string{"accounts"})

	if os.IsNotExist(err) {
		return &accounts.GetAccountResponse{}, nil
	}

	accID := strings.TrimRight(string(accountIDRaw), "\r\n")
	accL.Infof("GetAccount responding: %v", accID)
	return &accounts.GetAccountResponse{AccountId: accID}, err
}
