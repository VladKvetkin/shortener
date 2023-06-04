package main

import (
	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/router"
	"github.com/VladKvetkin/shortener/internal/app/server"
	"github.com/VladKvetkin/shortener/internal/app/storage"
)

func main() {
	config := config.NewConfig("localhost", "8080")
	storage := storage.NewStorage()
	router := router.NewRouter(storage)

	server := server.NewServer(config, router)

	err := server.Start()
	if err != nil {
		panic(err)
	}
}
