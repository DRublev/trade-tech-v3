package server

import (
	"context"
	"fmt"
	accounts "main/server/contracts/contracts.accounts"
	auth "main/server/contracts/contracts.auth"
	marketdata "main/server/contracts/contracts.marketdata"
	shares "main/server/contracts/contracts.shares"
	trade "main/server/contracts/contracts.trade"
	"main/server/controllers"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start(ctx context.Context, port int) error {
	s := grpc.NewServer()

	log.Info("Starting server")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	srv := &controllers.Server{}

	accounts.RegisterAccountsServer(s, srv)
	auth.RegisterAuthServer(s, srv)
	marketdata.RegisterMarketDataServer(s, srv)
	shares.RegisterSharesServer(s, srv)
	trade.RegisterTradeServer(s, srv)

	log.Infof("Server listening at: %v", lis.Addr())
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Error listening to server %v", err)
		return err
	}

	go func() {
		<-ctx.Done()
		s.Stop()
		lis.Close()
	}()

	return nil
}
