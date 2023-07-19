package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/entities"
	"github.com/VladKvetkin/shortener/internal/app/middleware"
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

func (h *Handler) GetUserUrls(res http.ResponseWriter, req *http.Request) {
	userID, ok := req.Context().Value(middleware.UserIDKey{}).(string)
	if !ok {
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	userURLs, err := h.storage.GetUserURLs(req.Context(), userID)
	if err != nil {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	if len(userURLs) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	responseModel := make([]models.APIUserURLResponse, 0, len(userURLs))
	for _, userURL := range userURLs {
		responseModel = append(
			responseModel,
			models.APIUserURLResponse{
				ShortURL:    userURL.ShortURL,
				OriginalURL: userURL.OriginalURL,
			},
		)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	jsonEncoder := json.NewEncoder(res)
	if err := jsonEncoder.Encode(responseModel); err != nil {
		http.Error(res, "Cannot encode response JSON body", http.StatusInternalServerError)
		return
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
}

func (h *Handler) GetHandler(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	url, err := h.storage.ReadByID(req.Context(), id)
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

	userID, ok := req.Context().Value(middleware.UserIDKey{}).(string)
	if !ok {
		userID = uuid.NewString()
	}

	id, err := h.createAndAddID(req.Context(), stringBody, userID)
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

	userID, ok := req.Context().Value(middleware.UserIDKey{}).(string)
	if !ok {
		userID = uuid.NewString()
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
				OriginalURL: batchData.OriginalURL,
				ShortURL:    shortURL,
				UserID:      userID,
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

	userID, ok := req.Context().Value(middleware.UserIDKey{}).(string)
	if !ok {
		userID = uuid.NewString()
	}

	id, err := h.createAndAddID(req.Context(), requestModel.URL, userID)
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

func (h *Handler) createAndAddID(ctx context.Context, URL string, userID string) (string, error) {
	id, err := shortener.CreateID(URL)
	if err != nil {
		return "", err
	}

	if _, err := h.storage.ReadByID(ctx, id); err != nil {
		if errors.Is(err, storage.ErrIDNotExists) {
			err := h.storage.Add(entities.URL{
				ShortURL:    id,
				OriginalURL: URL,
				UserID:      userID,
			})
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
