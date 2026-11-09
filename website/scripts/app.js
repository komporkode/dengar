const REPO_OWNER = "komporkode";
const REPO_NAME = "dengar";
const GITHUB_REPO_URL = `https://github.com/${REPO_OWNER}/${REPO_NAME}`;
const GITHUB_RELEASES_URL = `${GITHUB_REPO_URL}/releases`;
const GITHUB_LATEST_PAGE_URL = `${GITHUB_RELEASES_URL}/latest`;
const GITHUB_API_LATEST_URL = `https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest`;
const ALLOWED_RELEASE_PREFIX = `https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/`;

const FALLBACK_RELEASE = {
  tag_name: "v1.4.1",
  published_at: "2026-11-08T00:00:00Z",
  html_url: `${GITHUB_RELEASES_URL}/tag/v1.4.1`,
  assets: [
    {
      name: "dengerin-windows-amd64-installer.exe",
      size: null,
      browser_download_url: `${GITHUB_RELEASES_URL}/download/v1.4.1/dengerin-windows-amd64-installer.exe`
    },
    {
      name: "dengerin-windows-amd64-portable.zip",
      size: null,
      browser_download_url: `${GITHUB_RELEASES_URL}/download/v1.4.1/dengerin-windows-amd64-portable.zip`
    },
    {
      name: "dengerin-linux-amd64.tar.gz",
      size: null,
      browser_download_url: `${GITHUB_RELEASES_URL}/download/v1.4.1/dengerin-linux-amd64.tar.gz`
    },
    {
      name: "dengerin-macos-universal.zip",
      size: null,
      browser_download_url: `${GITHUB_RELEASES_URL}/download/v1.4.1/dengerin-macos-universal.zip`
    }
  ]
};

const PLATFORM_DEFINITIONS = {
  windows: {
    labelKey: "os.windows",
    artifacts: [
      {
        match: "installer.exe",
        titleKey: "art.winInstaller.title",
        descKey: "art.winInstaller.desc",
        recommended: true
      },
      {
        match: "portable.zip",
        titleKey: "art.winPortable.title",
        descKey: "art.winPortable.desc",
        recommended: false
      }
    ],
    notesTitleKey: "notes.win.title",
    notesKeys: ["notes.win.1", "notes.win.2", "notes.win.3"]
  },
  linux: {
    labelKey: "os.linux",
    artifacts: [
      {
        match: "linux-amd64.tar.gz",
        titleKey: "art.linuxBin.title",
        descKey: "art.linuxBin.desc",
        recommended: true
      }
    ],
    notesTitleKey: "notes.linux.title",
    notesKeys: ["notes.linux.1", "notes.linux.2", "notes.linux.3"]
  },
  macos: {
    labelKey: "os.macos",
    artifacts: [
      {
        match: "macos-universal.zip",
        titleKey: "art.macBin.title",
        descKey: "art.macBin.desc",
        recommended: true
      }
    ],
    notesTitleKey: "notes.mac.title",
    notesKeys: ["notes.mac.1", "notes.mac.2", "notes.mac.3"]
  }
};

const SCREENSHOTS = [
  {
    id: "explore",
    src: "/assets/screenshots/screen-explore.webp",
    width: 1024,
    height: 682,
    titleKey: "shot.s1.title",
    descKey: "shot.s1.desc",
    altKey: "shot.s1.alt"
  },
  {
    id: "foryou",
    src: "/assets/screenshots/screen-foryou.webp",
    width: 1024,
    height: 682,
    titleKey: "shot.s2.title",
    descKey: "shot.s2.desc",
    altKey: "shot.s2.alt"
  },
  {
    id: "equalizer",
    src: "/assets/screenshots/screen-equalizer.webp",
    width: 1024,
    height: 682,
    titleKey: "shot.s3.title",
    descKey: "shot.s3.desc",
    altKey: "shot.s3.alt"
  },
  {
    id: "sleeptimer",
    src: "/assets/screenshots/screen-sleeptimer.webp",
    width: 1024,
    height: 682,
    titleKey: "shot.s4.title",
    descKey: "shot.s4.desc",
    altKey: "shot.s4.alt"
  },
  {
    id: "miniplayer",
    src: "/assets/screenshots/screen-miniplayer.webp",
    width: 420,
    height: 175,
    titleKey: "shot.s5.title",
    descKey: "shot.s5.desc",
    altKey: "shot.s5.alt"
  }
];

const I18N = {
  id: {
    "meta.title": "Dengerin Desktop - Pemutar Musik YouTube Bebas Iklan",
    "meta.desc": "Aplikasi desktop Go dan Wails ringan untuk memutar katalog YouTube Music bebas iklan dengan cache lokal offline, equalizer 5-band, dan playlist P2P.",
    "skip": "Lewati ke konten utama",
    "nav.download": "Unduh",
    "nav.screens": "Tampilan",
    "nav.features": "Fitur",
    "nav.releases": "Rilis",
    "nav.specs": "Spesifikasi",
    "nav.github": "GitHub",
    "nav.primaryAria": "Utama",
    "nav.langAria": "Pilihan bahasa",
    "nav.drawerAria": "Menu seluler",
    "nav.menuOpen": "Buka menu navigasi",
    "nav.menuClose": "Tutup menu",
    "hero.title": "Putar musik YouTube bebas iklan, langsung dari desktop.",
    "hero.lead": "Dibangun dengan Go dan Wails. Putar katalog YouTube Music dengan cache lokal, equalizer 5-band, dan update otomatis.",
    "hero.ctaPrimary": "Unduh Dengerin",
    "hero.ctaSecondary": "Lihat Kode di GitHub",
    "hero.zoomAria": "Perbesar tangkapan layar utama Dengerin Desktop",
    "hero.zoomHint": "Klik gambar untuk melihat ukuran penuh",
    "hero.mediaAlt": "Tangkapan layar asli Dengerin Desktop menampilkan halaman Eksplorasi & Cari beserta pemutar audio mengambang di bawah.",
    "shot.kicker": "Antarmuka Asli",
    "shot.title": "Cuplikan layar aplikasi desktop",
    "shot.desc": "Tangkapan layar asli dari jendela Dengerin Desktop di Windows. Klik gambar mana saja untuk melihat ukuran penuh.",
    "shot.expandLabel": "Perbesar",
    "shot.miniKicker": "Mode Ringkas",
    "shot.zoomBtn": "Perbesar Mini Player",
    "shot.prevAria": "Gambar sebelumnya",
    "shot.nextAria": "Gambar berikutnya",
    "shot.closeAria": "Tutup pratinjau",
    "shot.s1.title": "Eksplorasi & Cari",
    "shot.s1.desc": "Pencarian langsung YouTube Music, lagu sering diputar, dan bar audio mengambang.",
    "shot.s1.aria": "Perbesar tangkapan layar Eksplorasi & Cari",
    "shot.s1.alt": "Tangkapan layar halaman Eksplorasi & Cari dengan pencarian musik, daftar sering diputar, dan bilah pemutar Eco Audio.",
    "shot.s2.title": "Rekomendasi Untuk Anda",
    "shot.s2.desc": "Mix personal, kurasi pilihan studio, serta kontrol Putar Semua dan Putar Acak.",
    "shot.s2.aria": "Perbesar tangkapan layar Untuk Anda",
    "shot.s2.alt": "Tangkapan layar tab Untuk Anda menampilkan Personal Mix, Koleksi Populer & Pilihan Studio, serta Sering Kamu Putar.",
    "shot.s3.title": "Equalizer Audio 5-Band",
    "shot.s3.desc": "Filter 60 Hz hingga 14 kHz secara lokal dengan preset cepat seperti Vocal Booster.",
    "shot.s3.aria": "Perbesar tangkapan layar Equalizer Audio 5-Band",
    "shot.s3.alt": "Tangkapan layar dialog Equalizer Audio 5-Band menampilkan preset suara dan slider 60 Hz hingga 14 kHz.",
    "shot.s4.title": "Sleep Timer",
    "shot.s4.desc": "Timer cepat 15 hingga 60 menit, berhenti di akhir lagu, dan fade-out halus 5 detik.",
    "shot.s4.aria": "Perbesar tangkapan layar Sleep Timer",
    "shot.s4.alt": "Tangkapan layar dialog Sleep Timer dengan tombol durasi cepat 15 sampai 60 menit dan opsi berhenti di akhir lagu.",
    "shot.s5.title": "Mini Player Mengambang (Always on Top)",
    "shot.s5.desc": "Jendela kecil yang dapat disematkan di atas aplikasi lain untuk kontrol jeda, lagu berikutnya, progres putar, dan volume cepat.",
    "shot.s5.aria": "Perbesar tangkapan layar Mini Player Mengambang",
    "shot.s5.alt": "Tangkapan layar jendela Mini Player Dengerin dengan lencana HQ, status Pinned, kontrol playback, dan slider volume.",
    "dl.kicker": "Rilis Resmi GitHub",
    "dl.title": "Pilih paket sesuai sistem operasi",
    "dl.desc": "Tautan unduhan ditarik langsung dari rilis terbaru repositori komporkode/dengar atau mengarah ke arsip GitHub Releases.",
    "dl.statusLoading": "Memeriksa rilis terbaru...",
    "dl.statusReady": "Sinkron dengan GitHub Releases",
    "dl.statusFallback": "Mode cadangan (GitHub API tidak terjangkau)",
    "dl.metaPrefix": "Versi",
    "dl.published": "Dirilis",
    "dl.fallbackBanner": "Metadata langsung dari API GitHub gagal dimuat (batasan kuota atau gangguan jaringan). Tautan tetap mengarah ke berkas rilis resmi v1.4.1 di GitHub.",
    "dl.retry": "Coba Muat Ulang",
    "dl.allReleases": "Semua Rilis di GitHub",
    "dl.tablistAria": "Pilih platform desktop",
    "dl.detected": "(Sistem Anda)",
    "dl.downloadBtn": "Unduh",
    "dl.sizeUnknown": "Ukuran lihat di GitHub",
    "dl.unavailable": "Belum tersedia pada rilis ini. Periksa halaman GitHub Releases.",
    "os.windows": "Windows (64-bit)",
    "os.linux": "Linux (x64)",
    "os.macos": "macOS (Universal)",
    "art.winInstaller.title": "Installer Windows (NSIS .exe)",
    "art.winInstaller.desc": "Mendukung pilihan folder instalasi dan pembaruan otomatis in-place.",
    "art.winPortable.title": "Portable ZIP Windows (.zip)",
    "art.winPortable.desc": "Ekstrak dan jalankan langsung tanpa proses instalasi ke sistem.",
    "art.linuxBin.title": "Biner Linux x86_64 (.tar.gz)",
    "art.linuxBin.desc": "Arsip executable terkompilasi Wails untuk distribusi desktop Linux.",
    "art.macBin.title": "Aplikasi macOS Universal (.zip)",
    "art.macBin.desc": "Bundel .app universal untuk arsitektur Apple Silicon dan Intel.",
    "notes.win.title": "Catatan Windows 10 & 11",
    "notes.win.1": "Membutuhkan runtime Microsoft Edge WebView2 (terpasang bawaan di Windows 11).",
    "notes.win.2": "Mendukung tombol media hardware (SMTC), System Tray, dan pembaruan 1 klik.",
    "notes.win.3": "Data konfigurasi dan playlist tersimpan di %APPDATA%\\dengar.",
    "notes.linux.title": "Catatan Linux Desktop",
    "notes.linux.1": "Memerlukan pustaka GTK 3 dan WebKit2GTK 4.0 untuk merender antarmuka desktop.",
    "notes.linux.2": "Konfigurasi lokal dan cache disk tersimpan di ~/.config/dengar.",
    "notes.linux.3": "Dapat juga dikompilasi dari kode sumber memakai perintah wails build.",
    "notes.mac.title": "Catatan macOS",
    "notes.mac.1": "Dikemas sebagai arsip .zip universal (.app) dari GitHub Actions macOS.",
    "notes.mac.2": "Jika Gatekeeper menahan aplikasi baru, buka klik kanan lalu pilih Open.",
    "notes.mac.3": "Kode sumber terbuka GPL-3.0 tersedia penuh untuk audit atau build manual.",
    "feat.kicker": "Arsitektur Audio",
    "feat.title": "Dirancang untuk pemutaran stabil tanpa pelacakan",
    "feat.desc": "Dengerin memisahkan mesin streaming dan cache di backend Go agar antarmuka Svelte tetap responsif saat memutar audio jangka panjang.",
    "feat.mediaCap": "Kendali frekuensi 5-band (60Hz hingga 14kHz) berjalan lokal via Web Audio API dan internal CORS streaming proxy.",
    "feat.c1.title": "Failover Multi-Sumber",
    "feat.c1.desc": "Menarik stream audio Opus/AAC langsung lewat katalog YouTube Music dengan rotasi instance otomatis ke Piped dan Invidious bila server lambat, serta menyaring podcast atau vlog.",
    "feat.c2.title": "Cache Lokal & Equalizer",
    "feat.c2.desc": "Menyimpan lagu putaran ke disk dengan batas kuota LRU yang dapat diatur (termasuk pindah drive), ditambah equalizer studio 5-band, 8 preset suara, crossfade, dan mode hemat kuota.",
    "feat.c3.title": "Playlist Mosaik & Berbagi P2P",
    "feat.c3.desc": "Membuat cover art kolase 2x2 otomatis dari empat lagu teratas, mengimpor playlist publik YouTube atau Spotify, serta bertukar daftar putar melalui kode Base64 atau berkas JSON.",
    "feat.c4.title": "Integrasi Native Desktop",
    "feat.c4.desc": "Menyediakan Floating Mini Player (Always on Top), kontrol media keyboard (Windows SMTC), System Tray, Discord Rich Presence, timer tidur smooth fade-out, dan 1-click updater.",
    "rel.kicker": "Catatan Rilis",
    "rel.title": "Riwayat pembaruan terbaru",
    "rel.desc": "Disarikan langsung dari berkas CHANGELOG.md repositori proyek. Tekan Ctrl+F untuk mencari catatan versi tertentu secara instan.",
    "rel.141.badge": "Perbaikan Bug",
    "rel.141.title": "Isolasi Lagu Populer & Sinkronisasi Versi",
    "rel.141.date": "8 November 2026",
    "rel.141.lead": "Memisahkan state katalog tren beranda dari hasil pencarian agar query artis tidak menimpa daftar Lagu Populer saat input dibersihkan.",
    "rel.141.p1": "State daftar Lagu Populer terlindungi saat melakukan pencarian musisi.",
    "rel.141.p2": "Menghapus teks pencarian langsung mengembalikan katalog tren beragam.",
    "rel.141.p3": "Label versi pada dialog Pengaturan dan tombol Apa yang Baru terhubung dinamis ke konstanta rilis.",
    "rel.140.badge": "Fitur Utama",
    "rel.140.title": "Redesign Antrean, Aksesibilitas & Auto-Update",
    "rel.140.date": "7 November 2026",
    "rel.140.lead": "Pembaruan antarmuka antrean 1 baris, progress bar updater live (MB/s), dukungan navigasi Escape dan keyboard penuh, serta stabilitas cold-start.",
    "rel.140.p1": "Auto-updater menampilkan byte terunduh, kecepatan transfer, dan restart in-place otomatis.",
    "rel.140.p2": "Dukungan aksesibilitas Escape, aria-label lengkap, focus-visible kontras tinggi, dan prefers-reduced-motion.",
    "rel.140.p3": "Perombakan laci antrean dengan aksi Bersihkan Sisa dan perbaikan konflik file .tmp saat cold start.",
    "rel.132.badge": "Perbaikan",
    "rel.132.title": "Infinite Scroll Pencarian & Pipeline Installer",
    "rel.132.date": "3 November 2026",
    "rel.132.lead": "Memulihkan token kontinuitas pencarian tanpa batas, persistensi batas disk cache LRU, serta kelengkapan build NSIS installer Windows pada CI/CD.",
    "rel.132.p1": "Pencarian memuat 20+ lagu lanjutan secara otomatis saat halaman digulir.",
    "rel.132.p2": "Kuota penyimpanan LRU (500 MB hingga Tanpa Batas) tersimpan permanen setelah restart.",
    "rel.132.p3": "Pipeline GitHub Actions menghasilkan artefak Multi-OS lengkap.",
    "rel.130.badge": "Fitur & Cache",
    "rel.130.title": "Lokasi Folder Cache Kustom & Mode Eco Default",
    "rel.130.date": "3 November 2026",
    "rel.130.lead": "Menambahkan kebebasan memilih direktori cache di drive eksternal dengan migrasi file audio otomatis dan pengingat direktori upgrade di installer.",
    "rel.130.p1": "Pindahkan direktori cache audio antar-drive tanpa kehilangan lagu tersimpan offline.",
    "rel.130.p2": "Mode Eco aktif sebagai bawaan untuk hemat kuota data.",
    "rel.120.badge": "Ekosistem OS",
    "rel.120.title": "SMTC Windows, Discord RPC & Impor Playlist",
    "rel.120.date": "3 November 2026",
    "rel.120.lead": "Integrasi tombol media hardware, Discord Rich Presence via Named Pipe, tab lagu offline terpin, dan algoritma rekomendasi FYP multi-artis.",
    "rel.120.p1": "Impor tautan playlist publik YouTube maupun Spotify langsung ke pustaka lokal.",
    "rel.120.p2": "Lagu favorit dan offline terlindungi agar tidak ikut terhapus pembersihan cache LRU.",
    "spec.kicker": "Transparansi & Kode Sumber",
    "spec.title": "Komponen teknis & kompilasi mandiri",
    "spec.desc": "Seluruh kode backend dan frontend terbuka di bawah lisensi GNU GPL-3.0 untuk diperiksa atau dikompilasi langsung di perangkat Anda.",
    "spec.dt1": "Backend & Runtime",
    "spec.dd1": "Go 1.25 + Wails v2.17",
    "spec.dt2": "Antarmuka Desktop",
    "spec.dd2": "Svelte 5, TypeScript, Tailwind CSS v4",
    "spec.dt3": "Penyimpanan Lokal",
    "spec.dd3": "Atomic JSON lokal tanpa server akun cloud",
    "spec.dt4": "Lisensi Perangkat Lunak",
    "spec.dd4": "GNU General Public License v3.0 (GPL-3.0)",
    "spec.buildTitle": "Bangun langsung dari source code",
    "spec.copyBtn": "Salin Perintah",
    "spec.copiedBtn": "Tersalin",
    "spec.copyReq": "Prasyarat: Go 1.25+, Node.js 20+, dan Wails CLI v2.",
    "footer.note": "Aplikasi pemutar audio YouTube desktop open-source berlisensi GNU GPL-3.0. Copyright (C) 2026 Komporkode.",
    "footer.linksAria": "Tautan proyek",
    "footer.repo": "Repositori GitHub",
    "footer.releases": "Daftar Rilis",
    "footer.changelog": "Changelog",
    "footer.issues": "Laporkan Bug",
    "footer.license": "Lisensi GPL-3.0",
    "announce.lang": "Bahasa situs diubah ke Bahasa Indonesia.",
    "announce.tab": "Menampilkan paket unduhan untuk"
  },
  en: {
    "meta.title": "Dengerin Desktop - Ad-Free Desktop YouTube Music Player",
    "meta.desc": "Lightweight desktop music player powered by Go and Wails with offline local disk cache, 5-band equalizer, and multi-source audio failover.",
    "skip": "Skip to main content",
    "nav.download": "Download",
    "nav.screens": "Screens",
    "nav.features": "Features",
    "nav.releases": "Releases",
    "nav.specs": "Specs",
    "nav.github": "GitHub",
    "nav.primaryAria": "Primary",
    "nav.langAria": "Language selector",
    "nav.drawerAria": "Mobile menu",
    "nav.menuOpen": "Open navigation menu",
    "nav.menuClose": "Close menu",
    "hero.title": "Stream YouTube music ad-free, right from your desktop.",
    "hero.lead": "Built with Go and Wails. Stream YouTube Music with local disk cache, 5-band equalizer, and 1-click updates.",
    "hero.ctaPrimary": "Download Dengerin",
    "hero.ctaSecondary": "View Source on GitHub",
    "hero.zoomAria": "Enlarge main Dengerin Desktop screenshot",
    "hero.zoomHint": "Click image to view full size",
    "hero.mediaAlt": "Real Dengerin Desktop screenshot showing the Explore & Search screen with the floating player bar below.",
    "shot.kicker": "Real Interface",
    "shot.title": "Desktop interface screenshots",
    "shot.desc": "Actual captures taken from Dengerin Desktop on Windows. Select any screenshot to view full size.",
    "shot.expandLabel": "Enlarge",
    "shot.miniKicker": "Compact Mode",
    "shot.zoomBtn": "Inspect Mini Player",
    "shot.prevAria": "Previous screenshot",
    "shot.nextAria": "Next screenshot",
    "shot.closeAria": "Close preview",
    "shot.s1.title": "Explore & Search",
    "shot.s1.desc": "Direct YouTube Music search, frequent tracks, and the floating audio bar.",
    "shot.s1.aria": "Enlarge Explore & Search screenshot",
    "shot.s1.alt": "Explore & Search view with catalog search bar, frequently played tracks, and Eco Audio bottom player.",
    "shot.s2.title": "For You Recommendations",
    "shot.s2.desc": "Personal mix hero, popular studio picks, and Play All / Shuffle controls.",
    "shot.s2.aria": "Enlarge For You Recommendations screenshot",
    "shot.s2.alt": "For You tab showing Personal Mix banner, Popular & Studio selections, and Frequent Tracks.",
    "shot.s3.title": "5-Band Audio Equalizer",
    "shot.s3.desc": "Local 60 Hz to 14 kHz band sliders paired with presets like Vocal Booster.",
    "shot.s3.aria": "Enlarge 5-Band Audio Equalizer screenshot",
    "shot.s3.alt": "5-Band Audio Equalizer modal showing sound presets and 60 Hz to 14 kHz gain sliders.",
    "shot.s4.title": "Sleep Timer",
    "shot.s4.desc": "Quick 15 to 60 minute timers, stop after current track, and 5-second fade-out.",
    "shot.s4.aria": "Enlarge Sleep Timer screenshot",
    "shot.s4.alt": "Sleep Timer modal with 15 to 60 minute quick buttons and stop at end of track control.",
    "shot.s5.title": "Floating Mini Player (Always on Top)",
    "shot.s5.desc": "Compact window pinned over other apps for instant pause, track skip, seek, and volume controls.",
    "shot.s5.aria": "Enlarge Floating Mini Player screenshot",
    "shot.s5.alt": "Compact Dengerin Mini Player window showing HQ badge, Pinned status, playback controls, and volume slider.",
    "dl.kicker": "Official GitHub Releases",
    "dl.title": "Choose the package for your system",
    "dl.desc": "Download links are fetched directly from the komporkode/dengar release feed with fallback links to GitHub Releases.",
    "dl.statusLoading": "Checking latest release...",
    "dl.statusReady": "Synced with GitHub Releases",
    "dl.statusFallback": "Fallback mode (GitHub API unreachable)",
    "dl.metaPrefix": "Version",
    "dl.published": "Published",
    "dl.fallbackBanner": "Could not load live metadata from GitHub API (rate limit or network interruption). Direct links point to the baseline v1.4.1 GitHub Release assets.",
    "dl.retry": "Retry Sync",
    "dl.allReleases": "All GitHub Releases",
    "dl.tablistAria": "Select desktop platform",
    "dl.detected": "(Detected)",
    "dl.downloadBtn": "Download",
    "dl.sizeUnknown": "Check size on GitHub",
    "dl.unavailable": "Not attached to this release yet. Check GitHub Releases directly.",
    "os.windows": "Windows (64-bit)",
    "os.linux": "Linux (x64)",
    "os.macos": "macOS (Universal)",
    "art.winInstaller.title": "Windows Installer (NSIS .exe)",
    "art.winInstaller.desc": "Supports custom install directory selection and silent in-place updates.",
    "art.winPortable.title": "Windows Portable Archive (.zip)",
    "art.winPortable.desc": "Extract and launch directly without writing system installation entries.",
    "art.linuxBin.title": "Linux x86_64 Binary (.tar.gz)",
    "art.linuxBin.desc": "Compiled Wails desktop executable archive for Linux distributions.",
    "art.macBin.title": "macOS Universal App (.zip)",
    "art.macBin.desc": "Universal .app bundle supporting both Apple Silicon and Intel Macs.",
    "notes.win.title": "Windows 10 & 11 Notes",
    "notes.win.1": "Requires Microsoft Edge WebView2 runtime (built into Windows 11).",
    "notes.win.2": "Supports hardware media keys (SMTC), System Tray controls, and 1-click updates.",
    "notes.win.3": "Configuration, playlists, and state persist under %APPDATA%\\dengar.",
    "notes.linux.title": "Linux Desktop Notes",
    "notes.linux.1": "Requires GTK 3 and WebKit2GTK 4.0 shared libraries on the host.",
    "notes.linux.2": "Local settings and audio disk cache are stored inside ~/.config/dengar.",
    "notes.linux.3": "You can also compile from source using the wails build CLI.",
    "notes.mac.title": "macOS Notes",
    "notes.mac.1": "Packaged as a universal .app archive built by GitHub Actions.",
    "notes.mac.2": "If Gatekeeper blocks unsigned binaries on first launch, right-click and choose Open.",
    "notes.mac.3": "Full GPL-3.0 source code is open for local inspection and manual builds.",
    "feat.kicker": "Audio Architecture",
    "feat.title": "Built for steady playback without account tracking",
    "feat.desc": "Dengerin runs stream resolution and local caching inside a Go backend so the Svelte interface stays fast during long listening sessions.",
    "feat.mediaCap": "5-band frequency control (60Hz to 14kHz) powered locally via Web Audio API and an internal CORS stream proxy.",
    "feat.c1.title": "Multi-Source Failover",
    "feat.c1.desc": "Resolves Opus/AAC streams via YouTube Music with automatic instance rotation across Piped and Invidious plus non-music video filtering.",
    "feat.c2.title": "Disk Cache & Equalizer",
    "feat.c2.desc": "Caches played tracks locally with configurable LRU quota (custom drive supported), a 5-band studio EQ with 8 presets, crossfade, and Eco mode.",
    "feat.c3.title": "Mosaic Playlists & P2P Share",
    "feat.c3.desc": "Generates 2x2 mosaic covers automatically, imports public YouTube or Spotify playlists, and shares playlists P2P as Base64 text or JSON files.",
    "feat.c4.title": "Native Desktop Controls",
    "feat.c4.desc": "Includes an Always-on-Top Floating Mini Player, hardware media keys (SMTC), System Tray mode, Discord Rich Presence, sleep timer fade-out, and 1-click self-updater.",
    "rel.kicker": "Release History",
    "rel.title": "Recent updates",
    "rel.desc": "Taken directly from the repository CHANGELOG.md. Use Ctrl+F to search across versions immediately.",
    "rel.141.badge": "Bug Fix",
    "rel.141.title": "Popular Tracks Isolation & Version Sync",
    "rel.141.date": "November 8, 2026",
    "rel.141.lead": "Separates home trending catalog state from artist search results so clearing the query restores diverse Popular Tracks.",
    "rel.141.p1": "Popular Tracks state remains isolated during artist or song searches.",
    "rel.141.p2": "Clearing the search box automatically restores the home trending feed.",
    "rel.141.p3": "Release badges and What's New triggers stay synced to CURRENT_APP_VERSION.",
    "rel.140.badge": "Major Update",
    "rel.140.title": "Queue Redesign, Accessibility & Auto-Updater",
    "rel.140.date": "November 7, 2026",
    "rel.140.lead": "Introduces the single-row Queue Drawer, live download speed updater UI, full Escape/keyboard coverage, and cold-start playback fixes.",
    "rel.140.p1": "Real-time updater progress bar with transfer rate (MB/s) and silent in-place restart.",
    "rel.140.p2": "Full Escape dismiss support, aria-labels, accessible focus rings, and prefers-reduced-motion.",
    "rel.140.p3": "Structured queue sections with Clear Upcoming action and resolved cold-start file locks.",
    "rel.132.badge": "Fixes",
    "rel.132.title": "Infinite Search Pagination & CI Installer",
    "rel.132.date": "November 3, 2026",
    "rel.132.lead": "Restores infinite search continuation tokens, persistent LRU cache limits, and automated NSIS Windows installer builds.",
    "rel.132.p1": "Search scrolling loads 20+ additional tracks reliably without breaking non-music filters.",
    "rel.132.p2": "LRU disk quota (500 MB to Unlimited) persists cleanly across app restarts.",
    "rel.132.p3": "Multi-OS GitHub Actions workflow builds Windows, Linux, and macOS packages.",
    "rel.130.badge": "Storage",
    "rel.130.title": "Custom Cache Drive & Smart Upgrade Paths",
    "rel.130.date": "November 3, 2026",
    "rel.130.lead": "Allows choosing any drive for the offline audio cache with automatic file migration and upgrade path memory in the installer.",
    "rel.130.p1": "Move cached tracks across drives directly from Settings without re-downloading.",
    "rel.130.p2": "Eco quality mode enabled as the responsive data-saving default.",
    "rel.120.badge": "OS Integration",
    "rel.120.title": "Windows SMTC, Discord RPC & Playlist Import",
    "rel.120.date": "November 3, 2026",
    "rel.120.lead": "Adds hardware media keys, Discord Rich Presence via Named Pipe, pinned offline library tab, and multi-artist FYP mix balancing.",
    "rel.120.p1": "Import public YouTube or Spotify playlists straight into local Dengerin collections.",
    "rel.120.p2": "Offline and favorite tracks stay pinned and protected from LRU cache eviction.",
    "spec.kicker": "Open Source Specs",
    "spec.title": "Technical stack & local builds",
    "spec.desc": "Both backend and frontend source code are published under the GNU GPL-3.0 license for inspection or building on your own machine.",
    "spec.dt1": "Backend & Runtime",
    "spec.dd1": "Go 1.25 + Wails v2.17",
    "spec.dt2": "Desktop UI",
    "spec.dd2": "Svelte 5, TypeScript, Tailwind CSS v4",
    "spec.dt3": "Local Data",
    "spec.dd3": "Local atomic JSON storage with no cloud account",
    "spec.dt4": "License",
    "spec.dd4": "GNU General Public License v3.0 (GPL-3.0)",
    "spec.buildTitle": "Compile from source",
    "spec.copyBtn": "Copy Commands",
    "spec.copiedBtn": "Copied",
    "spec.copyReq": "Prerequisites: Go 1.25+, Node.js 20+, and Wails CLI v2.",
    "footer.note": "Open-source desktop YouTube audio player licensed under GNU GPL-3.0. Copyright (C) 2026 Komporkode.",
    "footer.linksAria": "Project links",
    "footer.repo": "GitHub Repo",
    "footer.releases": "Releases",
    "footer.changelog": "Changelog",
    "footer.issues": "Report Issue",
    "footer.license": "GPL-3.0 License",
    "announce.lang": "Language switched to English.",
    "announce.tab": "Showing downloads for"
  }
};

function detectPreferredOs() {
  const hint = (navigator.userAgentData?.platform || navigator.platform || "").toLowerCase();
  if (hint.includes("mac") || hint.includes("darwin")) return "macos";
  if (hint.includes("linux") || hint.includes("x11")) return "linux";
  if (hint.includes("win")) return "windows";
  return "windows";
}

function isSafeReleaseUrl(rawUrl) {
  try {
    const parsed = new URL(rawUrl);
    return parsed.protocol === "https:" && rawUrl.startsWith(ALLOWED_RELEASE_PREFIX);
  } catch {
    return false;
  }
}

function formatFileSize(bytes, lang) {
  if (typeof bytes !== "number" || bytes <= 0) {
    return I18N[lang]["dl.sizeUnknown"];
  }
  const mb = bytes / (1024 * 1024);
  return `${mb.toFixed(1)} MB`;
}

function formatDate(isoString, lang) {
  if (!isoString) return "";
  try {
    const dt = new Date(isoString);
    return new Intl.DateTimeFormat(lang === "id" ? "id-ID" : "en-US", {
      year: "numeric",
      month: "long",
      day: "numeric"
    }).format(dt);
  } catch {
    return "";
  }
}

function t(lang, key) {
  return I18N[lang]?.[key] ?? I18N.id[key] ?? key;
}

class DengarSiteController {
  constructor() {
    this.lang = "id";
    this.detectedOs = detectPreferredOs();
    this.activeOs = this.detectedOs;
    this.releaseState = "loading";
    this.releaseData = FALLBACK_RELEASE;

    this.liveRegion = document.getElementById("live-announcer");
    this.heroDownloadBtn = document.getElementById("hero-download-btn");
    this.heroDownloadLabel = document.getElementById("hero-download-label");
    this.statusPill = document.getElementById("release-status-pill");
    this.statusText = document.getElementById("release-status-text");
    this.releaseMeta = document.getElementById("release-meta");
    this.errorBanner = document.getElementById("release-error-banner");
    this.errorText = document.getElementById("release-error-text");
    this.retryBtn = document.getElementById("release-retry-btn");
    this.artifactList = document.getElementById("artifact-list");
    this.notesTitle = document.getElementById("platform-notes-title");
    this.notesList = document.getElementById("platform-notes-list");

    this.tabs = Array.from(document.querySelectorAll('[role="tab"][data-os]'));
    this.langBtns = Array.from(document.querySelectorAll("button[data-lang]"));
    this.copyBtn = document.getElementById("copy-build-btn");
    this.buildCode = document.getElementById("build-code-block");

    this.lightboxDialog = document.getElementById("screenshot-lightbox");
    this.lightboxImg = document.getElementById("lightbox-img");
    this.lightboxTitle = document.getElementById("lightbox-title");
    this.lightboxCaption = document.getElementById("lightbox-caption");
    this.lightboxCounter = document.getElementById("lightbox-counter");
    this.lightboxPrev = document.getElementById("lightbox-prev");
    this.lightboxNext = document.getElementById("lightbox-next");
    this.lightboxClose = document.getElementById("lightbox-close");
    this.activeShotIndex = 0;
  }

  init() {
    this.bindLangSwitch();
    this.bindTabs();
    this.bindRetry();
    this.bindCopyCommand();
    this.setupNavigationDrawer();
    this.setupLightbox();
    this.applyLanguage("id", false);
    this.selectPlatformTab(this.activeOs, false);
    this.fetchLatestRelease();
  }

  announce(message) {
    if (!this.liveRegion) return;
    this.liveRegion.textContent = "";
    window.requestAnimationFrame(() => {
      this.liveRegion.textContent = message;
    });
  }

  bindLangSwitch() {
    this.langBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const nextLang = btn.dataset.lang;
        if (nextLang && nextLang !== this.lang && I18N[nextLang]) {
          this.applyLanguage(nextLang, true);
        }
      });
    });
  }

  applyLanguage(lang, shouldAnnounce) {
    this.lang = lang;
    document.documentElement.lang = lang;
    document.title = t(lang, "meta.title");

    const metaDesc = document.querySelector('meta[name="description"]');
    if (metaDesc) {
      metaDesc.setAttribute("content", t(lang, "meta.desc"));
    }

    this.langBtns.forEach((btn) => {
      btn.setAttribute("aria-pressed", btn.dataset.lang === lang ? "true" : "false");
    });

    document.querySelectorAll("[data-i18n]").forEach((el) => {
      const key = el.getAttribute("data-i18n");
      if (key) {
        el.textContent = t(lang, key);
      }
    });

    document.querySelectorAll("[data-i18n-alt]").forEach((el) => {
      const key = el.getAttribute("data-i18n-alt");
      if (key) {
        el.setAttribute("alt", t(lang, key));
      }
    });

    document.querySelectorAll("[data-i18n-aria]").forEach((el) => {
      const key = el.getAttribute("data-i18n-aria");
      if (key) {
        el.setAttribute("aria-label", t(lang, key));
      }
    });

    this.renderTabsLabels();
    this.renderReleaseHeader();
    this.renderPlatformPanel();
    this.renderLightbox();
    this.updateHeroCta();

    if (shouldAnnounce) {
      this.announce(t(lang, "announce.lang"));
    }
  }

  updateHeroCta() {
    if (!this.heroDownloadLabel) return;
    const osShort =
      this.detectedOs === "windows"
        ? "Windows"
        : this.detectedOs === "macos"
          ? "macOS"
          : "Linux";
    const version = this.releaseData?.tag_name || FALLBACK_RELEASE.tag_name;
    this.heroDownloadLabel.textContent = `${t(this.lang, "hero.ctaPrimary")} ${version} (${osShort})`;
  }

  bindTabs() {
    this.tabs.forEach((tab, idx) => {
      tab.addEventListener("click", () => {
        const os = tab.dataset.os;
        if (os) this.selectPlatformTab(os, true);
      });

      tab.addEventListener("keydown", (event) => {
        let nextIdx = null;
        if (event.key === "ArrowRight") nextIdx = (idx + 1) % this.tabs.length;
        else if (event.key === "ArrowLeft") nextIdx = (idx - 1 + this.tabs.length) % this.tabs.length;
        else if (event.key === "Home") nextIdx = 0;
        else if (event.key === "End") nextIdx = this.tabs.length - 1;

        if (nextIdx !== null) {
          event.preventDefault();
          const targetTab = this.tabs[nextIdx];
          targetTab.focus();
          const os = targetTab.dataset.os;
          if (os) this.selectPlatformTab(os, true);
        }
      });
    });
  }

  renderTabsLabels() {
    this.tabs.forEach((tab) => {
      const os = tab.dataset.os;
      const def = PLATFORM_DEFINITIONS[os];
      if (!def) return;
      const labelEl = tab.querySelector("[data-tab-label]");
      const badgeEl = tab.querySelector("[data-tab-detected]");
      if (labelEl) {
        labelEl.textContent = t(this.lang, def.labelKey);
      }
      if (badgeEl) {
        badgeEl.textContent = os === this.detectedOs ? t(this.lang, "dl.detected") : "";
        badgeEl.hidden = os !== this.detectedOs;
      }
    });
  }

  selectPlatformTab(os, shouldAnnounce) {
    if (!PLATFORM_DEFINITIONS[os]) return;
    this.activeOs = os;
    this.tabs.forEach((tab) => {
      const active = tab.dataset.os === os;
      tab.setAttribute("aria-selected", active ? "true" : "false");
      tab.setAttribute("tabindex", active ? "0" : "-1");
    });

    const panel = document.getElementById("platform-panel");
    const activeTab = this.tabs.find((tab) => tab.dataset.os === os);
    if (panel && activeTab) {
      panel.setAttribute("aria-labelledby", activeTab.id);
    }

    this.renderPlatformPanel();

    if (shouldAnnounce) {
      this.announce(`${t(this.lang, "announce.tab")} ${t(this.lang, PLATFORM_DEFINITIONS[os].labelKey)}`);
    }
  }

  bindRetry() {
    if (!this.retryBtn) return;
    this.retryBtn.addEventListener("click", () => {
      this.fetchLatestRelease();
    });
  }

  async fetchLatestRelease() {
    this.releaseState = "loading";
    this.renderReleaseHeader();
    if (this.errorBanner) this.errorBanner.hidden = true;

    const controller = new AbortController();
    const timeoutId = window.setTimeout(() => controller.abort(), 8000);

    try {
      const res = await fetch(GITHUB_API_LATEST_URL, {
        headers: { Accept: "application/vnd.github+json" },
        signal: controller.signal
      });
      if (!res.ok) {
        throw new Error(`HTTP ${res.status}`);
      }
      const data = await res.json();
      if (!data || typeof data.tag_name !== "string" || !Array.isArray(data.assets)) {
        throw new Error("Invalid release schema");
      }

      const safeAssets = data.assets
        .filter((asset) => asset && typeof asset.name === "string" && isSafeReleaseUrl(asset.browser_download_url))
        .map((asset) => ({
          name: asset.name,
          size: typeof asset.size === "number" ? asset.size : null,
          browser_download_url: asset.browser_download_url
        }));

      this.releaseData = {
        tag_name: data.tag_name,
        published_at: typeof data.published_at === "string" ? data.published_at : "",
        html_url: isSafeReleaseUrl(data.html_url) ? data.html_url : GITHUB_LATEST_PAGE_URL,
        assets: safeAssets
      };
      this.releaseState = "ready";
      if (this.errorBanner) this.errorBanner.hidden = true;
    } catch {
      this.releaseData = FALLBACK_RELEASE;
      this.releaseState = "fallback";
      if (this.errorBanner && this.errorText) {
        this.errorText.textContent = t(this.lang, "dl.fallbackBanner");
        this.errorBanner.hidden = false;
      }
    } finally {
      window.clearTimeout(timeoutId);
      this.renderReleaseHeader();
      this.renderPlatformPanel();
      this.updateHeroCta();
    }
  }

  renderReleaseHeader() {
    if (!this.statusPill || !this.statusText || !this.releaseMeta) return;

    this.statusPill.dataset.state = this.releaseState;
    if (this.releaseState === "loading") {
      this.statusText.textContent = t(this.lang, "dl.statusLoading");
    } else if (this.releaseState === "ready") {
      this.statusText.textContent = t(this.lang, "dl.statusReady");
    } else {
      this.statusText.textContent = t(this.lang, "dl.statusFallback");
    }

    const version = this.releaseData?.tag_name || FALLBACK_RELEASE.tag_name;
    const pubDate = formatDate(this.releaseData?.published_at, this.lang);

    const versionChunk = document.createElement("span");
    versionChunk.className = "release-meta-item";
    versionChunk.textContent = `${t(this.lang, "dl.metaPrefix")} ${version}`;

    if (pubDate) {
      const dateChunk = document.createElement("span");
      dateChunk.className = "release-meta-item";
      dateChunk.textContent = `${t(this.lang, "dl.published")} ${pubDate}`;
      this.releaseMeta.replaceChildren(versionChunk, dateChunk);
    } else {
      this.releaseMeta.replaceChildren(versionChunk);
    }
  }

  renderPlatformPanel() {
    if (!this.artifactList || !this.notesTitle || !this.notesList) return;
    const def = PLATFORM_DEFINITIONS[this.activeOs];
    if (!def) return;

    const assets = Array.isArray(this.releaseData?.assets) ? this.releaseData.assets : [];

    this.artifactList.replaceChildren();
    def.artifacts.forEach((spec) => {
      const matched = assets.find((a) => a.name.toLowerCase().includes(spec.match.toLowerCase()));
      const card = document.createElement("div");
      card.className = spec.recommended ? "artifact-card is-recommended" : "artifact-card";

      const info = document.createElement("div");
      info.className = "artifact-info";

      const title = document.createElement("div");
      title.className = "artifact-title";
      title.textContent = t(this.lang, spec.titleKey);

      const desc = document.createElement("p");
      desc.className = "artifact-desc";
      desc.textContent = t(this.lang, spec.descKey);

      const sub = document.createElement("div");
      sub.className = "artifact-sub";

      info.append(title, desc, sub);

      if (matched && isSafeReleaseUrl(matched.browser_download_url)) {
        const sizeLabel = formatFileSize(matched.size, this.lang);
        sub.textContent = `${matched.name} (${sizeLabel})`;

        const link = document.createElement("a");
        link.className = spec.recommended ? "btn btn-primary btn-compact" : "btn btn-secondary btn-compact";
        link.href = matched.browser_download_url;
        link.rel = "noopener";
        link.setAttribute("download", matched.name);
        link.textContent = t(this.lang, "dl.downloadBtn");

        const icon = this.createSpriteIcon("#icon-download");
        link.prepend(icon);

        card.append(info, link);
      } else {
        sub.textContent = t(this.lang, "dl.unavailable");
        const fallbackLink = document.createElement("a");
        fallbackLink.className = "btn btn-secondary btn-compact";
        fallbackLink.href = this.releaseData?.html_url || GITHUB_LATEST_PAGE_URL;
        fallbackLink.target = "_blank";
        fallbackLink.rel = "noopener noreferrer";
        fallbackLink.textContent = t(this.lang, "dl.allReleases");
        card.append(info, fallbackLink);
      }

      this.artifactList.appendChild(card);
    });

    this.notesTitle.textContent = t(this.lang, def.notesTitleKey);
    this.notesList.replaceChildren();
    def.notesKeys.forEach((key) => {
      const li = document.createElement("li");
      li.textContent = t(this.lang, key);
      this.notesList.appendChild(li);
    });
  }

  createSpriteIcon(symbolId) {
    const svgNS = "http://www.w3.org/2000/svg";
    const svg = document.createElementNS(svgNS, "svg");
    svg.setAttribute("class", "ui-icon");
    svg.setAttribute("aria-hidden", "true");
    svg.setAttribute("width", "18");
    svg.setAttribute("height", "18");
    svg.setAttribute("viewBox", "0 0 24 24");
    const use = document.createElementNS(svgNS, "use");
    use.setAttribute("href", symbolId);
    svg.appendChild(use);
    return svg;
  }

  bindCopyCommand() {
    if (!this.copyBtn || !this.buildCode) return;
    let resetTimer = null;

    this.copyBtn.addEventListener("click", async () => {
      const text = this.buildCode.textContent?.trim() || "";
      if (!text) return;

      try {
        if (navigator.clipboard?.writeText) {
          await navigator.clipboard.writeText(text);
        } else {
          const area = document.createElement("textarea");
          area.value = text;
          area.setAttribute("readonly", "");
          area.className = "visually-hidden";
          document.body.appendChild(area);
          area.select();
          document.execCommand("copy");
          document.body.removeChild(area);
        }
        const label = this.copyBtn.querySelector("[data-copy-label]");
        if (label) {
          label.textContent = t(this.lang, "spec.copiedBtn");
        }
        this.announce(t(this.lang, "spec.copiedBtn"));

        window.clearTimeout(resetTimer);
        resetTimer = window.setTimeout(() => {
          if (label) {
            label.textContent = t(this.lang, "spec.copyBtn");
          }
        }, 2000);
      } catch {
        // Silent failure avoids false success indicator
      }
    });
  }

  setupNavigationDrawer() {
    const drawer = document.getElementById("nav-drawer");
    const openBtn = document.getElementById("drawer-open");
    const closeBtn = document.getElementById("drawer-close");
    const mainEl = document.getElementById("main-shell");
    const headerEl = document.querySelector(".site-header");
    const footerEl = document.querySelector(".site-footer");
    if (!drawer || !openBtn || !mainEl) return;

    const scroller = drawer.querySelector(".nav-drawer-scroller");
    const sheet = drawer.querySelector(".nav-drawer-sheet");
    if (!scroller || !sheet) return;

    if (!CSS.supports("animation-timeline: scroll()")) {
      scroller.addEventListener("scroll", () => {
        const width = sheet.offsetWidth || 1;
        const ratio = Math.max(0, Math.min(1, 1 - scroller.scrollLeft / width));
        drawer.style.setProperty("--drawer-backdrop", ratio.toFixed(3));
      }, { passive: true });
    }

    const supportsPopover = typeof drawer.showPopover === "function";
    if (!supportsPopover) {
      drawer.hidden = true;
    }

    const backgroundLandmarks = [headerEl, mainEl, footerEl].filter(Boolean);
    let isOpen = false;
    let shouldRestoreFocus = true;
    let closeFallbackTimer = null;

    const setBackgroundInert = (state) => {
      backgroundLandmarks.forEach((el) => {
        el.inert = state;
      });
    };

    const onOpened = () => {
      window.clearTimeout(closeFallbackTimer);
      if (isOpen) return;
      isOpen = true;
      setBackgroundInert(true);
      openBtn.setAttribute("aria-expanded", "true");
      sheet.focus();
    };

    const onClosed = (restoreFocus = shouldRestoreFocus) => {
      window.clearTimeout(closeFallbackTimer);
      if (!isOpen) return;
      isOpen = false;
      if (supportsPopover && drawer.matches(":popover-open")) {
        drawer.hidePopover();
      } else {
        drawer.hidden = true;
      }
      setBackgroundInert(false);
      openBtn.setAttribute("aria-expanded", "false");
      if (restoreFocus) {
        openBtn.focus();
      }
    };

    const openDrawer = async () => {
      window.clearTimeout(closeFallbackTimer);
      shouldRestoreFocus = true;
      if (supportsPopover) {
        drawer.showPopover();
      } else {
        drawer.hidden = false;
      }

      if (!CSS.supports("scroll-initial-target", "nearest")) {
        scroller.scrollTo({ left: scroller.offsetWidth, behavior: "instant" });
        await new Promise((resolve) =>
          window.requestAnimationFrame(() => window.requestAnimationFrame(resolve))
        );
      }
      scroller.scrollTo({ left: 0, behavior: "auto" });
      onOpened();
    };

    const closeDrawer = (restoreFocus = true) => {
      if (!isOpen) return;
      shouldRestoreFocus = restoreFocus;
      scroller.scrollTo({ left: scroller.offsetWidth, behavior: "auto" });
      if (window.matchMedia("(prefers-reduced-motion: reduce)").matches || !supportsPopover) {
        onClosed(restoreFocus);
        return;
      }
      window.clearTimeout(closeFallbackTimer);
      closeFallbackTimer = window.setTimeout(() => {
        if (isOpen) onClosed(shouldRestoreFocus);
      }, 450);
    };

    const visibleThreshold = 1 / Math.max(320, window.innerWidth);
    const observer = new IntersectionObserver(
      (entries) => {
        const entry = entries.at(-1);
        if (!entry) return;
        if (isOpen && entry.intersectionRatio < visibleThreshold) {
          onClosed(shouldRestoreFocus);
        }
        if (entry.intersectionRatio === 1 && ((supportsPopover && drawer.matches(":popover-open")) || !drawer.hidden)) {
          onOpened();
        }
      },
      { root: drawer, threshold: [visibleThreshold, 1] }
    );
    observer.observe(sheet);

    openBtn.addEventListener("click", openDrawer);
    closeBtn?.addEventListener("click", () => closeDrawer(true));

    drawer.addEventListener("click", (event) => {
      if (!sheet.contains(event.target)) {
        closeDrawer(true);
      }
    });

    sheet.querySelectorAll("a").forEach((link) => {
      link.addEventListener("click", () => {
        const href = link.getAttribute("href") || "";
        onClosed(!href.startsWith("#"));
      });
    });

    window.matchMedia("(min-width: 56rem)").addEventListener("change", (event) => {
      if (event.matches && isOpen) {
        onClosed(false);
      }
    });

    document.addEventListener("keydown", (event) => {
      if (event.key === "Escape" && isOpen) {
        closeDrawer(true);
      }
    });
  }

  setupLightbox() {
    const dialog = this.lightboxDialog;
    if (!dialog || !this.lightboxImg) return;

    const openShotById = (id, triggerEl) => {
      const index = SCREENSHOTS.findIndex((item) => item.id === id);
      if (index === -1) return;
      if (triggerEl && document.activeElement !== triggerEl) {
        triggerEl.focus();
      }
      this.activeShotIndex = index;
      this.renderLightbox();
      if (!dialog.open) {
        dialog.showModal();
      }
      this.lightboxClose?.focus();
    };

    const triggers = Array.from(document.querySelectorAll("[data-shot]"));
    triggers.forEach((trigger) => {
      trigger.addEventListener("click", () => {
        openShotById(trigger.getAttribute("data-shot"), trigger);
      });
    });

    document.querySelectorAll(".hero-media-caption, .shot-zoom-hint").forEach((hint) => {
      hint.addEventListener("click", () => {
        const trigger = hint.closest("figure")?.querySelector("button[data-shot]");
        if (trigger) {
          openShotById(trigger.getAttribute("data-shot"), trigger);
        }
      });
    });

    const step = (delta) => {
      this.activeShotIndex = (this.activeShotIndex + delta + SCREENSHOTS.length) % SCREENSHOTS.length;
      this.renderLightbox();
    };

    this.lightboxPrev?.addEventListener("click", () => step(-1));
    this.lightboxNext?.addEventListener("click", () => step(1));
    this.lightboxClose?.addEventListener("click", () => dialog.close());

    dialog.addEventListener("keydown", (event) => {
      if (!dialog.open) return;
      if (event.key === "ArrowLeft") {
        event.preventDefault();
        step(-1);
      } else if (event.key === "ArrowRight") {
        event.preventDefault();
        step(1);
      }
    });

    // Fallback for browsers lacking <dialog closedby="any"> (e.g. Safari)
    if (!("closedBy" in HTMLDialogElement.prototype)) {
      const isOutside = (event) => {
        const rect = dialog.getBoundingClientRect();
        return (
          event.clientX < rect.left ||
          event.clientX > rect.right ||
          event.clientY < rect.top ||
          event.clientY > rect.bottom
        );
      };
      let pressedOutside = false;
      dialog.addEventListener("pointerdown", (event) => {
        pressedOutside = isOutside(event);
      });
      dialog.addEventListener("pointerup", (event) => {
        if (pressedOutside && isOutside(event)) {
          dialog.close();
        }
      });
    }
  }

  renderLightbox() {
    const item = SCREENSHOTS[this.activeShotIndex];
    if (!item || !this.lightboxImg) return;

    if (this.lightboxDialog) {
      this.lightboxDialog.dataset.shot = item.id;
    }
    this.lightboxImg.src = item.src;
    this.lightboxImg.width = item.width;
    this.lightboxImg.height = item.height;
    this.lightboxImg.alt = t(this.lang, item.altKey);

    if (this.lightboxTitle) {
      this.lightboxTitle.textContent = t(this.lang, item.titleKey);
    }
    if (this.lightboxCaption) {
      this.lightboxCaption.textContent = t(this.lang, item.descKey);
    }
    if (this.lightboxCounter) {
      this.lightboxCounter.textContent = `${this.activeShotIndex + 1} / ${SCREENSHOTS.length}`;
    }
  }
}

document.addEventListener("DOMContentLoaded", () => {
  const app = new DengarSiteController();
  app.init();
});
