package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/savioruz/goth/internal/delivery/http/middleware"
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
		if reqID, ok := c.UserContext().Value(middleware.RequestIDKey).(string); ok {
			res.RequestID = reqID
		}
	}

	return c.Status(code).JSON(res)
}

func newErrorValidationResponse(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return newErrorResponse(c, fiber.StatusBadRequest, errorMsg{
			"validation": {err.Error()},
		})
	}

	fieldErrors := make(errorMsg)
	for _, e := range validationErrors {
		field := e.Field()
		message := getValidationMessage(e.Tag())
		fieldErrors[field] = append(fieldErrors[field], message)
	}

	return newErrorResponse(c, fiber.StatusUnprocessableEntity, fieldErrors)
}

func getValidationMessage(tag string) string {
	switch tag {
	case "required":
		return "REQUIRED"
	case "required_if":
		return "REQUIRED"
	case "boolean":
		return "MUST_BE_BOOLEAN"
	case "email":
		return "INVALID_EMAIL"
	case "min":
		return "TOO_SHORT"
	case "max":
		return "TOO_LONG"
	case "numeric":
		return "MUST_BE_NUMERIC"
	case "alphanum":
		return "MUST_BE_ALPHANUMERIC"
	default:
		return tag
	}
}
