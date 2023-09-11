package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kupriyanovkk/shortener/internal/generator"
	"github.com/kupriyanovkk/shortener/internal/storage"
)

func PostHandler(w http.ResponseWriter, r *http.Request, s storage.Storage, baseURL string) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
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
	s.AddValue(id, parsedURL.String())
	result := fmt.Sprintf("%s/%s", baseURL, id)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", result)
	w.Write([]byte(result))
}

func GetHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	id := r.URL.String()
	origURL, err := s.GetValue(id[1:])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, origURL, http.StatusTemporaryRedirect)
}
