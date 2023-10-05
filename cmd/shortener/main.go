package main

import (
	_ "github.com/lib/pq"

	"go.uber.org/zap"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/handler"
	"github.com/VladKvetkin/shortener/internal/app/router"
	"github.com/VladKvetkin/shortener/internal/app/server"
	"github.com/VladKvetkin/shortener/internal/app/storage"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	storage, err := storage.GetStorage(config)
	if err != nil {
		panic(err)
	}

	defer storage.Close()

	router := router.NewRouter(handler.NewHandler(storage, config))
	server := server.NewServer(config, router.Router)

	zap.L().Info("Running server", zap.String("Address", config.Address))

	err = server.Start()
	if err != nil {
		panic(err)
	}
}
