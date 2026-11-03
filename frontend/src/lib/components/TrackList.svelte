<script lang="ts">
  import { createEventDispatcher, onDestroy } from 'svelte'
  import type { SearchResultItem } from '../types'
  import TrackItem from './TrackItem.svelte'
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'
  import { t } from '../stores/i18nStore'
  import HeroBanner from './HeroBanner.svelte'

  export let title = 'Daftar Lagu'
  export let tracks: SearchResultItem[] = []
  export let isLoading = false
  export let hasMore = false
  export let isLoadingMore = false
  export let emptyMessage = 'Tidak ada lagu yang ditemukan.'
  export let isHistory = false
  export let showBanner = false
  export let bannerTag = ''
  export let bannerTitle = ''
  export let bannerSubtitle = ''
  export let flush = false
  export let onRemoveItem: ((track: SearchResultItem) => void) | null = null
  export let onClearAll: (() => void) | null = null

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
  <!-- Cinematic Hero Banner (for Favorites, History, etc.) -->
  {#if showBanner || isHistory}
    <div class="px-4 sm:px-8 pt-6 pb-2">
      <HeroBanner
        tag={bannerTag || (isHistory ? "AKTIVITAS TERAKHIR • RECENTLY PLAYED" : "KOLEKSI PRIBADI • FAVORITES")}
        title={bannerTitle || title}
        subtitle={bannerSubtitle || (tracks.length > 0 ? `${tracks.length} ${$t('items_count')}${totalDurationFormatted ? ` • ${totalDurationFormatted}` : ''} • ${isHistory ? 'Dengarkan kembali alunan musik yang baru saja kamu nikmati' : 'Semua lagu favorit pilihan terbaik yang kamu sukai'}` : emptyMessage)}
        imageUrl={tracks[0]?.thumbnailUrl || ''}
        trackId={tracks[0]?.id}
        playLabel={isCurrentListPlaying ? $t('pause_all') : $t('play_all')}
        isPlaying={isCurrentListPlaying}
        showPlay={tracks.length > 0}
        showShuffle={tracks.length > 0}
        shuffleTitle={$t('shuffle')}
        on:play={handlePlayAll}
        on:shuffle={handleShuffleAll}
      >
        {#if isHistory && tracks.length > 0 && onClearAll}
          <button
            type="button"
            class="w-11 h-11 sm:w-12 sm:h-12 rounded-full bg-white/[0.08] hover:bg-red-500/20 border border-white/10 hover:border-red-500/30 text-zinc-300 hover:text-red-400 flex items-center justify-center transition-all hover:scale-105 active:scale-95 cursor-pointer shrink-0"
            on:click={onClearAll}
            title={$t('clear_history')}
          >
            <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <polyline points="3 6 5 6 21 6"></polyline>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
            </svg>
          </button>
        {/if}
      </HeroBanner>
    </div>
  {:else}
    <!-- Search Results Header -->
    <section class="{flush ? 'px-0 pt-2 pb-3' : 'px-4 sm:px-8 pt-4 sm:pt-6 pb-3 sm:pb-4'}">
      <div class="flex flex-col md:flex-row md:items-end justify-between gap-3 sm:gap-4">
        <div>
          <div class="flex items-center gap-2 mb-1.5">
            <span class="text-[11px] uppercase tracking-wider font-semibold text-crimson-400">{$t('audio_exploration')}</span>
            <span class="text-zinc-600 text-xs">•</span>
            <span class="text-[11px] text-zinc-400 font-medium">{$t('stream_engine')}</span>
          </div>
          <h1 class="text-xl sm:text-2xl md:text-3xl font-extrabold text-white tracking-tight flex items-center gap-2.5 sm:gap-3 flex-wrap">
            <span>{title}</span>
            {#if !isLoading && tracks.length > 0}
              <span class="text-xs font-mono font-medium px-2.5 py-1 rounded-full bg-white/[0.08] text-zinc-300 border border-white/10 inline-flex items-center gap-1.5 align-middle">
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
          <div class="flex items-center gap-2.5 shrink-0">
            <button
              class="h-10 px-5 rounded-full bg-gradient-to-r from-crimson-600 to-crimson-500 hover:from-crimson-500 hover:to-crimson-400 border border-transparent text-white font-semibold text-xs sm:text-sm inline-flex items-center gap-2 shadow-lg shadow-crimson-500/25 hover:shadow-crimson-500/40 transition-all transform active:scale-95 cursor-pointer whitespace-nowrap shrink-0"
              on:click={handlePlayAll}
              title={isCurrentListPlaying ? $t('pause_all') : $t('play_all')}
            >
              {#if isCurrentListPlaying}
                <svg class="w-4 h-4 fill-current shrink-0" viewBox="0 0 24 24">
                  <rect height="14" width="3.5" x="6" y="5"></rect>
                  <rect height="14" width="3.5" x="14.5" y="5"></rect>
                </svg>
                <span class="whitespace-nowrap">{$t('pause_all')}</span>
              {:else}
                <svg class="w-4 h-4 fill-current shrink-0" viewBox="0 0 24 24">
                  <polygon points="5 3 19 12 5 21 5 3"></polygon>
                </svg>
                <span class="whitespace-nowrap">{$t('play_all')}</span>
              {/if}
            </button>

            <button
              class="h-10 px-5 rounded-full bg-white/[0.06] hover:bg-white/[0.12] border border-white/10 hover:border-white/20 text-zinc-300 hover:text-white font-semibold text-xs sm:text-sm inline-flex items-center gap-2 transition-all active:scale-95 cursor-pointer whitespace-nowrap shrink-0"
              on:click={handleShuffleAll}
              title={$t('shuffle')}
            >
              <svg class="w-4 h-4 text-zinc-400 shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <polyline points="16 3 21 3 21 8"></polyline>
                <line x1="4" y1="20" x2="21" y2="3"></line>
                <polyline points="21 16 21 21 16 21"></polyline>
                <line x1="15" y1="15" x2="21" y2="21"></line>
                <line x1="4" y1="4" x2="9" y2="9"></line>
              </svg>
              <span class="whitespace-nowrap">{$t('shuffle')}</span>
            </button>
          </div>
        {/if}
      </div>
    </section>
  {/if}

  <!-- Track Table Section -->
  <section class="{flush ? 'px-0 mt-2' : 'px-3 sm:px-6 mt-2'}">
    <!-- Table Column Headers (Ramping & Rata Kiri) -->
    <div class="flex items-center gap-2.5 sm:gap-3.5 px-2.5 sm:px-3 py-2 text-[11px] font-semibold text-zinc-400 uppercase tracking-wider border-b border-white/[0.08]">
      <div class="w-8 sm:w-9 shrink-0 text-center">#</div>
      <div class="flex-1 min-w-0">{$t('title_and_artist')}</div>
      <div class="hidden lg:block w-44 xl:w-56 shrink-0 pr-3">{$t('channel_or_artist')}</div>
      <div class="w-24 sm:w-28 shrink-0 text-right pr-1">{$t('quality_and_duration')}</div>
      <div class="w-16 sm:w-20 shrink-0 text-right pr-1">{$t('actions')}</div>
    </div>

    <!-- Track List Items Container -->
    {#if isLoading}
      <div class="divide-y divide-white/[0.03] mt-2 font-medium">
        {#each Array(6) as _}
          <div class="flex items-center gap-2.5 sm:gap-3.5 px-2.5 sm:px-3 py-3">
            <div class="w-8 sm:w-9 shrink-0 flex justify-center">
              <div class="w-4 h-4 rounded bg-white/[0.05] animate-pulse"></div>
            </div>
            <div class="flex-1 min-w-0 flex items-center gap-2.5 sm:gap-3">
              <div class="w-10 h-10 sm:w-11 sm:h-11 rounded-xl bg-white/[0.05] animate-pulse shrink-0"></div>
              <div class="flex-1 space-y-1.5">
                <div class="w-3/4 h-3.5 rounded bg-white/[0.05] animate-pulse"></div>
                <div class="w-1/2 h-2.5 rounded bg-white/[0.05] animate-pulse"></div>
              </div>
            </div>
            <div class="hidden lg:block w-44 xl:w-56 shrink-0 pr-3">
              <div class="w-2/3 h-3 rounded bg-white/[0.05] animate-pulse"></div>
            </div>
            <div class="w-24 sm:w-28 shrink-0 flex justify-end">
              <div class="w-14 h-3 rounded bg-white/[0.05] animate-pulse"></div>
            </div>
            <div class="w-16 sm:w-20 shrink-0 flex justify-end">
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
          <TrackItem
            {track}
            index={i}
            tracksContext={tracks}
            onRemove={onRemoveItem ? () => onRemoveItem(track) : null}
            removeTitle={isHistory ? $t('remove_from_history') : ''}
          />
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
