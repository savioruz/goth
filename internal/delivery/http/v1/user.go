package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/savioruz/goth/internal/dto/request"
	"github.com/savioruz/goth/internal/dto/response"
	"github.com/savioruz/goth/internal/service"
	"github.com/savioruz/goth/pkg/logger"
)

type userRoutes struct {
	u service.UserService
	l logger.Interface
	v *validator.Validate
}

func NewUserRoutes(group fiber.Router, l logger.Interface, u service.UserService) {
	r := &userRoutes{u, l, validator.New(validator.WithRequiredStructEnabled())}

	userGroup := group.Group("/user")
	{
		userGroup.Post("/", r.CreateUser)
	}
}

func (r *userRoutes) CreateUser(ctx *fiber.Ctx) error {
	var req request.CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return newErrorResponse(ctx, fiber.StatusBadRequest, errorMsg{"body": {"failed to parse request body"}})
	}

	if err := r.v.Struct(req); err != nil {
		return newErrorResponse(ctx, fiber.StatusUnprocessableEntity, errorMsg{"validation": {err.Error()}})
	}

	data, err := r.u.Create(ctx.UserContext(), req)
	if err != nil {
		r.l.Error(
			"failed to create user",
			"request_id", ctx.Get("X-Request-ID"),
			"error", err,
		)
		return newErrorResponse(ctx, fiber.StatusInternalServerError, errorMsg{"server": {"failed to create user"}})
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		response.NewResponse(data, nil),
	)
}
