package db

import (
	util_error "arvan-challenge/pkg/utils/errors"
	"arvan-challenge/services/coupon/pkg/env"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rs/zerolog"
)

type DBHandler interface {
	CreateCoupon(c *Coupon, ctx context.Context) (*Coupon, error)
	GetCoupon(code string, ctx context.Context) (*Coupon, error)
	GetUsersByCoupon(code string, ctx context.Context) (*[]string, error)
	UseCoupon(code, phone_number string, ctx context.Context) (*Coupon, error)
	GetOutboxes(ctx context.Context) (*[]Outbox, error)
	DeleteOutbox(ids *[]int, ctx context.Context) error
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
	// defer dbpool.Close()

	return &dbHandler{
		db:  dbpool,
		cfg: cfg,
		l:   l,
	}
}

func (h *dbHandler) CreateCoupon(c *Coupon, ctx context.Context) (*Coupon, error) {
	var coupon Coupon

	ct, err := h.GetCoupon(c.Code, ctx)

	if err != nil {
		switch t := err.(type) {
		case *util_error.InternalServerError:
			h.l.Error().Msgf("An error occured while executing query: %v", t)
			return &coupon, util_error.NewInternalServerError("Somthing went wrong")
		case *util_error.BadRequestError:
			break
		}

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

	if row.Next() {
		err = row.Scan(&coupon.Id, &coupon.Code, &coupon.ExpireDate, &coupon.ChargeAmount, &coupon.AllowedTimes)

		if err != nil {
			h.l.Error().Msgf("An error occured while executing query: %v", err)
			return &coupon, util_error.NewInternalServerError("Somthing went wrong")
		}
	} else {
		return nil, util_error.NewBadRequestError("This code is not a valid coupon")
	}

	return &coupon, nil

}

func (h *dbHandler) UseCoupon(code, phone_number string, ctx context.Context) (*Coupon, error) {

	ct, err := h.GetCoupon(code, ctx)

	if err != nil {
		h.l.Error().Msgf("An error occured while executing query: %v", err)
		return nil, err
	}

	tx, err := h.db.Begin(ctx)
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

	_, err = tx.Exec(ctx, "update coupons set allowed_times = allowed_times -1 where code = $1 and allowed_times > 0", ct.Code)
	if err != nil {
		h.l.Error().Msgf("An error occured while executing update query: %v", err)
		return nil, util_error.NewInternalServerError("Somthing went wrong")
	}

	id, _ := uuid.NewV7()

	_, err = tx.Exec(ctx, "INSERT INTO couponsـused (id, phone_number, coupon_id) VALUES ($1, $2, $3)", id, phone_number, ct.Id)
	if err != nil {
		h.l.Error().Msgf("An error occured while executing insert query: %v", err)
		return nil, util_error.NewInternalServerError("Somthing went wrong")
	}

	_, err = tx.Exec(ctx, "INSERT INTO outbox (phone_number, amount) VALUES ($1, $2)", phone_number, ct.ChargeAmount)
	if err != nil {
		h.l.Error().Msgf("An error occured while executing insert query to outbox: %v", err)
		return nil, util_error.NewInternalServerError("Somthing went wrong")
	}

	return ct, nil

}

func (h *dbHandler) GetUsersByCoupon(code string, ctx context.Context) (*[]string, error) {

	_, err := h.GetCoupon(code, ctx)

	if err != nil {
		h.l.Error().Msgf("An error occured while executing query: %v", err)
		return nil, err
	}

	var coupons []string

	row, err := h.db.Query(ctx, `SELECT phone_number FROM couponsـused WHERE coupon_id = (select coupon_id from coupons where code = $1 )`, code)

	if err != nil {
		h.l.Error().Msgf("An error occured while executing query: %v", err)
		return nil, util_error.NewInternalServerError("Somthing went wrong")
	}

	for row.Next() {
		var coupon string
		err = row.Scan(&coupon)

		if err != nil {
			h.l.Error().Msgf("An error occured while executing query: %v", err)
			return nil, util_error.NewInternalServerError("Somthing went wrong")
		}
		coupons = append(coupons, coupon)
	}

	return &coupons, nil

}

func (h *dbHandler) GetOutboxes(ctx context.Context) (*[]Outbox, error) {

	var outboxes []Outbox

	row, err := h.db.Query(ctx, `SELECT * FROM outbox`)

	if err != nil {
		h.l.Error().Msgf("An error occured while executing query: %v", err)
		return nil, util_error.NewInternalServerError("Somthing went wrong")
	}

	for row.Next() {
		var outbox Outbox
		err = row.Scan(&outbox.Id, &outbox.PhoneNumber, &outbox.Amount)

		if err != nil {
			h.l.Error().Msgf("An error occured while executing query: %v", err)
			return nil, util_error.NewInternalServerError("Somthing went wrong")
		}

		outboxes = append(outboxes, outbox)
	}

	return &outboxes, nil
}

func (h *dbHandler) DeleteOutbox(ids *[]int, ctx context.Context) error {

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	for _, v := range *ids {
		_, err = tx.Exec(ctx, "DELETE FROM outbox WHERE id=$1", v)
		if err != nil {
			h.l.Error().Msgf("An error occured while executing delete query: %v", err)
			return util_error.NewInternalServerError("Somthing went wrong")
		}
	}

	return nil
}

func (h *dbHandler) CloseConnection() {
	h.db.Close()
}
