package handlers

import (
	"net/http"

	"github.com/kupriyanovkk/shortener/internal/config"
)

func GetID(w http.ResponseWriter, r *http.Request, env *config.Env) {
	id := r.URL.String()
	origURL, err := env.Storage.GetValue(r.Context(), id[1:])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, origURL, http.StatusTemporaryRedirect)
}
