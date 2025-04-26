package service

import (
	"context"
	"github.com/savioruz/goth/pkg/failure"
	"github.com/savioruz/goth/pkg/logger"
	"github.com/savioruz/goth/pkg/postgres"

	"github.com/savioruz/goth/internal/domains/user/dto"
	"github.com/savioruz/goth/internal/domains/user/repository"
)

type UserService interface {
	Profile(ctx context.Context, email string) (*dto.UserProfileResponse, error)
}

type userService struct {
	db     postgres.PgxIface
	repo   repository.Querier
	logger logger.Interface
}

func New(db postgres.PgxIface, repo repository.Querier, l logger.Interface) UserService {
	return &userService{
		db:     db,
		repo:   repo,
		logger: l,
	}
}

func (s *userService) Profile(ctx context.Context, email string) (res *dto.UserProfileResponse, err error) {
	user, err := s.repo.GetUserByEmail(ctx, s.db, email)
	if err != nil {
		s.logger.Error("service - user - profile - failed to get user by email", err)

		return nil, failure.InternalError(err)
	}

	if user == (repository.User{}) {
		s.logger.Error("service - user - profile - user not found")

		return nil, failure.NotFound("user not found")
	}

	res = new(dto.UserProfileResponse).ToProfileResponse(user)

	return res, nil
}
