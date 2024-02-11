package cache

import (
	util_error "arvan-challenge/pkg/utils/errors"
	"arvan-challenge/services/coupon/internal/db"
	"arvan-challenge/services/coupon/pkg/env"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var ALLOWED_TIME_KEY = "allowed-times"

type Cache interface {
	SetKeys(key string, fv map[string]any, ctx context.Context) error
	SetKeyForCoupon(coupon *db.Coupon, ctx context.Context) error
	FieldExists(key, field string, ctx context.Context) error
	CheckAndUseCoupon(key, phone_number string, ctx context.Context) error
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

func (c *cache) SetKeys(key string, fv map[string]any, ctx context.Context) error {

	err := c.redisClient.HSet(ctx, key, fv).Err()

	if err != nil {
		c.l.Error().Msgf("An error occured while set key to redis: %v", err)
		return util_error.NewInternalServerError("Somthing went wrong")
	}

	return nil

}

func (c *cache) SetKeyForCoupon(coupon *db.Coupon, ctx context.Context) error {

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

func (c *cache) FieldExists(key, field string, ctx context.Context) error {

	err := c.redisClient.HExists(ctx, key, field).Err()

	if err != nil {
		c.l.Error().Msgf("An error occured while set key to redis: %v", err)
		return util_error.NewInternalServerError("Somthing went wrong")
	}

	return nil

}

func (c *cache) CheckAndUseCoupon(key, phone_number string, ctx context.Context) error {

	var incrBy = redis.NewScript(`
		local current_value = redis.call('HGET', KEYS[1], ARGV[1])
		if current_value and tonumber(current_value) > 0 then
			-- Check if phone number exists as a field
			local phone_exists = redis.call('HEXISTS', KEYS[1], ARGV[3])
			if tonumber(phone_exists) == 0 then
				-- Phone number doesn't exist, increment the field and add the phone number
				redis.call('HINCRBY', KEYS[1], ARGV[1], ARGV[2])
				redis.call('HSET', KEYS[1], ARGV[3], 1) -- Set phone number to 1 (or any value)
				return 1 -- Indicate new phone number added
			else
				-- Phone number already exists, skip increment
				return -1
			end
		else
		-- Field value not positive, skip increment
		return 0
		end
	`)

	keys := []string{key}
	values := []interface{}{ALLOWED_TIME_KEY, -1, phone_number}
	num, err := incrBy.Run(ctx, c.redisClient, keys, values...).Int()

	if err != nil {
		c.l.Error().Msgf("An error occured while running lua function: %v", err)
		return util_error.NewInternalServerError("Somthing went wrong")
	}

	if num == -1 {
		return util_error.NewBadRequestError("This user already used this coupon")
	} else if num == 0 {
		return util_error.NewBadRequestError("The number of times allowed to use this coupon has ended")
	} else if num != 1 {
		return util_error.NewInternalServerError("Somthing went wrong")
	}

	return nil

}
