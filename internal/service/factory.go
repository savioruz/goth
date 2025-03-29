package service

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/savioruz/goth/pkg/logger"
)

type Factory struct {
	db *pgxpool.Pool
	l  logger.Interface
}

func NewFactory(db *pgxpool.Pool, l logger.Interface) *Factory {
	return &Factory{db, l}
}

type Services struct {
	AuthService AuthService
}

func (f *Factory) NewServices() *Services {
	return &Services{
		AuthService: NewAuthService(f.db),
	}
}
