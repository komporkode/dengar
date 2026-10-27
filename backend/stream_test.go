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

	pool := NewInstancePool([]string{server.URL}, nil)
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

	pool := NewInstancePool(nil, []string{server.URL})
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
