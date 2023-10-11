// Package shortener отвечает за генерацию сокращенной ссылки.

package shortener

import (
	"crypto/sha256"
	"encoding/base64"
)

// CreateID - функция, которая сокращает url.
func CreateID(url string) (string, error) {
	hasher := sha256.New()
	if _, err := hasher.Write([]byte(url)); err != nil {
		return "", nil
	}

	hash := hasher.Sum(nil)

	id := base64.URLEncoding.EncodeToString(hash)

	return id[:8], nil
}
