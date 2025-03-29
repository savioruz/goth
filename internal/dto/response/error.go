package response

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorMsg map[string][]string

type ErrorResponse struct {
	RequestID string   `json:"request_id,omitempty"`
	Errors    ErrorMsg `json:"errors"`
}

func NewErrorResponse(c *fiber.Ctx, code int, errors ErrorMsg) error {
	if len(errors) == 0 {
		return c.Status(code).JSON(ErrorResponse{
			Errors: ErrorMsg{"error": {"something went wrong"}},
		})
	}

	res := ErrorResponse{
		Errors: errors,
	}

	if code >= 500 {
		if reqID, ok := c.Locals("request_id").(string); ok {
			res.RequestID = reqID
		}
	}

	return c.Status(code).JSON(res)
}

func NewErrorValidationResponse(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return NewErrorResponse(c, fiber.StatusBadRequest, ErrorMsg{
			"validation": {err.Error()},
		})
	}

	fieldErrors := make(ErrorMsg)
	for _, e := range validationErrors {
		field := e.Field()
		message := getValidationMessage(e.Tag())
		fieldErrors[field] = append(fieldErrors[field], message)
	}

	return NewErrorResponse(c, fiber.StatusUnprocessableEntity, fieldErrors)
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
