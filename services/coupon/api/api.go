package api

import (
	"fmt"

	fiber "github.com/gofiber/fiber/v2"
	fiber_logger "github.com/gofiber/fiber/v2/middleware/logger"
	fiber_recover "github.com/gofiber/fiber/v2/middleware/recover"

	zerolog "github.com/rs/zerolog"

	"arvan-challenge/services/coupon/internal/cache"
	"arvan-challenge/services/coupon/internal/db"
	env "arvan-challenge/services/coupon/pkg/env"

	util_error "arvan-challenge/pkg/utils/errors"

	v1 "arvan-challenge/services/coupon/api/v1"
)

type Api interface {
	Init()
	Run()
}

type api struct {
	app       *fiber.App
	logger    *zerolog.Logger
	config    *env.Config
	dbHandler db.DBHandler
	cache     cache.Cache
}

func NewApi(l *zerolog.Logger, cfg *env.Config) Api {
	return &api{
		logger: l,
		config: cfg,
	}
}

func (a *api) Init() {
	cfg := fiber.Config{
		EnablePrintRoutes:     false,
		DisableStartupMessage: true,
		ErrorHandler:          util_error.ErrorHandler,
	}

	a.app = fiber.New(cfg)

	a.dbHandler = db.NewDBHandler(a.config, a.logger)

	a.cache = cache.NewCache(a.config, a.logger)

	// register recover middleware for catching error
	a.app.Use(fiber_recover.New())

	// // register logger middleware for logging request
	a.app.Use(fiber_logger.New())

	// register controllers
	a.AddControllers()
}

func (a *api) AddControllers() {
	api := a.app.Group("/api")
	api_v1 := api.Group("/v1")
	healthController := v1.NewHealthController()

	api_v1.Get("/health1", healthController.Health1)
	api_v1.Get("/health2", healthController.Health2)

	couponController := v1.NewCouponController(a.dbHandler, a.cache, a.logger, a.config)

	api_v1.Post("/create-coupon", couponController.CreateCoupon) // POST /api/v1/create-coupon
	api_v1.Post("/use-coupon", couponController.UseCoupon)       // POST /api/v1/use-coupon

	api_v1.Get("/users-used-coupon/:code", couponController.GetUsersByCoupons) // Get /api/v1/users-used-coupon
}

func (a *api) Run() {
	// run fiber
	addr := fmt.Sprintf("0.0.0.0:%v", a.config.Port)

	a.logger.Info().Msgf("running fiber router with address: %v", addr)
	err := a.app.Listen(addr)

	if err != nil {
		a.logger.Error().Msgf(err.Error())
	}
}
