package v1

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/savioruz/goth/internal/dto/request"
	"github.com/savioruz/goth/internal/dto/response"
	"github.com/savioruz/goth/internal/service"
	"github.com/savioruz/goth/pkg/logger"
)

type authRoutes struct {
	u service.AuthService
	l logger.Interface
	v *validator.Validate
}

func NewAuthRoutes(group fiber.Router, l logger.Interface, u service.AuthService) {
	r := &authRoutes{u, l, validator.New(validator.WithRequiredStructEnabled())}

	userGroup := group.Group("/auth")
	{
		userGroup.Post("/register", r.Register)
		userGroup.Post("/login", r.Login)
	}
}

// Register godoc
// @Summary Register new user
// @Description Register new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body request.UserRegisterRequest true "User register request"
// @Success 201 {object} response.Response[response.UserRegisterResponse]
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/register [post]
func (r *authRoutes) Register(ctx *fiber.Ctx) error {
	var req request.UserRegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, response.ErrorMsg{"BODY": {"failed to parse request body"}})
	}

	if err := r.v.Struct(req); err != nil {
		return response.NewErrorValidationResponse(ctx, err)
	}

	data, err := r.u.Register(ctx.UserContext(), req)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("user with email %s already exist", req.Email)) {
			return response.NewErrorResponse(ctx, fiber.StatusConflict, response.ErrorMsg{"CONFLICT": {"user already exist"}})
		}

		reqID := "unknown"
		if id, ok := ctx.Locals("request_id").(string); ok {
			reqID = id
		}
		r.l.Error("http - v1 - auth - register - request_id: " + reqID + " - " + err.Error())
		return response.NewErrorResponse(ctx, fiber.StatusInternalServerError, response.ErrorMsg{"SERVER": {"failed to register"}})
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		response.NewResponse(data, nil),
	)
}

// Login godoc
// @Summary Login user
// @Description Login user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body request.UserLoginRequest true "User login request"
// @Success 201 {object} response.Response[response.UserLoginResponse]
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/login [post]
func (r *authRoutes) Login(ctx *fiber.Ctx) error {
	var req request.UserLoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, response.ErrorMsg{"BODY": {"failed to parse request body"}})
	}

	if err := r.v.Struct(req); err != nil {
		return response.NewErrorValidationResponse(ctx, err)
	}

	data, err := r.u.Login(ctx.UserContext(), req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NewErrorResponse(ctx, fiber.StatusBadRequest, response.ErrorMsg{"NOT_FOUND": {"user not found"}})
		}

		if strings.Contains(err.Error(), "unauthorized") {
			return response.NewErrorResponse(ctx, fiber.StatusUnauthorized, response.ErrorMsg{"UNAUTHORIZED": {"invalid email or password"}})
		}

		reqID := "unknown"
		if id, ok := ctx.Locals("request_id").(string); ok {
			reqID = id
		}
		r.l.Error("http - v1 - auth - register - request_id: " + reqID + " - " + err.Error())
		return response.NewErrorResponse(ctx, fiber.StatusInternalServerError, response.ErrorMsg{"SERVER": {"failed to login"}})
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		response.NewResponse(data, nil),
	)
}
