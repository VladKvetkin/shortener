package grpchandlers

import (
	context "context"
	"errors"
	"fmt"
	"net"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/entities"
	"github.com/VladKvetkin/shortener/internal/app/shortener"
	"github.com/VladKvetkin/shortener/internal/app/storage"
	grpc "google.golang.org/grpc"
)

var (
	// ErrOriginalURLAlreadyExists - ошибка, которая означает, что оригинальный URL уже существует в базе данных.
	ErrOriginalURLAlreadyExists = errors.New("original URL already exists")
)

type GRPCServer struct {
	UnimplementedShortenerProtoServer
	storage storage.Storage
	config  config.Config
	server  *grpc.Server
}

func NewServer(storage storage.Storage, config config.Config) *GRPCServer {
	server := &GRPCServer{
		storage: storage,
		config:  config,
	}

	return server
}

func (s *GRPCServer) Start() error {
	listen, err := net.Listen("tcp", ":3333")
	if err != nil {
		return err
	}

	gRPCServer := grpc.NewServer()
	RegisterShortenerProtoServer(gRPCServer, s)
	s.server = gRPCServer

	return gRPCServer.Serve(listen)
}

func (s *GRPCServer) Stop() error {
	s.server.GracefulStop()
	return nil
}

func (s *GRPCServer) formatShortURL(id string) string {
	return fmt.Sprintf("%s/%s", s.config.BaseShortURLAddress, id)
}

func (s *GRPCServer) createAndAddID(ctx context.Context, URL string, userID string) (string, error) {
	id, err := shortener.CreateID(URL)
	if err != nil {
		return "", err
	}

	if _, err := s.storage.ReadByID(ctx, id); err != nil {
		if errors.Is(err, storage.ErrIDNotExists) {
			err := s.storage.Add(entities.URL{
				ShortURL:    id,
				OriginalURL: URL,
				UserID:      userID,
			})
			if err != nil {
				return "", err
			}

			return id, nil
		}

		return "", err
	}

	return id, ErrOriginalURLAlreadyExists
}
