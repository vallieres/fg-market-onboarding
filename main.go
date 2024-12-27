package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html/v2"
	"github.com/rollbar/rollbar-go"

	"github.com/vallieres/fg-market-onboarding/handler"
	customtemplate "github.com/vallieres/fg-market-onboarding/internal/customTemplate"
	"github.com/vallieres/fg-market-onboarding/routes"
	"github.com/vallieres/fg-market-onboarding/services"
)

//go:embed views/*
var viewsFS embed.FS

//go:embed public/*
var publicFS embed.FS

func main() {
	env := os.Getenv("FGONBOARDING_ENVIRONMENT")
	environment := "local"
	if len(env) > 0 {
		environment = env
	}

	// Initialize standard Go html template engine
	subViewsFS, _ := fs.Sub(viewsFS, "views")
	engine := html.NewFileSystem(http.FS(subViewsFS), ".html")
	engine.AddFunc("unescape", customtemplate.Unescape)
	engine.AddFunc("inc", customtemplate.Inc)

	app := fiber.New(fiber.Config{
		Views:     engine,
		Prefork:   false,
		BodyLimit: 8 * 1024 * 1024, //nolint:mnd // 8MB
	})
	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: 1,
	}))

	// Security Headers
	app.Use(helmet.New())

	app.Use(requestid.New(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return "static-id"
		},
	}))

	shopifyAppToken := os.Getenv("FGONBOARDING_SHOPIFY_TOKEN")
	shopifyStorefrontToken := os.Getenv("FGONBOARDING_SHOPIFY_STOREFRONT_TOKEN")

	rateLimiterService := services.NewRateLimiterService()
	customerService := services.NewCustomerService(shopifyAppToken, shopifyStorefrontToken)

	app.Use(limiter.New(limiter.Config{
		Max:                    30, //nolint:mnd
		Expiration:             1 * time.Minute,
		KeyGenerator:           rateLimiterService.KeyGenerator,
		LimitReached:           rateLimiterService.LimitReached,
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      limiter.FixedWindow{},
		Next:                   rateLimiterService.Next,
	}))

	router := &routes.Router{
		Engine: app,
		Public: &handler.PublicHandlers{
			CustomerService: *customerService,
		},
		Common:   &handler.CommonHandlers{},
		PublicFS: publicFS,
	}

	initError := router.Init()
	if initError != nil {
		errorDetails := fmt.Errorf("unable init router: %w", initError)
		rollbar.Error(errorDetails)
		log.Fatal(errorDetails)
	}

	port := os.Getenv("FGONBOARDING_SERVER_PORT")

	if strings.EqualFold(environment, "local") {
		log.Fatal(app.ListenTLS(":443", "./market.furrygarden.io.pem", "./market.furrygarden.io-key.pem"))
	}
	log.Fatal(app.Listen(":" + port))
}
