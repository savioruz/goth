package v1

import (
	"fmt"
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
		userGroup.Get("/profile", middleware.Jwt(), r.UserProfile)
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
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/register [post]
func (r *userRoutes) UserRegister(ctx *fiber.Ctx) error {
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
		r.l.Error("http - v1 - user - register - request_id: " + reqID + " - " + err.Error())
		return response.NewErrorResponse(ctx, fiber.StatusInternalServerError, response.ErrorMsg{"SERVER": {"failed to register"}})
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
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/login [post]
func (r *userRoutes) UserLogin(ctx *fiber.Ctx) error {
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
		r.l.Error("http - v1 - user - register - request_id: " + reqID + " - " + err.Error())
		return response.NewErrorResponse(ctx, fiber.StatusInternalServerError, response.ErrorMsg{"SERVER": {"failed to login"}})
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		response.NewResponse(data, nil),
	)
}

// UserProfile godoc
// @Summary Get user profile
// @Description Get user profile
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} response.Response[response.UserProfileResponse]
// @Failure 500 {object} response.ErrorResponse
// @Router /users/profile [get]
// @Security BearerAuth
func (r *userRoutes) UserProfile(ctx *fiber.Ctx) error {
	email, ok := ctx.Locals("email").(string)

	data, err := r.u.Profile(ctx.UserContext(), email)
	if err != nil || !ok {
		reqID := "unknown"
		if id, ok := ctx.Locals("request_id").(string); ok {
			reqID = id
		}
		r.l.Error("http - v1 - user - profile - request_id: " + reqID + " - " + err.Error())
		return response.NewErrorResponse(ctx, fiber.StatusInternalServerError, response.ErrorMsg{"SERVER": {"failed to get profile"}})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		response.NewResponse(data, nil),
	)
}
