# PRD — Aplikasi Desktop Pemutar YouTube / YouTube Music Tanpa Iklan (Personal Use)

**Versi:** 1.0
**Status:** Draft
**Target Pengguna:** 1 orang (pemakaian pribadi)
**Biaya Operasional Target:** Rp0 (tanpa server berbayar)

---

## 1. Latar Belakang & Tujuan

### 1.1 Masalah
Menonton YouTube / YouTube Music di browser resmi selalu disertai iklan, rekomendasi yang mengganggu, dan konsumsi resource browser yang berat. Pengguna ingin pengalaman menonton/mendengarkan yang bersih, ringan, dan tanpa iklan.

### 1.2 Tujuan
Membangun aplikasi desktop lintas platform (Windows/Linux/macOS) berbasis **Go + Wails 2** yang:
- Memutar video/audio YouTube **tanpa iklan**.
- Memiliki antarmuka mirip YouTube Music (pencarian, daftar putar, pemutar).
- **Tanpa biaya operasional** — menggunakan instance publik Invidious & Piped dengan mekanisme fallback timbal balik.
- Cukup andal untuk pemakaian pribadi sehari-hari.

### 1.3 Non-Tujuan (Out of Scope)
- Tidak untuk distribusi publik / komersial.
- Tidak ada fitur login akun Google, sinkronisasi langganan, atau riwayat cloud.
- Tidak ada unduhan massal / batch download (bisa ditambahkan nanti sebagai stretch goal).
- Tidak menjamin kualitas stream maksimal (tergantung instance publik).

---

## 2. Target Pengguna & Skenario

### 2.1 Persona
Pengguna tunggal (pemilik aplikasi) yang nyaman dengan build dari source, menggunakan PC desktop, ingin mendengarkan musik / menonton video tanpa iklan dan tanpa berlangganan.

### 2.2 Skenario Utama
1. **Cari & Putar Musik**: Buka aplikasi → ketik judul lagu → pilih hasil → audio langsung diputar.
2. **Putar Video**: Sama seperti musik, tetapi dengan tampilan `<video>`.
3. **Daftar Putar Lokal**: Tambahkan lagu ke queue, putar berurutan.
4. **Fallback Otomatis**: Jika instance utama gagal, aplikasi diam-diam beralih ke instance lain tanpa mengganggu pengguna.
5. **Pemutaran Ulang Cepat**: Video yang pernah diputar dimuat dari cache (tanpa request API baru).

---

## 3. Fitur (Functional Requirements)

### 3.1 Fitur Wajib (MVP)
| ID | Fitur | Deskripsi | Prioritas |
|----|-------|-----------|-----------|
| F-01 | Pencarian | Cari video/lagu via API Piped (utama) → Invidious (fallback) | Wajib |
| F-02 | Daftar Hasil | Tampilkan thumbnail, judul, channel/artis, durasi | Wajib |
| F-03 | Pemutar Audio | Putar stream audio (m4a/webm) via `<audio>` | Wajib |
| F-04 | Pemutar Video | Putar stream video via `<video>` | Wajib |
| F-05 | Fallback Instance | Coba instance berikutnya otomatis jika gagal | Wajib |
| F-06 | Cache Stream | Simpan URL stream per `videoID` di memori (TTL 30 menit) | Wajib |
| F-07 | Queue / Playlist | Tambah/hapus/urutkan lagu dalam antrean | Wajib |
| F-08 | Kontrol Pemutar | Play/pause, seek, volume, next/prev | Wajib |
| F-09 | Mode Gelap | Tema gelap default (mirip YT Music) | Wajib |

### 3.2 Fitur Opsional (Nice to Have)
| ID | Fitur | Deskripsi |
|----|-------|-----------|
| F-10 | Riwayat Lokal | Simpan riwayat putar di file JSON lokal |
| F-11 | Favorit | Tandai lagu favorit, simpan lokal |
| F-12 | Unduh Audio | Simpan stream ke file lokal (m4a) |
| F-13 | Pencarian Playlist | Cari playlist YouTube dan putar seluruhnya |
| F-14 | Global Hotkey | Media key (play/pause/next) di OS |

### 3.3 Kriteria Penerimaan (Acceptance Criteria)
- [ ] Pencarian mengembalikan hasil dalam < 5 detik pada koneksi normal.
- [ ] Jika instance pertama gagal, hasil tetap muncul dari instance kedua.
- [ ] Audio mulai diputar dalam < 3 detik setelah klik hasil.
- [ ] Tidak ada iklan yang muncul dalam bentuk apa pun.
- [ ] Aplikasi tidak crash jika semua instance gagal (tampilkan pesan error ramah).

---

## 4. Arsitektur Teknis

### 4.1 Stack
| Layer | Teknologi |
|-------|-----------|
| Backend | Go 1.21+ |
| Desktop Framework | Wails 2 |
| Frontend | Svelte (atau React/Vue) + TypeScript |
| Styling | TailwindCSS (opsional) |
| HTTP Client | `net/http` bawaan Go |
| Penyimpanan lokal | File JSON di `os.UserConfigDir()` |

### 4.2 Diagram Alur (Konseptual)
