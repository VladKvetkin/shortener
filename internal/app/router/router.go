package router

import (
	"compress/gzip"
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

	chiRouter.Use(
		middleware.DecompressBodyReader,
		middleware.JWTCookie,
		middleware.Logger,
		chiMiddleware.Compress(gzip.BestCompression, "application/json", "text/html"),
	)

	chiRouter.Route("/", func(r chi.Router) {
		r.Post("/", http.HandlerFunc(handler.PostHandler))
		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", http.HandlerFunc(handler.APIShortenHandler))
				r.Post("/batch", http.HandlerFunc(handler.APIShortenBatchHandler))
			})

			r.Get("/user/urls", http.HandlerFunc(handler.GetUserUrlsHandler))
			r.Delete("/user/urls", http.HandlerFunc(handler.DeleteUserUrlsHandler))
		})
		r.Get("/{id}", http.HandlerFunc(handler.GetHandler))
		r.Get("/ping", http.HandlerFunc(handler.PingHandler))
	})

	return router
}
