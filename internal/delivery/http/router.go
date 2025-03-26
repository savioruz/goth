package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/savioruz/goth/config"
	_ "github.com/savioruz/goth/docs" // Swagger docs

	"github.com/savioruz/goth/internal/delivery/http/middleware"
	v1 "github.com/savioruz/goth/internal/delivery/http/v1"
	"github.com/savioruz/goth/internal/service"
	"github.com/savioruz/goth/pkg/logger"
)

// Swagger spec:
// @title       Goth API
// @description This is a sample server Goth API.
// @version     1.0
// @host        localhost:3000
// @BasePath    /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(app *fiber.App, cfg *config.Config, l logger.Interface, s service.Services) {
	// Options
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))
	app.Use(middleware.RequestID())

	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	apiV1Group := app.Group("/v1")
	{
		v1.NewUserRoutes(apiV1Group, l, s.UserService)
	}

	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "route not found",
		})
	})
}
