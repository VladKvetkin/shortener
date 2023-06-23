package router

import (
	"net/http"

	"github.com/VladKvetkin/shortener/internal/app/handler"
	"github.com/VladKvetkin/shortener/internal/app/middleware"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
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

	chiRouter.Use(chiMiddleware.Compress(5, "application/json", "text/html"), middleware.Logger)
	chiRouter.Route("/", func(r chi.Router) {
		r.Post("/", http.HandlerFunc(handler.PostHandler))
		r.Post("/api/shorten", http.HandlerFunc(handler.APIShortenHandler))
		r.Get("/{id}", http.HandlerFunc(handler.GetHandler))
	})

	return router
}
