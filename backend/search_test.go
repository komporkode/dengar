package backend

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParsePipedSearch(t *testing.T) {
	mockJSON := `{
		"items": [
			{
				"url": "/watch?v=dQw4w9WgXcQ",
				"title": "Never Gonna Give You Up",
				"uploaderName": "Rick Astley",
				"duration": 213,
				"thumbnail": "https://pipedproxy.example.com/vi/dQw4w9WgXcQ/hqdefault.jpg"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockJSON))
	}))
	defer server.Close()

	svc := NewSearchService([]string{server.URL}, nil)
	results, err := svc.Search("rick astley", "music_songs")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	item := results[0]
	if item.ID != "dQw4w9WgXcQ" {
		t.Errorf("Expected ID 'dQw4w9WgXcQ', got '%s'", item.ID)
	}
	if item.Title != "Never Gonna Give You Up" {
		t.Errorf("Expected Title 'Never Gonna Give You Up', got '%s'", item.Title)
	}
	if item.UploaderName != "Rick Astley" {
		t.Errorf("Expected UploaderName 'Rick Astley', got '%s'", item.UploaderName)
	}
	if item.DurationFormatted != "03:33" {
		t.Errorf("Expected DurationFormatted '03:33', got '%s'", item.DurationFormatted)
	}
}

func TestSearchFailover(t *testing.T) {
	// First server fails with 500
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server1.Close()

	// Second server succeeds
	mockJSON := `{
		"items": [
			{
				"url": "/watch?v=second123",
				"title": "Failover Track",
				"uploaderName": "Resilient Band",
				"duration": 180,
				"thumbnail": "https://example.com/thumb.jpg"
			}
		]
	}`
	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockJSON))
	}))
	defer server2.Close()

	svc := NewSearchService([]string{server1.URL, server2.URL}, nil)
	results, err := svc.Search("query", "all")
	if err != nil {
		t.Fatalf("Expected failover to succeed, got error: %v", err)
	}

	if len(results) != 1 || results[0].ID != "second123" {
		t.Errorf("Expected result from second server, got: %+v", results)
	}
}

func TestInnertubeSearch(t *testing.T) {
	mockResponse := `{
		"contents": {
			"twoColumnSearchResultsRenderer": {
				"primaryContents": {
					"sectionListRenderer": {
						"contents": [
							{
								"itemSectionRenderer": {
									"contents": [
										{
											"videoRenderer": {
												"videoId": "innertube123",
												"title": {
													"simpleText": "Direct YouTube Track"
												},
												"ownerText": {
													"runs": [
														{"text": "Official Artist"}
													]
												},
												"lengthText": {
													"simpleText": "03:45"
												},
												"thumbnail": {
													"thumbnails": [
														{"url": "https://i.ytimg.com/vi/innertube123/hqdefault.jpg"}
													]
												}
											}
										}
									]
								}
							}
						]
					}
				}
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	svc := NewSearchService(nil, nil)
	svc.SetInnertubeURL(server.URL)

	results, err := svc.Search("Direct Track", "music_songs")
	if err != nil {
		t.Fatalf("Innertube search failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.ID != "innertube123" {
		t.Errorf("Expected ID 'innertube123', got '%s'", r.ID)
	}
	if r.Title != "Direct YouTube Track" {
		t.Errorf("Expected Title 'Direct YouTube Track', got '%s'", r.Title)
	}
	if r.UploaderName != "Official Artist" {
		t.Errorf("Expected UploaderName 'Official Artist', got '%s'", r.UploaderName)
	}
	if r.Duration != 225 {
		t.Errorf("Expected Duration 225s, got %d", r.Duration)
	}
	if r.DurationFormatted != "03:45" {
		t.Errorf("Expected DurationFormatted '03:45', got '%s'", r.DurationFormatted)
	}
}
