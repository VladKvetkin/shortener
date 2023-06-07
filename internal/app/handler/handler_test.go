package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/handler"
	"github.com/VladKvetkin/shortener/internal/app/router"
	"github.com/VladKvetkin/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouterPostHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        *regexp.Regexp
	}

	tests := []struct {
		name    string
		request string
		method  string
		body    string
		storage storage.Repositories
		config  config.Config
		headers map[string]string
		want    want
	}{
		{
			name:    "post request without body",
			request: "/",
			method:  http.MethodPost,
			body:    "",
			storage: storage.NewStorage(),
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        regexp.MustCompile(`^Invalid request\s*$`),
			},
		},
		{
			name:    "post request with body",
			request: "/",
			method:  http.MethodPost,
			body:    "https://practicum.yandex.ru/",
			storage: storage.NewStorage(),
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				body:        regexp.MustCompile(`^http://localhost/.{8}$`),
			},
		},
		{
			name:    "post request with body short URL already exists",
			request: "/",
			method:  http.MethodPost,
			body:    "https://practicum.yandex.ru/",
			storage: func() *storage.Storage {
				storage := storage.NewStorage()
				storage.Add("EwHXdJfB", "https://practicum.yandex.ru/")
				return storage
			}(),
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				body:        regexp.MustCompile(`^http://localhost/EwHXdJfB$`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))

			for header, value := range tt.headers {
				request.Header.Add(header, value)
			}

			recorder := httptest.NewRecorder()
			router := router.NewRouter(handler.NewHandler(tt.storage, tt.config))

			router.Router.ServeHTTP(recorder, request)

			result := recorder.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Regexp(t, tt.want.body, string(body))
		})
	}
}

func TestRouterGetHandler(t *testing.T) {
	type want struct {
		location   string
		statusCode int
		body       *regexp.Regexp
	}

	tests := []struct {
		name    string
		request string
		method  string
		body    string
		storage storage.Repositories
		config  config.Config
		headers map[string]string
		want    want
	}{
		{
			name:    "get request without short URL",
			request: "/",
			method:  http.MethodGet,
			storage: storage.NewStorage(),
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				location:   "",
				body:       regexp.MustCompile(`^$`),
			},
		},
		{
			name:    "get request with short URL, which not in storage",
			request: "/EwHXdJfB",
			method:  http.MethodGet,
			storage: storage.NewStorage(),
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				location:   "",
				body:       regexp.MustCompile(`^Invalid request\s*$`),
			},
		},
		{
			name:    "get request with short URL, which in storage",
			request: "/EwHXdJfB",
			method:  "GET",
			body:    "",
			storage: func() *storage.Storage {
				storage := storage.NewStorage()
				storage.Add("EwHXdJfB", "https://practicum.yandex.ru/")
				return storage
			}(),
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				statusCode: 307,
				location:   "https://practicum.yandex.ru/",
				body:       regexp.MustCompile(`^$`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			for header, value := range tt.headers {
				request.Header.Add(header, value)
			}

			recorder := httptest.NewRecorder()
			router := router.NewRouter(handler.NewHandler(tt.storage, tt.config))

			router.Router.ServeHTTP(recorder, request)

			result := recorder.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Regexp(t, tt.want.body, string(body))
		})
	}
}