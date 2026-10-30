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

func TestIsVersionNewer(t *testing.T) {
	tests := []struct {
		candidate string
		current   string
		expected  bool
	}{
		// Downgrade or same version should NOT trigger update
		{"v1.1.0", "v1.1.2", false},
		{"v1.1.1", "v1.1.2", false},
		{"v1.1.2", "v1.1.2", false},
		{"1.1.1", "v1.1.2", false},
		{"1.0.9", "1.1.2", false},
		{"", "v1.1.2", false},

		// Upgrade SHOULD trigger update
		{"v1.1.3", "v1.1.2", true},
		{"v1.2.0", "v1.1.2", true},
		{"v2.0.0", "v1.1.2", true},
		{"1.1.3", "v1.1.2", true},
		{"v1.1.10", "v1.1.9", true},
	}

	for _, tc := range tests {
		actual := isVersionNewer(tc.candidate, tc.current)
		if actual != tc.expected {
			t.Errorf("isVersionNewer(%q, %q) = %v; expected %v", tc.candidate, tc.current, actual, tc.expected)
		}
	}
}
