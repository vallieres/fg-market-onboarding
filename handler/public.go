package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

func (v1 *PublicHandlers) PreparePlanGET(c *fiber.Ctx) error {
	planID := c.Params("planID")

	fmt.Fprintf(c, "%s\n", planID)

	// TODO Retrieve plan.go details

	// TODO Create Cart and generate URL

	// and pass it to the render
	return c.Render("prepare-plan", fiber.Map{
		"PlanID":      planID,
		"PlanDetails": nil,
		"CheckoutURL": "",
	})
}

func (v1 *PublicHandlers) OnboardPOST(c *fiber.Ctx) error {
	var customerDetails model.OnboardPostBody

	ctx := context.Background()

	errParse := c.BodyParser(&customerDetails)
	if errParse != nil {
		return c.Render("index", fiber.Map{
			"ErrorMessage":    errParse.Error(),
			"CustomerDetails": customerDetails,
		})
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

	planID, errCreatePlan := v1.PlanService.CreateBasicPlan(ctx, customerDetails)
	if errCreatePlan != nil {
		rollbar.Error("Unable to create plan : ", errCreatePlan.Error())
		return c.Render("index", fiber.Map{
			"ErrorMessage":    "Unable to Create Plan : " + errCreatePlan.Error(),
			"CustomerDetails": customerDetails,
		})
	}

	domain := os.Getenv("FGONBOARDING_DOMAIN")
	redirectTo := fmt.Sprintf("https://%s/prepare-plan/%d", domain, planID)

	return c.Render("plan.go", fiber.Map{
		"Message":    "Preparing plan.go for " + customerDetails.PetName + "...",
		"RedirectTo": redirectTo,
	})
}

func (v1 *PublicHandlers) ResetGET(c *fiber.Ctx) error {
	return c.Render("reset", fiber.Map{})
}

func (v1 *PublicHandlers) HomeGET(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Home | ",
	})
	/**
	* This was used before, probably not needed anymore.
	domain := os.Getenv("FGONBOARDING_DOMAIN")
	redirectTo := fmt.Sprintf("https://%s/", domain)
	if os.Getenv("FGONBOARDING_ENVIRONMENT") == EnvLocal {
		return c.Render("index", fiber.Map{
			"Title": "Home | " + redirectTo,
		})
	}

	return c.Redirect(redirectTo)
	**/
}

func (v1 *PublicHandlers) RESTCitiesGET(c *fiber.Ctx) error {
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

func (v1 *PublicHandlers) RESTTestGET(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"hello": "world",
	})
}

func (v1 *PublicHandlers) RESTPlansGET(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"hello": "world",
	})
}

func (v1 *PublicHandlers) RESTIsPlanReadyGET(c *fiber.Ctx) error {
	planID := c.Params("planID")
	planIDInt, errConv := strconv.Atoi(planID)
	if errConv != nil {
		rollbar.Error(fmt.Errorf("error converting planID string to int64: %s, %w", planID, errConv))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "error checking if plan is ready",
		})
	}
	isPlanReady, errIsPlanReady := v1.PlanService.IsPlanReady(int64(planIDInt))
	if errIsPlanReady != nil {
		if strings.Contains(errIsPlanReady.Error(), "plan is not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"is_ready": "false",
			})
		}
		rollbar.Error(fmt.Errorf("error checking if plan is ready: %s, %w", planID, errIsPlanReady))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "error checking if plan is ready",
		})
	}
	isPlanReadyStr := strconv.FormatBool(isPlanReady)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"is_ready": isPlanReadyStr,
	})
}
