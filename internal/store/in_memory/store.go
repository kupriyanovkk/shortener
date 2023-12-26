package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/kupriyanovkk/shortener/internal/failure"
	"github.com/kupriyanovkk/shortener/internal/models"
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

func (s Store) AddValue(ctx context.Context, opts models.AddValueOptions) (string, error) {
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

func (s Store) Ping() error {
	return nil
}

func (s Store) GetUserURLs(ctx context.Context, opts models.GetUserURLsOptions) ([]models.UserURL, error) {
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

func (s Store) DeleteURLs(ctx context.Context, opts []models.DeletedURLs) error {
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

func NewStore() models.Store {
	return Store{
		values: make(map[string]models.URL, 100),
	}
}
