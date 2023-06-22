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

	_, err = h.storage.ReadByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrIDNotExists) {
			h.storage.Add(id, stringBody)
		} else {
			http.Error(res, "Invalid request", http.StatusBadRequest)
			return
		}
	}

	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("%s/%s", h.config.BaseShortURLAddress, id)))
}

func (h *Handler) ApiShortenHandler(res http.ResponseWriter, req *http.Request) {
	var requestModel models.ApiShortenRequest

	jsonDecoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	if err := jsonDecoder.Decode(&requestModel); err != nil {
		http.Error(res, "Cannot decode request JSON body", http.StatusInternalServerError)
		return
	}

	if requestModel.Url == "" {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := shortener.CreateID(requestModel.Url)
	if err != nil {
		http.Error(res, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err = h.storage.ReadByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrIDNotExists) {
			h.storage.Add(id, requestModel.Url)
		} else {
			http.Error(res, "Invalid request", http.StatusBadRequest)
			return
		}
	}

	responseModel := models.ApiShortenResponse{
		Result: fmt.Sprintf("%s/%s", h.config.BaseShortURLAddress, id),
	}

	jsonEncoder := json.NewEncoder(res)
	if err := jsonEncoder.Encode(responseModel); err != nil {
		http.Error(res, "Cannot encode response JSON body", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(http.StatusOK)
}
