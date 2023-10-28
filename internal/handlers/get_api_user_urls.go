package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/contextkey"
	"github.com/kupriyanovkk/shortener/internal/store"
)

func GetAPIUserURLs(w http.ResponseWriter, r *http.Request, env *config.Env) {
	userID := fmt.Sprint(r.Context().Value(contextkey.ContextUserKey))
	_, err := r.Cookie("UserID")

	if err != nil {
		http.Error(w, errors.New("missing user id").Error(), http.StatusUnauthorized)
		return
	}

	URLs, err := env.Store.GetUserURLs(r.Context(), store.GetUserURLsOptions{
		UserID:  userID,
		BaseURL: env.Flags.B,
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
