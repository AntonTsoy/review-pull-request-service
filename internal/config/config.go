package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Database struct {
	Host     string `env:"DB_HOST" env-default:"db"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	Name     string `env:"DB_NAME" env-default:"postgres"`
	SSL      string `env:"DB_SSL" env-default:"disable"`
	Pool     int    `env:"DB_POOL" env-default:"10"`
}

type Config struct {
	Host string `env:"HOST" env-default:""`
	Port string `env:"PORT" env-default:"8080"`
	DB   Database
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}
