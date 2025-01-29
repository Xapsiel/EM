package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseConfig
	HostConfig
	APIConfig
}
type DatabaseConfig struct {
	Host     string `env:"db_host"`
	Port     int    `env:"db_port"`
	User     string `env:"db_user"`
	Password string `env:"db_password"`
	Name     string `env:"db_name"`
	Sslmode  string `env:"db_sslmode"`
}
type HostConfig struct {
	Port string `env:"host_port"`
}
type APIConfig struct {
	Domain string `env:"domain"`
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
