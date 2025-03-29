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

type authRoutes struct {
	a service.AuthService
	o service.OAuthService
	l logger.Interface
	v *validator.Validate
}

func NewAuthRoutes(group fiber.Router, l logger.Interface, u service.AuthService, o service.OAuthService) {
	r := &authRoutes{u, o, l, validator.New(validator.WithRequiredStructEnabled())}

	authGroup := group.Group("/auth")
	{
		authGroup.Post("/register", r.Register)
		authGroup.Post("/login", r.Login)
		authGroup.Get("/google/login", r.GoogleLogin)
		authGroup.Get("/google/callback", r.GoogleCallback)

		// Private routes
		authGroup.Get("/profile", middleware.Jwt(), r.Profile)
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

	data, err := r.a.Register(ctx.UserContext(), req)
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

	data, err := r.a.Login(ctx.UserContext(), req)
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

// GoogleLogin godoc
// @Summary Login with Google
// @Description Redirects to Google OAuth consent screen
// @Tags auth
// @Accept json
// @Produce json
// @Success 302 {string} string "Redirect to Google"
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/google/login [get]
func (r *authRoutes) GoogleLogin(c *fiber.Ctx) error {
	url := r.o.GetGoogleAuthURL()
	return c.Redirect(url)
}

// GoogleCallback godoc
// @Summary Google OAuth callback
// @Description Handle the Google OAuth callback and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Success 200 {object} response.Response[response.UserLoginResponse]
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/google/callback [get]
func (r *authRoutes) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return response.NewErrorResponse(c, fiber.StatusBadRequest, response.ErrorMsg{"CODE": {"missing authorization code"}})
	}

	data, err := r.o.HandleGoogleCallback(c.Context(), code)
	if err != nil {
		if strings.Contains(err.Error(), "failed to exchange code") {
			return response.NewErrorResponse(c, fiber.StatusBadRequest, response.ErrorMsg{"CODE": {"invalid authorization code"}})
		}

		reqID := "unknown"
		if id, ok := c.Locals("request_id").(string); ok {
			reqID = id
		}
		r.l.Error("http - v1 - auth - google callback - request_id: " + reqID + " - " + err.Error())
		return response.NewErrorResponse(c, fiber.StatusInternalServerError, response.ErrorMsg{"SERVER": {"failed to process OAuth callback"}})
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse(data, nil),
	)
}

// Profile godoc
// @Summary Get user profile
// @Description Get user profile
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response[response.UserProfileResponse]
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/profile [get]
// @Security BearerAuth
func (r *authRoutes) Profile(ctx *fiber.Ctx) error {
	email := ctx.Locals("email").(string)
	data, err := r.a.Profile(ctx.UserContext(), email)
	if err != nil {
		reqID := "unknown"
		if id, ok := ctx.Locals("request_id").(string); ok {
			reqID = id
		}
		r.l.Error("http - v1 - auth - profile - request_id: " + reqID + " - " + err.Error())
		return response.NewErrorResponse(ctx, fiber.StatusInternalServerError, response.ErrorMsg{"SERVER": {"failed to get profile"}})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		response.NewResponse(data, nil),
	)
}
