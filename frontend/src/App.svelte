<script lang="ts">
  import { onMount } from 'svelte'
  import type { NavigationTab, SearchResultItem, SearchFilter } from './lib/types'
  import { Search } from '../wailsjs/go/main/App'
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

<div class="app-layout">
  <!-- Left Sidebar -->
  <Sidebar bind:activeTab />

  <!-- Main Content Area -->
  <main class="main-viewport">
    <div class="content-scroll">
      {#if activeTab === 'search'}
        <SearchBar
          bind:query={searchQuery}
          filter={currentFilter}
          isLoading={isSearching}
          on:search={handleSearchEvent}
        />

        {#if searchError}
          <div class="error-banner">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <span>{searchError}</span>
          </div>
        {/if}

        <TrackList
          title={searchQuery ? `Hasil Pencarian: "${searchQuery}"` : 'Eksplorasi Lagu Populer'}
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
    </div>
  </main>

  <!-- Video Floating Player Overlay -->
  <VideoModal />

  <!-- Queue Drawer -->
  <QueueDrawer />

  <!-- Bottom Player Bar -->
  <PlayerBar />
</div>

<style>
  .app-layout {
    display: flex;
    width: 100vw;
    height: 100vh;
    background-color: var(--bg-primary);
    color: var(--text-primary);
    overflow: hidden;
    position: relative;
  }

  .main-viewport {
    flex: 1;
    height: calc(100vh - var(--player-height));
    overflow-y: auto;
    position: relative;
  }

  .content-scroll {
    max-width: 1080px;
    margin: 0 auto;
    padding: 32px 36px 60px 36px;
    display: flex;
    flex-direction: column;
  }

  .error-banner {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 18px;
    background: rgba(239, 68, 68, 0.12);
    border: 1px solid rgba(239, 68, 68, 0.3);
    border-radius: 10px;
    color: #f87171;
    font-size: 13px;
    margin-bottom: 20px;
  }
</style>
