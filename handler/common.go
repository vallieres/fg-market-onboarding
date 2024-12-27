package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/vallieres/fg-market-onboarding/services"
)

// Handlers connecting all the controllers.
type Handlers struct {
	CustomerService services.CustomerService
	ZipCodeService  services.ZipCodeService
}

type (
	CommonHandlers Handlers
)

func ShowErrorPage(c *fiber.Ctx, errorMsg string) error {
	return c.Render("error", fiber.Map{
		"ErrorMessage": errorMsg,
	})
}

func (v1 *CommonHandlers) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("OK")
}
