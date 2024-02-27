package server

import (
	"context"
	"fmt"
	"log"
	accounts "main/grpcGW/grpcGW.accounts"
	auth "main/grpcGW/grpcGW.auth"
	marketdata "main/grpcGW/grpcGW.marketdata"
	shares "main/grpcGW/grpcGW.shares"
	trade "main/grpcGW/grpcGW.trade"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	accounts.UnimplementedAccountsServer
	auth.UnimplementedAuthServer
	marketdata.UnimplementedMarketDataServer
	shares.UnimplementedSharesServer
	trade.UnimplementedTradeServer
}

func Start(ctx context.Context, port int) {
	s := grpc.NewServer()
	defer s.Stop()
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
	shares.RegisterSharesServer(s, srv)
	trade.RegisterTradeServer(s, srv)

	fmt.Println("Starting server", lis.Addr())
	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("Error listening to server ", err)
	}

	<-ctx.Done()
}
