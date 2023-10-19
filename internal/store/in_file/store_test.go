package infile

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/kupriyanovkk/shortener/internal/store"
)

func TestAddValue(t *testing.T) {
	testCases := []struct {
		description string
		filename    string
		opts        store.AddValueOptions
		expectedURL string
		expectedErr error
	}{
		{
			description: "Add valid value",
			filename:    "testfile.txt",
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
			filename:    "testfile.txt",
			opts: store.AddValueOptions{
				Original: "",
				Short:    "def",
				BaseURL:  "https://short.ly",
			},
			expectedURL: "",
			expectedErr: errors.New("original URL cannot be empty"),
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
}

func TestGetValue(t *testing.T) {
	testCases := []struct {
		description string
		storeValues map[string]string
		short       string
		expectedURL string
		expectedErr error
	}{
		{
			description: "Get existing value",
			storeValues: map[string]string{
				"abc": "https://example.com",
			},
			short:       "abc",
			expectedURL: "https://example.com",
			expectedErr: nil,
		},
		{
			description: "Get non-existing value",
			storeValues: map[string]string{},
			short:       "def",
			expectedURL: "",
			expectedErr: fmt.Errorf("value doesn't exist by key %s", "def"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			store := Store{values: testCase.storeValues}
			url, err := store.GetValue(context.Background(), testCase.short)

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
	testCases := []struct {
		name     string
		input    string
		expected map[string]string
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
			{"uuid": 1, "short_url": "abc", "original_url": "https://example.com"}
			{"uuid": 2, "short_url": "def", "original_url": "https://example.org"}`,
			expected: map[string]string{
				"abc": "https://example.com",
				"def": "https://example.org",
			},
			err: nil,
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

func compareStringMaps(a, b map[string]string) bool {
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
