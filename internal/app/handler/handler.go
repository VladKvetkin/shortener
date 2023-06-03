package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/VladKvetkin/shortener/internal/app/shortener"
	"github.com/VladKvetkin/shortener/internal/app/storage"
)

func MainHandler(storage storage.Repositories) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			getHandler(res, req, storage)
		} else if req.Method == http.MethodPost {
			postHandler(res, req, storage)
		} else {
			http.Error(res, "Invalid request", http.StatusBadRequest)
		}
	})
}

func getHandler(res http.ResponseWriter, req *http.Request, storage storage.Repositories) {
	shortURL := strings.TrimLeft(req.URL.Path, "/")

	if shortURL == "" {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	url, err := storage.ReadByShortURL(shortURL)
	if err != nil {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func postHandler(res http.ResponseWriter, req *http.Request, storage storage.Repositories) {
	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	if string(body) == "" {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	shortURL := storage.ReadByURL(string(body))
	if shortURL == "" {
		shortURL = shortener.CreateShortURL()
		storage.Add(shortURL, string(body))
	}

	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}

	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("%s://%s/%s", scheme, req.Host, shortURL)))
}
