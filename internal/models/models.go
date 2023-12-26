package models

import (
	"context"
	"database/sql"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type URL struct {
	UUID        int    `json:"uuid"`
	Short       string `json:"short_url"`
	Original    string `json:"original_url"`
	UserID      string `json:"user_id"`
	DeletedFlag bool   `json:"is_deleted"`
}

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UserURL struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

type Store interface {
	GetOriginalURL(ctx context.Context, short string) (string, error)
	AddValue(ctx context.Context, opts AddValueOptions) (string, error)
	GetUserURLs(ctx context.Context, opts GetUserURLsOptions) ([]UserURL, error)
	Ping() error
	DeleteURLs(ctx context.Context, opts []DeletedURLs) error
}

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

type DeletedURLs struct {
	UserID string
	URLs   []string
}

type DatabaseConnection interface {
	Ping() error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
