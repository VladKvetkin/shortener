package config

import (
	"errors"
	"flag"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	ErrorInvalidURLAddress = errors.New("need valid address for short URL in a form scheme://host:port/")
)

type Config struct {
	Host                string
	Port                int
	BaseShortURLAddress string
}

func NewConfig() (Config, error) {
	config := Config{}

	address, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		host, port, err := config.parseAddress(address)
		if err != nil {
			return config, err
		}

		config.Host = host
		config.Port = port
	}

	flag.Func("a", "HTTP server address", func(address string) error {
		if config.Host != "" && config.Port != 0 {
			return nil
		}

		host, port, err := config.parseAddress(address)
		if err != nil {
			return err
		}

		config.Host = host
		config.Port = port

		return nil
	})

	baseAddressForShortURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		_, err := url.ParseRequestURI(baseAddressForShortURL)
		if err != nil {
			return config, ErrorInvalidURLAddress
		}

		config.BaseShortURLAddress = baseAddressForShortURL
	}

	flag.Func("b", "Base address for short URL", func(flagValue string) error {
		if config.BaseShortURLAddress != "" {
			return nil
		}

		_, err := url.ParseRequestURI(flagValue)
		if err != nil {
			return ErrorInvalidURLAddress
		}

		config.BaseShortURLAddress = flagValue

		return nil
	})

	flag.Parse()

	if config.Host == "" {
		config.Host = "localhost"
	}

	if config.Port == 0 {
		config.Port = 8080
	}

	if config.BaseShortURLAddress == "" {
		config.BaseShortURLAddress = "http://localhost:8080/"
	}

	return config, nil
}

func (c *Config) GetAddress() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func (c *Config) parseAddress(address string) (host string, port int, err error) {
	splitAddress := strings.Split(address, ":")
	if len(splitAddress) != 2 {
		return "", 0, errors.New("need HTTP server address in a form host:port")
	}

	port, err = strconv.Atoi(splitAddress[1])
	if err != nil {
		return "", 0, err
	}

	return splitAddress[0], port, nil
}
