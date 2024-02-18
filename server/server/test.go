package server

import (
	"context"
	"fmt"
	test "main/grpcGW/grpcGW.test"
)

func (s *Server) Ping(ctx context.Context, in *test.PingRequest) (*test.PingResponse, error) {
	fmt.Println("test Ping from", in.Content)
	return &test.PingResponse{Content: "hi test from server"}, nil
}
