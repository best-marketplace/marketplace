package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	ProxyPort     string `env:"PROXY_PORT" default:"8000"`
	ServiceRoutes ServiceRoutesConfig
}

type ServiceRoutesConfig struct {
	CatalogService string `env:"CATALOG_SERVICE_URL" default:"http://backend:8080"`
}

func Load() *Config {
	cfgApp := &Config{}
	parseConfig(cfgApp)
	return cfgApp
}

func parseConfig(cfg *Config) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	err := env.Parse(cfg)
	if err != nil {
		log.Fatalf("unable to parse environment variables: %v", err)
	}
}

func (c *Config) GetRoutes() map[string]string {
	return map[string]string{
		"/search":   c.ServiceRoutes.CatalogService,
		"/products": c.ServiceRoutes.CatalogService,
		"/product":  c.ServiceRoutes.CatalogService,
		"/comment":  c.ServiceRoutes.CatalogService,
		"/comments": c.ServiceRoutes.CatalogService,
	}
}
