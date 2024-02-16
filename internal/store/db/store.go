package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/kupriyanovkk/shortener/internal/failure"
	"github.com/kupriyanovkk/shortener/internal/models"
	storeInterface "github.com/kupriyanovkk/shortener/internal/store/interface"
	"github.com/lib/pq"
)

// Store structure
type Store struct {
	db storeInterface.DatabaseConnection
}

// Bootstrap function create table shortener and
// set unique index for 'original' field.
func (s Store) Bootstrap(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	tx.ExecContext(ctx, `
		CREATE TABLE shortener(
			id serial PRIMARY KEY,
			short varchar(128),
			original TEXT,
			user_id varchar(128) NOT NULL,
			is_deleted BOOLEAN NOT NULL
		)
	`)

	tx.ExecContext(ctx, "CREATE UNIQUE INDEX url_id ON shortener (original)")

	return tx.Commit()
}

// FindOriginalURL using for search original URL by short.
func (s Store) FindOriginalURL(ctx context.Context, short string) (models.URL, error) {
	var (
		original  string
		isDeleted bool
	)
	row := s.db.QueryRowContext(ctx, `SELECT original, is_deleted FROM shortener WHERE short = $1`, short)
	err := row.Scan(&original, &isDeleted)

	return models.URL{
		Original:    original,
		DeletedFlag: isDeleted,
	}, err
}

// FindShortURL using for search short URL by original.
func (s Store) FindShortURL(ctx context.Context, original string) (shortURL string, err error) {
	row := s.db.QueryRowContext(ctx, `SELECT short FROM shortener WHERE original = $1`, original)
	err = row.Scan(&shortURL)
	return
}

// InsertURL inserts new URL into a table.
func (s Store) InsertURL(ctx context.Context, short, original, userID string) error {
	_, err := s.db.ExecContext(ctx, `
			INSERT INTO shortener
			(short, original, user_id, is_deleted)
			VALUES
			($1, $2, $3, $4);
	`, short, original, userID, false)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = failure.ErrConflict
		}
	}

	return err
}

// GetOriginalURL using for search original URL by short.
func (s Store) GetOriginalURL(ctx context.Context, short string) (string, error) {
	URL, err := s.FindOriginalURL(ctx, short)

	if URL.DeletedFlag {
		return "", errors.New("URL is deleted")
	}

	return URL.Original, err
}

// AddValue adding new URL into database.
func (s Store) AddValue(ctx context.Context, opts storeInterface.AddValueOptions) (string, error) {
	if opts.Original == "" {
		return "", failure.ErrEmptyOrigURL
	}

	err := s.InsertURL(ctx, opts.Short, opts.Original, opts.UserID)

	if err != nil && errors.Is(err, failure.ErrConflict) {
		short, _ := s.FindShortURL(ctx, opts.Original)
		result := fmt.Sprintf("%s/%s", opts.BaseURL, short)

		return result, err
	}

	return fmt.Sprintf("%s/%s", opts.BaseURL, opts.Short), nil
}

// GetUserURLs returning all URLs by particular user.
func (s Store) GetUserURLs(ctx context.Context, opts storeInterface.GetUserURLsOptions) ([]models.UserURL, error) {
	limit := 100
	result := make([]models.UserURL, 0, limit)

	rows, err := s.db.QueryContext(ctx, `SELECT original, short FROM shortener WHERE user_id = $1 LIMIT $2`, opts.UserID, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u models.UserURL
		err = rows.Scan(&u.Original, &u.Short)
		if err != nil {
			return nil, err
		}

		result = append(result, models.UserURL{
			Short:    fmt.Sprintf("%s/%s", opts.BaseURL, u.Short),
			Original: u.Original,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteURLs marked URLs as deleted.
func (s Store) DeleteURLs(ctx context.Context, opts []storeInterface.DeletedURLs) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, o := range opts {
		for _, u := range o.URLs {
			_, err := tx.ExecContext(ctx, `
			UPDATE shortener SET is_deleted = TRUE
				WHERE short = $1 AND user_id = $2
		`, u, o.UserID)

			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

// Ping checks database connection.
func (s Store) Ping() error {
	err := s.db.Ping()
	return err
}

// GetInternalStats returning internal statistics
func (s Store) GetInternalStats(ctx context.Context) (models.InternalStats, error) {
	var (
		urlCount  int
		userCount int
	)
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM shortener`).Scan(&urlCount)
	if err != nil {
		return models.InternalStats{}, err
	}
	err = s.db.QueryRowContext(ctx, `SELECT user_id, COUNT(user_id) FROM shortener GROUP BY user_id`).Scan(&userCount)
	if err != nil {
		return models.InternalStats{}, err
	}

	return models.InternalStats{
		URLs:  urlCount,
		Users: userCount,
	}, nil
}

// NewStore return Store for working with DB
func NewStore(dbDSN string) storeInterface.Store {
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		panic(err)
	}

	store := Store{
		db: db,
	}
	store.Bootstrap(context.Background())

	return store
}
