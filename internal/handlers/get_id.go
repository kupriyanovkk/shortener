package handlers

import (
	"net/http"

	"github.com/kupriyanovkk/shortener/internal/config"
)

// GetID process requests for getting original URL
func GetID(w http.ResponseWriter, r *http.Request, app *config.App) {
	id := r.URL.String()
	origURL, err := app.Store.GetOriginalURL(r.Context(), id[1:])

	if err != nil {
		if err.Error() == "URL is deleted" {
			http.Error(w, err.Error(), http.StatusGone)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	http.Redirect(w, r, origURL, http.StatusTemporaryRedirect)
}
