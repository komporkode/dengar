<script lang="ts">
  import type { SearchResultItem } from '../types'
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'
  import { userDataStore } from '../stores/userDataStore'

  export let track: SearchResultItem
  export let index: number

  $: isCurrent = $playerStore.currentTrack?.id === track.id
  $: isPlaying = isCurrent && $playerStore.isPlaying
  $: isFav = $userDataStore.favorites.some(f => f.id === track.id)

  function handlePlay() {
    if (isCurrent) {
      playerStore.togglePlay()
    } else {
      playerStore.playTrack(track)
      queueStore.addToQueue(track)
    }
  }

  function handleAddToQueue(e: MouseEvent) {
    e.stopPropagation()
    queueStore.addToQueue(track)
  }

  function handleToggleFav(e: MouseEvent) {
    e.stopPropagation()
    userDataStore.toggleFavorite(track)
  }
</script>

<div
  class="track-row"
  class:current={isCurrent}
  on:click={handlePlay}
  role="button"
  tabindex="0"
  on:keydown={(e) => e.key === 'Enter' && handlePlay()}
>
  <div class="index-col">
    {#if isPlaying}
      <div class="equalizer">
        <span></span>
        <span></span>
        <span></span>
      </div>
    {:else}
      <span class="index-num">{index + 1}</span>
      <button class="hover-play-btn" on:click|stopPropagation={handlePlay}>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
          <polygon points="5 3 19 12 5 21 5 3"></polygon>
        </svg>
      </button>
    {/if}
  </div>

  <div class="thumb-col">
    {#if track.thumbnailUrl}
      <img src={track.thumbnailUrl} alt={track.title} loading="lazy" />
    {:else}
      <div class="thumb-placeholder">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
          <path d="M9 18V5l12-2v13"></path>
        </svg>
      </div>
    {/if}
  </div>

  <div class="meta-col">
    <div class="track-title text-truncate" class:active-title={isCurrent}>
      {track.title}
    </div>
    <div class="track-artist text-truncate text-secondary">
      {track.uploaderName}
    </div>
  </div>

  <div class="actions-col">
    <button
      class="action-btn fav-btn"
      class:is-fav={isFav}
      on:click={handleToggleFav}
      title={isFav ? 'Hapus dari favorit' : 'Tambah ke favorit'}
    >
      <svg width="16" height="16" viewBox="0 0 24 24" fill={isFav ? '#ff0000' : 'none'} stroke="currentColor" stroke-width="2">
        <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
      </svg>
    </button>

    <button
      class="action-btn"
      on:click={handleAddToQueue}
      title="Tambah ke Antrean"
    >
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <line x1="12" y1="5" x2="12" y2="19"></line>
        <line x1="5" y1="12" x2="19" y2="12"></line>
      </svg>
    </button>

    <div class="duration-col text-muted">
      {track.durationFormatted || '00:00'}
    </div>
  </div>
</div>

<style>
  .track-row {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px 12px;
    border-radius: 8px;
    cursor: pointer;
    transition: background-color 0.15s ease;
    user-select: none;
  }

  .track-row:hover {
    background: var(--bg-hover);
  }

  .track-row.current {
    background: var(--bg-card);
  }

  .index-col {
    width: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-muted);
    font-size: 13px;
    position: relative;
  }

  .hover-play-btn {
    display: none;
    color: var(--text-primary);
  }

  .track-row:hover .index-num {
    display: none;
  }

  .track-row:hover .hover-play-btn {
    display: inline-flex;
  }

  .thumb-col {
    width: 44px;
    height: 44px;
    border-radius: 6px;
    overflow: hidden;
    background: var(--bg-card);
    flex-shrink: 0;
  }

  .thumb-col img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .thumb-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-muted);
    background: var(--bg-surface);
  }

  .meta-col {
    flex: 1;
    min-width: 0;
  }

  .track-title {
    font-size: 14px;
    font-weight: 500;
    color: var(--text-primary);
  }

  .active-title {
    color: var(--accent-red);
    font-weight: 600;
  }

  .track-artist {
    font-size: 12px;
  }

  .actions-col {
    display: flex;
    align-items: center;
    gap: 10px;
    flex-shrink: 0;
  }

  .action-btn {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    color: var(--text-muted);
    opacity: 0.7;
    transition: all 0.15s ease;
  }

  .action-btn:hover {
    color: var(--text-primary);
    opacity: 1;
    background: rgba(255, 255, 255, 0.08);
  }

  .fav-btn.is-fav {
    color: var(--accent-red);
    opacity: 1;
  }

  .duration-col {
    font-size: 12px;
    width: 45px;
    text-align: right;
  }

  /* Equalizer Animation */
  .equalizer {
    display: flex;
    align-items: flex-end;
    gap: 2px;
    height: 14px;
  }

  .equalizer span {
    width: 3px;
    background: var(--accent-red);
    border-radius: 1px;
    animation: bounce 0.8s ease-in-out infinite alternate;
  }

  .equalizer span:nth-child(1) { height: 60%; animation-delay: 0.1s; }
  .equalizer span:nth-child(2) { height: 100%; animation-delay: 0.3s; }
  .equalizer span:nth-child(3) { height: 40%; animation-delay: 0.2s; }

  @keyframes bounce {
    0% { height: 20%; }
    100% { height: 100%; }
  }
</style>
