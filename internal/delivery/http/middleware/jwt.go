package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/savioruz/goth/internal/dto/response"
	"github.com/savioruz/goth/pkg/jwt"
)

func Jwt() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.NewErrorResponse(c, fiber.StatusUnauthorized, response.ErrorMsg{"UNAUTHORIZED": {"authorization header not found"}})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.NewErrorResponse(c, fiber.StatusUnauthorized, response.ErrorMsg{"UNAUTHORIZED": {"invalid authorization header"}})
		}

		claims, err := jwt.ValidateToken(parts[1])
		if err != nil {
			return response.NewErrorResponse(c, fiber.StatusUnauthorized, response.ErrorMsg{"UNAUTHORIZED": {"invalid token"}})
		}

		if claims != nil {
			fmt.Printf("JWT Claims: %+v\n", claims)
			c.Locals("user_id", claims.ID)
			c.Locals("email", claims.Email)
			c.Locals("level", claims.Level)
		}

		return c.Next()
	}
}
