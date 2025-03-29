package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/savioruz/goth/internal/dto/request"
	"github.com/savioruz/goth/internal/dto/response"
	"github.com/savioruz/goth/internal/repository"
	"github.com/savioruz/goth/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req request.UserRegisterRequest) (*response.UserRegisterResponse, error)
	Login(ctx context.Context, req request.UserLoginRequest) (*response.UserLoginResponse, error)
	Profile(ctx context.Context, email string) (*response.UserProfileResponse, error)
}

type userService struct {
	db *pgxpool.Pool
	q  *repository.Queries
}

func NewUserService(db *pgxpool.Pool) UserService {
	return &userService{
		db: db,
		q:  repository.New(),
	}
}

func (s *userService) Register(ctx context.Context, req request.UserRegisterRequest) (*response.UserRegisterResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("Register - service - failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	exist, err := s.q.GetUserByEmail(ctx, tx, req.Email)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("Register - service - failed to get user by email: %w", err)
	}

	if exist.Email != "" {
		return nil, fmt.Errorf("Register - service - user with email %s already exist", req.Email)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Register - service - failed to hash password: %w", err)
	}

	newUser, err := s.q.CreateUser(ctx, tx, repository.CreateUserParams{
		Email: req.Email,
		Password: pgtype.Text{
			String: string(password),
			Valid:  true,
		},
		Level: "1",
		FullName: pgtype.Text{
			String: req.Name,
			Valid:  true,
		},
		IsVerified: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Register - service - failed to create user: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("Register - service - failed to commit transaction: %w", err)
	}

	return &response.UserRegisterResponse{
		ID:    newUser.ID.String(),
		Email: newUser.Email,
	}, nil
}

func (s *userService) Login(ctx context.Context, req request.UserLoginRequest) (*response.UserLoginResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("Login - service - failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	user, err := s.q.GetUserByEmail(ctx, tx, req.Email)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("Login - service - failed to get user by email: %w", err)
	}

	if user.Email == "" {
		return nil, fmt.Errorf("Login - service - user with email %s not found", req.Email)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("Login - service - unauthorized")
	}

	// TODO: check if user is verified

	_, err = s.q.UpdateLastLogin(ctx, tx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("Login - service - failed to update last login: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("Login - service - failed to commit transaction: %w", err)
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String(), user.Email, user.Level)
	if err != nil {
		return nil, fmt.Errorf("Login - service - failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID.String(), user.Email, user.Level)
	if err != nil {
		return nil, fmt.Errorf("Login - service - failed to generate refresh token: %w", err)
	}

	return &response.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) Profile(ctx context.Context, email string) (*response.UserProfileResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("Profile - service - failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	user, err := s.q.GetUserByEmail(ctx, tx, email)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("Profile - service - failed to get user by id: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("Profile - service - failed to commit transaction: %w", err)
	}

	return &response.UserProfileResponse{
		Email:        user.Email,
		Name:         user.FullName.String,
		ProfileImage: user.ProfileImage.String,
	}, nil
}
