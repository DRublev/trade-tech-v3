package tinkoff

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

var instance *investgo.Client = nil
var once sync.Once

func Init(ctx context.Context, config investgo.Config, l investgo.Logger) *investgo.Client {
	once.Do(func() {
		inst, err := investgo.NewClient(ctx, config, l)
		if err != nil {
			log.Fatalln("Cannot init sdk!" + err.Error())
			return
		}
		instance = inst
		fmt.Println("Instance created", inst != nil)
	})

	return instance
}

func IsInited() bool {
	return instance != nil
}

func GetInstance() *investgo.Client {
	return instance
}
