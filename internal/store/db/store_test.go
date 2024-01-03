package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kupriyanovkk/shortener/internal/failure"
	"github.com/kupriyanovkk/shortener/internal/models"
	storeInterface "github.com/kupriyanovkk/shortener/internal/store/interface"
)

func TestAddValue(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred while creating mock database: %s", err)
	}
	defer db.Close()

	storage := Store{
		db: db,
	}

	successExpectation := func(short, original, user string) {
		mock.ExpectExec("INSERT INTO shortener").WithArgs(short, original, user, false).WillReturnResult(sqlmock.NewResult(1, 1))
	}

	testCases := []struct {
		name          string
		short         string
		original      string
		user          string
		dbExpectation func(short, original, user string)
		expectedURL   string
		expectedErr   error
	}{
		{
			name:          "AddValue success",
			short:         "example",
			original:      "https://example.com",
			user:          "123",
			dbExpectation: successExpectation,
			expectedURL:   "https://example.com/example",
			expectedErr:   nil,
		},
		{
			name:     "AddValue conflict",
			short:    "example",
			original: "https://example.com",
			user:     "123",
			dbExpectation: func(short, original, user string) {
				mock.ExpectExec("INSERT INTO shortener").WithArgs(short, original, user, false).WillReturnError(failure.ErrConflict)
				mock.ExpectQuery("SELECT short FROM shortener").WithArgs(original).WillReturnRows(sqlmock.NewRows([]string{"short"}).AddRow(short))
			},
			expectedURL: "https://example.com/example",
			expectedErr: failure.ErrConflict,
		},
		{
			name:          "AddValue with FileWriterStorage",
			short:         "example",
			original:      "https://example.com",
			user:          "123",
			dbExpectation: func(short, original, user string) {},
			expectedURL:   "https://example.com/example",
			expectedErr:   nil,
		},
		{
			name:          "AddValue with MemoryStorage",
			short:         "example",
			original:      "https://example.com",
			user:          "123",
			dbExpectation: func(short, original, user string) {},
			expectedURL:   "https://example.com/example",
			expectedErr:   nil,
		},
		{
			name:     "AddValue conflict - Short URL already exists",
			short:    "example",
			original: "https://example.com",
			user:     "123",
			dbExpectation: func(short, original, user string) {
				mock.ExpectExec("INSERT INTO shortener").WithArgs(short, original, user, false).WillReturnError(failure.ErrConflict)
				mock.ExpectQuery("SELECT short FROM shortener").WithArgs(original).WillReturnRows(sqlmock.NewRows([]string{"short"}).AddRow(short))
			},
			expectedURL: "https://example.com/example",
			expectedErr: failure.ErrConflict,
		},
		{
			name:          "AddValue with an empty original URL",
			short:         "example",
			user:          "123",
			original:      "",
			dbExpectation: func(short, original, user string) {},
			expectedURL:   "",
			expectedErr:   failure.ErrEmptyOrigURL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.dbExpectation(tc.short, tc.original, tc.user)

			url, err := storage.AddValue(context.Background(), storeInterface.AddValueOptions{
				Original: tc.original,
				BaseURL:  "https://example.com",
				Short:    tc.short,
				UserID:   tc.user,
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
		db: db,
	}

	successExpectation := func(short, original, user string) {
		mock.ExpectExec("INSERT INTO shortener").WithArgs(short, original, user, false).WillReturnResult(sqlmock.NewResult(1, 1))
	}

	conflictExpectation := func(short, original, user string) {
		mock.ExpectExec("INSERT INTO shortener").WithArgs(short, original, user, false).WillReturnError(failure.ErrConflict)
	}

	testCases := []struct {
		name        string
		short       string
		original    string
		user        string
		expectation func(short string, original string, user string)
		expectedErr error
	}{
		{
			name:        "SaveURL success",
			short:       "example",
			original:    "https://example.com",
			user:        "123",
			expectation: successExpectation,
			expectedErr: nil,
		},
		{
			name:        "SaveURL conflict",
			short:       "example",
			original:    "https://example.com",
			user:        "123",
			expectation: conflictExpectation,
			expectedErr: failure.ErrConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectation(tc.short, tc.original, tc.user)

			err := storage.InsertURL(context.Background(), tc.short, tc.original, tc.user)

			if err != tc.expectedErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectations were not met: %s", err)
			}
		})
	}
}

func TestGetUserURLs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred while creating mock database: %s", err)
	}
	defer db.Close()

	s := Store{
		db: db,
	}

	userID := "testUserID"
	baseURL := "http://example.com"
	opts := storeInterface.GetUserURLsOptions{
		UserID:  userID,
		BaseURL: baseURL,
	}

	rows := sqlmock.NewRows([]string{"original", "short"}).
		AddRow("http://example.com/original1", "short1").
		AddRow("http://example.com/original2", "short2")

	mock.ExpectQuery("SELECT original, short FROM shortener").
		WithArgs(userID, 100).
		WillReturnRows(rows)

	urls, err := s.GetUserURLs(context.Background(), opts)

	if err != nil {
		t.Errorf("Error was not expected, got: %v", err)
	}

	expectedURLs := []models.UserURL{
		{
			Short:    "http://example.com/short1",
			Original: "http://example.com/original1",
		},
		{
			Short:    "http://example.com/short2",
			Original: "http://example.com/original2",
		},
	}

	if len(urls) != len(expectedURLs) {
		t.Errorf("Expected %d URLs, but got %d", len(expectedURLs), len(urls))
	}

	for i, u := range urls {
		if u != expectedURLs[i] {
			t.Errorf("Expected URL %d to be %v, but got %v", i, expectedURLs[i], u)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
