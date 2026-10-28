# Personalized Recommendations & Spotify-Style Shareable Playlists Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement a privacy-friendly, zero-server personalized music recommendation engine ("Sering Kamu Putar" & "Rekomendasi Untuk Kamu") and Spotify-style shareable playlist system (Create, Manage, Auto-Collage Cover Art, Dual-Sharing via Base64 code & JSON file export/import) in Dengar.

**Architecture:** 
- **Backend (Go/Wails):** Extend atomic local storage `%APPDATA%/dengar/userdata.json` to store `Playlists` and `PlayCounts`. Implement recommendation heuristics (top played track frequency + top artist discovery using existing search/innertube backend).
- **Frontend (Svelte 4/Tailwind v4):** Multi-section homepage (*Sering Kamu Putar*, *Rekomendasi Untuk Kamu*, *Trending Hari Ini*), Playlist management in Sidebar, context menus on `TrackItem` to add songs to playlists, dedicated Spotify-style `PlaylistView` with dynamic 2x2 collage covers, and a `PlaylistShareModal` supporting clipboard code and JSON file export/import.
- **Constraints:** All Git commits must remain strictly local (**never push to remote**). Zero external paid services; fully offline-first storage.

**Tech Stack:** Go 1.23, Wails v2, Svelte 4, TypeScript, Tailwind CSS v4.

**Spec:** Grilling design agreement recorded on 2026-10-28 (Hybrid local play count + artist recommendations; Multi-section Spotify home; Auto-collage covers; Dual code/file sharing).

---

## Global Constraints
- **Zero Remote Push:** Git commits must be strictly local (`git add` & `git commit`). NEVER run `git push`.
- **Atomic File Storage:** All writes to `%APPDATA%/dengar/userdata.json` must use atomic temp file write + rename to guarantee zero data corruption.
- **Valid Play Threshold:** Play counts increment ONLY after a track plays for >= 20 seconds or completes playback to prevent spurious clicks from skewing recommendations.
- **Backward Compatibility:** Existing `userdata.json` with only `favorites` and `history` must parse seamlessly without data loss.
- **Zero Server Dependency:** Sharing must work peer-to-peer using compressed URL-safe Base64 strings or `.dengarplaylist` JSON files.

## Review Focus
1. `userdata.json` migration: existing user files without `playlists` or `playCounts` fields must load cleanly with non-nil slices and maps.
2. Anti-spam play count: skipping 5 songs rapidly within 5 seconds must NOT increment any play count.
3. Empty playlists: an empty playlist must display an aesthetic default gradient and prompt to add songs without crashing the auto-collage renderer.
4. Corrupted share code: pasting an invalid or malformed Base64 share code into the import dialog must gracefully show an error toast without throwing unhandled exceptions.
5. Large playlist import: importing a playlist with 100+ songs must parse and save smoothly without UI freeze.

---

### Task 1: Backend Storage Extension (Playlists & Play Tracking Models)

**Files:**
- Modify: `backend/models.go`
- Modify: `backend/storage.go`
- Test: `backend/storage_test.go`

**Interfaces:**
- Produces:
  - `Playlist struct { ID, Name, Description, Tracks, CreatedAt, UpdatedAt }`
  - `StorageService.CreatePlaylist(name, description string) (*Playlist, error)`
  - `StorageService.UpdatePlaylist(id, name, description string) (*Playlist, error)`
  - `StorageService.DeletePlaylist(id string) error`
  - `StorageService.AddTrackToPlaylist(playlistID string, track SearchResultItem) error`
  - `StorageService.RemoveTrackFromPlaylist(playlistID string, trackID string) error`
  - `StorageService.RecordPlay(track SearchResultItem) error`
  - `StorageService.GetTopPlayed(limit int) []SearchResultItem`

- [ ] **Step 1: Write failing unit tests for Playlist CRUD and Play Tracking**
  Add tests in `backend/storage_test.go` verifying:
  - Playlist creation with UUID/timestamp generation
  - Adding tracks to playlist and preventing duplicates if already in playlist
  - Removing a track from a playlist
  - Deleting a playlist
  - Recording play count and verifying `GetTopPlayed(limit)` returns tracks ordered by play frequency.

- [ ] **Step 2: Run test to verify failure**
  Run: `go test -v ./backend -run TestPlaylistAndPlayTracking`
  Expected: FAIL (types and methods undefined)

- [ ] **Step 3: Implement Playlist and PlayCounts in `backend/models.go` and `backend/storage.go`**
  - Add `Playlist` struct to `backend/models.go`.
  - Update `UserData` struct in `backend/models.go`:
    ```go
    type UserData struct {
        Favorites  []SearchResultItem `json:"favorites"`
        History    []SearchResultItem `json:"history"`
        Playlists  []Playlist         `json:"playlists"`
        PlayCounts map[string]int     `json:"playCounts"`
    }
    ```
  - Implement CRUD and play count logic with mutex locks and atomic disk saving in `backend/storage.go`.

- [ ] **Step 4: Run tests to verify pass**
  Run: `go test -v ./backend -run TestPlaylistAndPlayTracking`
  Expected: PASS

- [ ] **Step 5: Local Commit**
  Run:
  ```bash
  git add backend/models.go backend/storage.go backend/storage_test.go
  git commit -m "feat(backend): implement playlist CRUD and play frequency tracking in storage"
  ```

---

### Task 2: Recommendation Generation Engine & Wails Bindings

**Files:**
- Modify: `backend/storage.go`
- Modify: `app.go`
- Test: `backend/storage_test.go`

**Interfaces:**
- Consumes: `SearchService`, `StorageService`
- Produces:
  - `App.CreatePlaylist(name, description string) (*backend.Playlist, error)`
  - `App.UpdatePlaylist(id, name, description string) (*backend.Playlist, error)`
  - `App.DeletePlaylist(id string) error`
  - `App.AddTrackToPlaylist(playlistID string, track backend.SearchResultItem) error`
  - `App.RemoveTrackFromPlaylist(playlistID string, trackID string) error`
  - `App.RecordPlay(track backend.SearchResultItem) error`
  - `App.GetHomeData() (map[string][]backend.SearchResultItem, error)`
  - `App.ExportPlaylistJSON(playlistID string) (string, error)`
  - `App.ImportPlaylistJSON(jsonPayload string) (*backend.Playlist, error)`

- [ ] **Step 1: Write test for recommendation synthesis and playlist import/export**
  Add unit tests in `backend/storage_test.go` testing:
  - Exporting a playlist to JSON and re-importing it.
  - Generating top artist discovery seeds from play stats and history.

- [ ] **Step 2: Run test to verify failure**
  Run: `go test -v ./backend -run TestPlaylistExportImport`
  Expected: FAIL

- [ ] **Step 3: Implement Wails bridge methods in `app.go`**
  - Connect all StorageService playlist methods to Wails runtime.
  - Implement `GetHomeData()` returning:
    - `"topPlayed"`: songs with highest play count.
    - `"recommended"`: related songs fetched via artist discovery seeds.
    - `"favorites"`: list of user favorites.

- [ ] **Step 4: Run tests & verify compilation**
  Run: `go test -v ./backend`
  Expected: PASS

- [ ] **Step 5: Local Commit**
  Run:
  ```bash
  git add backend/storage.go backend/storage_test.go app.go
  git commit -m "feat(app): expose playlist management and recommendation engine via Wails"
  ```

---

### Task 3: Frontend Stores & Threshold Play Tracking

**Files:**
- Modify: `frontend/src/lib/types.ts`
- Modify: `frontend/src/lib/stores/userDataStore.ts`
- Modify: `frontend/src/lib/stores/playerStore.ts`
- Modify: `frontend/src/lib/stores/i18nStore.ts`

**Interfaces:**
- Produces:
  - `userDataStore.playlists`: Svelte store array of `Playlist`
  - `userDataStore.createPlaylist(name, desc)`
  - `userDataStore.addTrackToPlaylist(playlistId, track)`
  - `userDataStore.removeTrackFromPlaylist(playlistId, trackId)`
  - `userDataStore.deletePlaylist(playlistId)`
  - `userDataStore.importPlaylist(codeOrJson)`
  - `playerStore`: tracks elapsed seconds and triggers `RecordPlay(track)` when `currentTime >= 20` seconds or on `ended`.

- [ ] **Step 1: Update `types.ts` and `i18nStore.ts`**
  - Add TypeScript interfaces for `Playlist`.
  - Add i18n dictionary keys for Indonesian & English (e.g. `playlist_title`, `create_playlist`, `add_to_playlist`, `share_playlist`, `import_playlist`, `for_you`, `most_played`, `share_copied`).

- [ ] **Step 2: Update `userDataStore.ts`**
  - Add reactive methods wrapping Wails bindings for playlist CRUD and import/export.
  - Ensure local optimistic UI updates followed by backend sync.

- [ ] **Step 3: Implement 20-second play threshold listener in `playerStore.ts`**
  - When a song plays, start tracking elapsed duration. If `currentTime >= 20` seconds (and hasn't been counted for this session), call `RecordPlay(currentTrack)`.
  - Reset counter when changing track.

- [ ] **Step 4: Verify frontend builds cleanly**
  Run: `npm run build` in `frontend/`
  Expected: PASS with 0 TypeScript/Svelte errors.

- [ ] **Step 5: Local Commit**
  Run:
  ```bash
  git add frontend/src/lib/types.ts frontend/src/lib/stores/userDataStore.ts frontend/src/lib/stores/playerStore.ts frontend/src/lib/stores/i18nStore.ts
  git commit -m "feat(store): add playlist state management and 20s playback threshold tracking"
  ```

---

### Task 4: UI Modals (Create/Edit Playlist, Add to Playlist, and Dual Sharing)

**Files:**
- Create: `frontend/src/lib/components/PlaylistModal.svelte`
- Create: `frontend/src/lib/components/AddToPlaylistModal.svelte`
- Create: `frontend/src/lib/components/PlaylistShareModal.svelte`
- Modify: `frontend/src/lib/components/TrackItem.svelte`

**Interfaces:**
- Produces:
  - `PlaylistModal`: Dialog to create/rename playlist (Name, Description).
  - `AddToPlaylistModal`: Popup showing existing playlists with checkmarks to add/remove selected track.
  - `PlaylistShareModal`: Tabbed modal with:
    1. "Salin Kode Playlist" (Base64 URL-safe string to copy).
    2. "Ekspor File (.json)" (Download file).
    3. "Impor Playlist" (Input box to paste code or drag/drop JSON file).
  - `TrackItem.svelte`: Dropdown menu button `⋮` with options: *Tambah ke Antrean*, *Tambah ke Playlist*, *Salin Tautan Lagu*, *Favorit*.

- [ ] **Step 1: Create `PlaylistModal.svelte`**
  - Sleek dark glass modal with input fields, autofocus, validation, and Submit button.

- [ ] **Step 2: Create `AddToPlaylistModal.svelte`**
  - Displays user's playlists with track counts and quick tap-to-add action with confirmation toast.

- [ ] **Step 3: Create `PlaylistShareModal.svelte`**
  - Dual-mode sharing interface with one-click copy and file download/upload handlers.

- [ ] **Step 4: Add Context Action Menu `⋮` to `TrackItem.svelte`**
  - Add sleek 3-dots button on hover, opening contextual popup to trigger *Tambah ke Playlist*.

- [ ] **Step 5: Verify frontend builds cleanly**
  Run: `npm run build` in `frontend/`
  Expected: PASS

- [ ] **Step 6: Local Commit**
  Run:
  ```bash
  git add frontend/src/lib/components/PlaylistModal.svelte frontend/src/lib/components/AddToPlaylistModal.svelte frontend/src/lib/components/PlaylistShareModal.svelte frontend/src/lib/components/TrackItem.svelte
  git commit -m "feat(ui): add playlist creation, song assignment, and dual-mode sharing modals"
  ```

---

### Task 5: Spotify-Style Playlist View & Dynamic 2x2 Auto-Collage Cover Art

**Files:**
- Create: `frontend/src/lib/components/PlaylistView.svelte`
- Modify: `frontend/src/lib/components/Sidebar.svelte`
- Modify: `frontend/src/App.svelte`

**Interfaces:**
- Produces:
  - `PlaylistView.svelte`:
    - Header with 2x2 collage thumbnail (using covers of first 4 songs in playlist).
    - Title, Description, Total Duration, Track Count.
    - Action buttons: *Putar Semua* (Hero CTA), *Putar Acak*, *Bagikan* (opens Share Modal), *Ubah*, *Hapus*.
    - Track table with track index, song title, artist, duration, and delete icon.
  - `Sidebar.svelte`:
    - Under *Koleksi Saya*, render list of user's playlists.
    - `+ Buat Playlist` button.
    - Tooltips showing playlist name when sidebar is collapsed.
  - `App.svelte`:
    - Support selecting a playlist from sidebar (`activeTab = { type: 'playlist', id: playlist.id }`).

- [ ] **Step 1: Implement `PlaylistView.svelte`**
  - Render auto-collage 2x2 grid from `tracks.slice(0, 4)`.
  - Handle playback of entire playlist, shuffle play, remove track, edit, and share.

- [ ] **Step 2: Update `Sidebar.svelte`**
  - Render user playlists with smooth expand/collapse transitions and z-index 60 tooltips.

- [ ] **Step 3: Wire into `App.svelte`**
  - Handle navigation between Search, Favorites, History, and Playlist views.

- [ ] **Step 4: Verify frontend builds cleanly**
  Run: `npm run build` in `frontend/`
  Expected: PASS

- [ ] **Step 5: Local Commit**
  Run:
  ```bash
  git add frontend/src/lib/components/PlaylistView.svelte frontend/src/lib/components/Sidebar.svelte frontend/src/App.svelte
  git commit -m "feat(ui): implement dedicated playlist view with 2x2 auto-collage and sidebar integration"
  ```

---

### Task 6: Multi-Section Homepage (Sering Kamu Putar, Rekomendasi Untuk Kamu, Trending)

**Files:**
- Create: `frontend/src/lib/components/HomeSections.svelte`
- Modify: `frontend/src/App.svelte`

**Interfaces:**
- Produces:
  - `HomeSections.svelte`:
    - Section 1: *Sering Kamu Putar* (Horizontal cards / quick pick pills of top played tracks).
    - Section 2: *Rekomendasi Untuk Kamu* (Recommended tracks based on top artists).
    - Section 3: *Trending Hari Ini* (Existing popular track list with filter pills).
  - Appears on the Search/Explore tab when there is no active search query.

- [ ] **Step 1: Create `HomeSections.svelte`**
  - Layout top-played tracks in a responsive grid.
  - Layout personalized recommendation tracks in clean streaming cards with play buttons.
  - Seamlessly embed the trending track list below.

- [ ] **Step 2: Integrate into `App.svelte`**
  - Show `HomeSections` when `searchQuery` is empty; show search results when a query is entered.

- [ ] **Step 3: Test full Wails application build**
  Run: `& "D:\APP\gopath\bin\wails.exe" build -clean`
  Expected: Success, binary `build/bin/dengar.exe` created.

- [ ] **Step 4: Local Commit**
  Run:
  ```bash
  git add frontend/src/lib/components/HomeSections.svelte frontend/src/App.svelte
  git commit -m "feat(ui): implement multi-section personalized homepage with top played and recommendations"
  ```

---

## Verification Plan

### Automated Tests
1. **Backend Storage & Algorithm Tests:**
   ```powershell
   go test -v ./backend
   ```
   *Expected:* All tests in `storage_test.go`, `search_test.go`, and `stream_test.go` pass (100% OK).

2. **Frontend Vite Build Validation:**
   ```powershell
   npm --prefix frontend run build
   ```
   *Expected:* Exit code 0, bundled in < 3s with 0 errors.

3. **Complete Wails Desktop Binary Build:**
   ```powershell
   & "D:\APP\gopath\bin\wails.exe" build -clean
   ```
   *Expected:* Compiles `build/bin/dengar.exe` successfully with 0 errors.

### Manual Verification
1. **Play Tracking Verification:**
   - Play a song for 25 seconds; verify that its play count increments in `%APPDATA%/dengar/userdata.json`.
   - Skip another song after 3 seconds; verify that its play count does NOT increment.
2. **Homepage Multi-Section:**
   - Launch app on clean home view; verify *Sering Kamu Putar*, *Rekomendasi Untuk Kamu*, and *Trending Hari Ini* render in harmonious visual hierarchy.
3. **Playlist Creation & Song Addition:**
   - Click `+ Buat Playlist` in sidebar; create "Lagu Santai Malam".
   - Open a song, click `⋮` $\to$ *Tambah ke Playlist*, select "Lagu Santai Malam".
   - Open the playlist from the sidebar; verify song appears and collage cover renders.
4. **Dual Sharing:**
   - In playlist detail, click *Bagikan*; click *Salin Kode Playlist*.
   - Click *Impor Playlist*, paste the code; verify duplicate playlist is created instantly.
   - Click *Ekspor File (.json)*; verify JSON file downloads with playlist data.
