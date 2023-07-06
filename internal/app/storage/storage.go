package storage

import (
	"errors"

	"github.com/VladKvetkin/shortener/internal/app/entities"
	"go.uber.org/zap"
)

var (
	ErrIDNotExists = errors.New("id not exists")
)

type Storage interface {
	ReadByID(id string) (string, error)
	Add(id string, url string) error
	Ping() error
	AddBatch([]entities.URL) error
}

type MemStorage struct {
	storage   map[string]string
	persister Persister
}

func NewMemStorage(persister Persister) Storage {
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

func (s *MemStorage) ReadByID(id string) (string, error) {
	url, ok := s.storage[id]
	if !ok {
		return "", ErrIDNotExists
	}

	return url, nil
}

func (s *MemStorage) AddBatch(urls []entities.URL) error {
	for _, url := range urls {
		s.Add(url.ShortURL, url.OriginalURL)
	}

	return nil
}

func (s *MemStorage) Add(id string, url string) error {
	s.storage[id] = url

	if err := s.persister.Save(id, url); err != nil {
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

func (s *MemStorage) AddWithoutPersisterSave(id string, url string) error {
	s.storage[id] = url

	return nil
}
