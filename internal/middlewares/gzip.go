package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// Header method return compressWriter header
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write method return zw.Write
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader method set Content-Encoding header
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close method call zw.Close
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// NewCompressWriter return pointer to compressWriter
func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// Read method return zr.Read result
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close method call zr.Close
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// NewCompressReader return pointer to compressReader
func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Gzip is middleware for zip/unzip requests body.
func Gzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		compressibleTypes := []string{"application/json", "text/html"}
		contentType := r.Header.Get("Content-Type")

		supportsContentType := contains(compressibleTypes, contentType)
		supportsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

		if supportsGzip && supportsContentType {
			cw := NewCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}
