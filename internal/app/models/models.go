// Пакет models используется для описаний моделей приложения.

package models

// APIShortenRequest - структура, которая описывает тело запроса для обработчика APIShortenHandler.
type APIShortenRequest struct {
	URL string `json:"url"`
}

// APIShortenResponse - структура, которая описывает тело ответа обработчика APIShortenHandler.
type APIShortenResponse struct {
	Result string `json:"result"`
}

// FileStorageRecord - структура, которая описывает формат сохранения сокращенных ссылок пользователя в файл.
type FileStorageRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// APIShortenBatchRequest - структура, которая описывает тело запроса для обработчика APIShortenBatchHandler.
type APIShortenBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// APIShortenBatchResponse - структура, которая описывает тело ответа обработчика APIShortenBatchHandler.
type APIShortenBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// APIUserURLResponse - структура, которая описывает тело ответа обработчика APIUserURLHandler.
type APIUserURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// APIUserDeleteURLRequest - тип, который описывает тело запроса для обработчика APIUserDeleteURLHandler.
type APIUserDeleteURLRequest []string
