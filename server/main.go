package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	test "main/proto/test"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	test.UnimplementedTestServer
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	test.RegisterTestServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error listening to server ", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	<-ctx.Done()

	os.Exit(1)
}

func (s *server) Ping(ctx context.Context, in *test.PingRequest) (*test.PingResponse, error) {
	fmt.Println("test Ping from", in.Content)
	return &test.PingResponse{Content: "hi test from server"}, nil
}
