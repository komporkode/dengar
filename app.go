package main

import (
	"context"
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
}

// Search queries Piped/Invidious instances for songs or videos.
func (a *App) Search(query string, filterType string) ([]backend.SearchResultItem, error) {
	return a.search.Search(query, filterType)
}

// GetStream resolves streaming media URLs with 30-min TTL caching.
func (a *App) GetStream(videoID string, forceRefresh bool) (*backend.StreamResult, error) {
	return a.resolver.Resolve(videoID, forceRefresh)
}

// GetUserData loads bookmarks and playback history.
func (a *App) GetUserData() (*backend.UserData, error) {
	if a.storage == nil {
		return &backend.UserData{
			Favorites: []backend.SearchResultItem{},
			History:   []backend.SearchResultItem{},
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
