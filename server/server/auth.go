package server

import (
	"context"
	"errors"
	"fmt"
	"main/db"
	auth "main/grpcGW/grpcGW.auth"
	"main/utils"
	"os"
)

var dbInstance = db.DB{}

func (s *Server) SetToken(ctx context.Context, in *auth.SetTokenRequest) (*auth.SetTokenResponse, error) {
	fmt.Println("auth SetToken ", in)

	sercret, exists := os.LookupEnv("SECRET")
	if !exists {
		return &auth.SetTokenResponse{}, errors.New("missing key for encryption")
	}

	encrypted, err := utils.Encrypt(in.Token, sercret)
	if err != nil {
		fmt.Println("Cannot encrypt token", err)
		return &auth.SetTokenResponse{}, errors.New("cannot encrypt token")
	}

	err = dbInstance.Append([]string{"auth"}, []byte(encrypted))

	return &auth.SetTokenResponse{}, err
}
