package middlewares

import (
	"net/http"
	"strings"

	"github.com/kupriyanovkk/shortener/internal/compress"
	"github.com/kupriyanovkk/shortener/internal/utils"
)

func Gzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		compressibleTypes := []string{"application/json", "text/html"}
		contentType := r.Header.Get("Content-Type")

		supportsContentType := utils.Contains(compressibleTypes, contentType)
		supportsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

		if supportsGzip && supportsContentType {
			cw := compress.NewWriter(w)
			ow = cw
			defer cw.Close()
		}

		h.ServeHTTP(ow, r)
	})
}
