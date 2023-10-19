// Package server отвечает за сервер приложения.

package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/services/mycert"
)

// Server - структура сервера.
type Server struct {
	config config.Config
	server *http.Server
}

// NewServer - конструктор Server.
func NewServer(config config.Config, router http.Handler) *Server {
	return &Server{
		config: config,
		server: &http.Server{
			Addr:              config.Address,
			Handler:           router,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

// Start - функция, которая запускает сервер.
func (s *Server) Start() error {
	if s.config.EnableHTTPS {
		certFile, keyFile, err := mycert.GetCertAndKey()
		if err != nil {
			return err
		}

		return s.server.ListenAndServeTLS(certFile, keyFile)
	}

	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}
