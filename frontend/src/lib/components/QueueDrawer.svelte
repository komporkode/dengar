<script lang="ts">
  import { queueStore } from '../stores/queueStore'
  import { playerStore } from '../stores/playerStore'

  $: items = $queueStore.items
  $: currentIndex = $queueStore.currentIndex
  $: isOpen = $queueStore.isDrawerOpen

  function handlePlayQueueItem(index: number) {
    const track = items[index]
    if (track) {
      queueStore.setQueue(items, index)
      playerStore.playTrack(track)
    }
  }

  function handleRemoveItem(e: MouseEvent, index: number) {
    e.stopPropagation()
    queueStore.removeFromQueue(index)
  }
</script>

{#if isOpen}
  <!-- Backdrop for click outside -->
  <div class="backdrop" on:click={() => queueStore.toggleDrawer(false)} role="presentation"></div>

  <aside class="queue-drawer">
    <div class="drawer-header">
      <div class="header-title">
        <h3>Antrean Lagu</h3>
        <span class="count-badge">{items.length} lagu</span>
      </div>

      <div class="header-actions">
        {#if items.length > 0}
          <button class="clear-btn" on:click={queueStore.clearQueue}>Hapus Semua</button>
        {/if}
        <button class="close-btn" on:click={() => queueStore.toggleDrawer(false)} title="Tutup Antrean" aria-label="Tutup Antrean">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </div>

    <div class="drawer-body">
      {#if items.length === 0}
        <div class="empty-queue">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <line x1="8" y1="6" x2="21" y2="6"></line>
            <line x1="8" y1="12" x2="21" y2="12"></line>
            <line x1="8" y1="18" x2="21" y2="18"></line>
          </svg>
          <p>Antrean kosong</p>
          <span class="text-muted">Tambahkan lagu dari pencarian atau favorit</span>
        </div>
      {:else}
        <div class="queue-items-list">
          {#each items as item, i (item.id + '-' + i)}
            <div
              class="queue-item"
              class:current-playing={i === currentIndex}
              on:click={() => handlePlayQueueItem(i)}
              role="button"
              tabindex="0"
              on:keydown={(e) => e.key === 'Enter' && handlePlayQueueItem(i)}
            >
              <div class="item-status">
                {#if i === currentIndex}
                  <div class="playing-pulse"></div>
                {:else}
                  <span class="item-idx">{i + 1}</span>
                {/if}
              </div>

              <div class="item-thumb">
                {#if item.thumbnailUrl}
                  <img src={item.thumbnailUrl} alt={item.title} />
                {:else}
                  <div class="no-thumb">🎵</div>
                {/if}
              </div>

              <div class="item-meta">
                <span class="item-title text-truncate" class:active={i === currentIndex}>{item.title}</span>
                <span class="item-artist text-truncate text-secondary">{item.uploaderName}</span>
              </div>

              <button
                class="remove-btn"
                on:click={(e) => handleRemoveItem(e, i)}
                title="Hapus dari antrean"
              >
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </aside>
{/if}

<style>
  .backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: var(--player-height);
    background: rgba(0, 0, 0, 0.4);
    z-index: 105;
  }

  .queue-drawer {
    position: fixed;
    top: 0;
    right: 0;
    bottom: var(--player-height);
    width: 360px;
    background: var(--bg-surface);
    border-left: 1px solid var(--border-subtle);
    box-shadow: -8px 0 24px rgba(0, 0, 0, 0.5);
    display: flex;
    flex-direction: column;
    z-index: 110;
    animation: slideIn 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  }

  @keyframes slideIn {
    from { transform: translateX(100%); }
    to { transform: translateX(0); }
  }

  .drawer-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 18px 20px;
    border-bottom: 1px solid var(--border-subtle);
  }

  .header-title {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .header-title h3 {
    font-size: 16px;
    font-weight: 700;
  }

  .count-badge {
    font-size: 11px;
    padding: 2px 8px;
    border-radius: 10px;
    background: rgba(255, 255, 255, 0.08);
    color: var(--text-muted);
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .clear-btn {
    font-size: 12px;
    color: var(--accent-red);
    padding: 4px 8px;
    border-radius: 6px;
  }

  .clear-btn:hover {
    background: rgba(255, 0, 0, 0.1);
  }

  .close-btn {
    width: 30px;
    height: 30px;
    border-radius: 50%;
    color: var(--text-muted);
  }

  .close-btn:hover {
    color: var(--text-primary);
    background: rgba(255, 255, 255, 0.1);
  }

  .drawer-body {
    flex: 1;
    overflow-y: auto;
    padding: 12px 14px;
  }

  .empty-queue {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: var(--text-muted);
    gap: 8px;
    text-align: center;
  }

  .queue-items-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .queue-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 10px;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.15s ease;
  }

  .queue-item:hover {
    background: var(--bg-hover);
  }

  .queue-item.current-playing {
    background: var(--bg-card);
  }

  .item-status {
    width: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .item-idx {
    font-size: 12px;
    color: var(--text-muted);
  }

  .playing-pulse {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--accent-red);
    box-shadow: 0 0 8px var(--accent-red);
  }

  .item-thumb {
    width: 38px;
    height: 38px;
    border-radius: 6px;
    overflow: hidden;
    background: var(--bg-card);
    flex-shrink: 0;
  }

  .item-thumb img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .no-thumb {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
  }

  .item-meta {
    flex: 1;
    min-width: 0;
  }

  .item-title {
    font-size: 13px;
    font-weight: 500;
  }

  .item-title.active {
    color: var(--accent-red);
    font-weight: 600;
  }

  .item-artist {
    font-size: 11px;
  }

  .remove-btn {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    color: var(--text-muted);
    opacity: 0;
    transition: all 0.15s ease;
  }

  .queue-item:hover .remove-btn {
    opacity: 1;
  }

  .remove-btn:hover {
    color: var(--accent-red);
    background: rgba(255, 0, 0, 0.1);
  }
</style>
