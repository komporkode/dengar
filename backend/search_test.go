package backend

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
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

	if len(results.Items) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results.Items))
	}

	item := results.Items[0]
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

	if len(results.Items) != 1 || results.Items[0].ID != "second123" {
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
							},
							{
								"continuationItemRenderer": {
									"continuationEndpoint": {
										"continuationCommand": {
											"token": "token_page_2"
										}
									}
								}
							}
						]
					}
				}
			}
		}
	}`

	var hitCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hitCount, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	svc := NewSearchService(nil, nil)
	svc.SetYTMURL("")
	svc.SetInnertubeURL(server.URL)

	results, err := svc.Search("Direct Track", "music_songs")
	if err != nil {
		t.Fatalf("Innertube search failed: %v", err)
	}

	if len(results.Items) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results.Items))
	}

	if results.ContinuationToken != "token_page_2" {
		t.Errorf("Expected ContinuationToken 'token_page_2', got '%s'", results.ContinuationToken)
	}

	r := results.Items[0]
	if r.ID != "innertube123" {
		t.Errorf("Expected ID 'innertube123', got '%s'", r.ID)
	}
	if r.Title != "Direct YouTube Track" {
		t.Errorf("Expected Title 'Direct YouTube Track', got '%s'", r.Title)
	}

	// Test Cache: Repeated search should NOT hit the HTTP server again!
	cachedResults, err := svc.Search("Direct Track", "music_songs")
	if err != nil {
		t.Fatalf("Cached search failed: %v", err)
	}
	if len(cachedResults.Items) != 1 {
		t.Errorf("Expected cached items, got %d", len(cachedResults.Items))
	}
	if atomic.LoadInt32(&hitCount) != 1 {
		t.Errorf("Expected hitCount to be 1 (served from cache), got %d", atomic.LoadInt32(&hitCount))
	}
}

func TestInnertubeContinuation(t *testing.T) {
	mockContinuationResponse := `{
		"onResponseReceivedCommands": [
			{
				"appendContinuationItemsAction": {
					"continuationItems": [
						{
							"itemSectionRenderer": {
								"contents": [
									{
										"videoRenderer": {
											"videoId": "page2_track",
											"title": {"simpleText": "Page 2 Song"},
											"ownerText": {"simpleText": "Artist 2"},
											"lengthText": {"simpleText": "04:20"},
											"thumbnail": {"thumbnails": [{"url": "https://img.test/p2.jpg"}]}
										}
									}
								]
							}
						},
						{
							"continuationItemRenderer": {
								"continuationEndpoint": {
									"continuationCommand": {
										"token": "token_page_3"
									}
								}
							}
						}
					]
				}
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockContinuationResponse))
	}))
	defer server.Close()

	svc := NewSearchService(nil, nil)
	svc.SetYTMURL("")
	svc.SetInnertubeURL(server.URL)

	res, err := svc.SearchMore("token_page_2")
	if err != nil {
		t.Fatalf("SearchMore failed: %v", err)
	}

	if len(res.Items) != 1 {
		t.Fatalf("Expected 1 continuation item, got %d", len(res.Items))
	}
	if res.Items[0].ID != "page2_track" {
		t.Errorf("Expected ID 'page2_track', got '%s'", res.Items[0].ID)
	}
	if res.ContinuationToken != "token_page_3" {
		t.Errorf("Expected next token 'token_page_3', got '%s'", res.ContinuationToken)
	}
}

func TestParseDurationString(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"03:45", 225},
		{"3.45", 225},
		{"04:20", 260},
		{"4.20", 260},
		{"1:02:15", 3735},
		{"1.02.15", 3735},
		{"50", 50},
		{"", 0},
		{"  ", 0},
	}

	for _, tt := range tests {
		got := parseDurationString(tt.input)
		if got != tt.expected {
			t.Errorf("parseDurationString(%q) = %d; want %d", tt.input, got, tt.expected)
		}
	}
}

func TestIsMusicTrack(t *testing.T) {
	cases := []struct {
		title    string
		uploader string
		dur      int64
		query    string
		expected bool
	}{
		{"Sheila On 7 - Dan (Official Audio)", "Sony Music Indonesia", 280, "Sheila On 7", true},
		{"Avenged Sevenfold - Dear God (Official Video)", "Warner Records", 245, "Avenged Sevenfold", true},
		{"GUE KORBAN OM!! DITIPU KELAS HABIS2AN", "Deddy Corbuzier", 4920, "Deddy Corbuzier", false},
		{"Tutorial Cara Install Windows 11 Cepat", "Tech Indo", 600, "Windows", false},
		{"Gameplay GTA 5 Mod Indonesia", "Gamer Pro", 1200, "GTA 5", false},
		{"Vlog Liburan ke Bali Bareng Pacar", "Daily Vlogger", 850, "Bali", false},
		{"Reaksi Bule Mendengarkan Denny Caknan", "Bule Reaction", 450, "Denny Caknan", false},
		{"Alur Cerita Film Horor Terseram", "Kamar Nonton", 950, "Film Horor", false},
		{"Meme Pendek Lucu", "Shorts Channel", 15, "Meme", false}, // < 30s
		{"BEST OF SHEILA ON 7 - 1 JAM NONSTOP", "SETH.", 3600, "Sheila On 7", true}, // music compilation video is kept
		{"NINA - KAMANA CINTANA, BAGJA JEUNG CINTA, JAYANTI - POP SUNDA GASENTRA PAJAMPANGAN TERPOPULER 2025", "DangDing Sunda", 0, "Nina Gasentra", true}, // music compilation is kept
		{"NINA GASENTRA - KUMPULAN LAGU POP SUNDA TERBARU", "GASENTRA PAJAMPANGAN", 0, "Nina Gasentra", true}, // music compilation is kept
		{"AVENGED SEVENFOLD BEST OF ROCK SONG COLLECTION 2025", "Rock Music", 0, "Avenged Sevenfold", true}, // music compilation is kept
		{"Video Random Non Musik Dua Jam", "Random Channel", 7200, "Sheila On 7", false}, // > 20m and not compilation
		{"Deddy Corbuzier Podcast Theme Song", "Deddy Corbuzier", 210, "Deddy Corbuzier podcast", true}, // user explicitly asked for podcast keyword
	}

	for _, c := range cases {
		got := IsMusicTrack(c.title, c.uploader, c.dur, c.query)
		if got != c.expected {
			t.Errorf("IsMusicTrack(title=%q, dur=%d, query=%q) = %v; want %v", c.title, c.dur, c.query, got, c.expected)
		}
	}
}

func TestParseYTMResponse(t *testing.T) {
	mockYTM := `{
		"contents": {
			"tabbedSearchResultsRenderer": {
				"tabs": [
					{
						"tabRenderer": {
							"content": {
								"sectionListRenderer": {
									"contents": [
										{
											"musicShelfRenderer": {
												"contents": [
													{
														"musicResponsiveListItemRenderer": {
															"flexColumns": [
																{
																	"musicResponsiveListItemFlexColumnRenderer": {
																		"text": {
																			"runs": [
																				{
																					"text": "Sephia",
																					"navigationEndpoint": {
																						"watchEndpoint": {
																							"videoId": "ytm_sephia_123"
																						}
																					}
																				}
																			]
																		}
																	}
																},
																{
																	"musicResponsiveListItemFlexColumnRenderer": {
																		"text": {
																			"runs": [
																				{"text": "Lagu"},
																				{"text": " • "},
																				{"text": "Sheila On 7"},
																				{"text": " • "},
																				{"text": "4.56"},
																				{"text": " • "},
																				{"text": "10 jt ditonton"}
																			]
																		}
																	}
																}
															],
															"thumbnail": {
																"musicThumbnailRenderer": {
																	"thumbnail": {
																		"thumbnails": [
																			{"url": "https://img.ytm/sephia.jpg"}
																		]
																	}
																}
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
				]
			}
		}
	}`

	items := parseYTMResponse([]byte(mockYTM), "Sheila On 7")
	if len(items) != 1 {
		t.Fatalf("Expected 1 parsed YTM item, got %d", len(items))
	}

	it := items[0]
	if it.ID != "ytm_sephia_123" {
		t.Errorf("Expected ID 'ytm_sephia_123', got '%s'", it.ID)
	}
	if it.Title != "Sephia" {
		t.Errorf("Expected Title 'Sephia', got '%s'", it.Title)
	}
	if it.UploaderName != "Sheila On 7" {
		t.Errorf("Expected UploaderName 'Sheila On 7', got '%s'", it.UploaderName)
	}
	if it.Duration != 296 {
		t.Errorf("Expected Duration 296s (4:56), got %d", it.Duration)
	}
	if it.ThumbnailURL != "https://img.ytm/sephia.jpg" {
		t.Errorf("Expected ThumbnailURL 'https://img.ytm/sephia.jpg', got '%s'", it.ThumbnailURL)
	}
}

func TestYouTubeMusicSearch(t *testing.T) {
	mockYTM := `{
		"contents": {
			"tabbedSearchResultsRenderer": {
				"tabs": [
					{
						"tabRenderer": {
							"content": {
								"sectionListRenderer": {
									"contents": [
										{
											"musicShelfRenderer": {
												"contents": [
													{
														"musicResponsiveListItemRenderer": {
															"flexColumns": [
																{
																	"musicResponsiveListItemFlexColumnRenderer": {
																		"text": {
																			"runs": [
																				{
																					"text": "Song 1",
																					"navigationEndpoint": {"watchEndpoint": {"videoId": "s1"}}
																				}
																			]
																		}
																	}
																},
																{
																	"musicResponsiveListItemFlexColumnRenderer": {
																		"text": {
																			"runs": [
																				{"text": "Artist A"},
																				{"text": " • "},
																				{"text": "3.30"}
																			]
																		}
																	}
																}
															]
														}
													},
													{
														"musicResponsiveListItemRenderer": {
															"flexColumns": [
																{
																	"musicResponsiveListItemFlexColumnRenderer": {
																		"text": {
																			"runs": [
																				{
																					"text": "Song 2",
																					"navigationEndpoint": {"watchEndpoint": {"videoId": "s2"}}
																				}
																			]
																		}
																	}
																},
																{
																	"musicResponsiveListItemFlexColumnRenderer": {
																		"text": {
																			"runs": [
																				{"text": "Artist A"},
																				{"text": " • "},
																				{"text": "4.00"}
																			]
																		}
																	}
																}
															]
														}
													},
													{
														"musicResponsiveListItemRenderer": {
															"flexColumns": [
																{
																	"musicResponsiveListItemFlexColumnRenderer": {
																		"text": {
																			"runs": [
																				{
																					"text": "Song 3",
																					"navigationEndpoint": {"watchEndpoint": {"videoId": "s3"}}
																				}
																			]
																		}
																	}
																},
																{
																	"musicResponsiveListItemFlexColumnRenderer": {
																		"text": {
																			"runs": [
																				{"text": "Artist A"},
																				{"text": " • "},
																				{"text": "3.45"}
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
									]
								}
							}
						}
					}
				]
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockYTM))
	}))
	defer server.Close()

	svc := NewSearchService(nil, nil)
	svc.SetInnertubeURL("")
	svc.SetYTMURL(server.URL)

	res, err := svc.Search("Artist A", "music_songs")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(res.Items) != 3 {
		t.Fatalf("Expected 3 items from YTM, got %d", len(res.Items))
	}
	if res.Items[0].Title != "Song 1" || res.Items[1].Title != "Song 2" || res.Items[2].Title != "Song 3" {
		t.Errorf("Unexpected items: %+v", res.Items)
	}
}
