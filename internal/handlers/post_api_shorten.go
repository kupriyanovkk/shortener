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
	storeInterface "github.com/kupriyanovkk/shortener/internal/store/interface"
)

// PostAPIShorten process requests for shorten URL.
func PostAPIShorten(w http.ResponseWriter, r *http.Request, app *config.App) {
	var req models.Request
	dec := json.NewDecoder(r.Body)
	baseURL := app.Flags.BaseURL
	userID := fmt.Sprint(r.Context().Value(contextkey.ContextUserKey))

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil {
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	id, _ := generator.GetRandomStr(10)
	short, saveErr := app.Store.AddValue(r.Context(), storeInterface.AddValueOptions{
		Original: parsedURL.String(),
		BaseURL:  baseURL,
		Short:    id,
		UserID:   userID,
	})

	resp := models.Response{
		Result: short,
	}
	w.Header().Set("Content-Type", "application/json")
	if errors.Is(saveErr, failure.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		return
	}

	w.Header().Set("Location", short)
}
