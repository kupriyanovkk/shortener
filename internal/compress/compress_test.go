package compress

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompressWriterAndReader(t *testing.T) {
	testCases := []struct {
		name        string
		inputData   []byte
		expectedErr error
	}{
		{
			name:        "Valid data",
			inputData:   []byte("Hello, world!"),
			expectedErr: nil,
		},
		{
			name:        "Empty data",
			inputData:   []byte(""),
			expectedErr: nil,
		},
		{
			name:        "Invalid data",
			inputData:   []byte("Invalid gzip data"),
			expectedErr: io.ErrUnexpectedEOF,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			cw := NewWriter(w)

			_, err := cw.Write(tc.inputData)
			if err != nil {
				t.Fatalf("Error writing data: %v", err)
			}

			err = cw.Close()
			if err != nil {
				t.Fatalf("Error closing CompressWriter: %v", err)
			}

			cr, err := NewReader(io.NopCloser(bytes.NewReader(w.Body.Bytes())))
			if err != nil {
				if tc.expectedErr != err {
					t.Fatalf("Expected error: %v, got: %v", tc.expectedErr, err)
				}
				return
			}
			defer cr.Close()

			decompressed, err := io.ReadAll(cr)
			if err != nil {
				t.Fatalf("Error reading decompressed data: %v", err)
			}

			if !bytes.Equal(decompressed, tc.inputData) {
				t.Fatalf("Decompressed data does not match original data")
			}
		})
	}
}

func TestCompressWriter_Header(t *testing.T) {
	w := httptest.NewRecorder()

	cw := NewWriter(w)

	cw.WriteHeader(http.StatusOK)

	if encoding := w.Header().Get("Content-Encoding"); encoding != "gzip" {
		t.Fatalf("Content-Encoding header not set correctly")
	}
}

func TestNewReader_Error(t *testing.T) {
	invalidData := []byte("Invalid gzip data")
	_, err := NewReader(io.NopCloser(bytes.NewReader(invalidData)))

	if err == nil {
		t.Fatalf("Expected an error when creating CompressReader with invalid data")
	}
}
