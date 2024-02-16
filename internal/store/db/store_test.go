package db

import (
	"context"
	"errors"
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

func TestFindOriginalURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := Store{db: db}

	cases := []struct {
		short       string
		expectedURL models.URL
		expectedErr error
	}{
		{
			short: "example_short",
			expectedURL: models.URL{
				Original:    "example_original",
				DeletedFlag: false,
			},
			expectedErr: nil,
		},
	}

	for _, c := range cases {
		mock.ExpectQuery("SELECT original, is_deleted FROM shortener WHERE short = ?").
			WithArgs(c.short).
			WillReturnRows(sqlmock.NewRows([]string{"original", "is_deleted"}).AddRow(c.expectedURL.Original, c.expectedURL.DeletedFlag))

		url, err := s.FindOriginalURL(context.Background(), c.short)
		if err != c.expectedErr {
			t.Errorf("expected error %v, but got %v", c.expectedErr, err)
		}
		if url != c.expectedURL {
			t.Errorf("expected URL %+v, but got %+v", c.expectedURL, url)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteURLs(t *testing.T) {
	// Test case for deleting a single URL
	t.Run("DeleteSingleURL", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		s := Store{db: db}

		opts := []storeInterface.DeletedURLs{
			{
				URLs:   []string{"example1.com"},
				UserID: "user1",
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE shortener SET is_deleted = TRUE").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := s.DeleteURLs(context.Background(), opts)
		if err != nil {
			t.Errorf("Failed to delete single URL: %v", err)
		}
	})

	// Test case for deleting multiple URLs
	t.Run("DeleteMultipleURLs", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		s := Store{db: db}

		opts := []storeInterface.DeletedURLs{
			{
				URLs:   []string{"example2.com", "example3.com"},
				UserID: "user2",
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE shortener").WithArgs("example2.com", "user2").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("UPDATE shortener").WithArgs("example3.com", "user2").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := s.DeleteURLs(context.Background(), opts)
		if err != nil {
			t.Errorf("Failed to delete multiple URLs: %v", err)
		}
	})
}

func TestGetInternalStats(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	store := Store{db: db}

	mock.ExpectQuery("SELECT COUNT(*) FROM shortener").WillReturnError(errors.New("database error"))
	_, err = store.GetInternalStats(ctx)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	mock.ExpectQuery("SELECT COUNT(*) FROM shortener").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
	mock.ExpectQuery("SELECT user_id, COUNT(user_id) FROM shortener GROUP BY user_id").WillReturnError(errors.New("database error"))
	_, err = store.GetInternalStats(ctx)
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}
