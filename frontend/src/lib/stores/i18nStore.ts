import { writable, derived } from 'svelte/store'

export type Locale = 'id' | 'en'

const savedLocale = (typeof localStorage !== 'undefined' ? localStorage.getItem('dengar_locale') : null) as Locale | null
export const currentLocale = writable<Locale>(savedLocale === 'en' ? 'en' : 'id')

currentLocale.subscribe((loc) => {
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('dengar_locale', loc)
  }
})

const translations: Record<Locale, Record<string, string>> = {
  id: {
    // Nav
    nav_explore: 'Eksplorasi & Cari',
    nav_favorites: 'Lagu Favorit',
    nav_history: 'Riwayat Putar',
    nav_queue: 'Antrean Lagu',
    nav_active: 'Aktif',
    section_browse: 'Menjelajah',
    section_collection: 'Koleksi Saya',

    // Search
    search_placeholder: 'Cari lagu, musisi, album, atau lirik...',
    filter_songs: 'Lagu',
    filter_videos: 'Video',
    filter_all: 'Semua Hasil',
    search_results_prefix: 'Hasil Pencarian: ',
    popular_songs: 'Eksplorasi Lagu Populer',
    no_results: 'Tidak ada hasil yang ditemukan untuk kata kunci ini.',
    empty_search_tip: 'Coba ketik kata kunci lain pada kolom pencarian di atas.',
    audio_exploration: 'Audio Eksplorasi',
    stream_engine: 'Stream Master Audio',

    // Controls
    play_all: 'Putar Semua',
    pause_all: 'Jeda Semua',
    shuffle: 'Acak',
    items_count: 'item',
    playing: 'Memutar',
    title_and_artist: 'Judul Lagu & Penyanyi',
    channel_or_artist: 'Channel / Penyanyi',
    quality_and_duration: 'Kualitas & Durasi',
    actions: 'Aksi',

    // Player
    no_track_playing: 'Belum ada lagu yang diputar',
    select_track_prompt: 'Pilih lagu untuk mulai memutar',
    bit_perfect: 'Stream Langsung',
    audio_mode_hq: 'HQ Audio',
    video_mode: 'Video MV',
    queue_title: 'Antrean Pemutaran',
    queue_empty: 'Antrean lagu masih kosong',
    queue_empty_tip: 'Pilih lagu dari pencarian untuk menambahkan',
    clear_all: 'Hapus Semua',

    // Actions & Toasts
    added_to_queue: 'Ditambahkan ke antrean',
    added_to_favorites: 'Disimpan ke lagu favorit',
    removed_from_favorites: 'Dihapus dari lagu favorit',
    toast_cleared_cache: 'Cache audio berhasil dibersihkan',

    // Settings & About
    settings_title: 'Pengaturan & Tentang',
    tab_general: 'Umum & Bahasa',
    tab_audio: 'Kualitas Audio',
    tab_about: 'Tentang Aplikasi',
    lang_select: 'Pilihan Bahasa Antarmuka',
    lang_id: 'Bahasa Indonesia',
    lang_en: 'English (US)',
    quality_setting: 'Preferensi Kualitas Audio',
    quality_high: 'Kualitas Tertinggi (Opus 160kbps / AAC)',
    quality_saver: 'Hemat Kuota / Cepat (~70-128kbps)',
    quality_notice: 'Catatan Kualitas: Audio bersumber dari YouTube Music dengan resolusi terbaik hingga 160kbps Opus / 256kbps AAC.',
    clear_cache_btn: 'Bersihkan Cache Stream',
    about_desc: 'Dengar adalah aplikasi pemutar musik desktop sumber terbuka yang ringan, modern, tanpa iklan, dan gratis selamanya (Rp0).',
    built_with: 'Dibuat dengan Go, Svelte, Tailwind CSS, dan Wails.',
    credits_title: 'Kredit & Kontributor',
    creator: 'Pengembang',
    close: 'Tutup',
    mini_player: 'Mini Player',
    restore_player: 'Kembali ke Layar Penuh'
  },
  en: {
    // Nav
    nav_explore: 'Explore & Search',
    nav_favorites: 'Favorites',
    nav_history: 'History',
    nav_queue: 'Playback Queue',
    nav_active: 'Active',
    section_browse: 'Browse',
    section_collection: 'My Collection',

    // Search
    search_placeholder: 'Search songs, artists, albums, or lyrics...',
    filter_songs: 'Songs',
    filter_videos: 'Videos',
    filter_all: 'All Results',
    search_results_prefix: 'Search Results: ',
    popular_songs: 'Trending & Popular Songs',
    no_results: 'No results found for this query.',
    empty_search_tip: 'Try searching with different keywords above.',
    audio_exploration: 'Audio Exploration',
    stream_engine: 'Master Audio Stream',

    // Controls
    play_all: 'Play All',
    pause_all: 'Pause All',
    shuffle: 'Shuffle',
    items_count: 'items',
    playing: 'Playing',
    title_and_artist: 'Song Title & Artist',
    channel_or_artist: 'Channel / Artist',
    quality_and_duration: 'Quality & Duration',
    actions: 'Actions',

    // Player
    no_track_playing: 'No track currently playing',
    select_track_prompt: 'Select a track to start playback',
    bit_perfect: 'Direct Stream',
    audio_mode_hq: 'HQ Audio',
    video_mode: 'Video MV',
    queue_title: 'Playback Queue',
    queue_empty: 'Queue is currently empty',
    queue_empty_tip: 'Add songs from search results to start queue',
    clear_all: 'Clear All',

    // Actions & Toasts
    added_to_queue: 'Added to playback queue',
    added_to_favorites: 'Saved to favorites',
    removed_from_favorites: 'Removed from favorites',
    toast_cleared_cache: 'Audio stream cache successfully cleared',

    // Settings & About
    settings_title: 'Settings & About',
    tab_general: 'General & Language',
    tab_audio: 'Audio Quality',
    tab_about: 'About App',
    lang_select: 'Interface Language',
    lang_id: 'Bahasa Indonesia',
    lang_en: 'English (US)',
    quality_setting: 'Audio Stream Quality',
    quality_high: 'Highest Quality (Opus 160kbps / AAC)',
    quality_saver: 'Data Saver / Fast (~70-128kbps)',
    quality_notice: 'Quality Notice: Audio is streamed directly from YouTube Music at up to 160kbps Opus / 256kbps AAC.',
    clear_cache_btn: 'Clear Audio Cache',
    about_desc: 'Dengar is an open-source, modern, lightweight desktop music player with zero ads and completely free forever.',
    built_with: 'Built with Go, Svelte, Tailwind CSS, and Wails.',
    credits_title: 'Credits & Contributors',
    creator: 'Developer',
    close: 'Close',
    mini_player: 'Mini Player',
    restore_player: 'Restore Full Player'
  }
}

export const t = derived(currentLocale, ($locale) => {
  return (key: string): string => {
    return translations[$locale]?.[key] || translations['id']?.[key] || key
  }
})
