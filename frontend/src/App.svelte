<script lang="ts">
  import { onMount } from 'svelte'
  import type { NavigationTab, SearchResultItem, SearchFilter } from './lib/types'
  import { Search } from '../wailsjs/go/main/App'
  import { Quit, WindowMinimise, WindowToggleMaximise } from '../wailsjs/runtime/runtime.js'
  import { playerStore } from './lib/stores/playerStore'
  import { queueStore } from './lib/stores/queueStore'
  import { userDataStore } from './lib/stores/userDataStore'

  import Sidebar from './lib/components/Sidebar.svelte'
  import SearchBar from './lib/components/SearchBar.svelte'
  import TrackList from './lib/components/TrackList.svelte'
  import PlayerBar from './lib/components/PlayerBar.svelte'
  import VideoModal from './lib/components/VideoModal.svelte'
  import QueueDrawer from './lib/components/QueueDrawer.svelte'

  let activeTab: NavigationTab = 'search'
  let searchQuery = ''
  let currentFilter: SearchFilter = 'music_songs'
  let searchResults: SearchResultItem[] = []
  let isSearching = false
  let searchError: string | null = null

  $: currentTrack = $playerStore.currentTrack

  onMount(async () => {
    // 1. Load user favorites & history from backend JSON storage
    await userDataStore.load()

    // 2. Perform initial search so app is instantly ready with tracks
    performSearch('Trending Music', 'music_songs')

    // 3. Register global in-app keyboard shortcuts
    window.addEventListener('keydown', handleGlobalKeydown)
    return () => {
      window.removeEventListener('keydown', handleGlobalKeydown)
    }
  })

  async function performSearch(query: string, filter: SearchFilter) {
    if (!query) {
      searchResults = []
      return
    }
    isSearching = true
    searchError = null
    try {
      const results = await Search(query, filter)
      searchResults = results || []
    } catch (err: any) {
      console.error('Search failed:', err)
      searchError = err?.message || 'Gagal mencari lagu. Periksa koneksi internet.'
      searchResults = []
    } finally {
      isSearching = false
    }
  }

  function handleSearchEvent(e: CustomEvent<{ query: string; filter: SearchFilter }>) {
    searchQuery = e.detail.query
    currentFilter = e.detail.filter
    performSearch(searchQuery || 'Popular Music', currentFilter)
  }

  function handleWindowClose() {
    try {
      Quit()
    } catch (e) {
      console.warn('Wails Quit error:', e)
    }
  }

  function handleWindowMinimise() {
    try {
      WindowMinimise()
    } catch (e) {
      console.warn('Wails WindowMinimise error:', e)
    }
  }

  function handleWindowToggleMaximise() {
    try {
      WindowToggleMaximise()
    } catch (e) {
      console.warn('Wails WindowToggleMaximise error:', e)
    }
  }

  function handleGlobalKeydown(e: KeyboardEvent) {
    const target = e.target as HTMLElement
    // Avoid intercepting keystrokes when user is typing in search input
    if (target && (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA')) {
      return
    }

    switch (e.code) {
      case 'Space':
        e.preventDefault()
        playerStore.togglePlay()
        break
      case 'ArrowLeft':
        e.preventDefault()
        playerStore.seek(Math.max(0, $playerStore.currentTime - 5))
        break
      case 'ArrowRight':
        e.preventDefault()
        playerStore.seek(Math.min($playerStore.duration, $playerStore.currentTime + 5))
        break
      case 'ArrowUp':
        e.preventDefault()
        playerStore.setVolume(Math.min(1, $playerStore.volume + 0.05))
        break
      case 'ArrowDown':
        e.preventDefault()
        playerStore.setVolume(Math.max(0, $playerStore.volume - 0.05))
        break
      case 'KeyM':
        playerStore.toggleMute()
        break
      case 'KeyV':
        playerStore.toggleMode()
        break
      case 'KeyJ':
        playerStore.playPrev()
        break
      case 'KeyL':
        playerStore.playNext()
        break
      case 'KeyF':
        if ($playerStore.currentTrack) {
          userDataStore.toggleFavorite($playerStore.currentTrack)
        }
        break
    }
  }
</script>

<div class="w-full h-full flex flex-col bg-obsidian text-zinc-100 font-sans antialiased select-none overflow-hidden">
  <!-- Top Native Window Title Bar -->
  <header
    class="h-9 px-4 flex items-center justify-between bg-[#08080C] border-b border-white/[0.05] text-zinc-400 text-xs shrink-0 select-none z-50"
    style="--wails-draggable:drag"
  >
    <div class="flex items-center gap-2" style="--wails-draggable:no-drag">
      <!-- Window Controls (Traffic lights) -->
      <div class="flex items-center gap-1.5 mr-2">
        <button
          class="w-3 h-3 rounded-full bg-[#FF5F56]/80 hover:bg-[#FF5F56] transition-colors cursor-pointer border-none p-0"
          on:click={handleWindowClose}
          title="Tutup Aplikasi"
          aria-label="Tutup Aplikasi"
        ></button>
        <button
          class="w-3 h-3 rounded-full bg-[#FFBD2E]/80 hover:bg-[#FFBD2E] transition-colors cursor-pointer border-none p-0"
          on:click={handleWindowMinimise}
          title="Minimalisir"
          aria-label="Minimalisir"
        ></button>
        <button
          class="w-3 h-3 rounded-full bg-[#27C93F]/80 hover:bg-[#27C93F] transition-colors cursor-pointer border-none p-0"
          on:click={handleWindowToggleMaximise}
          title="Maksimalkan"
          aria-label="Maksimalkan"
        ></button>
      </div>

      <div class="flex items-center gap-2 pl-2 border-l border-white/10 text-zinc-400 font-medium">
        <div class="w-4 h-4 rounded-md bg-gradient-to-tr from-crimson-600 to-rose-500 flex items-center justify-center">
          <svg class="w-2.5 h-2.5 text-white" viewBox="0 0 269.22 267.55" fill="none">
            <path fill="currentColor" d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"/>
          </svg>
        </div>
        <span class="text-zinc-200 font-semibold tracking-wide">Dengar</span>
        <span class="text-[9px] text-zinc-500 uppercase tracking-wider font-mono">Hi-Fi Audio Player</span>
      </div>
    </div>

    <!-- Center Window Title -->
    <div class="text-[11px] font-medium text-zinc-400 flex items-center gap-2 truncate max-w-md">
      {#if currentTrack}
        <span class="text-zinc-300 truncate">{currentTrack.title} — {currentTrack.uploaderName}</span>
      {:else}
        <span class="text-zinc-500">Master Audio Streaming Engine</span>
      {/if}
    </div>

    <!-- Right Window Controls -->
    <div class="flex items-center gap-2 text-zinc-400" style="--wails-draggable:no-drag">
      <button
        class="hover:text-zinc-200 transition-colors p-1 rounded hover:bg-white/[0.06] cursor-pointer"
        on:click={handleWindowMinimise}
        title="Minimalisir"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M19 12H5"></path></svg>
      </button>
      <button
        class="hover:text-zinc-200 transition-colors p-1 rounded hover:bg-white/[0.06] cursor-pointer"
        on:click={handleWindowToggleMaximise}
        title="Maksimalkan"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><rect height="16" rx="2" width="16" x="4" y="4"></rect></svg>
      </button>
      <button
        class="hover:text-crimson-400 transition-colors p-1 rounded hover:bg-crimson-500/20 cursor-pointer"
        on:click={handleWindowClose}
        title="Tutup"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M18 6L6 18M6 6l12 12"></path></svg>
      </button>
    </div>
  </header>

  <!-- Application Body Layout -->
  <div class="flex flex-1 overflow-hidden relative">
    <!-- Left Sidebar -->
    <Sidebar bind:activeTab />

    <!-- Main Content Viewport -->
    <main class="flex-1 flex flex-col overflow-y-auto bg-gradient-to-b from-[#12141F] via-obsidian to-obsidian pb-32">
      {#if activeTab === 'search'}
        <SearchBar
          bind:query={searchQuery}
          filter={currentFilter}
          isLoading={isSearching}
          on:search={handleSearchEvent}
        />

        {#if searchError}
          <div class="mx-8 mt-4 flex items-center gap-3 px-4 py-3 rounded-xl bg-crimson-500/10 border border-crimson-500/30 text-crimson-400 text-xs">
            <svg class="w-4 h-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <span>{searchError}</span>
          </div>
        {/if}

        <TrackList
          title={searchQuery ? `Hasil Pencarian: “${searchQuery}”` : 'Eksplorasi Lagu Populer'}
          tracks={searchResults}
          isLoading={isSearching}
          emptyMessage="Tidak ada hasil yang ditemukan untuk kata kunci ini."
        />

      {:else if activeTab === 'favorites'}
        <TrackList
          title="Lagu Favorit Saya"
          tracks={$userDataStore.favorites}
          isLoading={$userDataStore.isLoading}
          emptyMessage="Belum ada lagu yang ditandai favorit. Klik ikon hati pada lagu untuk menyimpannya!"
        />

      {:else if activeTab === 'history'}
        <TrackList
          title="Riwayat Pemutaran Terakhir"
          tracks={$userDataStore.history}
          isLoading={$userDataStore.isLoading}
          emptyMessage="Belum ada riwayat putar. Mulai putar lagu untuk mencatat riwayat."
        />
      {/if}
    </main>
  </div>

  <!-- Video Floating Player Overlay -->
  <VideoModal />

  <!-- Queue Drawer -->
  <QueueDrawer />

  <!-- Bottom Floating Island Player -->
  <PlayerBar />
</div>
