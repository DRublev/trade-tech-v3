package controllers

import (
	"context"
	"errors"
	"main/db"
	auth "main/server/contracts/contracts.auth"
	"main/utils"
	"os"

	log "github.com/sirupsen/logrus"
)

var authL = log.WithFields(log.Fields{
	"controller": "auth",
})

var dbInstance = db.DB{}

var authStore = []string{"auth"}

func (s *Server) SetToken(ctx context.Context, in *auth.SetTokenRequest) (*auth.SetTokenResponse, error) {
	authL.Info("SetToken requested")

	sercret, exists := os.LookupEnv("SECRET")
	if !exists {
		authL.Error("Missing secret key for token encryption")
		return &auth.SetTokenResponse{}, errors.New("missing key for encryption")
	}

	encrypted, err := utils.Encrypt(in.Token, sercret)
	if err != nil {
		authL.Errorf("Failed encrypting token: %v", err)
		return &auth.SetTokenResponse{}, errors.New("cannot encrypt token")
	}

	err = dbInstance.CreateStore(authStore)

	authL.Trace("Saving token to storage")
	err = dbInstance.Append(authStore, []byte(encrypted))

	authL.Infof("SetToken responding. Err: %v", err)
	return &auth.SetTokenResponse{}, err
}

func (s *Server) PruneTokens(ctx context.Context, in *auth.PruneTokensRequest) (*auth.PruneTokensResponse, error) {
	authL.Trace("Prune tokens")
	err := dbInstance.Prune(authStore)

	authL.Info("PruneToken responding")
	return &auth.PruneTokensResponse{}, err
}

func (s *Server) HasToken(ctx context.Context, in *auth.HasTokenRequest) (*auth.HasTokenResponse, error) {
	authL.Info("HasToken requesting")

	encrypted, err := dbInstance.Get(authStore)
	if err != nil && !os.IsNotExist(err) {
		authL.Errorf("Failed getting token from storage: %v", err)
		return &auth.HasTokenResponse{HasToken: false}, err
	}

	authL.Info("HasToken responding")
	return &auth.HasTokenResponse{HasToken: encrypted != nil}, nil
}
