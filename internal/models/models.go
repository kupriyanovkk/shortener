package models

import (
	"context"
	"database/sql"
)

// Request struct
type Request struct {
	URL string `json:"url"`
}

// Response struct
type Response struct {
	Result string `json:"result"`
}

// URL is a structure contains all URL data
type URL struct {
	UUID        int    `json:"uuid"`
	Short       string `json:"short_url"`
	Original    string `json:"original_url"`
	UserID      string `json:"user_id"`
	DeletedFlag bool   `json:"is_deleted"`
}

// BatchRequest is a structure for URL batching
type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResponse is a structure for URL batching
type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserURL is a structure for user
type UserURL struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

// Store interface for storage working
type Store interface {
	GetOriginalURL(ctx context.Context, short string) (string, error)
	AddValue(ctx context.Context, opts AddValueOptions) (string, error)
	GetUserURLs(ctx context.Context, opts GetUserURLsOptions) ([]UserURL, error)
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
