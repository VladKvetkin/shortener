// Package server отвечает за сервер приложения.

package server

import (
	"net/http"

	"github.com/VladKvetkin/shortener/internal/app/config"
)

// Server - структура сервера.
type Server struct {
	config config.Config
	router http.Handler
}

// NewServer - конструктор Server.
func NewServer(config config.Config, router http.Handler) *Server {
	return &Server{
		config: config,
		router: router,
	}
}

// Start - функция, которая запускает сервер.
func (s *Server) Start() error {
	return http.ListenAndServe(s.config.Address, s.router)
}
