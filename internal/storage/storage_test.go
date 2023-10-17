package storage

import (
	"bufio"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestStorage_GetValue(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	s := Storage{
		mode:   DataBaseStorage,
		db:     db,
		values: make(map[string]string),
	}

	testCases := []struct {
		short         string
		expectedValue string
		expectedErr   error
		dbMockExpect  func()
	}{
		{
			short:         "example",
			expectedValue: "http://example.com",
			expectedErr:   nil,
			dbMockExpect: func() {
				mock.ExpectQuery("SELECT original FROM shortener WHERE short = ?").
					WithArgs("example").
					WillReturnRows(sqlmock.NewRows([]string{"original"}).AddRow("http://example.com"))
			},
		},
		{
			short:         "nonexistent",
			expectedValue: "",
			expectedErr:   errors.New("not found"),
			dbMockExpect: func() {
				mock.ExpectQuery("SELECT original FROM shortener WHERE short = ?").
					WithArgs("nonexistent").
					WillReturnError(errors.New("not found"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.short, func(t *testing.T) {
			tc.dbMockExpect()
			value, err := s.GetValue(context.Background(), tc.short)
			assert.Equal(t, tc.expectedValue, value)
			assert.Equal(t, tc.expectedErr, err)
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_AddValue(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred while creating mock database: %s", err)
	}
	defer db.Close()

	storage := Storage{
		mode:   DataBaseStorage,
		db:     db,
		values: make(map[string]string),
	}

	testCases := []struct {
		name          string
		short         string
		original      string
		dbExpectation func()
		expectedURL   string
		expectedErr   error
	}{
		{
			name:     "AddValue success",
			short:    "example",
			original: "https://example.com",
			dbExpectation: func() {
				mock.ExpectExec("INSERT INTO shortener").WithArgs("example", "https://example.com").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedURL: "https://example.com/example",
			expectedErr: nil,
		},
		{
			name:     "AddValue conflict",
			short:    "example",
			original: "https://example.com",
			dbExpectation: func() {
				mock.ExpectExec("INSERT INTO shortener").WithArgs("example", "https://example.com").WillReturnError(ErrConflict)
				mock.ExpectQuery("SELECT short FROM shortener").WithArgs("https://example.com").WillReturnRows(sqlmock.NewRows([]string{"short"}).AddRow("example"))
			},
			expectedURL: "https://example.com/example",
			expectedErr: ErrConflict,
		},
		{
			name:          "AddValue with FileWriterStorage",
			short:         "example",
			original:      "https://example.com",
			dbExpectation: func() {},
			expectedURL:   "https://example.com/example",
			expectedErr:   nil,
		},
		{
			name:          "AddValue with MemoryStorage",
			short:         "example",
			original:      "https://example.com",
			dbExpectation: func() {},
			expectedURL:   "https://example.com/example",
			expectedErr:   nil,
		},
		{
			name:     "AddValue conflict - Short URL already exists",
			short:    "example",
			original: "https://example.com",
			dbExpectation: func() {
				mock.ExpectExec("INSERT INTO shortener").WithArgs("example", "https://example.com").WillReturnError(ErrConflict)
				mock.ExpectQuery("SELECT short FROM shortener").WithArgs("https://example.com").WillReturnRows(sqlmock.NewRows([]string{"short"}).AddRow("example"))
			},
			expectedURL: "https://example.com/example",
			expectedErr: ErrConflict,
		},
		{
			name:          "AddValue with an empty original URL",
			short:         "example",
			original:      "",
			dbExpectation: func() {},
			expectedURL:   "",
			expectedErr:   errors.New("original URL cannot be empty"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.dbExpectation()

			url, err := storage.AddValue(context.Background(), AddValueOptions{
				Original: tc.original,
				BaseURL:  "https://example.com",
				Short:    tc.short,
			})

			if url != tc.expectedURL {
				t.Errorf("Expected URL: %s, got: %s", tc.expectedURL, url)
			}

			if (err == nil && tc.expectedErr != nil) || (err != nil && err.Error() != tc.expectedErr.Error()) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectations were not met: %s", err)
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

func TestSaveURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred while creating mock database: %s", err)
	}
	defer db.Close()

	storage := Storage{
		mode:   DataBaseStorage,
		db:     db,
		values: make(map[string]string),
	}

	successExpectation := func(short, original string) {
		mock.ExpectExec("INSERT INTO shortener").WithArgs(short, original).WillReturnResult(sqlmock.NewResult(1, 1))
	}

	conflictExpectation := func(short, original string) {
		mock.ExpectExec("INSERT INTO shortener").WithArgs(short, original).WillReturnError(ErrConflict)
	}

	testCases := []struct {
		name        string
		short       string
		original    string
		expectation func(short string, original string)
		expectedErr error
	}{
		{
			name:        "SaveURL success",
			short:       "example",
			original:    "https://example.com",
			expectation: successExpectation,
			expectedErr: nil,
		},
		{
			name:        "SaveURL conflict",
			short:       "example",
			original:    "https://example.com",
			expectation: conflictExpectation,
			expectedErr: ErrConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectation(tc.short, tc.original)

			err := storage.SaveURL(context.Background(), tc.short, tc.original)

			if err != tc.expectedErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectations were not met: %s", err)
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
