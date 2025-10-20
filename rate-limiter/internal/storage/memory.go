package storage

import (
	"context"
	"sync"
	"time"

	"github.com/felipeosantos/goexpert/rate-limiter/config"
)

// Item represents a rate limiter item with expiration
type Item struct {
	Count     int
	ExpiresAt time.Time
}

// MemoryStorage implements the Storage interface using in-memory maps
type MemoryStorage struct {
	mu        sync.RWMutex
	items     map[string]Item
	blocklist map[string]time.Time
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		items:     make(map[string]Item),
		blocklist: make(map[string]time.Time),
	}
}

// Get returns the current count for a key
func (s *MemoryStorage) Get(ctx context.Context, key string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Clean expired items
	s.cleanExpired(key)

	if item, found := s.items[key]; found {
		return item.Count, nil
	}
	return 0, nil
}

// Increment increments the counter for a key and returns the new value
func (s *MemoryStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Clean expired items
	s.cleanExpired(key)

	item, found := s.items[key]
	if !found {
		s.items[key] = Item{
			Count:     1,
			ExpiresAt: time.Now().Add(expiration),
		}
		return 1, nil
	}

	item.Count++
	s.items[key] = item
	return item.Count, nil
}

// Reset resets the counter for a key
func (s *MemoryStorage) Reset(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key)
	return nil
}

// IsBlocked checks if a key is in the blocklist
func (s *MemoryStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Clean expired blocks
	s.cleanExpiredBlocks(key)

	_, blocked := s.blocklist[key]
	return blocked, nil
}

// Block adds a key to the blocklist with the given expiration
func (s *MemoryStorage) Block(ctx context.Context, key string, expiration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blocklist[key] = time.Now().Add(expiration)
	return nil
}

// Close closes the storage (no-op for memory storage)
func (s *MemoryStorage) Close() error {
	return nil
}

// cleanExpired removes expired items
func (s *MemoryStorage) cleanExpired(key string) {
	if item, found := s.items[key]; found && time.Now().After(item.ExpiresAt) {
		delete(s.items, key)
	}
}

// cleanExpiredBlocks removes expired blocks
func (s *MemoryStorage) cleanExpiredBlocks(key string) {
	if expireTime, found := s.blocklist[key]; found && time.Now().After(expireTime) {
		delete(s.blocklist, key)
	}
}

func init() {
	Register("memory", func(cfg config.StorageConfig) (Storage, error) {
		return NewMemoryStorage(), nil
	})
}
