package main

import "C"
import (
	"context"
	"flag"
	"os"
	"os/signal"

	"main/configuration"
	"main/identity"
	server "main/server"

	"github.com/joho/godotenv"
	"github.com/magnetde/loki"
	log "github.com/sirupsen/logrus"
)

var (
	port         = flag.Int("port", 50051, "The server port")
	logsAddress  = flag.String("logsAddress", "http://87.242.100.16:3100", "The server port")
	yamlConfPath = flag.String("conf", "./config.yaml", "Path for config")
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

var envFromBuild string = "CI"
var secretFromBuild string = "trade-tech-secret-for-encryption"

var version = "24.08.1"

func prepareFlagsAndEnv() {
	flag.Parse()

	env, ok := os.LookupEnv("ENV")
	if !ok {
		env = envFromBuild
	}

	if len(env) == 0 || env == "DEV" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Cannot load env!")
		}
	}

	if env == "PROD" {
		uid := identity.GetId()
		hook := loki.NewHook(
			*logsAddress,
			loki.WithName("trade-tech"),
			loki.WithLabel("app", "server"),
			loki.WithLabel("uid", uid),
			loki.WithLabel("version", version),
			loki.WithLevel(log.TraceLevel),
			loki.WithLabelsEnabled(loki.LevelLabel, loki.FieldsLabel, loki.MessageLabel),
		)
		defer hook.Close()
		log.AddHook(hook)
	}

	if len(secretFromBuild) > 0 {
		os.Setenv("SECRET", secretFromBuild)
	}

	conf := configuration.Configuration{
		TinkoffEndpoint: "invest-public-api.tinkoff.ru:443",
	}
	conf.Load(*yamlConfPath)
}

func main() {
	prepareFlagsAndEnv()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go LaunchServer()

	<-ctx.Done()
	os.Exit(1)

}

//export LaunchServer
func LaunchServer() int {
	prepareFlagsAndEnv()

	go func() {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()
		err := server.Start(ctx, *port)

		if err != nil {
			log.Fatalf("Error starting server: %v", err)
			os.Exit(1)
			return
		}
		<-ctx.Done()

		os.Exit(1)
	}()

	return *port
}
