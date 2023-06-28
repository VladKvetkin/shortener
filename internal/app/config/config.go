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
}

func NewConfig() (Config, error) {
	config := Config{
		Address:             "localhost:8080",
		BaseShortURLAddress: "http://localhost:8080/",
		FileStoragePath:     "/tmp/short-url-db.json",
	}

	flag.StringVar(&config.Address, "a", config.Address, "HTTP server address")
	flag.StringVar(&config.BaseShortURLAddress, "b", config.BaseShortURLAddress, "Base address for short URL")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "File storage path for short URLs")

	flag.Parse()

	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	if err := config.validateConfig(); err != nil {
		return Config{}, err
	}

	config.normalizeConfig()

	return config, nil
}

func (c *Config) validateConfig() error {
	_, err := url.ParseRequestURI(c.Address)
	if err != nil {
		return err
	}

	_, err = url.ParseRequestURI(c.BaseShortURLAddress)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) normalizeConfig() {
	c.BaseShortURLAddress = strings.TrimRight(c.BaseShortURLAddress, "/")
}
