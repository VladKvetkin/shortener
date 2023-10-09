// Package storage отвечает за работу с базой данных в приложении.

package storage

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/VladKvetkin/shortener/internal/app/entities"
)

var (
	// ErrIDNotExists - ошибка, которая означает, что сокращенная ссылка не найдена в базе данных.
	ErrIDNotExists = errors.New("id not exists")
)

// Storage - интерфейс базы данных приложения.
type Storage interface {
	// ReadByID - функция для получения entities.URL из базы данных.
	ReadByID(context.Context, string) (entities.URL, error)
	// Add - функция для добавления entities.URL в базу данных.
	Add(entities.URL) error
	// Ping - функция для проверки работоспособности базы данных.
	Ping() error
	// AddBatch - функция для добавления массива entities.URL в базу данных.
	AddBatch(context.Context, []entities.URL) error
	// DeleteBatch - функция для удаления сокращенных ссылок из базы данных.
	DeleteBatch(context.Context, []string, string) error
	// Close - функция для закрытия соединения с базой данных.
	Close() error
	// ReadByID - функция для получения массива entities.URL из базы данных.
	GetUserURLs(context.Context, string) ([]entities.URL, error)
}

// MemStorage - структура базы данных, которая хранит данные в мапе.
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

func (s *MemStorage) ReadByID(ctx context.Context, id string) (entities.URL, error) {
	url, ok := s.storage[id]
	if !ok {
		return entities.URL{}, ErrIDNotExists
	}

	return entities.URL{
		ShortURL:    id,
		OriginalURL: url,
	}, nil
}

func (s *MemStorage) AddBatch(ctx context.Context, urls []entities.URL) error {
	for _, url := range urls {
		s.Add(url)
	}

	return nil
}

func (s *MemStorage) DeleteBatch(ctx context.Context, shortURLs []string, userID string) error {
	for _, shortURL := range shortURLs {
		delete(s.storage, shortURL)
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
