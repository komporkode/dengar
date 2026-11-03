package backend

import (
	"strings"
	"testing"
)

func TestParsePlaylistURL(t *testing.T) {
	tests := []struct {
		url        string
		wantSource string
		wantID     string
		wantValid  bool
	}{
		{"https://www.youtube.com/playlist?list=PL12345ABC", "youtube", "PL12345ABC", true},
		{"https://music.youtube.com/playlist?list=RDCLAK5uy_xyz", "youtube", "RDCLAK5uy_xyz", true},
		{"https://open.spotify.com/playlist/37i9dQZF1DXcBWIGoYBM5M?si=abcd", "spotify", "37i9dQZF1DXcBWIGoYBM5M", true},
		{"https://open.spotify.com/intl-id/playlist/37i9dQZF1DXcBWIGoYBM5M", "spotify", "37i9dQZF1DXcBWIGoYBM5M", true},
		{"https://example.com/other", "", "", false},
	}

	for _, tc := range tests {
		source, id, valid := ParsePlaylistURL(tc.url)
		if valid != tc.wantValid || source != tc.wantSource || id != tc.wantID {
			t.Errorf("ParsePlaylistURL(%q) = (%q, %q, %v); want (%q, %q, %v)",
				tc.url, source, id, valid, tc.wantSource, tc.wantID, tc.wantValid)
		}
	}
}

func TestParseSpotifyEmbedHTML(t *testing.T) {
	mockHTML := `
<!DOCTYPE html>
<html>
<head>
<title>Spotify - Top Hits 2026</title>
<meta property="og:title" content="Top Hits 2026">
<meta property="og:description" content="Koleksi lagu terpopuler">
</head>
<body>
<script id="__NEXT_DATA__" type="application/json">
{
  "props": {
    "pageProps": {
      "state": {
        "data": {
          "entity": {
            "name": "Top Hits 2026",
            "trackList": [
              {"title": "Melompat Lebih Tinggi", "subtitle": "Sheila on 7", "duration": 190000},
              {"title": "Fix You", "subtitle": "Coldplay", "duration": 290000}
            ]
          }
        }
      }
    }
  }
}
</script>
</body>
</html>`

	name, desc, tracks, err := ParseSpotifyEmbedHTML(strings.NewReader(mockHTML))
	if err != nil {
		t.Fatalf("ParseSpotifyEmbedHTML error: %v", err)
	}
	_ = desc

	if name != "Top Hits 2026" {
		t.Errorf("Expected name 'Top Hits 2026', got %q", name)
	}

	if len(tracks) != 2 {
		t.Fatalf("Expected 2 tracks, got %d", len(tracks))
	}

	if tracks[0].Title != "Melompat Lebih Tinggi" || tracks[0].UploaderName != "Sheila on 7" {
		t.Errorf("Track 0 mismatch: %+v", tracks[0])
	}
	if tracks[1].Title != "Fix You" || tracks[1].UploaderName != "Coldplay" {
		t.Errorf("Track 1 mismatch: %+v", tracks[1])
	}
}
