package infile

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kupriyanovkk/shortener/internal/failure"
	"github.com/kupriyanovkk/shortener/internal/models"
	storeInterface "github.com/kupriyanovkk/shortener/internal/store/interface"
)

func TestAddValue(t *testing.T) {
	fileName := "testfile.txt"
	testCases := []struct {
		description string
		filename    string
		opts        storeInterface.AddValueOptions
		expectedURL string
		expectedErr error
	}{
		{
			description: "Add valid value",
			filename:    fileName,
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
			filename:    fileName,
			opts: storeInterface.AddValueOptions{
				Original: "",
				Short:    "def",
				BaseURL:  "https://short.ly",
			},
			expectedURL: "",
			expectedErr: failure.ErrEmptyOrigURL,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			store := NewStore(testCase.filename)
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

	os.Remove(fileName)
}

func TestGetValue(t *testing.T) {
	values := make(map[string]models.URL)
	values["abc"] = models.URL{
		Short:    "abc",
		Original: "https://example.com",
		UserID:   "123",
	}

	testCases := []struct {
		description string
		storeValues map[string]models.URL
		short       string
		expectedURL string
		expectedErr error
	}{
		{
			description: "Get existing value",
			storeValues: values,
			short:       "abc",
			expectedURL: "https://example.com",
			expectedErr: nil,
		},
		{
			description: "Get non-existing value",
			storeValues: make(map[string]models.URL),
			short:       "def",
			expectedURL: "",
			expectedErr: fmt.Errorf("value doesn't exist by key %s", "def"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			store := Store{values: testCase.storeValues}
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

func TestReadValuesFromFile(t *testing.T) {
	var errParsingJSON = errors.New("unexpected end of JSON input")
	values := make(map[string]models.URL)
	values["abc"] = models.URL{
		UUID:     1,
		Short:    "abc",
		Original: "https://example.com",
		UserID:   "123",
	}
	values["def"] = models.URL{
		UUID:     2,
		Short:    "def",
		Original: "https://example.org",
		UserID:   "123",
	}

	testCases := []struct {
		name     string
		input    string
		expected map[string]models.URL
		err      error
	}{
		{
			name:     "Read from empty file",
			input:    "",
			expected: nil,
			err:      nil,
		},
		{
			name: "Read from valid JSON file",
			input: `
			{"uuid": 1, "short_url": "abc", "original_url": "https://example.com", "user_id": "123"}
			{"uuid": 2, "short_url": "def", "original_url": "https://example.org", "user_id": "123"}`,
			expected: values,
			err:      nil,
		},
		{
			name: "Error while parsing JSON",
			input: `{"short": "abc", "original": "https://example.com"}
												{"short": "def", "original": "https://example.org"`,
			expected: nil,
			err:      errParsingJSON,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputReader := strings.NewReader(tc.input)
			scanner := bufio.NewScanner(inputReader)
			values, err := ReadValuesFromFile(scanner)

			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Expected error: %v, got: %v", tc.err, err)
			}

			if !compareStringMaps(values, tc.expected) {
				t.Errorf("Expected values: %v, got: %v", tc.expected, values)
			}
		})
	}
}

func compareStringMaps(a, b map[string]models.URL) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valA := range a {
		valB, exists := b[key]
		if !exists || valA != valB {
			return false
		}
	}

	return true
}

func TestDeleteURLs(t *testing.T) {
	s := Store{
		values: map[string]models.URL{
			"short1": {Short: "short1", Original: "original1", UserID: "user1", DeletedFlag: false},
			"short2": {Short: "short2", Original: "original2", UserID: "user2", DeletedFlag: false},
		},
	}

	deletedURLs := []storeInterface.DeletedURLs{
		{UserID: "user1", URLs: []string{"short1"}},
		{UserID: "user2", URLs: []string{"short2"}},
	}

	err := s.DeleteURLs(context.Background(), deletedURLs)
	if err != nil {
		t.Errorf("DeleteURLs returned an error: %v", err)
	}
}

func TestStore_GetInternalStats(t *testing.T) {
	fileName := "testfile.txt"
	ctx := context.Background()
	store := NewStore(fileName)

	expectedStatsEmpty := models.InternalStats{
		URLs:  0,
		Users: 0,
	}

	expectedStatsNonEmpty := models.InternalStats{
		URLs:  3,
		Users: 2,
	}

	t.Run("EmptyValues", func(t *testing.T) {
		stats, _ := store.GetInternalStats(context.Background())
		if stats != expectedStatsEmpty {
			t.Errorf("Expected %v, but got %v", expectedStatsEmpty, stats)
		}
	})

	t.Run("NonEmptyValues", func(t *testing.T) {
		store.AddValue(ctx, storeInterface.AddValueOptions{
			Original: "original1",
			Short:    "short1",
			UserID:   "user1",
		})
		store.AddValue(ctx, storeInterface.AddValueOptions{
			Original: "original2",
			Short:    "short2",
			UserID:   "user2",
		})
		store.AddValue(ctx, storeInterface.AddValueOptions{
			Original: "original3",
			Short:    "short3",
			UserID:   "user2",
		})

		stats, _ := store.GetInternalStats(context.Background())
		if stats != expectedStatsNonEmpty {
			t.Errorf("Expected %v, but got %v", expectedStatsNonEmpty, stats)
		}
	})

	os.Remove(fileName)
}
