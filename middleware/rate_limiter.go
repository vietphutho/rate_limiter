package middleware

import (
	"net/http"

	"rate-limit.vietnt/limiter"
)

// RateLimitMiddleware returns a middleware that applies the provided rate limiter.
func RateLimitMiddleware(rateLimiter limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get userID from header or fallback to remote IP
			userID := r.Header.Get("X-User-ID")
			if userID == "" {
				userID = r.RemoteAddr
			}

			// Apply limiter
			if !rateLimiter.Allow(userID) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Pass to next handler
			next.ServeHTTP(w, r)
		})
	}
}
