package server

import (
	"net/http"

	"github.com/VladKvetkin/shortener/internal/app/config"
)

type Server struct {
	config  config.Config
	handler http.Handler
}

func NewServer(config config.Config, handler http.Handler) *Server {
	return &Server{
		config:  config,
		handler: handler,
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.config.GetAddress(), s.handler)
}
