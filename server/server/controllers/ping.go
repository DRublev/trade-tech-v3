package controllers

import (
	"context"
	ping "main/server/contracts/contracts.ping"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var maxMissedPings = 6
var missedPings = 0
var pingMx = &sync.Mutex{}

var pingDuration = 10 * time.Second

func init() {
	go trackPings()
}

func (s *Server) Ping(ctx context.Context, in *ping.PingRequest) (*ping.PingResponse, error) {
	pingMx.Lock()
	missedPings = 0
	pingMx.Unlock()
	return &ping.PingResponse{}, nil
}

func trackPings() {
	for {
		select {
		case <-time.After(pingDuration):
			pingMx.Lock()
			missedPings++
			pingMx.Unlock()
			if missedPings >= maxMissedPings {
				log.Error("Max ping missed. Shutting down...")
				os.Exit(0)
			}
		}
	}
}
