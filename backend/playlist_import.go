package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// ImportedPlaylistResult holds parsed playlist information and tracks.
type ImportedPlaylistResult struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Source      string             `json:"source"`
	Tracks      []SearchResultItem `json:"tracks"`
}

// ParsePlaylistURL determines whether the URL is a YouTube or Spotify playlist and extracts its ID.
func ParsePlaylistURL(rawURL string) (source string, id string, valid bool) {
	clean := strings.TrimSpace(rawURL)
	u, err := url.Parse(clean)
	if err != nil {
		return "", "", false
	}

	host := strings.ToLower(u.Hostname())

	// 1. YouTube & YouTube Music Playlist
	if strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be") {
		listID := u.Query().Get("list")
		if listID != "" {
			return "youtube", listID, true
		}
	}

	// 2. Spotify Playlist
	if strings.Contains(host, "spotify.com") {
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		for i, part := range parts {
			if part == "playlist" && i+1 < len(parts) {
				playlistID := parts[i+1]
				if playlistID != "" {
					return "spotify", playlistID, true
				}
			}
		}
	}

	return "", "", false
}

// ParseSpotifyEmbedHTML extracts metadata and tracklist from Spotify public embed HTML.
func ParseSpotifyEmbedHTML(r io.Reader) (string, string, []SearchResultItem, error) {
	bodyBytes, err := io.ReadAll(r)
	if err != nil {
		return "", "", nil, err
	}
	bodyStr := string(bodyBytes)

	// Extract title from <title> or og:title
	name := "Spotify Playlist"
	reTitle := regexp.MustCompile(`(?i)<meta\s+property=["']og:title["']\s+content=["']([^"']+)["']`)
	match := reTitle.FindStringSubmatch(bodyStr)
	if len(match) > 1 {
		name = strings.TrimSpace(match[1])
	} else {
		reHTMLTitle := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
		matchTitle := reHTMLTitle.FindStringSubmatch(bodyStr)
		if len(matchTitle) > 1 {
			name = strings.TrimSpace(matchTitle[1])
			name = strings.TrimPrefix(name, "Spotify - ")
		}
	}

	description := ""
	reDesc := regexp.MustCompile(`(?i)<meta\s+property=["']og:description["']\s+content=["']([^"']+)["']`)
	matchDesc := reDesc.FindStringSubmatch(bodyStr)
	if len(matchDesc) > 1 {
		description = strings.TrimSpace(matchDesc[1])
	}

	var tracks []SearchResultItem

	// Try extracting NEXT_DATA json
	reNext := regexp.MustCompile(`(?s)<script\s+id=["']__NEXT_DATA__["'][^>]*>(.*?)</script>`)
	matchNext := reNext.FindStringSubmatch(bodyStr)
	if len(matchNext) > 1 {
		var nextData struct {
			Props struct {
				PageProps struct {
					State struct {
						Data struct {
							Entity struct {
								Name      string `json:"name"`
								TrackList []struct {
									Title    string `json:"title"`
									Subtitle string `json:"subtitle"`
									Duration int64  `json:"duration"`
								} `json:"trackList"`
							} `json:"entity"`
						} `json:"data"`
					} `json:"state"`
				} `json:"pageProps"`
			} `json:"props"`
		}
		if err := json.Unmarshal([]byte(matchNext[1]), &nextData); err == nil {
			if nextData.Props.PageProps.State.Data.Entity.Name != "" {
				name = nextData.Props.PageProps.State.Data.Entity.Name
			}
			for i, t := range nextData.Props.PageProps.State.Data.Entity.TrackList {
				durSec := t.Duration / 1000
				tracks = append(tracks, SearchResultItem{
					ID:           fmt.Sprintf("sp_%d", i),
					Title:        t.Title,
					UploaderName: t.Subtitle,
					Duration:     durSec,
				})
			}
		}
	}

	return name, description, tracks, nil
}

// PlaylistImportService resolves external playlists from YouTube and Spotify.
type PlaylistImportService struct {
	client *http.Client
}

func NewPlaylistImportService() *PlaylistImportService {
	return &PlaylistImportService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *PlaylistImportService) Import(rawURL string) (*ImportedPlaylistResult, error) {
	source, id, valid := ParsePlaylistURL(rawURL)
	if !valid {
		return nil, fmt.Errorf("URL tidak dikenali. Masukkan tautan playlist YouTube atau Spotify yang valid.")
	}

	if source == "spotify" {
		embedURL := fmt.Sprintf("https://open.spotify.com/embed/playlist/%s", id)
		req, err := http.NewRequest("GET", embedURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
		resp, err := s.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("gagal menghubungi Spotify embed: %w", err)
		}
		defer resp.Body.Close()

		name, desc, tracks, err := ParseSpotifyEmbedHTML(resp.Body)
		if err != nil {
			return nil, err
		}
		if len(tracks) == 0 {
			return nil, fmt.Errorf("tidak dapat menemukan daftar lagu pada playlist Spotify ini")
		}
		return &ImportedPlaylistResult{
			Name:        name,
			Description: desc,
			Source:      "spotify",
			Tracks:      tracks,
		}, nil
	}

	if source == "youtube" {
		return s.fetchYouTubePlaylist(id)
	}

	return nil, fmt.Errorf("sumber playlist tidak didukung")
}

func (s *PlaylistImportService) fetchYouTubePlaylist(playlistID string) (*ImportedPlaylistResult, error) {
	instances := []string{"https://pipedapi.kavin.rocks", "https://api.piped.privacydev.net", "https://pipedapi.adminforge.de"}
	for _, inst := range instances {
		apiURL := fmt.Sprintf("%s/playlists/%s", inst, playlistID)
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")
		resp, err := s.client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			if resp != nil {
				resp.Body.Close()
			}
			continue
		}

		var data struct {
			Name           string `json:"name"`
			Description    string `json:"description"`
			RelatedStreams []struct {
				URL          string `json:"url"`
				Title        string `json:"title"`
				UploaderName string `json:"uploaderName"`
				Duration     int64  `json:"duration"`
				Thumbnail    string `json:"thumbnail"`
			} `json:"relatedStreams"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil && len(data.RelatedStreams) > 0 {
			resp.Body.Close()
			var tracks []SearchResultItem
			for _, item := range data.RelatedStreams {
				videoID := extractVideoID(item.URL)
				tracks = append(tracks, SearchResultItem{
					ID:           videoID,
					Title:        item.Title,
					UploaderName: item.UploaderName,
					Duration:     item.Duration,
					ThumbnailURL: item.Thumbnail,
				})
			}
			return &ImportedPlaylistResult{
				Name:        data.Name,
				Description: data.Description,
				Source:      "youtube",
				Tracks:      tracks,
			}, nil
		}
		resp.Body.Close()
	}

	return nil, fmt.Errorf("gagal mengunduh daftar lagu playlist YouTube dari server publik")
}
