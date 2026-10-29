<div align="center">

# 🎵 Dengar Desktop

**Aplikasi Pemutar Musik YouTube Modern, Ringan, dan 100% Bebas Iklan**

[![Release](https://img.shields.io/github/v/release/komporkode/dengar?style=flat-square&color=E61E38)](https://github.com/komporkode/dengar/releases)
[![License: GPL-3.0](https://img.shields.io/badge/License-GPL--3.0-blue.svg?style=flat-square)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Wails](https://img.shields.io/badge/Wails-v2-DF1A29?style=flat-square)](https://wails.io)
[![Svelte](https://img.shields.io/badge/Svelte-4-FF3E00?style=flat-square&logo=svelte)](https://svelte.dev/)
[![TailwindCSS](https://img.shields.io/badge/TailwindCSS-v4-38B2AC?style=flat-square&logo=tailwind-css)](https://tailwindcss.com/)

[Fitur](#-fitur-unggulan) • [Unduh](#-unduhan-dan-instalasi) • [Cara Menjalankan](#-panduan-pengembangan) • [Lisensi](#-lisensi)

</div>

---

## ✨ Fitur Unggulan

- 🚫 **Tanpa Iklan & Bebas Pelacakan** — Streaming audio musik jernih tanpa gangguan jeda iklan pihak ketiga.
- 🎧 **Audio Berkualitas Tinggi (Hi-Fi Stream)** — Menggunakan stream Opus/AAC dengan konsumsi RAM dan bandwidth yang sangat efisien.
- 📚 **Sistem Playlist Spotify-Style** — Buat playlist kustom dengan cover art mosaik (*auto-collage 2x2*) dinamis dari 4 lagu pertama di dalamnya.
- 🔗 **Dual-Sharing Playlist P2P** — Bagikan playlist antar-teman secara langsung lewat kode teks Base64 atau ekspor/impor berkas JSON tanpa memerlukan akun cloud terpusat.
- ⚡ **Rekomendasi Cerdas & Cache Lokal**:
  - **Sering Kamu Putar**: 100% dihitung dari riwayat putar lokal di PC pengguna (*zero network latency*).
  - **Rekomendasi Untuk Kamu**: Dilengkapi *in-memory cache* 2 jam dan tombol segarkan manual untuk meminimalkan beban jaringan.
- 🔊 **Sistem Volume Persisten & Robust**:
  - Nilai volume dan status bisu (*mute*) tersimpan otomatis di perangkat.
  - Mendukung kontrol *mouse scroll wheel* pada slider volume.
  - Memiliki *Smart Unmute Memory* yang otomatis memulihkan level suara sebelumnya saat dinyalakan dari posisi 0%.
- 🪟 **Floating Mini Player**:
  - Mode ringkas mengambang dengan tombol *Pin (Always on Top)*.
  - Mendukung kontrol scrubber cepat, visualizer animasi, dan shortcut keyboard.
- 🔄 **Pemeriksa Pembaruan Terintegrasi** — Tombol *Cek Pembaruan* langsung di menu Pengaturan yang terhubung ke rilis GitHub.

---

## 📦 Unduhan dan Instalasi

### Windows (x64)
1. Unduh rilis portabel terbaru di [GitHub Releases](https://github.com/komporkode/dengar/releases/latest).
2. Ekstrak file zip `dengar-windows-amd64.zip` (atau sesuai rilis versi terbaru).
3. Jalankan `dengar.exe` langsung (portabel, tanpa perlu instalasi rumit).

### Linux & macOS
Aplikasi dapat di-build langsung dari source code dengan Go dan Wails CLI:
```bash
# Clone repositori
git clone https://github.com/komporkode/dengar.git
cd dengar

# Build native binary
wails build -clean
```

---

## 🛠️ Panduan Pengembangan

### Kebutuhan Sistem:
- [Go 1.23+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/) & npm
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation):
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

### Menjalankan Mode Pengembangan (Live Reload):
```bash
# Di direktori root proyek:
wails dev
```
Perubahan kode pada frontend Svelte dan backend Go akan langsung dimuat ulang secara otomatis (*hot-reload*).

---

## 🏛️ Arsitektur Teknologi

- **Backend:** Go 1.23, Wails v2.17 Runtime
- **Frontend:** Svelte 4, TypeScript, Tailwind CSS v4, Vite
- **Penyimpanan:** Atomic local storage di `%APPDATA%/dengar/userdata.json` (Windows) / `~/.config/dengar/userdata.json` (Linux)
- **Komunikasi:** IPC Wails Binding terkompilasi aman

---

## 📄 Lisensi

Proyek ini dilisensikan di bawah **GNU General Public License v3.0 (GPL-3.0)**. 
Silakan lihat berkas [LICENSE](LICENSE) untuk informasi hak cipta dan ketentuan redistribusi.

Copyright (C) 2026 **Komporkode**.
