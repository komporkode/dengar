package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// StreamResolver resolves streaming audio and video URLs with failover and caching.
type StreamResolver struct {
	pool   *InstancePool
	cache  *StreamCache
	client *http.Client
}

// NewStreamResolver creates a new StreamResolver.
func NewStreamResolver(pool *InstancePool, ttl time.Duration) *StreamResolver {
	return &StreamResolver{
		pool:  pool,
		cache: NewStreamCache(ttl),
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

// Resolve returns StreamResult for videoID. If forceRefresh is false, cached value is preferred.
func (r *StreamResolver) Resolve(videoID string, forceRefresh bool) (*StreamResult, error) {
	vid := strings.TrimSpace(videoID)
	if vid == "" {
		return nil, fmt.Errorf("video ID kosong")
	}

	if !forceRefresh {
		if cached, ok := r.cache.Get(vid); ok {
			return cached, nil
		}
	}

	// 1. Try Piped instances
	for _, inst := range r.pool.GetPipedInstances() {
		res, err := r.resolvePiped(inst, vid)
		if err == nil && res.AudioURL != "" {
			r.cache.Set(vid, res)
			return res, nil
		}
		r.pool.RotatePiped()
	}

	// 2. Fallback to Invidious instances
	for _, inst := range r.pool.GetInvidiousInstances() {
		res, err := r.resolveInvidious(inst, vid)
		if err == nil && res.AudioURL != "" {
			r.cache.Set(vid, res)
			return res, nil
		}
		r.pool.RotateInvidious()
	}

	return nil, fmt.Errorf("gagal memutar media dari seluruh instance")
}

type pipedStreamAudio struct {
	URL      string `json:"url"`
	MimeType string `json:"mimeType"`
	Bitrate  int64  `json:"bitrate"`
}

type pipedStreamVideo struct {
	URL       string `json:"url"`
	Quality   string `json:"quality"`
	MimeType  string `json:"mimeType"`
	VideoOnly bool   `json:"videoOnly"`
}

type pipedStreamResponse struct {
	Title        string             `json:"title"`
	Uploader     string             `json:"uploader"`
	ThumbnailURL string             `json:"thumbnailUrl"`
	Duration     int64              `json:"duration"`
	AudioStreams []pipedStreamAudio `json:"audioStreams"`
	VideoStreams []pipedStreamVideo `json:"videoStreams"`
}

func (r *StreamResolver) resolvePiped(instanceURL, videoID string) (*StreamResult, error) {
	apiURL := fmt.Sprintf("%s/streams/%s", strings.TrimRight(instanceURL, "/"), videoID)

	resp, err := r.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("piped stream status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data pipedStreamResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	// Pick best audio: prefer m4a/mp4 highest bitrate, fallback to webm
	bestAudioURL := ""
	var maxBitrate int64 = -1
	for _, a := range data.AudioStreams {
		if strings.Contains(a.MimeType, "mp4") || strings.Contains(a.MimeType, "m4a") {
			if a.Bitrate > maxBitrate {
				maxBitrate = a.Bitrate
				bestAudioURL = a.URL
			}
		}
	}
	if bestAudioURL == "" {
		for _, a := range data.AudioStreams {
			if a.Bitrate > maxBitrate {
				maxBitrate = a.Bitrate
				bestAudioURL = a.URL
			}
		}
	}

	// Pick muxed video (videoOnly == false)
	bestVideoURL := ""
	for _, v := range data.VideoStreams {
		if !v.VideoOnly && strings.Contains(v.MimeType, "mp4") {
			bestVideoURL = v.URL
			if strings.Contains(v.Quality, "720") {
				break
			}
		}
	}
	// Fallback to any non-empty video URL if muxed not found
	if bestVideoURL == "" && len(data.VideoStreams) > 0 {
		for _, v := range data.VideoStreams {
			if !v.VideoOnly {
				bestVideoURL = v.URL
				break
			}
		}
	}

	return &StreamResult{
		VideoID:   videoID,
		Title:     data.Title,
		Artist:    data.Uploader,
		Thumbnail: data.ThumbnailURL,
		Duration:  data.Duration,
		AudioURL:  bestAudioURL,
		VideoURL:  bestVideoURL,
	}, nil
}

type invidiousAdaptiveFormat struct {
	URL      string `json:"url"`
	Type     string `json:"type"`
	Bitrate  int64  `json:"bitrate,string"`
	Encoding string `json:"encoding"`
}

type invidiousFormatStream struct {
	URL        string `json:"url"`
	Resolution string `json:"resolution"`
	Quality    string `json:"quality"`
}

type invidiousVideoResponse struct {
	Title           string                    `json:"title"`
	Author          string                    `json:"author"`
	LengthSeconds   int64                     `json:"lengthSeconds"`
	VideoThumbnails []struct{ URL string }    `json:"videoThumbnails"`
	AdaptiveFormats []invidiousAdaptiveFormat `json:"adaptiveFormats"`
	FormatStreams   []invidiousFormatStream   `json:"formatStreams"`
}

func (r *StreamResolver) resolveInvidious(instanceURL, videoID string) (*StreamResult, error) {
	apiURL := fmt.Sprintf("%s/api/v1/videos/%s", strings.TrimRight(instanceURL, "/"), videoID)

	resp, err := r.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invidious video status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data invidiousVideoResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	bestAudioURL := ""
	var maxBitrate int64 = -1
	for _, f := range data.AdaptiveFormats {
		if strings.HasPrefix(f.Type, "audio/") {
			if f.Bitrate > maxBitrate {
				maxBitrate = f.Bitrate
				bestAudioURL = f.URL
			}
		}
	}

	bestVideoURL := ""
	if len(data.FormatStreams) > 0 {
		bestVideoURL = data.FormatStreams[0].URL
	}

	thumb := ""
	if len(data.VideoThumbnails) > 0 {
		thumb = data.VideoThumbnails[0].URL
	}

	return &StreamResult{
		VideoID:   videoID,
		Title:     data.Title,
		Artist:    data.Author,
		Thumbnail: thumb,
		Duration:  data.LengthSeconds,
		AudioURL:  bestAudioURL,
		VideoURL:  bestVideoURL,
	}, nil
}
