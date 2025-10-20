package limiter

import (
	"context"

	"github.com/felipeosantos/goexpert/rate-limiter/config"
	"github.com/felipeosantos/goexpert/rate-limiter/internal/storage"
)

// Config holds rate limiter configuration
// type Config struct {
// 	// IPRateLimit defines the maximum requests per second for an IP
// 	IPRateLimit int
// 	// TokenRateLimit defines the maximum requests per second for a token
// 	TokenRateLimit int
// 	// BlockDuration defines how long an IP or token will be blocked after exceeding limits
// 	BlockDuration time.Duration
// 	// Window defines the time window for rate limiting
// 	Window time.Duration
// 	// TokenConfigs maps token names to their specific configurations
// 	TokenConfigs map[string]config.TokenConfig
// }

type Config struct {
	IP    config.LimiterConfig
	Token map[string]config.LimiterConfig
}

// RateLimiter manages rate limiting logic
type RateLimiter struct {
	storage storage.Storage
	config  Config
}

// New creates a new rate limiter with the provided storage and configuration
func New(storage storage.Storage, config Config) *RateLimiter {
	return &RateLimiter{
		storage: storage,
		config:  config,
	}
}

// Allow checks if a request is allowed based on IP and token
func (rl *RateLimiter) Allow(ctx context.Context, ip string, token string) (bool, error) {
	// First check if IP or token is blocked
	ipBlocked, err := rl.storage.IsBlocked(ctx, "ip:"+ip)
	if err != nil {
		return false, err
	}

	if ipBlocked {
		return false, nil
	}

	// If token is provided, check if it's blocked
	if token != "" {
		tokenBlocked, err := rl.storage.IsBlocked(ctx, "token:"+token)
		if err != nil {
			return false, err
		}

		if tokenBlocked {
			return false, nil
		}

		// If token provided and not blocked, check token limit
		return rl.checkTokenLimit(ctx, token, ip)
	}

	// If no token, check IP limit
	return rl.checkIPLimit(ctx, ip)
}

// checkIPLimit checks if the IP has exceeded its limit
func (rl *RateLimiter) checkIPLimit(ctx context.Context, ip string) (bool, error) {
	ipKey := "ip:" + ip
	count, err := rl.storage.Increment(ctx, ipKey, rl.config.IP.RateWindow)
	if err != nil {
		return false, err
	}

	// If IP exceeds rate limit, block it
	if count > rl.config.IP.RateLimit {
		err = rl.storage.Block(ctx, ipKey, rl.config.IP.BlockDuration)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

// checkTokenLimit checks if the token has exceeded its limit
func (rl *RateLimiter) checkTokenLimit(ctx context.Context, token, ip string) (bool, error) {
	tokenKey := "token:" + token

	// Check if this token has specific configurations
	tokenConfig, hasCustomConfig := rl.config.Token[token]

	// Determine which rate limit to use
	rateLimit := rl.config.IP.RateLimit
	rateWindow := rl.config.IP.RateWindow
	blockDuration := rl.config.IP.BlockDuration

	if hasCustomConfig {
		rateLimit = tokenConfig.RateLimit
		rateWindow = tokenConfig.RateWindow
		blockDuration = tokenConfig.BlockDuration
	}

	count, err := rl.storage.Increment(ctx, tokenKey, rateWindow)
	if err != nil {
		return false, err
	}

	// If token exceeds rate limit, block both token and IP
	if count > rateLimit {
		err = rl.storage.Block(ctx, tokenKey, blockDuration)
		if err != nil {
			return false, err
		}

		err = rl.storage.Block(ctx, "ip:"+ip, blockDuration)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// Close closes the underlying storage
func (rl *RateLimiter) Close() error {
	return rl.storage.Close()
}
