package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// AudioCacheMeta stores metadata for a cached audio file.
type AudioCacheMeta struct {
	VideoID        string `json:"videoId"`
	Title          string `json:"title"`
	Artist         string `json:"artist"`
	Duration       int64  `json:"duration"`
	Format         string `json:"format"`
	MimeType       string `json:"mimeType"`
	FileSize       int64  `json:"fileSize"`
	CachedAt       int64  `json:"cachedAt"`
	LastAccessedAt int64  `json:"lastAccessedAt"`
}

// AudioCacheStats holds usage information for disk audio caching.
type AudioCacheStats struct {
	TotalBytes int64  `json:"totalBytes"`
	TotalFiles int    `json:"totalFiles"`
	MaxBytes   int64  `json:"maxBytes"`
	CacheDir   string `json:"cacheDir"`
}

// AudioCacheManager handles caching full audio files on local disk for offline and instant replay.
type AudioCacheManager struct {
	dir      string
	maxBytes int64
	mu       sync.RWMutex
	inFlight sync.Map
	client   *http.Client
}

// DefaultCacheDir returns the standard cache directory under "dengerin" folder.
func DefaultCacheDir() string {
	baseDir, err := os.UserCacheDir()
	if err != nil {
		baseDir, err = os.UserConfigDir()
		if err != nil {
			baseDir = "."
		}
	}
	return filepath.Join(baseDir, "dengerin", "audiocache")
}

// NewAudioCacheManager creates a new AudioCacheManager.
// If customDir is empty, DefaultCacheDir() is used.
// If maxBytes <= 0, default 1 GB (1024*1024*1024) is used.
func NewAudioCacheManager(customDir string, maxBytes int64) (*AudioCacheManager, error) {
	dir := customDir
	if dir == "" {
		dir = DefaultCacheDir()
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("gagal membuat folder cache audio di %s: %w", dir, err)
	}

	if maxBytes <= 0 {
		maxBytes = 1024 * 1024 * 1024 // 1 GB default
	}

	return &AudioCacheManager{
		dir:      dir,
		maxBytes: maxBytes,
		client: &http.Client{
			Timeout: 45 * time.Second,
		},
	}, nil
}

// sanitizeVideoID prevents path traversal attacks.
func sanitizeVideoID(id string) string {
	clean := strings.TrimSpace(id)
	clean = strings.ReplaceAll(clean, "/", "_")
	clean = strings.ReplaceAll(clean, "\\", "_")
	clean = strings.ReplaceAll(clean, "..", "_")
	return clean
}

// Get checks if a track is cached on local disk and ready for instant/offline playback.
func (m *AudioCacheManager) Get(videoID string) (filePath string, mimeType string, found bool) {
	vid := sanitizeVideoID(videoID)
	if vid == "" {
		return "", "", false
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	cachePath := filepath.Join(m.dir, vid+".cache")
	info, err := os.Stat(cachePath)
	if err != nil || info.Size() == 0 {
		return "", "", false
	}

	// Read metadata if available
	metaPath := filepath.Join(m.dir, vid+".meta")
	mime := "audio/mp4"
	if metaBytes, err := os.ReadFile(metaPath); err == nil {
		var meta AudioCacheMeta
		if json.Unmarshal(metaBytes, &meta) == nil && meta.MimeType != "" {
			mime = meta.MimeType
		}
	}

	// Update touch time for LRU tracking
	now := time.Now()
	_ = os.Chtimes(cachePath, now, now)

	return cachePath, mime, true
}

// QueueBackgroundCache initiates an atomic background download for a stream if not already cached.
func (m *AudioCacheManager) QueueBackgroundCache(videoID string, streamURL string, title string, artist string, duration int64, format string) {
	vid := sanitizeVideoID(videoID)
	if vid == "" || streamURL == "" {
		return
	}

	// Don't cache if already local proxy or relative URL
	if strings.HasPrefix(streamURL, "/api/") || strings.Contains(streamURL, "cache_id=") {
		return
	}

	// Skip if already in cache
	if _, _, ok := m.Get(vid); ok {
		return
	}

	// Deduplicate in-flight requests
	if _, loaded := m.inFlight.LoadOrStore(vid, true); loaded {
		return
	}

	go func() {
		defer m.inFlight.Delete(vid)
		m.downloadAndCommit(vid, streamURL, title, artist, duration, format)
	}()
}

func (m *AudioCacheManager) downloadAndCommit(vid string, streamURL string, title string, artist string, duration int64, format string) {
	tmpPath := filepath.Join(m.dir, vid+".tmp")
	cachePath := filepath.Join(m.dir, vid+".cache")
	metaPath := filepath.Join(m.dir, vid+".meta")

	// 1. Prepare temp file
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return
	}
	defer func() {
		tmpFile.Close()
		_ = os.Remove(tmpPath) // Clean up temp file on failure
	}()

	// 2. Fetch stream from upstream
	req, err := http.NewRequest("GET", streamURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := m.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		if strings.Contains(strings.ToLower(format), "webm") || strings.Contains(strings.ToLower(format), "opus") {
			mimeType = "audio/webm"
		} else {
			mimeType = "audio/mp4"
		}
	}

	// 3. Stream data to temp file
	written, err := io.Copy(tmpFile, resp.Body)
	if err != nil || written < 64*1024 {
		return // Discard incomplete or truncated downloads
	}
	_ = tmpFile.Close()

	// 4. Atomic rename temp file to cache file
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := os.Rename(tmpPath, cachePath); err != nil {
		return
	}

	// 5. Write metadata
	meta := AudioCacheMeta{
		VideoID:        vid,
		Title:          title,
		Artist:         artist,
		Duration:       duration,
		Format:         format,
		MimeType:       mimeType,
		FileSize:       written,
		CachedAt:       time.Now().Unix(),
		LastAccessedAt: time.Now().Unix(),
	}
	if metaJSON, err := json.Marshal(meta); err == nil {
		_ = os.WriteFile(metaPath, metaJSON, 0644)
	}

	// 6. Enforce storage limit (LRU eviction)
	m.enforceLimitLocked()
}

// EnforceLimit cleans up oldest accessed files when total size exceeds maxBytes.
func (m *AudioCacheManager) EnforceLimit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enforceLimitLocked()
}

func (m *AudioCacheManager) enforceLimitLocked() {
	if m.maxBytes <= 0 {
		return // Unlimited
	}

	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return
	}

	type fileItem struct {
		name    string
		size    int64
		modTime time.Time
	}

	var cacheFiles []fileItem
	var totalSize int64

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".cache") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		cacheFiles = append(cacheFiles, fileItem{
			name:    entry.Name(),
			size:    info.Size(),
			modTime: info.ModTime(),
		})
		totalSize += info.Size()
	}

	if totalSize <= m.maxBytes {
		return
	}

	// Sort oldest accessed first
	sort.Slice(cacheFiles, func(i, j int) bool {
		return cacheFiles[i].modTime.Before(cacheFiles[j].modTime)
	})

	for _, f := range cacheFiles {
		if totalSize <= m.maxBytes {
			break
		}
		vid := strings.TrimSuffix(f.name, ".cache")
		_ = os.Remove(filepath.Join(m.dir, f.name))
		_ = os.Remove(filepath.Join(m.dir, vid+".meta"))
		totalSize -= f.size
	}
}

// GetStats returns usage statistics for the cache directory.
func (m *AudioCacheManager) GetStats() AudioCacheStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := AudioCacheStats{
		MaxBytes: m.maxBytes,
		CacheDir: m.dir,
	}

	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return stats
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".cache") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		stats.TotalBytes += info.Size()
		stats.TotalFiles++
	}

	return stats
}

// Clear removes all cached audio files and metadata.
func (m *AudioCacheManager) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".cache") || strings.HasSuffix(name, ".meta") || strings.HasSuffix(name, ".tmp") {
			_ = os.Remove(filepath.Join(m.dir, name))
		}
	}
	return nil
}

// SetLimit updates the maximum allowed storage bytes and triggers eviction if needed.
func (m *AudioCacheManager) SetLimit(maxBytes int64) {
	m.mu.Lock()
	m.maxBytes = maxBytes
	m.enforceLimitLocked()
	m.mu.Unlock()
}
