package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

// Start initializes the application, sets up the router, parses flags, sets up the store,
// creates an application instance, and starts the server.
func Start() {
	router := chi.NewRouter()
	flags, err := config.ParseFlags(os.Args[0], os.Args[1:])
	if err != nil {
		panic(err)
	}

	store := getStore(flags)

	app := &config.App{
		Flags:   flags,
		Store:   store,
		URLChan: make(chan storeInterface.DeletedURLs, 10),
	}

	setupMiddlewares(router)
	setupRoutes(router, app)

	runServer(flags, router, app)
}

// getStore returns a store based on the provided flags.
func getStore(flags *config.ConfigFlags) storeInterface.Store {
	if flags.DatabaseDSN != "" {
		return db.NewStore(flags.DatabaseDSN)
	} else if flags.FileStoragePath != "" {
		return infile.NewStore(flags.FileStoragePath)
	}
	return inmemory.NewStore()
}

// setupMiddlewares sets up middleware for the router.
func setupMiddlewares(router *chi.Mux) {
	router.Use(
		middlewares.Logger,
		middlewares.Gzip,
		middlewares.Auth,
	)
	router.Mount("/debug", middleware.Profiler())
}

// setupRoutes sets up routes for the router.
func setupRoutes(router *chi.Mux, app *config.App) {
	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetID(w, r, app)
	})
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostRoot(w, r, app)
	})
	setupAPIRoutes(router, app)
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetPing(w, r, app)
	})
}

// setupAPIRoutes sets up API routes for the router.
func setupAPIRoutes(router *chi.Mux, app *config.App) {
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

		r.Route("/internal/stats", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				handlers.GetInternalStats(w, r, app)
			})
		})
	})
}

func runServer(flags *config.ConfigFlags, router http.Handler, app *config.App) {
	shutdownTimeout := 5 * time.Second

	server := &http.Server{
		Addr:    flags.ServerAddress,
		Handler: router,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	go func() {
		defer wg.Done()

		handlers.FlushDeletedURLs(app, ctx)
	}()

	go func() {
		defer wg.Done()

		if flags.EnableHTTPS {
			manager := &autocert.Manager{
				Cache:      autocert.DirCache("assets"),
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist("localhost"),
			}
			server.TLSConfig = manager.TLSConfig()

			if err := server.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
				log.Fatalf("HTTP server ListenAndServeTLS: %v", err)
			}
		} else {
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalf("HTTP server ListenAndServe: %v", err)
			}
		}
	}()

	<-ctx.Done()

	cancel()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	wg.Wait()
}
