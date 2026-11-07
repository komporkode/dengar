# Changelog

Semua perubahan penting pada proyek **Dengerin Desktop** akan dicatat di berkas ini.

Format berkas ini mengacu pada [Keep a Changelog](https://keepachangelog.com/id/1.0.0/),
dan proyek ini mematuhi [Semantic Versioning (SemVer)](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.0] - 2026-11-07

### Added
- **Penyempurnaan Sistem Auto-Update & Installer:**
  - Progress bar unduhan pembaruan real-time lengkap dengan persentase dan kecepatan unduhan (MB/s).
  - Pemasangan pembaruan otomatis & restart instan (*in-place silent upgrade*) untuk installer `.exe` dan versi portabel `.zip`.
  - Installer NSIS kini menyertakan opsi *Jalankan Dengerin* di halaman akhir untuk kemudahan pengguna.
  - Deteksi akurat antara edisi Terpasang (Installer) dan Portabel via keberadaan `uninstall.exe`.
  - Pratinjau catatan rilis (*changelog*) langsung di halaman Pengaturan.
- **Craftsmanship & Aksesibilitas Frontend (Anti-Slop):**
  - Dukungan tombol keyboard `Escape` untuk menutup seluruh laci antrean dan jendela dialog/modal.
  - Cakupan penuh atribut `aria-label` dan `:focus-visible` untuk navigasi pembaca layar dan keyboard.
  - Dukungan preferensi aksesibilitas `@media (prefers-reduced-motion: reduce)`.
  - Eliminasi ornamen visual berlebih (radial blur orb) dan buzzword pemasaran buatan AI.
- **Spotify-Style Player Empty State:** PlayerBar dan MiniPlayer menampilkan status kosong yang elegan saat belum ada lagu yang dipilih (tombol kendali dinonaktifkan redup, volume & equalizer tetap aktif, tanpa layout shift atau spinner macet).
- **Redesign Minimalis Queue Drawer:**
  - Header antrean ringkas 1 baris (*single-tier*, 56px) menggantikan header bertingkat yang padat tanpa ikon hamburger besar.
  - Pembagian visual antrean terstruktur: *Sedang Diputar* (kartu crimson glowing dengan animasi equalizer), *Berikutnya*, dan *Sebelumnya* lengkap dengan hitungan lagu.
  - Aksi baru *Hapus Sisa* (Clear Upcoming) untuk mengosongkan lagu berikutnya tanpa menghapus lagu yang sedang diputar.
  - Toggle *Autoplay Serupa* kompak terintegrasi di toolbar kanan.
- **Lokalisasi Penuh Bilingual (ID & EN):** Dukungan terjemahan 100% di seluruh komponen antrean, pencarian, menu opsi lagu, modal playlist, dan pengaturan.

### Changed
- Navigasi sidebar kini menggunakan penanda visual aktif yang elegan (`bg-white/[0.1]` dan ikon aksen merah) tanpa teks label "Aktif".
- Toleransi watchdog pemuatan stream dinaikkan ke 15 detik dan stall timer ke 12 detik dengan proteksi status buffer aktif.

### Fixed
- **Error Retry Pemutaran Pertama (Cold Start):**
  - Mengeliminasi konflik unduhan audio ganda di latar belakang yang memicu perebutan bandwidth dan *file lock* temporary (`.tmp`).
  - Mengabaikan `AbortError` dan `NotAllowedError` siklus hidup browser saat pergantian track agar tidak keliru didiagnosis sebagai kegagalan server.
  - Mengaktifkan AudioContext secara sinkron tepat saat gestur klik pengguna untuk mencegah batasan pemblokiran audio browser.
- Timeout HTTP Range probe upstream ditingkatkan ke 2.5 detik untuk mencegah *false-alarm failover* akibat latensi jaringan internasional.

## [1.3.2] - 2026-11-03

### Fixed
- Infinite scroll pencarian yang terhenti di 11 lagu akibat token kontinuitas yang hilang.
- Indikator titik merah yang menghalangi ikon pada sidebar collapsed.
- Persistensi kuota kapasitas disk cache (LRU) di Pengaturan.
- Build NSIS Windows installer pada GitHub Actions CI/CD.

## [1.3.1] - 2026-11-03

### Added
- Pencarian utama langsung dari katalog resmi YouTube Music (`WEB_REMIX`).
- Penyaringan cerdas konten non-musik (podcast, gameplay, tutorial, vlog).
- Dukungan parsing durasi format regional dengan tanda titik maupun titik dua.

## [1.3.0] - 2026-11-03

### Added
- Kustomisasi lokasi folder penyimpanan cache audio ke drive lain dengan migrasi otomatis.
- Installer Windows fleksibel dengan pemilihan direktori instalasi dan deteksi upgrade otomatis.
- Mode kualitas audio hemat kuota (*Eco*) sebagai default.

