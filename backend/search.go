package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// SearchService performs searches with instance rotation, caching, singleflight deduplication, and pagination.
type SearchService struct {
	pool         *InstancePool
	client       *http.Client
	innertubeURL string
	ytmURL       string
	cache        *SearchCache
	flight       *SingleFlightGroup
}

// NewSearchService creates a new SearchService instance with bounded search caching.
func NewSearchService(customPiped, customInvidious []string) *SearchService {
	itURL := "https://www.youtube.com/youtubei/v1/search?prettyPrint=false"
	ytmURL := "https://music.youtube.com/youtubei/v1/search?prettyPrint=false"
	if len(customPiped) > 0 || len(customInvidious) > 0 {
		itURL = ""
		ytmURL = ""
	}
	return &SearchService{
		pool: NewInstancePool(customPiped, customInvidious),
		client: &http.Client{
			Timeout: 3500 * time.Millisecond,
		},
		innertubeURL: itURL,
		ytmURL:       ytmURL,
		cache:        NewSearchCache(1*time.Hour, 60), // Caps RAM to max 60 queries (~1MB)
		flight:       NewSingleFlightGroup(),
	}
}

// SetInnertubeURL overrides the Innertube API endpoint (useful for testing or custom proxies).
func (s *SearchService) SetInnertubeURL(url string) {
	s.innertubeURL = url
}

// SetYTMURL overrides the YouTube Music API endpoint (useful for testing or custom proxies).
func (s *SearchService) SetYTMURL(url string) {
	s.ytmURL = url
}

// ResetInstances forces failover rotation of primary instances and refreshes live instance pool.
func (s *SearchService) ResetInstances() {
	if s.pool != nil {
		s.pool.RotatePiped()
		s.pool.RotateInvidious()
		s.pool.RefreshInvidiousInstances()
	}
}

// pipedItem represents item structure from Piped search response.
type pipedItem struct {
	URL          string `json:"url"`
	Title        string `json:"title"`
	UploaderName string `json:"uploaderName"`
	Duration     int64  `json:"duration"`
	Thumbnail    string `json:"thumbnail"`
}

type pipedSearchResponse struct {
	Items []pipedItem `json:"items"`
}

// invidiousItem represents item structure from Invidious search response.
type invidiousItem struct {
	Type            string `json:"type"`
	VideoID         string `json:"videoId"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	LengthSeconds   int64  `json:"lengthSeconds"`
	VideoThumbnails []struct {
		URL string `json:"url"`
	} `json:"videoThumbnails"`
}

type innertubeVideoRenderer struct {
	VideoID string `json:"videoId"`
	Title   struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
		SimpleText string `json:"simpleText"`
	} `json:"title"`
	OwnerText struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
		SimpleText string `json:"simpleText"`
	} `json:"ownerText"`
	LengthText struct {
		SimpleText string `json:"simpleText"`
		Runs       []struct {
			Text string `json:"text"`
		} `json:"runs"`
		Accessibility struct {
			AccessibilityData struct {
				Label string `json:"label"`
			} `json:"accessibilityData"`
		} `json:"accessibility"`
	} `json:"lengthText"`
	Thumbnail struct {
		Thumbnails []struct {
			URL string `json:"url"`
		} `json:"thumbnails"`
	} `json:"thumbnail"`
}

type innertubeContinuationItem struct {
	ContinuationEndpoint struct {
		ContinuationCommand struct {
			Token string `json:"token"`
		} `json:"continuationCommand"`
	} `json:"continuationEndpoint"`
}

type innertubeSearchItem struct {
	VideoRenderer *innertubeVideoRenderer `json:"videoRenderer,omitempty"`
	ShelfRenderer *struct {
		Content *struct {
			VerticalListRenderer *struct {
				Items []struct {
					VideoRenderer *innertubeVideoRenderer `json:"videoRenderer,omitempty"`
				} `json:"items"`
			} `json:"verticalListRenderer"`
		} `json:"content"`
	} `json:"shelfRenderer,omitempty"`
}

type innertubeSearchResponse struct {
	Contents struct {
		TwoColumnSearchResultsRenderer struct {
			PrimaryContents struct {
				SectionListRenderer struct {
					Contents []struct {
						ItemSectionRenderer *struct {
							Contents []innertubeSearchItem `json:"contents"`
						} `json:"itemSectionRenderer,omitempty"`
						ContinuationItemRenderer *innertubeContinuationItem `json:"continuationItemRenderer,omitempty"`
					} `json:"contents"`
				} `json:"sectionListRenderer"`
			} `json:"primaryContents"`
		} `json:"twoColumnSearchResultsRenderer"`
	} `json:"contents"`
}

type innertubeContinuationResponse struct {
	OnResponseReceivedCommands []struct {
		AppendContinuationItemsAction *struct {
			ContinuationItems []struct {
				ItemSectionRenderer *struct {
					Contents []innertubeSearchItem `json:"contents"`
				} `json:"itemSectionRenderer,omitempty"`
				ContinuationItemRenderer *innertubeContinuationItem `json:"continuationItemRenderer,omitempty"`
			} `json:"continuationItems"`
		} `json:"appendContinuationItemsAction,omitempty"`
	} `json:"onResponseReceivedCommands"`
}

// Search searches across direct YouTube Innertube first with deduplication and caching, with Piped/Invidious fallback.
func (s *SearchService) Search(query string, filterType string) (*SearchResponse, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return &SearchResponse{Items: []SearchResultItem{}}, nil
	}

	cacheKey := strings.ToLower(q) + ":" + filterType
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached, nil
	}

	// Use SingleFlight to deduplicate simultaneous duplicate requests
	val, err := s.flight.Do(cacheKey, func() (interface{}, error) {
		isMusicSearch := filterType != "videos" && filterType != "video"

		// 1. Primary: Direct YouTube Innertube with music filtering & pagination support
		if resp, err := s.searchYouTubeInnertube(q); err == nil && len(resp.Items) > 0 {
			var filteredItems []SearchResultItem
			for _, it := range resp.Items {
				if !isMusicSearch || IsMusicTrack(it.Title, it.UploaderName, it.Duration, q) {
					filteredItems = append(filteredItems, it)
				}
			}
			if len(filteredItems) > 0 {
				resp.Items = filteredItems
				s.cache.Set(cacheKey, resp)
				return resp, nil
			}
		}

		// 2. Secondary fallback: Direct YouTube Music search (100% music guarantee)
		if isMusicSearch {
			if resp, err := s.searchYouTubeMusic(q); err == nil && len(resp.Items) > 0 {
				s.cache.Set(cacheKey, resp)
				return resp, nil
			}
		}

		filterParam := "music_songs"
		switch filterType {
		case "videos", "video":
			filterParam = "videos"
		case "all":
			filterParam = "all"
		}

		// 3. Fallback to Piped instances
		for _, inst := range s.pool.GetPipedInstances() {
			results, err := s.searchPiped(inst, q, filterParam)
			if err == nil && len(results) > 0 {
				var filteredItems []SearchResultItem
				for _, it := range results {
					if !isMusicSearch || IsMusicTrack(it.Title, it.UploaderName, it.Duration, q) {
						filteredItems = append(filteredItems, it)
					}
				}
				if len(filteredItems) > 0 {
					resp := &SearchResponse{Items: filteredItems}
					s.cache.Set(cacheKey, resp)
					return resp, nil
				}
			}
			s.pool.RotatePiped()
		}

		// 4. Fallback to Invidious instances
		for _, inst := range s.pool.GetInvidiousInstances() {
			results, err := s.searchInvidious(inst, q)
			if err == nil && len(results) > 0 {
				var filteredItems []SearchResultItem
				for _, it := range results {
					if !isMusicSearch || IsMusicTrack(it.Title, it.UploaderName, it.Duration, q) {
						filteredItems = append(filteredItems, it)
					}
				}
				if len(filteredItems) > 0 {
					resp := &SearchResponse{Items: filteredItems}
					s.cache.Set(cacheKey, resp)
					return resp, nil
				}
			}
			s.pool.RotateInvidious()
		}

		return nil, fmt.Errorf("semua instance pencarian gagal merespons")
	})

	if err != nil {
		return nil, err
	}
	return val.(*SearchResponse), nil
}

// SearchMore loads the next page of search results using a continuation token.
func (s *SearchService) SearchMore(continuationToken string) (*SearchResponse, error) {
	tok := strings.TrimSpace(continuationToken)
	if tok == "" {
		return &SearchResponse{Items: []SearchResultItem{}}, nil
	}

	val, err := s.flight.Do("cont:"+tok, func() (interface{}, error) {
		resp, err := s.searchYouTubeInnertubeContinuation(tok)
		if err != nil {
			return nil, err
		}
		var filtered []SearchResultItem
		for _, it := range resp.Items {
			if IsMusicTrack(it.Title, it.UploaderName, it.Duration, "") {
				filtered = append(filtered, it)
			}
		}
		resp.Items = filtered
		return resp, nil
	})

	if err != nil {
		return nil, err
	}
	return val.(*SearchResponse), nil
}

func (s *SearchService) searchYouTubeInnertube(query string) (*SearchResponse, error) {
	if s.innertubeURL == "" {
		return nil, fmt.Errorf("innertube disabled")
	}

	payload := map[string]interface{}{
		"context": map[string]interface{}{
			"client": map[string]interface{}{
				"clientName":    "WEB",
				"clientVersion": "2.20240101.01.00",
				"hl":            "id",
				"gl":            "ID",
			},
		},
		"query": query,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, s.innertubeURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("innertube status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data innertubeSearchResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var items []SearchResultItem
	var continuationToken string

	sections := data.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents
	for _, sec := range sections {
		if sec.ContinuationItemRenderer != nil {
			if tok := sec.ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token; tok != "" {
				continuationToken = tok
			}
		}
		if sec.ItemSectionRenderer != nil {
			for _, it := range sec.ItemSectionRenderer.Contents {
				if it.VideoRenderer != nil && it.VideoRenderer.VideoID != "" {
					items = append(items, extractVideoItem(it.VideoRenderer))
				}
				if it.ShelfRenderer != nil && it.ShelfRenderer.Content != nil && it.ShelfRenderer.Content.VerticalListRenderer != nil {
					for _, sub := range it.ShelfRenderer.Content.VerticalListRenderer.Items {
						if sub.VideoRenderer != nil && sub.VideoRenderer.VideoID != "" {
							items = append(items, extractVideoItem(sub.VideoRenderer))
						}
					}
				}
			}
		}
	}

	return &SearchResponse{
		Items:             items,
		ContinuationToken: continuationToken,
	}, nil
}

func (s *SearchService) searchYouTubeInnertubeContinuation(continuationToken string) (*SearchResponse, error) {
	if s.innertubeURL == "" {
		return nil, fmt.Errorf("innertube disabled")
	}

	payload := map[string]interface{}{
		"context": map[string]interface{}{
			"client": map[string]interface{}{
				"clientName":    "WEB",
				"clientVersion": "2.20240101.01.00",
				"hl":            "id",
				"gl":            "ID",
			},
		},
		"continuation": continuationToken,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, s.innertubeURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("innertube continuation status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data innertubeContinuationResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var items []SearchResultItem
	var nextToken string

	for _, cmd := range data.OnResponseReceivedCommands {
		if cmd.AppendContinuationItemsAction != nil {
			for _, ci := range cmd.AppendContinuationItemsAction.ContinuationItems {
				if ci.ContinuationItemRenderer != nil {
					if tok := ci.ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token; tok != "" {
						nextToken = tok
					}
				}
				if ci.ItemSectionRenderer != nil {
					for _, it := range ci.ItemSectionRenderer.Contents {
						if it.VideoRenderer != nil && it.VideoRenderer.VideoID != "" {
							items = append(items, extractVideoItem(it.VideoRenderer))
						}
						if it.ShelfRenderer != nil && it.ShelfRenderer.Content != nil && it.ShelfRenderer.Content.VerticalListRenderer != nil {
							for _, sub := range it.ShelfRenderer.Content.VerticalListRenderer.Items {
								if sub.VideoRenderer != nil && sub.VideoRenderer.VideoID != "" {
									items = append(items, extractVideoItem(sub.VideoRenderer))
								}
							}
						}
					}
				}
			}
		}
	}

	return &SearchResponse{
		Items:             items,
		ContinuationToken: nextToken,
	}, nil
}

func extractVideoItem(vr *innertubeVideoRenderer) SearchResultItem {
	title := vr.Title.SimpleText
	if title == "" && len(vr.Title.Runs) > 0 {
		var sb strings.Builder
		for _, r := range vr.Title.Runs {
			sb.WriteString(r.Text)
		}
		title = sb.String()
	}

	uploader := vr.OwnerText.SimpleText
	if uploader == "" && len(vr.OwnerText.Runs) > 0 {
		var sb strings.Builder
		for _, r := range vr.OwnerText.Runs {
			sb.WriteString(r.Text)
		}
		uploader = sb.String()
	}

	rawDuration := vr.LengthText.SimpleText
	if rawDuration == "" && len(vr.LengthText.Runs) > 0 {
		var sb strings.Builder
		for _, r := range vr.LengthText.Runs {
			sb.WriteString(r.Text)
		}
		rawDuration = sb.String()
	}
	durationSec := parseDurationString(rawDuration)
	thumb := ""
	if len(vr.Thumbnail.Thumbnails) > 0 {
		thumb = vr.Thumbnail.Thumbnails[len(vr.Thumbnail.Thumbnails)-1].URL
	}

	return SearchResultItem{
		ID:                vr.VideoID,
		Title:             title,
		UploaderName:      uploader,
		Duration:          durationSec,
		DurationFormatted: FormatDuration(durationSec),
		ThumbnailURL:      thumb,
	}
}

func parseDurationString(s string) int64 {
	clean := strings.TrimSpace(s)
	clean = strings.ReplaceAll(clean, ".", ":")
	parts := strings.Split(clean, ":")
	if len(parts) == 0 || parts[0] == "" {
		return 0
	}
	var total int64
	for _, p := range parts {
		var val int64
		for _, c := range p {
			if c >= '0' && c <= '9' {
				val = val*10 + int64(c-'0')
			}
		}
		total = total*60 + val
	}
	return total
}

func (s *SearchService) searchPiped(instanceURL, query, filter string) ([]SearchResultItem, error) {
	apiURL := fmt.Sprintf("%s/search?q=%s&filter=%s",
		strings.TrimRight(instanceURL, "/"),
		url.QueryEscape(query),
		url.QueryEscape(filter),
	)

	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("piped status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data pipedSearchResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	items := make([]SearchResultItem, 0, len(data.Items))
	for _, it := range data.Items {
		vidID := extractVideoID(it.URL)
		if vidID == "" {
			continue
		}
		items = append(items, SearchResultItem{
			ID:                vidID,
			Title:             it.Title,
			UploaderName:      it.UploaderName,
			Duration:          it.Duration,
			DurationFormatted: FormatDuration(it.Duration),
			ThumbnailURL:      it.Thumbnail,
		})
	}
	return items, nil
}

func (s *SearchService) searchInvidious(instanceURL, query string) ([]SearchResultItem, error) {
	apiURL := fmt.Sprintf("%s/api/v1/search?q=%s&type=video",
		strings.TrimRight(instanceURL, "/"),
		url.QueryEscape(query),
	)

	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invidious status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rawItems []invidiousItem
	if err := json.Unmarshal(body, &rawItems); err != nil {
		return nil, err
	}

	items := make([]SearchResultItem, 0, len(rawItems))
	for _, it := range rawItems {
		if it.VideoID == "" {
			continue
		}
		thumb := ""
		if len(it.VideoThumbnails) > 0 {
			thumb = it.VideoThumbnails[0].URL
		}
		items = append(items, SearchResultItem{
			ID:                it.VideoID,
			Title:             it.Title,
			UploaderName:      it.Author,
			Duration:          it.LengthSeconds,
			DurationFormatted: FormatDuration(it.LengthSeconds),
			ThumbnailURL:      thumb,
		})
	}
	return items, nil
}

// extractVideoID extracts YouTube ID from paths like "/watch?v=ID" or returns ID directly.
func extractVideoID(rawURL string) string {
	if strings.HasPrefix(rawURL, "/watch?v=") {
		return strings.TrimPrefix(rawURL, "/watch?v=")
	}
	if u, err := url.Parse(rawURL); err == nil && u.Query().Get("v") != "" {
		return u.Query().Get("v")
	}
	return strings.TrimPrefix(rawURL, "/")
}

// FormatDuration formats seconds into MM:SS or HH:MM:SS.
func FormatDuration(sec int64) string {
	if sec <= 0 {
		return "00:00"
	}
	h := sec / 3600
	m := (sec % 3600) / 60
	s := sec % 60
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

var durationRegex = regexp.MustCompile(`^\d+[\.:]\d+`)
var compilationRegex = regexp.MustCompile(`(?i)(kumpulan lagu|kompilasi|compilation|best of|greatest hits|full album|nonstop|non stop|non-stop|\b\d+\s*(?:jam|hours?)\b|album lengkap|semua lagu|pilihan terbaik|lagu pilihan|lagu lawas pilihan|terpopuler\s*20\d\d)`)

// searchYouTubeMusic queries YouTube Music (WEB_REMIX) for 100% music content.
func (s *SearchService) searchYouTubeMusic(query string) (*SearchResponse, error) {
	if s.ytmURL == "" {
		return nil, fmt.Errorf("ytm disabled")
	}

	payload := map[string]interface{}{
		"context": map[string]interface{}{
			"client": map[string]interface{}{
				"clientName":    "WEB_REMIX",
				"clientVersion": "1.20240101.01.00",
				"hl":            "id",
				"gl":            "ID",
			},
		},
		"query": query,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, s.ytmURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://music.youtube.com/")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ytm status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	items := parseYTMResponse(body, query)
	if len(items) == 0 {
		return nil, fmt.Errorf("no items parsed from youtube music")
	}

	return &SearchResponse{
		Items: items,
	}, nil
}

func parseYTMResponse(body []byte, query string) []SearchResultItem {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil
	}

	contents, ok := data["contents"].(map[string]interface{})
	if !ok {
		return nil
	}
	tsrr, ok := contents["tabbedSearchResultsRenderer"].(map[string]interface{})
	if !ok {
		return nil
	}
	tabs, ok := tsrr["tabs"].([]interface{})
	if !ok || len(tabs) == 0 {
		return nil
	}
	tab0, ok := tabs[0].(map[string]interface{})["tabRenderer"].(map[string]interface{})
	if !ok {
		return nil
	}
	tabContent, ok := tab0["content"].(map[string]interface{})
	if !ok {
		return nil
	}
	sectionList, ok := tabContent["sectionListRenderer"].(map[string]interface{})["contents"].([]interface{})
	if !ok {
		return nil
	}

	var results []SearchResultItem
	seen := make(map[string]bool)

	for _, s := range sectionList {
		sMap, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		for _, v := range sMap {
			vMap, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			cList, ok := vMap["contents"].([]interface{})
			if !ok {
				continue
			}
			for _, item := range cList {
				iMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				mrli, ok := iMap["musicResponsiveListItemRenderer"].(map[string]interface{})
				if !ok {
					continue
				}

				flexCols, ok := mrli["flexColumns"].([]interface{})
				if !ok || len(flexCols) == 0 {
					continue
				}

				// 1. Title
				col0, ok := flexCols[0].(map[string]interface{})["musicResponsiveListItemFlexColumnRenderer"].(map[string]interface{})
				if !ok {
					continue
				}
				text0, ok := col0["text"].(map[string]interface{})
				if !ok {
					continue
				}
				runs0, ok := text0["runs"].([]interface{})
				if !ok || len(runs0) == 0 {
					continue
				}
				title, _ := runs0[0].(map[string]interface{})["text"].(string)
				if title == "" {
					continue
				}

				// 2. Video ID
				var videoID string
				if nav, ok := runs0[0].(map[string]interface{})["navigationEndpoint"].(map[string]interface{}); ok {
					if we, ok := nav["watchEndpoint"].(map[string]interface{}); ok {
						if vid, ok := we["videoId"].(string); ok {
							videoID = vid
						}
					}
				}
				if videoID == "" {
					if pnav, ok := mrli["overlay"].(map[string]interface{}); ok {
						if mibo, ok := pnav["musicItemThumbnailOverlayRenderer"].(map[string]interface{}); ok {
							if cb, ok := mibo["content"].(map[string]interface{}); ok {
								if mpb, ok := cb["musicPlayButtonRenderer"].(map[string]interface{}); ok {
									if pe, ok := mpb["playNavigationEndpoint"].(map[string]interface{}); ok {
										if we, ok := pe["watchEndpoint"].(map[string]interface{}); ok {
											if vid, ok := we["videoId"].(string); ok {
												videoID = vid
											}
										}
									}
								}
							}
						}
					}
				}

				if videoID == "" || seen[videoID] {
					continue
				}

				// Discard non-track sections (Playlist, Album, Episode, Podcast, Artis, Profil)
				var itemType string
				if len(flexCols) > 1 {
					if col1, ok := flexCols[1].(map[string]interface{})["musicResponsiveListItemFlexColumnRenderer"].(map[string]interface{}); ok {
						if text1, ok := col1["text"].(map[string]interface{}); ok {
							if runs1, ok := text1["runs"].([]interface{}); ok && len(runs1) > 0 {
								if r0, ok := runs1[0].(map[string]interface{}); ok {
									itemType, _ = r0["text"].(string)
								}
							}
						}
					}
				}

				lowerType := strings.ToLower(itemType)
				if lowerType == "playlist" || lowerType == "daftar putar" ||
					lowerType == "album" || lowerType == "episode" ||
					lowerType == "podcast" || lowerType == "artis" ||
					lowerType == "artist" || lowerType == "profil" || lowerType == "profile" {
					continue
				}

				// 3. Artist & Duration from flexCols[1]
				var artist string
				var durationSec int64
				if len(flexCols) > 1 {
					col1, ok := flexCols[1].(map[string]interface{})["musicResponsiveListItemFlexColumnRenderer"].(map[string]interface{})
					if ok {
						if text1, ok := col1["text"].(map[string]interface{}); ok {
							if runs1, ok := text1["runs"].([]interface{}); ok {
								var artistParts []string
								for _, r := range runs1 {
									rMap, ok := r.(map[string]interface{})
									if !ok {
										continue
									}
									t, _ := rMap["text"].(string)
									t = strings.TrimSpace(t)
									if t == "•" || t == "" {
										continue
									}
									lowerT := strings.ToLower(t)
									if lowerT == "lagu" || lowerT == "song" || lowerT == "video" || lowerT == "playlist" || lowerT == "album" {
										continue
									}
									if durationRegex.MatchString(t) {
										durationSec = parseDurationString(t)
										continue
									}
									if strings.HasSuffix(lowerT, "ditonton") || strings.HasSuffix(lowerT, "views") {
										continue
									}
									artistParts = append(artistParts, t)
								}
								if len(artistParts) > 0 {
									artist = artistParts[0]
								}
							}
						}
					}
				}

				if artist == "" {
					artist = query
				}

				// 4. Thumbnail
				var thumb string
				if thumbRenderer, ok := mrli["thumbnail"].(map[string]interface{}); ok {
					if mtr, ok := thumbRenderer["musicThumbnailRenderer"].(map[string]interface{}); ok {
						if tObj, ok := mtr["thumbnail"].(map[string]interface{}); ok {
							if thumbs, ok := tObj["thumbnails"].([]interface{}); ok && len(thumbs) > 0 {
								last, _ := thumbs[len(thumbs)-1].(map[string]interface{})
								if last != nil {
									thumb, _ = last["url"].(string)
								}
							}
						}
					}
				}

				// Filter out non-music if any managed to slip in
				if !IsMusicTrack(title, artist, durationSec, query) {
					continue
				}

				seen[videoID] = true
				results = append(results, SearchResultItem{
					ID:                videoID,
					Title:             title,
					UploaderName:      artist,
					Duration:          durationSec,
					DurationFormatted: FormatDuration(durationSec),
					ThumbnailURL:      thumb,
				})
			}
		}
	}
	return results
}

// IsMusicTrack checks whether an item is likely a music track based on keywords and duration heuristics.
func IsMusicTrack(title, uploader string, durationSec int64, query string) bool {
	tLower := strings.ToLower(title)
	uLower := strings.ToLower(uploader)
	qLower := strings.ToLower(query)

	// Non-music keywords blacklist
	blacklist := []string{
		"podcast", "talkshow", "interview", "wawancara", "klarifikasi",
		"gameplay", "walkthrough", "playthrough", "let's play",
		"tutorial", "cara membuat", "cara install", "cara mudah", "unboxing", "review hp",
		"vlog", "prank", "mukbang", "reaction", "reaksi",
		"breaking news", "berita terkini", "live stream", "highlight liga", "full match",
		"alur cerita", "rekap film", "sinopsis film", "trailer resmi", "official trailer",
		"full movie", "film pendek", "short movie", "alur film", "spoiler", "bedah film",
	}

	for _, kw := range blacklist {
		// If query explicitly asked for this keyword, allow it
		if strings.Contains(qLower, kw) {
			continue
		}
		if strings.Contains(tLower, kw) || strings.Contains(uLower, kw) {
			return false
		}
	}

	// Identify whether the item is a music compilation or long-form music mix
	isMusicCompilation := compilationRegex.MatchString(title) ||
		strings.Count(title, ",") >= 2 ||
		strings.Contains(qLower, "album") ||
		strings.Contains(qLower, "kompilasi") ||
		strings.Contains(qLower, "compilation") ||
		strings.Contains(qLower, "kumpulan") ||
		strings.Contains(qLower, "best of") ||
		strings.Contains(qLower, "discography") ||
		strings.Contains(qLower, "playlist") ||
		strings.Contains(qLower, "mix") ||
		strings.Contains(qLower, "full") ||
		strings.Contains(qLower, "dj")

	// Duration filter:
	// Memes/shorts under 30s are filtered out.
	// Long videos (> 20 mins) are allowed if they are music compilations, albums, or mixes,
	// but filtered out if they are non-compilation video uploads.
	if durationSec > 0 {
		if durationSec < 30 {
			return false
		}
		if durationSec > 1200 && !isMusicCompilation {
			return false
		}
	}

	return true
}
