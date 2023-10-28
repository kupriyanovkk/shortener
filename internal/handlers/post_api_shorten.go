package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/contextkey"
	"github.com/kupriyanovkk/shortener/internal/generator"
	"github.com/kupriyanovkk/shortener/internal/models"
	"github.com/kupriyanovkk/shortener/internal/store"
)

func PostAPIShorten(w http.ResponseWriter, r *http.Request, env *config.Env) {
	var req models.Request
	dec := json.NewDecoder(r.Body)
	baseURL := env.Flags.B
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
	short, saveErr := env.Store.AddValue(r.Context(), store.AddValueOptions{
		Original: parsedURL.String(),
		BaseURL:  baseURL,
		Short:    id,
		UserID:   userID,
	})

	resp := models.Response{
		Result: short,
	}
	w.Header().Set("Content-Type", "application/json")
	if errors.Is(saveErr, store.ErrConflict) {
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
