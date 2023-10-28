package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kupriyanovkk/shortener/internal/models"
)

type Store interface {
	GetValue(ctx context.Context, short string) (string, error)
	AddValue(ctx context.Context, opts AddValueOptions) (string, error)
	GetUserURLs(ctx context.Context, opts GetUserURLsOptions) ([]models.UserURL, error)
	Ping() error
}

var ErrConflict = errors.New("data conflict")

type AddValueOptions struct {
	Original string
	BaseURL  string
	Short    string
	UserID   string
}

type GetUserURLsOptions struct {
	UserID  string
	BaseURL string
}

type DatabaseConnection interface {
	Ping() error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
