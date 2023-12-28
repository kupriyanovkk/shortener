package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/contextkey"
	"github.com/kupriyanovkk/shortener/internal/models"
)

// GetAPIUserURLs processes requests for getting user URLs
func GetAPIUserURLs(w http.ResponseWriter, r *http.Request, app *config.App) {
	userID := fmt.Sprint(r.Context().Value(contextkey.ContextUserKey))
	_, err := r.Cookie("UserID")

	if err != nil {
		http.Error(w, errors.New("missing user id").Error(), http.StatusUnauthorized)
		return
	}

	URLs, err := app.Store.GetUserURLs(r.Context(), models.GetUserURLsOptions{
		UserID:  userID,
		BaseURL: app.Flags.BaseURL,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(URLs) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(URLs); err != nil {
		return
	}
}
