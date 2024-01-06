package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/kupriyanovkk/shortener/internal/failure"
	"github.com/kupriyanovkk/shortener/internal/models"
	storeInterface "github.com/kupriyanovkk/shortener/internal/store/interface"
)

// Store structure
type Store struct {
	values map[string]models.URL
}

// GetOriginalURL using for search original URL by short.
func (s Store) GetOriginalURL(ctx context.Context, short string) (string, error) {
	if value, ok := s.values[short]; ok {
		if value.DeletedFlag {
			return "", errors.New("URL is deleted")
		}

		return value.Original, nil
	}

	return "", fmt.Errorf("value doesn't exist by key %s", short)
}

// AddValue adding new URL into database.
func (s Store) AddValue(ctx context.Context, opts storeInterface.AddValueOptions) (string, error) {
	if opts.Original == "" {
		return "", failure.ErrEmptyOrigURL
	}

	s.values[opts.Short] = models.URL{
		Short:       opts.Short,
		Original:    opts.Original,
		UserID:      opts.UserID,
		DeletedFlag: false,
	}

	return fmt.Sprintf("%s/%s", opts.BaseURL, opts.Short), nil
}

// Ping checks database connection.
func (s Store) Ping() error {
	return nil
}

// GetUserURLs returning all URLs by particular user.
func (s Store) GetUserURLs(ctx context.Context, opts storeInterface.GetUserURLsOptions) ([]models.UserURL, error) {
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

// DeleteURLs marked URLs as deleted.
func (s Store) DeleteURLs(ctx context.Context, opts []storeInterface.DeletedURLs) error {
	for _, o := range opts {
		for _, value := range s.values {
			if value.UserID == o.UserID {
				for _, u := range o.URLs {
					if u == value.Short {
						s.values[value.Short] = models.URL{
							Short:       value.Short,
							Original:    value.Original,
							UserID:      value.UserID,
							DeletedFlag: true,
						}
					}
				}
			}
		}
	}

	return nil
}

// NewStore return Store for working with memory
func NewStore() storeInterface.Store {
	return Store{
		values: make(map[string]models.URL, 100),
	}
}
