package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/kupriyanovkk/shortener/internal/store"
	"github.com/lib/pq"
)

type Store struct {
	values map[string]string
	db     store.DatabaseConnection
}

func (s Store) Bootstrap(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS shortener(
			uuid serial PRIMARY KEY,
			short varchar(128),
			original TEXT
		)
	`)

	tx.ExecContext(ctx, "CREATE UNIQUE INDEX url_id ON shortener (original)")

	return tx.Commit()
}

func (s Store) FindOriginalURL(ctx context.Context, short string) (originalURL string, err error) {
	row := s.db.QueryRowContext(ctx, `SELECT original FROM shortener WHERE short = $1`, short)
	err = row.Scan(&originalURL)
	return
}

func (s Store) FindShortURL(ctx context.Context, original string) (shortURL string, err error) {
	row := s.db.QueryRowContext(ctx, `SELECT short FROM shortener WHERE original = $1`, original)
	err = row.Scan(&shortURL)
	return
}

func (s Store) SaveURL(ctx context.Context, short, original string) error {
	_, err := s.db.ExecContext(ctx, `
			INSERT INTO shortener
			(short, original)
			VALUES
			($1, $2);
	`, short, original)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = store.ErrConflict
		}
	}

	return err
}

func (s Store) GetValue(ctx context.Context, short string) (string, error) {
	value, err := s.FindOriginalURL(ctx, short)
	return value, err
}

func (s Store) AddValue(ctx context.Context, opts store.AddValueOptions) (string, error) {
	if opts.Original == "" {
		return "", errors.New("original URL cannot be empty")
	}

	s.values[opts.Short] = opts.Original

	err := s.SaveURL(ctx, opts.Short, opts.Original)

	if err != nil && errors.Is(err, store.ErrConflict) {
		short, _ := s.FindShortURL(ctx, opts.Original)
		result := fmt.Sprintf("%s/%s", opts.BaseURL, short)

		return result, err
	}

	return fmt.Sprintf("%s/%s", opts.BaseURL, opts.Short), nil
}

func (s Store) Ping() error {
	err := s.db.Ping()
	return err
}

func NewStore(dbDSN string) store.Store {
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		panic(err)
	}

	store := Store{
		db:     db,
		values: make(map[string]string),
	}
	store.Bootstrap(context.Background())

	return store
}
