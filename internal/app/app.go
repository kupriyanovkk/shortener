package app

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/handlers"
	"github.com/kupriyanovkk/shortener/internal/middlewares"
	"github.com/kupriyanovkk/shortener/internal/store/db"
	infile "github.com/kupriyanovkk/shortener/internal/store/in_file"
	inmemory "github.com/kupriyanovkk/shortener/internal/store/in_memory"
	storeInterface "github.com/kupriyanovkk/shortener/internal/store/interface"
	"golang.org/x/crypto/acme/autocert"
)

// Start it function witch init all API handlers,
// and init Storage.
func Start() {
	router := chi.NewRouter()
	flags := config.ParseFlags()

	var Store storeInterface.Store
	if flags.DatabaseDSN != "" {
		Store = db.NewStore(flags.DatabaseDSN)
	} else if flags.FileStoragePath != "" {
		Store = infile.NewStore(flags.FileStoragePath)
	} else {
		Store = inmemory.NewStore()
	}

	app := &config.App{
		Flags:   flags,
		Store:   Store,
		URLChan: make(chan storeInterface.DeletedURLs, 10),
	}

	go handlers.FlushDeletedURLs(app)

	router.Use(
		middlewares.Logger,
		middlewares.Gzip,
		middlewares.Auth,
	)
	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetID(w, r, app)
	})
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostRoot(w, r, app)
	})
	router.Route("/api", func(r chi.Router) {
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
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetPing(w, r, app)
	})

	router.Mount("/debug", middleware.Profiler())

	if flags.EnableHTTPS {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("assets"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("localhost"),
		}
		server := &http.Server{
			Addr:      flags.ServerAddress,
			Handler:   router,
			TLSConfig: manager.TLSConfig(),
		}

		err := server.ListenAndServeTLS("", "")
		if err != nil {
			panic(err)
		}
	} else {
		err := http.ListenAndServe(flags.ServerAddress, router)
		if err != nil {
			panic(err)
		}
	}
}
