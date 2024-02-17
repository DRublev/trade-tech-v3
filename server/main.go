package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	server "main/server"

	"github.com/joho/godotenv"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatalln("Cannot load env!")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	go server.Start(ctx, *port)

	<-ctx.Done()

	os.Exit(1)
}
