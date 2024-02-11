package main

import (
	"arvan-challenge/pkg/logger"
	"arvan-challenge/services/coupon/api"
	"arvan-challenge/services/coupon/pkg/env"
	"arvan-challenge/services/coupon/pkg/job"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
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

	jh := job.NewJob(cfg, l)
	defer jh.Close()

	// create a Scheduler
	s, _ := gocron.NewScheduler()
	defer func() { _ = s.Shutdown() }()

	_, _ = s.NewJob(
		gocron.DurationJob(
			time.Minute,
		),
		gocron.NewTask(jh.Do, context.Background()),
	)

	// start the scheduler
	s.Start()

	// initialize api handler instance
	a = api.NewApi(l, cfg)

	a.Init()

	a.Run()
}
