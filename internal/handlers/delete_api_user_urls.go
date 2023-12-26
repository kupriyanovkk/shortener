package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/contextkey"
	"github.com/kupriyanovkk/shortener/internal/models"
)

func DeleteAPIUserURLs(w http.ResponseWriter, r *http.Request, app *config.App) {
	var URLs []string
	dec := json.NewDecoder(r.Body)
	userID := fmt.Sprint(r.Context().Value(contextkey.ContextUserKey))
	_, err := r.Cookie("UserID")

	if err != nil {
		http.Error(w, errors.New("missing user id").Error(), http.StatusUnauthorized)
		return
	}

	if err := dec.Decode(&URLs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(URLs) == 0 {
		http.Error(w, "empty request", http.StatusBadRequest)
		return
	}

	app.URLChan <- models.DeletedURLs{
		UserID: userID,
		URLs:   URLs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func FlushDeletedURLs(app *config.App) {
	ticker := time.NewTicker(10 * time.Second)

	var URLs []models.DeletedURLs

	for {
		select {
		case u := <-app.URLChan:
			URLs = append(URLs, u)
		case <-ticker.C:
			if len(URLs) == 0 {
				continue
			}

			err := app.Store.DeleteURLs(context.TODO(), URLs)
			if err != nil {
				fmt.Println("cannot save urls", err)
				continue
			}
			URLs = nil
		}
	}
}
