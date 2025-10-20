package middleware

import (
	"net"
	"net/http"

	"github.com/felipeosantos/goexpert/rate-limiter/internal/limiter"
)

const (
	// APIKeyHeader is the header name for the API key
	APIKeyHeader = "API_KEY"

	// RateLimitExceededMessage is the message shown when rate limit is exceeded
	RateLimitExceededMessage = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

// RateLimiterMiddleware creates a middleware for rate limiting
func RateLimiterMiddleware(limiter *limiter.RateLimiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get IP address
			ip := getIPAddress(r)

			// Get token from header
			token := r.Header.Get(APIKeyHeader)

			// Check if request is allowed
			allowed, err := limiter.Allow(r.Context(), ip, token)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				w.WriteHeader(http.StatusTooManyRequests) // 429 Too Many Requests
				w.Write([]byte(RateLimitExceededMessage))
				return
			}

			// Pass to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// getIPAddress returns the client's IP address from the request
func getIPAddress(r *http.Request) string {
	// Check for X-Forwarded-For header first (when behind a proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// If not behind a proxy, use RemoteAddr but strip the port
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If there's an error splitting (perhaps because there's no port),
		// just use the original RemoteAddr
		return r.RemoteAddr
	}
	return ip
}
