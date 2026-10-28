package backend

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestStorageFavoritesAndHistory(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "userdata.json")

	svc, err := NewStorageService(filePath)
	if err != nil {
		t.Fatalf("NewStorageService failed: %v", err)
	}

	// 1. Test Add Favorite
	fav1 := SearchResultItem{ID: "f1", Title: "Fav Track 1", UploaderName: "Artist 1"}
	if err := svc.AddFavorite(fav1); err != nil {
		t.Fatalf("AddFavorite failed: %v", err)
	}

	data, err := svc.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(data.Favorites) != 1 || data.Favorites[0].ID != "f1" {
		t.Errorf("Expected 1 favorite with ID f1, got: %+v", data.Favorites)
	}

	// 2. Test Remove Favorite
	if err := svc.RemoveFavorite("f1"); err != nil {
		t.Fatalf("RemoveFavorite failed: %v", err)
	}
	data, _ = svc.Load()
	if len(data.Favorites) != 0 {
		t.Errorf("Expected 0 favorites, got: %d", len(data.Favorites))
	}

	// 3. Test Add History with capping at 100
	for i := 0; i < 110; i++ {
		item := SearchResultItem{
			ID:           fmt.Sprintf("hist-%d", i),
			Title:        fmt.Sprintf("History Track %d", i),
			UploaderName: "Artist",
		}
		if err := svc.AddHistory(item); err != nil {
			t.Fatalf("AddHistory failed at %d: %v", i, err)
		}
	}

	data, _ = svc.Load()
	if len(data.History) != 100 {
		t.Errorf("Expected history capped at 100 items, got: %d", len(data.History))
	}

	// The most recent item should be at index 0 (hist-109)
	if data.History[0].ID != "hist-109" {
		t.Errorf("Expected index 0 to be hist-109, got: %s", data.History[0].ID)
	}
}

func TestPlaylistAndPlayTracking(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "userdata.json")

	svc, err := NewStorageService(filePath)
	if err != nil {
		t.Fatalf("NewStorageService failed: %v", err)
	}

	// 1. Create Playlist
	pl, err := svc.CreatePlaylist("My Rocks", "Fav rock songs")
	if err != nil {
		t.Fatalf("CreatePlaylist failed: %v", err)
	}
	if pl.ID == "" || pl.Name != "My Rocks" {
		t.Fatalf("Unexpected playlist created: %+v", pl)
	}

	// 2. Add Tracks
	t1 := SearchResultItem{ID: "t1", Title: "Rock Song 1", UploaderName: "Band A"}
	t2 := SearchResultItem{ID: "t2", Title: "Rock Song 2", UploaderName: "Band B"}
	if err := svc.AddTrackToPlaylist(pl.ID, t1); err != nil {
		t.Fatalf("AddTrackToPlaylist t1 failed: %v", err)
	}
	if err := svc.AddTrackToPlaylist(pl.ID, t2); err != nil {
		t.Fatalf("AddTrackToPlaylist t2 failed: %v", err)
	}

	data, _ := svc.Load()
	if len(data.Playlists) != 1 || len(data.Playlists[0].Tracks) != 2 {
		t.Fatalf("Expected 1 playlist with 2 tracks, got: %+v", data.Playlists)
	}

	// 3. Remove track t1
	if err := svc.RemoveTrackFromPlaylist(pl.ID, "t1"); err != nil {
		t.Fatalf("RemoveTrackFromPlaylist failed: %v", err)
	}
	data, _ = svc.Load()
	if len(data.Playlists[0].Tracks) != 1 || data.Playlists[0].Tracks[0].ID != "t2" {
		t.Fatalf("Expected 1 track (t2) remaining, got: %+v", data.Playlists[0].Tracks)
	}

	// 4. Update playlist
	updated, err := svc.UpdatePlaylist(pl.ID, "My Rock Hits", "Updated description")
	if err != nil || updated.Name != "My Rock Hits" {
		t.Fatalf("UpdatePlaylist failed: %v, updated: %+v", err, updated)
	}

	// 5. Play tracking
	t3 := SearchResultItem{ID: "t3", Title: "Top Song 3", UploaderName: "Artist C"}
	for i := 0; i < 3; i++ {
		_ = svc.RecordPlay(t2)
	}
	for i := 0; i < 5; i++ {
		_ = svc.RecordPlay(t3)
	}

	top := svc.GetTopPlayed(5)
	if len(top) < 2 || top[0].ID != "t3" || top[1].ID != "t2" {
		t.Fatalf("Expected top played [t3, t2], got: %+v", top)
	}

	// 6. Delete Playlist
	if err := svc.DeletePlaylist(pl.ID); err != nil {
		t.Fatalf("DeletePlaylist failed: %v", err)
	}
	data, _ = svc.Load()
	if len(data.Playlists) != 0 {
		t.Fatalf("Expected 0 playlists after delete, got: %d", len(data.Playlists))
	}
}

