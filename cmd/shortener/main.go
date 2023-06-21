package main

import (
	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/handler"
	"github.com/VladKvetkin/shortener/internal/app/logger"
	"github.com/VladKvetkin/shortener/internal/app/router"
	"github.com/VladKvetkin/shortener/internal/app/server"
	"github.com/VladKvetkin/shortener/internal/app/storage"
	"go.uber.org/zap"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	if err := logger.Initialize(config.LogLevel); err != nil {
		panic(err)
	}

	router := router.NewRouter(handler.NewHandler(storage.NewStorage(), config))
	server := server.NewServer(config, router.Router)

	logger.Log.Info("Running server", zap.String("Address", config.Address))

	err = server.Start()
	if err != nil {
		panic(err)
	}
}
