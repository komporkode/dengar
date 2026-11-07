<script lang="ts">
  import { onMount } from 'svelte'
  import type { SearchResultItem } from '../types'
  import { GetOfflineTracks, RemoveOfflineTrack, GetAudioCacheStats } from '../../../wailsjs/go/main/App'
  import TrackItem from './TrackItem.svelte'
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'
  import { showToast } from '../stores/settingsStore'
  import { t } from '../stores/i18nStore'
  import { offlineStore } from '../stores/offlineStore'
  import HeroBanner from './HeroBanner.svelte'

  let tracks: SearchResultItem[] = []
  let isLoading = true
  let totalBytes = 0

  onMount(() => {
    loadOfflineTracks()
  })

  async function loadOfflineTracks() {
    isLoading = true
    try {
      await offlineStore.load()
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
      await offlineStore.removeTrack(track.id)
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
  <!-- Cinematic Full-Bleed Music Hero -->
  <div class="px-4 sm:px-8 pt-6 pb-2">
    <HeroBanner
      tag="BEBAS KUOTA • OFFLINE READY"
      title={$t('offline_heading')}
      subtitle={`${tracks.length} ${$t('offline_songs_count')} • ${formatBytes(totalBytes)}${formatTotalDuration(tracks) ? ` • ${formatTotalDuration(tracks)}` : ''} • ${$t('offline_subheading')}`}
      imageUrl={tracks[0]?.thumbnailUrl || ''}
      trackId={tracks[0]?.id}
      playLabel={isCurrentListPlaying ? $t('offline_pause') : $t('offline_play_all')}
      isPlaying={isCurrentListPlaying}
      showPlay={tracks.length > 0}
      showShuffle={tracks.length > 0}
      shuffleTitle={$t('offline_shuffle')}
      on:play={handlePlayAll}
      on:shuffle={handleShuffleAll}
    >
      <button
        type="button"
        class="w-11 h-11 sm:w-12 sm:h-12 rounded-full bg-white/[0.08] hover:bg-white/[0.14] border border-white/10 hover:border-white/20 text-zinc-200 hover:text-white flex items-center justify-center transition-all hover:scale-105 active:scale-95 cursor-pointer disabled:opacity-50 shrink-0"
        on:click={loadOfflineTracks}
        disabled={isLoading}
        title={$t('offline_refresh')}
        aria-label={$t('offline_refresh')}
      >
        <svg
          class="w-4 h-4 {isLoading ? 'animate-spin text-crimson-400' : 'text-zinc-200'} shrink-0"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67" />
        </svg>
      </button>
    </HeroBanner>
  </div>

  <!-- Track List Body -->
  <section class="px-3 sm:px-6 py-4">
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
      <!-- Header Bar (Ramping & Rata Kiri) -->
      <div class="flex items-center gap-2.5 sm:gap-3.5 px-2.5 sm:px-3 py-2 text-[11px] font-bold tracking-wider text-zinc-400 uppercase border-b border-white/[0.06] mb-1">
        <div class="w-8 sm:w-9 shrink-0 text-center">#</div>
        <div class="flex-1 min-w-0">{$t('title_and_artist')}</div>
        <div class="hidden lg:block w-44 xl:w-56 shrink-0 pr-3">{$t('channel_or_artist')}</div>
        <div class="w-24 sm:w-28 shrink-0 text-right pr-1">{$t('quality_and_duration')}</div>
        <div class="w-16 sm:w-20 shrink-0 text-right pr-1">{$t('actions')}</div>
      </div>

      <div class="space-y-1">
        {#each tracks as track, index (track.id)}
          <TrackItem
            {track}
            {index}
            tracksContext={tracks}
            isOffline={true}
            onRemove={() => handleRemoveTrack(track)}
            removeTitle={$t('offline_delete_title')}
          />
        {/each}
      </div>
    {/if}
  </section>
</div>
