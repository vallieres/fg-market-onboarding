package routes

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/rollbar/rollbar-go"

	"github.com/vallieres/fg-market-onboarding/handler"
)

type Router struct {
	Engine   *fiber.App
	Public   *handler.PublicHandlers
	Common   *handler.CommonHandlers
	PublicFS embed.FS
}

// Init maps route and actions.
func (r *Router) Init() error {
	if r.Engine == nil {
		return errors.New("Router.Engine top level echo framework instance is nil, cannot be nil")
	}
	if r.Public == nil {
		return errors.New("Router.Public needs to be set, cannot be nil")
	}
	if r.Common == nil {
		return errors.New("Router.Common needs to be set, cannot be nil")
	}

	r.Engine.Get("/", r.Public.HomeGET)

	// Common Routes
	subPublicFS, _ := fs.Sub(r.PublicFS, "public")
	r.Engine.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(subPublicFS),
		PathPrefix: "",
		Browse:     false,
		MaxAge:     60 * 60 * 24 * 7, //nolint:mnd  // 7 days
	}))

	r.Engine.Get("/css/*.css", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "text/css")
		if strings.Contains(rollbar.Environment(), "local") {
			c.Set(fiber.HeaderCacheControl, "no-cache")
		}
		return c.Next()
	})

	r.Engine.Get("/js/*.js", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "text/javascript")
		if strings.Contains(rollbar.Environment(), "local") {
			c.Set(fiber.HeaderCacheControl, "no-cache")
		}
		return c.Next()
	})
	r.Engine.Get("/img/*.jpg", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "image/jpeg")
		if strings.Contains(rollbar.Environment(), "local") {
			c.Set(fiber.HeaderCacheControl, "no-cache")
		}
		return c.Next()
	})
	r.Engine.Get("/img/*.gif", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "image/gif")
		if strings.Contains(rollbar.Environment(), "local") {
			c.Set(fiber.HeaderCacheControl, "no-cache")
		}
		return c.Next()
	})
	r.Engine.Get("/img/*.png", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "image/png")
		if strings.Contains(rollbar.Environment(), "local") {
			c.Set(fiber.HeaderCacheControl, "no-cache")
		}
		return c.Next()
	})

	r.Engine.Get("/health-check", r.Common.HealthCheck)

	// Public Routes
	r.Engine.Post("/onboard", r.Public.OnboardPOST)
	r.Engine.Get("/onboard", r.Public.OnboardGET)
	r.Engine.Get("/plan-result/:planID", r.Public.PlanResultGET)
	r.Engine.Get("/reset", r.Public.ResetGET)

	// REST Calls
	r.Engine.Get("/rest/cities/:zipCode", r.Public.RESTCitiesGET)
	r.Engine.Get("/rest/test", r.Public.RESTTestGET)
	r.Engine.Get("/rest/plans/:email", r.Public.RESTPlansGET)

	r.Engine.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).SendString("Page not found")
		},
	)

	return nil
}
