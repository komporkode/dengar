package main

import (
	"path/filepath"
	"testing"

	"dengar/backend"
)

func TestPlaylistExportImport(t *testing.T) {
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "userdata.json")

	storageSvc, err := backend.NewStorageService(storagePath)
	if err != nil {
		t.Fatalf("NewStorageService failed: %v", err)
	}

	app := &App{
		storage: storageSvc,
	}

	// Create playlist
	pl, err := app.CreatePlaylist("My Jams", "Best collection")
	if err != nil {
		t.Fatalf("CreatePlaylist failed: %v", err)
	}

	track := backend.SearchResultItem{
		ID:           "test1234",
		Title:        "Awesome Track",
		UploaderName: "Star Artist",
		Duration:     180,
	}
	if err := app.AddTrackToPlaylist(pl.ID, track); err != nil {
		t.Fatalf("AddTrackToPlaylist failed: %v", err)
	}

	// Export
	b64, err := app.ExportPlaylist(pl.ID)
	if err != nil {
		t.Fatalf("ExportPlaylist failed: %v", err)
	}
	if b64 == "" {
		t.Fatal("Exported code is empty")
	}

	// Import into new app
	imported, err := app.ImportPlaylist(b64)
	if err != nil {
		t.Fatalf("ImportPlaylist failed: %v", err)
	}
	if imported.Name != "My Jams" || len(imported.Tracks) != 1 || imported.Tracks[0].ID != "test1234" {
		t.Fatalf("Imported playlist mismatch: %+v", imported)
	}
}
