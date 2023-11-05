// Package config отвечает за конфигурацию приложения.
// Конфигурировать приложение можно как флагами в командной строке, так и через переменные окружения.

package config

import (
	"encoding/json"
	"flag"
	"net/url"
	"os"
	"strings"

	"dario.cat/mergo"

	"github.com/caarlos0/env/v8"
)

// Config - структура конфига, содержит в себе настройки приложения.
type Config struct {
	// Address - адрес, на котором запускается приложение.
	Address string `env:"SERVER_ADDRESS" json:"server_address"`
	// BaseShortURLAddress - адрес, который используется для генерации сокращенной ссылки.
	BaseShortURLAddress string `env:"BASE_URL" json:"base_url"`
	// FileStoragePath - путь к файлу, который используется для сохранения сокращенных ссылок.
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	// DatabaseDSN - DSN для базы данных.
	DatabaseDSN string `env:"DATABASE_DSN" json:"database_dsn"`
	// EnableHTTPS - запускает сервер с поддержкой HTTPS
	EnableHTTPS bool `env:"ENABLE_HTTPS" json:"enable_https"`
	// ConfigPath - путь к файлу JSON-конфигурации
	ConfigPath string `env:"CONFIG"`
}

// NewConfig – конструктор Config.
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

	if config.ConfigPath != "" {
		if err := config.parseConfig(&config); err != nil {
			return Config{}, err
		}
	}

	if err := config.validateConfig(); err != nil {
		return Config{}, err
	}

	config.normalizeConfig()

	return config, nil
}

func (c *Config) parseConfig(config *Config) error {
	tmpConf := &Config{}

	configFile, err := os.Open(config.ConfigPath)
	if err != nil {
		return err
	}

	err = json.NewDecoder(configFile).Decode(tmpConf)
	if err != nil {
		return err
	}

	if err = mergo.Merge(config, tmpConf); err != nil {
		return err
	}

	return nil
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.Address, "a", c.Address, "HTTP server address")
	flag.StringVar(&c.BaseShortURLAddress, "b", c.BaseShortURLAddress, "Base address for short URL")
	flag.StringVar(&c.FileStoragePath, "f", c.FileStoragePath, "File storage path for short URLs")
	flag.StringVar(&c.DatabaseDSN, "d", c.DatabaseDSN, "Database data source name")
	flag.BoolVar(&c.EnableHTTPS, "s", c.EnableHTTPS, "Enable HTTPS")
	flag.StringVar(&c.ConfigPath, "c", c.ConfigPath, "JSON config path")
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
