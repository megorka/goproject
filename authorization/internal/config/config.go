package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/megorka/goproject/authorization/internal/transport/http"
	"github.com/megorka/goproject/authorization/pkg/postgres"
)

type Config struct {
	Router   router.Config   `yaml:"ROUTER"`
	Postgres postgres.Config `yaml:"POSTGRES"`
}

func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig("/app/config/config.yaml", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
