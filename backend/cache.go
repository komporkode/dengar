package backend

import (
	"sync"
	"time"
)

type cacheEntry struct {
	result    *StreamResult
	expiresAt time.Time
}

// StreamCache provides in-memory thread-safe caching with a fixed TTL.
type StreamCache struct {
	mu      sync.RWMutex
	ttl     time.Duration
	entries map[string]cacheEntry
}

// NewStreamCache creates a new StreamCache.
func NewStreamCache(ttl time.Duration) *StreamCache {
	return &StreamCache{
		ttl:     ttl,
		entries: make(map[string]cacheEntry),
	}
}

// Get retrieves a cached StreamResult if it exists and has not expired.
func (c *StreamCache) Get(videoID string) (*StreamResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[videoID]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.result, true
}

// Set saves a StreamResult to the cache with the configured TTL.
func (c *StreamCache) Set(videoID string, result *StreamResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[videoID] = cacheEntry{
		result:    result,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Clear removes all cached items.
func (c *StreamCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]cacheEntry)
}
