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
