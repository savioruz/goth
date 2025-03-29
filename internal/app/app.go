package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/savioruz/goth/config"
	v1 "github.com/savioruz/goth/internal/delivery/http"
	"github.com/savioruz/goth/internal/service"
	"github.com/savioruz/goth/pkg/httpserver"
	"github.com/savioruz/goth/pkg/jwt"
	"github.com/savioruz/goth/pkg/logger"
	"github.com/savioruz/goth/pkg/oauth"
	"github.com/savioruz/goth/pkg/postgres"
	"github.com/savioruz/goth/pkg/redis"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// JWT
	jwt.Initialize(cfg.App.Name, cfg.JWT.Secret, jwt.ParseDuration(cfg.JWT.AccessTokenExpiry), jwt.ParseDuration(cfg.JWT.RefreshTokenExpiry))

	// Postgres
	dsn := postgres.ConnectionBuilder(cfg.Pg.Host, cfg.Pg.Port, cfg.Pg.User, cfg.Pg.Password, cfg.Pg.Dbname, cfg.Pg.SSLMode)
	pg, err := postgres.New(dsn, postgres.MaxPoolSize(cfg.Pg.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Pool.Close()

	if err := pg.Ping(context.Background()); err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.Ping: %w", err))
	}

	// Redis
	addr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	redis, err := redis.New(addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - redis.New: %w", err))
	}
	defer redis.Close()

	if err := redis.Ping(context.Background()); err != nil {
		l.Fatal(fmt.Errorf("app - Run - redis.Ping: %w", err))
	}

	// OAuth
	googleProvider := oauth.NewGoogleProvider(cfg.OAuth.Google.ClientID, cfg.OAuth.Google.ClientSecret, cfg.OAuth.Google.RedirectURL)

	// Service
	serviceFactory := service.NewFactory(pg.Pool, l, googleProvider)
	services := serviceFactory.NewServices()

	// HTTP Server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port))
	v1.NewRouter(httpServer.App, cfg, l, *services)

	// Start Server
	httpServer.Start()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown Server
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
