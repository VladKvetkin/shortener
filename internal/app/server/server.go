package server

import (
	"net/http"

	"github.com/VladKvetkin/shortener/internal/app/config"
)

type Server struct {
	config config.Config
	router http.Handler
}

func NewServer(config config.Config, router http.Handler) *Server {
	return &Server{
		config: config,
		router: router,
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.config.GetAddress(), s.router)
}
