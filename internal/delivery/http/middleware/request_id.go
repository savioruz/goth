package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set("X-Request-ID", requestID)

		ctx := context.WithValue(c.Context(), requestIDKey, requestID)
		c.SetUserContext(ctx)

		return c.Next()
	}
}
