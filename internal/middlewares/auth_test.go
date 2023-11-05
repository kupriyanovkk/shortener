package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	t.Run("Test Auth middleware with valid UserID cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/user/urls", nil)
		rr := httptest.NewRecorder()

		Auth(handler).ServeHTTP(rr, req)

		resp := rr.Result()
		resp.Body.Close()

		cookies := resp.Cookies()

		if len(cookies) != 1 {
			t.Error("Expected one cookie to be set, but none found.")
		} else {
			cookie := cookies[0]
			if cookie.Name != "UserID" {
				t.Errorf("Expected cookie name 'UserID', but got '%s'", cookie.Name)
			}
		}
	})
}
