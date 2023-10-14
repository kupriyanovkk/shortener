package storage

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	storageFile := "/tmp/short-url-db.json"

	testCases := []struct {
		description    string
		initialStorage map[string]string
		keyToAdd       string
		valueToAdd     string
		keyToGet       string
		expectedValue  string
		expectedError  error
	}{
		{
			description:    "Get existing value",
			initialStorage: map[string]string{"key1": "value1", "key2": "value2"},
			keyToGet:       "key2",
			expectedValue:  "value2",
			expectedError:  nil,
		},
		{
			description:    "Get non-existent value",
			initialStorage: map[string]string{"key1": "value1"},
			keyToGet:       "key2",
			expectedValue:  "",
			expectedError:  fmt.Errorf("value doesn't exist by key key2"),
		},
		{
			description:    "Add new value",
			initialStorage: map[string]string{"key1": "value1"},
			keyToAdd:       "key2",
			valueToAdd:     "value2",
			keyToGet:       "key2",
			expectedValue:  "value2",
			expectedError:  nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			storage := NewStorage(storageFile, "")
			storage.values = testCase.initialStorage

			if testCase.keyToAdd != "" {
				storage.AddValue(context.Background(), testCase.keyToAdd, testCase.valueToAdd)
			}

			value, err := storage.GetValue(context.Background(), testCase.keyToGet)

			assert.Equal(t, testCase.expectedValue, value, "Value mismatch")
			assert.Equal(t, testCase.expectedError, err, "Error mismatch")
		})
	}
}

func TestReadValues(t *testing.T) {
	tests := []struct {
		input     string
		expected  map[string]string
		expectErr bool
	}{
		{
			input: `
				{"uuid": 1, "short_url": "abc", "original_url": "https://example.com/1"}
				{"uuid": 2, "short_url": "def", "original_url": "https://example.com/2"}
				{"uuid": 3, "short_url": "ghi", "original_url": "https://example.com/3"}`,
			expected: map[string]string{
				"abc": "https://example.com/1",
				"def": "https://example.com/2",
				"ghi": "https://example.com/3",
			},
			expectErr: false,
		},
		{
			input:     "",
			expected:  map[string]string{},
			expectErr: false,
		},
		{
			input: `
				{"uuid": 1, "short_url": "abc", "original_url": "https://example.com/1"}
				{"uuid": 2, "short_url": "def", "original_url": "https://example.com/2"}
				{"uuid": 3, "short_url": "ghi", "original_url": "https://example.com/3"
			`,
			expected:  nil,
			expectErr: true,
		},
	}

	for _, test := range tests {
		scanner := bufio.NewScanner(strings.NewReader(test.input))
		result, err := ReadValues(scanner)

		if test.expectErr {
			if err == nil {
				t.Errorf("ReadValues did not return an expected error for input:\n%s", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("ReadValues returned an unexpected error for input:\n%s\nError: %v", test.input, err)
			}

			if !compareStringMaps(result, test.expected) {
				t.Errorf("ReadValues returned incorrect result for input:\n%s\nExpected: %v\nGot: %v", test.input, test.expected, result)
			}
		}
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
