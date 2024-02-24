package server

import (
	"context"
	"fmt"
	"log"
	accounts "main/grpcGW/grpcGW.accounts"
	auth "main/grpcGW/grpcGW.auth"
	marketdata "main/grpcGW/grpcGW.marketdata"
	instruments "main/grpcGW/grpcGW.instruments"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	accounts.UnimplementedAccountsServer
	auth.UnimplementedAuthServer
	marketdata.UnimplementedMarketDataServer
	instruments.UnimplementedSharesServer
}

func Start(ctx context.Context, port int) {
	s := grpc.NewServer()
	fmt.Println("Starting server")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error listening port", err)
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	srv := &Server{}

	accounts.RegisterAccountsServer(s, srv)
	auth.RegisterAuthServer(s, srv)
	marketdata.RegisterMarketDataServer(s, srv)
	instruments.RegisterSharesServer(s, srv)

	fmt.Println("Starting server", lis.Addr())
	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("Error listening to server ", err)
	}

	<-ctx.Done()
}
