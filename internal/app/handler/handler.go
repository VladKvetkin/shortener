package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/models"
	"github.com/VladKvetkin/shortener/internal/app/shortener"
	"github.com/VladKvetkin/shortener/internal/app/storage"
	"github.com/go-chi/chi"
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

func (h *Handler) GetHandler(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
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

	defer req.Body.Close()

	if err != nil {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	stringBody := string(body)

	if stringBody == "" {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := shortener.CreateID(stringBody)
	if err != nil {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	if _, err = h.storage.ReadByID(id); err != nil {
		if errors.Is(err, storage.ErrIDNotExists) {
			h.storage.Add(id, stringBody)
		} else {
			http.Error(res, "Invalid request", http.StatusBadRequest)
			return
		}
	}

	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(h.formatShortURL(id)))
}

func (h *Handler) APIShortenHandler(res http.ResponseWriter, req *http.Request) {
	var requestModel models.APIShortenRequest

	jsonDecoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	if err := jsonDecoder.Decode(&requestModel); err != nil {
		http.Error(res, "Cannot decode request JSON body", http.StatusInternalServerError)
		return
	}

	if requestModel.URL == "" {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := shortener.CreateID(requestModel.URL)
	if err != nil {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	if _, err = h.storage.ReadByID(id); err != nil {
		if errors.Is(err, storage.ErrIDNotExists) {
			h.storage.Add(id, requestModel.URL)
		} else {
			http.Error(res, "Invalid request", http.StatusBadRequest)
			return
		}
	}

	responseModel := models.APIShortenResponse{
		Result: h.formatShortURL(id),
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	jsonEncoder := json.NewEncoder(res)
	if err := jsonEncoder.Encode(responseModel); err != nil {
		http.Error(res, "Cannot encode response JSON body", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) formatShortURL(id string) string {
	return fmt.Sprintf("%s/%s", h.config.BaseShortURLAddress, id)
}
