package storage

import (
	"context"
	"errors"

	"github.com/VladKvetkin/shortener/internal/app/entities"
	"go.uber.org/zap"
)

var (
	ErrIDNotExists = errors.New("id not exists")
)

type Storage interface {
	ReadByID(context.Context, string) (string, error)
	Add(entities.URL) error
	Ping() error
	AddBatch([]entities.URL) error
	Close() error
	GetUserURLs(context.Context, string) ([]entities.URL, error)
}

type MemStorage struct {
	storage   map[string]string
	persister Persister
}

func newMemStorage(persister Persister) Storage {
	storage := &MemStorage{
		storage:   make(map[string]string),
		persister: persister,
	}

	if err := persister.Restore(storage); err != nil {
		zap.L().Sugar().Errorw(
			"Cannot restore storage",
			"err", err,
		)
	}

	return storage
}

func (s *MemStorage) GetUserURLs(ctx context.Context, userID string) ([]entities.URL, error) {
	return nil, nil
}

func (s *MemStorage) ReadByID(ctx context.Context, id string) (string, error) {
	url, ok := s.storage[id]
	if !ok {
		return "", ErrIDNotExists
	}

	return url, nil
}

func (s *MemStorage) AddBatch(urls []entities.URL) error {
	for _, url := range urls {
		s.Add(url)
	}

	return nil
}

func (s *MemStorage) Add(url entities.URL) error {
	s.storage[url.ShortURL] = url.OriginalURL

	if err := s.persister.Save(url); err != nil {
		zap.L().Sugar().Errorw(
			"Cannot save data to persister",
			"err", err,
		)
	}

	return nil
}

func (s *MemStorage) Ping() error {
	return nil
}

func (s *MemStorage) Close() error {
	s.storage = nil

	return nil
}

func (s *MemStorage) AddWithoutPersisterSave(url entities.URL) error {
	s.storage[url.ShortURL] = url.OriginalURL

	return nil
}
