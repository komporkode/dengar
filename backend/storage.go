package backend

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

const maxHistoryItems = 100

// StorageService manages persistence for user favorites, history, playlists, and play counts.
type StorageService struct {
	mu       sync.RWMutex
	filePath string
	data     UserData
}

// NewStorageService creates a storage service at the given path (or default %APPDATA%/dengar).
func NewStorageService(customPath string) (*StorageService, error) {
	path := customPath
	if path == "" {
		cfgDir, err := os.UserConfigDir()
		if err != nil {
			cfgDir = "."
		}
		dir := filepath.Join(cfgDir, "dengar")
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
		path = filepath.Join(dir, "userdata.json")
	} else {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err
		}
	}

	s := &StorageService{
		filePath: path,
		data: UserData{
			Favorites:  []SearchResultItem{},
			History:    []SearchResultItem{},
			Playlists:  []Playlist{},
			PlayCounts: make(map[string]int),
			TrackMeta:  make(map[string]SearchResultItem),
		},
	}

	_ = s.loadFromDisk() // If file doesn't exist yet, empty data is initialized
	return s, nil
}

func (s *StorageService) loadFromDisk() error {
	raw, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}
	var loaded UserData
	if err := json.Unmarshal(raw, &loaded); err != nil {
		return err
	}
	if loaded.Favorites == nil {
		loaded.Favorites = []SearchResultItem{}
	}
	if loaded.History == nil {
		loaded.History = []SearchResultItem{}
	}
	if loaded.Playlists == nil {
		loaded.Playlists = []Playlist{}
	}
	if loaded.PlayCounts == nil {
		loaded.PlayCounts = make(map[string]int)
	}
	if loaded.TrackMeta == nil {
		loaded.TrackMeta = make(map[string]SearchResultItem)
	}
	s.data = loaded
	return nil
}

func (s *StorageService) saveToDisk() error {
	raw, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	// Atomic write using temporary file
	tmpPath := s.filePath + ".tmp"
	if err := os.WriteFile(tmpPath, raw, 0644); err != nil {
		return err
	}
	return os.Rename(tmpPath, s.filePath)
}

// Load returns the current UserData snapshot.
func (s *StorageService) Load() (*UserData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := UserData{
		Favorites:  append([]SearchResultItem{}, s.data.Favorites...),
		History:    append([]SearchResultItem{}, s.data.History...),
		Playlists:  append([]Playlist{}, s.data.Playlists...),
		PlayCounts: make(map[string]int, len(s.data.PlayCounts)),
		TrackMeta:  make(map[string]SearchResultItem, len(s.data.TrackMeta)),
	}
	for k, v := range s.data.PlayCounts {
		res.PlayCounts[k] = v
	}
	for k, v := range s.data.TrackMeta {
		res.TrackMeta[k] = v
	}
	return &res, nil
}

// AddFavorite adds an item to favorites if not already present.
func (s *StorageService) AddFavorite(item SearchResultItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, f := range s.data.Favorites {
		if f.ID == item.ID {
			return nil // already favorited
		}
	}
	s.data.Favorites = append([]SearchResultItem{item}, s.data.Favorites...)
	return s.saveToDisk()
}

// RemoveFavorite removes an item from favorites by ID.
func (s *StorageService) RemoveFavorite(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filtered := make([]SearchResultItem, 0, len(s.data.Favorites))
	for _, f := range s.data.Favorites {
		if f.ID != id {
			filtered = append(filtered, f)
		}
	}
	s.data.Favorites = filtered
	return s.saveToDisk()
}

// ToggleFavorite toggles favorite status for an item and returns the new isFavorite status.
func (s *StorageService) ToggleFavorite(item SearchResultItem) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	isFav := false
	filtered := make([]SearchResultItem, 0, len(s.data.Favorites))
	for _, f := range s.data.Favorites {
		if f.ID == item.ID {
			isFav = true
		} else {
			filtered = append(filtered, f)
		}
	}

	if isFav {
		// remove it
		s.data.Favorites = filtered
		return false, s.saveToDisk()
	}

	// add it
	s.data.Favorites = append([]SearchResultItem{item}, s.data.Favorites...)
	return true, s.saveToDisk()
}

// AddHistory prepends an item to playback history, removes earlier duplicates, and caps at maxHistoryItems.
func (s *StorageService) AddHistory(item SearchResultItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Filter out previous occurrence of this item
	filtered := make([]SearchResultItem, 0, len(s.data.History)+1)
	filtered = append(filtered, item)
	for _, h := range s.data.History {
		if h.ID != item.ID {
			filtered = append(filtered, h)
		}
		if len(filtered) >= maxHistoryItems {
			break
		}
	}
	s.data.History = filtered
	return s.saveToDisk()
}

// CreatePlaylist creates a new custom playlist and persists it.
func (s *StorageService) CreatePlaylist(name, description string) (*Playlist, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().Unix()
	pl := Playlist{
		ID:          fmt.Sprintf("pl-%d", time.Now().UnixNano()),
		Name:        name,
		Description: description,
		Tracks:      []SearchResultItem{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.data.Playlists = append(s.data.Playlists, pl)
	if err := s.saveToDisk(); err != nil {
		return nil, err
	}
	return &pl, nil
}

// UpdatePlaylist updates a playlist's name and description.
func (s *StorageService) UpdatePlaylist(id, name, description string) (*Playlist, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.data.Playlists {
		if s.data.Playlists[i].ID == id {
			s.data.Playlists[i].Name = name
			s.data.Playlists[i].Description = description
			s.data.Playlists[i].UpdatedAt = time.Now().Unix()
			if err := s.saveToDisk(); err != nil {
				return nil, err
			}
			res := s.data.Playlists[i]
			return &res, nil
		}
	}
	return nil, fmt.Errorf("playlist not found: %s", id)
}

// DeletePlaylist removes a playlist by ID.
func (s *StorageService) DeletePlaylist(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filtered := make([]Playlist, 0, len(s.data.Playlists))
	for _, pl := range s.data.Playlists {
		if pl.ID != id {
			filtered = append(filtered, pl)
		}
	}
	s.data.Playlists = filtered
	return s.saveToDisk()
}

// AddTrackToPlaylist adds a track to a playlist if not already present.
func (s *StorageService) AddTrackToPlaylist(playlistID string, track SearchResultItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.data.Playlists {
		if s.data.Playlists[i].ID == playlistID {
			for _, t := range s.data.Playlists[i].Tracks {
				if t.ID == track.ID {
					return nil // already in playlist
				}
			}
			s.data.Playlists[i].Tracks = append(s.data.Playlists[i].Tracks, track)
			s.data.Playlists[i].UpdatedAt = time.Now().Unix()
			return s.saveToDisk()
		}
	}
	return fmt.Errorf("playlist not found: %s", playlistID)
}

// RemoveTrackFromPlaylist removes a track from a playlist by ID.
func (s *StorageService) RemoveTrackFromPlaylist(playlistID string, trackID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.data.Playlists {
		if s.data.Playlists[i].ID == playlistID {
			filtered := make([]SearchResultItem, 0, len(s.data.Playlists[i].Tracks))
			for _, t := range s.data.Playlists[i].Tracks {
				if t.ID != trackID {
					filtered = append(filtered, t)
				}
			}
			s.data.Playlists[i].Tracks = filtered
			s.data.Playlists[i].UpdatedAt = time.Now().Unix()
			return s.saveToDisk()
		}
	}
	return fmt.Errorf("playlist not found: %s", playlistID)
}

// RecordPlay increments play count for a track and caches its metadata.
func (s *StorageService) RecordPlay(track SearchResultItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data.PlayCounts[track.ID]++
	s.data.TrackMeta[track.ID] = track
	return s.saveToDisk()
}

// GetTopPlayed returns the most frequently played tracks, sorted descending.
func (s *StorageService) GetTopPlayed(limit int) []SearchResultItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type playEntry struct {
		item  SearchResultItem
		count int
	}

	entries := make([]playEntry, 0, len(s.data.PlayCounts))
	for id, count := range s.data.PlayCounts {
		if item, ok := s.data.TrackMeta[id]; ok {
			entries = append(entries, playEntry{item: item, count: count})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].count > entries[j].count
	})

	if limit > len(entries) {
		limit = len(entries)
	}
	res := make([]SearchResultItem, limit)
	for i := 0; i < limit; i++ {
		res[i] = entries[i].item
	}
	return res
}
