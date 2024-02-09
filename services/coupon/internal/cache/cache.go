package cache

import (
	util_error "arvan-challenge/pkg/utils/errors"
	"arvan-challenge/services/coupon/internal/db"
	"arvan-challenge/services/coupon/pkg/env"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

var ALLOWED_TIME_KEY = "allowed-times"

type Cache interface {
	SetKeys(key string, fv map[string]any, ctx *fasthttp.RequestCtx) error
	SetKeyForCoupon(coupon *db.Coupon, ctx *fasthttp.RequestCtx) error
}

type cache struct {
	cfg         *env.Config
	l           *zerolog.Logger
	redisClient *redis.Client
}

func NewCache(cfg *env.Config, l *zerolog.Logger) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			l.Info().Msg("new redis connection established \n")
			return nil
		},
	})

	return &cache{
		redisClient: rdb,
		l:           l,
		cfg:         cfg,
	}
}

func (c *cache) SetKeys(key string, fv map[string]any, ctx *fasthttp.RequestCtx) error {

	err := c.redisClient.HSet(ctx, key, fv).Err()

	if err != nil {
		c.l.Error().Msgf("An error occured while set key to redis: %v", err)
		return util_error.NewInternalServerError("Somthing went wrong")
	}

	return nil

}

func (c *cache) SetKeyForCoupon(coupon *db.Coupon, ctx *fasthttp.RequestCtx) error {

	err := c.redisClient.HSet(ctx, coupon.Code, ALLOWED_TIME_KEY, coupon.AllowedTimes).Err()

	if err != nil {
		c.l.Error().Msgf("An error occured while set key in redis: %v", err)
		return util_error.NewInternalServerError("Somthing went wrong")
	}

	err = c.redisClient.Expire(ctx, coupon.Code, time.Until(coupon.ExpireDate)).Err()

	if err != nil {
		c.l.Error().Msgf("An error occured while set expire for key in redis: %v", err)
		return util_error.NewInternalServerError("Somthing went wrong")
	}

	return nil

}

func (c *cache) FieldExists(key, field string, ctx *fasthttp.RequestCtx) error {

	err := c.redisClient.HExists(ctx, key, field).Err()

	if err != nil {
		c.l.Error().Msgf("An error occured while set key to redis: %v", err)
		return util_error.NewInternalServerError("Somthing went wrong")
	}

	return nil

}

func (c *cache) CheckAndUseCoupon(key string, ctx *fasthttp.RequestCtx) error {

	// define pipeline
	pipe := c.redisClient.Pipeline()

	commands := make([]*redis.BoolCmd, 2)

	commands = append(commands, pipe.HExists(ctx, key, ALLOWED_TIME_KEY))

	// if err != nil {
	// 	c.l.Error().Msgf("An error occured while set key in redis: %v", err)
	// 	return util_error.NewInternalServerError("Somthing went wrong")
	// }

	return nil

}
