package config

import (
	"context"

	"backend/internal/config/types"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Database types.DatabaseConfig `env:",prefix=DB_"`
}

func Load() (*Config, error) {
	ctx := context.Background()
	var config Config

	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func MustLoad() *Config {
	config, err := Load()
	if err != nil {
		panic("failed to load configuration: " + err.Error())
	}
	return config
}
