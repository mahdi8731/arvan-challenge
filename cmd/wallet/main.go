package main

import (
	"arvan-challenge/pkg/logger"
	"arvan-challenge/services/wallet/api"
	"arvan-challenge/services/wallet/pkg/env"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

var (
	cfg *env.Config
	l   *zerolog.Logger
	a   api.Api
)

func main() {
	cfg = env.ParseConfig()

	fmt.Println(os.Getwd())

	l = logger.NewLogger(cfg.LogLevel)

	// initialize api handler instance
	a = api.NewApi(l, cfg)

	a.Init()

	a.Run()
}
