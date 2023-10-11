package main

import (
	"fmt"

	_ "github.com/lib/pq"

	"go.uber.org/zap"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/handler"
	"github.com/VladKvetkin/shortener/internal/app/router"
	"github.com/VladKvetkin/shortener/internal/app/server"
	"github.com/VladKvetkin/shortener/internal/app/storage"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %v\n", buildVersion)
	fmt.Printf("Build date: %v\n", buildDate)
	fmt.Printf("Build commit: %v\n", buildCommit)

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
