package storage

import "errors"

var (
	ErrIDNotExists = errors.New("id not exists")
)

type Storage interface {
	ReadByURL(url string) (string, bool, error)
	ReadByID(id string) (string, error)
	Add(id string, url string) error
}

type MemStorage struct {
	storage map[string]string
}

func NewStorage() Storage {
	return &MemStorage{
		storage: make(map[string]string),
	}
}

func (s *MemStorage) ReadByURL(url string) (string, bool, error) {
	for key, value := range s.storage {
		if value == url {
			return key, true, nil
		}
	}

	return "", false, nil
}

func (s *MemStorage) ReadByID(id string) (string, error) {
	url, ok := s.storage[id]
	if !ok {
		return "", ErrIDNotExists
	}

	return url, nil
}

func (s *MemStorage) Add(id string, url string) error {
	s.storage[id] = url

	return nil
}
