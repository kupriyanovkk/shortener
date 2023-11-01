package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/kupriyanovkk/shortener/internal/models"
	"github.com/kupriyanovkk/shortener/internal/store"
)

type Store struct {
	values map[string]models.URL
}

func (s Store) GetOriginalURL(ctx context.Context, short string) (string, error) {
	if value, ok := s.values[short]; ok {
		if value.DeletedFlag {
			return "", errors.New("URL is deleted")
		}

		return value.Original, nil
	}

	return "", fmt.Errorf("value doesn't exist by key %s", short)
}

func (s Store) AddValue(ctx context.Context, opts store.AddValueOptions) (string, error) {
	if opts.Original == "" {
		return "", errors.New("original URL cannot be empty")
	}

	s.values[opts.Short] = models.URL{
		Short:       opts.Short,
		Original:    opts.Original,
		UserID:      opts.UserID,
		DeletedFlag: false,
	}

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

func (s Store) DeleteURLs(ctx context.Context, opts []store.DeletedURLs) error {
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

func NewStore() store.Store {
	return Store{
		values: make(map[string]models.URL),
	}
}
