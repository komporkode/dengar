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
			Timeout: 60 * time.Second, // Long timeout for streaming media chunks
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
				if err == nil {
					w.Header().Set("Content-Type", mime)
					w.Header().Set("Accept-Ranges", "bytes")
					http.ServeContent(w, r, filepath.Base(filePath), stat.ModTime(), f)
					return
				}
			}
		}
	}

	rawURL := strings.TrimSpace(r.URL.Query().Get("url"))
	if rawURL == "" {
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

	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}
