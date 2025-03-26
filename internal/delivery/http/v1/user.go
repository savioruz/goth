package v1

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/savioruz/goth/internal/delivery/http/middleware"
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

	userGroup := group.Group("/users")
	{
		userGroup.Post("/register", r.UserRegister)
		userGroup.Post("/login", r.UserLogin)
	}
}

// UserRegister godoc
// @Summary Register new user
// @Description Register new user
// @Tags user
// @Accept json
// @Produce json
// @Param body body request.UserRegisterRequest true "User register request"
// @Success 201 {object} response.Response[response.UserRegisterResponse]
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 422 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/register [post]
func (r *userRoutes) UserRegister(ctx *fiber.Ctx) error {
	var req request.UserRegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return newErrorResponse(ctx, fiber.StatusBadRequest, errorMsg{"BODY": {"failed to parse request body"}})
	}

	if err := r.v.Struct(req); err != nil {
		return newErrorValidationResponse(ctx, err)
	}

	data, err := r.u.Register(ctx.UserContext(), req)
	if err != nil {
		if strings.Contains(err.Error(), "exist") {
			return newErrorResponse(ctx, fiber.StatusConflict, errorMsg{"EXIST": {"user already exist"}})
		}

		reqID := "unknown"
		if id, ok := ctx.UserContext().Value(middleware.RequestIDKey).(string); ok {
			reqID = id
		}
		r.l.Error("http - v1 - user - register - request_id: " + reqID + " - " + err.Error())
		return newErrorResponse(ctx, fiber.StatusInternalServerError, errorMsg{"SERVER": {"failed to register"}})
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		response.NewResponse(data, nil),
	)
}

// UserLogin godoc
// @Summary Login user
// @Description Login user
// @Tags user
// @Accept json
// @Produce json
// @Param body body request.UserLoginRequest true "User login request"
// @Success 201 {object} response.Response[response.UserLoginResponse]
// @Failure 400 {object} errorResponse
// @Failure 422 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/login [post]
func (r *userRoutes) UserLogin(ctx *fiber.Ctx) error {
	var req request.UserLoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return newErrorResponse(ctx, fiber.StatusBadRequest, errorMsg{"BODY": {"failed to parse request body"}})
	}

	if err := r.v.Struct(req); err != nil {
		return newErrorValidationResponse(ctx, err)
	}

	data, err := r.u.Login(ctx.UserContext(), req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return newErrorResponse(ctx, fiber.StatusBadRequest, errorMsg{"NOT_FOUND": {"user not found"}})
		}

		if strings.Contains(err.Error(), "unauthorized") {
			return newErrorResponse(ctx, fiber.StatusUnauthorized, errorMsg{"UNAUTHORIZED": {"invalid email or password"}})
		}

		reqID := "unknown"
		if id, ok := ctx.UserContext().Value(middleware.RequestIDKey).(string); ok {
			reqID = id
		}
		r.l.Error("http - v1 - user - register - request_id: " + reqID + " - " + err.Error())
		return newErrorResponse(ctx, fiber.StatusInternalServerError, errorMsg{"SERVER": {"failed to login"}})
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		response.NewResponse(data, nil),
	)
}
