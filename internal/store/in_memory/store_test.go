package inmemory

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kupriyanovkk/shortener/internal/failure"
	"github.com/kupriyanovkk/shortener/internal/models"
	storeInterface "github.com/kupriyanovkk/shortener/internal/store/interface"
)

func TestAddValue(t *testing.T) {
	testCases := []struct {
		description string
		opts        storeInterface.AddValueOptions
		expectedURL string
		expectedErr error
	}{
		{
			description: "Add valid value",
			opts: storeInterface.AddValueOptions{
				Original: "https://example.com",
				Short:    "abc",
				BaseURL:  "https://short.ly",
			},
			expectedURL: "https://short.ly/abc",
			expectedErr: nil,
		},
		{
			description: "Add value with empty Original",
			opts: storeInterface.AddValueOptions{
				Original: "",
				Short:    "def",
				BaseURL:  "https://short.ly",
			},
			expectedURL: "",
			expectedErr: failure.ErrEmptyOrigURL,
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
		addOpts     storeInterface.AddValueOptions
		short       string
		expectedURL string
		expectedErr error
	}{
		{
			description: "Get existing value",
			addOpts: storeInterface.AddValueOptions{
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
			addOpts: storeInterface.AddValueOptions{
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

func TestStore_DeleteURLs(t *testing.T) {
	ctx := context.Background()
	s := Store{
		values: map[string]models.URL{
			"short1": {Short: "short1", Original: "original1", UserID: "user1"},
			"short2": {Short: "short2", Original: "original2", UserID: "user2"},
		},
	}

	tests := []struct {
		name     string
		opts     []storeInterface.DeletedURLs
		expected map[string]models.URL
	}{
		{
			name: "Mark URL as deleted for matching UserID and Short URL",
			opts: []storeInterface.DeletedURLs{
				{UserID: "user1", URLs: []string{"short1"}},
			},
			expected: map[string]models.URL{
				"short1": {Short: "short1", Original: "original1", UserID: "user1", DeletedFlag: true},
				"short2": {Short: "short2", Original: "original2", UserID: "user2"},
			},
		},
		{
			name: "No URLs to mark as deleted",
			opts: []storeInterface.DeletedURLs{
				{UserID: "user3", URLs: []string{"short3"}},
			},
			expected: map[string]models.URL{
				"short1": {Short: "short1", Original: "original1", UserID: "user1", DeletedFlag: true},
				"short2": {Short: "short2", Original: "original2", UserID: "user2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.DeleteURLs(ctx, tt.opts)
			if !reflect.DeepEqual(s.values, tt.expected) {
				t.Errorf("unexpected values map after %s test; got %v, want %v", tt.name, s.values, tt.expected)
			}
		})
	}
}
