<script lang="ts">
  import type { SearchResultItem } from '../types'
  import TrackItem from './TrackItem.svelte'
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'

  export let title = 'Daftar Lagu'
  export let tracks: SearchResultItem[] = []
  export let isLoading = false
  export let emptyMessage = 'Tidak ada lagu yang ditemukan.'

  function handlePlayAll() {
    if (tracks.length === 0) return
    queueStore.setQueue(tracks, 0)
    playerStore.playTrack(tracks[0])
  }

  function handleShuffleAll() {
    if (tracks.length === 0) return
    const shuffled = [...tracks].sort(() => Math.random() - 0.5)
    queueStore.setQueue(shuffled, 0)
    playerStore.playTrack(shuffled[0])
  }
</script>

<div class="w-full flex flex-col select-none">
  <!-- Search Results Header -->
  <section class="px-8 pt-6 pb-4">
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
      <div>
        <div class="flex items-center gap-2 mb-1">
          <span class="text-xs uppercase tracking-widest font-semibold text-crimson-400">Audio Eksplorasi</span>
          <span class="text-zinc-600">•</span>
          <span class="text-xs text-zinc-400">Stream Master Audio</span>
        </div>
        <h1 class="text-2xl md:text-3xl font-extrabold text-white tracking-tight flex items-center gap-3">
          <span>{title}</span>
          {#if !isLoading && tracks.length > 0}
            <span class="text-xs font-mono font-medium px-2.5 py-1 rounded-full bg-white/[0.08] text-zinc-300 border border-white/10">
              {tracks.length} item
            </span>
          {/if}
        </h1>
      </div>

      <!-- Play All & Shuffle Controls -->
      {#if !isLoading && tracks.length > 0}
        <div class="flex items-center gap-2.5">
          <button
            class="px-5 py-2.5 rounded-full bg-crimson-500 hover:bg-crimson-600 text-white font-bold text-sm flex items-center gap-2 shadow-lg shadow-crimson-500/30 transition transform active:scale-95 cursor-pointer"
            on:click={handlePlayAll}
          >
            <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
            <span>Putar Semua</span>
          </button>
          <button
            class="px-4 py-2.5 rounded-full bg-white/[0.06] hover:bg-white/[0.1] border border-white/10 text-zinc-300 font-semibold text-sm flex items-center gap-2 transition active:scale-95 cursor-pointer"
            on:click={handleShuffleAll}
          >
            <svg class="w-4 h-4 text-zinc-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <polyline points="16 3 21 3 21 8"></polyline>
              <line x1="4" x2="21" y1="20" y2="3"></line>
              <polyline points="21 16 21 21 16 21"></polyline>
              <line x1="15" x2="21" y1="15" y2="21"></line>
              <line x1="4" x2="9" y1="4" y2="9"></line>
            </svg>
            <span>Acak</span>
          </button>
        </div>
      {/if}
    </div>
  </section>

  <!-- Track Table Section -->
  <section class="px-8 mt-2">
    <!-- Table Column Headers -->
    <div class="grid grid-cols-12 gap-4 px-4 py-2.5 text-[11px] font-semibold text-zinc-400 uppercase tracking-wider border-b border-white/[0.08]">
      <div class="col-span-1 text-center">#</div>
      <div class="col-span-6 md:col-span-5">Judul Lagu & Penyanyi</div>
      <div class="hidden md:block col-span-3">Channel / Penyanyi</div>
      <div class="col-span-3 md:col-span-2 text-right">Kualitas & Durasi</div>
      <div class="col-span-2 md:col-span-1 text-right pr-2">Aksi</div>
    </div>

    <!-- Track List Items Container -->
    {#if isLoading}
      <div class="divide-y divide-white/[0.03] mt-2 font-medium">
        {#each Array(6) as _}
          <div class="grid grid-cols-12 gap-4 px-4 py-3.5 items-center">
            <div class="col-span-1 flex justify-center">
              <div class="w-4 h-4 rounded bg-white/[0.05] animate-pulse"></div>
            </div>
            <div class="col-span-6 md:col-span-5 flex items-center gap-3.5">
              <div class="w-11 h-11 rounded-xl bg-white/[0.05] animate-pulse shrink-0"></div>
              <div class="flex-1 space-y-1.5">
                <div class="w-3/4 h-3.5 rounded bg-white/[0.05] animate-pulse"></div>
                <div class="w-1/2 h-2.5 rounded bg-white/[0.05] animate-pulse"></div>
              </div>
            </div>
            <div class="hidden md:block col-span-3">
              <div class="w-2/3 h-3 rounded bg-white/[0.05] animate-pulse"></div>
            </div>
            <div class="col-span-3 md:col-span-2 flex justify-end">
              <div class="w-16 h-3 rounded bg-white/[0.05] animate-pulse"></div>
            </div>
            <div class="col-span-2 md:col-span-1 flex justify-end">
              <div class="w-8 h-3 rounded bg-white/[0.05] animate-pulse"></div>
            </div>
          </div>
        {/each}
      </div>
    {:else if tracks.length === 0}
      <div class="flex flex-col items-center justify-center py-20 text-center gap-3">
        <div class="w-12 h-12 rounded-2xl bg-white/[0.04] border border-white/10 flex items-center justify-center text-zinc-500">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" stroke-width="1.8" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="8" y1="12" x2="16" y2="12"></line>
          </svg>
        </div>
        <p class="text-sm font-medium text-zinc-300">{emptyMessage}</p>
        <span class="text-xs text-zinc-500">Coba ketik kata kunci lain pada kolom pencarian di atas</span>
      </div>
    {:else}
      <div class="divide-y divide-white/[0.03] mt-1 font-medium">
        {#each tracks as track, i (track.id + '-' + i)}
          <TrackItem {track} index={i} />
        {/each}
      </div>
    {/if}
  </section>
</div>
