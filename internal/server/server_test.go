package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kupriyanovkk/shortener/internal/models"
	"github.com/kupriyanovkk/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandleFunc(t *testing.T) {
	defaultURL := "http://localhost:8080/"

	t.Run("Valid POST Request", func(t *testing.T) {
		body := []byte("https://example.com")
		s := storage.NewStorage()
		req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { PostRootHandler(w, r, s, defaultURL) })

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Contains(t, rr.Header().Get("Location"), "http://localhost:8080/")
	})

	t.Run("Invalid POST Request", func(t *testing.T) {
		body := []byte("invalid-url")
		s := storage.NewStorage()
		req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { PostRootHandler(w, r, s, defaultURL) })

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Error parsing URL")
	})

	t.Run("Valid GET Request", func(t *testing.T) {
		id := "abc123"
		s := storage.NewStorage()
		s.AddValue(id, "http://example.com")

		req, err := http.NewRequest(http.MethodGet, "/"+id, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { GetHandler(w, r, s) })

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, "http://example.com", rr.Header().Get("Location"))
	})

	t.Run("Invalid GET Request (Not Found)", func(t *testing.T) {
		s := storage.NewStorage()
		req, err := http.NewRequest(http.MethodGet, "/nonexistent", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { GetHandler(w, r, s) })

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "value doesn't exist by key nonexistent")
	})

	t.Run("Valid POST Request to API", func(t *testing.T) {
		s := storage.NewStorage()
		body := []byte(`{"url":"http://example.com/"}`)
		req, err := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { PostAPIHandler(w, r, s, defaultURL) })

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Contains(t, rr.Header().Get("Location"), "http://localhost:8080/")

		var resp models.Response
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Errorf("Error decoding response JSON: %v", err)
		}
	})
}
