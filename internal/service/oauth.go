package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/savioruz/goth/internal/dto/response"
	"github.com/savioruz/goth/internal/repository"
	"github.com/savioruz/goth/pkg/jwt"
	"github.com/savioruz/goth/pkg/oauth"
)

type OAuthService interface {
	GetGoogleAuthURL() string
	HandleGoogleCallback(ctx context.Context, code string) (*response.UserLoginResponse, error)
}

type oauthService struct {
	db             *pgxpool.Pool
	q              *repository.Queries
	googleProvider *oauth.GoogleProvider
}

func NewOAuthService(db *pgxpool.Pool, googleProvider *oauth.GoogleProvider) OAuthService {
	return &oauthService{
		db:             db,
		q:              repository.New(),
		googleProvider: googleProvider,
	}
}

func (s *oauthService) GetGoogleAuthURL() string {
	return s.googleProvider.GetAuthURL()
}

func (s *oauthService) HandleGoogleCallback(ctx context.Context, code string) (*response.UserLoginResponse, error) {
	token, err := s.googleProvider.Exchange(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	userInfo, err := s.googleProvider.GetUserInfo(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if user exists
	user, err := s.q.GetUserByEmail(ctx, tx, userInfo.Email)
	if err != nil {
		params := repository.CreateUserParams{
			Email:        userInfo.Email,
			FullName:     pgtype.Text{String: userInfo.Name, Valid: true},
			IsVerified:   pgtype.Bool{Bool: userInfo.VerifiedEmail, Valid: true},
			Level:        "1",
			ProfileImage: pgtype.Text{String: userInfo.Picture, Valid: true},
		}
		user, err = s.q.CreateUser(ctx, tx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String(), user.Email, user.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID.String(), user.Email, user.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &response.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
