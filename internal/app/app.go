package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/handlers"
	"github.com/kupriyanovkk/shortener/internal/middlewares"
	"github.com/kupriyanovkk/shortener/internal/storage"
)

func Start() {
	r := chi.NewRouter()
	f := config.ParseFlags()
	s := storage.NewStorage(f.F, f.D)
	env := &config.Env{Flags: f, Storage: s}

	r.Use(
		middlewares.Logger,
		middlewares.Gzip,
	)
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetID(w, r, env)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostRoot(w, r, env)
	})
	r.Route("/api", func(r chi.Router) {
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				handlers.PostAPIShorten(w, r, env)
			})

			r.Post("/batch", func(w http.ResponseWriter, r *http.Request) {
				handlers.PostAPIShortenBatch(w, r, env)
			})
		})
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetPing(w, r, env)
	})

	err := http.ListenAndServe(f.A, r)
	if err != nil {
		panic(err)
	}
}
