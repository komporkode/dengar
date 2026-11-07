<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import { handleCoverError } from '../stores/offlineStore'

  export let tag = ''
  export let title = ''
  export let subtitle = ''
  export let imageUrl = ''
  export let playLabel = 'Putar Sekarang'
  export let isPlaying = false
  export let showPlay = true
  export let showShuffle = true
  export let shuffleTitle = 'Putar Acak'
  export let trackId: string | undefined = undefined
  export let heightClass = 'h-[250px] sm:h-[265px] md:h-[275px]'

  const dispatch = createEventDispatcher<{
    play: void
    shuffle: void
  }>()

  let imgFailed = false

  function onImgError(e: Event) {
    imgFailed = true
    handleCoverError(e, trackId)
  }
</script>

<div
  class="relative overflow-hidden rounded-3xl {heightClass} bg-[#0A0B14] border border-white/[0.08] shadow-2xl p-6 sm:p-7 md:p-8 select-none transition-all group"
>
  <!-- Top-edge luminous accent line -->
  <div class="absolute inset-x-0 top-0 h-[1px] bg-gradient-to-r from-transparent via-white/15 to-transparent pointer-events-none z-20"></div>

  <!-- Ambient background glow if no image or fallback (Consistent Crimson brand palette) -->
  <div class="absolute -right-16 -top-16 w-96 h-96 bg-crimson-600/10 rounded-full blur-[140px] pointer-events-none z-0"></div>

  <!-- Full-height background cover image pinned to right, 100% height fit with zero crop & seamless dissolve -->
  {#if imageUrl && !imgFailed}
    <div
      class="absolute right-0 top-0 bottom-0 h-full pointer-events-none select-none overflow-hidden z-0 flex justify-end"
    >
      <div class="relative h-full flex items-center justify-end banner-img-container">
        <img
          src={imageUrl}
          alt=""
          class="block h-full w-auto max-w-none object-contain opacity-85 sm:opacity-95 filter blur-[0.1px] transition-transform duration-500 group-hover:scale-[1.01]"
          loading="lazy"
          on:error={onImgError}
        />
        <!-- Inner image vignette so the left edge dissolves smoothly into the card background -->
        <div class="absolute inset-0 banner-image-vignette pointer-events-none"></div>
      </div>
    </div>
  {/if}

  <!-- Banner Content (Left Column) with fixed standard height layout -->
  <div class="relative z-10 flex flex-col justify-between h-full max-w-xl sm:max-w-2xl">
    <div class="space-y-1.5 sm:space-y-2">
      {#if tag}
        <div class="inline-flex items-center gap-2">
          <span class="text-[10px] sm:text-[11px] font-bold tracking-[0.16em] uppercase text-zinc-400 font-mono">
            {tag}
          </span>
        </div>
      {/if}

      <h1 class="text-2xl sm:text-3xl md:text-3.5xl font-black text-white tracking-tight leading-[1.18] drop-shadow-md line-clamp-2">
        {title}
      </h1>

      {#if subtitle}
        <p class="text-xs sm:text-sm text-zinc-300 font-medium leading-relaxed line-clamp-1 max-w-lg">
          {subtitle}
        </p>
      {/if}
    </div>

    <!-- Minimalist Controls Suite (Pill + Circle Glass Buttons) -->
    <div class="flex items-center gap-3 pt-2 flex-wrap">
      {#if showPlay}
        <button
          type="button"
          class="h-11 sm:h-12 px-6 sm:px-7 rounded-full bg-crimson-600 hover:bg-crimson-500 border border-transparent text-white font-bold text-xs sm:text-sm inline-flex items-center gap-2.5 shadow-xl shadow-crimson-600/35 hover:shadow-crimson-600/50 hover:scale-[1.02] active:scale-95 transition-all cursor-pointer whitespace-nowrap shrink-0"
          on:click={() => dispatch('play')}
          title={playLabel}
          aria-label={playLabel}
        >
          {#if isPlaying}
            <svg class="w-4 h-4 fill-current shrink-0" viewBox="0 0 24 24">
              <rect height="14" width="3.5" x="6" y="5"></rect>
              <rect height="14" width="3.5" x="14.5" y="5"></rect>
            </svg>
            <span>Jeda</span>
          {:else}
            <svg class="w-4 h-4 fill-current shrink-0 ml-0.5" viewBox="0 0 24 24">
              <polygon points="5 3 19 12 5 21 5 3"></polygon>
            </svg>
            <span>{playLabel}</span>
          {/if}
        </button>
      {/if}

      {#if showShuffle}
        <button
          type="button"
          class="w-11 h-11 sm:w-12 sm:h-12 rounded-full bg-white/[0.08] hover:bg-white/[0.14] border border-white/10 hover:border-white/20 text-zinc-200 hover:text-white flex items-center justify-center transition-all hover:scale-105 active:scale-95 cursor-pointer shrink-0"
          on:click={() => dispatch('shuffle')}
          title={shuffleTitle}
          aria-label={shuffleTitle}
        >
          <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <polyline points="16 3 21 3 21 8"></polyline>
            <line x1="4" y1="20" x2="21" y2="3"></line>
            <polyline points="21 16 21 21 16 21"></polyline>
            <line x1="15" y1="15" x2="21" y2="21"></line>
            <line x1="4" y1="4" x2="9" y2="9"></line>
          </svg>
        </button>
      {/if}

      <!-- Slot for custom extra actions (e.g. Refresh, Delete, Share) -->
      <slot />
    </div>
  </div>
</div>

<style>
  .banner-img-container {
    -webkit-mask-image: linear-gradient(
      90deg,
      transparent 0%,
      rgba(0, 0, 0, 0.25) 8%,
      rgba(0, 0, 0, 0.85) 20%,
      black 32%
    );
    mask-image: linear-gradient(
      90deg,
      transparent 0%,
      rgba(0, 0, 0, 0.25) 8%,
      rgba(0, 0, 0, 0.85) 20%,
      black 32%
    );
  }

  .banner-image-vignette {
    background: linear-gradient(
      90deg,
      #0A0B14 0%,
      rgba(10, 11, 20, 0.5) 10%,
      transparent 24%
    ),
    linear-gradient(
      180deg,
      rgba(10, 11, 20, 0.25) 0%,
      transparent 15%,
      transparent 85%,
      rgba(10, 11, 20, 0.35) 100%
    );
  }
</style>
