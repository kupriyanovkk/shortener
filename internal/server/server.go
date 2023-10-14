package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/generator"
	"github.com/kupriyanovkk/shortener/internal/models"
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
	env.Storage.AddValue(r.Context(), id, parsedURL.String())
	result := fmt.Sprintf("%s/%s", baseURL, id)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", result)
	w.Write([]byte(result))
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
	env.Storage.AddValue(r.Context(), id, parsedURL.String())
	result := fmt.Sprintf("%s/%s", baseURL, id)

	resp := models.Response{
		Result: result,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		return
	}

	w.Header().Set("Location", result)
}

func GetPingHandler(w http.ResponseWriter, r *http.Request, env *config.Env) {
	err := env.Storage.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
