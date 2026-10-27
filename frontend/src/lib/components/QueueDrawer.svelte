<script lang="ts">
  import { queueStore } from '../stores/queueStore'
  import { playerStore } from '../stores/playerStore'
  import { t } from '../stores/i18nStore'

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
  <div class="fixed inset-0 bg-black/50 z-40 backdrop-blur-xs" on:click={() => queueStore.toggleDrawer(false)} role="presentation"></div>

  <aside class="fixed top-0 right-0 bottom-24 w-84 sm:w-96 bg-surface/95 backdrop-blur-xl border-l border-white/[0.08] shadow-2xl flex flex-col z-50 animate-slide-in select-none">
    <div class="flex items-center justify-between px-5 py-4 border-b border-white/[0.08]">
      <div class="flex items-center gap-2.5">
        <h3 class="font-bold text-sm tracking-wide text-white">{$t('queue_title')}</h3>
        <span class="text-[10px] font-mono px-2 py-0.5 rounded-full bg-white/[0.08] text-zinc-300 border border-white/10">
          {items.length} {$t('items_count')}
        </span>
      </div>

      <div class="flex items-center gap-2">
        {#if items.length > 0}
          <button class="text-xs text-crimson-400 hover:text-crimson-300 px-2 py-1 rounded hover:bg-crimson-500/10 transition cursor-pointer" on:click={queueStore.clearQueue}>
            {$t('clear_all')}
          </button>
        {/if}
        <button class="p-1 rounded-lg text-zinc-400 hover:text-white hover:bg-white/[0.08] transition cursor-pointer" on:click={() => queueStore.toggleDrawer(false)} title="Tutup Antrean">
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto p-3 space-y-1">
      {#if items.length === 0}
        <div class="h-full flex flex-col items-center justify-center text-center p-6 text-zinc-500 gap-2">
          <div class="w-10 h-10 rounded-xl bg-white/[0.03] border border-white/10 flex items-center justify-center">
            <svg class="w-5 h-5 text-zinc-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
              <line x1="8" y1="6" x2="21" y2="6"></line>
              <line x1="8" y1="12" x2="21" y2="12"></line>
              <line x1="8" y1="18" x2="21" y2="18"></line>
            </svg>
          </div>
          <p class="text-xs font-medium text-zinc-400">{$t('queue_empty')}</p>
          <span class="text-[11px] text-zinc-500">{$t('queue_empty_tip')}</span>
        </div>
      {:else}
        {#each items as item, i (item.id + '-' + i)}
          <div
            class="group flex items-center gap-3 px-3 py-2 rounded-xl cursor-pointer transition {i === currentIndex ? 'bg-crimson-500/15 border border-crimson-500/30' : 'hover:bg-white/[0.04]'}"
            on:click={() => handlePlayQueueItem(i)}
            role="button"
            tabindex="0"
            on:keydown={(e) => e.key === 'Enter' && handlePlayQueueItem(i)}
          >
            <div class="w-5 flex items-center justify-center shrink-0">
              {#if i === currentIndex}
                <span class="w-2 h-2 rounded-full bg-crimson-400 shadow-[0_0_8px_#FF4D62]"></span>
              {:else}
                <span class="text-xs font-mono text-zinc-500">{i + 1}</span>
              {/if}
            </div>

            <div class="w-9 h-9 rounded-lg overflow-hidden bg-surface-highlight shrink-0 flex items-center justify-center border border-white/10">
              {#if item.thumbnailUrl}
                <img src={item.thumbnailUrl} alt={item.title} class="w-full h-full object-cover" />
              {:else}
                <svg class="w-4 h-4 text-zinc-500" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
                </svg>
              {/if}
            </div>

            <div class="flex-1 min-w-0">
              <span class="text-xs font-medium truncate block {i === currentIndex ? 'text-crimson-300 font-bold' : 'text-zinc-200'}">
                {item.title}
              </span>
              <span class="text-[11px] text-zinc-400 truncate block">{item.uploaderName}</span>
            </div>

            <button
              class="p-1 text-zinc-500 hover:text-crimson-400 hover:bg-crimson-500/10 rounded opacity-0 group-hover:opacity-100 transition cursor-pointer"
              on:click={(e) => handleRemoveItem(e, i)}
              title="Hapus dari antrean"
            >
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>
        {/each}
      {/if}
    </div>
  </aside>
{/if}
