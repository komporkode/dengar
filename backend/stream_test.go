package backend

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStreamCache(t *testing.T) {
	cache := NewStreamCache(100 * time.Millisecond)

	sample := &StreamResult{
		VideoID:  "vid123",
		Title:    "Test Song",
		AudioURL: "https://stream.audio/1",
	}

	cache.Set("vid123", sample)

	got, found := cache.Get("vid123")
	if !found || got.AudioURL != "https://stream.audio/1" {
		t.Fatalf("Expected cached item to be found, got: %+v, found: %v", got, found)
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	_, foundAfterExp := cache.Get("vid123")
	if foundAfterExp {
		t.Errorf("Expected expired item to not be found")
	}
}

func TestStreamResolverPiped(t *testing.T) {
	mockPipedStream := `{
		"title": "Bohemian Rhapsody",
		"uploader": "Queen",
		"thumbnailUrl": "https://img.youtube.com/vi/bohemian/hqdefault.jpg",
		"duration": 354,
		"audioStreams": [
			{
				"url": "https://stream.example/audio-low.webm",
				"mimeType": "audio/webm",
				"bitrate": 64000
			},
			{
				"url": "https://stream.example/audio-high.m4a",
				"mimeType": "audio/mp4",
				"bitrate": 160000
			}
		],
		"videoStreams": [
			{
				"url": "https://stream.example/video-muxed-720p.mp4",
				"quality": "720p",
				"mimeType": "video/mp4",
				"videoOnly": false
			},
			{
				"url": "https://stream.example/video-only-1080p.mp4",
				"quality": "1080p",
				"mimeType": "video/mp4",
				"videoOnly": true
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockPipedStream))
	}))
	defer server.Close()

	pool := NewInstancePool([]string{server.URL}, []string{})
	resolver := NewStreamResolver(pool, 30*time.Minute)

	res, err := resolver.Resolve("bohemian", false)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if res.Title != "Bohemian Rhapsody" {
		t.Errorf("Expected title 'Bohemian Rhapsody', got '%s'", res.Title)
	}
	if res.AudioURL != "https://stream.example/audio-high.m4a" {
		t.Errorf("Expected highest bitrate m4a stream, got '%s'", res.AudioURL)
	}
	if res.VideoURL != "https://stream.example/video-muxed-720p.mp4" {
		t.Errorf("Expected muxed 720p video stream, got '%s'", res.VideoURL)
	}
}

func TestStreamResolverInvidious(t *testing.T) {
	mockInvidiousStream := `{
		"title": "Fix You",
		"author": "Coldplay",
		"lengthSeconds": 295,
		"videoThumbnails": [
			{"url": "https://img.youtube.com/vi/fixyou/hqdefault.jpg"}
		],
		"adaptiveFormats": [
			{
				"type": "audio/webm; codecs=\"opus\"",
				"url": "https://googlevideo.com/audio-low.webm",
				"bitrate": "70000"
			},
			{
				"type": "audio/mp4; codecs=\"mp4a.40.2\"",
				"url": "https://googlevideo.com/audio-high.m4a",
				"bitrate": "130000"
			}
		],
		"formatStreams": [
			{
				"quality": "medium",
				"resolution": "720p",
				"url": "https://googlevideo.com/video-720p.mp4"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockInvidiousStream))
	}))
	defer server.Close()

	pool := NewInstancePool([]string{}, []string{server.URL})
	resolver := NewStreamResolver(pool, 30*time.Minute)

	res, err := resolver.Resolve("fixyou", false)
	if err != nil {
		t.Fatalf("Invidious resolve failed: %v", err)
	}

	if res.Title != "Fix You" {
		t.Errorf("Expected title 'Fix You', got '%s'", res.Title)
	}
	if res.Artist != "Coldplay" {
		t.Errorf("Expected artist 'Coldplay', got '%s'", res.Artist)
	}
	if res.AudioURL != "https://googlevideo.com/audio-high.m4a" {
		t.Errorf("Expected highest bitrate audio stream, got '%s'", res.AudioURL)
	}
	if res.VideoURL != "https://googlevideo.com/video-720p.mp4" {
		t.Errorf("Expected 720p format stream, got '%s'", res.VideoURL)
	}
}

func TestStreamResolverCompanionNormalization(t *testing.T) {
	mockInvidiousSABR := `{
		"title": "Mawar Bodas",
		"author": "Fazhri official",
		"lengthSeconds": 300,
		"videoThumbnails": [
			{"url": "https://img.youtube.com/vi/mawar/hqdefault.jpg"}
		],
		"adaptiveFormats": [
			{
				"type": "audio/webm; codecs=\"opus\"",
				"url": "https://yt-comp1.example/companion/sabr-segment?itag=251",
				"bitrate": "160000"
			}
		],
		"formatStreams": [
			{
				"quality": "medium",
				"resolution": "720p",
				"url": "https://yt-comp1.example/companion/sabr-segment?itag=22"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockInvidiousSABR))
	}))
	defer server.Close()

	pool := NewInstancePool([]string{}, []string{server.URL})
	resolver := NewStreamResolver(pool, 30*time.Minute)

	res, err := resolver.Resolve("mawar", false)
	if err != nil {
		t.Fatalf("Invidious companion resolve failed: %v", err)
	}

	expectedAudio := server.URL + "/companion/sabr-segment?itag=251"
	if res.AudioURL != expectedAudio {
		t.Errorf("Expected normalized companion audio URL '%s', got '%s'", expectedAudio, res.AudioURL)
	}
}

func TestStreamResolverQualitySaver(t *testing.T) {
	mockPipedStream := `{
		"title": "Starboy",
		"uploader": "The Weeknd",
		"duration": 230,
		"audioStreams": [
			{
				"url": "https://stream.example/audio-low.webm",
				"mimeType": "audio/webm",
				"bitrate": 64000
			},
			{
				"url": "https://stream.example/audio-high.m4a",
				"mimeType": "audio/mp4",
				"bitrate": 160000
			}
		],
		"videoStreams": []
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockPipedStream))
	}))
	defer server.Close()

	pool := NewInstancePool([]string{server.URL}, []string{})
	resolver := NewStreamResolver(pool, 30*time.Minute)

	// Test saver quality
	resSaver, err := resolver.Resolve("starboy", false, "saver")
	if err != nil {
		t.Fatalf("Resolve with saver failed: %v", err)
	}
	if resSaver.AudioURL != "https://stream.example/audio-low.webm" {
		t.Errorf("Expected low bitrate stream for saver, got %s", resSaver.AudioURL)
	}
	if resSaver.Quality != "saver" {
		t.Errorf("Expected Quality 'saver', got %s", resSaver.Quality)
	}

	// Test high quality
	resHigh, err := resolver.Resolve("starboy", false, "high")
	if err != nil {
		t.Fatalf("Resolve with high failed: %v", err)
	}
	if resHigh.AudioURL != "https://stream.example/audio-high.m4a" {
		t.Errorf("Expected high bitrate stream for high, got %s", resHigh.AudioURL)
	}
	if resHigh.Quality != "high" {
		t.Errorf("Expected Quality 'high', got %s", resHigh.Quality)
	}
}

func TestSmartServerMode(t *testing.T) {
	invidiousServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invidious 500 error", http.StatusInternalServerError)
	}))
	defer invidiousServer.Close()

	pipedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"title": "Piped Track",
			"uploader": "Artist",
			"audioStreams": [{"url": "https://stream.example/audio.m4a", "bitrate": 160000, "mimeType": "audio/mp4"}]
		}`))
	}))
	defer pipedServer.Close()

	pool := NewInstancePool([]string{pipedServer.URL}, []string{invidiousServer.URL})
	resolver := NewStreamResolver(pool, 1*time.Minute)

	if resolver.GetPrimaryProvider() != "invidious" {
		t.Fatalf("Expected initial primary provider to be 'invidious', got '%s'", resolver.GetPrimaryProvider())
	}

	// Track 1: Invidious fails, Piped succeeds -> fails=1, primary remains invidious
	_, err := resolver.Resolve("track1", true)
	if err != nil {
		t.Fatalf("Resolve track1 failed: %v", err)
	}
	if resolver.GetPrimaryProvider() != "invidious" {
		t.Errorf("Expected primary to remain 'invidious' after 1 fail, got '%s'", resolver.GetPrimaryProvider())
	}

	// Retry Track 1: should NOT increase failure counter
	_, _ = resolver.Resolve("track1", true)
	if resolver.GetPrimaryProvider() != "invidious" {
		t.Errorf("Expected primary to remain 'invidious' after repeating same track, got '%s'", resolver.GetPrimaryProvider())
	}

	// Track 2: Invidious fails, Piped succeeds -> fails=2, primary remains invidious
	_, err = resolver.Resolve("track2", true)
	if err != nil {
		t.Fatalf("Resolve track2 failed: %v", err)
	}
	if resolver.GetPrimaryProvider() != "invidious" {
		t.Errorf("Expected primary to remain 'invidious' after 2 fails, got '%s'", resolver.GetPrimaryProvider())
	}

	// Track 3: Invidious fails, Piped succeeds -> fails=3 -> PROMOTED to piped!
	_, err = resolver.Resolve("track3", true)
	if err != nil {
		t.Fatalf("Resolve track3 failed: %v", err)
	}
	if resolver.GetPrimaryProvider() != "piped" {
		t.Errorf("Expected primary to be promoted to 'piped' after 3 consecutive song failovers, got '%s'", resolver.GetPrimaryProvider())
	}

	// Track 4: Now Piped is primary, resolves directly
	res4, err := resolver.Resolve("track4", true)
	if err != nil {
		t.Fatalf("Resolve track4 failed: %v", err)
	}
	if res4.Title != "Piped Track" {
		t.Errorf("Expected track from Piped, got title: '%s'", res4.Title)
	}
}

func TestSmartServerSuccessReset(t *testing.T) {
	invidiousFail := true
	invidiousServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if invidiousFail {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"title": "Invidious Track",
			"author": "Author",
			"adaptiveFormats": [{"type": "audio/mp4", "url": "https://stream.example/invidious.m4a", "bitrate": "130000"}]
		}`))
	}))
	defer invidiousServer.Close()

	pipedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"title": "Piped Track",
			"uploader": "Artist",
			"audioStreams": [{"url": "https://stream.example/audio.m4a", "bitrate": 160000, "mimeType": "audio/mp4"}]
		}`))
	}))
	defer pipedServer.Close()

	pool := NewInstancePool([]string{pipedServer.URL}, []string{invidiousServer.URL})
	resolver := NewStreamResolver(pool, 1*time.Minute)

	// Song 1 fails on Invidious, succeeds on Piped (1/3)
	invidiousFail = true
	_, err := resolver.Resolve("song1", true)
	if err != nil {
		t.Fatalf("Resolve song1 failed: %v", err)
	}

	// Song 2 succeeds on Invidious -> resets consecutive fails to 0
	invidiousFail = false
	_, err = resolver.Resolve("song2", true)
	if err != nil {
		t.Fatalf("Resolve song2 failed: %v", err)
	}

	// Song 3 fails on Invidious, succeeds on Piped (1/3 again)
	invidiousFail = true
	_, err = resolver.Resolve("song3", true)
	if err != nil {
		t.Fatalf("Resolve song3 failed: %v", err)
	}

	// Song 4 fails on Invidious, succeeds on Piped (2/3)
	_, err = resolver.Resolve("song4", true)
	if err != nil {
		t.Fatalf("Resolve song4 failed: %v", err)
	}

	// Provider must STILL be invidious (not promoted yet because counter was reset by song2)
	if resolver.GetPrimaryProvider() != "invidious" {
		t.Errorf("Expected primary provider to still be 'invidious', got '%s'", resolver.GetPrimaryProvider())
	}
}
