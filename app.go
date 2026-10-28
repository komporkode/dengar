package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"dengar/backend"
)

// App struct manages runtime context and backend services.
type App struct {
	ctx      context.Context
	search   *backend.SearchService
	resolver *backend.StreamResolver
	storage  *backend.StorageService
}

// HomeData holds sections for the personalized home feed.
type HomeData struct {
	TopPlayed   []backend.SearchResultItem `json:"topPlayed"`
	Recommended []backend.SearchResultItem `json:"recommended"`
}

// NewApp creates a new App application struct.
func NewApp() *App {
	pool := backend.NewInstancePool(nil, nil)
	searchSvc := backend.NewSearchService(nil, nil)
	resolverSvc := backend.NewStreamResolver(pool, 30*time.Minute)
	storageSvc, err := backend.NewStorageService("")
	if err != nil {
		println("Warning: failed to initialize storage:", err.Error())
	}

	return &App{
		search:   searchSvc,
		resolver: resolverSvc,
		storage:  storageSvc,
	}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	setupWindowsAppIcon()
}

// Search queries YouTube/Piped/Invidious instances for songs or videos with query caching.
func (a *App) Search(query string, filterType string) (*backend.SearchResponse, error) {
	return a.search.Search(query, filterType)
}

// SearchMore fetches subsequent pages of tracks using the continuation token.
func (a *App) SearchMore(continuationToken string) (*backend.SearchResponse, error) {
	return a.search.SearchMore(continuationToken)
}

// GetStream resolves streaming media URLs with 30-min TTL caching.
func (a *App) GetStream(videoID string, forceRefresh bool) (*backend.StreamResult, error) {
	return a.resolver.Resolve(videoID, forceRefresh)
}

// GetUserData loads bookmarks and playback history.
func (a *App) GetUserData() (*backend.UserData, error) {
	if a.storage == nil {
		return &backend.UserData{
			Favorites:  []backend.SearchResultItem{},
			History:    []backend.SearchResultItem{},
			Playlists:  []backend.Playlist{},
			PlayCounts: make(map[string]int),
		}, nil
	}
	return a.storage.Load()
}

// ToggleFavorite toggles favorite status for an item and returns true if now favorited.
func (a *App) ToggleFavorite(item backend.SearchResultItem) (bool, error) {
	if a.storage == nil {
		return false, nil
	}
	return a.storage.ToggleFavorite(item)
}

// AddToHistory adds a played track to local history.
func (a *App) AddToHistory(item backend.SearchResultItem) error {
	if a.storage == nil {
		return nil
	}
	return a.storage.AddHistory(item)
}

// RecordPlay records that a track was played and increments its play count.
func (a *App) RecordPlay(item backend.SearchResultItem) error {
	if a.storage == nil {
		return nil
	}
	return a.storage.RecordPlay(item)
}

// CreatePlaylist creates a new custom playlist.
func (a *App) CreatePlaylist(name, description string) (*backend.Playlist, error) {
	if a.storage == nil {
		return nil, fmt.Errorf("storage service unavailable")
	}
	return a.storage.CreatePlaylist(name, description)
}

// UpdatePlaylist updates a playlist's title and description.
func (a *App) UpdatePlaylist(id, name, description string) (*backend.Playlist, error) {
	if a.storage == nil {
		return nil, fmt.Errorf("storage service unavailable")
	}
	return a.storage.UpdatePlaylist(id, name, description)
}

// DeletePlaylist deletes a playlist by ID.
func (a *App) DeletePlaylist(id string) error {
	if a.storage == nil {
		return fmt.Errorf("storage service unavailable")
	}
	return a.storage.DeletePlaylist(id)
}

// AddTrackToPlaylist adds a track to a playlist.
func (a *App) AddTrackToPlaylist(playlistID string, track backend.SearchResultItem) error {
	if a.storage == nil {
		return fmt.Errorf("storage service unavailable")
	}
	return a.storage.AddTrackToPlaylist(playlistID, track)
}

// RemoveTrackFromPlaylist removes a track from a playlist.
func (a *App) RemoveTrackFromPlaylist(playlistID string, trackID string) error {
	if a.storage == nil {
		return fmt.Errorf("storage service unavailable")
	}
	return a.storage.RemoveTrackFromPlaylist(playlistID, trackID)
}

// ExportPlaylist exports a playlist as a portable Base64 share code.
func (a *App) ExportPlaylist(playlistID string) (string, error) {
	if a.storage == nil {
		return "", fmt.Errorf("storage service unavailable")
	}
	data, err := a.storage.Load()
	if err != nil {
		return "", err
	}
	for _, pl := range data.Playlists {
		if pl.ID == playlistID {
			exportData := struct {
				Name        string                     `json:"name"`
				Description string                     `json:"description"`
				Tracks      []backend.SearchResultItem `json:"tracks"`
			}{
				Name:        pl.Name,
				Description: pl.Description,
				Tracks:      pl.Tracks,
			}
			raw, err := json.Marshal(exportData)
			if err != nil {
				return "", err
			}
			return base64.StdEncoding.EncodeToString(raw), nil
		}
	}
	return "", fmt.Errorf("playlist not found: %s", playlistID)
}

// ImportPlaylist imports a playlist from a Base64 string or JSON payload.
func (a *App) ImportPlaylist(payload string) (*backend.Playlist, error) {
	if a.storage == nil {
		return nil, fmt.Errorf("storage service unavailable")
	}
	trimmed := strings.TrimSpace(payload)
	trimmed = strings.TrimPrefix(trimmed, "dengar://playlist/")

	raw, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		raw = []byte(trimmed)
	}

	var imported struct {
		Name        string                     `json:"name"`
		Description string                     `json:"description"`
		Tracks      []backend.SearchResultItem `json:"tracks"`
	}
	if err := json.Unmarshal(raw, &imported); err != nil {
		return nil, fmt.Errorf("format playlist tidak valid: %v", err)
	}

	if imported.Name == "" {
		imported.Name = "Playlist Impor"
	}

	pl, err := a.storage.CreatePlaylist(imported.Name, imported.Description)
	if err != nil {
		return nil, err
	}

	for _, track := range imported.Tracks {
		_ = a.storage.AddTrackToPlaylist(pl.ID, track)
	}

	updatedData, err := a.storage.Load()
	if err == nil {
		for _, p := range updatedData.Playlists {
			if p.ID == pl.ID {
				return &p, nil
			}
		}
	}

	return pl, nil
}

// GetHomeData returns personalized recommendations and frequently played tracks.
func (a *App) GetHomeData() (*HomeData, error) {
	result := &HomeData{
		TopPlayed:   []backend.SearchResultItem{},
		Recommended: []backend.SearchResultItem{},
	}
	if a.storage == nil {
		return result, nil
	}

	top := a.storage.GetTopPlayed(8)
	result.TopPlayed = top

	userData, _ := a.storage.Load()
	artistMap := make(map[string]int)

	for _, item := range top {
		if item.UploaderName != "" {
			artistMap[item.UploaderName] += 2
		}
	}
	if userData != nil {
		for _, item := range userData.Favorites {
			if item.UploaderName != "" {
				artistMap[item.UploaderName] += 1
			}
		}
	}

	var topArtists []string
	for artist := range artistMap {
		topArtists = append(topArtists, artist)
		if len(topArtists) >= 2 {
			break
		}
	}

	var recTracks []backend.SearchResultItem
	seenIDs := make(map[string]bool)
	for _, item := range top {
		seenIDs[item.ID] = true
	}

	if len(topArtists) > 0 {
		for _, artist := range topArtists {
			resp, err := a.search.Search(artist, "music")
			if err == nil && resp != nil {
				for _, item := range resp.Items {
					if !seenIDs[item.ID] {
						seenIDs[item.ID] = true
						recTracks = append(recTracks, item)
						if len(recTracks) >= 12 {
							break
						}
					}
				}
			}
			if len(recTracks) >= 12 {
				break
			}
		}
	}

	if len(recTracks) < 6 {
		resp, err := a.search.Search("Trending Music", "music")
		if err == nil && resp != nil {
			for _, item := range resp.Items {
				if !seenIDs[item.ID] {
					seenIDs[item.ID] = true
					recTracks = append(recTracks, item)
					if len(recTracks) >= 12 {
						break
					}
				}
			}
		}
	}

	result.Recommended = recTracks
	return result, nil
}

