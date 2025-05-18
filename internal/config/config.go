package config

// Package config provides the configuration structure and initialization logic.
// It uses the cleanenv library to load configuration from environment variables or YAML files.

import (
	"TestRest/pkg/postgres"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres postgres.Config `yaml:"POSTGRES" env:"POSTGRES"`

	RESTHost string `yaml:"REST_HOST" env:"REST_HOST" env-default:"localhost"`
	RESTPort int    `yaml:"REST_PORT" env:"REST_PORT" env-default:"8080"`
}

// New initializes the application configuration.
// @Summary Initialize configuration
// @Description Loads the application configuration from environment variables or YAML files.
// @Tags config
// @Success 200 {object} config.Config "Application configuration"
// @Failure 500 {string} string "Failed to load configuration"
// @Router /config/new [get]
func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
