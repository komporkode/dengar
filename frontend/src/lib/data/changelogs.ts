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

export const CHANGELOG_DATA: ChangelogRelease[] = [
  {
    version: 'v1.2.0',
    title: 'Integrasi OS Desktop, Tab Offline & Perombakan Algoritma FYP',
    date: '3 November 2026',
    badge: 'Pembaruan Besar',
    summary: 'Pembaruan besar menghadirkan integrasi mendalam OS Windows (Hardware Media Keys / SMTC & System Tray), Discord Rich Presence, Ekosistem Playlist YouTube & Spotify, Tab Offline baru, Sleep Timer dengan smooth fade-out, serta algoritma rekomendasi FYP multi-artis berimbang.',
    sections: [
      {
        title: 'Integrasi OS & Ekosistem Desktop',
        icon: '🖥️',
        color: 'text-indigo-400',
        items: [
          'Windows SMTC & Media Keys: Kontrol putar, jeda, lagu berikutnya, dan scrubber langsung via keyboard hardware dan taskbar thumbnail.',
          'Windows System Tray: Ikon tray di taskbar dengan pilihan dialog konfirmasi saat tombol tutup (X) ditekan.',
          'Discord Rich Presence (RPC): Status aktivitas lagu yang sedang diputar di Discord dengan cover dan durasi via Windows Named Pipe.',
        ],
      },
      {
        title: 'Ekosistem Playlist & Mode Offline',
        icon: '💾',
        color: 'text-crimson-400',
        items: [
          'Impor Playlist YouTube & Spotify: Masukkan tautan playlist publik YouTube atau Spotify embed untuk langsung disimpan ke pustaka Dengerin Anda.',
          'Tab Tersimpan / Offline: Simpan lagu langsung ke penyimpanan lokal perangkat tanpa batas kuota internet.',
          'Perlindungan Pinning Cache: Lagu tersimpan dan lagu favorit tidak akan terhapus oleh mekanisme pembersihan LRU.',
        ],
      },
      {
        title: 'Audio Polish & Algoritma FYP',
        icon: '✨',
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
        icon: '🎛️',
        color: 'text-crimson-400',
        items: [
          'Equalizer Studio 5-Band: Pengaturan frekuensi manual 60Hz, 230Hz, 910Hz, 3.6kHz, 14kHz (-12dB hingga +12dB).',
          '8 Preset Suara Siap Pakai: Bass Boost, Bass Reducer, Vocal Booster, Treble Boost, Rock, Pop, Electronic, dan Flat.',
          'Internal CORS Streaming Proxy: Menghubungkan Web Audio API tanpa risiko pembisuan audio oleh browser security policy.',
        ],
      },
      {
        title: 'Pembaruan Tampilan & UI/UX',
        icon: '✨',
        color: 'text-amber-400',
        items: [
          'Cinematic Full-Bleed Hero Banner: Tampilan utama halaman Untuk Anda yang menyatu mulus ke latar gelap dengan tipografi tebal dan cover artistik.',
          'Respon Taktil Instan: Efek micro-press saat lagu ditekan serta indikator gelombang/spinner seketika tanpa jeda.',
          'Margin Player Bar Simetris & Animasi Halus: Bar pemutar dan sidebar kini bergerak serentak dengan transisi hardware-accelerated.',
        ],
      },
      {
        title: 'Keandalan & Perbaikan Bug',
        icon: '🛡️',
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
        icon: '🚀',
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

