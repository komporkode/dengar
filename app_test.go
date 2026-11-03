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
		{"v1.1.0", "v1.2.0", false},
		{"v1.1.2", "v1.2.0", false},
		{"v1.2.0", "v1.2.0", false},
		{"1.1.2", "v1.2.0", false},
		{"1.0.9", "1.2.0", false},
		{"", "v1.2.0", false},

		// Upgrade SHOULD trigger update
		{"v1.2.1", "v1.2.0", true},
		{"v1.3.0", "v1.2.0", true},
		{"v2.0.0", "v1.2.0", true},
		{"1.2.1", "v1.2.0", true},
		{"v1.1.10", "v1.1.9", true},
	}

	for _, tc := range tests {
		actual := isVersionNewer(tc.candidate, tc.current)
		if actual != tc.expected {
			t.Errorf("isVersionNewer(%q, %q) = %v; expected %v", tc.candidate, tc.current, actual, tc.expected)
		}
	}
}

func TestRecommendedMultiArtistDiversity(t *testing.T) {
	artists := []string{"Sheila on 7", "Coldplay", "Tulus", "Hindia"}
	mockCatalog := map[string][]backend.SearchResultItem{
		"Sheila on 7": {
			{ID: "s1", Title: "Dan", UploaderName: "Sheila on 7"},
			{ID: "s2", Title: "Sephia", UploaderName: "Sheila on 7"},
			{ID: "s3", Title: "Melompat Lebih Tinggi", UploaderName: "Sheila on 7"},
		},
		"Coldplay": {
			{ID: "c1", Title: "Yellow", UploaderName: "Coldplay"},
			{ID: "c2", Title: "Fix You", UploaderName: "Coldplay"},
			{ID: "c3", Title: "Viva La Vida", UploaderName: "Coldplay"},
		},
		"Tulus": {
			{ID: "t1", Title: "Hati-Hati di Jalan", UploaderName: "Tulus"},
			{ID: "t2", Title: "Monokrom", UploaderName: "Tulus"},
			{ID: "t3", Title: "Diri", UploaderName: "Tulus"},
		},
		"Hindia": {
			{ID: "h1", Title: "Secukupnya", UploaderName: "Hindia"},
			{ID: "h2", Title: "Evaluasi", UploaderName: "Hindia"},
			{ID: "h3", Title: "Rumah ke Rumah", UploaderName: "Hindia"},
		},
	}

	mockSearch := func(query, filterType string) (*backend.SearchResponse, error) {
		for a, items := range mockCatalog {
			if a == query || query == a+" lagu" || query == a+" mix" {
				return &backend.SearchResponse{Items: items}, nil
			}
		}
		return &backend.SearchResponse{Items: []backend.SearchResultItem{}}, nil
	}

	seenIDs := map[string]bool{"s1": true} // Simulated already in top played
	results := buildRecommendedTracks(artists, mockSearch, seenIDs)

	if len(results) == 0 {
		t.Fatal("Expected recommended tracks, got 0")
	}

	artistCounts := make(map[string]int)
	for _, track := range results {
		if track.ID == "s1" {
			t.Errorf("Track s1 should have been excluded because it was in seenIDs")
		}
		artistCounts[track.UploaderName]++
		if artistCounts[track.UploaderName] > 2 {
			t.Errorf("Artist %s exceeded cap of 2 tracks (count: %d)", track.UploaderName, artistCounts[track.UploaderName])
		}
	}

	if len(artistCounts) < 3 {
		t.Errorf("Expected at least 3 distinct artists, got %d: %+v", len(artistCounts), artistCounts)
	}
}
