package shortener

import (
	"crypto/sha256"
	"encoding/base64"
)

func CreateID(url string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	hash := hasher.Sum(nil)

	id := base64.URLEncoding.EncodeToString(hash)

	return id[:8], nil
}
