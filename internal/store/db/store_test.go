package db

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kupriyanovkk/shortener/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	s := Store{
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

func TestAddValue(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred while creating mock database: %s", err)
	}
	defer db.Close()

	storage := Store{
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
				mock.ExpectExec("INSERT INTO shortener").WithArgs("example", "https://example.com").WillReturnError(store.ErrConflict)
				mock.ExpectQuery("SELECT short FROM shortener").WithArgs("https://example.com").WillReturnRows(sqlmock.NewRows([]string{"short"}).AddRow("example"))
			},
			expectedURL: "https://example.com/example",
			expectedErr: store.ErrConflict,
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
				mock.ExpectExec("INSERT INTO shortener").WithArgs("example", "https://example.com").WillReturnError(store.ErrConflict)
				mock.ExpectQuery("SELECT short FROM shortener").WithArgs("https://example.com").WillReturnRows(sqlmock.NewRows([]string{"short"}).AddRow("example"))
			},
			expectedURL: "https://example.com/example",
			expectedErr: store.ErrConflict,
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

			url, err := storage.AddValue(context.Background(), store.AddValueOptions{
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

func TestSaveURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred while creating mock database: %s", err)
	}
	defer db.Close()

	storage := Store{
		db:     db,
		values: make(map[string]string),
	}

	successExpectation := func(short, original string) {
		mock.ExpectExec("INSERT INTO shortener").WithArgs(short, original).WillReturnResult(sqlmock.NewResult(1, 1))
	}

	conflictExpectation := func(short, original string) {
		mock.ExpectExec("INSERT INTO shortener").WithArgs(short, original).WillReturnError(store.ErrConflict)
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
			expectedErr: store.ErrConflict,
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
