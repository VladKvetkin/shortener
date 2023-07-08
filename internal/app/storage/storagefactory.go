package storage

import (
	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type StorageFactory struct{}

func (sf *StorageFactory) GetStorage(config config.Config) (Storage, error) {
	if config.DatabaseDSN != "" {
		db, err := sqlx.Connect("postgres", config.DatabaseDSN)
		if err != nil {
			return nil, err
		}

		storage, err := NewPostgresStorage(db)
		if err != nil {
			return nil, err
		}

		zap.L().Info("Create database storage", zap.String("DatabaseDSN", config.DatabaseDSN))

		return storage, nil
	}

	storage := NewMemStorage(NewPersister(config.FileStoragePath))

	zap.L().Info("Create memory storage")

	return storage, nil
}
