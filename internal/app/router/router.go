package router

import (
	"fmt"
	"io"
	"net/http"

	"github.com/VladKvetkin/shortener/internal/app/shortener"
	"github.com/VladKvetkin/shortener/internal/app/storage"
	"github.com/go-chi/chi"
)

func NewRouter(storage storage.Repositories) http.Handler {
	router := chi.NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
			postHandler(writer, request, storage)
		})

		r.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
			getHandler(writer, request, storage)
		})
	})

	return router
}

func getHandler(res http.ResponseWriter, req *http.Request, storage storage.Repositories) {
	shortURL := chi.URLParam(req, "id")
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
