package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/contextkey"
	"github.com/kupriyanovkk/shortener/internal/generator"
	"github.com/kupriyanovkk/shortener/internal/store"
)

func PostRoot(w http.ResponseWriter, r *http.Request, app *config.App) {
	body, err := io.ReadAll(r.Body)
	baseURL := app.Flags.B
	userID := fmt.Sprint(r.Context().Value(contextkey.ContextUserKey))

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	bodyString := string(body)
	parsedURL, err := url.ParseRequestURI(bodyString)

	if err != nil {
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	id, _ := generator.GetRandomStr(10)
	short, saveErr := app.Store.AddValue(r.Context(), store.AddValueOptions{
		Original: parsedURL.String(),
		BaseURL:  baseURL,
		Short:    id,
		UserID:   userID,
	})

	if errors.Is(saveErr, store.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", short)
	w.Write([]byte(short))
}
