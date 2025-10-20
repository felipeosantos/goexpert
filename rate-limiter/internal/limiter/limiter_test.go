package limiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/felipeosantos/goexpert/rate-limiter/config"
	"github.com/felipeosantos/goexpert/rate-limiter/internal/limiter"
	"github.com/felipeosantos/goexpert/rate-limiter/internal/storage"
)

func TestRateLimiter(t *testing.T) {
	// Create a memory storage for testing
	store := storage.NewMemoryStorage()

	// Import config package for TokenConfig
	customTokenConfigs := map[string]config.LimiterConfig{
		"premium-token": {
			RateLimit:     10,
			RateWindow:    1 * time.Second,
			BlockDuration: 30 * time.Second,
		},
		"premium-token2": {
			RateLimit:     10,
			RateWindow:    1 * time.Second,
			BlockDuration: 30 * time.Second,
		},
		"basic-token": {
			RateLimit:     3,
			RateWindow:    1 * time.Second,
			BlockDuration: 2 * time.Minute,
		},
	}

	// Create a rate limiter with a small window and limits
	rl := limiter.New(store, limiter.Config{
		IP: config.LimiterConfig{
			RateLimit:     3,
			RateWindow:    time.Second,
			BlockDuration: time.Minute,
		},
		Token: customTokenConfigs,
	})

	// Test IP rate limiting
	t.Run("IP rate limiting", func(t *testing.T) {
		ip := "192.168.1.1"

		// First three requests should be allowed
		for i := 0; i < 3; i++ {
			allowed, err := rl.Allow(context.Background(), ip, "")
			if err != nil {
				t.Fatalf("Error checking rate limit: %v", err)
			}
			if !allowed {
				t.Errorf("Request %d should be allowed", i+1)
			}
		}

		// Fourth request should be blocked
		allowed, err := rl.Allow(context.Background(), ip, "")
		if err != nil {
			t.Fatalf("Error checking rate limit: %v", err)
		}
		if allowed {
			t.Errorf("Request should be blocked after exceeding limit")
		}
	})

	// Test token rate limiting
	t.Run("Token rate limiting", func(t *testing.T) {
		ip := "192.168.1.2"
		token := "test-token"

		// First three requests should be allowed
		for i := 0; i < 3; i++ {
			allowed, err := rl.Allow(context.Background(), ip, token)
			if err != nil {
				t.Fatalf("Error checking rate limit: %v", err)
			}
			if !allowed {
				t.Errorf("Request %d should be allowed", i+1)
			}
		}

		// Fourth request should be blocked
		allowed, err := rl.Allow(context.Background(), ip, token)
		if err != nil {
			t.Fatalf("Error checking rate limit: %v", err)
		}
		if allowed {
			t.Errorf("Request should be blocked after exceeding token limit")
		}
	})

	// Test token precedence over IP
	t.Run("Token precedence over IP", func(t *testing.T) {
		ip := "192.168.1.3"
		token := "premium-token"

		// Make some IP requests first
		for i := 0; i < 2; i++ {
			_, _ = rl.Allow(context.Background(), ip, "")
		}

		// Now use the token - should still have full token limit
		for i := 0; i < 10; i++ {
			allowed, err := rl.Allow(context.Background(), ip, token)
			if err != nil {
				t.Fatalf("Error checking rate limit: %v", err)
			}
			if !allowed {
				t.Errorf("Request %d with token should be allowed despite IP usage", i+1)
			}
		}
	})

	// Test token-specific configurations
	t.Run("Token-specific configurations", func(t *testing.T) {
		// Test premium token with higher limit
		t.Run("Premium token", func(t *testing.T) {
			ip := "192.168.1.4"
			token := "premium-token2"

			// First 10 requests should be allowed (custom limit)
			for i := 0; i < 10; i++ {
				allowed, err := rl.Allow(context.Background(), ip, token)
				if err != nil {
					t.Fatalf("Error checking rate limit: %v", err)
				}
				if !allowed {
					t.Errorf("Request %d should be allowed for premium token", i+1)
				}
			}

			// 11th request should be blocked
			allowed, err := rl.Allow(context.Background(), ip, token)
			if err != nil {
				t.Fatalf("Error checking rate limit: %v", err)
			}
			if allowed {
				t.Errorf("Request should be blocked after exceeding custom token limit")
			}
		})

		// Test basic token with lower limit
		t.Run("Basic token", func(t *testing.T) {
			ip := "192.168.1.5"
			token := "basic-token"

			// First 3 requests should be allowed (custom limit)
			for i := 0; i < 3; i++ {
				allowed, err := rl.Allow(context.Background(), ip, token)
				if err != nil {
					t.Fatalf("Error checking rate limit: %v", err)
				}
				if !allowed {
					t.Errorf("Request %d should be allowed for basic token", i+1)
				}
			}

			// 4th request should be blocked
			allowed, err := rl.Allow(context.Background(), ip, token)
			if err != nil {
				t.Fatalf("Error checking rate limit: %v", err)
			}
			if allowed {
				t.Errorf("Request should be blocked after exceeding custom token limit")
			}
		})
	})
}
