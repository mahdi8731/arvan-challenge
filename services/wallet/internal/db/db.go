package db

import (
	util_error "arvan-challenge/pkg/utils/errors"
	"arvan-challenge/services/wallet/pkg/env"
	"context"
	"fmt"
	"os"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rs/zerolog"
)

type DBHandler interface {
	AddTransaction(t *Transaction, w *Wallet, ctx context.Context) (*Wallet, error)
	CloseConnection()
}

type dbHandler struct {
	cfg *env.Config
	l   *zerolog.Logger
	db  *pgxpool.Pool
}

func NewDBHandler(cfg *env.Config, l *zerolog.Logger) DBHandler {

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	// Connect to database
	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		l.Fatal().Msgf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	dbconfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		l.Fatal().Msgf("Unable to parse config: %v\n", err)
		os.Exit(1)
	}
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	return &dbHandler{
		db:  dbpool,
		cfg: cfg,
		l:   l,
	}
}

func (h *dbHandler) AddTransaction(t *Transaction, w *Wallet, ctx context.Context) (*Wallet, error) {
	var wallet Wallet

	tx, err := h.db.Begin(ctx)

	fmt.Print("#################")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	err = tx.QueryRow(ctx, `insert into wallet (wallet_id, last_modified, phone_number, inventory) values ($1, $2, $3, $4)
	on CONFLICT(phone_number) DO update set inventory = wallet.inventory + $4 , last_modified = $2 where wallet.phone_number = $3 RETURNING *`,
		w.Id, w.LastModefied, w.PhoneNumber, t.Amount).Scan(&wallet.Id, &wallet.PhoneNumber, &wallet.LastModefied, &wallet.Inventory)
	if err != nil {
		h.l.Error().Msgf("An error occured while executing insert query: %v", err)
		return nil, util_error.NewInternalServerError("Somthing went wrong")
	}

	_, err = tx.Exec(ctx, `
			INSERT INTO transactions (id, description, date, amount, wallet_id)
			VALUES ($1, $2, $3, $4, $5)`,
		t.Id, t.Description, t.Date, t.Amount, wallet.Id)

	if err != nil {
		h.l.Error().Msgf("An error occured while executing update query: %v", err)
		return nil, util_error.NewInternalServerError("Somthing went wrong")
	}

	return &wallet, nil
}

func (h *dbHandler) CloseConnection() {
	h.db.Close()
}
