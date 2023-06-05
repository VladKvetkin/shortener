package config

import (
	"errors"
	"flag"
	"net/url"
	"strconv"
	"strings"
)

type Config struct {
	Host                string
	Port                int
	BaseShortURLAddress string
}

func NewConfig() Config {
	config := Config{
		Host:                "localhost",
		Port:                8080,
		BaseShortURLAddress: "http://localhost:8080/",
	}

	flag.Func("a", "HTTP server address", func(flagValue string) error {
		splitAddress := strings.Split(flagValue, ":")
		if len(splitAddress) != 2 {
			return errors.New("need HTTP server address in a form host:port")
		}

		port, err := strconv.Atoi(splitAddress[1])
		if err != nil {
			return err
		}

		config.Host = splitAddress[0]
		config.Port = port

		return nil
	})

	flag.Func("b", "Base address for short URL", func(flagValue string) error {
		_, err := url.ParseRequestURI(flagValue)
		if err != nil {
			return errors.New("need valid address for short URL in a form scheme://host:port/")
		}

		config.BaseShortURLAddress = flagValue

		return nil
	})

	flag.Parse()

	return config
}

func (c *Config) GetAddress() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func (c *Config) GetBaseShortURLAddress() string {
	return c.BaseShortURLAddress
}
