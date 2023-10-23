package handlers

import (
	"net/http"

	"github.com/kupriyanovkk/shortener/internal/config"
)

func GetPing(w http.ResponseWriter, r *http.Request, env *config.Env) {
	err := env.Store.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}