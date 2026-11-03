package backend

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAudioCacheManager_SaveAndGet(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dengerin_cache_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mgr, err := NewAudioCacheManager(tempDir, 10*1024*1024)
	if err != nil {
		t.Fatalf("Failed to create AudioCacheManager: %v", err)
	}

	// Mock upstream HTTP server serving audio data
	mockAudioContent := make([]byte, 128*1024) // 128 KB
	for i := range mockAudioContent {
		mockAudioContent[i] = byte(i % 256)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/mp4")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockAudioContent)
	}))
	defer server.Close()

	// 1. Initial check - not found
	_, _, found := mgr.Get("test_song_1")
	if found {
		t.Errorf("Expected not found initially, got true")
	}

	// 2. Queue background cache
	mgr.QueueBackgroundCache("test_song_1", server.URL, "Test Title", "Test Artist", 180, "mp4")

	// Wait for background download to finish
	var cachedPath string
	var cachedMime string
	for i := 0; i < 50; i++ {
		time.Sleep(50 * time.Millisecond)
		cachedPath, cachedMime, found = mgr.Get("test_song_1")
		if found {
			break
		}
	}

	if !found {
		t.Fatalf("Track test_song_1 should have been cached, but was not found")
	}

	if cachedMime != "audio/mp4" {
		t.Errorf("Expected mimeType audio/mp4, got %s", cachedMime)
	}

	// Verify file size
	info, err := os.Stat(cachedPath)
	if err != nil {
		t.Fatalf("Failed to stat cached file: %v", err)
	}
	if info.Size() != int64(len(mockAudioContent)) {
		t.Errorf("Expected size %d, got %d", len(mockAudioContent), info.Size())
	}

	// 3. Check stats
	stats := mgr.GetStats()
	if stats.TotalFiles != 1 {
		t.Errorf("Expected 1 cached file, got %d", stats.TotalFiles)
	}
	if stats.TotalBytes != int64(len(mockAudioContent)) {
		t.Errorf("Expected %d bytes, got %d", len(mockAudioContent), stats.TotalBytes)
	}
}

func TestAudioCacheManager_LRUEnforce(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dengerin_cache_lru_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set limit to 350 KB
	mgr, err := NewAudioCacheManager(tempDir, 350*1024)
	if err != nil {
		t.Fatalf("Failed to create AudioCacheManager: %v", err)
	}

	// Create 3 fake cache files directly: 150 KB each (Total 450 KB > 350 KB)
	buf150KB := make([]byte, 150*1024)

	for i := 1; i <= 3; i++ {
		name := fmt.Sprintf("song_%d", i)
		cachePath := filepath.Join(tempDir, name+".cache")
		if err := os.WriteFile(cachePath, buf150KB, 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		// Adjust mod time so song_1 is oldest, song_3 newest
		modTime := time.Now().Add(time.Duration(i-5) * time.Hour)
		_ = os.Chtimes(cachePath, modTime, modTime)
	}

	// Enforce limit
	mgr.EnforceLimit()

	// song_1 should be evicted since it's the oldest
	_, _, found1 := mgr.Get("song_1")
	if found1 {
		t.Errorf("Expected song_1 to be evicted by LRU")
	}

	// song_3 should still be present
	_, _, found3 := mgr.Get("song_3")
	if !found3 {
		t.Errorf("Expected song_3 to be kept")
	}

	stats := mgr.GetStats()
	if stats.TotalBytes > 350*1024 {
		t.Errorf("Total bytes %d exceeds maxBytes 350KB", stats.TotalBytes)
	}
}

func TestAudioCacheManager_Clear(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dengerin_cache_clear_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mgr, err := NewAudioCacheManager(tempDir, 10*1024*1024)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Create dummy file
	dummy := filepath.Join(tempDir, "song_x.cache")
	_ = os.WriteFile(dummy, []byte("test audio"), 0644)

	if err := mgr.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	stats := mgr.GetStats()
	if stats.TotalFiles != 0 || stats.TotalBytes != 0 {
		t.Errorf("Expected empty cache after Clear, got: %+v", stats)
	}
}

func TestAudioProxyHandler_WithDiskCache(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dengerin_proxy_cache_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mgr, err := NewAudioCacheManager(tempDir, 10*1024*1024)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Write mock audio file larger than MinValidAudioFileSize
	mockContent := make([]byte, 130*1024)
	copy(mockContent, []byte("audio payload byte stream for range seeking test"))
	cachePath := filepath.Join(tempDir, "cached_vid_123.cache")
	_ = os.WriteFile(cachePath, mockContent, 0644)

	handler := NewAudioProxyHandlerWithCache(mgr)

	// Test request with cache_id
	req := httptest.NewRequest("GET", "/api/audio-proxy?cache_id=cached_vid_123", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", rr.Code)
	}

	body, _ := io.ReadAll(rr.Body)
	if len(body) != len(mockContent) {
		t.Errorf("Expected size %d, got %d", len(mockContent), len(body))
	}

	// Test Range request on cached file (Seek test)
	rangeReq := httptest.NewRequest("GET", "/api/audio-proxy?cache_id=cached_vid_123", nil)
	rangeReq.Header.Set("Range", "bytes=0-12")
	rangeRR := httptest.NewRecorder()

	handler.ServeHTTP(rangeRR, rangeReq)

	if rangeRR.Code != http.StatusPartialContent {
		t.Fatalf("Expected status 206 Partial Content, got %d", rangeRR.Code)
	}

	rangeBody, _ := io.ReadAll(rangeRR.Body)
	expectedSlice := string(mockContent[0:13])
	if string(rangeBody) != expectedSlice {
		t.Errorf("Expected range slice %q, got %q", expectedSlice, string(rangeBody))
	}
}

func TestAudioCacheManager_Purge(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dengerin_purge_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mgr, err := NewAudioCacheManager(tempDir, 10*1024*1024)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	mockContent := make([]byte, 130*1024)
	cachePath := filepath.Join(tempDir, "bad_song.cache")
	metaPath := filepath.Join(tempDir, "bad_song.meta")
	_ = os.WriteFile(cachePath, mockContent, 0644)
	_ = os.WriteFile(metaPath, []byte("{}"), 0644)

	if _, _, found := mgr.Get("bad_song"); !found {
		t.Fatalf("Expected bad_song to be found initially")
	}

	mgr.Purge("bad_song")

	if _, _, found := mgr.Get("bad_song"); found {
		t.Errorf("Expected bad_song to be purged, but was still found")
	}
}

func TestDetectAudioMimeType(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dengerin_mime_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test WebM magic bytes (EBML 1A 45 DF A3)
	webMPath := filepath.Join(tempDir, "test.webm")
	_ = os.WriteFile(webMPath, []byte{0x1A, 0x45, 0xDF, 0xA3, 0x01, 0x02}, 0644)
	if mime := DetectAudioMimeType(webMPath, "video/webm"); mime != "audio/webm" {
		t.Errorf("Expected audio/webm, got %s", mime)
	}

	// Test MP4 ftyp
	mp4Path := filepath.Join(tempDir, "test.mp4")
	_ = os.WriteFile(mp4Path, []byte{0x00, 0x00, 0x00, 0x20, 'f', 't', 'y', 'p', 'M', '4', 'A', ' '}, 0644)
	if mime := DetectAudioMimeType(mp4Path, "video/mp4"); mime != "audio/mp4" {
		t.Errorf("Expected audio/mp4, got %s", mime)
	}

	// Test Fallback sanitization
	if mime := DetectAudioMimeType(filepath.Join(tempDir, "nonexistent"), "video/webm; codecs=opus"); mime != "audio/webm" {
		t.Errorf("Expected sanitized audio/webm, got %s", mime)
	}
}

func TestAudioProxyHandler_CacheNotFound(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dengerin_notfound_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mgr, err := NewAudioCacheManager(tempDir, 10*1024*1024)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	handler := NewAudioProxyHandlerWithCache(mgr)

	req := httptest.NewRequest("GET", "/api/audio-proxy?cache_id=nonexistent_track", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404 Not Found for missing cache entry, got %d", rr.Code)
	}
}

func TestAudioCacheManager_PinnedProtection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dengerin_pinned_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set limit to 350 KB
	mgr, err := NewAudioCacheManager(tempDir, 350*1024)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Pin song_1 (Favorites / Playlist track)
	mgr.SetPinnedFunc(func() map[string]bool {
		return map[string]bool{
			"song_1": true,
		}
	})

	// Create 3 fake cache files: 150 KB each (Total 450 KB > 350 KB)
	buf150KB := make([]byte, 150*1024)
	for i := 1; i <= 3; i++ {
		name := fmt.Sprintf("song_%d", i)
		cachePath := filepath.Join(tempDir, name+".cache")
		if err := os.WriteFile(cachePath, buf150KB, 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		// Adjust mod time so song_1 is oldest, song_2 is middle, song_3 is newest
		modTime := time.Now().Add(time.Duration(i-5) * time.Hour)
		_ = os.Chtimes(cachePath, modTime, modTime)
	}

	// Enforce limit
	mgr.EnforceLimit()

	// song_1 is oldest but PINNED -> must NOT be evicted!
	_, _, found1 := mgr.Get("song_1")
	if !found1 {
		t.Errorf("Expected song_1 to be protected from LRU eviction because it is pinned")
	}

	// song_2 is not pinned -> should be evicted instead of song_1
	_, _, found2 := mgr.Get("song_2")
	if found2 {
		t.Errorf("Expected unpinned song_2 to be evicted by LRU")
	}

	// song_3 is newest -> should be kept
	_, _, found3 := mgr.Get("song_3")
	if !found3 {
		t.Errorf("Expected song_3 to be kept")
	}
}

func TestAudioCacheManager_CommitStream(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dengerin_commit_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mgr, err := NewAudioCacheManager(tempDir, 10*1024*1024)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Set pending metadata
	mgr.SetPendingMeta("commit_track", AudioCacheMeta{
		Title:  "Stream Song",
		Artist: "Stream Artist",
	})

	// Write mock temporary file (> 128 KB)
	tmpPath := filepath.Join(tempDir, "commit_track.tmp")
	data := make([]byte, 150*1024)
	// EBML header for WebM
	data[0] = 0x1A
	data[1] = 0x45
	data[2] = 0xDF
	data[3] = 0xA3
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		t.Fatalf("Failed to write tmp file: %v", err)
	}

	committed := mgr.CommitStream("commit_track", tmpPath, int64(len(data)), "video/webm")
	if !committed {
		t.Fatalf("CommitStream failed")
	}

	// Verify track is in cache
	filePath, mime, found := mgr.Get("commit_track")
	if !found {
		t.Fatalf("Expected commit_track to be found in cache")
	}
	if mime != "audio/webm" {
		t.Errorf("Expected audio/webm, got %s", mime)
	}
	if _, err := os.Stat(filePath); err != nil {
		t.Errorf("Cached file does not exist on disk: %v", err)
	}

	// Verify metadata was saved
	meta := mgr.GetMeta("commit_track")
	if meta.Title != "Stream Song" || meta.Artist != "Stream Artist" {
		t.Errorf("Metadata mismatch: %+v", meta)
	}

	// Verify GetCachedTracks
	cachedList := mgr.GetCachedTracks()
	if len(cachedList) == 0 {
		t.Fatalf("Expected at least 1 cached track in GetCachedTracks")
	}
	if cachedList[0].VideoID != "commit_track" {
		t.Errorf("Expected videoId commit_track, got %s", cachedList[0].VideoID)
	}
}

func TestAudioCacheManager_ThumbnailCache(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dengerin_thumb_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mgr, err := NewAudioCacheManager(tempDir, 10*1024*1024)
	if err != nil {
		t.Fatalf("Failed to create AudioCacheManager: %v", err)
	}

	mockThumb := make([]byte, 1024)
	mockThumb[0] = 0xFF
	mockThumb[1] = 0xD8
	mockThumb[2] = 0xFF

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockThumb)
	}))
	defer server.Close()

	// 1. Initially thumbnail not found
	if _, found := mgr.GetThumbnail("song_123"); found {
		t.Errorf("Expected thumbnail not found initially")
	}

	// 2. Cache thumbnail
	err = mgr.CacheThumbnail("song_123", server.URL)
	if err != nil {
		t.Fatalf("CacheThumbnail failed: %v", err)
	}

	// 3. Verify thumbnail found
	thumbPath, found := mgr.GetThumbnail("song_123")
	if !found {
		t.Fatalf("Expected thumbnail to be found after caching")
	}
	info, err := os.Stat(thumbPath)
	if err != nil || info.Size() != int64(len(mockThumb)) {
		t.Errorf("Thumbnail file size mismatch: %v", err)
	}

	// 4. Test Purge cleans up thumbnail
	mgr.Purge("song_123")
	if _, found := mgr.GetThumbnail("song_123"); found {
		t.Errorf("Expected thumbnail to be removed after Purge")
	}
}

func TestCoverProxyHandler_ServeCached(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dengerin_cover_proxy_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mgr, err := NewAudioCacheManager(tempDir, 10*1024*1024)
	if err != nil {
		t.Fatalf("Failed to create AudioCacheManager: %v", err)
	}

	mockThumb := make([]byte, 512)
	mockThumb[0] = 0xFF
	mockThumb[1] = 0xD8
	mockThumb[2] = 0xFF
	thumbPath := filepath.Join(tempDir, "cached_vid.thumb")
	_ = os.WriteFile(thumbPath, mockThumb, 0644)

	handler := NewCoverProxyHandlerWithCache(mgr)

	// 1. Request cached thumbnail
	req := httptest.NewRequest("GET", "/api/cover-proxy?cache_id=cached_vid", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
	if rr.Header().Get("Content-Type") != "image/jpeg" {
		t.Errorf("Expected Content-Type image/jpeg, got %s", rr.Header().Get("Content-Type"))
	}
	if rr.Body.Len() != len(mockThumb) {
		t.Errorf("Expected body length %d, got %d", len(mockThumb), rr.Body.Len())
	}

	// 2. Request nonexistent thumbnail with invalid upstream URL - should serve fallback SVG
	reqFallback := httptest.NewRequest("GET", "/api/cover-proxy?cache_id=nonexistent&url=http://invalid.local.test/404", nil)
	rrFallback := httptest.NewRecorder()
	handler.ServeHTTP(rrFallback, reqFallback)

	if rrFallback.Code != http.StatusOK {
		t.Errorf("Expected status 200 for fallback, got %d", rrFallback.Code)
	}
	if rrFallback.Header().Get("Content-Type") != "image/svg+xml" {
		t.Errorf("Expected Content-Type image/svg+xml for fallback, got %s", rrFallback.Header().Get("Content-Type"))
	}
	if !strings.Contains(rrFallback.Body.String(), "<svg") {
		t.Errorf("Expected fallback body to contain svg, got: %s", rrFallback.Body.String())
	}
}

func TestAudioCacheManager_ConfigPersistence(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dengerin_cache_cfg_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 1. Initial creation with default 10MB
	mgr1, err := NewAudioCacheManager(tempDir, 10*1024*1024)
	if err != nil {
		t.Fatalf("Failed to create AudioCacheManager: %v", err)
	}
	if mgr1.GetStats().MaxBytes != 10*1024*1024 {
		t.Errorf("Expected initial maxBytes %d, got %d", 10*1024*1024, mgr1.GetStats().MaxBytes)
	}

	// 2. Change limit to 2GB (2048 MB)
	targetLimit := int64(2048) * 1024 * 1024
	mgr1.SetLimit(targetLimit)
	if mgr1.GetStats().MaxBytes != targetLimit {
		t.Errorf("Expected updated maxBytes %d, got %d", targetLimit, mgr1.GetStats().MaxBytes)
	}

	// 3. Re-instantiate AudioCacheManager on the same dir (simulating app restart)
	mgr2, err := NewAudioCacheManager(tempDir, 1024*1024*1024)
	if err != nil {
		t.Fatalf("Failed to re-instantiate AudioCacheManager: %v", err)
	}
	if mgr2.GetStats().MaxBytes != targetLimit {
		t.Errorf("Expected persistent restored maxBytes %d, got %d", targetLimit, mgr2.GetStats().MaxBytes)
	}

	// 4. Test setting to 0 (Unlimited)
	mgr2.SetLimit(0)
	if mgr2.GetStats().MaxBytes != 0 {
		t.Errorf("Expected maxBytes 0, got %d", mgr2.GetStats().MaxBytes)
	}

	// 5. Re-instantiate again and confirm 0 is restored
	mgr3, err := NewAudioCacheManager(tempDir, 1024*1024*1024)
	if err != nil {
		t.Fatalf("Failed to re-instantiate AudioCacheManager for unlimited: %v", err)
	}
	if mgr3.GetStats().MaxBytes != 0 {
		t.Errorf("Expected restored maxBytes 0, got %d", mgr3.GetStats().MaxBytes)
	}
}

func TestAudioCacheManager_SetCacheDir_Migration(t *testing.T) {
	dir1, err := os.MkdirTemp("", "dengerin_cachedir1_*")
	if err != nil {
		t.Fatalf("Failed to create dir1: %v", err)
	}
	defer os.RemoveAll(dir1)

	dir2, err := os.MkdirTemp("", "dengerin_cachedir2_*")
	if err != nil {
		t.Fatalf("Failed to create dir2: %v", err)
	}
	defer os.RemoveAll(dir2)

	mgr, err := NewAudioCacheManager(dir1, 50*1024*1024)
	if err != nil {
		t.Fatalf("Failed to create AudioCacheManager: %v", err)
	}

	if mgr.GetDir() != dir1 {
		t.Errorf("Expected active dir %s, got %s", dir1, mgr.GetDir())
	}

	// Create dummy .cache and .meta files in dir1
	dummyAudio := make([]byte, 150*1024)
	_ = os.WriteFile(filepath.Join(dir1, "song123.cache"), dummyAudio, 0644)
	_ = os.WriteFile(filepath.Join(dir1, "song123.meta"), []byte(`{"videoId":"song123","title":"Song 123"}`), 0644)

	// Move to dir2 with moveFiles = true
	if err := mgr.SetCacheDir(dir2, true); err != nil {
		t.Fatalf("SetCacheDir failed: %v", err)
	}

	if mgr.GetDir() != dir2 {
		t.Errorf("Expected active dir %s, got %s", dir2, mgr.GetDir())
	}

	// Verify files moved to dir2
	if _, err := os.Stat(filepath.Join(dir2, "song123.cache")); err != nil {
		t.Errorf("song123.cache should exist in dir2: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir2, "song123.meta")); err != nil {
		t.Errorf("song123.meta should exist in dir2: %v", err)
	}

	// Verify old files no longer in dir1
	if _, err := os.Stat(filepath.Join(dir1, "song123.cache")); !os.IsNotExist(err) {
		t.Errorf("song123.cache should have been removed from dir1")
	}

	// Stats should reflect moved file
	stats := mgr.GetStats()
	if stats.TotalFiles != 1 {
		t.Errorf("Expected 1 file in stats, got %d", stats.TotalFiles)
	}
	if stats.CacheDir != dir2 {
		t.Errorf("Expected stats.CacheDir %s, got %s", dir2, stats.CacheDir)
	}
}




