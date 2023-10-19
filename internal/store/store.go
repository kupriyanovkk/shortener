package store

import (
	"context"
	"database/sql"
	"errors"
)

type Store interface {
	GetValue(ctx context.Context, short string) (string, error)
	AddValue(ctx context.Context, opts AddValueOptions) (string, error)
	Ping() error
}

var ErrConflict = errors.New("data conflict")

type AddValueOptions struct {
	Original string
	BaseURL  string
	Short    string
}

type DatabaseConnection interface {
	Ping() error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
