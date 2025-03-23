package v1

import (
	"github.com/gofiber/fiber/v2"
)

type errorMsg map[string][]string

type errorResponse struct {
	RequestID string   `json:"request_id,omitempty"`
	Errors    errorMsg `json:"errors"`
}

func newErrorResponse(c *fiber.Ctx, code int, errors errorMsg) error {
	if len(errors) == 0 {
		return c.Status(code).JSON(errorResponse{
			Errors: errorMsg{"error": {"something went wrong"}},
		})
	}

	res := errorResponse{
		Errors: errors,
	}

	if code >= 500 {
		res.RequestID = c.Get("X-Request-ID")
	}

	return c.Status(code).JSON(res)
}
