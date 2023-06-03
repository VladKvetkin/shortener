package storage

import "errors"

type Repositories interface {
	ReadByURL(url string) string
	ReadByShortURL(shortURL string) (string, error)
	Add(shortURL string, url string) bool
}

type Storage struct {
	storage map[string]string
}

var (
	ErrShortURLNotExists = errors.New("short url not exists")
)

func NewStorage() *Storage {
	return &Storage{
		storage: make(map[string]string),
	}
}

func (s *Storage) ReadByURL(url string) string {
	shortURL, ok := s.storage[url]
	if !ok {
		return ""
	}

	return shortURL
}

func (s *Storage) ReadByShortURL(shortURL string) (string, error) {
	for key, value := range s.storage {
		if value == shortURL {
			return key, nil
		}
	}

	return "", ErrShortURLNotExists
}

func (s *Storage) Add(shortURL string, url string) bool {
	s.storage[url] = shortURL

	return true
}
