package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/google/uuid"

	"github.com/VladKvetkin/shortener/internal/app/models"
)

type Persister interface {
	Restore(storage Storage) error
	Save(id string, url string) error
}

type FilePersister struct {
	filePath string
}

func NewPersister(filePath string) Persister {
	return &FilePersister{
		filePath: filePath,
	}
}

func (fr *FilePersister) Restore(storage Storage) error {
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

		storage.Add(record.ShortURL, record.OriginalURL, false)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (fr *FilePersister) Save(id string, url string) error {
	file, err := os.OpenFile(fr.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	jsonRecord, err := json.Marshal(
		models.FileStorageRecord{
			UUID:        uuid.NewString(),
			ShortURL:    id,
			OriginalURL: url,
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
