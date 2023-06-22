package router

import (
	"net/http"

	"github.com/VladKvetkin/shortener/internal/app/handler"
	"github.com/VladKvetkin/shortener/internal/app/middleware"
	"github.com/go-chi/chi"
)

type Router struct {
	Router  *chi.Mux
	handler *handler.Handler
}

func NewRouter(handler *handler.Handler) *Router {
	chiRouter := chi.NewRouter()

	router := &Router{
		Router:  chiRouter,
		handler: handler,
	}

	chiRouter.Use(middleware.Logger)
	chiRouter.Route("/", func(r chi.Router) {
		r.Post("/", http.HandlerFunc(handler.PostHandler))
		r.Post("/api/shorten", http.HandlerFunc(handler.ApiShortenHandler))
		r.Get("/{id}", http.HandlerFunc(handler.GetHandler))
	})

	return router
}
