package middleware

import "net/http"

// constants for the X-Frame-Options header
const (
	XFrameOptionsSameOrigin = "sameorigin"
	XFrameOptionsDeny       = "deny"
)

// FrameOptions is middleware that sets the X-Frame-Options header. Policy should
// be one of deny|sameorigin|allow-from https://example.com/
func FrameOptions(policy string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Frame-Options", policy)
			next.ServeHTTP(w, r)
		})
	}
}
