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
		err := config.setAddress(address)
		if err != nil {
			return config, err
		}
	}

	flag.Func("a", "HTTP server address", func(address string) error {
		if config.Host != "" && config.Port != 0 {
			return nil
		}

		return config.setAddress(address)
	})

	baseAddressForShortURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		err := config.setBaseShortURL(baseAddressForShortURL)
		if err != nil {
			return config, err
		}
	}

	flag.Func("b", "Base address for short URL", func(flagValue string) error {
		if config.BaseShortURLAddress != "" {
			return nil
		}

		return config.setBaseShortURL(flagValue)
	})

	flag.Parse()

	config.setDefaultValues()

	return config, nil
}

func (c *Config) setAddress(address string) error {
	host, port, err := c.parseAddress(address)
	if err != nil {
		return err
	}

	c.Host = host
	c.Port = port

	return nil
}

func (c *Config) setBaseShortURL(baseShortURL string) error {
	_, err := url.ParseRequestURI(baseShortURL)
	if err != nil {
		return ErrorInvalidURLAddress
	}

	c.BaseShortURLAddress = baseShortURL

	return nil
}

func (c *Config) setDefaultValues() {
	if c.Host == "" {
		c.Host = "localhost"
	}

	if c.Port == 0 {
		c.Port = 8080
	}

	if c.BaseShortURLAddress == "" {
		c.BaseShortURLAddress = "http://localhost:8080/"
	}
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
