package handler

import (
	"context"
	"fmt"
	"os"
	"strconv"

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
		return c.Render("signup", fiber.Map{
			"ErrorMessage": err.Error(),
		})
	}

	userEmail, errCreate := v1.CustomerService.Create(ctx, customerDetails)
	if errCreate != nil {
		rollbar.Error("Unable to create user : ", errCreate.Error())
		return c.Render("signup", fiber.Map{
			"ErrorMessage": "Unable to Create User : " + errCreate.Error(),
		})
	}
	if userEmail != "" {
		rollbar.Warning("User Already Exists, UserID: " + strconv.FormatInt(0, 10))
		return c.Render("signup", fiber.Map{
			"ErrorMessage": "User already exists, please login instead.",
		})
	}

	redirectTo := ""

	return c.Render("index", fiber.Map{
		"Message":    "You have successfully signed up! Check your email for the verification link and then login below!",
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
