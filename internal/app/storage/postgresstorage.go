package storage

import (
	"context"
	"database/sql"

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
	var originalURL sql.NullString

	row := s.db.QueryRowxContext(context.Background(), "SELECT original_url FROM url WHERE short_url = $1;", id)
	err := row.Scan(&originalURL)
	if err != nil {
		return "", err
	}

	if originalURL.Valid {
		return originalURL.String, nil
	}

	return "", ErrIDNotExists
}

func (s *PostgresStorage) Add(id string, url string, saveToPersister bool) error {
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
			original_url TEXT NOT NULL
		);
		`,
	)

	if err != nil {
		return err
	}

	return nil
}
