<script lang="ts">
  import type { SearchResultItem } from '../types'
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'
  import { userDataStore } from '../stores/userDataStore'

  export let track: SearchResultItem
  export let index: number

  $: isCurrent = $playerStore.currentTrack?.id === track.id
  $: isPlaying = isCurrent && $playerStore.isPlaying
  $: isFav = $userDataStore.favorites.some(f => f.id === track.id)

  function handlePlay() {
    if (isCurrent) {
      playerStore.togglePlay()
    } else {
      playerStore.playTrack(track)
      queueStore.addToQueue(track)
    }
  }

  function handleAddToQueue(e: MouseEvent) {
    e.stopPropagation()
    queueStore.addToQueue(track)
  }

  function handleToggleFav(e: MouseEvent) {
    e.stopPropagation()
    userDataStore.toggleFavorite(track)
  }
</script>

<div
  class="track-row group grid grid-cols-12 gap-4 px-4 py-3 rounded-2xl items-center transition-all duration-150 cursor-pointer {isCurrent ? 'bg-crimson-500/[0.08] border border-crimson-500/25 hover:bg-crimson-500/[0.12]' : 'hover:bg-white/[0.04] border border-transparent'}"
  on:click={handlePlay}
  role="button"
  tabindex="0"
  on:keydown={(e) => e.key === 'Enter' && handlePlay()}
>
  <!-- Col 1: Index / Play Icon / Equalizer -->
  <div class="col-span-1 flex items-center justify-center relative">
    {#if isPlaying}
      <span class="track-index text-crimson-400 font-mono text-sm font-bold flex items-center gap-0.5">
        <span class="inline-block w-1 h-3.5 bg-crimson-500 rounded-sm animate-pulse"></span>
        <span class="inline-block w-1 h-4 bg-crimson-400 rounded-sm animate-bounce"></span>
        <span class="inline-block w-1 h-2 bg-crimson-500 rounded-sm animate-pulse"></span>
      </span>
      <button
        class="play-hover-btn absolute inset-0 m-auto w-7 h-7 rounded-full bg-crimson-500 text-white flex items-center justify-center opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all shadow-md cursor-pointer"
        on:click|stopPropagation={handlePlay}
        title="Jeda Lagu"
      >
        <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24">
          <rect height="16" width="4" x="6" y="4"></rect>
          <rect height="16" width="4" x="14" y="4"></rect>
        </svg>
      </button>
    {:else}
      <span class="track-index text-zinc-400 font-mono text-xs group-hover:opacity-0 group-hover:invisible transition-opacity">{index + 1}</span>
      <button
        class="play-hover-btn absolute inset-0 m-auto w-7 h-7 rounded-full bg-crimson-500 text-white flex items-center justify-center opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all shadow-md hover:scale-105 cursor-pointer"
        on:click|stopPropagation={handlePlay}
        title="Putar Lagu"
      >
        <svg class="w-3.5 h-3.5 fill-current ml-0.5" viewBox="0 0 24 24">
          <polygon points="5 3 19 12 5 21 5 3"></polygon>
        </svg>
      </button>
    {/if}
  </div>

  <!-- Col 2: Title, Artist, Thumbnail -->
  <div class="col-span-6 md:col-span-5 flex items-center gap-3.5 min-w-0">
    <div class="relative shrink-0 w-11 h-11 rounded-xl overflow-hidden shadow-md border border-white/10 group-hover:scale-105 transition-transform duration-200 bg-surface-highlight flex items-center justify-center">
      {#if track.thumbnailUrl}
        <img alt={track.title} class="w-full h-full object-cover" src={track.thumbnailUrl} loading="lazy" />
      {:else}
        <svg class="w-5 h-5 text-zinc-500" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
        </svg>
      {/if}
    </div>
    <div class="min-w-0 flex flex-col">
      <div class="flex items-center gap-2">
        <span class="text-sm font-semibold truncate {isCurrent ? 'text-crimson-400 font-bold' : 'text-zinc-100 group-hover:text-crimson-400 transition'}">
          {track.title}
        </span>
        {#if isPlaying}
          <span class="text-[9px] font-bold uppercase tracking-wider px-1.5 py-0.5 rounded bg-crimson-500/20 text-crimson-300 border border-crimson-500/30 shrink-0">Memutar</span>
        {/if}
      </div>
      <span class="text-xs text-zinc-400 truncate">{track.uploaderName}</span>
    </div>
  </div>

  <!-- Col 3: Channel / Uploader -->
  <div class="hidden md:block col-span-3 text-xs text-zinc-400 truncate">
    {track.uploaderName}
  </div>

  <!-- Col 4: Quality & Duration -->
  <div class="col-span-3 md:col-span-2 flex items-center justify-end gap-2 text-right">
    <span class="text-[10px] font-mono px-1.5 py-0.5 rounded {isCurrent ? 'bg-amber-500/10 text-amber-300 border border-amber-500/20' : 'bg-white/[0.06] text-zinc-300'}">
      {track.duration > 300 ? 'FLAC 24b' : 'Lossless'}
    </span>
    <span class="text-xs font-mono {isCurrent ? 'text-crimson-400 font-bold' : 'text-zinc-400'}">
      {track.durationFormatted || '00:00'}
    </span>
  </div>

  <!-- Col 5: Actions -->
  <div class="col-span-2 md:col-span-1 flex items-center justify-end gap-1.5 text-zinc-400">
    <button
      class="p-1.5 hover:text-crimson-400 transition cursor-pointer"
      on:click={handleToggleFav}
      title={isFav ? 'Hapus dari favorit' : 'Tambah ke favorit'}
    >
      <svg class="w-4 h-4 {isFav ? 'fill-crimson-500 text-crimson-500' : 'text-zinc-400'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
      </svg>
    </button>
    <button
      class="p-1.5 hover:text-white transition cursor-pointer"
      on:click={handleAddToQueue}
      title="Tambah ke Antrean"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <line x1="12" y1="5" x2="12" y2="19"></line>
        <line x1="5" y1="12" x2="19" y2="12"></line>
      </svg>
    </button>
  </div>
</div>
