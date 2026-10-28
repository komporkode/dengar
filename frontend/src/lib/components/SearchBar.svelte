<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import type { SearchFilter } from '../types'
  import { queueStore } from '../stores/queueStore'
  import { enterMiniPlayer } from '../stores/settingsStore'
  import { t } from '../stores/i18nStore'

  export let query = ''
  export let filter: SearchFilter = 'music_songs'
  export let isLoading = false

  const dispatch = createEventDispatcher<{
    search: { query: string; filter: SearchFilter }
  }>()

  let debounceTimer: any
  $: queueCount = $queueStore.items.length

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

<section class="sticky top-0 z-30 px-8 pt-5 pb-4 glass-panel border-b border-white/[0.06] flex flex-col gap-3.5 select-none">
  <!-- Search Input Bar & Top Right Actions -->
  <div class="flex items-center justify-between gap-4">
    <!-- Search Bar Input -->
    <div class="relative flex-1 max-w-2xl">
      <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-zinc-400">
        <svg class="w-4 h-4 text-zinc-400" fill="none" stroke="currentColor" stroke-width="2.2" viewBox="0 0 24 24">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
        </svg>
      </div>
      <input
        class="w-full pl-10 pr-24 py-2.5 bg-white/[0.06] hover:bg-white/[0.08] focus:bg-white/[0.1] border border-white/10 focus:border-crimson-500 rounded-2xl text-sm font-medium text-white placeholder-zinc-400 outline-none transition duration-150 backdrop-blur-md shadow-inner"
        placeholder={$t('search_placeholder')}
        type="text"
        bind:value={query}
        on:input={handleInput}
        on:keydown={handleKeydown}
      />
      <div class="absolute inset-y-0 right-0 pr-3 flex items-center gap-1.5">
        {#if isLoading}
          <div class="w-4 h-4 border-2 border-white/20 border-t-crimson-500 rounded-full animate-spin"></div>
        {:else if query.length > 0}
          <button
            class="p-1 text-zinc-400 hover:text-white transition cursor-pointer"
            on:click={clearQuery}
            title="Hapus teks"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path d="M18 6L6 18M6 6l12 12"></path>
            </svg>
          </button>
        {/if}
        <kbd class="hidden sm:inline-block px-1.5 py-0.5 text-[10px] font-mono font-medium text-zinc-400 bg-white/[0.08] border border-white/10 rounded-md">Enter</kbd>
      </div>
    </div>

    <!-- Top Right Action Controls: Queue, Mini Player, Settings -->
    <div class="flex items-center gap-2">
      <!-- Queue Button with Badge -->
      <button
        class="flex items-center gap-2 px-3 py-2 rounded-xl bg-white/[0.06] hover:bg-white/[0.1] border border-white/10 text-zinc-200 hover:text-white transition cursor-pointer relative {$queueStore.isDrawerOpen ? 'bg-white/[0.12] text-white border-crimson-500/50' : ''}"
        on:click={() => queueStore.toggleDrawer()}
        title={$t('nav_queue')}
      >
        <svg class="w-4 h-4 text-zinc-300" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <line x1="8" y1="6" x2="21" y2="6"></line>
          <line x1="8" y1="12" x2="21" y2="12"></line>
          <line x1="8" y1="18" x2="21" y2="18"></line>
          <line x1="3" y1="6" x2="3.01" y2="6"></line>
          <line x1="3" y1="12" x2="3.01" y2="12"></line>
          <line x1="3" y1="18" x2="3.01" y2="18"></line>
        </svg>
        <span class="text-xs font-semibold hidden md:inline">{$t('nav_queue')}</span>
        {#if queueCount > 0}
          <span class="px-1.5 py-0.2 rounded-full bg-crimson-500 text-[10px] font-bold text-white font-mono shadow-sm">
            {queueCount}
          </span>
        {/if}
      </button>

      <!-- Mini Player Button -->
      <button
        class="p-2 rounded-xl bg-white/[0.06] hover:bg-white/[0.1] border border-white/10 text-zinc-300 hover:text-white transition cursor-pointer"
        on:click={enterMiniPlayer}
        title={$t('mini_player')}
      >
        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
          <rect x="13" y="13" width="6" height="6" rx="1" fill="currentColor"></rect>
        </svg>
      </button>
    </div>
  </div>

  <!-- Filter Tags / Chips (Emoji replaced with pure SVG) -->
  <div class="flex items-center gap-2 overflow-x-auto pb-0.5 text-xs font-medium">
    <!-- Lagu / Songs -->
    <button
      class="px-4 py-1.5 rounded-full transition cursor-pointer flex items-center gap-1.5 {filter === 'music_songs' ? 'bg-crimson-500 text-white font-semibold shadow-md shadow-crimson-500/25' : 'bg-white/[0.06] hover:bg-white/[0.1] text-zinc-300 border border-white/[0.08]'}"
      on:click={() => setFilter('music_songs')}
    >
      <svg class="w-3.5 h-3.5 {filter === 'music_songs' ? 'text-white' : 'text-zinc-400'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path d="M9 18V5l12-2v13"></path>
        <circle cx="6" cy="18" r="3" fill="currentColor"></circle>
        <circle cx="18" cy="16" r="3" fill="currentColor"></circle>
      </svg>
      <span>{$t('filter_songs')}</span>
    </button>

    <!-- Video / Videos -->
    <button
      class="px-4 py-1.5 rounded-full transition cursor-pointer flex items-center gap-1.5 {filter === 'videos' ? 'bg-crimson-500 text-white font-semibold shadow-md shadow-crimson-500/25' : 'bg-white/[0.06] hover:bg-white/[0.1] text-zinc-300 border border-white/[0.08]'}"
      on:click={() => setFilter('videos')}
    >
      <svg class="w-3.5 h-3.5 {filter === 'videos' ? 'text-white' : 'text-zinc-400'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <polygon points="23 7 16 12 23 17 23 7"></polygon>
        <rect x="1" y="5" width="15" height="14" rx="2" ry="2"></rect>
      </svg>
      <span>{$t('filter_videos')}</span>
    </button>

    <!-- Semua / All -->
    <button
      class="px-4 py-1.5 rounded-full transition cursor-pointer flex items-center gap-1.5 {filter === 'all' ? 'bg-crimson-500 text-white font-semibold shadow-md shadow-crimson-500/25' : 'bg-white/[0.06] hover:bg-white/[0.1] text-zinc-300 border border-white/[0.08]'}"
      on:click={() => setFilter('all')}
    >
      <svg class="w-3.5 h-3.5 {filter === 'all' ? 'text-white' : 'text-zinc-400'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="2" y1="12" x2="22" y2="12"></line>
        <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
      </svg>
      <span>{$t('filter_all')}</span>
    </button>
  </div>
</section>
