package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT" default:"8080"`
	Secret     string `env:"SECRET"`
	Database   DatabaseConfig
}

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST"`
	Name     string `env:"DATABASE_NAME"`
	Port     string `env:"DATABASE_PORT"`
	User     string `env:"DATABASE_USER"`
	Password string `env:"DATABASE_PASSWORD"`
}

func Load() *Config {
	cfgApp := &Config{}
	parseConfig(cfgApp)
	return cfgApp
}

func parseConfig(cfg *Config) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	err := env.Parse(cfg)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}
}
