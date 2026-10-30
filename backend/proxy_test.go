package backend

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestAudioProxyHandler(t *testing.T) {
	// Mock upstream media server
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rangeHeader := r.Header.Get("Range")
		if rangeHeader == "bytes=0-4" {
			w.Header().Set("Content-Type", "audio/webm")
			w.Header().Set("Content-Range", "bytes 0-4/10")
			w.Header().Set("Content-Length", "5")
			w.WriteHeader(http.StatusPartialContent)
			w.Write([]byte("HELLO"))
			return
		}

		w.Header().Set("Content-Type", "audio/webm")
		w.Header().Set("Content-Length", "10")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HELLOWORLD"))
	}))
	defer upstream.Close()

	proxy := NewAudioProxyHandler()

	// 1. Test CORS preflight OPTIONS
	optReq := httptest.NewRequest(http.MethodOptions, "/api/audio-proxy", nil)
	optRec := httptest.NewRecorder()
	proxy.ServeHTTP(optRec, optReq)

	if optRec.Code != http.StatusNoContent {
		t.Errorf("Expected 204 for OPTIONS, got: %d", optRec.Code)
	}
	if optRec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected CORS header '*', got: %s", optRec.Header().Get("Access-Control-Allow-Origin"))
	}

	// 2. Test GET stream with range
	getURL := "/api/audio-proxy?url=" + url.QueryEscape(upstream.URL)
	getReq := httptest.NewRequest(http.MethodGet, getURL, nil)
	getReq.Header.Set("Range", "bytes=0-4")
	getRec := httptest.NewRecorder()
	proxy.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusPartialContent {
		t.Errorf("Expected 206 Partial Content, got: %d", getRec.Code)
	}
	if getRec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected CORS header '*', got: %s", getRec.Header().Get("Access-Control-Allow-Origin"))
	}
	body, _ := io.ReadAll(getRec.Body)
	if string(body) != "HELLO" {
		t.Errorf("Expected 'HELLO', got: '%s'", string(body))
	}
}

