package v1

import "github.com/gofiber/fiber/v2"

type HealthController interface {
	Health1(ctx *fiber.Ctx) error
	Health2(ctx *fiber.Ctx) error
}

type healthController struct {
}

func NewHealthController() HealthController {
	return &healthController{}
}

func (c *healthController) Health1(ctx *fiber.Ctx) error {
	return ctx.SendString("OK")
}

func (c *healthController) Health2(ctx *fiber.Ctx) error {
	return nil
}
