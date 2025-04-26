package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/savioruz/goth/internal/domains/user/mock"
	"github.com/savioruz/goth/internal/domains/user/repository"
	"github.com/savioruz/goth/pkg/failure"
	log "github.com/savioruz/goth/pkg/logger/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
	"time"
)

func TestUserService_Profile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockQuerier := mock.NewMockQuerier(ctrl)
	mockPgx, _ := pgxmock.NewPool()
	mockLogger := log.NewMockInterface(ctrl)
	mockError := errors.New("error")

	service := New(mockPgx, mockQuerier, mockLogger)

	mockID := uuid.New()
	profileMock := repository.User{
		ID:           pgtype.UUID{Bytes: mockID, Valid: true},
		Email:        "string@gmail.com",
		Password:     pgtype.Text{String: "strongpassword", Valid: true},
		Level:        "user",
		GoogleID:     pgtype.Text{String: "google123", Valid: true},
		FullName:     pgtype.Text{String: "Test User", Valid: true},
		ProfileImage: pgtype.Text{String: "https://example.com/profile.jpg", Valid: true},
		IsVerified:   pgtype.Bool{Bool: true, Valid: true},
		LastLogin:    pgtype.Timestamp{Time: time.Now(), Valid: true},
		CreatedAt:    pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamp{Time: time.Now(), Valid: true},
		DeletedAt:    pgtype.Timestamp{Valid: false},
	}

	t.Run("error: failure getting user by email", func(t *testing.T) {
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())

		mockQuerier.EXPECT().
			GetUserByEmail(gomock.Any(), gomock.Any(), "error@gmail.com").
			Return(repository.User{}, mockError).
			Times(1)

		res, err := service.Profile(ctx, "error@gmail.com")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, http.StatusInternalServerError, failure.GetCode(err))
	})

	t.Run("error: user not found", func(t *testing.T) {
		mockLogger.EXPECT().Error(gomock.Any())

		mockQuerier.EXPECT().
			GetUserByEmail(gomock.Any(), gomock.Any(), "notfound@gmail.com").
			Return(repository.User{}, nil).
			Times(1)

		res, err := service.Profile(ctx, "notfound@gmail.com")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, http.StatusNotFound, failure.GetCode(err))
	})

	t.Run("success: ", func(t *testing.T) {
		mockQuerier.EXPECT().
			GetUserByEmail(gomock.Any(), gomock.Any(), "string@gmail.com").
			Return(profileMock, nil).
			Times(1)

		res, err := service.Profile(ctx, "string@gmail.com")

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "string@gmail.com", res.Email)
		assert.Equal(t, "Test User", res.Name)
		assert.Equal(t, "https://example.com/profile.jpg", res.ProfileImage)
	})
}
