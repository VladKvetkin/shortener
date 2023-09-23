package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/handler"
	"github.com/VladKvetkin/shortener/internal/app/storage"
)

func ExampleHandler_APIShortenHandler() {
	defaultStorage, err := storage.GetStorage(config.Config{})
	if err != nil {
		panic(err)
	}

	config := config.Config{
		Address:             "localhost:8080",
		BaseShortURLAddress: "http://localhost",
	}

	handler := handler.NewHandler(defaultStorage, config)

	recorder := httptest.NewRecorder()

	request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url": "https://practicum.yandex.ru"}`))

	handler.APIShortenHandler(recorder, request)
}

func ExampleHandler_PostHandler() {
	defaultStorage, err := storage.GetStorage(config.Config{})
	if err != nil {
		panic(err)
	}

	config := config.Config{
		Address:             "localhost:8080",
		BaseShortURLAddress: "http://localhost",
	}

	handler := handler.NewHandler(defaultStorage, config)

	recorder := httptest.NewRecorder()

	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`https://practicum.yandex.ru/`))

	handler.PostHandler(recorder, request)
}

func ExampleHandler_GetHandler() {
	defaultStorage, err := storage.GetStorage(config.Config{})
	if err != nil {
		panic(err)
	}

	config := config.Config{
		Address:             "localhost:8080",
		BaseShortURLAddress: "http://localhost",
	}

	handler := handler.NewHandler(defaultStorage, config)

	recorder := httptest.NewRecorder()

	request := httptest.NewRequest(http.MethodGet, "/EwHXdJfB", strings.NewReader(``))

	handler.PostHandler(recorder, request)
}
