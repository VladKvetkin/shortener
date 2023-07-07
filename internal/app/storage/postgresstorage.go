package storage

import (
	"context"

	"github.com/VladKvetkin/shortener/internal/app/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) (Storage, error) {
	storage := &PostgresStorage{
		db: db,
	}

	err := storage.createTables()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *PostgresStorage) ReadByID(id string) (string, error) {
	var originalURL string

	row := s.db.QueryRowxContext(context.Background(), "SELECT original_url FROM url WHERE short_url = $1;", id)

	err := row.Scan(&originalURL)
	if err != nil {
		return "", ErrIDNotExists
	}

	return originalURL, nil
}

func (s *PostgresStorage) AddBatch(urls []entities.URL) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	for _, url := range urls {
		_, err := tx.ExecContext(
			context.Background(),
			`
				INSERT INTO url (id, short_url, original_url)
				VALUES ($1, $2, $3);
			`,
			url.UUID, url.ShortURL, url.OriginalURL,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *PostgresStorage) Add(id string, url string) error {
	_, err := s.db.ExecContext(
		context.Background(),
		`
			INSERT INTO url (id, short_url, original_url)
			VALUES ($1, $2, $3);
		`,
		uuid.NewString(), id, url,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) Ping() error {
	return s.db.Ping()
}

func (s *PostgresStorage) createTables() error {
	_, err := s.db.ExecContext(
		context.Background(),
		`
		CREATE TABLE IF NOT EXISTS url (
			id VARCHAR(36) PRIMARY KEY,
			short_url VARCHAR(255) NOT NULL,
			original_url TEXT NOT NULL UNIQUE
		);
		`,
	)

	if err != nil {
		return err
	}

	return nil
}
