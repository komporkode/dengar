<script lang="ts">
  import { fly, fade } from "svelte/transition";
  import { cubicOut } from "svelte/easing";
  import {
    queueStore,
    isAutoplayEnabled,
    toggleAutoplay,
  } from "../stores/queueStore";
  import { playerStore } from "../stores/playerStore";
  import { showToast } from "../stores/settingsStore";
  import { t } from "../stores/i18nStore";
  import { handleCoverError } from "../stores/offlineStore";

  $: items = $queueStore.items;
  $: currentIndex = $queueStore.currentIndex;
  $: isOpen = $queueStore.isDrawerOpen;
  $: isPlaying = $playerStore.isPlaying;

  function handlePlayQueueItem(index: number) {
    const track = queueStore.jumpToIndex(index);
    if (track) {
      playerStore.playTrack(track, undefined, true);
    }
  }

  function handleRemoveItem(e: MouseEvent, index: number) {
    e.stopPropagation();
    const removedItem = items[index];
    queueStore.removeFromQueue(index);
    if (removedItem) {
      showToast($t("queue_toast_removed"), () => {
        queueStore.addToQueue(removedItem);
      });
    }
  }

  let draggedIndex: number | null = null;

  function handleDragStart(index: number, e: DragEvent) {
    draggedIndex = index;
    if (e.dataTransfer) {
      e.dataTransfer.effectAllowed = "move";
      e.dataTransfer.setData("text/plain", String(index));
    }
  }

  function handleDragOver(index: number, e: DragEvent) {
    e.preventDefault();
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = "move";
    }
  }

  function handleDrop(targetIndex: number, e: DragEvent) {
    e.preventDefault();
    if (draggedIndex !== null && draggedIndex !== targetIndex) {
      queueStore.reorderQueue(draggedIndex, targetIndex);
    }
    draggedIndex = null;
  }

  function handleDragEnd() {
    draggedIndex = null;
  }
</script>

<svelte:window
  on:keydown={(e) => isOpen && e.key === "Escape" && queueStore.toggleDrawer(false)}
/>

{#if isOpen}
  <!-- Subtle click-outside backdrop with smooth fade animation -->
  <div
    class="fixed inset-0 bg-black/50 z-50 backdrop-blur-xs transition-opacity"
    transition:fade={{ duration: 200 }}
    on:click={() => queueStore.toggleDrawer(false)}
    role="presentation"
  ></div>

  <!-- Ultra-Premium Slide-over Queue Drawer with comfortable width and breathing room -->
  <aside
    class="fixed top-9 right-0 bottom-0 w-full sm:w-[440px] md:w-[460px] max-w-[480px] bg-[#0E1017]/95 backdrop-blur-2xl border-l border-white/[0.08] shadow-2xl flex flex-col z-50 select-none overflow-hidden"
    transition:fly={{ x: 480, duration: 280, easing: cubicOut }}
  >
    <!-- Ultra-Minimalist Header (Single Tier, Clean, Spacious) -->
    <div
      class="h-14 px-4 border-b border-white/[0.06] bg-[#0A0C12]/95 flex items-center justify-between gap-3 shrink-0"
    >
      <!-- Left: Title + Song Count -->
      <div class="flex items-center gap-2 min-w-0">
        <h3 class="font-bold text-sm tracking-wide text-white truncate">
          {$t("queue_title")}
        </h3>
        <span class="text-xs font-mono text-zinc-500 shrink-0">
          ({items.length})
        </span>
      </div>

      <!-- Right: Minimalist Controls & Actions Toolbar -->
      <div class="flex items-center gap-1.5 shrink-0">
        <!-- Autoplay Radio Toggle -->
        <button
          type="button"
          class="flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-xs font-medium transition cursor-pointer border {$isAutoplayEnabled
            ? 'bg-crimson-500/15 text-crimson-300 border-crimson-500/30 hover:bg-crimson-500/20'
            : 'bg-white/[0.04] text-zinc-400 border-white/5 hover:text-white hover:bg-white/[0.08]'}"
          on:click={toggleAutoplay}
          title={$t("queue_autoplay_tip")}
        >
          <span
            class="w-1.5 h-1.5 rounded-full {$isAutoplayEnabled
              ? 'bg-crimson-400 shadow-[0_0_6px_#E61E38]'
              : 'bg-zinc-600'}"
          ></span>
          <span>{$t("queue_autoplay")}</span>
        </button>

        {#if items.length > 1 && currentIndex >= 0}
          <button
            type="button"
            class="text-xs font-medium text-zinc-400 hover:text-zinc-200 px-2 py-1 rounded-lg hover:bg-white/[0.06] transition cursor-pointer"
            on:click={queueStore.clearUpcomingQueue}
            title={$t("queue_clear_upcoming_tip")}
          >
            {$t("queue_clear_upcoming")}
          </button>
        {/if}

        {#if items.length > 0}
          <button
            type="button"
            class="text-xs font-medium text-crimson-400 hover:text-crimson-300 px-2 py-1 rounded-lg hover:bg-crimson-500/10 transition cursor-pointer"
            on:click={queueStore.clearQueue}
            title={$t("queue_clear_all_tip")}
          >
            {$t("clear_all")}
          </button>
        {/if}

        <span class="w-px h-3.5 bg-white/10 mx-0.5"></span>

        <!-- Close Drawer Button -->
        <button
          type="button"
          class="w-7 h-7 rounded-lg text-zinc-400 hover:text-white hover:bg-white/[0.08] flex items-center justify-center transition cursor-pointer"
          on:click={() => queueStore.toggleDrawer(false)}
          title={$t("close")}
          aria-label={$t("close")}
        >
          <svg
            class="w-4 h-4"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </div>

    <!-- Items List with comfortable bottom scrolling space -->
    <div class="flex-1 overflow-y-auto px-4 py-3 space-y-1 pb-28">
      {#if items.length === 0}
        <div
          class="h-full min-h-[300px] flex flex-col items-center justify-center text-center p-6 text-zinc-500 gap-3"
        >
          <div
            class="w-12 h-12 rounded-2xl bg-white/[0.03] border border-white/10 flex items-center justify-center shadow-inner"
          >
            <svg
              class="w-6 h-6 text-zinc-500"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.8"
            >
              <line x1="8" y1="6" x2="21" y2="6"></line>
              <line x1="8" y1="12" x2="21" y2="12"></line>
              <line x1="8" y1="18" x2="21" y2="18"></line>
            </svg>
          </div>
          <div>
            <p class="text-xs font-semibold text-zinc-300">
              {$t("queue_empty")}
            </p>
            <span class="text-[11px] text-zinc-500 mt-1 block"
              >{$t("queue_empty_tip")}</span
            >
          </div>
        </div>
      {:else}
        {#each items as item, i (item.id + "-" + i)}
          <!-- Section Divider: Sebelumnya -->
          {#if currentIndex > 0 && i === 0}
            <div
              class="text-[10px] font-bold uppercase tracking-wider text-zinc-500 px-1 pt-1 pb-1"
            >
              {$t("queue_previous")} ({currentIndex})
            </div>
          {/if}

          <!-- Section Divider: Sedang Diputar -->
          {#if currentIndex >= 0 && i === currentIndex}
            <div
              class="text-[10px] font-bold uppercase tracking-wider text-crimson-400 px-1 pt-2.5 pb-1 flex items-center gap-1.5"
            >
              <span
                class="w-1.5 h-1.5 rounded-full bg-crimson-400 animate-pulse"
              ></span>
              <span>{$t("queue_now_playing")}</span>
            </div>
          {/if}

          <!-- Section Divider: Berikutnya -->
          {#if currentIndex >= 0 && i === currentIndex + 1}
            <div
              class="text-[10px] font-bold uppercase tracking-wider text-zinc-400 px-1 pt-3.5 pb-1 flex items-center justify-between"
            >
              <span>{$t("queue_upcoming")}</span>
              <span class="text-[10px] font-normal lowercase text-zinc-500"
                >{items.length - currentIndex - 1}
                {$t("queue_songs_remaining")}</span
              >
            </div>
          {/if}

          <div
            class="group flex items-center gap-3 px-3 py-2.5 rounded-2xl cursor-pointer transition-all duration-150 {i ===
            currentIndex
              ? 'bg-gradient-to-r from-crimson-500/20 via-crimson-500/10 to-transparent border border-crimson-500/35 shadow-sm shadow-crimson-950/20'
              : 'hover:bg-white/[0.05] border border-transparent'} {draggedIndex ===
            i
              ? 'opacity-40 scale-[0.98]'
              : ''}"
            draggable="true"
            on:dragstart={(e) => handleDragStart(i, e)}
            on:dragover={(e) => handleDragOver(i, e)}
            on:drop={(e) => handleDrop(i, e)}
            on:dragend={handleDragEnd}
            on:click={() => handlePlayQueueItem(i)}
            role="button"
            tabindex="0"
            on:keydown={(e) => e.key === "Enter" && handlePlayQueueItem(i)}
          >
            <!-- Index / Animated Equalizer / Play icon on hover -->
            <div class="w-6 flex items-center justify-center shrink-0">
              {#if i === currentIndex && isPlaying}
                <div class="flex items-center gap-0.5">
                  <span
                    class="inline-block w-0.5 h-3 bg-crimson-400 rounded-xs animate-pulse"
                  ></span>
                  <span
                    class="inline-block w-0.5 h-4 bg-crimson-500 rounded-xs animate-bounce"
                  ></span>
                  <span
                    class="inline-block w-0.5 h-2 bg-crimson-400 rounded-xs animate-pulse"
                  ></span>
                </div>
              {:else if i === currentIndex}
                <span
                  class="w-2.5 h-2.5 rounded-full bg-crimson-400 shadow-[0_0_8px_#FF4D62]"
                ></span>
              {:else}
                <span
                  class="text-xs font-mono text-zinc-500 group-hover:hidden transition-colors"
                  >{i + 1}</span
                >
                <svg
                  class="w-3.5 h-3.5 text-zinc-300 hidden group-hover:block transition-colors"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <polygon points="6 3 20 12 6 21 6 3"></polygon>
                </svg>
              {/if}
            </div>

            <!-- Thumbnail -->
            <div
              class="w-12 h-12 rounded-xl overflow-hidden bg-surface-highlight shrink-0 flex items-center justify-center border {i ===
              currentIndex
                ? 'border-crimson-500/40 shadow-sm shadow-crimson-500/10'
                : 'border-white/10 shadow-sm'} group-hover:scale-105 transition-transform duration-200"
            >
              {#if item.thumbnailUrl}
                <img
                  src={item.thumbnailUrl}
                  alt={item.title}
                  class="w-full h-full object-cover"
                  on:error={(e) => handleCoverError(e, item.id)}
                />
              {:else}
                <svg
                  class="w-5 h-5 text-zinc-500"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"
                  />
                </svg>
              {/if}
            </div>

            <!-- Track Info -->
            <div class="flex-1 min-w-0 pr-1">
              <span
                class="text-xs font-semibold truncate block {i === currentIndex
                  ? 'text-white font-bold'
                  : 'text-zinc-200 group-hover:text-white transition-colors'}"
                title={item.title}
              >
                {item.title}
              </span>
              <span class="text-[11px] text-zinc-400 truncate block mt-0.5"
                >{item.uploaderName}</span
              >
            </div>

            <!-- Duration & Actions -->
            <div class="flex items-center gap-1 shrink-0">
              <span
                class="text-[11px] font-mono text-zinc-500 group-hover:hidden"
              >
                {item.durationFormatted || ""}
              </span>
              <!-- Grip Drag Indicator -->
              <span
                class="p-1.5 text-zinc-500 hover:text-white rounded-lg hidden group-hover:flex items-center justify-center cursor-grab active:cursor-grabbing transition"
                title={$t("queue_drag_tip")}
              >
                <svg
                  class="w-3.5 h-3.5"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  viewBox="0 0 24 24"
                >
                  <circle cx="9" cy="7" r="1.2" fill="currentColor"></circle>
                  <circle cx="15" cy="7" r="1.2" fill="currentColor"></circle>
                  <circle cx="9" cy="12" r="1.2" fill="currentColor"></circle>
                  <circle cx="15" cy="12" r="1.2" fill="currentColor"></circle>
                  <circle cx="9" cy="17" r="1.2" fill="currentColor"></circle>
                  <circle cx="15" cy="17" r="1.2" fill="currentColor"></circle>
                </svg>
              </span>
              <button
                type="button"
                class="p-1.5 text-zinc-500 hover:text-crimson-400 hover:bg-crimson-500/10 rounded-lg hidden group-hover:flex transition cursor-pointer"
                on:click={(e) => handleRemoveItem(e, i)}
                title={$t("queue_remove_tip")}
                aria-label={$t("queue_remove_tip")}
              >
                <svg
                  class="w-3.5 h-3.5"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
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
