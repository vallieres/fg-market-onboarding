package app

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/vallieres/fg-market-onboarding/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html/v2"
	"github.com/pressly/goose/v3"
	"github.com/rollbar/rollbar-go"

	"github.com/vallieres/fg-market-onboarding/handler"
	"github.com/vallieres/fg-market-onboarding/internal/customtemplate"
	"github.com/vallieres/fg-market-onboarding/internal/database"
	"github.com/vallieres/fg-market-onboarding/routes"
	"github.com/vallieres/fg-market-onboarding/services"
)

//go:embed db/*.sql
var migrationsFS embed.FS

//go:embed views/*
var viewsFS embed.FS

//go:embed public/*
var publicFS embed.FS

//nolint:funlen
func Run() {
	env := os.Getenv("FGONBOARDING_ENVIRONMENT")

	// Setup Rollbar
	rollbar.SetToken(os.Getenv("FGONBOARDING_ROLLBAR_TOKEN"))
	rollbar.SetEnvironment(env)
	rollbar.SetCodeVersion("main")
	rollbar.SetServerRoot("github.com/vallieres/fg-market-onboarding")
	rollbar.Wait()
	if os.Getenv("FGONBOARDING_DISABLE_ROLLBAR") == "true" {
		rollbar.SetEnabled(false)
	}

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

	db, err := database.MySQLConnection()
	if err != nil {
		errorDetails := fmt.Errorf("unable to connect to the database: %w", err)
		rollbar.Error(errorDetails)
		log.Fatal(errorDetails)
	}

	// Setup Goose Migrations
	goose.SetBaseFS(migrationsFS)

	// Update Database
	if errSetDialect := goose.SetDialect("mysql"); errSetDialect != nil {
		errorDetails := fmt.Errorf("unable to set database dialiect: %w", errSetDialect)
		rollbar.Error(errorDetails)
		panic(errSetDialect)
	}

	if errUp := goose.Up(db.DB, "db"); errUp != nil {
		errorDetails := fmt.Errorf("unable upgrade migrationsFS: %w", errUp)
		rollbar.Error(errorDetails)
		panic(errorDetails)
	}

	shopifyAppToken := os.Getenv("FGONBOARDING_SHOPIFY_TOKEN")
	shopifyStorefrontToken := os.Getenv("FGONBOARDING_SHOPIFY_STOREFRONT_TOKEN")

	zipCodeRepository := repository.NewZipCodeRepository(db)
	planRepository := repository.NewPlanRepository(db)

	rateLimiterService := services.NewRateLimiterService()
	customerService := services.NewCustomerService(shopifyAppToken, shopifyStorefrontToken)
	zipCodeService := services.NewZipCodeService(zipCodeRepository)
	planService := services.NewPlanService(planRepository, shopifyAppToken, shopifyStorefrontToken)

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
			ZipCodeService:  *zipCodeService,
			PlanService:     *planService,
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
