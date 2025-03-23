package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savioruz/goth/config"
	v1 "github.com/savioruz/goth/internal/delivery/http/v1"
	"github.com/savioruz/goth/internal/service"
	"github.com/savioruz/goth/pkg/logger"
)

func NewRouter(app *fiber.App, cfg *config.Config, l logger.Interface, s service.Services) {
	apiV1Group := app.Group("/v1")
	{
		v1.NewUserRoutes(apiV1Group, l, s.UserService)
	}
}
