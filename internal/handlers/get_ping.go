package handlers

import (
	"net/http"

	"github.com/kupriyanovkk/shortener/internal/config"
)

// GetPing process request for ping DB.
func GetPing(w http.ResponseWriter, r *http.Request, app *config.App) {
	err := app.Store.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
