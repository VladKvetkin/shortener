// Package router отвечает за маршрутизацию в приложении.

package router

import (
	"compress/gzip"
	"net/http"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/VladKvetkin/shortener/internal/app/handler"
	"github.com/VladKvetkin/shortener/internal/app/middleware"
)

// Router - структура маршрутизатора.
type Router struct {
	// Router - контроллер из пакета github.com/go-chi/chi.
	Router  *chi.Mux
	handler *handler.Handler
}

// NewRouter - конструктор Router.
// Задает middleware и маршруты в приложении.
func NewRouter(handler *handler.Handler) *Router {
	chiRouter := chi.NewRouter()

	router := &Router{
		Router:  chiRouter,
		handler: handler,
	}

	chiRouter.Use(
		middleware.DecompressBodyReader,
		middleware.Logger,
		chiMiddleware.Compress(gzip.BestSpeed, "application/json", "text/html"),
	)

	chiRouter.Route("/", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTCookie)

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

		r.Group(func(r chi.Router) {
			r.Use(middleware.IPChecker(handler.Config.TrustedSubnet))

			r.Get("/api/internal/stats", http.HandlerFunc(handler.GetInternalStats))
		})
	})

	return router
}
