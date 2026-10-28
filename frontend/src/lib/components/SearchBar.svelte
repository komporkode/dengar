<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import type { SearchFilter } from '../types'
  import { queueStore } from '../stores/queueStore'
  import { enterMiniPlayer } from '../stores/settingsStore'
  import { t } from '../stores/i18nStore'

  export let query = ''
  export let isLoading = false

  const dispatch = createEventDispatcher<{
    search: { query: string; filter: SearchFilter }
  }>()

  let queueCount = 0
  $: queueCount = $queueStore.items.length

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter') {
      triggerSearch()
    }
  }

  function clearQuery() {
    query = ''
  }

  function triggerSearch() {
    dispatch('search', { query: query.trim(), filter: 'music_songs' })
  }
</script>

<section class="sticky top-0 z-30 px-4 sm:px-8 py-3.5 bg-[#0E1017]/85 backdrop-blur-xl border-b border-white/[0.04] select-none">
  <!-- Search Input Bar & Top Right Actions -->
  <div class="flex items-center justify-between gap-2.5 sm:gap-4">
    <!-- Search Bar Input with Clickable Magnifying Glass -->
    <div class="relative flex-1 max-w-2xl">
      <button
        type="button"
        class="absolute inset-y-0 left-0 pl-3.5 pr-2 flex items-center text-zinc-400 hover:text-crimson-400 active:scale-90 transition cursor-pointer"
        on:click={triggerSearch}
        title="Cari (Enter)"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2.2" viewBox="0 0 24 24">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
        </svg>
      </button>

      <input
        class="w-full pl-10 pr-24 py-2.5 bg-white/[0.06] hover:bg-white/[0.08] focus:bg-white/[0.1] border border-white/10 focus:border-crimson-500 rounded-2xl text-sm font-medium text-white placeholder-zinc-400 outline-none transition duration-150 backdrop-blur-md shadow-inner"
        placeholder={$t('search_placeholder')}
        type="text"
        bind:value={query}
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

    <!-- Top Right Action Controls: Unified Sleek Segmented Capsule Dock -->
    <div class="flex items-center bg-white/[0.04] p-1 rounded-2xl border border-white/[0.08] backdrop-blur-md shadow-xs shrink-0">
      <!-- Queue Button with Badge -->
      <button
        class="flex items-center gap-2 px-3 py-1.5 rounded-xl transition-all cursor-pointer {$queueStore.isDrawerOpen ? 'bg-white/[0.14] text-white shadow-xs' : 'text-zinc-300 hover:text-white hover:bg-white/[0.06]'}"
        on:click={() => queueStore.toggleDrawer()}
        title={$t('nav_queue')}
      >
        <svg class="w-4 h-4 {$queueStore.isDrawerOpen ? 'text-crimson-400' : 'text-zinc-300'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <line x1="8" y1="6" x2="21" y2="6"></line>
          <line x1="8" y1="12" x2="21" y2="12"></line>
          <line x1="8" y1="18" x2="21" y2="18"></line>
          <line x1="3" y1="6" x2="3.01" y2="6"></line>
          <line x1="3" y1="12" x2="3.01" y2="12"></line>
          <line x1="3" y1="18" x2="3.01" y2="18"></line>
        </svg>
        <span class="text-xs font-semibold hidden md:inline">{$t('nav_queue')}</span>
        {#if queueCount > 0}
          <span class="px-1.5 py-0.2 rounded-full bg-crimson-500 text-[10px] font-bold text-white font-mono shadow-xs">
            {queueCount}
          </span>
        {/if}
      </button>

      <div class="w-px h-3.5 bg-white/10 my-auto mx-0.5"></div>

      <!-- Mini Player Button -->
      <button
        class="p-1.5 rounded-xl text-zinc-300 hover:text-white hover:bg-white/[0.06] transition-all cursor-pointer"
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
</section>
