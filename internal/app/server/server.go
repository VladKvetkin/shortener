package server

import (
	"net/http"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/router"
)

type Server struct {
	config config.Config
	mux    *http.ServeMux
	router router.Router
}

func NewServer(config config.Config, router router.Router) *Server {
	return &Server{
		mux:    http.NewServeMux(),
		config: config,
		router: router,
	}
}

func (s *Server) Start() error {
	for pattern, handler := range s.router.Routes {
		s.mux.Handle(pattern, handler)
	}

	return http.ListenAndServe(s.config.GetAddress(), s.mux)
}
