package job

import (
	"arvan-challenge/services/coupon/internal/db"
	"arvan-challenge/services/coupon/pkg/env"
	"context"
	"encoding/json"

	"github.com/rs/zerolog"

	"github.com/nats-io/nats.go"
)

type Job interface {
	Do(ctx context.Context) error
}

type job struct {
	cfg       *env.Config
	l         *zerolog.Logger
	dbHandler db.DBHandler
}

func NewJob(cfg *env.Config, l *zerolog.Logger) Job {
	dbHandler := db.NewDBHandler(cfg, l)
	return &job{
		cfg:       cfg,
		l:         l,
		dbHandler: dbHandler,
	}
}

func (job *job) Do(ctx context.Context) error {

	job.l.Info().Msg("Job started")

	messages, err := job.dbHandler.GetOutboxes(ctx)

	if err != nil {
		return err
	}

	var ids []int

	if len(*messages) > 0 {

		nc, err := nats.Connect("nats://localhost:4222")
		if err != nil {
			job.l.Error().Msgf("can not connect to nats: %v", err)
			return err
		}

		defer nc.Close()

		for _, v := range *messages {
			j, err := json.Marshal(v)
			if err != nil {
				job.l.Error().Msgf("ERROR when marshalling message: %v", err)
				continue
			}
			err = nc.Publish("wallet", j)
			if err != nil {
				job.l.Error().Msgf("ERROR when publishing message: %v", err)
				break
			}

			job.l.Info().Msgf("Published message: %v", v.Id)

			ids = append(ids, v.Id)

		}

		err = job.dbHandler.DeleteOutbox(&ids, ctx)

		if err != nil {
			job.l.Error().Msgf("can not delete outboxes: %v", err)
			return err
		}

	}

	return nil
}
