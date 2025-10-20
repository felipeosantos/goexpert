package storage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/felipeosantos/goexpert/rate-limiter/config"
)

var (
	ErrStorageNotFound = errors.New("storage not found")
)

var (
	registry   = make(map[string]func(config.StorageConfig) (Storage, error))
	registryMu sync.RWMutex
)

// Storage Strategy defines the interface for rate limiter storage backends
type Storage interface {
	// Get returns the current count for a key
	Get(ctx context.Context, key string) (int, error)

	// Increment increments the counter for a key and returns the new value
	// If the key doesn't exist, it creates it with the given expiration
	Increment(ctx context.Context, key string, expiration time.Duration) (int, error)

	// Reset resets the counter for a key
	Reset(ctx context.Context, key string) error

	// IsBlocked checks if a key is in the blocklist
	IsBlocked(ctx context.Context, key string) (bool, error)

	// Block adds a key to the blocklist with the given expiration
	Block(ctx context.Context, key string, expiration time.Duration) error

	// Close closes the storage connection
	Close() error
}

func Register(storageName string, storageConstructor func(config.StorageConfig) (Storage, error)) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[storageName] = storageConstructor
}

func New(storageName string, storageCfg config.StorageConfig) (Storage, error) {
	registryMu.RLock()
	factory, exists := registry[storageName]
	registryMu.RUnlock()

	if !exists {
		return nil, ErrStorageNotFound
	}
	return factory(storageCfg)
}
