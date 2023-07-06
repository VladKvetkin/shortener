package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/handler"
	"github.com/VladKvetkin/shortener/internal/app/router"
	"github.com/VladKvetkin/shortener/internal/app/server"
	"github.com/VladKvetkin/shortener/internal/app/storage"
	"go.uber.org/zap"
)

func main() {
	var dataStorage storage.Storage

	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	if config.DatabaseDSN != "" {
		db, err := sqlx.Connect("postgres", config.DatabaseDSN)
		if err != nil {
			panic(err)
		}

		defer db.Close()

		dataStorage, err = storage.NewPostgresStorage(db)
		if err != nil {
			panic(err)
		}
	} else {
		dataStorage = storage.NewMemStorage(storage.NewPersister(config.FileStoragePath))
	}

	router := router.NewRouter(handler.NewHandler(dataStorage, config))
	server := server.NewServer(config, router.Router)

	zap.L().Info("Running server", zap.String("Address", config.Address))

	err = server.Start()
	if err != nil {
		panic(err)
	}
}
