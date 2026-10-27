<script lang="ts">
  import { onMount } from 'svelte'
  import type { NavigationTab, SearchResultItem, SearchFilter } from './lib/types'
  import { Search } from '../wailsjs/go/main/App'
  import { playerStore } from './lib/stores/playerStore'
  import { userDataStore } from './lib/stores/userDataStore'
  import { isMiniPlayer, toastMessage } from './lib/stores/settingsStore'
  import { t } from './lib/stores/i18nStore'

  import Sidebar from './lib/components/Sidebar.svelte'
  import SearchBar from './lib/components/SearchBar.svelte'
  import TrackList from './lib/components/TrackList.svelte'
  import PlayerBar from './lib/components/PlayerBar.svelte'
  import VideoModal from './lib/components/VideoModal.svelte'
  import QueueDrawer from './lib/components/QueueDrawer.svelte'
  import SettingsModal from './lib/components/SettingsModal.svelte'
  import MiniPlayer from './lib/components/MiniPlayer.svelte'

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

{#if $isMiniPlayer}
  <!-- Mini Music Player View -->
  <MiniPlayer />
{:else}
  <!-- Standard Desktop Application View -->
  <div class="w-full h-full flex flex-col bg-obsidian text-zinc-100 font-sans antialiased select-none overflow-hidden relative">
    
    <!-- Application Body Layout (Duplicate Mac Titlebar removed cleanly!) -->
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
            title={searchQuery ? `${$t('search_results_prefix')}“${searchQuery}”` : $t('popular_songs')}
            tracks={searchResults}
            isLoading={isSearching}
            emptyMessage={$t('no_results')}
          />

        {:else if activeTab === 'favorites'}
          <TrackList
            title={$t('nav_favorites')}
            tracks={$userDataStore.favorites}
            isLoading={$userDataStore.isLoading}
            emptyMessage="Belum ada lagu yang ditandai favorit. Klik ikon hati pada lagu untuk menyimpannya!"
          />

        {:else if activeTab === 'history'}
          <TrackList
            title={$t('nav_history')}
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

    <!-- Settings & About Modal -->
    <SettingsModal />

    <!-- Toast Notification -->
    {#if $toastMessage}
      <div class="fixed top-5 right-5 z-50 pointer-events-none animate-in fade-in slide-in-from-top-3 duration-200">
        <div class="px-4 py-2.5 rounded-xl bg-[#181A24]/95 border border-white/10 shadow-2xl text-xs font-medium text-white flex items-center gap-2.5 backdrop-blur-md">
          <span class="w-2 h-2 rounded-full bg-crimson-500 shadow-[0_0_8px_#E61E38]"></span>
          <span>{$toastMessage}</span>
        </div>
      </div>
    {/if}
  </div>
{/if}
