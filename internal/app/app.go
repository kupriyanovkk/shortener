package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/server"
	"github.com/kupriyanovkk/shortener/internal/storage"
)

func Start() {
	var s = storage.NewStorage()

	r := chi.NewRouter()
	f := config.ParseFlags()

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		server.GetHandler(w, r, s)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		server.PostHandler(w, r, s, f.B)
	})

	err := http.ListenAndServe(f.A, r)
	if err != nil {
		panic(err)
	}
}
