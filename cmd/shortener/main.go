package main

import (
	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/router"
	"github.com/VladKvetkin/shortener/internal/app/server"
	"github.com/VladKvetkin/shortener/internal/app/storage"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	storage := storage.NewStorage()
	router := router.NewRouter(storage, config)

	server := server.NewServer(config, router.Router)

	err = server.Start()
	if err != nil {
		panic(err)
	}
}
