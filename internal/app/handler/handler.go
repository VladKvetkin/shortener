package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/entities"
	"github.com/VladKvetkin/shortener/internal/app/models"
	"github.com/VladKvetkin/shortener/internal/app/shortener"
	"github.com/VladKvetkin/shortener/internal/app/storage"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

var (
	ErrOriginalURLAlreadyExists = errors.New("original URL already exists")
)

type Handler struct {
	storage storage.Storage
	config  config.Config
}

func NewHandler(storage storage.Storage, config config.Config) *Handler {
	return &Handler{
		config:  config,
		storage: storage,
	}
}

func (h *Handler) PingHandler(res http.ResponseWriter, req *http.Request) {
	err := h.storage.Ping()
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(http.StatusText(http.StatusOK)))
}

func (h *Handler) GetHandler(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	fmt.Print(id)
	if id == "" {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	url, err := h.storage.ReadByID(id)
	if err != nil {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) PostHandler(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	stringBody := string(body)
	if stringBody == "" {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := h.createAndAddID(stringBody)
	if err != nil {
		if errors.Is(err, ErrOriginalURLAlreadyExists) {
			res.Header().Set("Content-type", "text/plain")
			res.WriteHeader(http.StatusConflict)
			res.Write([]byte(h.formatShortURL(id)))
			return
		}

		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(h.formatShortURL(id)))
}

func (h *Handler) APIShortenBatchHandler(res http.ResponseWriter, req *http.Request) {
	var requestModel []models.APIShortenBatchRequest

	jsonDecoder := json.NewDecoder(req.Body)

	if err := jsonDecoder.Decode(&requestModel); err != nil {
		http.Error(res, "Cannot decode request JSON body", http.StatusBadRequest)
		return
	}

	urls := make([]entities.URL, 0, len(requestModel))
	responseModel := make([]models.APIShortenBatchResponse, 0, len(requestModel))

	for _, batchData := range requestModel {
		shortURL, err := shortener.CreateID(batchData.OriginalURL)
		if err != nil {
			http.Error(res, "Invalid request", http.StatusBadRequest)
			return
		}

		urls = append(
			urls,
			entities.URL{
				UUID:        uuid.NewString(),
				OriginalURL: batchData.OriginalURL,
				ShortURL:    shortURL,
			},
		)

		responseModel = append(
			responseModel,
			models.APIShortenBatchResponse{
				CorrelationID: batchData.CorrelationID,
				ShortURL:      h.formatShortURL(shortURL),
			},
		)
	}

	err := h.storage.AddBatch(urls)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	jsonEncoder := json.NewEncoder(res)
	if err := jsonEncoder.Encode(responseModel); err != nil {
		http.Error(res, "Cannot encode response JSON body", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) APIShortenHandler(res http.ResponseWriter, req *http.Request) {
	var requestModel models.APIShortenRequest

	jsonDecoder := json.NewDecoder(req.Body)

	if err := jsonDecoder.Decode(&requestModel); err != nil {
		http.Error(res, "Cannot decode request JSON body", http.StatusBadRequest)
		return
	}

	if requestModel.URL == "" {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := h.createAndAddID(requestModel.URL)
	if err != nil {
		if errors.Is(err, ErrOriginalURLAlreadyExists) {
			h.sendJSONShortURL(res, id, http.StatusConflict)
			return
		}

		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	h.sendJSONShortURL(res, id, http.StatusCreated)
}

func (h *Handler) formatShortURL(id string) string {
	return fmt.Sprintf("%s/%s", h.config.BaseShortURLAddress, id)
}

func (h *Handler) createAndAddID(URL string) (string, error) {
	id, err := shortener.CreateID(URL)
	if err != nil {
		return "", err
	}

	if _, err := h.storage.ReadByID(id); err != nil {
		if errors.Is(err, storage.ErrIDNotExists) {
			err := h.storage.Add(id, URL)
			if err != nil {
				return "", err
			}

			return id, nil
		}

		return "", err
	}

	return id, ErrOriginalURLAlreadyExists
}

func (h *Handler) sendJSONShortURL(res http.ResponseWriter, id string, httpStatus int) {
	responseModel := models.APIShortenResponse{
		Result: h.formatShortURL(id),
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatus)

	jsonEncoder := json.NewEncoder(res)
	if err := jsonEncoder.Encode(responseModel); err != nil {
		http.Error(res, "Cannot encode response JSON body", http.StatusInternalServerError)
	}
}
