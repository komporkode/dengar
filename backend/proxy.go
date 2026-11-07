package backend

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AudioProxyHandler handles proxying audio streams to WebView with permissive CORS and range support.
type AudioProxyHandler struct {
	client       *http.Client
	cacheManager *AudioCacheManager
}

// NewAudioProxyHandler creates a new AudioProxyHandler without disk cache.
func NewAudioProxyHandler() *AudioProxyHandler {
	return NewAudioProxyHandlerWithCache(nil)
}

// NewAudioProxyHandlerWithCache creates a new AudioProxyHandler with disk cache support.
func NewAudioProxyHandlerWithCache(cacheManager *AudioCacheManager) *AudioProxyHandler {
	return &AudioProxyHandler{
		client: &http.Client{
			Timeout: 0, // No full-body timeout for continuous audio media streaming
			Transport: &http.Transport{
				ResponseHeaderTimeout: 15 * time.Second,
				IdleConnTimeout:       90 * time.Second,
			},
		},
		cacheManager: cacheManager,
	}
}

func (h *AudioProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Always set CORS headers so Web Audio API never blocks or mutes audio
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Range, Content-Type, Accept, Origin")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Range, Content-Length, Accept-Ranges")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// 2. Check local disk cache first if cache_id is provided
	cacheID := strings.TrimSpace(r.URL.Query().Get("cache_id"))
	if cacheID != "" && h.cacheManager != nil {
		if filePath, mime, found := h.cacheManager.Get(cacheID); found {
			f, err := os.Open(filePath)
			if err == nil {
				defer f.Close()
				stat, err := f.Stat()
				if err == nil && stat.Size() >= MinValidAudioFileSize {
					realMime := DetectAudioMimeType(filePath, mime)
					w.Header().Set("Content-Type", realMime)
					w.Header().Set("Accept-Ranges", "bytes")
					// Use proper audio file extension in synthetic filename so browser decoder identifies it correctly
					ext := ".mp4"
					if strings.Contains(realMime, "webm") {
						ext = ".webm"
					} else if strings.Contains(realMime, "ogg") {
						ext = ".ogg"
					} else if strings.Contains(realMime, "mpeg") {
						ext = ".mp3"
					}
					http.ServeContent(w, r, cacheID+ext, stat.ModTime(), f)
					return
				}
			}
		}
	}

	rawURL := strings.TrimSpace(r.URL.Query().Get("url"))
	if rawURL == "" {
		if cacheID != "" {
			http.Error(w, "cache entry not found or incomplete", http.StatusNotFound)
			return
		}
		http.Error(w, "missing stream url or cache id", http.StatusBadRequest)
		return
	}

	// r.URL.Query().Get("url") has already unescaped the query string once.
	// Only unescape if it still looks percentage-encoded (does not begin with http).
	targetURL := rawURL
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		if unescaped, err := url.QueryUnescape(rawURL); err == nil {
			targetURL = unescaped
		}
	}

	u, err := url.Parse(targetURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		http.Error(w, "invalid stream url", http.StatusBadRequest)
		return
	}

	// 3. Build upstream request forwarding Range header for audio seeking/scrubbing
	req, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create upstream request: %v", err), http.StatusInternalServerError)
		return
	}

	if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// 4. Execute upstream request
	resp, err := h.client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("stream request error: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 5. Mirror critical audio streaming headers
	for _, key := range []string{"Content-Type", "Content-Length", "Content-Range", "Accept-Ranges"} {
		if val := resp.Header.Get(key); val != "" {
			w.Header().Set(key, val)
		}
	}
	if w.Header().Get("Accept-Ranges") == "" {
		w.Header().Set("Accept-Ranges", "bytes")
	}

	videoID := strings.TrimSpace(r.URL.Query().Get("video_id"))
	if videoID != "" && h.cacheManager != nil {
		if _, _, found := h.cacheManager.Get(videoID); !found {
			rangeHeader := r.Header.Get("Range")
			isFromStart := (rangeHeader == "" || strings.HasPrefix(rangeHeader, "bytes=0-"))
			if isFromStart && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent) {
				tmpPath := filepath.Join(h.cacheManager.dir, sanitizeVideoID(videoID)+".tmp")
				// shortcut: if another writer is creating the tmp file, fall back to plain proxying
				if tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644); err == nil {
					w.WriteHeader(resp.StatusCode)
					tee := io.TeeReader(resp.Body, tmpFile)
					written, copyErr := io.Copy(w, tee)
					_ = tmpFile.Close()
					if copyErr == nil && (resp.ContentLength <= 0 || written >= resp.ContentLength) && written >= MinValidAudioFileSize {
						h.cacheManager.CommitStream(videoID, tmpPath, written, resp.Header.Get("Content-Type"))
					} else {
						_ = os.Remove(tmpPath)
					}
					return
				}
			}
		}
	}

	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

// DetectImageMimeType sniffs magic bytes from an image buffer to determine the correct image/* MIME type.
func DetectImageMimeType(data []byte, fallback string) string {
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "image/jpeg"
	}
	if len(data) >= 8 && string(data[:8]) == "\x89PNG\r\n\x1a\n" {
		return "image/png"
	}
	if len(data) >= 12 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return "image/webp"
	}
	if len(data) >= 6 && (string(data[:6]) == "GIF87a" || string(data[:6]) == "GIF89a") {
		return "image/gif"
	}
	if fallback != "" && strings.HasPrefix(fallback, "image/") {
		return fallback
	}
	return "image/jpeg"
}

// FallbackCoverSVG is a dark-mode SVG placeholder returned when an image cannot be loaded.
const FallbackCoverSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="300" height="300" viewBox="0 0 300 300" fill="none"><rect width="300" height="300" fill="#18181b"/><path d="M120 190a30 30 0 1 1-20-28.3V90l80-16v78.3a30 30 0 1 1-20-28.3V104l-40 8v78z" fill="#71717a"/></svg>`

// CoverProxyHandler proxies album artwork/thumbnails and caches them locally for offline display.
type CoverProxyHandler struct {
	client       *http.Client
	cacheManager *AudioCacheManager
}

// NewCoverProxyHandler creates a new CoverProxyHandler without disk cache.
func NewCoverProxyHandler() *CoverProxyHandler {
	return NewCoverProxyHandlerWithCache(nil)
}

// NewCoverProxyHandlerWithCache creates a new CoverProxyHandler with disk cache support.
func NewCoverProxyHandlerWithCache(cacheManager *AudioCacheManager) *CoverProxyHandler {
	return &CoverProxyHandler{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		cacheManager: cacheManager,
	}
}

func (h *CoverProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Set CORS headers and client caching
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// 2. Serve from disk cache if cache_id is provided
	cacheID := strings.TrimSpace(r.URL.Query().Get("cache_id"))
	if cacheID != "" && h.cacheManager != nil {
		if thumbPath, found := h.cacheManager.GetThumbnail(cacheID); found {
			f, err := os.Open(thumbPath)
			if err == nil {
				defer f.Close()
				stat, err := f.Stat()
				if err == nil && stat.Size() >= 100 {
					var header [16]byte
					n, _ := f.Read(header[:])
					_, _ = f.Seek(0, io.SeekStart)
					mime := DetectImageMimeType(header[:n], "image/jpeg")
					w.Header().Set("Content-Type", mime)
					http.ServeContent(w, r, cacheID+".jpg", stat.ModTime(), f)
					return
				}
			}
		}
	}

	// 3. Fall back to upstream URL or YouTube standard thumbnail URL
	rawURL := strings.TrimSpace(r.URL.Query().Get("url"))
	targetURL := rawURL
	if targetURL == "" && cacheID != "" {
		targetURL = fmt.Sprintf("https://i.ytimg.com/vi/%s/hqdefault.jpg", cacheID)
	}

	if targetURL != "" && !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		if unescaped, err := url.QueryUnescape(targetURL); err == nil {
			targetURL = unescaped
		}
	}

	if targetURL != "" {
		u, err := url.Parse(targetURL)
		if err == nil && (u.Scheme == "http" || u.Scheme == "https") {
			req, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, nil)
			if err == nil {
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
				resp, err := h.client.Do(req)
				if err == nil {
					defer resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						data, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024))
						if err == nil && len(data) >= 100 {
							mime := DetectImageMimeType(data, resp.Header.Get("Content-Type"))
							w.Header().Set("Content-Type", mime)

							// Cache thumbnail locally if cache_id is provided
							if cacheID != "" && h.cacheManager != nil {
								_ = os.WriteFile(filepath.Join(h.cacheManager.dir, sanitizeVideoID(cacheID)+".thumb"), data, 0644)
							}

							w.WriteHeader(http.StatusOK)
							_, _ = w.Write(data)
							return
						}
					}
				}
			}
		}
	}

	// 4. Upstream failed or device offline: return clean SVG placeholder
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(FallbackCoverSVG))
}

