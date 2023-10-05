package storage

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/VladKvetkin/shortener/internal/app/config"
)

// GetStorage - функция, которая возвращает Storage в зависимости от конфигурации приложения.
func GetStorage(config config.Config) (Storage, error) {
	if config.DatabaseDSN != "" {
		db, err := sqlx.Connect("postgres", config.DatabaseDSN)
		if err != nil {
			return nil, err
		}

		storage, err := newPostgresStorage(db)
		if err != nil {
			return nil, err
		}

		zap.L().Info("Create database storage", zap.String("DatabaseDSN", config.DatabaseDSN))

		return storage, nil
	}

	storage := newMemStorage(newPersister(config.FileStoragePath))

	zap.L().Info("Create memory storage")

	return storage, nil
}
