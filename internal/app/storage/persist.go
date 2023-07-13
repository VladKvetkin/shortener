package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/VladKvetkin/shortener/internal/app/entities"
	"github.com/VladKvetkin/shortener/internal/app/models"
	"github.com/google/uuid"
)

type Persister interface {
	Restore(storage *MemStorage) error
	Save(entities.URL) error
}

type FilePersister struct {
	filePath string
}

func newPersister(filePath string) Persister {
	return &FilePersister{
		filePath: filePath,
	}
}

func (fr *FilePersister) Restore(storage *MemStorage) error {
	file, err := os.OpenFile(fr.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var record models.FileStorageRecord

		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return err
		}

		storage.AddWithoutPersisterSave(entities.URL{
			UUID:        record.UUID,
			ShortURL:    record.ShortURL,
			OriginalURL: record.OriginalURL,
		})
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (fr *FilePersister) Save(url entities.URL) error {
	file, err := os.OpenFile(fr.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	jsonRecord, err := json.Marshal(
		models.FileStorageRecord{
			UUID:        uuid.NewString(),
			ShortURL:    url.ShortURL,
			OriginalURL: url.OriginalURL,
		},
	)

	if err != nil {
		return err
	}

	if _, err := file.Write(jsonRecord); err != nil {
		return err
	}

	if _, err := file.Write([]byte{'\n'}); err != nil {
		return err
	}

	return nil
}
