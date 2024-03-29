package subscriber

import (
	"arvan-challenge/services/wallet/internal/db"
	"arvan-challenge/services/wallet/pkg/env"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type Subscriber interface {
	Subscribe(ctx *context.Context)
}

type subscriber struct {
	cfg       *env.Config
	l         *zerolog.Logger
	dbHandler db.DBHandler
}

func NewSubscriber(cfg *env.Config, l *zerolog.Logger) Subscriber {
	dbHandler := db.NewDBHandler(cfg, l)
	return &subscriber{
		cfg:       cfg,
		l:         l,
		dbHandler: dbHandler,
	}
}

func (s *subscriber) Subscribe(ctx *context.Context) {

	// Use default server (localhost:4222) or specify custom address
	nc, err := nats.Connect(s.cfg.NATSUrl)
	if err != nil {
		log.Fatal(err)
	}
	// defer nc.Close()

	// Subscribe to a specific topic
	_, err = nc.Subscribe("wallet", func(m *nats.Msg) {
		s.l.Info().Msgf("Received message: %s", string(m.Data))

		var msg Message

		err := json.Unmarshal(m.Data, &msg)

		if err != nil {
			s.l.Fatal().Msgf("Error decoding message: %v", err)
		}

		tid, _ := uuid.NewV7()
		wid, _ := uuid.NewV7()
		tt := time.Now().UTC()

		transaction := &db.Transaction{
			Id:          tid,
			Date:        tt,
			Description: "From Code",
			Amount:      msg.Amount,
		}

		wallet := &db.Wallet{
			Id:           wid,
			PhoneNumber:  msg.PhoneNumber,
			LastModefied: tt,
		}

		s.dbHandler.AddTransaction(transaction, wallet, context.Background())
	})
	if err != nil {
		log.Fatal(err)
	}
	// defer sub.Unsubscribe()

}
