export interface ChangelogRelease {
  version: string
  title: string
  date: string
  badge: string
  summary: string
  sections: {
    title: string
    icon: string
    color: string
    items: string[]
  }[]
}

export const CURRENT_APP_VERSION = 'v1.4.0'

export const CHANGELOG_DATA: ChangelogRelease[] = [
  {
    version: 'v1.4.0',
    title: 'Redesign Antrean, Craftsmanship Aksesibilitas & Auto-Update Sempurna',
    date: '7 November 2026',
    badge: 'Fitur Baru & Perbaikan',
    summary: 'Rilis v1.4.0 menghadirkan antarmuka antrean minimalis satu baris dengan pembagian visual terstruktur, status pemutar elegan ala Spotify Desktop, sistem auto-update & installer NSIS yang telah disempurnakan (progress bar live & in-place restart), craftsmanship anti-slop dengan aksesibilitas keyboard lengkap, serta stabilitas pemutaran cold-start.',
    sections: [
      {
        title: 'Sistem Auto-Update & Installer Disempurnakan',
        icon: 'sparkles',
        color: 'text-crimson-400',
        items: [
          'Progress Bar Unduhan Real-Time: Indikator visual live dengan byte terunduh, persentase akurat, dan kecepatan transfer (MB/s).',
          'Pemasangan & Mulai Ulang Otomatis: Upgrade in-place instan untuk installer (.exe) dan versi portabel (.zip) dengan proteksi file lock.',
          'Opsi Jalankan Dengerin di Installer: Wizard NSIS kini menyertakan opsi langsung membuka Dengerin di akhir instalasi.',
          'Deteksi Akurat Tipe Paket: Membedakan edisi Terpasang dan Portabel dengan verifikasi keberadaan uninstaller.',
          'Pratinjau Catatan Rilis: Akordeon interaktif untuk membaca changelog pembaruan langsung di halaman Pengaturan.',
        ],
      },
      {
        title: 'Craftsmanship, Aksesibilitas & Anti-Slop',
        icon: 'sparkles',
        color: 'text-amber-400',
        items: [
          'Navigasi Keyboard Escape Menyeluruh: Tekan Escape untuk menutup laci antrean, equalizer, timer tidur, dan seluruh modal dialog.',
          'Cakupan Aksesibilitas Penuh: Atribut aria-label lengkap pada semua tombol kendali dan indikator focus-visible crimson yang kontras.',
          'Dukungan Motion Halus: Integrasi prefers-reduced-motion untuk kenyamanan pengguna dengan sensitivitas gerak.',
          'Pembersihan Visual & Salinan: Menghapus radial blur orb asing dan mengganti buzzword pemasaran dengan terminologi audio yang membumi.',
        ],
      },
      {
        title: 'Antarmuka Antrean Minimalis & Terstruktur',
        icon: 'sparkles',
        color: 'text-crimson-400',
        items: [
          'Header Tunggal Ramping: Mengganti header bertingkat yang padat dengan toolbar minimalis satu baris (56px) tanpa ikon hamburger besar.',
          'Pemisahan Visual Section: Membagi daftar antrean menjadi Sedang Diputar (kartu crimson glowing dengan equalizer aktif), Berikutnya, dan Sebelumnya lengkap dengan badge hitung lagu.',
          'Aksi Bersihkan Sisa (Clear Upcoming): Tombol cepat untuk mengosongkan antrean lagu berikutnya tanpa menghapus riwayat yang sedang diputar.',
          'Autoplay Kompak: Toggle autoplay serupa kini menyatu rapi di toolbar kanan.',
        ],
      },
      {
        title: 'Pengalaman Pemutar & Navigasi Bersih',
        icon: 'sparkles',
        color: 'text-indigo-400',
        items: [
          'Spotify-Style Player Empty State: Saat belum ada lagu yang dipilih, tombol kendali (Play, Shuffle, Prev, Next, Repeat) dinonaktifkan dengan transisi redup halus tanpa pergeseran tata letak dan tanpa spinner macet. Slider volume dan equalizer tetap aktif.',
          'Sidebar Ramping: Menghilangkan teks indikator "Aktif" pada tombol navigasi, digantikan dengan sorotan visual halus (background highlight dan ikon merah).',
          'Lokalisasi Penuh (Bilingual ID/EN): 100% teks di antrean, pencarian, menu opsi lagu, modal playlist, dan pengaturan telah diterjemahkan dengan benar.',
        ],
      },
      {
        title: 'Perbaikan Mesin Pemutar & Keandalan Jaringan',
        icon: 'sparkles',
        color: 'text-emerald-400',
        items: [
          'Bebas Error Retry Pemutaran Pertama: Menghilangkan konflik unduhan duplikat di background yang memicu perebutan bandwidth dan file lock saat memutar lagu pertama.',
          'Watchdog & Stall Timer Adaptif: Mengakomodasi latensi cold-start HTTPS/DNS internasional dengan timer pemantau 15 detik dan stall 12 detik berproteksi buffer.',
          'Filter Lifecycle Browser: Mencegah AbortError alami dari Chromium WebView2 ditafsirkan sebagai kegagalan server.',
          'AudioContext Instan: Membuka kunci Web Audio API secara sinkron saat klik awal untuk mencegah pembatasan pemblokiran audio browser.',
        ],
      },
    ],
  },
  {
    version: 'v1.3.2',
    title: 'Perbaikan Infinite Scroll, Tampilan Sidebar & Installer Windows',
    date: '3 November 2026',
    badge: 'Bug Fix & Release',
    summary: 'Rilis v1.3.2 memperbaiki infinite scroll pencarian yang terhenti di 11 lagu, membersihkan titik merah yang menutupi ikon pada sidebar collapsed, menjaga persistensi batas kuota cache LRU, dan mengaktifkan build installer Windows (.exe) pada CI/CD Multi-OS.',
    sections: [
      {
        title: 'Perbaikan Pencarian & Infinite Scroll',
        icon: 'sparkles',
        color: 'text-crimson-400',
        items: [
          'Infinite Scroll Berfungsi Normal: Memulihkan token kontinuitas sehingga halaman pencarian otomatis memuat 20+ lagu berikutnya saat di-scroll tanpa batas.',
          'Penyaringan Non-Musik Tetap Aktif: Menyaring podcast, talkshow, gameplay, dan tutorial tanpa merusak pagination.',
          'Dukungan Kompilasi Musik Penuh: Lagu satuan dan video kompilasi musik (Best of, Kumpulan Lagu, Full Album) tetap muncul.',
        ],
      },
      {
        title: 'Penyempurnaan Tampilan & Setelan Cache',
        icon: 'sparkles',
        color: 'text-emerald-400',
        items: [
          'Sidebar Collapsed Bersih: Menghilangkan badge titik merah yang menghalangi ikon Favorit dan Riwayat Putar, serta merapikan header logo.',
          'Persistensi Kuota Disk Cache (LRU): Pilihan batas kapasitas disk (500 MB, 1 GB, 2 GB, atau Tanpa Batas) tersimpan permanen saat aplikasi ditutup dan dibuka kembali.',
          'Mode Hemat Data (Eco) Default: Menjamin mode audio hemat kuota aktif secara default.',
        ],
      },
      {
        title: 'Build & Installer Multi-OS',
        icon: 'desktop',
        color: 'text-indigo-400',
        items: [
          'Installer Windows (.exe): Menambahkan dependensi NSIS ke pipeline build GitHub Actions sehingga installer resmi terbuat sempurna.',
          'Rilis Lengkap Multi-OS: Menyediakan Installer Windows, Portable ZIP Windows, Linux tar.gz, dan macOS zip.',
        ],
      },
    ],
  },
  {
    version: 'v1.3.1',
    title: 'Pembersihan Konten Non-Musik & Integrasi YouTube Music',
    date: '3 November 2026',
    badge: 'Pembaruan',
    summary: 'Rilis v1.3.1 menghadirkan pencarian utama langsung dari katalog resmi YouTube Music, penyaringan otomatis konten non-musik (podcast, gameplay, vlog, tutorial), dan penyempurnaan rekomendasi lagu.',
    sections: [
      {
        title: 'Katalog Murni YouTube Music',
        icon: 'sparkles',
        color: 'text-crimson-400',
        items: [
          'Pencarian Utama YouTube Music: Menghubungkan langsung ke katalog resmi YouTube Music (WEB_REMIX) untuk hasil lagu yang bersih dan akurat.',
          'Penyaringan Konten Non-Musik Cerdas: Menyaring otomatis konten non-musik seperti podcast, wawancara, tutorial, gameplay, vlog, dan trailer film.',
          'Dukungan Video Kompilasi Musik: Tetap mempertahankan dan mengizinkan video kompilasi musik resmi (Best of, Kumpulan Lagu, Full Album).',
          'Perbaikan Format Durasi Regional: Parsing durasi akurat untuk format bertanda titik maupun titik dua.',
        ],
      },
    ],
  },
  {
    version: 'v1.3.0',
    title: 'Kustomisasi Lokasi Cache, Installer Fleksibel & Peningkatan Upgrade',
    date: '3 November 2026',
    badge: 'Fitur Baru',
    summary: 'Rilis v1.3.0 menghadirkan pemilihan folder penyimpanan cache kustom ke drive lain dengan migrasi otomatis, installer Windows yang fleksibel dengan deteksi upgrade otomatis, mode kualitas Eco default, dan peningkatan performa.',
    sections: [
      {
        title: 'Manajemen Penyimpanan & Cache Kustom',
        icon: 'database',
        color: 'text-crimson-400',
        items: [
          'Kustomisasi Folder Cache: Pilih folder penyimpanan cache audio di mana saja (termasuk drive D:, E:, atau eksternal) langsung dari Pengaturan.',
          'Migrasi Otomatis: File audio yang sudah di-cache otomatis dipindahkan ke direktori baru tanpa kehilangan lagu offline.',
          'Buka & Reset Cepat: Tombol cepat untuk membuka folder cache di File Explorer atau mereset ke folder default OS.',
        ],
      },
      {
        title: 'Installer Windows & Smart Upgrade',
        icon: 'desktop',
        color: 'text-indigo-400',
        items: [
          'Pemilihan Direktori Instalasi: Bebas menentukan folder tujuan saat menginstal aplikasi.',
          'Deteksi Upgrade Cerdas: Installer otomatis mengingat lokasi instalasi sebelumnya saat melakukan pembaruan versi.',
          'Penutupan Aman: Proses aplikasi lama ditutup secara otomatis sebelum file baru diekstrak untuk mencegah kendala "file in use".',
        ],
      },
      {
        title: 'Pengaturan & UI Polish',
        icon: 'sparkles',
        color: 'text-emerald-400',
        items: [
          'Default Eco Mode: Kualitas streaming audio otomatis disetel ke Saver/Eco yang hemat kuota dan responsif.',
          'Batas Kuota Cache Persisten: Pengaturan limit LRU cache tersimpan permanen di file konfigurasi lokal.',
          'Badge Sidebar Rapi: Indikator titik merah pada tombol favorit & riwayat diposisikan tepat tanpa menghalangi ikon.',
        ],
      },
    ],
  },
  {
    version: 'v1.2.0',
    title: 'Integrasi OS Desktop, Tab Offline & Perombakan Algoritma FYP',
    date: '3 November 2026',
    badge: 'Pembaruan Besar',
    summary: 'Pembaruan besar menghadirkan integrasi mendalam OS Windows (Hardware Media Keys / SMTC & System Tray), Discord Rich Presence, Ekosistem Playlist YouTube & Spotify, Tab Offline baru, Sleep Timer dengan smooth fade-out, serta algoritma rekomendasi FYP multi-artis berimbang.',
    sections: [
      {
        title: 'Integrasi OS & Ekosistem Desktop',
        icon: 'desktop',
        color: 'text-indigo-400',
        items: [
          'Windows SMTC & Media Keys: Kontrol putar, jeda, lagu berikutnya, dan scrubber langsung via keyboard hardware dan taskbar thumbnail.',
          'Windows System Tray: Ikon tray di taskbar dengan pilihan dialog konfirmasi saat tombol tutup (X) ditekan.',
          'Discord Rich Presence (RPC): Status aktivitas lagu yang sedang diputar di Discord dengan cover dan durasi via Windows Named Pipe.',
        ],
      },
      {
        title: 'Ekosistem Playlist & Mode Offline',
        icon: 'database',
        color: 'text-crimson-400',
        items: [
          'Impor Playlist YouTube & Spotify: Masukkan tautan playlist publik YouTube atau Spotify embed untuk langsung disimpan ke pustaka Dengerin Anda.',
          'Tab Tersimpan / Offline: Simpan lagu langsung ke penyimpanan lokal perangkat tanpa batas kuota internet.',
          'Perlindungan Pinning Cache: Lagu tersimpan dan lagu favorit tidak akan terhapus oleh mekanisme pembersihan LRU.',
        ],
      },
      {
        title: 'Audio Polish & Algoritma FYP',
        icon: 'sparkles',
        color: 'text-emerald-400',
        items: [
          'Perombakan Total Algoritma FYP: Rekomendasi kini menggabungkan hingga 6 artis teratas dengan batas maksimal 2 lagu per artis, discovery mix, dan pengacakan segar.',
          'Sleep Timer Nyaman: Pilihan timer durasi atau akhir lagu dengan penurunan volume bertahap (5 detik smooth fade-out).',
          'Zero Latency & Crossfade: Pre-fetch otomatis 15 detik sebelum lagu usai dan pengaturan transisi crossfade 0-10 detik.',
        ],
      },
    ],
  },
  {
    version: 'v1.1.2',
    title: 'Equalizer Studio, Desain Sinematik & Stabilitas Audio',
    date: '30 Oktober 2026',
    badge: 'Rilis Pembaruan',
    summary: 'Pembaruan besar yang menghadirkan mesin audio equalizer 5-band, pemutaran lebih tangguh anti-gagal, desain hero baru yang sinematik, dan responsivitas taktil.',
    sections: [
      {
        title: 'Fitur Baru Audio',
        icon: 'sliders',
        color: 'text-crimson-400',
        items: [
          'Equalizer Studio 5-Band: Pengaturan frekuensi manual 60Hz, 230Hz, 910Hz, 3.6kHz, 14kHz (-12dB hingga +12dB).',
          '8 Preset Suara Siap Pakai: Bass Boost, Bass Reducer, Vocal Booster, Treble Boost, Rock, Pop, Electronic, dan Flat.',
          'Internal CORS Streaming Proxy: Menghubungkan Web Audio API tanpa risiko pembisuan audio oleh browser security policy.',
        ],
      },
      {
        title: 'Pembaruan Tampilan & UI/UX',
        icon: 'sparkles',
        color: 'text-amber-400',
        items: [
          'Cinematic Full-Bleed Hero Banner: Tampilan utama halaman Untuk Anda yang menyatu mulus ke latar gelap dengan tipografi tebal dan cover artistik.',
          'Respon Taktil Instan: Efek micro-press saat lagu ditekan serta indikator gelombang/spinner seketika tanpa jeda.',
          'Margin Player Bar Simetris & Animasi Halus: Bar pemutar dan sidebar kini bergerak serentak dengan transisi hardware-accelerated.',
        ],
      },
      {
        title: 'Keandalan & Perbaikan Bug',
        icon: 'shield',
        color: 'text-emerald-400',
        items: [
          'Pencegahan Gagal Putar: Perbaikan validasi domain stream YouTube/Googlevideo dan rotasi instance otomatis.',
          'Stall Watchdog: Otomatis memulihkan dan merotasi server dalam 7 detik jika buffering audio mengalami kebuntuan.',
          'Pembaruan Mandiri Aman: Pembaruan in-place kini mengekstrak ke folder temp dengan elevasi hak akses admin otomatis.',
        ],
      },
    ],
  },
  {
    version: 'v1.1.1',
    title: 'Dukungan Multi-Installer & Identitas Dengerin',
    date: '29 Oktober 2026',
    badge: 'Minor Update',
    summary: 'Pembaruan identitas baru aplikasi menjadi Dengerin, penyempurnaan installer NSIS, dan perbaikan tata letak jendela.',
    sections: [
      {
        title: 'Perubahan Utama',
        icon: 'rocket',
        color: 'text-sky-400',
        items: [
          'Branding baru Dengerin Desktop dengan lisensi sumber terbuka.',
          'Tersedia dalam dua varian: Installer Windows (NSIS) dan Portable standalone.',
          'Peningkatan tata letak kartu pengaturan dan pembacaan informasi sistem.',
        ],
      },
    ],
  },
]

