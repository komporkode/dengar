<script lang="ts">
  import { queueStore } from '../stores/queueStore'
  import { playerStore } from '../stores/playerStore'
  import { t } from '../stores/i18nStore'

  $: items = $queueStore.items
  $: currentIndex = $queueStore.currentIndex
  $: isOpen = $queueStore.isDrawerOpen
  $: isPlaying = $playerStore.isPlaying

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
  <!-- Subtle click-outside backdrop (not pitch black, lets user see the music app) -->
  <div
    class="fixed inset-0 bg-black/30 z-50 backdrop-blur-xs transition-opacity duration-300"
    on:click={() => queueStore.toggleDrawer(false)}
    role="presentation"
  ></div>

  <!-- Ultra-Premium Slide-over Queue Drawer (Full height top-9 to bottom-0, on top of player) -->
  <aside class="fixed top-9 right-0 bottom-0 w-88 sm:w-96 bg-[#10121A]/95 backdrop-blur-2xl border-l border-white/[0.08] shadow-2xl flex flex-col z-50 animate-spring-slide select-none overflow-hidden">
    <!-- Header -->
    <div class="flex items-center justify-between px-5 py-4 border-b border-white/[0.08] bg-[#0A0C12]/90 shrink-0">
      <div class="flex items-center gap-2.5">
        <div class="w-8 h-8 rounded-xl bg-crimson-500/20 border border-crimson-500/30 flex items-center justify-center text-crimson-400 shadow-xs">
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
            <line x1="8" y1="6" x2="21" y2="6"></line>
            <line x1="8" y1="12" x2="21" y2="12"></line>
            <line x1="8" y1="18" x2="21" y2="18"></line>
            <line x1="3" y1="6" x2="3.01" y2="6"></line>
            <line x1="3" y1="12" x2="3.01" y2="12"></line>
            <line x1="3" y1="18" x2="3.01" y2="18"></line>
          </svg>
        </div>
        <div>
          <h3 class="font-bold text-sm tracking-wide text-white">{$t('queue_title')}</h3>
          <p class="text-[11px] text-zinc-400 font-mono">{items.length} {$t('items_count')}</p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        {#if items.length > 0}
          <button
            class="text-xs font-semibold text-crimson-400 hover:text-crimson-300 px-2.5 py-1.5 rounded-lg hover:bg-crimson-500/10 transition cursor-pointer"
            on:click={queueStore.clearQueue}
          >
            {$t('clear_all')}
          </button>
        {/if}
        <button
          class="p-2 rounded-xl text-zinc-400 hover:text-white hover:bg-white/[0.08] transition cursor-pointer"
          on:click={() => queueStore.toggleDrawer(false)}
          title="Tutup Antrean"
        >
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </div>

    <!-- Items List with comfortable bottom scrolling space -->
    <div class="flex-1 overflow-y-auto p-3 space-y-1.5 pb-24">
      {#if items.length === 0}
        <div class="h-full min-h-[300px] flex flex-col items-center justify-center text-center p-6 text-zinc-500 gap-3">
          <div class="w-14 h-14 rounded-2xl bg-white/[0.03] border border-white/10 flex items-center justify-center shadow-inner">
            <svg class="w-7 h-7 text-zinc-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
              <line x1="8" y1="6" x2="21" y2="6"></line>
              <line x1="8" y1="12" x2="21" y2="12"></line>
              <line x1="8" y1="18" x2="21" y2="18"></line>
            </svg>
          </div>
          <div>
            <p class="text-sm font-semibold text-zinc-300">{$t('queue_empty')}</p>
            <span class="text-xs text-zinc-500 mt-1 block">{$t('queue_empty_tip')}</span>
          </div>
        </div>
      {:else}
        {#each items as item, i (item.id + '-' + i)}
          <div
            class="group flex items-center gap-3 px-3 py-2.5 rounded-2xl cursor-pointer transition-all duration-150 {i === currentIndex ? 'bg-crimson-500/15 border border-crimson-500/35 shadow-sm shadow-crimson-500/10' : 'hover:bg-white/[0.05] border border-transparent'}"
            on:click={() => handlePlayQueueItem(i)}
            role="button"
            tabindex="0"
            on:keydown={(e) => e.key === 'Enter' && handlePlayQueueItem(i)}
          >
            <!-- Index / Animated Equalizer -->
            <div class="w-6 flex items-center justify-center shrink-0">
              {#if i === currentIndex && isPlaying}
                <div class="flex items-center gap-0.5">
                  <span class="inline-block w-0.5 h-3 bg-crimson-400 rounded-xs animate-pulse"></span>
                  <span class="inline-block w-0.5 h-4 bg-crimson-500 rounded-xs animate-bounce"></span>
                  <span class="inline-block w-0.5 h-2 bg-crimson-400 rounded-xs animate-pulse"></span>
                </div>
              {:else if i === currentIndex}
                <span class="w-2.5 h-2.5 rounded-full bg-crimson-400 shadow-[0_0_8px_#FF4D62]"></span>
              {:else}
                <span class="text-xs font-mono text-zinc-500 group-hover:text-zinc-300 transition-colors">{i + 1}</span>
              {/if}
            </div>

            <!-- Thumbnail -->
            <div class="w-11 h-11 rounded-xl overflow-hidden bg-surface-highlight shrink-0 flex items-center justify-center border border-white/10 shadow-sm group-hover:scale-105 transition-transform duration-200">
              {#if item.thumbnailUrl}
                <img src={item.thumbnailUrl} alt={item.title} class="w-full h-full object-cover" />
              {:else}
                <svg class="w-5 h-5 text-zinc-500" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
                </svg>
              {/if}
            </div>

            <!-- Track Info -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-1.5">
                <span class="text-xs font-semibold truncate block {i === currentIndex ? 'text-crimson-400 font-bold' : 'text-zinc-200 group-hover:text-white transition-colors'}">
                  {item.title}
                </span>
                {#if i === currentIndex}
                  <span class="text-[8px] font-bold uppercase tracking-wider px-1.5 py-0.2 rounded bg-crimson-500/20 text-crimson-300 border border-crimson-500/30 shrink-0">
                    {$t('playing')}
                  </span>
                {/if}
              </div>
              <span class="text-[11px] text-zinc-400 truncate block mt-0.5">{item.uploaderName}</span>
            </div>

            <!-- Duration & Remove button -->
            <div class="flex items-center gap-1 shrink-0">
              <span class="text-[10px] font-mono text-zinc-500 group-hover:hidden">
                {item.durationFormatted || ''}
              </span>
              <button
                class="p-1.5 text-zinc-500 hover:text-crimson-400 hover:bg-crimson-500/10 rounded-lg hidden group-hover:flex transition cursor-pointer"
                on:click={(e) => handleRemoveItem(e, i)}
                title="Hapus dari antrean"
              >
                <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>
          </div>
        {/each}
      {/if}
    </div>
  </aside>
{/if}
