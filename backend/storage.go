package backend

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

const maxHistoryItems = 100

// StorageService manages persistence for user favorites and playback history.
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
			Favorites: []SearchResultItem{},
			History:   []SearchResultItem{},
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
		Favorites: append([]SearchResultItem{}, s.data.Favorites...),
		History:   append([]SearchResultItem{}, s.data.History...),
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
