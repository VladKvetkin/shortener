package storage

import (
	"errors"

	"go.uber.org/zap"
)

var (
	ErrIDNotExists = errors.New("id not exists")
)

type Storage interface {
	ReadByID(id string) (string, error)
	Add(id string, url string, saveToRestorer bool) error
}

type MemStorage struct {
	storage  map[string]string
	restorer Restorer
}

func NewStorage(restorer Restorer) Storage {
	memStorage := &MemStorage{
		storage:  make(map[string]string),
		restorer: restorer,
	}

	if err := restorer.Restore(memStorage); err != nil {
		zap.L().Sugar().Errorw(
			"Cannot restore storage",
			"err", err,
		)
	}

	return memStorage
}

func (s *MemStorage) ReadByID(id string) (string, error) {
	url, ok := s.storage[id]
	if !ok {
		return "", ErrIDNotExists
	}

	return url, nil
}

func (s *MemStorage) Add(id string, url string, saveToRestorer bool) error {
	s.storage[id] = url

	if saveToRestorer {
		if err := s.restorer.Save(id, url); err != nil {
			zap.L().Sugar().Errorw(
				"Cannot save data to restorer",
				"err", err,
			)
		}
	}

	return nil
}
