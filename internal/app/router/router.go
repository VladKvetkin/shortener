package router

import (
	"github.com/VladKvetkin/shortener/internal/app/handler"
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
		r.Post("/", handler.PostHandler)
		r.Get("/{id}", handler.GetHandler)
	})

	return router
}
