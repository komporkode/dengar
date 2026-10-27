package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SearchService performs searches with instance rotation and fallback.
type SearchService struct {
	pool   *InstancePool
	client *http.Client
}

// NewSearchService creates a new SearchService instance.
func NewSearchService(customPiped, customInvidious []string) *SearchService {
	return &SearchService{
		pool: NewInstancePool(customPiped, customInvidious),
		client: &http.Client{
			Timeout: 6 * time.Second,
		},
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

// Search searches across Piped instances with Invidious fallback.
func (s *SearchService) Search(query string, filterType string) ([]SearchResultItem, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return []SearchResultItem{}, nil
	}

	filterParam := "music_songs"
	switch filterType {
	case "videos", "video":
		filterParam = "videos"
	case "all":
		filterParam = "all"
	}

	// 1. Try Piped instances
	for _, inst := range s.pool.GetPipedInstances() {
		results, err := s.searchPiped(inst, q, filterParam)
		if err == nil && len(results) > 0 {
			return results, nil
		}
		s.pool.RotatePiped()
	}

	// 2. Fallback to Invidious instances
	for _, inst := range s.pool.GetInvidiousInstances() {
		results, err := s.searchInvidious(inst, q)
		if err == nil && len(results) > 0 {
			return results, nil
		}
		s.pool.RotateInvidious()
	}

	return nil, fmt.Errorf("semua instance pencarian gagal merespons")
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
