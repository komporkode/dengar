<script lang="ts">
  import { createEventDispatcher, onDestroy } from 'svelte'
  import type { SearchResultItem } from '../types'
  import TrackItem from './TrackItem.svelte'
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'
  import { t } from '../stores/i18nStore'

  export let title = 'Daftar Lagu'
  export let tracks: SearchResultItem[] = []
  export let isLoading = false
  export let hasMore = false
  export let isLoadingMore = false
  export let emptyMessage = 'Tidak ada lagu yang ditemukan.'

  const dispatch = createEventDispatcher<{
    loadMore: void
  }>()

  let observer: IntersectionObserver | null = null

  function setupObserver(node: HTMLElement) {
    if (observer) observer.disconnect()
    observer = new IntersectionObserver((entries) => {
      const first = entries[0]
      if (first.isIntersecting && hasMore && !isLoadingMore && !isLoading) {
        dispatch('loadMore')
      }
    }, { rootMargin: '300px' })
    observer.observe(node)
  }

  onDestroy(() => {
    if (observer) observer.disconnect()
  })

  function handleTriggerLoadMore() {
    if (hasMore && !isLoadingMore && !isLoading) {
      dispatch('loadMore')
    }
  }

  $: isCurrentListPlaying = tracks.length > 0 &&
    $playerStore.isPlaying &&
    tracks.some(t => t.id === $playerStore.currentTrack?.id)

  $: totalDurationFormatted = formatTotalDuration(tracks)

  function formatTotalDuration(items: SearchResultItem[]): string {
    if (!items || items.length === 0) return ''
    const totalSec = items.reduce((acc, item) => acc + (item.duration || 0), 0)
    if (totalSec <= 0) return ''
    const hrs = Math.floor(totalSec / 3600)
    const mins = Math.floor((totalSec % 3600) / 60)
    if (hrs > 0) {
      return `~${hrs} jam ${mins} mnt`
    }
    return `~${mins} mnt`
  }

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
</script>

<div class="w-full flex flex-col select-none">
  <!-- Search Results Header -->
  <section class="px-4 sm:px-8 pt-4 sm:pt-6 pb-3 sm:pb-4">
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-3 sm:gap-4">
      <div>
        <div class="flex items-center gap-2 mb-1">
          <span class="text-xs uppercase tracking-widest font-semibold text-crimson-400">{$t('audio_exploration')}</span>
          <span class="text-zinc-600">•</span>
          <span class="text-xs text-zinc-400">{$t('stream_engine')}</span>
        </div>
        <h1 class="text-xl sm:text-2xl md:text-3xl font-extrabold text-white tracking-tight flex items-center gap-2 sm:gap-3 flex-wrap">
          <span>{title}</span>
          {#if !isLoading && tracks.length > 0}
            <span class="text-xs font-mono font-medium px-2.5 py-1 rounded-full bg-white/[0.08] text-zinc-300 border border-white/10 flex items-center gap-1.5">
              <span>{tracks.length} {$t('items_count')}</span>
              {#if totalDurationFormatted}
                <span class="text-zinc-500">•</span>
                <span class="text-zinc-400">{totalDurationFormatted}</span>
              {/if}
            </span>
          {/if}
        </h1>
      </div>

      <!-- Play All & Shuffle Controls (Dynamic state) -->
      {#if !isLoading && tracks.length > 0}
        <div class="flex items-center gap-2 sm:gap-2.5">
          <button
            class="px-4 sm:px-5 py-2 sm:py-2.5 rounded-full bg-crimson-500 hover:bg-crimson-600 text-white font-bold text-xs sm:text-sm flex items-center gap-2 shadow-lg shadow-crimson-500/30 transition transform active:scale-95 cursor-pointer"
            on:click={handlePlayAll}
            title={isCurrentListPlaying ? $t('pause_all') : $t('play_all')}
          >
            {#if isCurrentListPlaying}
              <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 fill-current" viewBox="0 0 24 24">
                <rect height="14" width="3.5" x="6" y="5"></rect>
                <rect height="14" width="3.5" x="14.5" y="5"></rect>
              </svg>
              <span>{$t('pause_all')}</span>
            {:else}
              <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 fill-current" viewBox="0 0 24 24">
                <polygon points="5 3 19 12 5 21 5 3"></polygon>
              </svg>
              <span>{$t('play_all')}</span>
            {/if}
          </button>

          <button
            class="px-3.5 sm:px-4 py-2 sm:py-2.5 rounded-full bg-white/[0.06] hover:bg-white/[0.1] border border-white/10 text-zinc-300 font-semibold text-xs sm:text-sm flex items-center gap-2 transition active:scale-95 cursor-pointer"
            on:click={handleShuffleAll}
            title={$t('shuffle')}
          >
            <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 text-zinc-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <polyline points="16 3 21 3 21 8"></polyline>
              <line x1="4" y1="20" x2="21" y2="3"></line>
              <polyline points="21 16 21 21 16 21"></polyline>
              <line x1="15" y1="15" x2="21" y2="21"></line>
              <line x1="4" y1="4" x2="9" y2="9"></line>
            </svg>
            <span>{$t('shuffle')}</span>
          </button>
        </div>
      {/if}
    </div>
  </section>

  <!-- Track Table Section -->
  <section class="px-3 sm:px-8 mt-2">
    <!-- Table Column Headers -->
    <div class="grid grid-cols-12 gap-2 sm:gap-4 px-3 sm:px-4 py-2.5 text-[11px] font-semibold text-zinc-400 uppercase tracking-wider border-b border-white/[0.08]">
      <div class="col-span-1 text-center">#</div>
      <div class="col-span-7 sm:col-span-6 md:col-span-5">{$t('title_and_artist')}</div>
      <div class="hidden md:block col-span-3">{$t('channel_or_artist')}</div>
      <div class="col-span-2 md:col-span-2 text-right">{$t('quality_and_duration')}</div>
      <div class="col-span-2 md:col-span-1 text-right pr-1 sm:pr-2">{$t('actions')}</div>
    </div>

    <!-- Track List Items Container -->
    {#if isLoading}
      <div class="divide-y divide-white/[0.03] mt-2 font-medium">
        {#each Array(6) as _}
          <div class="grid grid-cols-12 gap-4 px-4 py-3.5 items-center">
            <div class="col-span-1 flex justify-center">
              <div class="w-4 h-4 rounded bg-white/[0.05] animate-pulse"></div>
            </div>
            <div class="col-span-6 md:col-span-5 flex items-center gap-3.5">
              <div class="w-11 h-11 rounded-xl bg-white/[0.05] animate-pulse shrink-0"></div>
              <div class="flex-1 space-y-1.5">
                <div class="w-3/4 h-3.5 rounded bg-white/[0.05] animate-pulse"></div>
                <div class="w-1/2 h-2.5 rounded bg-white/[0.05] animate-pulse"></div>
              </div>
            </div>
            <div class="hidden md:block col-span-3">
              <div class="w-2/3 h-3 rounded bg-white/[0.05] animate-pulse"></div>
            </div>
            <div class="col-span-3 md:col-span-2 flex justify-end">
              <div class="w-16 h-3 rounded bg-white/[0.05] animate-pulse"></div>
            </div>
            <div class="col-span-2 md:col-span-1 flex justify-end">
              <div class="w-8 h-3 rounded bg-white/[0.05] animate-pulse"></div>
            </div>
          </div>
        {/each}
      </div>
    {:else if tracks.length === 0}
      <div class="flex flex-col items-center justify-center py-20 text-center gap-3">
        <div class="w-12 h-12 rounded-2xl bg-white/[0.04] border border-white/10 flex items-center justify-center text-zinc-500">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" stroke-width="1.8" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="8" y1="12" x2="16" y2="12"></line>
          </svg>
        </div>
        <p class="text-sm font-medium text-zinc-300">{emptyMessage}</p>
        <span class="text-xs text-zinc-500">{$t('empty_search_tip')}</span>
      </div>
    {:else}
      <div class="divide-y divide-white/[0.03] mt-1 font-medium">
        {#each tracks as track, i (track.id + '-' + i)}
          <TrackItem {track} index={i} tracksContext={tracks} />
        {/each}
      </div>

      <!-- Infinite Scroll & Load More Section -->
      {#if hasMore}
        <div
          use:setupObserver
          class="pt-6 pb-6 flex flex-col items-center justify-center gap-2"
        >
          {#if isLoadingMore}
            <div class="flex items-center gap-2.5 px-4 py-2.5 rounded-full bg-white/[0.05] border border-white/10 text-xs text-zinc-300 shadow-sm animate-pulse">
              <div class="w-3.5 h-3.5 border-2 border-white/20 border-t-crimson-500 rounded-full animate-spin"></div>
              <span>{$t('loading_more')}</span>
            </div>
          {:else}
            <button
              class="px-5 py-2.5 rounded-full bg-white/[0.05] hover:bg-white/[0.09] active:bg-white/[0.12] border border-white/10 text-xs font-semibold text-zinc-300 hover:text-white transition-all cursor-pointer flex items-center gap-2 shadow-sm hover:scale-[1.02] active:scale-95"
              on:click={handleTriggerLoadMore}
            >
              <svg class="w-3.5 h-3.5 text-crimson-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <polyline points="6 9 12 15 18 9"></polyline>
              </svg>
              <span>{$t('load_more')}</span>
            </button>
          {/if}
        </div>
      {/if}
    {/if}
  </section>
</div>
