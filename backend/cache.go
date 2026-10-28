package backend

import (
	"sync"
	"time"
)

type streamCacheEntry struct {
	result    *StreamResult
	expiresAt time.Time
}

// StreamCache provides in-memory thread-safe caching with a fixed TTL and bounded capacity.
type StreamCache struct {
	mu         sync.RWMutex
	ttl        time.Duration
	maxEntries int
	entries    map[string]streamCacheEntry
	order      []string
}

// NewStreamCache creates a new bounded StreamCache.
func NewStreamCache(ttl time.Duration) *StreamCache {
	return &StreamCache{
		ttl:        ttl,
		maxEntries: 100, // Caps memory footprint to ~100 items (under 1MB)
		entries:    make(map[string]streamCacheEntry),
		order:      make([]string, 0, 100),
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

// Set saves a StreamResult to the cache with the configured TTL and LRU-style eviction.
func (c *StreamCache) Set(videoID string, result *StreamResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If entry doesn't exist and we're at max capacity, evict oldest
	if _, exists := c.entries[videoID]; !exists && len(c.order) >= c.maxEntries {
		oldest := c.order[0]
		c.order = c.order[1:]
		delete(c.entries, oldest)
	}

	c.entries[videoID] = streamCacheEntry{
		result:    result,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.order = append(c.order, videoID)
}

// Clear removes all cached items.
func (c *StreamCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]streamCacheEntry)
	c.order = make([]string, 0, c.maxEntries)
}

type searchCacheEntry struct {
	response  *SearchResponse
	expiresAt time.Time
}

// SearchCache provides bounded, in-memory caching for search queries to save RPM and prevent IP bans.
type SearchCache struct {
	mu         sync.RWMutex
	ttl        time.Duration
	maxEntries int
	entries    map[string]searchCacheEntry
	order      []string
}

// NewSearchCache creates a new bounded SearchCache.
func NewSearchCache(ttl time.Duration, maxEntries int) *SearchCache {
	if maxEntries <= 0 {
		maxEntries = 60 // Default 60 queries max (prevents memory bloat, < 1.5MB RAM)
	}
	return &SearchCache{
		ttl:        ttl,
		maxEntries: maxEntries,
		entries:    make(map[string]searchCacheEntry),
		order:      make([]string, 0, maxEntries),
	}
}

// Get retrieves a cached SearchResponse if valid and not expired.
func (c *SearchCache) Get(key string) (*SearchResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.response, true
}

// Set stores a SearchResponse in the cache, evicting the oldest when max capacity is reached.
func (c *SearchCache) Set(key string, response *SearchResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.entries[key]; !exists && len(c.order) >= c.maxEntries {
		oldest := c.order[0]
		c.order = c.order[1:]
		delete(c.entries, oldest)
	}

	c.entries[key] = searchCacheEntry{
		response:  response,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.order = append(c.order, key)
}

// Clear empties the search cache.
func (c *SearchCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]searchCacheEntry)
	c.order = make([]string, 0, c.maxEntries)
}

// call represents an in-flight or completed singleflight request.
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// SingleFlightGroup suppresses duplicate in-flight requests for the same key.
type SingleFlightGroup struct {
	mu sync.Mutex
	m  map[string]*call
}

// NewSingleFlightGroup creates a new SingleFlightGroup.
func NewSingleFlightGroup() *SingleFlightGroup {
	return &SingleFlightGroup{
		m: make(map[string]*call),
	}
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a time.
func (g *SingleFlightGroup) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
