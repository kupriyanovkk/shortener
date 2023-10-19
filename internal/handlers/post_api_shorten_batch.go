package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/generator"
	"github.com/kupriyanovkk/shortener/internal/models"
	"github.com/kupriyanovkk/shortener/internal/storage"
)

func PostAPIShortenBatch(w http.ResponseWriter, r *http.Request, env *config.Env) {
	var req []models.BatchRequest
	var result []models.BatchResponse
	baseURL := env.Flags.B
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req) == 0 {
		http.Error(w, "empty response", http.StatusBadRequest)
		return
	}

	for _, v := range req {
		parsedURL, err := url.ParseRequestURI(v.OriginalURL)
		if err != nil {
			http.Error(w, "Error parsing URL", http.StatusBadRequest)
			return
		}

		id, _ := generator.GetRandomStr(10)
		short, saveErr := env.Storage.AddValue(r.Context(), storage.AddValueOptions{
			Original: parsedURL.String(),
			BaseURL:  baseURL,
			Short:    id,
		})
		result = append(result, models.BatchResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      short,
		})
		if errors.Is(saveErr, storage.ErrConflict) {
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(result); err != nil {
		return
	}
}
