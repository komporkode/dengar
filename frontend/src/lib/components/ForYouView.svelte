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

  // Dynamic greeting based on user's local hour
  function getGreeting(): string {
    const hour = new Date().getHours()
    if (hour >= 4 && hour < 11) return 'Selamat Pagi'
    if (hour >= 11 && hour < 15) return 'Selamat Siang'
    if (hour >= 15 && hour < 19) return 'Selamat Sore'
    return 'Selamat Malam'
  }

  // Combined deduplicated mix pool for 1-click personal playback
  $: mixTracks = [...topPlayed, ...recommended].filter(
    (item, idx, arr) => arr.findIndex(t => t.id === item.id) === idx
  )

  // Top 4 artwork thumbnails for the vinyl/mosaic cover showcase
  $: sampleArtworks = (mixTracks.length > 0 ? mixTracks : topPlayed)
    .filter(t => !!t.thumbnailUrl)
    .slice(0, 4)

  $: isMixPlaying = $playerStore.isPlaying && mixTracks.some(t => t.id === $playerStore.currentTrack?.id)

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

  function handlePlayMix(shuffle = false) {
    const pool = mixTracks.length > 0 ? mixTracks : (topPlayed.length > 0 ? topPlayed : recommended)
    if (pool.length === 0) return
    queueStore.setQueue(pool, 0)
    if (shuffle) {
      queueStore.enableShuffle()
    }
    const current = queueStore.getCurrent()
    if (current) {
      playerStore.playTrack(current)
    }
  }

  function handlePlaySection(tracks: SearchResultItem[], shuffle = false) {
    if (!tracks || tracks.length === 0) return
    queueStore.setQueue(tracks, 0)
    if (shuffle) {
      queueStore.enableShuffle()
    }
    const current = queueStore.getCurrent()
    if (current) {
      playerStore.playTrack(current)
    }
  }
</script>

<div class="px-4 sm:px-8 py-6 space-y-8 select-none">
  <!-- Ultra-Premium Music Hero: Apple Music / Spotify Studio Style -->
  <div class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-[#1C1024] via-[#121320] to-[#0A0C14] border border-white/[0.08] p-6 sm:p-8 md:p-10 shadow-2xl transition-all">
    <!-- Top-edge luminous accent highlight -->
    <div class="absolute inset-x-0 top-0 h-[1px] bg-gradient-to-r from-transparent via-crimson-400/40 to-transparent pointer-events-none"></div>

    <!-- Radiant Ambient Light Orbs -->
    <div class="absolute -right-16 -top-16 w-96 h-96 bg-crimson-600/20 rounded-full blur-[110px] pointer-events-none"></div>
    <div class="absolute -left-16 -bottom-16 w-80 h-80 bg-indigo-600/15 rounded-full blur-[100px] pointer-events-none"></div>
    <div class="absolute right-1/3 bottom-0 w-64 h-64 bg-violet-600/10 rounded-full blur-[90px] pointer-events-none"></div>

    <div class="relative z-10 flex flex-col lg:flex-row lg:items-center justify-between gap-8">
      <!-- Left Column: Curation Info & Action Suite -->
      <div class="flex-1 max-w-2xl space-y-4 sm:space-y-5">
        <!-- Equalizer & Badge Pill -->
        <div class="inline-flex items-center gap-2.5 px-3 py-1.5 rounded-full bg-crimson-500/10 border border-crimson-500/25 backdrop-blur-md">
          <!-- Animated mini equalizer wave -->
          <div class="flex items-end gap-[3px] h-3.5 w-3.5" aria-hidden="true">
            <span class="w-[2.5px] bg-crimson-400 rounded-full sound-bar-1 h-2"></span>
            <span class="w-[2.5px] bg-crimson-400 rounded-full sound-bar-2 h-3.5"></span>
            <span class="w-[2.5px] bg-crimson-400 rounded-full sound-bar-3 h-1.5"></span>
            <span class="w-[2.5px] bg-crimson-400 rounded-full sound-bar-4 h-3"></span>
          </div>
          <span class="text-[11px] font-bold tracking-[0.18em] text-crimson-400 uppercase font-mono">Personal Mix</span>
          <span class="text-white/20 text-xs">•</span>
          <span class="text-[10px] font-semibold text-zinc-300 tracking-wider uppercase">Dirancang Untukmu</span>
        </div>

        <!-- Headline & Dynamic Greeting -->
        <div>
          <div class="flex items-center gap-2 text-zinc-400 text-xs sm:text-sm font-medium mb-1">
            <span>{getGreeting()}</span>
            <span class="text-zinc-600">•</span>
            <span class="text-zinc-400">Pilihan Musik Cerdas</span>
          </div>
          <h1 class="text-2xl sm:text-4xl md:text-5xl font-black text-white tracking-tight leading-tight">
            Untuk Anda
          </h1>
          <p class="text-xs sm:text-sm text-zinc-300/90 mt-2 leading-relaxed max-w-xl">
            Racikan spesial musik favorit, riwayat putar, dan penemuan artis baru yang disesuaikan dengan selera dengarmu hari ini.
          </p>
        </div>

        <!-- Controls Suite: Putar Mix, Putar Acak, Segarkan (Uniform Heights & Non-wrapping) -->
        <div class="flex flex-wrap items-center gap-3 pt-1">
          <button
            class="h-10 sm:h-11 px-5 sm:px-6 rounded-full bg-gradient-to-r from-crimson-600 to-crimson-500 hover:from-crimson-500 hover:to-crimson-400 border border-transparent text-white font-semibold text-xs sm:text-sm inline-flex items-center gap-2 shadow-lg shadow-crimson-500/30 hover:shadow-crimson-500/50 hover:scale-[1.02] active:scale-95 transition-all cursor-pointer whitespace-nowrap shrink-0"
            on:click={() => handlePlayMix(false)}
            title="Putar semua rekomendasi untuk Anda"
          >
            {#if isMixPlaying}
              <svg class="w-4 h-4 fill-current shrink-0" viewBox="0 0 24 24">
                <rect height="14" width="3.5" x="6" y="5"></rect>
                <rect height="14" width="3.5" x="14.5" y="5"></rect>
              </svg>
              <span class="whitespace-nowrap">Jeda Mix</span>
            {:else}
              <svg class="w-4 h-4 fill-current shrink-0" viewBox="0 0 24 24">
                <polygon points="5 3 19 12 5 21 5 3"></polygon>
              </svg>
              <span class="whitespace-nowrap">Putar Mix</span>
            {/if}
          </button>

          <button
            class="h-10 sm:h-11 px-4 sm:px-5 rounded-full bg-white/[0.08] hover:bg-white/[0.14] border border-white/10 hover:border-white/20 text-zinc-200 hover:text-white font-semibold text-xs sm:text-sm inline-flex items-center gap-2 transition-all active:scale-95 cursor-pointer whitespace-nowrap shrink-0"
            on:click={() => handlePlayMix(true)}
            title="Putar acak personal mix"
          >
            <svg class="w-4 h-4 text-zinc-300 shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <polyline points="16 3 21 3 21 8"></polyline>
              <line x1="4" y1="20" x2="21" y2="3"></line>
              <polyline points="21 16 21 21 16 21"></polyline>
              <line x1="15" y1="15" x2="21" y2="21"></line>
              <line x1="4" y1="4" x2="9" y2="9"></line>
            </svg>
            <span class="whitespace-nowrap">Putar Acak</span>
          </button>

          <button
            class="h-10 sm:h-11 px-4 sm:px-5 rounded-full bg-white/[0.04] hover:bg-white/[0.09] border border-white/10 hover:border-white/20 text-zinc-300 hover:text-white font-medium text-xs sm:text-sm inline-flex items-center gap-2 transition-all active:scale-95 cursor-pointer disabled:opacity-50 whitespace-nowrap shrink-0"
            on:click={handleRefresh}
            disabled={isRefreshing}
            title="Segarkan kurasi musik"
          >
            <svg class="w-3.5 h-3.5 {isRefreshing ? 'animate-spin text-crimson-400' : 'text-zinc-400'} shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67"/>
            </svg>
            <span class="whitespace-nowrap">{isRefreshing ? 'Menyegarkan...' : 'Segarkan'}</span>
          </button>
        </div>

        <!-- Meta Info Pills -->
        <div class="flex items-center gap-4 pt-1 text-[11px] text-zinc-400 font-medium flex-wrap">
          <span class="inline-flex items-center gap-1.5">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
            <span>{mixTracks.length > 0 ? `${mixTracks.length} Trek Pilihan` : 'Kurasi Otomatis'}</span>
          </span>
          <span class="text-zinc-600">•</span>
          <span class="inline-flex items-center gap-1.5">
            <svg class="w-3.5 h-3.5 text-crimson-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
            <span>Lossless Hi-Fi</span>
          </span>
          <span class="text-zinc-600">•</span>
          <span class="text-zinc-500">Bebas Iklan Selamanya</span>
        </div>
      </div>

      <!-- Right Column: Visual Artwork Mosaic / Vinyl Cover Stage -->
      <div class="shrink-0 flex items-center justify-center lg:justify-end">
        <div class="relative group cursor-pointer" on:click={() => handlePlayMix(false)} role="button" tabindex="0" on:keydown={(e) => e.key === 'Enter' && handlePlayMix(false)}>
          <!-- Outer Ambient Glow behind cover -->
          <div class="absolute -inset-3 bg-gradient-to-tr from-crimson-600/40 via-violet-600/30 to-transparent rounded-3xl blur-2xl opacity-70 group-hover:opacity-100 transition-opacity duration-500"></div>

          <!-- Mosaic Card Shell -->
          <div class="relative w-44 h-44 sm:w-52 sm:h-52 rounded-2xl sm:rounded-3xl overflow-hidden bg-black/60 border border-white/15 shadow-2xl p-1.5 group-hover:scale-[1.02] transition-transform duration-300">
            {#if sampleArtworks.length >= 4}
              <!-- 2x2 Album Mosaic Grid -->
              <div class="w-full h-full grid grid-cols-2 grid-rows-2 gap-1 rounded-xl sm:rounded-2xl overflow-hidden bg-zinc-900">
                {#each sampleArtworks.slice(0, 4) as item}
                  <img src={item.thumbnailUrl} alt={item.title} class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" loading="lazy" />
                {/each}
              </div>
            {:else if sampleArtworks.length > 0}
              <!-- Single/Dual Artwork with vinyl disc peek -->
              <div class="w-full h-full relative rounded-xl sm:rounded-2xl overflow-hidden bg-zinc-900">
                <img src={sampleArtworks[0].thumbnailUrl} alt={sampleArtworks[0].title} class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" loading="lazy" />
              </div>
            {:else}
              <!-- Vinyl Graphic Placeholder with Sleek Grooves -->
              <div class="w-full h-full rounded-xl sm:rounded-2xl bg-gradient-to-br from-zinc-900 via-black to-zinc-950 flex items-center justify-center relative overflow-hidden">
                <div class="absolute inset-4 rounded-full border border-white/5"></div>
                <div class="absolute inset-8 rounded-full border border-white/5"></div>
                <div class="absolute inset-12 rounded-full border border-white/5"></div>
                <div class="w-14 h-14 rounded-full bg-gradient-to-tr from-crimson-600 to-crimson-400 flex items-center justify-center shadow-lg shadow-crimson-500/40">
                  <svg class="w-6 h-6 text-white ml-0.5 fill-current" viewBox="0 0 24 24"><polygon points="5 3 19 12 5 21 5 3"/></svg>
                </div>
              </div>
            {/if}

            <!-- Center Hover / Play Button Overlay -->
            <div class="absolute inset-0 flex items-center justify-center bg-black/35 backdrop-blur-[2px] opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-2xl sm:rounded-3xl">
              <div class="w-12 h-12 rounded-full bg-crimson-500 text-white flex items-center justify-center shadow-xl shadow-crimson-500/50 transform group-hover:scale-110 transition-transform">
                {#if isMixPlaying}
                  <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24"><rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/></svg>
                {:else}
                  <svg class="w-5 h-5 fill-current ml-0.5" viewBox="0 0 24 24"><polygon points="5 3 19 12 5 21 5 3"/></svg>
                {/if}
              </div>
            </div>
          </div>
        </div>
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
          <div class="flex items-center gap-2 shrink-0">
            <button
              class="h-9 px-4 rounded-full bg-crimson-600 hover:bg-crimson-500 border border-transparent text-white font-semibold text-xs transition cursor-pointer flex items-center gap-1.5 shadow-sm shadow-crimson-600/30 whitespace-nowrap shrink-0"
              on:click={() => handlePlaySection(topPlayed, false)}
            >
              <svg class="w-3.5 h-3.5 fill-current shrink-0" viewBox="0 0 24 24"><polygon points="5 3 19 12 5 21 5 3"/></svg>
              <span class="whitespace-nowrap">Putar Semua</span>
            </button>
            <button
              class="h-9 px-4 rounded-full bg-white/[0.06] hover:bg-white/[0.12] border border-white/10 text-zinc-300 hover:text-white font-semibold text-xs transition cursor-pointer flex items-center gap-1.5 whitespace-nowrap shrink-0"
              on:click={() => handlePlaySection(topPlayed, true)}
            >
              <svg class="w-3.5 h-3.5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M16 3h5v5M4 20L21 3M21 16v5h-5M15 15l6 6M4 4l5 5"/>
              </svg>
              <span class="whitespace-nowrap">Putar Acak</span>
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

<style>
  @keyframes soundwave {
    0%, 100% { height: 4px; }
    50% { height: 14px; }
  }

  .sound-bar-1 {
    animation: soundwave 1.1s ease-in-out infinite;
  }
  .sound-bar-2 {
    animation: soundwave 0.8s ease-in-out infinite 0.2s;
  }
  .sound-bar-3 {
    animation: soundwave 1.4s ease-in-out infinite 0.4s;
  }
  .sound-bar-4 {
    animation: soundwave 0.95s ease-in-out infinite 0.1s;
  }
</style>
