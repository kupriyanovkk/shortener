package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/models"
	"github.com/kupriyanovkk/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultURL = "http://localhost:8080/"
var storageFile = "/tmp/short-url-db.json"
var dbDSN = ""
var f = config.ConfigFlags{
	B: defaultURL,
	F: storageFile,
	D: dbDSN,
}

func TestHandleFuncs(t *testing.T) {
	t.Run("Valid POST Request", func(t *testing.T) {
		body := []byte("https://example.com")
		s := storage.NewStorage(storageFile, dbDSN)
		env := &config.Env{Flags: f, Storage: s}
		req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { PostRootHandler(w, r, env) })

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Contains(t, rr.Header().Get("Location"), defaultURL)
	})

	t.Run("Invalid POST Request", func(t *testing.T) {
		body := []byte("invalid-url")
		s := storage.NewStorage(storageFile, dbDSN)
		env := &config.Env{Flags: f, Storage: s}
		req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { PostRootHandler(w, r, env) })

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Error parsing URL")
	})

	t.Run("Valid GET Request", func(t *testing.T) {
		id := "abc123"
		s := storage.NewStorage(storageFile, dbDSN)
		env := &config.Env{Flags: f, Storage: s}
		s.AddValue(context.Background(), id, "http://example.com")

		req, err := http.NewRequest(http.MethodGet, "/"+id, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { GetHandler(w, r, env) })

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "http://example.com", rr.Header().Get("Location"))
	})

	t.Run("Invalid GET Request (Not Found)", func(t *testing.T) {
		s := storage.NewStorage(storageFile, dbDSN)
		env := &config.Env{Flags: f, Storage: s}
		req, err := http.NewRequest(http.MethodGet, "/nonexistent", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { GetHandler(w, r, env) })

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "value doesn't exist by key nonexistent")
	})
}

func TestPostAPIHandler(t *testing.T) {
	s := storage.NewStorage(storageFile, dbDSN)
	env := &config.Env{Flags: f, Storage: s}
	body := []byte(`{"url":"http://example.com/"}`)
	req, err := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { PostAPIHandler(w, r, env) })

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Header().Get("Location"), defaultURL)

	var resp models.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("Error decoding response JSON: %v", err)
	}
	require.NotEmpty(t, resp, resp.Result)
}

func TestBatchHandler(t *testing.T) {
	s := storage.NewStorage(storageFile, dbDSN)
	env := &config.Env{Flags: f, Storage: s}

	testCases := []struct {
		Name         string
		Request      []models.BatchRequest
		ExpectedCode int
	}{
		{
			Name: "ValidInput",
			Request: []models.BatchRequest{
				{
					CorrelationID: "123",
					OriginalURL:   "https://example.org",
				},
			},
			ExpectedCode: http.StatusCreated,
		},
		{
			Name:         "EmptyRequestBody",
			Request:      []models.BatchRequest{},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:         "InvalidJSON",
			Request:      nil,
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name: "InvalidURL",
			Request: []models.BatchRequest{
				{
					CorrelationID: "123",
					OriginalURL:   "not_a_valid_url",
				},
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name: "InvalidBaseURL",
			Request: []models.BatchRequest{
				{
					CorrelationID: "123",
					OriginalURL:   "https://example.org",
				},
			},
			ExpectedCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.Request)
			req := httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer(reqBody))
			rec := httptest.NewRecorder()

			BatchHandler(rec, req, env)

			if rec.Code != tc.ExpectedCode {
				t.Errorf("Test case '%s' failed. Expected status code %d, but got %d", tc.Name, tc.ExpectedCode, rec.Code)
			}
		})
	}
}
