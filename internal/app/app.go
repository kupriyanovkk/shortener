package app

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/middlewares"
	"github.com/kupriyanovkk/shortener/internal/server"
	"github.com/kupriyanovkk/shortener/internal/storage"
	_ "github.com/lib/pq"
)

func Start() {
	r := chi.NewRouter()
	f := config.ParseFlags()
	s := storage.NewStorage(f.F)

	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, dbErr := sql.Open("postgres", connStr)
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	r.Use(
		middlewares.Logger,
		middlewares.Gzip,
	)
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		server.GetHandler(w, r, s)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		server.PostRootHandler(w, r, s, f.B)
	})
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		server.PostAPIHandler(w, r, s, f.B)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		server.GetPingHandler(w, r, db)
	})

	err := http.ListenAndServe(f.A, r)
	if err != nil {
		panic(err)
	}
}
