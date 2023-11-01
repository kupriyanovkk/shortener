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

	var Store store.Store
	if f.D != "" {
		Store = db.NewStore(f.D)
	} else if f.F != "" {
		Store = infile.NewStore(f.F)
	} else {
		Store = inmemory.NewStore()
	}

	app := &config.App{
		Flags: f,
		Store: Store,
		URLChan: make(chan store.DeletedURLs, 10),
	}

	go handlers.FlushDeletedURLs(app)

	r.Use(
		middlewares.Logger,
		middlewares.Gzip,
		middlewares.Auth,
	)
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetID(w, r, app)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostRoot(w, r, app)
	})
	r.Route("/api", func(r chi.Router) {
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				handlers.PostAPIShorten(w, r, app)
			})

			r.Post("/batch", func(w http.ResponseWriter, r *http.Request) {
				handlers.PostAPIShortenBatch(w, r, app)
			})
		})

		r.Route("/user", func(r chi.Router) {
			r.Route("/urls", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					handlers.GetAPIUserURLs(w, r, app)
				})
				r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
					handlers.DeleteAPIUserURLs(w, r, app)
				})
			})
		})
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetPing(w, r, app)
	})

	err := http.ListenAndServe(f.A, r)
	if err != nil {
		panic(err)
	}
}
