package config

import (
	"TestRest/pkg/postgres"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres postgres.Config `yaml:"POSTGRES" env:"POSTGRES"`

	RESTHost string `yaml:"REST_HOST" env:"REST_HOST" env-default:"localhost"`
	RESTPort int    `yaml:"REST_PORT" env:"REST_PORT" env-default:"8080"`

	ExternalAPIs struct {
		AgeURL         string `yaml:"AGE_API_URL" env:"AGE_API_URL" env-default:"https://api.agify.io"`
		GenderURL      string `yaml:"GENDER_API_URL" env:"GENDER_API_URL" env-default:"https://api.genderize.io"`
		NationalityURL string `yaml:"NATIONALITY_API_URL" env:"NATIONALITY_API_URL" env-default:"https://api.nationalize.io"`
	} `yaml:"EXTERNAL_APIS"`
}

func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
