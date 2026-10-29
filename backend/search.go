package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SearchService performs searches with instance rotation, caching, singleflight deduplication, and pagination.
type SearchService struct {
	pool         *InstancePool
	client       *http.Client
	innertubeURL string
	cache        *SearchCache
	flight       *SingleFlightGroup
}

// NewSearchService creates a new SearchService instance with bounded search caching.
func NewSearchService(customPiped, customInvidious []string) *SearchService {
	itURL := "https://www.youtube.com/youtubei/v1/search?prettyPrint=false"
	if len(customPiped) > 0 || len(customInvidious) > 0 {
		itURL = ""
	}
	return &SearchService{
		pool: NewInstancePool(customPiped, customInvidious),
		client: &http.Client{
			Timeout: 5000 * time.Millisecond,
		},
		innertubeURL: itURL,
		cache:        NewSearchCache(1*time.Hour, 60), // Caps RAM to max 60 queries (~1MB)
		flight:       NewSingleFlightGroup(),
	}
}

// SetInnertubeURL overrides the Innertube API endpoint (useful for testing or custom proxies).
func (s *SearchService) SetInnertubeURL(url string) {
	s.innertubeURL = url
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

type innertubeSearchResponse struct {
	Contents struct {
		TwoColumnSearchResultsRenderer struct {
			PrimaryContents struct {
				SectionListRenderer struct {
					Contents []struct {
						ItemSectionRenderer *struct {
							Contents []struct {
								VideoRenderer *innertubeVideoRenderer `json:"videoRenderer,omitempty"`
							} `json:"contents"`
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
					Contents []struct {
						VideoRenderer *innertubeVideoRenderer `json:"videoRenderer,omitempty"`
					} `json:"contents"`
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
		// 1. Primary: Direct YouTube Innertube search
		if resp, err := s.searchYouTubeInnertube(q); err == nil && len(resp.Items) > 0 {
			s.cache.Set(cacheKey, resp)
			return resp, nil
		}

		filterParam := "music_songs"
		switch filterType {
		case "videos", "video":
			filterParam = "videos"
		case "all":
			filterParam = "all"
		}

		// 2. Fallback to Piped instances
		for _, inst := range s.pool.GetPipedInstances() {
			results, err := s.searchPiped(inst, q, filterParam)
			if err == nil && len(results) > 0 {
				resp := &SearchResponse{Items: results}
				s.cache.Set(cacheKey, resp)
				return resp, nil
			}
			s.pool.RotatePiped()
		}

		// 3. Fallback to Invidious instances
		for _, inst := range s.pool.GetInvidiousInstances() {
			results, err := s.searchInvidious(inst, q)
			if err == nil && len(results) > 0 {
				resp := &SearchResponse{Items: results}
				s.cache.Set(cacheKey, resp)
				return resp, nil
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
		return s.searchYouTubeInnertubeContinuation(tok)
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
				if it.VideoRenderer == nil || it.VideoRenderer.VideoID == "" {
					continue
				}
				items = append(items, extractVideoItem(it.VideoRenderer))
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
						if it.VideoRenderer == nil || it.VideoRenderer.VideoID == "" {
							continue
						}
						items = append(items, extractVideoItem(it.VideoRenderer))
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
	parts := strings.Split(strings.TrimSpace(s), ":")
	if len(parts) == 0 || parts[0] == "" {
		return 0
	}
	var total int64
	for _, p := range parts {
		val, _ := strconv.ParseInt(p, 10, 64)
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
