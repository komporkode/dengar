package main

import (
	"archive/zip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"dengar/backend"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct manages runtime context and backend services.
type App struct {
	ctx           context.Context
	search        *backend.SearchService
	resolver      *backend.StreamResolver
	storage       *backend.StorageService
	recCacheMutex sync.Mutex
	cachedRec     []backend.SearchResultItem
	lastRecFetch  time.Time
}

// HomeData holds sections for the personalized home feed.
type HomeData struct {
	TopPlayed   []backend.SearchResultItem `json:"topPlayed"`
	Recommended []backend.SearchResultItem `json:"recommended"`
}

// UpdateAsset holds information about a specific downloadable release asset.
type UpdateAsset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"downloadUrl"`
	Size        int64  `json:"size"`
}

// UpdateInfo holds details for GitHub application release updates.
type UpdateInfo struct {
	CurrentVersion string        `json:"currentVersion"`
	LatestVersion  string        `json:"latestVersion"`
	HasUpdate      bool          `json:"hasUpdate"`
	ReleaseURL     string        `json:"releaseUrl"`
	ReleaseNotes   string        `json:"releaseNotes"`
	PublishedAt    string        `json:"publishedAt"`
	DownloadURL    string        `json:"downloadUrl"`
	AssetFileName  string        `json:"assetFileName"`
	AssetSize      int64         `json:"assetSize"`
	Assets         []UpdateAsset `json:"assets"`
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

// GetStream resolves streaming media URLs with 30-min TTL caching and quality preference ("high" or "saver").
func (a *App) GetStream(videoID string, forceRefresh bool, quality string) (*backend.StreamResult, error) {
	return a.resolver.Resolve(videoID, forceRefresh, quality)
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
	trimmed = strings.TrimPrefix(trimmed, "dengerin://playlist/")
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

// GetHomeData returns personalized recommendations and frequently played tracks with 2-hour cache.
func (a *App) GetHomeData() (*HomeData, error) {
	return a.getHomeDataInternal(false)
}

// RefreshHomeData forces fresh recommendations directly from network.
func (a *App) RefreshHomeData() (*HomeData, error) {
	return a.getHomeDataInternal(true)
}

func (a *App) getHomeDataInternal(forceRefresh bool) (*HomeData, error) {
	result := &HomeData{
		TopPlayed:   []backend.SearchResultItem{},
		Recommended: []backend.SearchResultItem{},
	}
	if a.storage == nil {
		return result, nil
	}

	// TopPlayed is 100% local from disk
	top := a.storage.GetTopPlayed(8)
	result.TopPlayed = top

	// Check recommendation in-memory cache (30-minute TTL, can be refreshed manually)
	a.recCacheMutex.Lock()
	if !forceRefresh && len(a.cachedRec) > 0 && time.Since(a.lastRecFetch) < 30*time.Minute {
		result.Recommended = a.cachedRec
		a.recCacheMutex.Unlock()
		return result, nil
	}
	a.recCacheMutex.Unlock()

	userData, _ := a.storage.Load()
	artistMap := make(map[string]int)

	// 1. Top Played: heaviest weight (+3 per track)
	for _, item := range top {
		if item.UploaderName != "" {
			artistMap[item.UploaderName] += 3
		}
	}

	// 2. Favorites: high affinity (+2 per favorited track)
	if userData != nil {
		for _, item := range userData.Favorites {
			if item.UploaderName != "" {
				artistMap[item.UploaderName] += 2
			}
		}

		// 3. Recent History: (+1 for most recent tracks)
		for i, item := range userData.History {
			if item.UploaderName != "" && i < 10 {
				artistMap[item.UploaderName] += 1
			}
		}
	}

	// Sort artists by affinity score descending (highest affinity first)
	type scoredArtist struct {
		name  string
		score int
	}
	var scoredList []scoredArtist
	for name, score := range artistMap {
		scoredList = append(scoredList, scoredArtist{name: name, score: score})
	}
	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].score > scoredList[j].score
	})

	var topArtists []string
	for i := 0; i < len(scoredList) && i < 4; i++ {
		topArtists = append(topArtists, scoredList[i].name)
	}

	var recTracks []backend.SearchResultItem
	seenIDs := make(map[string]bool)
	for _, item := range top {
		seenIDs[item.ID] = true
	}

	// Fetch tracks from user's actual favorite artists
	if len(topArtists) > 0 {
		for _, artist := range topArtists {
			resp, err := a.search.Search(artist, "music")
			if err == nil && resp != nil {
				for _, item := range resp.Items {
					if !seenIDs[item.ID] {
						seenIDs[item.ID] = true
						recTracks = append(recTracks, item)
						if len(recTracks) >= 16 {
							break
						}
					}
				}
			}
			if len(recTracks) >= 16 {
				break
			}
		}
	}

	// Cold-start fallback for new users with little/no history:
	// Use curated regional & popular queries rather than bare generic queries
	if len(recTracks) < 6 {
		fallbacks := []string{"Top Hits Indonesia", "Lagu Populer Indonesia", "Viral Music Hits"}
		for _, q := range fallbacks {
			resp, err := a.search.Search(q, "music")
			if err == nil && resp != nil {
				for _, item := range resp.Items {
					if !seenIDs[item.ID] {
						seenIDs[item.ID] = true
						recTracks = append(recTracks, item)
						if len(recTracks) >= 16 {
							break
						}
					}
				}
			}
			if len(recTracks) >= 16 {
				break
			}
		}
	}

	a.recCacheMutex.Lock()
	a.cachedRec = recTracks
	a.lastRecFetch = time.Now()
	a.recCacheMutex.Unlock()

	result.Recommended = recTracks
	return result, nil
}

// CheckForUpdates checks GitHub releases API for newer version of Dengerin.
func (a *App) CheckForUpdates() (*UpdateInfo, error) {
	const currentVersion = "v1.1.1"
	client := &http.Client{Timeout: 6 * time.Second}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/komporkode/dengar/releases/latest", nil)
	if err != nil {
		return &UpdateInfo{CurrentVersion: currentVersion, HasUpdate: false}, nil
	}
	req.Header.Set("User-Agent", "Dengerin-Desktop-App")

	resp, err := client.Do(req)
	if err != nil {
		return &UpdateInfo{CurrentVersion: currentVersion, HasUpdate: false}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &UpdateInfo{CurrentVersion: currentVersion, HasUpdate: false}, nil
	}

	var release struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
		Body    string `json:"body"`
		PubAt   string `json:"published_at"`
		Assets  []struct {
			Name               string `json:"name"`
			Size               int64  `json:"size"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return &UpdateInfo{CurrentVersion: currentVersion, HasUpdate: false}, nil
	}

	hasUpdate := release.TagName != "" && release.TagName != currentVersion

	var assets []UpdateAsset
	var matchedURL string
	var matchedFileName string
	var matchedSize int64

	osKeyword := "windows"
	if goruntime.GOOS == "darwin" {
		osKeyword = "macos"
	} else if goruntime.GOOS == "linux" {
		osKeyword = "linux"
	}

	for _, ast := range release.Assets {
		assets = append(assets, UpdateAsset{
			Name:        ast.Name,
			DownloadURL: ast.BrowserDownloadURL,
			Size:        ast.Size,
		})
		if strings.Contains(strings.ToLower(ast.Name), osKeyword) {
			matchedURL = ast.BrowserDownloadURL
			matchedFileName = ast.Name
			matchedSize = ast.Size
		}
	}

	if matchedURL == "" && len(assets) > 0 {
		matchedURL = assets[0].DownloadURL
		matchedFileName = assets[0].Name
		matchedSize = assets[0].Size
	}

	return &UpdateInfo{
		CurrentVersion: currentVersion,
		LatestVersion:  release.TagName,
		HasUpdate:      hasUpdate,
		ReleaseURL:     release.HTMLURL,
		ReleaseNotes:   release.Body,
		PublishedAt:    release.PubAt,
		DownloadURL:    matchedURL,
		AssetFileName:  matchedFileName,
		AssetSize:      matchedSize,
		Assets:         assets,
	}, nil
}

// DownloadUpdate downloads the update file directly to the user's Downloads folder
// and returns the full absolute path of the saved file.
func (a *App) DownloadUpdate(downloadURL, fileName string) (string, error) {
	if downloadURL == "" {
		return "", fmt.Errorf("URL unduhan tidak valid atau kosong")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("tidak dapat menemukan direktori pengguna: %w", err)
	}

	downloadsDir := filepath.Join(homeDir, "Downloads")
	if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
		downloadsDir = homeDir
	}

	if fileName == "" {
		fileName = "dengerin-update.zip"
	}
	targetPath := filepath.Join(downloadsDir, fileName)

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("gagal membuat request unduhan: %w", err)
	}
	req.Header.Set("User-Agent", "Dengerin-Desktop-App")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("koneksi unduhan gagal: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server merespons dengan status HTTP %d", resp.StatusCode)
	}

	outFile, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file di folder Downloads: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return "", fmt.Errorf("gagal menyimpan file unduhan: %w", err)
	}

	return targetPath, nil
}

// OpenInFolder reveals the downloaded update file in the system file manager (Explorer, Finder, xdg-open).
func (a *App) OpenInFolder(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path kosong")
	}

	cleanPath := filepath.Clean(filePath)

	switch goruntime.GOOS {
	case "windows":
		return exec.Command("explorer", "/select,", cleanPath).Start()
	case "darwin":
		return exec.Command("open", "-R", cleanPath).Start()
	default:
		dir := filepath.Dir(cleanPath)
		return exec.Command("xdg-open", dir).Start()
	}
}

// ApplyUpdateAndRestart extracts the updated binary or runs the installer, applies it in-place, and restarts the application.
func (a *App) ApplyUpdateAndRestart(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("path file pembaruan kosong")
	}

	cleanPath := filepath.Clean(filePath)
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return fmt.Errorf("file pembaruan tidak ditemukan: %s", cleanPath)
	}

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("gagal mendapatkan path eksekusi aplikasi: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(cleanPath))

	switch goruntime.GOOS {
	case "windows":
		// Case 1: Downloaded asset is an installer (.exe)
		if ext == ".exe" {
			cmd := exec.Command("cmd.exe", "/c", "start", "", cleanPath)
			if err := cmd.Start(); err != nil {
				return fmt.Errorf("gagal memulai installer: %w", err)
			}
			go func() {
				time.Sleep(500 * time.Millisecond)
				wailsRuntime.Quit(a.ctx)
			}()
			return nil
		}

		// Case 2: Downloaded asset is a .zip archive - extract dengar.exe in-place and run updater
		if ext == ".zip" {
			r, err := zip.OpenReader(cleanPath)
			if err != nil {
				return fmt.Errorf("gagal membuka file zip pembaruan: %w", err)
			}
			defer r.Close()

			appDir := filepath.Dir(execPath)
			updateExePath := filepath.Join(appDir, "dengerin_new.exe")

			var foundExe bool
			for _, f := range r.File {
				nameLower := strings.ToLower(filepath.Base(f.Name))
				if nameLower == "dengerin.exe" || nameLower == "dengar.exe" || strings.HasSuffix(nameLower, ".exe") {
					rc, err := f.Open()
					if err != nil {
						return fmt.Errorf("gagal membaca biner baru dari zip: %w", err)
					}
					outFile, err := os.OpenFile(updateExePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
					if err != nil {
						rc.Close()
						return fmt.Errorf("gagal mengekstrak biner baru: %w", err)
					}
					_, err = io.Copy(outFile, rc)
					outFile.Close()
					rc.Close()
					if err != nil {
						return fmt.Errorf("gagal menyimpan biner baru: %w", err)
					}
					foundExe = true
					break
				}
			}

			if !foundExe {
				// No exe found in zip, fallback to opening folder
				return a.OpenInFolder(cleanPath)
			}

			// Generate updater batch script to replace running binary after exit
			updaterBat := filepath.Join(appDir, "dengerin_updater.bat")
			batContent := fmt.Sprintf("@echo off\r\ntimeout /t 1 /nobreak > NUL\r\nmove /y \"%s\" \"%s\"\r\nstart \"\" \"%s\"\r\ndel \"%%~f0\"\r\n",
				updateExePath, execPath, execPath)
			if err := os.WriteFile(updaterBat, []byte(batContent), 0755); err != nil {
				return fmt.Errorf("gagal membuat script updater: %w", err)
			}

			// Run updater script in background detached
			cmd := exec.Command("cmd.exe", "/c", updaterBat)
			if err := cmd.Start(); err != nil {
				return fmt.Errorf("gagal mengeksekusi updater: %w", err)
			}

			go func() {
				time.Sleep(500 * time.Millisecond)
				wailsRuntime.Quit(a.ctx)
			}()
			return nil
		}

		return a.OpenInFolder(cleanPath)
	case "darwin":
		if ext == ".dmg" {
			exec.Command("open", cleanPath).Start()
			go func() {
				time.Sleep(500 * time.Millisecond)
				wailsRuntime.Quit(a.ctx)
			}()
			return nil
		}
		return a.OpenInFolder(cleanPath)
	default:
		return a.OpenInFolder(cleanPath)
	}
}

