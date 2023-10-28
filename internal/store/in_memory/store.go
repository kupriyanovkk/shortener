package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/kupriyanovkk/shortener/internal/models"
	"github.com/kupriyanovkk/shortener/internal/store"
)

type Store struct {
	values map[string]store.AddValueOptions
}

func (s Store) GetValue(ctx context.Context, short string) (string, error) {
	if value, ok := s.values[short]; ok {
		return value.Original, nil
	}

	return "", fmt.Errorf("value doesn't exist by key %s", short)
}

func (s Store) AddValue(ctx context.Context, opts store.AddValueOptions) (string, error) {
	if opts.Original == "" {
		return "", errors.New("original URL cannot be empty")
	}

	s.values[opts.Short] = opts

	return fmt.Sprintf("%s/%s", opts.BaseURL, opts.Short), nil
}

func (s Store) Ping() error {
	return nil
}

func (s Store) GetUserURLs(ctx context.Context, opts store.GetUserURLsOptions) ([]models.UserURL, error) {
	result := make([]models.UserURL, 0, 100)
	for _, value := range s.values {
		if value.UserID == opts.UserID {
			result = append(result, models.UserURL{
				Short:    fmt.Sprintf("%s/%s", opts.BaseURL, value.Short),
				Original: value.Original,
			})
		}
	}

	return result, nil
}

func NewStore() store.Store {
	return Store{
		values: make(map[string]store.AddValueOptions),
	}
}
