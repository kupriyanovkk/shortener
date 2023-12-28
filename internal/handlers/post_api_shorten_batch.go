package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/contextkey"
	"github.com/kupriyanovkk/shortener/internal/failure"
	"github.com/kupriyanovkk/shortener/internal/generator"
	"github.com/kupriyanovkk/shortener/internal/models"
)

// PostAPIShortenBatch process requests for shorten URLs by batches.
func PostAPIShortenBatch(w http.ResponseWriter, r *http.Request, app *config.App) {
	var req []models.BatchRequest
	var result []models.BatchResponse
	baseURL := app.Flags.BaseURL
	dec := json.NewDecoder(r.Body)
	userID := fmt.Sprint(r.Context().Value(contextkey.ContextUserKey))

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req) == 0 {
		http.Error(w, "empty request", http.StatusBadRequest)
		return
	}

	for _, v := range req {
		parsedURL, err := url.ParseRequestURI(v.OriginalURL)
		if err != nil {
			http.Error(w, "Error parsing URL", http.StatusBadRequest)
			return
		}

		id, _ := generator.GetRandomStr(10)
		short, saveErr := app.Store.AddValue(r.Context(), models.AddValueOptions{
			Original: parsedURL.String(),
			BaseURL:  baseURL,
			Short:    id,
			UserID:   userID,
		})
		result = append(result, models.BatchResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      short,
		})
		if errors.Is(saveErr, failure.ErrConflict) {
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
