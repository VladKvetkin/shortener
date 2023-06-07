package storage

import "errors"

var (
	ErrShortURLNotExists = errors.New("short url not exists")
)

type Repositories interface {
	ReadByURL(url string) (string, bool, error)
	ReadByShortURL(shortURL string) (string, error)
	Add(shortURL string, url string) error
}

type Storage struct {
	storage map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		storage: make(map[string]string),
	}
}

func (s *Storage) ReadByURL(url string) (string, bool, error) {
	for key, value := range s.storage {
		if value == url {
			return key, true, nil
		}
	}

	return "", false, nil
}

func (s *Storage) ReadByShortURL(shortURL string) (string, error) {
	url, ok := s.storage[shortURL]
	if !ok {
		return "", ErrShortURLNotExists
	}

	return url, nil
}

func (s *Storage) Add(shortURL string, url string) error {
	s.storage[shortURL] = url

	return nil
}
