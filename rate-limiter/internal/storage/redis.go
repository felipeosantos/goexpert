package storage

import (
	"context"
	"time"

	"github.com/felipeosantos/goexpert/rate-limiter/config"
	"github.com/redis/go-redis/v9"
)

// RedisStorage implements the Storage Strategy interface using Redis
type RedisStorage struct {
	client *redis.Client
}

// NewRedis creates a new Redis storage
func NewRedis(redisCfg config.StorageConfig) (*RedisStorage, error) {
	opts, err := redis.ParseURL(redisCfg.URL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	// Ping Redis to verify connection
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &RedisStorage{client: client}, nil
}

// Get returns the current count for a key
func (s *RedisStorage) Get(ctx context.Context, key string) (int, error) {
	val, err := s.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// Increment increments the counter for a key and returns the new value
func (s *RedisStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
	val, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// Set expiration only on the first increment
	if val == 1 {
		_, err = s.client.Expire(ctx, key, expiration).Result()
		if err != nil {
			return int(val), err
		}
	}

	return int(val), nil
}

// Reset resets the counter for a key
func (s *RedisStorage) Reset(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

// IsBlocked checks if a key is in the blocklist
func (s *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	blocklistKey := "blocklist:" + key
	exists, err := s.client.Exists(ctx, blocklistKey).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// Block adds a key to the blocklist with the given expiration
func (s *RedisStorage) Block(ctx context.Context, key string, expiration time.Duration) error {
	blocklistKey := "blocklist:" + key
	return s.client.Set(ctx, blocklistKey, 1, expiration).Err()
}

// Close closes the Redis connection
func (s *RedisStorage) Close() error {
	return s.client.Close()
}

func init() {
	Register("redis", func(cfg config.StorageConfig) (Storage, error) {
		return NewRedis(cfg)
	})
}
