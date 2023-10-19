package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/kupriyanovkk/shortener/internal/store"
)

type Store struct {
	values map[string]string
}

func (s Store) GetValue(ctx context.Context, short string) (string, error) {
	if value, ok := s.values[short]; ok {
		return value, nil
	}

	return "", fmt.Errorf("value doesn't exist by key %s", short)
}

func (s Store) AddValue(ctx context.Context, opts store.AddValueOptions) (string, error) {
	if opts.Original == "" {
		return "", errors.New("original URL cannot be empty")
	}

	s.values[opts.Short] = opts.Original

	return fmt.Sprintf("%s/%s", opts.BaseURL, opts.Short), nil
}

func (s Store) Ping() error {
	return nil
}

func NewStore() store.Store {
	return Store{
		values: make(map[string]string),
	}
}
