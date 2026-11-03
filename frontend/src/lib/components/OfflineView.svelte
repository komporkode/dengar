<script lang="ts">
  import { onMount } from 'svelte'
  import type { SearchResultItem } from '../types'
  import { GetOfflineTracks, RemoveOfflineTrack, GetAudioCacheStats } from '../../../wailsjs/go/main/App'
  import TrackItem from './TrackItem.svelte'
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'
  import { showToast } from '../stores/settingsStore'
  import { t } from '../stores/i18nStore'

  let tracks: SearchResultItem[] = []
  let isLoading = true
  let totalBytes = 0

  onMount(() => {
    loadOfflineTracks()
  })

  async function loadOfflineTracks() {
    isLoading = true
    try {
      tracks = (await GetOfflineTracks()) || []
      const stats = await GetAudioCacheStats()
      if (stats) {
        totalBytes = stats.totalBytes
      }
    } catch (e) {
      console.error('Failed to load offline tracks:', e)
    } finally {
      isLoading = false
    }
  }

  function formatBytes(bytes: number): string {
    if (bytes <= 0) return '0 MB'
    const mb = bytes / (1024 * 1024)
    if (mb >= 1024) {
      return `${(mb / 1024).toFixed(1)} GB`
    }
    return `${mb.toFixed(1)} MB`
  }

  function formatTotalDuration(items: SearchResultItem[]): string {
    if (!items || items.length === 0) return ''
    const totalSec = items.reduce((acc, item) => acc + (item.duration || 0), 0)
    if (totalSec <= 0) return ''
    const hrs = Math.floor(totalSec / 3600)
    const mins = Math.floor((totalSec % 3600) / 60)
    if (hrs > 0) return `~${hrs}h ${mins}m`
    return `~${mins}m`
  }

  $: isCurrentListPlaying = tracks.length > 0 &&
    $playerStore.isPlaying &&
    tracks.some(t => t.id === $playerStore.currentTrack?.id)

  function handlePlayAll() {
    if (tracks.length === 0) return
    if (isCurrentListPlaying) {
      playerStore.togglePlay()
    } else {
      queueStore.setQueue(tracks, 0)
      playerStore.playTrack(tracks[0])
    }
  }

  function handleShuffleAll() {
    if (tracks.length === 0) return
    queueStore.enableShuffle()
    queueStore.setQueue(tracks, 0)
    const current = queueStore.getCurrent()
    if (current) {
      playerStore.playTrack(current)
    }
  }

  async function handleRemoveTrack(track: SearchResultItem) {
    try {
      await RemoveOfflineTrack(track.id)
      tracks = tracks.filter(t => t.id !== track.id)
      showToast(`"${track.title}" ${$t('offline_delete_toast')}`)
      const stats = await GetAudioCacheStats()
      if (stats) {
        totalBytes = stats.totalBytes
      }
    } catch (e) {
      showToast($t('offline_delete_fail'))
    }
  }
</script>

<div class="w-full flex flex-col select-none">
  <!-- Hero Header -->
  <section class="px-4 sm:px-8 pt-6 pb-6 bg-gradient-to-b from-crimson-950/25 via-crimson-900/10 to-transparent border-b border-white/[0.06]">
    <div class="flex flex-col sm:flex-row items-center sm:items-end gap-6">
      <!-- Icon Cover Box -->
      <div class="w-36 h-36 sm:w-44 sm:h-44 shrink-0 rounded-2xl shadow-2xl border border-white/10 bg-gradient-to-tr from-[#161824] to-[#252839] flex items-center justify-center relative group">
        <div class="w-20 h-20 rounded-2xl bg-crimson-500/20 border border-crimson-500/30 text-crimson-400 flex items-center justify-center shadow-lg shadow-crimson-500/20">
          <svg class="w-10 h-10" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="7 10 12 15 17 10"></polyline>
            <line x1="12" y1="15" x2="12" y2="3"></line>
          </svg>
        </div>
      </div>

      <!-- Info Details -->
      <div class="flex flex-col gap-2 min-w-0 text-center sm:text-left flex-1">
        <div class="flex items-center justify-center sm:justify-start gap-2">
          <span class="text-[10px] font-bold uppercase tracking-widest text-crimson-400 bg-crimson-500/15 border border-crimson-500/25 px-2 py-0.5 rounded-md">
            {$t('offline_local_storage')}
          </span>
          <span class="text-xs text-zinc-400">{$t('offline_quota_free')}</span>
        </div>

        <h1 class="text-2xl sm:text-4xl font-extrabold text-white tracking-tight truncate">
          {$t('offline_heading')}
        </h1>

        <p class="text-xs sm:text-sm text-zinc-400 max-w-xl">
          {$t('offline_subheading')}
        </p>

        <!-- Stats Meta -->
        <div class="flex items-center justify-center sm:justify-start gap-3 mt-1 text-xs text-zinc-400">
          <span class="font-semibold text-white">{tracks.length} {$t('offline_songs_count')}</span>
          {#if formatTotalDuration(tracks)}
            <span>•</span>
            <span>{formatTotalDuration(tracks)}</span>
          {/if}
          <span>•</span>
          <span class="font-mono text-zinc-300 bg-white/[0.06] px-2 py-0.5 rounded-md">{formatBytes(totalBytes)} {$t('offline_used')}</span>
        </div>
      </div>
    </div>

    <!-- Action Bar -->
    <div class="flex items-center gap-3 mt-6">
      <button
        class="h-11 px-6 rounded-xl bg-crimson-500 hover:bg-crimson-600 active:scale-95 text-white text-xs sm:text-sm font-bold flex items-center gap-2 shadow-lg shadow-crimson-500/25 transition cursor-pointer disabled:opacity-50"
        on:click={handlePlayAll}
        disabled={tracks.length === 0}
      >
        {#if isCurrentListPlaying}
          <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
            <rect x="6" y="4" width="4" height="16" rx="1"></rect>
            <rect x="14" y="4" width="4" height="16" rx="1"></rect>
          </svg>
          <span>{$t('offline_pause')}</span>
        {:else}
          <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
            <polygon points="5 3 19 12 5 21 5 3"></polygon>
          </svg>
          <span>{$t('offline_play_all')}</span>
        {/if}
      </button>

      <button
        class="h-11 px-5 rounded-xl bg-white/[0.08] hover:bg-white/[0.12] active:scale-95 text-white text-xs sm:text-sm font-semibold flex items-center gap-2 border border-white/10 transition cursor-pointer disabled:opacity-50"
        on:click={handleShuffleAll}
        disabled={tracks.length === 0}
      >
        <svg class="w-4 h-4 text-zinc-300" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <polyline points="16 3 21 3 21 8"></polyline>
          <line x1="4" y1="20" x2="21" y2="3"></line>
          <polyline points="21 16 21 21 16 21"></polyline>
          <line x1="15" y1="15" x2="21" y2="21"></line>
          <line x1="4" y1="4" x2="9" y2="9"></line>
        </svg>
        <span>{$t('offline_shuffle')}</span>
      </button>

      <button
        class="h-11 px-4 rounded-xl bg-white/[0.04] hover:bg-white/[0.08] text-zinc-400 hover:text-white transition cursor-pointer border border-white/[0.06] ml-auto flex items-center gap-1.5 text-xs"
        on:click={loadOfflineTracks}
        title="{$t('offline_refresh')}"
      >
        <svg class="w-3.5 h-3.5 {isLoading ? 'animate-spin' : ''}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67"/>
        </svg>
        <span>{$t('offline_refresh')}</span>
      </button>
    </div>
  </section>

  <!-- Track List Body -->
  <section class="px-4 sm:px-8 py-6">
    {#if isLoading}
      <div class="flex items-center justify-center py-20">
        <span class="w-8 h-8 border-3 border-crimson-500/30 border-t-crimson-500 rounded-full animate-spin"></span>
      </div>
    {:else if tracks.length === 0}
      <div class="flex flex-col items-center justify-center py-20 text-center max-w-md mx-auto">
        <div class="w-16 h-16 rounded-2xl bg-white/[0.04] border border-white/[0.08] flex items-center justify-center text-zinc-500 mb-4">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="7 10 12 15 17 10"></polyline>
            <line x1="12" y1="15" x2="12" y2="3"></line>
          </svg>
        </div>
        <h3 class="text-base font-bold text-white mb-1">{$t('offline_empty_title')}</h3>
        <p class="text-xs text-zinc-400 leading-relaxed mb-4">
          {$t('offline_empty_desc')}
        </p>
      </div>
    {:else}
      <div class="space-y-1">
        {#each tracks as track, index (track.id)}
          <div class="group relative flex items-center rounded-xl hover:bg-white/[0.04] transition-colors pr-2">
            <div class="flex-1 min-w-0">
              <TrackItem {track} {index} />
            </div>
            <button
              class="p-2 text-zinc-500 hover:text-crimson-400 rounded-lg hover:bg-white/[0.06] transition opacity-0 group-hover:opacity-100 cursor-pointer shrink-0"
              on:click={() => handleRemoveTrack(track)}
              title="{$t('offline_delete_title')}"
              aria-label="{$t('offline_delete_title')}"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <polyline points="3 6 5 6 21 6"></polyline>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                <line x1="10" y1="11" x2="10" y2="17"></line>
                <line x1="14" y1="11" x2="14" y2="17"></line>
              </svg>
            </button>
          </div>
        {/each}
      </div>
    {/if}
  </section>
</div>
