package storage

import (
	"context"

	"github.com/VladKvetkin/shortener/internal/app/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func newPostgresStorage(db *sqlx.DB) (Storage, error) {
	storage := &PostgresStorage{
		db: db,
	}

	err := storage.createTables(context.TODO())
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *PostgresStorage) GetUserURLs(ctx context.Context, userID string) ([]entities.URL, error) {
	var userURLs []entities.URL

	err := s.db.SelectContext(ctx, &userURLs, "SELECT original_url, short_url FROM url WHERE user_id = $1;", userID)
	if err != nil {
		return nil, err
	}

	return userURLs, nil
}

func (s *PostgresStorage) ReadByID(ctx context.Context, id string) (entities.URL, error) {
	var url entities.URL

	row := s.db.QueryRowxContext(ctx, "SELECT id, short_url, original_url, user_id, is_deleted FROM url WHERE short_url = $1;", id)

	err := row.Scan(&url.UUID, &url.ShortURL, &url.OriginalURL, &url.UserID, &url.DeletedFlag)
	if err != nil {
		return entities.URL{}, ErrIDNotExists
	}

	return url, nil
}

func (s *PostgresStorage) AddBatch(ctx context.Context, urls []entities.URL) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	for _, url := range urls {
		_, err := tx.ExecContext(
			ctx,
			`
				INSERT INTO url (id, short_url, original_url, user_id)
				VALUES ($1, $2, $3, $4);
			`,
			uuid.NewString(), url.ShortURL, url.OriginalURL, url.UserID,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *PostgresStorage) DeleteBatch(ctx context.Context, shortURLs []string, userID string) error {
	_, err := s.db.ExecContext(
		ctx,
		`
			UPDATE url SET is_deleted = TRUE WHERE user_id = $1 AND short_url = ANY($2)
		`,
		userID, pq.Array(shortURLs),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) Add(url entities.URL) error {
	_, err := s.db.ExecContext(
		context.Background(),
		`
			INSERT INTO url (id, short_url, original_url, user_id)
			VALUES ($1, $2, $3, $4);
		`,
		uuid.NewString(), url.ShortURL, url.OriginalURL, url.UserID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) Ping() error {
	return s.db.Ping()
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

func (s *PostgresStorage) createTables(ctx context.Context) error {
	if err := s.createTableURL(ctx); err != nil {
		return err
	}

	return nil
}

func (s PostgresStorage) createTableURL(ctx context.Context) error {
	_, err := s.db.ExecContext(
		ctx,
		`
		CREATE TABLE IF NOT EXISTS url (
			id VARCHAR(36) PRIMARY KEY,
			short_url VARCHAR(255) NOT NULL,
			original_url TEXT NOT NULL UNIQUE,
			user_id VARCHAR(36) NOT NULL,
			is_deleted BOOLEAN DEFAULT FALSE
		);
		`,
	)

	if err != nil {
		return err
	}

	return nil
}
