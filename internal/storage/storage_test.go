package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
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
			storage := Storage{values: testCase.initialStorage}

			if testCase.keyToAdd != "" {
				storage.AddValue(testCase.keyToAdd, testCase.valueToAdd)
			}

			value, err := storage.GetValue(testCase.keyToGet)

			assert.Equal(t, testCase.expectedValue, value, "Value mismatch")
			assert.Equal(t, testCase.expectedError, err, "Error mismatch")
		})
	}
}
