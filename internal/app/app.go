package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/middlewares"
	"github.com/kupriyanovkk/shortener/internal/server"
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
		server.GetHandler(w, r, env)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		server.PostRootHandler(w, r, env)
	})
	r.Route("/api", func(r chi.Router) {
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				server.PostAPIHandler(w, r, env)
			})

			r.Post("/batch", func(w http.ResponseWriter, r *http.Request) {
				server.BatchHandler(w, r, env)
			})
		})
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		server.GetPingHandler(w, r, env)
	})

	err := http.ListenAndServe(f.A, r)
	if err != nil {
		panic(err)
	}
}
