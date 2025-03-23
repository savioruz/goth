package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/savioruz/goth/internal/dto/request"
	"github.com/savioruz/goth/internal/dto/response"
	"github.com/savioruz/goth/internal/repository"
)

type UserService interface {
	Create(ctx context.Context, arg request.CreateUserRequest) (*response.CreateUserResponse, error)
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

func (s *userService) Create(ctx context.Context, arg request.CreateUserRequest) (*response.CreateUserResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	exist, err := s.q.GetUser(ctx, tx, arg.Email)
	if err != nil {
		return nil, err
	}

	if exist.Email != "" {
		return nil, errors.New("email already exist")
	}

	createdUser, err := s.q.CreateUser(ctx, tx, repository.CreateUserParams{
		Email:    arg.Email,
		Password: arg.Password,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &response.CreateUserResponse{
		ID:    createdUser.ID.String(),
		Email: createdUser.Email,
	}, nil
}
