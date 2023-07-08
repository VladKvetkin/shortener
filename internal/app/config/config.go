package config

import (
	"flag"
	"net/url"
	"strings"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	Address             string `env:"SERVER_ADDRESS"`
	BaseShortURLAddress string `env:"BASE_URL"`
	FileStoragePath     string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN         string `env:"DATABASE_DSN"`
}

func NewConfig() (Config, error) {
	config := Config{
		Address:             "localhost:8080",
		BaseShortURLAddress: "http://localhost:8080/",
		FileStoragePath:     "/tmp/short-url-db.json",
	}

	config.parseFlags()

	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	if err := config.validateConfig(); err != nil {
		return Config{}, err
	}

	config.normalizeConfig()

	return config, nil
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.Address, "a", c.Address, "HTTP server address")
	flag.StringVar(&c.BaseShortURLAddress, "b", c.BaseShortURLAddress, "Base address for short URL")
	flag.StringVar(&c.FileStoragePath, "f", c.FileStoragePath, "File storage path for short URLs")
	flag.StringVar(&c.DatabaseDSN, "d", c.DatabaseDSN, "Database data source name")

	flag.Parse()
}

func (c *Config) validateConfig() error {
	for _, URI := range []string{c.Address, c.BaseShortURLAddress} {
		_, err := url.ParseRequestURI(URI)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) normalizeConfig() {
	c.BaseShortURLAddress = strings.TrimRight(c.BaseShortURLAddress, "/")
}
