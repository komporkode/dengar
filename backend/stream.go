package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
			Timeout: 12 * time.Second,
		},
	}
}

// Resolve returns StreamResult for videoID. If forceRefresh is false, cached value is preferred.
// quality can be "high" (default highest bitrate Opus 160k / AAC) or "saver" (conserve bandwidth ~70-128k).
func (r *StreamResolver) Resolve(videoID string, forceRefresh bool, quality ...string) (*StreamResult, error) {
	vid := strings.TrimSpace(videoID)
	if vid == "" {
		return nil, fmt.Errorf("video ID kosong")
	}

	q := "high"
	if len(quality) > 0 && strings.TrimSpace(quality[0]) != "" {
		normalized := strings.ToLower(strings.TrimSpace(quality[0]))
		if normalized == "saver" {
			q = "saver"
		}
	}

	cacheKey := fmt.Sprintf("%s:%s", vid, q)

	if !forceRefresh {
		if cached, ok := r.cache.Get(cacheKey); ok {
			return cached, nil
		}
	}

	// 1. Try Invidious instances first (highest availability for audio streams)
	for _, inst := range r.pool.GetInvidiousInstances() {
		res, err := r.resolveInvidious(inst, vid, q)
		if err == nil && res.AudioURL != "" {
			res.Quality = q
			r.cache.Set(cacheKey, res)
			return res, nil
		}
		r.pool.RotateInvidious()
	}

	// 2. Fallback to Piped instances
	for _, inst := range r.pool.GetPipedInstances() {
		res, err := r.resolvePiped(inst, vid, q)
		if err == nil && res.AudioURL != "" {
			res.Quality = q
			r.cache.Set(cacheKey, res)
			return res, nil
		}
		r.pool.RotatePiped()
	}

	// Trigger instance refresh for future requests if all failed
	r.pool.RefreshInvidiousInstances()
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

func formatAudioCodec(mime string, bitrate int64) string {
	codec := "Audio"
	if strings.Contains(mime, "opus") {
		codec = "Opus"
	} else if strings.Contains(mime, "mp4") || strings.Contains(mime, "m4a") || strings.Contains(mime, "aac") {
		codec = "AAC"
	}

	kbps := bitrate / 1000
	if kbps > 0 {
		return fmt.Sprintf("%s %dkbps", codec, kbps)
	}
	return codec
}

func (r *StreamResolver) resolvePiped(instanceURL, videoID, quality string) (*StreamResult, error) {
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

	bestAudioURL := ""
	var selectedBitrate int64 = 0
	var selectedFormat string = ""

	if len(data.AudioStreams) > 0 {
		if quality == "saver" {
			// Saver: choose stream with lower bitrate (<= 135kbps, e.g. ~64-128kbps) to conserve bandwidth.
			var candidate pipedStreamAudio
			foundSaver := false
			for _, a := range data.AudioStreams {
				if a.URL == "" {
					continue
				}
				if a.Bitrate > 0 && a.Bitrate <= 135000 {
					if !foundSaver || a.Bitrate > candidate.Bitrate {
						candidate = a
						foundSaver = true
					}
				}
			}
			if !foundSaver {
				// Fallback to lowest non-zero bitrate
				var minBitrate int64 = 1 << 60
				for _, a := range data.AudioStreams {
					if a.URL == "" {
						continue
					}
					if a.Bitrate > 0 && a.Bitrate < minBitrate {
						minBitrate = a.Bitrate
						candidate = a
						foundSaver = true
					}
				}
			}
			if foundSaver {
				bestAudioURL = candidate.URL
				selectedBitrate = candidate.Bitrate
				selectedFormat = formatAudioCodec(candidate.MimeType, candidate.Bitrate)
			}
		}

		if bestAudioURL == "" {
			// High quality: pick maximum bitrate available
			var maxBitrate int64 = -1
			var candidate pipedStreamAudio
			for _, a := range data.AudioStreams {
				if a.URL == "" {
					continue
				}
				if a.Bitrate > maxBitrate {
					maxBitrate = a.Bitrate
					candidate = a
				}
			}
			if candidate.URL != "" {
				bestAudioURL = candidate.URL
				selectedBitrate = candidate.Bitrate
				selectedFormat = formatAudioCodec(candidate.MimeType, candidate.Bitrate)
			} else {
				for _, a := range data.AudioStreams {
					if a.URL != "" {
						bestAudioURL = a.URL
						selectedBitrate = a.Bitrate
						selectedFormat = formatAudioCodec(a.MimeType, a.Bitrate)
						break
					}
				}
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

	if strings.HasPrefix(bestAudioURL, "/") {
		bestAudioURL = strings.TrimRight(instanceURL, "/") + bestAudioURL
	}
	if strings.HasPrefix(bestVideoURL, "/") {
		bestVideoURL = strings.TrimRight(instanceURL, "/") + bestVideoURL
	}

	return &StreamResult{
		VideoID:   videoID,
		Title:     data.Title,
		Artist:    data.Uploader,
		Thumbnail: data.ThumbnailURL,
		Duration:  data.Duration,
		AudioURL:  bestAudioURL,
		VideoURL:  bestVideoURL,
		Bitrate:   selectedBitrate,
		Format:    selectedFormat,
		Quality:   quality,
	}, nil
}

type invidiousAdaptiveFormat struct {
	URL      string          `json:"url"`
	Type     string          `json:"type"`
	Bitrate  json.RawMessage `json:"bitrate"`
	Encoding string          `json:"encoding"`
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

func parseRawBitrate(raw json.RawMessage) int64 {
	s := strings.Trim(string(raw), `"`)
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

func (r *StreamResolver) resolveInvidious(instanceURL, videoID, quality string) (*StreamResult, error) {
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
	var selectedBitrate int64 = 0
	var selectedFormat string = ""

	type audioCandidate struct {
		url     string
		bitrate int64
		mime    string
	}

	var audioList []audioCandidate
	for _, f := range data.AdaptiveFormats {
		if strings.Contains(f.Type, "audio") && f.URL != "" {
			b := parseRawBitrate(f.Bitrate)
			audioList = append(audioList, audioCandidate{
				url:     f.URL,
				bitrate: b,
				mime:    f.Type,
			})
		}
	}

	if len(audioList) > 0 {
		if quality == "saver" {
			var candidate audioCandidate
			foundSaver := false
			for _, a := range audioList {
				if a.bitrate > 0 && a.bitrate <= 135000 {
					if !foundSaver || a.bitrate > candidate.bitrate {
						candidate = a
						foundSaver = true
					}
				}
			}
			if !foundSaver {
				var minBitrate int64 = 1 << 60
				for _, a := range audioList {
					if a.bitrate > 0 && a.bitrate < minBitrate {
						minBitrate = a.bitrate
						candidate = a
						foundSaver = true
					}
				}
			}
			if foundSaver {
				bestAudioURL = candidate.url
				selectedBitrate = candidate.bitrate
				selectedFormat = formatAudioCodec(candidate.mime, candidate.bitrate)
			}
		}

		if bestAudioURL == "" {
			var maxBitrate int64 = -1
			var candidate audioCandidate
			for _, a := range audioList {
				if a.bitrate > maxBitrate {
					maxBitrate = a.bitrate
					candidate = a
				}
			}
			if candidate.url != "" {
				bestAudioURL = candidate.url
				selectedBitrate = candidate.bitrate
				selectedFormat = formatAudioCodec(candidate.mime, candidate.bitrate)
			} else {
				bestAudioURL = audioList[0].url
				selectedBitrate = audioList[0].bitrate
				selectedFormat = formatAudioCodec(audioList[0].mime, audioList[0].bitrate)
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

	if strings.HasPrefix(bestAudioURL, "/") {
		bestAudioURL = strings.TrimRight(instanceURL, "/") + bestAudioURL
	}
	if strings.HasPrefix(bestVideoURL, "/") {
		bestVideoURL = strings.TrimRight(instanceURL, "/") + bestVideoURL
	}

	return &StreamResult{
		VideoID:   videoID,
		Title:     data.Title,
		Artist:    data.Author,
		Thumbnail: thumb,
		Duration:  data.LengthSeconds,
		AudioURL:  bestAudioURL,
		VideoURL:  bestVideoURL,
		Bitrate:   selectedBitrate,
		Format:    selectedFormat,
		Quality:   quality,
	}, nil
}
