package main

import (
	"net/http"

	"github.com/kupriyanovkk/shortener/internal/server"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", server.HandleFunc)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
