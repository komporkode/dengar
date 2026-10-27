<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import type { SearchFilter } from '../types'

  export let query = ''
  export let filter: SearchFilter = 'music_songs'
  export let isLoading = false

  const dispatch = createEventDispatcher<{
    search: { query: string; filter: SearchFilter }
  }>()

  let debounceTimer: any

  function handleInput() {
    clearTimeout(debounceTimer)
    debounceTimer = setTimeout(() => {
      triggerSearch()
    }, 450)
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter') {
      clearTimeout(debounceTimer)
      triggerSearch()
    }
  }

  function setFilter(newFilter: SearchFilter) {
    if (filter === newFilter) return
    filter = newFilter
    triggerSearch()
  }

  function clearQuery() {
    query = ''
    triggerSearch()
  }

  function triggerSearch() {
    dispatch('search', { query: query.trim(), filter })
  }
</script>

<div class="search-container">
  <div class="search-input-wrapper">
    <svg class="search-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <circle cx="11" cy="11" r="8"></circle>
      <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
    </svg>

    <input
      type="search"
      placeholder="Cari lagu, artis, atau video musik YouTube..."
      bind:value={query}
      on:input={handleInput}
      on:keydown={handleKeydown}
    />

    {#if isLoading}
      <div class="spinner"></div>
    {:else if query.length > 0}
      <button class="clear-btn" on:click={clearQuery} aria-label="Hapus pencarian">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    {/if}
  </div>

  <div class="filter-chips">
    <button
      class="chip"
      class:active={filter === 'music_songs'}
      on:click={() => setFilter('music_songs')}
    >
      🎵 Lagu
    </button>
    <button
      class="chip"
      class:active={filter === 'videos'}
      on:click={() => setFilter('videos')}
    >
      🎬 Video
    </button>
    <button
      class="chip"
      class:active={filter === 'all'}
      on:click={() => setFilter('all')}
    >
      🌐 Semua
    </button>
  </div>
</div>

<style>
  .search-container {
    display: flex;
    flex-direction: column;
    gap: 14px;
    width: 100%;
    max-width: 680px;
    margin-bottom: 24px;
  }

  .search-input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
    width: 100%;
  }

  .search-icon {
    position: absolute;
    left: 18px;
    color: var(--text-muted);
    pointer-events: none;
  }

  input[type="search"] {
    width: 100%;
    padding-left: 48px;
    padding-right: 44px;
    height: 44px;
    border-radius: 24px;
    background: var(--bg-card);
    border: 1px solid var(--border-subtle);
    font-size: 14px;
    color: var(--text-primary);
  }

  input[type="search"]:focus {
    border-color: var(--border-focus);
    background: var(--bg-surface);
    box-shadow: 0 0 16px rgba(0, 0, 0, 0.4);
  }

  .clear-btn {
    position: absolute;
    right: 14px;
    width: 24px;
    height: 24px;
    border-radius: 50%;
    color: var(--text-muted);
  }

  .clear-btn:hover {
    color: var(--text-primary);
    background: rgba(255, 255, 255, 0.1);
  }

  .spinner {
    position: absolute;
    right: 16px;
    width: 16px;
    height: 16px;
    border: 2px solid rgba(255, 255, 255, 0.2);
    border-top-color: var(--accent-red);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .filter-chips {
    display: flex;
    gap: 8px;
  }

  .chip {
    padding: 6px 14px;
    border-radius: 16px;
    font-size: 13px;
    font-weight: 500;
    background: var(--bg-card);
    border: 1px solid var(--border-subtle);
    color: var(--text-secondary);
    transition: all 0.15s ease;
  }

  .chip:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  .chip.active {
    background: var(--text-primary);
    color: var(--bg-primary);
    font-weight: 600;
    border-color: var(--text-primary);
  }
</style>
