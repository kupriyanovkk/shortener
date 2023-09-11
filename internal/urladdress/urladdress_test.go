package urladdress

import (
	"testing"
)

func TestUrlAddress(t *testing.T) {
	tests := []struct {
		description    string
		input          string
		expectedString string
		expectedHost   string
		expectedPort   string
		expectedError  bool
	}{
		{
			description:    "Valid input",
			input:          "http://example.com:8080",
			expectedString: "example.com:8080",
			expectedHost:   "example.com",
			expectedPort:   "8080",
			expectedError:  false,
		},
		{
			description:   "Invalid input",
			input:         "not-a-valid-url",
			expectedError: true,
		},
		{
			description:   "Empty input",
			input:         "",
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			urlAddr := New()

			err := urlAddr.Set(test.input)

			if test.expectedError {
				if err == nil {
					t.Error("Expected Set() to return an error, but it didn't")
				}
			} else {
				if err != nil {
					t.Errorf("Expected Set() to return no error, but got an error: %v", err)
				}

				if urlAddr.String() != test.expectedString {
					t.Errorf("Expected String() to return '%s', but got '%s'", test.expectedString, urlAddr.String())
				}

				if urlAddr.Host != test.expectedHost {
					t.Errorf("Expected Host to be '%s', but got '%s'", test.expectedHost, urlAddr.Host)
				}

				if urlAddr.Port != test.expectedPort {
					t.Errorf("Expected Port to be '%s', but got '%s'", test.expectedPort, urlAddr.Port)
				}
			}
		})
	}
}
