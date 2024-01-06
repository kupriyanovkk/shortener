package store

import (
	"context"
	"database/sql"

	"github.com/kupriyanovkk/shortener/internal/models"
)

// Store interface for storage working
type Store interface {
	GetOriginalURL(ctx context.Context, short string) (string, error)
	AddValue(ctx context.Context, opts AddValueOptions) (string, error)
	GetUserURLs(ctx context.Context, opts GetUserURLsOptions) ([]models.UserURL, error)
	Ping() error
	DeleteURLs(ctx context.Context, opts []DeletedURLs) error
}

// AddValueOptions is a structure for AddValue method params
type AddValueOptions struct {
	Original string
	BaseURL  string
	Short    string
	UserID   string
}

// GetUserURLsOptions is a structure for getting user URLs
type GetUserURLsOptions struct {
	UserID  string
	BaseURL string
}

// DeletedURLs is a structure for deleting URLs
type DeletedURLs struct {
	UserID string
	URLs   []string
}

// Database interface
type DatabaseConnection interface {
	Ping() error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
