package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kupriyanovkk/shortener/internal/generator"
	"github.com/kupriyanovkk/shortener/internal/storage"
)

var store = storage.NewStorage()

func HandleFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postHandler(w, r)
	case http.MethodGet:
		getHandler(w, r)
	default:
		http.Error(w, "Only POST or GET requests are allowed!", http.StatusBadRequest)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	bodyString := string(body)
	fmt.Println("Request Body:", bodyString)

	parsedURL, err := url.Parse(bodyString)
	if err != nil {
		http.Error(w, "Error parsing URL", http.StatusBadRequest)
		return
	}

	id := generator.GetRandomStr(10)
	store.AddValue(id, parsedURL.String())
	result := fmt.Sprintf("http://localhost:8080/%s", id)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", result)
	w.Write([]byte(result))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.String()
	origURL, err := store.GetValue(id[1:])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, origURL, http.StatusTemporaryRedirect)
}
