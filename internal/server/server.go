package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/generator"
	"github.com/kupriyanovkk/shortener/internal/models"
	"github.com/kupriyanovkk/shortener/internal/storage"
)

func PostRootHandler(w http.ResponseWriter, r *http.Request, env *config.Env) {
	body, err := io.ReadAll(r.Body)
	baseURL := env.Flags.B

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
	short, saveErr := env.Storage.AddValue(r.Context(), storage.AddValueOptions{
		Original: parsedURL.String(),
		BaseURL:  baseURL,
		Short:    id,
	})

	if errors.Is(saveErr, storage.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", short)
	w.Write([]byte(short))
}

func GetHandler(w http.ResponseWriter, r *http.Request, env *config.Env) {
	id := r.URL.String()
	origURL, err := env.Storage.GetValue(r.Context(), id[1:])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, origURL, http.StatusTemporaryRedirect)
}

func PostAPIHandler(w http.ResponseWriter, r *http.Request, env *config.Env) {
	var req models.Request
	dec := json.NewDecoder(r.Body)
	baseURL := env.Flags.B

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
	short, saveErr := env.Storage.AddValue(r.Context(), storage.AddValueOptions{
		Original: parsedURL.String(),
		BaseURL:  baseURL,
		Short:    id,
	})

	resp := models.Response{
		Result: short,
	}
	w.Header().Set("Content-Type", "application/json")
	if errors.Is(saveErr, storage.ErrConflict) {
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

func GetPingHandler(w http.ResponseWriter, r *http.Request, env *config.Env) {
	err := env.Storage.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func BatchHandler(w http.ResponseWriter, r *http.Request, env *config.Env) {
	var req []models.BatchRequest
	var result []models.BatchResponse
	baseURL := env.Flags.B
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req) == 0 {
		http.Error(w, "empty response", http.StatusBadRequest)
		return
	}

	for _, v := range req {
		parsedURL, err := url.ParseRequestURI(v.OriginalURL)
		if err != nil {
			http.Error(w, "Error parsing URL", http.StatusBadRequest)
			return
		}

		id, _ := generator.GetRandomStr(10)
		short, saveErr := env.Storage.AddValue(r.Context(), storage.AddValueOptions{
			Original: parsedURL.String(),
			BaseURL:  baseURL,
			Short:    id,
		})
		result = append(result, models.BatchResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      short,
		})
		if errors.Is(saveErr, storage.ErrConflict) {
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(result); err != nil {
		return
	}
}
