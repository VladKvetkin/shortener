package config

import (
	"flag"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Address             string
	BaseShortURLAddress string
}

func NewConfig() (Config, error) {
	config := Config{
		Address:             "localhost:8080",
		BaseShortURLAddress: "http://localhost:8080/",
	}

	flag.StringVar(&config.Address, "a", config.Address, "HTTP server address")
	flag.StringVar(&config.BaseShortURLAddress, "b", config.BaseShortURLAddress, "Base address for short URL")

	flag.Parse()

	if address, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		config.Address = address
	}

	if baseAddressForShortURL, ok := os.LookupEnv("BASE_URL"); ok {
		config.BaseShortURLAddress = baseAddressForShortURL
	}

	err := config.validateConfig()
	if err != nil {
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
