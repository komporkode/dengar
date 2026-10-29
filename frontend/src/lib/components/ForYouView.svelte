<script lang="ts">
  import { onMount } from 'svelte'
  import type { SearchResultItem } from '../types'
  import { GetHomeData, RefreshHomeData } from '../../../wailsjs/go/main/App'
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'
  import { userDataStore } from '../stores/userDataStore'
  import { t } from '../stores/i18nStore'
  import TrackList from './TrackList.svelte'

  let topPlayed: SearchResultItem[] = []
  let recommended: SearchResultItem[] = []
  let isLoading = true
  let isRefreshing = false

  onMount(async () => {
    await loadData()
  })

  async function loadData() {
    isLoading = true
    try {
      const data = await GetHomeData()
      topPlayed = data?.topPlayed || []
      recommended = data?.recommended || []
    } catch (e) {
      console.warn('Failed to load For You data:', e)
    } finally {
      isLoading = false
    }
  }

  async function handleRefresh() {
    if (isRefreshing) return
    isRefreshing = true
    try {
      const data = await RefreshHomeData()
      topPlayed = data?.topPlayed || []
      recommended = data?.recommended || []
    } catch (e) {
      console.warn('Failed to refresh For You data:', e)
    } finally {
      isRefreshing = false
    }
  }

  function handlePlayAll(tracks: SearchResultItem[], shuffle = false) {
    if (!tracks || tracks.length === 0) return
    queueStore.setQueue(tracks, 0)
    if (shuffle) {
      queueStore.toggleShuffle()
    }
    const first = queueStore.getCurrent()
    if (first) {
      playerStore.playTrack(first, undefined, true)
    }
  }
</script>

<div class="px-4 sm:px-8 py-6 space-y-8 select-none">
  <!-- Hero Section -->
  <div class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-crimson-900/35 via-[#181A26]/80 to-[#0F111A] border border-white/[0.08] p-6 sm:p-8 shadow-2xl">
    <div class="absolute -right-12 -bottom-12 w-64 h-64 bg-crimson-600/15 rounded-full blur-3xl pointer-events-none"></div>
    <div class="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2 mb-2">
          <span class="w-2 h-2 rounded-full bg-crimson-500 shadow-[0_0_8px_#E61E38] animate-pulse"></span>
          <span class="text-xs font-bold tracking-widest text-crimson-400 uppercase font-mono">Personal Mix</span>
        </div>
        <h1 class="text-2xl sm:text-3xl font-extrabold text-white tracking-tight">
          Untuk Anda
        </h1>
        <p class="text-xs sm:text-sm text-zinc-400 mt-1 max-w-lg">
          Musik yang diracik khusus berdasarkan lagu yang sering kamu putar, riwayat, dan artis favoritmu.
        </p>
      </div>

      <div class="flex items-center gap-2 shrink-0">
        <button
          class="flex items-center gap-2 px-3.5 py-2 rounded-xl bg-white/[0.06] hover:bg-white/[0.12] border border-white/10 text-xs font-semibold text-zinc-200 transition cursor-pointer disabled:opacity-50"
          on:click={handleRefresh}
          disabled={isRefreshing}
        >
          <svg class="w-3.5 h-3.5 {isRefreshing ? 'animate-spin text-crimson-400' : 'text-zinc-400'}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67"/>
          </svg>
          <span>{isRefreshing ? 'Menyegarkan...' : 'Segarkan Rekomendasi'}</span>
        </button>
      </div>
    </div>
  </div>

  {#if isLoading}
    <div class="py-16 text-center text-zinc-500 flex flex-col items-center justify-center gap-3">
      <div class="w-6 h-6 border-2 border-white/15 border-t-crimson-500 rounded-full animate-spin"></div>
      <p class="text-xs font-medium">Meracik musik spesial untukmu...</p>
    </div>
  {:else}
    <!-- Section 1: Sering Kamu Putar -->
    {#if topPlayed.length > 0}
      <div>
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-2.5">
            <div class="w-2.5 h-2.5 rounded-full bg-crimson-500 shadow-[0_0_8px_#E61E38]"></div>
            <h2 class="text-lg font-bold text-white tracking-tight">
              Sering Kamu Putar
            </h2>
          </div>
          <div class="flex items-center gap-2">
            <button
              class="px-3 py-1.5 rounded-xl bg-crimson-600 hover:bg-crimson-500 text-white font-semibold text-xs transition cursor-pointer flex items-center gap-1.5 shadow-sm shadow-crimson-600/30"
              on:click={() => handlePlayAll(topPlayed, false)}
            >
              <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24"><polygon points="5 3 19 12 5 21 5 3"/></svg>
              <span>Putar Semua</span>
            </button>
            <button
              class="px-3 py-1.5 rounded-xl bg-white/[0.06] hover:bg-white/[0.12] border border-white/10 text-zinc-300 font-semibold text-xs transition cursor-pointer flex items-center gap-1.5"
              on:click={() => handlePlayAll(topPlayed, true)}
            >
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M16 3h5v5M4 20L21 3M21 16v5h-5M15 15l6 6M4 4l5 5"/>
              </svg>
              <span>Acak</span>
            </button>
          </div>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
          {#each topPlayed as track (track.id)}
            {@const isCurrent = $playerStore.currentTrack?.id === track.id}
            {@const isPlaying = isCurrent && $playerStore.isPlaying}
            <div
              class="group relative flex items-center gap-3 p-2.5 rounded-2xl bg-white/[0.03] hover:bg-white/[0.08] border border-white/[0.06] hover:border-crimson-500/30 transition-all cursor-pointer shadow-sm overflow-hidden {isCurrent ? 'bg-crimson-500/10 border-crimson-500/30' : ''}"
              on:click={() => {
                queueStore.setQueue(topPlayed, topPlayed.findIndex(t => t.id === track.id))
                playerStore.playTrack(track)
              }}
              role="button"
              tabindex="0"
              on:keydown={(e) => e.key === 'Enter'}
            >
              <div class="relative w-12 h-12 shrink-0 rounded-xl overflow-hidden bg-black/40 shadow-xs">
                {#if track.thumbnailUrl}
                  <img src={track.thumbnailUrl} alt={track.title} class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" loading="lazy" />
                {:else}
                  <div class="w-full h-full flex items-center justify-center bg-zinc-800 text-zinc-500">
                    <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>
                  </div>
                {/if}

                <div class="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                  <div class="w-6 h-6 rounded-full bg-crimson-500 text-white flex items-center justify-center shadow-md">
                    {#if isPlaying}
                      <svg class="w-3 h-3 fill-current" viewBox="0 0 24 24"><rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/></svg>
                    {:else}
                      <svg class="w-3 h-3 fill-current ml-0.5" viewBox="0 0 24 24"><polygon points="5 3 19 12 5 21 5 3"/></svg>
                    {/if}
                  </div>
                </div>
              </div>

              <div class="min-w-0 flex-1">
                <p class="text-xs font-bold truncate {isCurrent ? 'text-crimson-400' : 'text-zinc-200 group-hover:text-white'}">
                  {track.title}
                </p>
                <p class="text-[11px] text-zinc-400 truncate mt-0.5">
                  {track.uploaderName}
                </p>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- Section 2: Rekomendasi Lagu Populer & Musisi Favorit -->
    <div>
      <TrackList
        title="Rekomendasi Spesial Untukmu"
        tracks={recommended}
        isLoading={isLoading}
        emptyMessage="Mulai putar beberapa lagu favoritmu untuk melatih rekomendasi cerdas di tab Untuk Anda ini!"
      />
    </div>
  {/if}
</div>
