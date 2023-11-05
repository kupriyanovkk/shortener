package inmemory

import (
	"context"
	"errors"
	"testing"

	"github.com/kupriyanovkk/shortener/internal/store"
)

func TestAddValue(t *testing.T) {
	testCases := []struct {
		description string
		opts        store.AddValueOptions
		expectedURL string
		expectedErr error
	}{
		{
			description: "Add valid value",
			opts: store.AddValueOptions{
				Original: "https://example.com",
				Short:    "abc",
				BaseURL:  "https://short.ly",
			},
			expectedURL: "https://short.ly/abc",
			expectedErr: nil,
		},
		{
			description: "Add value with empty Original",
			opts: store.AddValueOptions{
				Original: "",
				Short:    "def",
				BaseURL:  "https://short.ly",
			},
			expectedURL: "",
			expectedErr: errors.New("original URL cannot be empty"),
		},
	}

	store := NewStore()

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			url, err := store.AddValue(context.Background(), testCase.opts)

			if err != nil && testCase.expectedErr == nil {
				t.Errorf("Expected no error, but got an error: %v", err)
			}

			if err == nil && testCase.expectedErr != nil {
				t.Errorf("Expected an error (%v), but got no error", testCase.expectedErr)
			}

			if url != testCase.expectedURL {
				t.Errorf("Expected URL: %s, but got URL: %s", testCase.expectedURL, url)
			}
		})
	}
}

func TestStore_GetValue(t *testing.T) {
	testCases := []struct {
		description string
		addOpts     store.AddValueOptions
		short       string
		expectedURL string
		expectedErr error
	}{
		{
			description: "Get existing value",
			addOpts: store.AddValueOptions{
				Original: "https://example.com",
				Short:    "abc",
				BaseURL:  "https://short.ly",
			},
			short:       "abc",
			expectedURL: "https://example.com",
			expectedErr: nil,
		},
		{
			description: "Get non-existing value",
			addOpts: store.AddValueOptions{
				Original: "https://example.com",
				Short:    "abc",
				BaseURL:  "https://short.ly",
			},
			short:       "def",
			expectedURL: "",
			expectedErr: errors.New("value doesn't exist by key def"),
		},
	}

	store := NewStore()

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			_, _ = store.AddValue(context.Background(), testCase.addOpts)
			url, err := store.GetOriginalURL(context.Background(), testCase.short)

			if err != nil && testCase.expectedErr == nil {
				t.Errorf("Expected no error, but got an error: %v", err)
			}

			if err == nil && testCase.expectedErr != nil {
				t.Errorf("Expected an error (%v), but got no error", testCase.expectedErr)
			}

			if url != testCase.expectedURL {
				t.Errorf("Expected URL: %s, but got URL: %s", testCase.expectedURL, url)
			}
		})
	}
}
