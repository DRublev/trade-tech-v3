package server

	import (
	"context"
	"fmt"
	auth "main/grpcGW/grpcGW.auth"
	)

func (s *Server) SetToken(ctx context.Context, in *auth.SetTokenRequest) (*auth.SetTokenResponse, error) {
	fmt.Println("auth SetToken ", in)
	return &auth.SetTokenResponse{}, nil
}