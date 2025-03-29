package service

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/savioruz/goth/pkg/logger"
	"github.com/savioruz/goth/pkg/oauth"
)

type Factory struct {
	db             *pgxpool.Pool
	l              logger.Interface
	googleProvider *oauth.GoogleProvider
}

func NewFactory(db *pgxpool.Pool, l logger.Interface, googleProvider *oauth.GoogleProvider) *Factory {
	return &Factory{db, l, googleProvider}
}

type Services struct {
	AuthService  AuthService
	OAuthService OAuthService
}

func (f *Factory) NewServices() *Services {
	return &Services{
		AuthService:  NewAuthService(f.db),
		OAuthService: NewOAuthService(f.db, f.googleProvider),
	}
}
