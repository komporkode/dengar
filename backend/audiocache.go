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
	dir         string
	maxBytes    int64
	mu          sync.RWMutex
	inFlight    sync.Map
	pendingMeta sync.Map
	pinnedFunc  func() map[string]bool
	client      *http.Client
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

// MinValidAudioFileSize is the minimum threshold (128 KB) for a file to be considered a valid, playable song.
const MinValidAudioFileSize int64 = 128 * 1024

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

	mgr := &AudioCacheManager{
		dir:      dir,
		maxBytes: maxBytes,
		client: &http.Client{
			Timeout: 45 * time.Second,
		},
	}

	// Self-healing: sweep and purge any leftover .tmp or corrupted < 128KB files from older versions
	go mgr.sanitizeCacheDir()

	return mgr, nil
}

// sanitizeVideoID prevents path traversal attacks.
func sanitizeVideoID(id string) string {
	clean := strings.TrimSpace(id)
	clean = strings.ReplaceAll(clean, "/", "_")
	clean = strings.ReplaceAll(clean, "\\", "_")
	clean = strings.ReplaceAll(clean, "..", "_")
	return clean
}

// DetectAudioMimeType sniffs magic bytes from an audio file to determine the correct audio/* MIME type.
func DetectAudioMimeType(filePath string, fallback string) string {
	cleanFallback := strings.TrimSpace(fallback)
	if strings.HasPrefix(cleanFallback, "video/webm") {
		cleanFallback = "audio/webm"
	} else if strings.HasPrefix(cleanFallback, "video/mp4") {
		cleanFallback = "audio/mp4"
	} else if !strings.HasPrefix(cleanFallback, "audio/") {
		cleanFallback = "audio/mp4"
	}

	f, err := os.Open(filePath)
	if err != nil {
		return cleanFallback
	}
	defer f.Close()

	var header [16]byte
	n, _ := io.ReadFull(f, header[:])
	if n >= 4 {
		// WebM / Matroska EBML: \x1A\x45\xDF\xA3
		if header[0] == 0x1A && header[1] == 0x45 && header[2] == 0xDF && header[3] == 0xA3 {
			return "audio/webm"
		}
		// OGG: "OggS"
		if string(header[0:4]) == "OggS" {
			return "audio/ogg"
		}
		// FLAC: "fLaC"
		if string(header[0:4]) == "fLaC" {
			return "audio/flac"
		}
		// MP3: "ID3"
		if header[0] == 'I' && header[1] == 'D' && header[2] == '3' {
			return "audio/mpeg"
		}
		// MP3: sync word (0xFF 0xFB, 0xFF 0xF3, 0xFF 0xF2, 0xFF 0xE3)
		if header[0] == 0xFF && (header[1]&0xE0) == 0xE0 {
			return "audio/mpeg"
		}
	}
	if n >= 12 {
		// MP4 / M4A: bytes 4..8 == "ftyp"
		if string(header[4:8]) == "ftyp" {
			return "audio/mp4"
		}
	}

	return cleanFallback
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
	if err != nil || info.Size() < MinValidAudioFileSize {
		// Corrupted or incomplete file from an aborted session - clean up immediately
		if err == nil && info.Size() < MinValidAudioFileSize {
			_ = os.Remove(cachePath)
			_ = os.Remove(filepath.Join(m.dir, vid+".meta"))
		}
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

	// Sniff real audio MIME from magic bytes (converts video/* to audio/*)
	detectedMime := DetectAudioMimeType(cachePath, mime)

	// Update touch time for LRU tracking only if last touched > 10 minutes ago
	// to avoid resetting ModTime during rapid Range requests by the player
	if time.Since(info.ModTime()) > 10*time.Minute {
		now := time.Now()
		_ = os.Chtimes(cachePath, now, now)
	}

	return cachePath, detectedMime, true
}

// SetPinnedFunc sets a provider function returning a map of pinned video IDs (Favorites & Playlists) immune to LRU eviction.
func (m *AudioCacheManager) SetPinnedFunc(fn func() map[string]bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pinnedFunc = fn
}

// SetPendingMeta temporarily caches track metadata so stream-through proxy can tag the cached file without extra requests.
func (m *AudioCacheManager) SetPendingMeta(videoID string, meta AudioCacheMeta) {
	vid := sanitizeVideoID(videoID)
	if vid != "" {
		m.pendingMeta.Store(vid, meta)
	}
}

// CommitStream atomically moves a completed .tmp file to .cache, saves metadata, and enforces LRU limits.
func (m *AudioCacheManager) CommitStream(videoID string, tmpPath string, written int64, mimeType string) bool {
	vid := sanitizeVideoID(videoID)
	if vid == "" || written < MinValidAudioFileSize {
		_ = os.Remove(tmpPath)
		return false
	}

	cachePath := filepath.Join(m.dir, vid+".cache")
	metaPath := filepath.Join(m.dir, vid+".meta")
	realMime := DetectAudioMimeType(tmpPath, mimeType)

	m.mu.Lock()
	defer m.mu.Unlock()

	if err := os.Rename(tmpPath, cachePath); err != nil {
		_ = os.Remove(tmpPath)
		return false
	}

	meta := AudioCacheMeta{
		VideoID:        vid,
		MimeType:       realMime,
		FileSize:       written,
		CachedAt:       time.Now().Unix(),
		LastAccessedAt: time.Now().Unix(),
	}
	if val, ok := m.pendingMeta.Load(vid); ok {
		if pm, ok := val.(AudioCacheMeta); ok {
			meta.Title = pm.Title
			meta.Artist = pm.Artist
			meta.Duration = pm.Duration
			meta.Format = pm.Format
		}
	}

	if metaJSON, err := json.Marshal(meta); err == nil {
		_ = os.WriteFile(metaPath, metaJSON, 0644)
	}

	m.enforceLimitLocked()
	return true
}

// GetMeta returns metadata for a cached track if available.
func (m *AudioCacheManager) GetMeta(videoID string) AudioCacheMeta {
	vid := sanitizeVideoID(videoID)
	if vid == "" {
		return AudioCacheMeta{}
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	metaPath := filepath.Join(m.dir, vid+".meta")
	var meta AudioCacheMeta
	if metaBytes, err := os.ReadFile(metaPath); err == nil {
		_ = json.Unmarshal(metaBytes, &meta)
	}
	if meta.VideoID == "" {
		if val, ok := m.pendingMeta.Load(vid); ok {
			if pm, ok := val.(AudioCacheMeta); ok {
				return pm
			}
		}
	}
	return meta
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
	if strings.HasPrefix(mimeType, "video/webm") {
		mimeType = "audio/webm"
	} else if strings.HasPrefix(mimeType, "video/mp4") {
		mimeType = "audio/mp4"
	} else if !strings.HasPrefix(mimeType, "audio/") {
		if strings.Contains(strings.ToLower(format), "webm") || strings.Contains(strings.ToLower(format), "opus") {
			mimeType = "audio/webm"
		} else {
			mimeType = "audio/mp4"
		}
	}

	// 3. Stream data to temp file
	written, err := io.Copy(tmpFile, resp.Body)
	if err != nil || written < MinValidAudioFileSize {
		return // Discard incomplete or truncated downloads
	}
	// Verify against Content-Length if provided by server
	if resp.ContentLength > 0 && written < resp.ContentLength {
		return // Truncated download, discard
	}
	_ = tmpFile.Close()

	// Sniff magic bytes from temp file to ensure real audio MIME
	realMime := DetectAudioMimeType(tmpPath, mimeType)

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
		MimeType:       realMime,
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

	var pinned map[string]bool
	if m.pinnedFunc != nil {
		pinned = m.pinnedFunc()
	}

	for _, f := range cacheFiles {
		if totalSize <= m.maxBytes {
			break
		}
		vid := strings.TrimSuffix(f.name, ".cache")
		if pinned != nil && pinned[vid] {
			continue // Pinned tracks (Favorites & Playlists) are immune to LRU eviction
		}
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

// sanitizeCacheDir cleans up orphaned .tmp files and corrupted .cache files smaller than MinValidAudioFileSize.
func (m *AudioCacheManager) sanitizeCacheDir() {
	m.mu.Lock()
	defer m.mu.Unlock()

	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		fullPath := filepath.Join(m.dir, name)

		if strings.HasSuffix(name, ".tmp") {
			_ = os.Remove(fullPath)
			continue
		}

		if strings.HasSuffix(name, ".cache") {
			info, err := entry.Info()
			if err != nil || info.Size() < MinValidAudioFileSize {
				_ = os.Remove(fullPath)
				vid := strings.TrimSuffix(name, ".cache")
				_ = os.Remove(filepath.Join(m.dir, vid+".meta"))
			}
		}
	}
}

// Purge removes cached files for a specific track (e.g. if playback reported media error).
func (m *AudioCacheManager) Purge(videoID string) {
	vid := sanitizeVideoID(videoID)
	if vid == "" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	_ = os.Remove(filepath.Join(m.dir, vid+".cache"))
	_ = os.Remove(filepath.Join(m.dir, vid+".meta"))
	_ = os.Remove(filepath.Join(m.dir, vid+".tmp"))
}
