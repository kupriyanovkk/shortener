package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/handlers"
	"github.com/kupriyanovkk/shortener/internal/middlewares"
	"github.com/kupriyanovkk/shortener/internal/store"
	"github.com/kupriyanovkk/shortener/internal/store/db"
	infile "github.com/kupriyanovkk/shortener/internal/store/in_file"
	inmemory "github.com/kupriyanovkk/shortener/internal/store/in_memory"
)

func Start() {
	r := chi.NewRouter()
	f := config.ParseFlags()

	var store store.Store
	if f.D != "" {
		store = db.NewStore(f.D)
	} else if f.F != "" {
		store = infile.NewStore(f.F)
	} else {
		store = inmemory.NewStore()
	}

	env := &config.Env{Flags: f, Store: store}

	r.Use(
		middlewares.Logger,
		middlewares.Gzip,
		middlewares.Auth,
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

		r.Route("/user", func(r chi.Router) {
			r.Route("/urls", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					handlers.GetAPIUserURLs(w, r, env)
				})
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
