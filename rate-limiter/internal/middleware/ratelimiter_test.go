package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/felipeosantos/goexpert/rate-limiter/config"
	"github.com/felipeosantos/goexpert/rate-limiter/internal/limiter"
	"github.com/felipeosantos/goexpert/rate-limiter/internal/middleware"
	"github.com/felipeosantos/goexpert/rate-limiter/internal/storage"
)

func TestRateLimiterMiddleware(t *testing.T) {
	// Create a memory storage for testing
	store := storage.NewMemoryStorage()

	// Create a rate limiter with a small window and limits
	rl := limiter.New(store, limiter.Config{
		IP: config.LimiterConfig{
			RateLimit:     2,
			RateWindow:    time.Second,
			BlockDuration: time.Minute,
		},
	})

	// Create test handler that just returns 200 OK
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Apply middleware
	middlewareHandler := middleware.RateLimiterMiddleware(rl)(handler)

	// Test IP rate limiting
	t.Run("IP rate limiting", func(t *testing.T) {
		// First two requests should be allowed
		for i := 0; i < 2; i++ {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = "192.168.1.100"
			rr := httptest.NewRecorder()

			middlewareHandler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Request %d should be allowed, got: %d", i+1, rr.Code)
			}
		}

		// Third request should be blocked
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.100"
		rr := httptest.NewRecorder()

		middlewareHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusTooManyRequests {
			t.Errorf("Request should be blocked, got: %d", rr.Code)
		}
	})

	// Test token rate limiting
	t.Run("Token rate limiting", func(t *testing.T) {
		// First two requests should be allowed
		for i := 0; i < 2; i++ {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = "192.168.1.101"
			req.Header.Set("API_KEY", "test-token-123")
			rr := httptest.NewRecorder()

			middlewareHandler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Request %d should be allowed, got: %d", i+1, rr.Code)
			}
		}

		// Third request should be blocked
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.101"
		req.Header.Set("API_KEY", "test-token-123")
		rr := httptest.NewRecorder()

		middlewareHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusTooManyRequests {
			t.Errorf("Request should be blocked, got: %d", rr.Code)
		}

		// Message should match expected
		if rr.Body.String() != middleware.RateLimitExceededMessage {
			t.Errorf("Expected message %q, got %q", middleware.RateLimitExceededMessage, rr.Body.String())
		}
	})

	// Test different IPs don't affect each other
	t.Run("Different IPs", func(t *testing.T) {
		// This IP should not be affected by previous tests
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.102"
		rr := httptest.NewRecorder()

		middlewareHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Request from new IP should be allowed, got: %d", rr.Code)
		}
	})
}
