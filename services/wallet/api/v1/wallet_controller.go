package v1

import (
	util_error "arvan-challenge/pkg/utils/errors"
	route_validator "arvan-challenge/pkg/utils/validator"
	"arvan-challenge/services/wallet/api/dtos"
	"arvan-challenge/services/wallet/internal/db"
	"arvan-challenge/services/wallet/pkg/env"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type WalletController interface {
	AddTransaction(ctx *fiber.Ctx) error
	GetUserTransactions(ctx *fiber.Ctx) error
}

type walletController struct {
	db  db.DBHandler
	cfg *env.Config
	l   *zerolog.Logger
}

func NewWalletController(db db.DBHandler, l *zerolog.Logger, cfg *env.Config) WalletController {
	return &walletController{
		db:  db,
		l:   l,
		cfg: cfg,
	}
}

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("IsDate", route_validator.IsDate)
}

func (c *walletController) AddTransaction(ctx *fiber.Ctx) error {

	transaction_dto := new(dtos.TransactionDto)

	// map body to dto instance
	if err := ctx.BodyParser(transaction_dto); err != nil {
		c.l.Warn().Msgf(err.Error())
		return util_error.NewInternalServerError(err.Error())
	}

	err := validate.Struct(transaction_dto)
	fmt.Println(err)
	if validationErrors, ok := err.(validator.ValidationErrors); err != nil && ok {
		c.l.Warn().Msg(validationErrors.Error())
		return util_error.NewInternalServerError(validationErrors.Error())
	}

	fmt.Println(transaction_dto)

	transaction, wallet := transaction_dto.ConvertToModel()

	val, err := c.db.AddTransaction(transaction, wallet, ctx.Context())

	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(val)
}

func (c *walletController) GetUserTransactions(ctx *fiber.Ctx) error {

	phone_number := ctx.Params("phone_number")

	// err := validate.StructFiltered(phone_number, validate.VarWithValue())
	// fmt.Println(err)
	// if validationErrors, ok := err.(validator.ValidationErrors); err != nil && ok {
	// 	c.l.Warn().Msg(validationErrors.Error())
	// 	return util_error.NewInternalServerError(validationErrors.Error())
	// }

	val, err := c.db.GetUserTransactions(phone_number, ctx.Context())

	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(val)

}
