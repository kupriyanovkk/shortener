package generator

import (
	"fmt"
	"testing"
)

func TestGetRandomStr(t *testing.T) {
	tests := []struct {
		size          int
		expectedError bool
	}{
		{0, false},
		{1, false},
		{10, false},
		{-1, true},
		{100, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Size %d", test.size), func(t *testing.T) {
			result, _ := GetRandomStr(test.size)

			if test.expectedError && result != "" {
				t.Errorf("Expected an error for size %d, but got: %s", test.size, result)
			}

			if !test.expectedError && len(result) != test.size {
				t.Errorf("Expected string of length %d for size %d, but got: %s", test.size, test.size, result)
			}
		})
	}
}
