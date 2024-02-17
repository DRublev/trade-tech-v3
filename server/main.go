package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	server "main/server"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	go server.Start(ctx, *port)

	<-ctx.Done()

	os.Exit(1)
}
