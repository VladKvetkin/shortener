package router

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/shortener"
	"github.com/VladKvetkin/shortener/internal/app/storage"
	"github.com/go-chi/chi"
)

type Router struct {
	Router  *chi.Mux
	storage storage.Repositories
	config  config.Config
}

func NewRouter(storage storage.Repositories, config config.Config) *Router {
	chiRouter := chi.NewRouter()

	router := &Router{
		Router:  chiRouter,
		config:  config,
		storage: storage,
	}

	chiRouter.Route("/", func(r chi.Router) {
		r.Post("/", router.PostHandler)
		r.Get("/{id}", router.GetHandler)
	})

	return router
}

func (r *Router) GetHandler(res http.ResponseWriter, req *http.Request) {
	shortURL := chi.URLParam(req, "id")
	if shortURL == "" {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	url, err := r.storage.ReadByShortURL(shortURL)
	if err != nil {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (r *Router) PostHandler(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	if string(body) == "" {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	shortURL := r.storage.ReadByURL(string(body))
	if shortURL == "" {
		shortURL = shortener.CreateShortURL()
		r.storage.Add(shortURL, string(body))
	}

	baseShortURLAddress := strings.TrimRight(r.config.BaseShortURLAddress, "/")

	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("%s/%s", baseShortURLAddress, shortURL)))
}
