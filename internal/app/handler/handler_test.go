package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/entities"
	"github.com/VladKvetkin/shortener/internal/app/handler"
	"github.com/VladKvetkin/shortener/internal/app/router"
	"github.com/VladKvetkin/shortener/internal/app/storage"
)

func TestRouterPostHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        *regexp.Regexp
	}

	defaultStorage, err := storage.GetStorage(config.Config{})
	if err != nil {
		panic(err)
	}
	shortURLAlreadyExistStorage, err := storage.GetStorage(config.Config{})
	if err != nil {
		panic(err)
	}

	shortURLAlreadyExistStorage.Add(entities.URL{
		ShortURL:    "QrPnX5IU",
		OriginalURL: "https://practicum.yandex.ru/",
	})

	tests := []struct {
		name    string
		request string
		method  string
		body    string
		storage storage.Storage
		config  config.Config
		headers map[string]string
		want    want
	}{
		{
			name:    "post request without body",
			request: "/",
			method:  http.MethodPost,
			body:    "",
			storage: defaultStorage,
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
			storage: defaultStorage,
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
			storage: shortURLAlreadyExistStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusConflict,
				body:        regexp.MustCompile(`^http://localhost/QrPnX5IU`),
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

func BenchmarkRouterPostHandler(b *testing.B) {
	type want struct {
		contentType string
		statusCode  int
		body        *regexp.Regexp
	}

	defaultStorage, err := storage.GetStorage(config.Config{})
	if err != nil {
		panic(err)
	}
	shortURLAlreadyExistStorage, err := storage.GetStorage(config.Config{})
	if err != nil {
		panic(err)
	}

	shortURLAlreadyExistStorage.Add(entities.URL{
		ShortURL:    "QrPnX5IU",
		OriginalURL: "https://practicum.yandex.ru/",
	})

	tests := []struct {
		name    string
		request string
		method  string
		body    string
		storage storage.Storage
		config  config.Config
		headers map[string]string
		want    want
	}{
		{
			name:    "post request without body",
			request: "/",
			method:  http.MethodPost,
			body:    "",
			storage: defaultStorage,
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
			storage: defaultStorage,
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
			storage: shortURLAlreadyExistStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusConflict,
				body:        regexp.MustCompile(`^http://localhost/QrPnX5IU`),
			},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))

				for header, value := range tt.headers {
					request.Header.Add(header, value)
				}

				recorder := httptest.NewRecorder()
				router := router.NewRouter(handler.NewHandler(tt.storage, tt.config))

				b.StartTimer()

				router.Router.ServeHTTP(recorder, request)

				result := recorder.Result()
				result.Body.Close()
			}
		})
	}
}

func TestRouterGetHandler(t *testing.T) {
	type want struct {
		location   string
		statusCode int
		body       *regexp.Regexp
	}

	defaultStorage, err := storage.GetStorage(config.Config{})
	if err != nil {
		panic(err)
	}
	shortURLAlreadyExistStorage, err := storage.GetStorage(config.Config{})
	if err != nil {
		panic(err)
	}

	shortURLAlreadyExistStorage.Add(entities.URL{
		ShortURL:    "EwHXdJfB",
		OriginalURL: "https://practicum.yandex.ru/",
	})

	tests := []struct {
		name    string
		request string
		method  string
		body    string
		storage storage.Storage
		config  config.Config
		headers map[string]string
		want    want
	}{
		{
			name:    "get request without short URL",
			request: "/",
			method:  http.MethodGet,
			storage: defaultStorage,
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
			storage: defaultStorage,
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
			storage: shortURLAlreadyExistStorage,
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

func BenchmarkRouterGetHandler(b *testing.B) {
	type want struct {
		location   string
		statusCode int
		body       *regexp.Regexp
	}

	defaultStorage, err := storage.GetStorage(config.Config{})
	if err != nil {
		panic(err)
	}
	shortURLAlreadyExistStorage, err := storage.GetStorage(config.Config{})
	if err != nil {
		panic(err)
	}

	shortURLAlreadyExistStorage.Add(entities.URL{
		ShortURL:    "EwHXdJfB",
		OriginalURL: "https://practicum.yandex.ru/",
	})

	tests := []struct {
		name    string
		request string
		method  string
		body    string
		storage storage.Storage
		config  config.Config
		headers map[string]string
		want    want
	}{
		{
			name:    "get request without short URL",
			request: "/",
			method:  http.MethodGet,
			storage: defaultStorage,
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
			storage: defaultStorage,
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
			storage: shortURLAlreadyExistStorage,
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
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
				for header, value := range tt.headers {
					request.Header.Add(header, value)
				}

				recorder := httptest.NewRecorder()
				router := router.NewRouter(handler.NewHandler(tt.storage, tt.config))

				b.StartTimer()

				router.Router.ServeHTTP(recorder, request)

				result := recorder.Result()
				result.Body.Close()
			}
		})
	}
}

func TestRouterAPIShortenHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
	}

	defaultStorage, err := storage.GetStorage(config.Config{})
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		request string
		method  string
		body    string
		storage storage.Storage
		config  config.Config
		headers map[string]string
		want    want
	}{
		{
			name:    "post request without body",
			request: "/api/shorten",
			method:  http.MethodPost,
			storage: defaultStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			body: "",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "Cannot decode request JSON body\n",
			},
		},
		{
			name:    "post request without URL in body",
			request: "/api/shorten",
			method:  http.MethodPost,
			storage: defaultStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body: `{"url": ""}`,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "Invalid request\n",
			},
		},
		{
			name:    "post request with URL",
			request: "/api/shorten",
			method:  http.MethodPost,
			storage: defaultStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body: `{"url": "https://practicum.yandex.ru"}`,
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "application/json",
				body: `{"result":"http://localhost/ipkjUVtE"}
`,
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

			assert.Equal(t, tt.want.body, string(body))
		})
	}
}

func BenchmarkRouterAPIShortenHandler(b *testing.B) {
	type want struct {
		contentType string
		statusCode  int
		body        string
	}

	defaultStorage, err := storage.GetStorage(config.Config{})
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		request string
		method  string
		body    string
		storage storage.Storage
		config  config.Config
		headers map[string]string
		want    want
	}{
		{
			name:    "post request without body",
			request: "/api/shorten",
			method:  http.MethodPost,
			storage: defaultStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			body: "",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "Cannot decode request JSON body\n",
			},
		},
		{
			name:    "post request without URL in body",
			request: "/api/shorten",
			method:  http.MethodPost,
			storage: defaultStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body: `{"url": ""}`,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "Invalid request\n",
			},
		},
		{
			name:    "post request with URL",
			request: "/api/shorten",
			method:  http.MethodPost,
			storage: defaultStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body: `{"url": "https://practicum.yandex.ru"}`,
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "application/json",
				body: `{"result":"http://localhost/ipkjUVtE"}
`,
			},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
				for header, value := range tt.headers {
					request.Header.Add(header, value)
				}

				recorder := httptest.NewRecorder()
				router := router.NewRouter(handler.NewHandler(tt.storage, tt.config))
				b.StartTimer()

				router.Router.ServeHTTP(recorder, request)

				result := recorder.Result()
				result.Body.Close()
			}
		})
	}
}

func TestRouterDeleteUserUrlsHandler(t *testing.T) {
	type want struct {
		statusCode int
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	deleteUrls := []string{"6qxTVvsy", "RTfd56hn", "Jlfd67ds"}

	mockStorage := storage.NewMockStorage(ctrl)
	mockStorage.EXPECT().DeleteBatch(gomock.Any(), deleteUrls, gomock.Any()).Return(nil).MinTimes(0)

	tests := []struct {
		name    string
		request string
		method  string
		body    string
		storage storage.Storage
		config  config.Config
		headers map[string]string
		want    want
	}{
		{
			name:    "delete request with empty body",
			request: "/api/user/urls",
			method:  http.MethodDelete,
			storage: mockStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			body: "",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "delete request with empty urls",
			request: "/api/user/urls",
			method:  http.MethodDelete,
			storage: mockStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			body: `[]`,
			want: want{
				statusCode: http.StatusAccepted,
			},
		},
		{
			name:    "delete request with urls",
			request: "/api/user/urls",
			method:  http.MethodDelete,
			storage: mockStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			body: `["6qxTVvsy", "RTfd56hn", "Jlfd67ds"]`,
			want: want{
				statusCode: http.StatusAccepted,
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
			result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func BenchmarkRouterDeleteUserUrlsHandler(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	deleteUrls := []string{"6qxTVvsy", "RTfd56hn", "Jlfd67ds"}

	mockStorage := storage.NewMockStorage(ctrl)
	mockStorage.EXPECT().DeleteBatch(gomock.Any(), deleteUrls, gomock.Any()).Return(nil).MinTimes(0)

	tests := []struct {
		name    string
		request string
		method  string
		body    string
		storage storage.Storage
		config  config.Config
		headers map[string]string
	}{
		{
			name:    "delete request with empty body",
			request: "/api/user/urls",
			method:  http.MethodDelete,
			storage: mockStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			body: "",
		},
		{
			name:    "delete request with empty urls",
			request: "/api/user/urls",
			method:  http.MethodDelete,
			storage: mockStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			body: `[]`,
		},
		{
			name:    "delete request with urls",
			request: "/api/user/urls",
			method:  http.MethodDelete,
			storage: mockStorage,
			config: config.Config{
				Address:             "localhost:8080",
				BaseShortURLAddress: "http://localhost",
			},
			headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			body: `["6qxTVvsy", "RTfd56hn", "Jlfd67ds"]`,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
				for header, value := range tt.headers {
					request.Header.Add(header, value)
				}

				recorder := httptest.NewRecorder()
				router := router.NewRouter(handler.NewHandler(tt.storage, tt.config))
				b.StartTimer()

				router.Router.ServeHTTP(recorder, request)

				result := recorder.Result()
				result.Body.Close()
			}
		})
	}
}
