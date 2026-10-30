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
