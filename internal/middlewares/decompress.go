package middlewares

import (
	"io"
	"net/http"
	"strings"

	"github.com/kupriyanovkk/shortener/internal/compress"
)

func Decompress(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			cr, err := compress.NewReader(r.Body)
			if (err) != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer cr.Close()

			body, err := io.ReadAll(cr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(strings.NewReader(string(body)))
		}

		h.ServeHTTP(w, r)
	})
}
