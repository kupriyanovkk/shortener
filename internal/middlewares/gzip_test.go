package middlewares

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kupriyanovkk/shortener/internal/config"
	"github.com/kupriyanovkk/shortener/internal/handlers"
	"github.com/kupriyanovkk/shortener/internal/models"
	"github.com/kupriyanovkk/shortener/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestGzip(t *testing.T) {
	defaultURL := "http://localhost:8080/"
	storageFile := "/tmp/short-url-db.json"
	dbDSN := ""

	f := config.ConfigFlags{
		B: defaultURL,
		F: storageFile,
		D: dbDSN,
	}

	// Helper function to create a compressed request body
	createCompressedRequestBody := func(input string) []byte {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(input))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)
		return buf.Bytes()
	}

	// Helper function to send a compressed request
	sendCompressedRequest := func(t *testing.T, srv *httptest.Server, requestBody []byte) (*http.Response, []byte) {
		r := httptest.NewRequest("POST", srv.URL, bytes.NewReader(requestBody))
		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")
		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		return resp, b
	}

	t.Run("sends gzip", func(t *testing.T) {
		s := storage.NewStorage(storageFile, dbDSN)
		env := &config.Env{Flags: f, Storage: s}
		handler := Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.PostAPIShorten(w, r, env)
		}))
		srv := httptest.NewServer(handler)
		defer srv.Close()

		requestBody := `{"url":"http://example.com/"}`
		compressedBody := createCompressedRequestBody(requestBody)

		resp, responseBody := sendCompressedRequest(t, srv, compressedBody)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var model models.Response
		require.NoError(t, json.Unmarshal(responseBody, &model))
		require.NotEmpty(t, model.Result)
	})

	t.Run("accepts gzip", func(t *testing.T) {
		s := storage.NewStorage(storageFile, dbDSN)
		env := &config.Env{Flags: f, Storage: s}
		handler := Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.PostAPIShorten(w, r, env)
		}))
		srv := httptest.NewServer(handler)
		defer srv.Close()

		requestBody := `{"url":"http://example.com/"}`
		compressedBody := createCompressedRequestBody(requestBody)

		r := httptest.NewRequest("POST", srv.URL, bytes.NewReader(compressedBody))
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")

		resp, responseBody := sendCompressedRequest(t, srv, compressedBody)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var model models.Response
		require.NoError(t, json.Unmarshal(responseBody, &model))
		require.NotEmpty(t, model.Result)
	})

	t.Run("no gzip", func(t *testing.T) {
		s := storage.NewStorage(storageFile, dbDSN)
		env := &config.Env{Flags: f, Storage: s}
		handler := Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.PostAPIShorten(w, r, env)
		}))
		srv := httptest.NewServer(handler)
		defer srv.Close()

		requestBody := `{"url":"http://example.com/"}`
		buf := bytes.NewBufferString(requestBody)
		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var model models.Response
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&model))
		require.NotEmpty(t, model.Result)
	})
}
