package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	server "main/server"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/magnetde/loki"
)

var (
	port = flag.Int("port", 50051, "The server port")
	logsAddress = flag.String("logsAddress", "http://79.174.80.98:3100", "The server port")
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.TraceLevel)
}

func main() {
	flag.Parse()

	hook := loki.NewHook(*logsAddress, loki.WithName("trade-tech"), loki.WithLabel("env", "dev"), loki.WithLabel("app", "server"), loki.WithLevel(log.InfoLevel))
	defer hook.Close()
	log.AddHook(hook)

	if env, ok := os.LookupEnv("ENV"); !ok || env != "PROD" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Cannot load env!")
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	go server.Start(ctx, *port)

	<-ctx.Done()

	os.Exit(1)
}
