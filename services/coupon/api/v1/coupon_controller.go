package v1

import (
	util_error "arvan-challenge/pkg/utils/errors"
	"arvan-challenge/services/coupon/api/dtos"
	"arvan-challenge/services/coupon/internal/cache"
	"arvan-challenge/services/coupon/internal/db"
	"arvan-challenge/services/coupon/pkg/env"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type CouponController interface {
	CreateCoupon(ctx *fiber.Ctx) error
}

type couponController struct {
	cache cache.Cache
	db    db.DBHandler
	cfg   *env.Config
	l     *zerolog.Logger
}

func NewCouponController(db db.DBHandler, cache cache.Cache, l *zerolog.Logger, cfg *env.Config) CouponController {
	return &couponController{
		db:    db,
		cache: cache,
		l:     l,
		cfg:   cfg,
	}
}

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func (c *couponController) CreateCoupon(ctx *fiber.Ctx) error {

	coupon_dto := new(dtos.CouponDto)

	// map body to dto instance
	if err := ctx.BodyParser(coupon_dto); err != nil {
		c.l.Warn().Msgf(err.Error())
		return util_error.NewInternalServerError(err.Error())
	}

	err := validate.Struct(coupon_dto)
	if validationErrors, ok := err.(validator.ValidationErrors); err != nil && !ok {
		fmt.Println(ok, err)
		c.l.Warn().Msg(validationErrors.Error())
		return util_error.NewInternalServerError(validationErrors)
	}

	coupon := coupon_dto.ConvertToModel()

	val, err := c.db.CreateCoupon(coupon, ctx.Context())

	if err != nil {
		return err
	}
	err = c.cache.SetKeyForCoupon(val, ctx.Context())

	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(val)
}

func (c *couponController) UseCoupon(ctx *fiber.Ctx) error {

	use_coupon_dto := new(dtos.UseCoupon)

	// map body to dto instance
	if err := ctx.BodyParser(use_coupon_dto); err != nil {
		c.l.Warn().Msgf(err.Error())
		return util_error.NewInternalServerError(err.Error())
	}

	err := validate.Struct(use_coupon_dto)
	if validationErrors, ok := err.(validator.ValidationErrors); err != nil && !ok {
		fmt.Println(ok, err)
		c.l.Warn().Msg(validationErrors.Error())
		return util_error.NewInternalServerError(validationErrors)
	}

	// coupon := coupon_dto.ConvertToModel()
	return nil
}
