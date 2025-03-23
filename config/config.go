package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App  App
		HTTP HTTP
		Log  Log
		Pg   Pg
	}

	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required" envDefault:"info"`
	}

	Pg struct {
		PoolMax  int    `env:"PG_POOL_MAX,required"`
		Host     string `env:"PG_HOST,required"`
		Port     int    `env:"PG_PORT,required"`
		User     string `env:"PG_USER"`
		Password string `env:"PG_PASSWORD"`
		Dbname   string `env:"PG_DATABASE,required"`
		SSLMode  string `env:"PG_SSLMODE,required"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config failed: %w", err)
	}

	return cfg, nil
}
