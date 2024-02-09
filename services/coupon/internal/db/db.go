package db

import (
	util_error "arvan-challenge/pkg/utils/errors"
	"arvan-challenge/services/coupon/pkg/env"
	"context"
	"fmt"
	"os"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"

	"github.com/rs/zerolog"
)

type DBHandler interface {
	CreateCoupon(c *Coupon, ctx *fasthttp.RequestCtx) (*Coupon, error)
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
	// defer dbpool.Close()

	return &dbHandler{
		db:  dbpool,
		cfg: cfg,
		l:   l,
	}
}

func (h *dbHandler) CreateCoupon(c *Coupon, ctx *fasthttp.RequestCtx) (*Coupon, error) {
	var coupon Coupon

	ct, err := h.GetCoupon(c.Code, ctx)

	if err != nil {
		h.l.Error().Msgf("An error occured while executing query: %v", err)
		return &coupon, util_error.NewInternalServerError("Somthing went wrong")
	} else if ct.Code == c.Code {
		return &coupon, util_error.NewBadRequestError("This code has already been defined")
	}

	err = h.db.QueryRow(ctx, `
			INSERT INTO coupons (coupon_id, code, expire_date, charge_amount, allowed_times)
			VALUES ($1, $2, $3, $4, $5) RETURNING *;`,
		c.Id, c.Code, c.ExpireDate, c.ChargeAmount, c.AllowedTimes).Scan(&coupon.Id, &coupon.Code, &coupon.ExpireDate, &coupon.ChargeAmount, &coupon.AllowedTimes)

	if err != nil {
		h.l.Error().Msgf("An error occured while executing query: %v", err)
		return &coupon, util_error.NewInternalServerError("Somthing went wrong")
	}

	return &coupon, nil
}

func (h *dbHandler) GetCoupon(code string, ctx context.Context) (*Coupon, error) {

	var coupon Coupon

	row, err := h.db.Query(ctx, `SELECT * FROM coupons WHERE code = $1`, code)

	if err != nil {
		h.l.Error().Msgf("An error occured while executing query: %v", err)
		return &coupon, util_error.NewInternalServerError("Somthing went wrong")
	}

	for row.Next() {
		err = row.Scan(&coupon.Id, &coupon.Code, &coupon.ExpireDate, &coupon.ChargeAmount, &coupon.AllowedTimes)

		if err != nil {
			h.l.Error().Msgf("An error occured while executing query: %v", err)
			return &coupon, util_error.NewInternalServerError("Somthing went wrong")
		}
	}

	return &coupon, nil

}

func (h *dbHandler) CloseConnection() {
	h.db.Close()
}
