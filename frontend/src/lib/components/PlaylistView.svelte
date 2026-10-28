<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import type { Playlist, SearchResultItem } from '../types'
  import TrackItem from './TrackItem.svelte'
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'
  import { userDataStore } from '../stores/userDataStore'
  import { openEditPlaylistModal, openSharePlaylistModal } from '../stores/playlistStore'
  import { showToast } from '../stores/settingsStore'
  import { t } from '../stores/i18nStore'

  export let playlist: Playlist

  const dispatch = createEventDispatcher<{
    close: void
  }>()

  $: tracks = playlist.tracks || []
  $: totalDurationFormatted = formatTotalDuration(tracks)

  $: thumbnails = getCollageThumbnails(tracks)

  function getCollageThumbnails(items: SearchResultItem[]): string[] {
    const list: string[] = []
    for (const t of items) {
      if (t.thumbnailUrl && !list.includes(t.thumbnailUrl)) {
        list.push(t.thumbnailUrl)
      }
      if (list.length >= 4) break
    }
    return list
  }

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

  async function handleDeletePlaylist() {
    if (window.confirm($t('confirm_delete_playlist'))) {
      await userDataStore.deletePlaylist(playlist.id)
      showToast('Playlist berhasil dihapus')
      dispatch('close')
    }
  }

  async function handleRemoveTrack(trackId: string) {
    await userDataStore.removeTrackFromPlaylist(playlist.id, trackId)
    showToast('Lagu dihapus dari playlist')
  }
</script>

<div class="w-full flex flex-col select-none">
  <!-- Spotify-Style Hero Header -->
  <section class="px-4 sm:px-8 pt-6 pb-6 bg-gradient-to-b from-crimson-950/20 to-transparent border-b border-white/[0.06]">
    <div class="flex flex-col sm:flex-row items-center sm:items-end gap-6">
      <!-- 2x2 Auto-Collage Cover Box -->
      <div class="w-40 h-40 sm:w-48 sm:h-48 shrink-0 rounded-2xl overflow-hidden shadow-2xl border border-white/10 bg-[#161824] flex items-center justify-center relative group">
        {#if thumbnails.length >= 4}
          <div class="w-full h-full grid grid-cols-2 grid-rows-2">
            {#each thumbnails.slice(0, 4) as thumb}
              <img src={thumb} alt="Cover" class="w-full h-full object-cover" />
            {/each}
          </div>
        {:else if thumbnails.length > 0}
          <img src={thumbnails[0]} alt="Cover" class="w-full h-full object-cover" />
        {:else}
          <div class="w-full h-full bg-gradient-to-tr from-[#1B1E2E] to-[#252A40] flex items-center justify-center">
            <svg class="w-16 h-16 text-zinc-600" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
            </svg>
          </div>
        {/if}
      </div>

      <!-- Playlist Metadata & Title -->
      <div class="flex-1 min-w-0 flex flex-col items-center sm:items-start text-center sm:text-left">
        <span class="text-[11px] font-bold tracking-widest uppercase text-crimson-400 mb-1">PLAYLIST</span>
        <h1 class="text-2xl sm:text-4xl font-extrabold text-white tracking-tight truncate max-w-full">
          {playlist.name}
        </h1>
        {#if playlist.description}
          <p class="text-xs sm:text-sm text-zinc-400 mt-1.5 line-clamp-2 max-w-xl">
            {playlist.description}
          </p>
        {/if}

        <div class="flex items-center gap-2 mt-3 text-xs text-zinc-400 font-medium">
          <span class="text-zinc-200 font-semibold">{tracks.length} {$t('items_count')}</span>
          {#if totalDurationFormatted}
            <span>•</span>
            <span>{totalDurationFormatted}</span>
          {/if}
        </div>
      </div>
    </div>

    <!-- Actions Row -->
    <div class="flex items-center gap-3 mt-6">
      <button
        class="h-11 px-6 rounded-full bg-crimson-500 hover:bg-crimson-600 text-white font-semibold text-xs tracking-wide flex items-center gap-2 shadow-lg shadow-crimson-500/25 transition cursor-pointer hover:scale-105 active:scale-95 disabled:opacity-40"
        disabled={tracks.length === 0}
        on:click={handlePlayAll}
      >
        <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
          {#if isCurrentListPlaying}
            <rect height="16" width="4" x="6" y="4"></rect>
            <rect height="16" width="4" x="14" y="4"></rect>
          {:else}
            <polygon points="5 3 19 12 5 21 5 3"></polygon>
          {/if}
        </svg>
        <span>{isCurrentListPlaying ? $t('pause_all') : $t('play_all')}</span>
      </button>

      <button
        class="h-11 w-11 rounded-full bg-white/[0.08] hover:bg-white/[0.12] text-zinc-200 flex items-center justify-center transition cursor-pointer hover:scale-105 active:scale-95 disabled:opacity-40"
        disabled={tracks.length === 0}
        on:click={handleShuffleAll}
        title={$t('shuffle')}
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <polyline points="16 3 21 3 21 8"></polyline>
          <line x1="4" y1="20" x2="21" y2="3"></line>
          <polyline points="21 16 21 21 16 21"></polyline>
          <line x1="15" y1="15" x2="21" y2="21"></line>
          <line x1="4" y1="4" x2="9" y2="9"></line>
        </svg>
      </button>

      <button
        class="h-11 w-11 rounded-full bg-white/[0.08] hover:bg-white/[0.12] text-zinc-200 flex items-center justify-center transition cursor-pointer hover:scale-105 active:scale-95"
        on:click={() => openSharePlaylistModal(playlist)}
        title={$t('share_playlist')}
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <circle cx="18" cy="5" r="3"></circle>
          <circle cx="6" cy="12" r="3"></circle>
          <circle cx="18" cy="19" r="3"></circle>
          <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"></line>
          <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"></line>
        </svg>
      </button>

      <button
        class="h-11 w-11 rounded-full bg-white/[0.08] hover:bg-white/[0.12] text-zinc-200 flex items-center justify-center transition cursor-pointer hover:scale-105 active:scale-95"
        on:click={() => openEditPlaylistModal(playlist)}
        title={$t('edit_playlist')}
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
        </svg>
      </button>

      <button
        class="h-11 w-11 rounded-full bg-white/[0.08] hover:bg-red-500/20 text-zinc-400 hover:text-red-400 flex items-center justify-center transition cursor-pointer hover:scale-105 active:scale-95 ml-auto"
        on:click={handleDeletePlaylist}
        title={$t('delete_playlist')}
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <polyline points="3 6 5 6 21 6"></polyline>
          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
        </svg>
      </button>
    </div>
  </section>

  <!-- Tracks Section -->
  <section class="px-4 sm:px-8 pt-4 pb-28">
    {#if tracks.length === 0}
      <div class="py-16 flex flex-col items-center justify-center text-center">
        <div class="w-14 h-14 rounded-2xl bg-white/[0.04] border border-white/[0.08] flex items-center justify-center mb-3 text-zinc-500">
          <svg class="w-7 h-7" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
        </div>
        <h4 class="text-sm font-bold text-zinc-300">{$t('playlist_empty')}</h4>
        <p class="text-xs text-zinc-500 mt-1 max-w-sm">{$t('playlist_empty_tip')}</p>
      </div>
    {:else}
      <!-- Header Bar -->
      <div class="grid grid-cols-12 gap-2 sm:gap-4 px-2.5 sm:px-4 py-2 text-[11px] font-bold tracking-wider text-zinc-400 uppercase border-b border-white/[0.06] mb-1">
        <div class="col-span-1 text-center">#</div>
        <div class="col-span-6 sm:col-span-6 md:col-span-5">{$t('title_and_artist')}</div>
        <div class="hidden md:block col-span-2 pr-2">{$t('channel_or_artist')}</div>
        <div class="col-span-2 md:col-span-2 text-right whitespace-nowrap">{$t('quality_and_duration')}</div>
        <div class="col-span-3 sm:col-span-3 md:col-span-2 text-right pr-1 sm:pr-2">{$t('actions')}</div>
      </div>

      <!-- Tracks List -->
      <div class="space-y-1">
        {#each tracks as track, idx (track.id)}
          <TrackItem
            {track}
            index={idx}
            tracksContext={tracks}
            onRemove={() => handleRemoveTrack(track.id)}
          />
        {/each}
      </div>
    {/if}
  </section>
</div>
