package main

import (
	"net/http"

	"github.com/kupriyanovkk/shortener/internal/server"
	"github.com/kupriyanovkk/shortener/internal/storage"
)

func main() {
	mux := http.NewServeMux()
	var s = storage.NewStorage()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { server.HandleFunc(w, r, s) })

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
