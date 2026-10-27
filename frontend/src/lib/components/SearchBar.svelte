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

<section class="sticky top-0 z-30 px-8 pt-5 pb-4 glass-panel border-b border-white/[0.06] flex flex-col gap-3.5 select-none">
  <!-- Search Input Bar -->
  <div class="flex items-center justify-between gap-4">
    <div class="relative flex-1 max-w-2xl">
      <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-zinc-400">
        <svg class="w-4 h-4 text-zinc-400" fill="none" stroke="currentColor" stroke-width="2.2" viewBox="0 0 24 24">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
        </svg>
      </div>
      <input
        class="w-full pl-10 pr-24 py-2.5 bg-white/[0.06] hover:bg-white/[0.08] focus:bg-white/[0.1] border border-white/10 focus:border-crimson-500 rounded-2xl text-sm font-medium text-white placeholder-zinc-400 outline-none transition duration-150 backdrop-blur-md shadow-inner"
        placeholder="Cari lagu, musisi, album, atau lirik..."
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
  </div>

  <!-- Filter Tags / Chips -->
  <div class="flex items-center gap-2 overflow-x-auto pb-0.5 text-xs font-medium">
    <button
      class="px-4 py-1.5 rounded-full transition cursor-pointer {filter === 'music_songs' ? 'bg-crimson-500 text-white font-semibold shadow-md shadow-crimson-500/25' : 'bg-white/[0.06] hover:bg-white/[0.1] text-zinc-300 border border-white/[0.08]'}"
      on:click={() => setFilter('music_songs')}
    >
      🎵 Lagu
    </button>

    <button
      class="px-4 py-1.5 rounded-full transition cursor-pointer flex items-center gap-1.5 {filter === 'videos' ? 'bg-crimson-500 text-white font-semibold shadow-md shadow-crimson-500/25' : 'bg-white/[0.06] hover:bg-white/[0.1] text-zinc-300 border border-white/[0.08]'}"
      on:click={() => setFilter('videos')}
    >
      <svg class="w-3.5 h-3.5 {filter === 'videos' ? 'text-white' : 'text-zinc-400'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <rect height="20" rx="2.18" ry="2.18" width="20" x="2" y="2"></rect>
        <line x1="7" x2="7" y1="2" y2="22"></line>
        <line x1="17" x2="17" y1="2" y2="22"></line>
        <line x1="2" x2="22" y1="12" y2="12"></line>
        <line x1="2" x2="7" y1="7" y2="7"></line>
        <line x1="2" x2="7" y1="17" y2="17"></line>
        <line x1="17" x2="22" y1="17" y2="17"></line>
        <line x1="17" x2="22" y1="7" y2="7"></line>
      </svg>
      🎬 Video
    </button>

    <button
      class="px-4 py-1.5 rounded-full transition cursor-pointer {filter === 'all' ? 'bg-crimson-500 text-white font-semibold shadow-md shadow-crimson-500/25' : 'bg-white/[0.06] hover:bg-white/[0.1] text-zinc-300 border border-white/[0.08]'}"
      on:click={() => setFilter('all')}
    >
      🌐 Semua Hasil
    </button>
  </div>
</section>
