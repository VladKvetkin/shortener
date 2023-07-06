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
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	db := &sqlx.DB{}
	if config.DatabaseDSN != "" {
		db, err = sqlx.Connect("postgres", config.DatabaseDSN)
		if err != nil {
			panic(err)
		}
	}

	router := router.NewRouter(handler.NewHandler(storage.NewStorage(storage.NewPersister(config.FileStoragePath)), config, db))
	server := server.NewServer(config, router.Router)

	zap.L().Info("Running server", zap.String("Address", config.Address))

	err = server.Start()
	if err != nil {
		panic(err)
	}
}
