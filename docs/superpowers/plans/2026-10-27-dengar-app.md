# Dengar Desktop App Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Membangun aplikasi desktop pemutar YouTube & YouTube Music tanpa iklan, ringan, personal-use, berbasis Go 1.23 + Wails v2 + Svelte & TypeScript dengan biaya operasional Rp0 mengandalkan failover instance publik Piped & Invidious.

**Architecture:** Arsitektur 2-layer ringkas (Go Backend Service + Svelte Frontend di WebView2). Backend mengelola curated instance pool, HTTP client dengan auto-failover, in-memory stream cache (TTL 30 menit), dan atomic JSON persistence di `os.UserConfigDir()`. Frontend menangani antarmuka 3-panel bergaya YouTube Music, HTML5 `<audio>` & `<video>` player dengan auto re-resolve stream saat error, dan antrean lagu responsif.

**Tech Stack:** Go 1.23+, Wails v2.17+, Svelte, TypeScript, HTML5 Audio/Video API, `net/http` bawaan Go, CSS modern (dark theme).

**Spec:** `d:\PROJECT\GAWEAN\PROJECT\APP\Dengar\denngar_prd.md`

## Global Constraints

- Biaya operasional target: Rp0 (tanpa server berbayar, tanpa cloud DB).
- Tanpa dependensi eksternal berat (tanpa FFmpeg binary bundled; gunakan stream native m4a dan muxed mp4).
- Bebas iklan YouTube dalam bentuk audio, video, maupun visual pop-up.
- Penyimpanan lokal user data (Riwayat & Favorit) disimpan di file JSON atomik pada `os.UserConfigDir()/dengar/userdata.json`.
- Kode backend mematuhi prinsip `/ponytail`: solusi paling ringkas, tanpa wrapper berlebihan, standard library first, direct type mapping.

## Review Focus

1. **Instance Rate-Limiting & HTTP 429/403:** Instance publik sering memblokir IP atau rate-limiting. Sistem harus otomatis mencoba instance berikutnya dalam pool tanpa membuat UI membeku.
2. **Stream URL Expiration:** URL stream YouTube memiliki batas waktu token expire. Player harus menangani event `error` pada `<audio>`/`<video>` secara otomatis dengan me-resolve URL baru dan resume di posisi `currentTime`.
3. **Data Corruption Resistance:** Jika aplikasi ditutup mendadak saat menulis file riwayat/favorit, file tidak boleh korup (wajib atomic write via temporary file + rename).
4. **Offline / Network Disconnection:** Jika koneksi internet putus, aplikasi menampilkan banner ramah dan tidak mengalami fatal crash / panic.
5. **Memory Leak / Zombie Streams:** Pergantian lagu beruntun tidak boleh meninggalkan stream dangling atau audio ganda yang berbunyi bersamaan.

---

### Task 1: Scaffolding Project Wails v2 + Svelte-TS

**Files:**
- Create: `wails.json`, `main.go`, `app.go`, `go.mod`
- Create: `frontend/package.json`, `frontend/src/App.svelte`, `frontend/src/main.ts`
- Test: Build scaffolding verification

**Interfaces:**
- Produces: Project skeleton dengan komunikasi IPC Wails berfungsi antara Go dan Svelte.

- [ ] **Step 1: Inisialisasi template Wails svelte-ts di workspace**
  Jalankan `wails init -n dengar -t svelte-ts` atau konfigurasi file project secara rapi sesuai struktur Wails v2.
- [ ] **Step 2: Verifikasi go.mod dan dependensi frontend**
  Pastikan `go mod tidy` dan `npm install` di folder `frontend` sukses tanpa error.
- [ ] **Step 3: Test compile backend Go**
  Run: `go build -o dengar_test.exe .`
  Expected: Kompilasi sukses tanpa error.
- [ ] **Step 4: Cleanup binary uji coba**
  Hapus `dengar_test.exe`.
- [ ] **Step 5: Commit**
  `git add . && git commit -m "chore: scaffold wails v2 svelte-ts project"`

---

### Task 2: Backend Instance Pool & Search Service

**Files:**
- Create: `backend/instances.go`
- Create: `backend/search.go`
- Create: `backend/models.go`
- Test: `backend/search_test.go`

**Interfaces:**
- Produces:
  - `type SearchResultItem struct { ID, Title, UploaderName, DurationFormatted, ThumbnailURL string }`
  - `func (s *SearchService) Search(query string, filterType string) ([]SearchResultItem, error)`
  - Pool rotasi instance Piped (`pipedapi.kavin.rocks`, `api.piped.privacydev.net`, dll.) dan Invidious (`inv.tux.pizza`, dll.).

- [ ] **Step 1: Tulis unit test untuk fallback pencarian**
  Buat `backend/search_test.go` yang menguji parser JSON respons Piped & Invidious dengan dummy HTTP response dan mekanisme failover ke instance kedua jika instance pertama mengembalikan status 500.
- [ ] **Step 2: Jalankan test untuk memastikan failing**
  Run: `go test ./backend/...`
  Expected: FAIL (belum ada implementasi).
- [ ] **Step 3: Implementasikan models.go, instances.go, dan search.go**
  Gunakan `net/http` dengan timeout 5 detik, rotasi pool instance otomatis, dan parsing hasil pencarian ke `[]SearchResultItem`.
- [ ] **Step 4: Jalankan test untuk memastikan passing**
  Run: `go test ./backend/... -v`
  Expected: PASS
- [ ] **Step 5: Commit**
  `git add backend/ && git commit -m "feat(backend): add instance pool and search service with failover"`

---

### Task 3: Backend Stream Resolver & In-Memory TTL Cache

**Files:**
- Create: `backend/stream.go`
- Create: `backend/cache.go`
- Test: `backend/stream_test.go`

**Interfaces:**
- Produces:
  - `type StreamResult struct { VideoID, AudioURL, VideoURL, Title, Artist, Thumbnail string, Duration int64 }`
  - `func (r *StreamResolver) Resolve(videoID string, forceRefresh bool) (*StreamResult, error)`
  - TTL Cache in-memory 30 menit thread-safe menggunakan `sync.RWMutex`.

- [ ] **Step 1: Tulis test untuk TTL Cache & Stream Resolver**
  Buat `backend/stream_test.go` menguji penyimpanan cache dengan expired duration dan pemanggilan ulang jika `forceRefresh = true`.
- [ ] **Step 2: Jalankan test untuk memastikan failing**
  Run: `go test -run TestStreamCache ./backend/...`
  Expected: FAIL
- [ ] **Step 3: Implementasikan cache.go dan stream.go**
  Ekstrak stream audio m4a (prioritas bitrate terbaik) dan muxed mp4 video (720p/360p) dari respons endpoint `/streams/:id` (Piped) atau `/api/v1/videos/:id` (Invidious).
- [ ] **Step 4: Jalankan test untuk memastikan passing**
  Run: `go test -run TestStreamCache ./backend/... -v`
  Expected: PASS
- [ ] **Step 5: Commit**
  `git add backend/ && git commit -m "feat(backend): add stream resolver and 30m TTL cache"`

---

### Task 4: Backend Atomic Local Storage (UserData: History & Favorites)

**Files:**
- Create: `backend/storage.go`
- Test: `backend/storage_test.go`

**Interfaces:**
- Produces:
  - `type UserData struct { Favorites []SearchResultItem; History []SearchResultItem }`
  - `func (s *StorageService) Load() (*UserData, error)`
  - `func (s *StorageService) AddFavorite(item SearchResultItem) error`
  - `func (s *StorageService) RemoveFavorite(id string) error`
  - `func (s *StorageService) AddHistory(item SearchResultItem) error` (capped 100 entri)
  - Lokasi: `%APPDATA%\dengar\userdata.json` dengan penulisan atomic write (`userdata.json.tmp` -> rename).

- [ ] **Step 1: Tulis test untuk atomic write dan pembatasan 100 riwayat**
  Buat `backend/storage_test.go` menguji penambahan riwayat hingga 110 item (memastikan tersisa 100 item terbaru) dan ketahanan thread-safety.
- [ ] **Step 2: Jalankan test untuk verifikasi fail**
  Run: `go test -run TestStorage ./backend/...`
  Expected: FAIL
- [ ] **Step 3: Implementasikan storage.go**
  Menggunakan `sync.RWMutex`, `os.UserConfigDir()`, dan `os.WriteFile` temp + `os.Rename`.
- [ ] **Step 4: Jalankan test untuk verifikasi pass**
  Run: `go test -run TestStorage ./backend/... -v`
  Expected: PASS
- [ ] **Step 5: Commit**
  `git add backend/storage* && git commit -m "feat(backend): add atomic local storage for history and favorites"`

---

### Task 5: Wails Application Binding & IPC Integration

**Files:**
- Modify: `app.go`
- Modify: `main.go`

**Interfaces:**
- Produces: Wails method bindings yang diekspos ke Svelte frontend:
  - `Search(query string, filterType string) ([]SearchResultItem, error)`
  - `GetStream(videoID string, forceRefresh bool) (*StreamResult, error)`
  - `GetUserData() (*UserData, error)`
  - `ToggleFavorite(item SearchResultItem) (bool, error)`
  - `AddToHistory(item SearchResultItem) error`

- [ ] **Step 1: Hubungkan service backend ke struct `App` di `app.go`**
  Inisialisasi `SearchService`, `StreamResolver`, dan `StorageService` pada fungsi `startup(ctx context.Context)`.
- [ ] **Step 2: Ekspos method-method Wails pada struct `App`**
- [ ] **Step 3: Verifikasi wails build dan generate runtime bindings**
  Run: `wails build` atau `go build .`
  Expected: Berhasil dikompilasi tanpa type conflict.
- [ ] **Step 4: Commit**
  `git add app.go main.go && git commit -m "feat: expose backend services to wails runtime"`

---

### Task 6: Frontend Design System & Theme Styling

**Files:**
- Create: `frontend/src/app.css`
- Modify: `frontend/src/App.svelte`

**Interfaces:**
- Produces: YouTube Music aesthetic tokens:
  - Warna: Background gelap murni `#030303`, surface `#121212`, hover `#212121`, aksen merah `#ff0000`, teks `#ffffff` & `#aaaaaa`.
  - Glassmorphic bottom player bar dengan blur effect.
  - Komponen slider seek dan volume kustom yang mulus.

- [ ] **Step 1: Tulis variabel CSS dan utilitas layout di `frontend/src/app.css`**
- [ ] **Step 2: Verifikasi styling di browser / vite preview**
  Run: `npm --prefix frontend run build`
  Expected: Build sukses tanpa error syntax CSS.
- [ ] **Step 3: Commit**
  `git add frontend/src/app.css && git commit -m "style(frontend): implement yt-music dark theme design system"`

---

### Task 7: Frontend State Management (Svelte Stores)

**Files:**
- Create: `frontend/src/lib/stores/playerStore.ts`
- Create: `frontend/src/lib/stores/queueStore.ts`
- Create: `frontend/src/lib/stores/userDataStore.ts`
- Create: `frontend/src/lib/types.ts`

**Interfaces:**
- Produces:
  - `playerStore`: status pemutaran (`isPlaying`, `currentTime`, `duration`, `volume`, `isMuted`, `mode: 'audio' | 'video'`, `currentTrack`, `isLoading`).
  - `queueStore`: daftar lagu antrean, index aktif, `repeatMode: 'off' | 'all' | 'one'`, `isShuffle`.
  - `userDataStore`: data sinkron riwayat dan favorit dari backend Wails.

- [ ] **Step 1: Implementasikan types.ts dan store reactive Svelte**
- [ ] **Step 2: Tambahkan logika antrean (Next, Prev, Shuffle array, Repeat toggle)**
- [ ] **Step 3: Verifikasi build frontend**
  Run: `npm --prefix frontend run build`
  Expected: Svelte TypeScript compiler valid tanpa error tipe.
- [ ] **Step 4: Commit**
  `git add frontend/src/lib/ && git commit -m "feat(frontend): create reactive player, queue, and user stores"`

---

### Task 8: UI Components - Sidebar, Search, Filter Tabs & Track List

**Files:**
- Create: `frontend/src/lib/components/Sidebar.svelte`
- Create: `frontend/src/lib/components/SearchBar.svelte`
- Create: `frontend/src/lib/components/TrackList.svelte`
- Create: `frontend/src/lib/components/TrackItem.svelte`

**Interfaces:**
- Produces: Komponen tampilan navigasi, input pencarian dengan debounce, filter chips (Lagu, Video, Semua), dan rendering list lagu dengan cover, judul, channel, durasi, serta tombol aksi cepat (Play, Antrean, Favorit).

- [ ] **Step 1: Implementasikan Sidebar dengan navigasi (Pencarian, Favorit, Riwayat)**
- [ ] **Step 2: Implementasikan SearchBar dengan debounce 400ms dan filter chips**
- [ ] **Step 3: Implementasikan TrackItem & TrackList dengan hover state dan action buttons**
- [ ] **Step 4: Uji integrasi tampilan di App.svelte**
  Run: `npm --prefix frontend run build`
  Expected: PASS
- [ ] **Step 5: Commit**
  `git add frontend/src/lib/components/ && git commit -m "feat(frontend): add sidebar, search bar, and track list components"`

---

### Task 9: Persistent Player Bar & Resilient Audio/Video Player

**Files:**
- Create: `frontend/src/lib/components/PlayerBar.svelte`
- Create: `frontend/src/lib/components/VideoModal.svelte`
- Modify: `frontend/src/App.svelte`

**Interfaces:**
- Produces:
  - `<audio>` & `<video>` element binding ke `playerStore`.
  - Auto-advance ke lagu berikutnya saat track `ended`.
  - **Auto Re-resolve Handler:** tangkap event `error` pada media element, panggil `GetStream(id, forceRefresh=true)`, pasang stream baru dan set `currentTime` ke posisi sebelum error.
  - Video switch: toggle antara Audio Mode dan Video Mode tanpa me-reset queue.

- [ ] **Step 1: Implementasikan PlayerBar dengan seek bar, play/pause, prev/next, volume, dan mode toggle**
- [ ] **Step 2: Implementasikan VideoModal / Video View yang tampil saat mode video aktif**
- [ ] **Step 3: Pasang event handler onerror untuk auto re-resolve stream failover**
- [ ] **Step 4: Verifikasi build frontend**
  Run: `npm --prefix frontend run build`
  Expected: PASS
- [ ] **Step 5: Commit**
  `git add frontend/src/ && git commit -m "feat(frontend): implement persistent player bar with resilient stream recovery"`

---

### Task 10: Queue Drawer & In-App Keyboard Shortcuts

**Files:**
- Create: `frontend/src/lib/components/QueueDrawer.svelte`
- Modify: `frontend/src/App.svelte`

**Interfaces:**
- Produces:
  - Panel drawer antrean dari kanan (lihat lagu saat ini, berikutnya, hapus item, reorder).
  - Window listener untuk keyboard shortcuts: Spasi (Play/Pause), Panah Kiri/Kanan (Seek ±5 detik), Panah Atas/Bawah (Volume ±5%), 'M' (Mute), 'F' (Toggle Favorit).

- [ ] **Step 1: Implementasikan QueueDrawer dengan visual track aktif**
- [ ] **Step 2: Implementasikan global keydown handler di App.svelte dengan guard saat mengetik di input search**
- [ ] **Step 3: Verifikasi build frontend**
  Run: `npm --prefix frontend run build`
  Expected: PASS
- [ ] **Step 4: Commit**
  `git add frontend/src/ && git commit -m "feat(frontend): add queue drawer and in-app keyboard shortcuts"`

---

### Task 11: End-to-End Build, Production Packaging & Verification

**Files:**
- Modify: `wails.json`
- Test: Full build Wails executable

- [ ] **Step 1: Jalankan seluruh Go unit test**
  Run: `go test -v ./...`
  Expected: Seluruh test backend PASS.
- [ ] **Step 2: Build aplikasi desktop production dengan Wails**
  Run: `wails build`
  Expected: Menghasilkan binary desktop executable di `build/bin/dengar.exe`.
- [ ] **Step 3: Verifikasi binary dan ukuran file**
  Pastikan binary siap dijalankan secara mandiri.
- [ ] **Step 4: Commit final**
  `git add . && git commit -m "chore: complete production build setup"`
