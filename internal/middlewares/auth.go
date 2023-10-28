package middlewares

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/kupriyanovkk/shortener/internal/contextkey"
	"github.com/kupriyanovkk/shortener/internal/encrypt"
	"github.com/kupriyanovkk/shortener/internal/random"
)

func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := random.Generate(10)
		cookie, cookieErr := r.Cookie("UserID")
		encrypt, err := encrypt.Get()

		if err != nil {
			fmt.Printf("encrypt.Get() error: %v\n", err)
		} else {
			if cookieErr != nil {
				encrypted := encrypt.AEAD.Seal(nil, encrypt.Nonce, userID, nil)
				cookie := &http.Cookie{
					Name:  "UserID",
					Value: hex.EncodeToString(encrypted),
					Path:  "/",
				}
				http.SetCookie(w, cookie)
			} else {
				decode, _ := hex.DecodeString(cookie.Value)
				decrypted, err := encrypt.AEAD.Open(nil, encrypt.Nonce, decode, nil)
				if err != nil {
					fmt.Printf("encrypt.AEAD.Open error: %v\n", err)
				}
				userID = decrypted
			}
		}

		ctx := context.WithValue(r.Context(), contextkey.ContextUserKey, hex.EncodeToString(userID))
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
