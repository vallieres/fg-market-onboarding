package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rollbar/rollbar-go"

	"github.com/vallieres/fg-market-onboarding/model"
)

type (
	PublicHandlers Handlers
)

// EnvLocal represents the local environment.
const EnvLocal = "local"

func (v1 *PublicHandlers) OnboardGET(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func (v1 *PublicHandlers) PlanResultGET(c *fiber.Ctx) error {
	fmt.Fprintf(c, "%s\n", c.Params("planID"))

	// TODO Retrieve plan details

	// TODO Create Cart and generate URL

	// and pass it to the render
	return c.Render("plan-result", fiber.Map{
		"PlanDetails": nil,
		"CheckoutURL": "",
	})
}

func (v1 *PublicHandlers) OnboardPOST(c *fiber.Ctx) error {
	var customerDetails model.OnboardPostBody

	ctx := context.Background()

	errParse := c.BodyParser(&customerDetails)
	if errParse != nil {
		return ShowErrorPage(c, errParse.Error())
	}

	// validation
	if err := customerDetails.Validate(); err != nil {
		rollbar.Warning("Unable to create user : ", err)
		return c.Render("index", fiber.Map{
			"ErrorMessage":    err.Error(),
			"CustomerDetails": customerDetails,
		})
	}

	_, errCreate := v1.CustomerService.Create(ctx, customerDetails)
	if errCreate != nil {
		rollbar.Error("Unable to create user : ", errCreate.Error())
		return c.Render("index", fiber.Map{
			"ErrorMessage":    "Unable to Create User : " + errCreate.Error(),
			"CustomerDetails": customerDetails,
		})
	}

	domain := os.Getenv("FGONBOARDING_DOMAIN")
	redirectTo := fmt.Sprintf("https://%s/plan-result/%s", domain, "unique-plan-id")

	return c.Render("plan", fiber.Map{
		"Message":    "Preparing plan for " + customerDetails.DogName + "...",
		"RedirectTo": redirectTo,
	})
}

func (v1 *PublicHandlers) ResetGET(c *fiber.Ctx) error {
	return c.Render("reset", fiber.Map{})
}

func (v1 *PublicHandlers) HomeGET(c *fiber.Ctx) error {
	domain := os.Getenv("FGONBOARDING_DOMAIN")
	redirectTo := fmt.Sprintf("https://%s/", domain)
	if os.Getenv("FGONBOARDING_ENVIRONMENT") == EnvLocal {
		return c.Render("index", fiber.Map{
			"Title": "Home | " + redirectTo,
		})
	}

	return c.Redirect(redirectTo)
}

func (v1 *PublicHandlers) CitiesGET(c *fiber.Ctx) error {
	cities, errGetCities := v1.ZipCodeService.GetCityByZipCode(c.Params("zipCode"))
	if errGetCities != nil {
		if strings.Contains(errGetCities.Error(), "no zip code entries found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "no zip code entries found",
			})
		}
		rollbar.Error(fmt.Errorf("error pulling cities for zip code: %s, %w", c.Params("zipCode"), errGetCities))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "error pulling cities for zip code",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"cities": cities,
	})
}
