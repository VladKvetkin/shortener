package router

import (
	"net/http"

	"github.com/VladKvetkin/shortener/internal/app/handler"
	"github.com/VladKvetkin/shortener/internal/app/logger"
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

	chiRouter.Route("/", func(r chi.Router) {
		r.Post("/", logger.WithLogging(http.HandlerFunc(handler.PostHandler)))
		r.Get("/{id}", logger.WithLogging(http.HandlerFunc(handler.GetHandler)))
	})

	return router
}
