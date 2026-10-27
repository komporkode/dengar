<script lang="ts">
  import type { SearchResultItem } from '../types'
  import TrackItem from './TrackItem.svelte'
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'

  export let title = 'Daftar Lagu'
  export let tracks: SearchResultItem[] = []
  export let isLoading = false
  export let emptyMessage = 'Tidak ada lagu yang ditemukan.'

  function handlePlayAll() {
    if (tracks.length === 0) return
    queueStore.setQueue(tracks, 0)
    playerStore.playTrack(tracks[0])
  }
</script>

<div class="track-list-wrapper">
  <div class="header-row">
    <div class="title-meta">
      <h2 class="section-title">{title}</h2>
      {#if !isLoading && tracks.length > 0}
        <span class="track-count">{tracks.length} item</span>
      {/if}
    </div>

    {#if !isLoading && tracks.length > 0}
      <button class="play-all-btn" on:click={handlePlayAll}>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
          <polygon points="5 3 19 12 5 21 5 3"></polygon>
        </svg>
        <span>Putar Semua</span>
      </button>
    {/if}
  </div>

  {#if isLoading}
    <div class="skeleton-list">
      {#each Array(6) as _}
        <div class="skeleton-row">
          <div class="skeleton-num"></div>
          <div class="skeleton-thumb"></div>
          <div class="skeleton-text">
            <div class="skeleton-line long"></div>
            <div class="skeleton-line short"></div>
          </div>
        </div>
      {/each}
    </div>
  {:else if tracks.length === 0}
    <div class="empty-state">
      <div class="empty-icon">
        <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <circle cx="12" cy="12" r="10"></circle>
          <path d="M8 12h8"></path>
        </svg>
      </div>
      <p>{emptyMessage}</p>
    </div>
  {:else}
    <div class="tracks-container">
      {#each tracks as track, i (track.id + '-' + i)}
        <TrackItem {track} index={i} />
      {/each}
    </div>
  {/if}
</div>

<style>
  .track-list-wrapper {
    display: flex;
    flex-direction: column;
    gap: 12px;
    width: 100%;
  }

  .header-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--border-subtle);
  }

  .title-meta {
    display: flex;
    align-items: baseline;
    gap: 10px;
  }

  .section-title {
    font-size: 20px;
    font-weight: 700;
    color: var(--text-primary);
  }

  .track-count {
    font-size: 12px;
    color: var(--text-muted);
  }

  .play-all-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    border-radius: 20px;
    background: var(--text-primary);
    color: var(--bg-primary);
    font-weight: 600;
    font-size: 13px;
    transition: transform 0.15s ease, background 0.15s ease;
  }

  .play-all-btn:hover {
    transform: scale(1.03);
    background: #ffffff;
  }

  .tracks-container {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  /* Skeleton Loading */
  .skeleton-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .skeleton-row {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px 12px;
  }

  .skeleton-num {
    width: 20px;
    height: 14px;
    border-radius: 4px;
    background: rgba(255, 255, 255, 0.05);
  }

  .skeleton-thumb {
    width: 44px;
    height: 44px;
    border-radius: 6px;
    background: rgba(255, 255, 255, 0.05);
  }

  .skeleton-text {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .skeleton-line {
    height: 12px;
    border-radius: 4px;
    background: rgba(255, 255, 255, 0.05);
    animation: shimmer 1.2s infinite ease-in-out;
  }

  .skeleton-line.long { width: 45%; }
  .skeleton-line.short { width: 25%; }

  @keyframes shimmer {
    0%, 100% { opacity: 0.4; }
    50% { opacity: 0.9; }
  }

  /* Empty State */
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
    color: var(--text-muted);
    text-align: center;
    gap: 12px;
  }

  .empty-icon {
    opacity: 0.6;
  }
</style>
