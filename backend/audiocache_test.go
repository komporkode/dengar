package backend

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

	// Set limit to 200 KB
	mgr, err := NewAudioCacheManager(tempDir, 200*1024)
	if err != nil {
		t.Fatalf("Failed to create AudioCacheManager: %v", err)
	}

	// Create 3 fake cache files directly: 100 KB each (Total 300 KB > 200 KB)
	buf100KB := make([]byte, 100*1024)

	for i := 1; i <= 3; i++ {
		name := fmt.Sprintf("song_%d", i)
		cachePath := filepath.Join(tempDir, name+".cache")
		if err := os.WriteFile(cachePath, buf100KB, 0644); err != nil {
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
	if stats.TotalBytes > 200*1024 {
		t.Errorf("Total bytes %d exceeds maxBytes 200KB", stats.TotalBytes)
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

	// Write mock audio file
	mockContent := []byte("audio payload byte stream for range seeking test")
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
	if string(body) != string(mockContent) {
		t.Errorf("Expected content %s, got %s", string(mockContent), string(body))
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
