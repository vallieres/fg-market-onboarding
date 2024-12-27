package services

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rollbar/rollbar-go"
)

type RateLimiterService struct{}

func NewRateLimiterService() *RateLimiterService {
	return &RateLimiterService{}
}

func (u *RateLimiterService) KeyGenerator(c *fiber.Ctx) string {
	return c.IP()
}

func (u *RateLimiterService) LimitReached(c *fiber.Ctx) error {
	rollbar.Warning(fmt.Errorf("rate limit reached for ip: %s", c.IP()))

	return c.Render("error", fiber.Map{
		"ErrorMessage": "Too many requests, please slow down.",
	})
}

func (u *RateLimiterService) Next(c *fiber.Ctx) bool {
	path := c.Context().URI().String()

	if strings.Contains(path, "/css/") ||
		strings.Contains(path, "/js/") ||
		strings.Contains(path, "/fonts/") ||
		strings.Contains(path, ".png") {
		return true
	}
	return false
}
