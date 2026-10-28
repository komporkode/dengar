<script lang="ts">
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'
  import { userDataStore } from '../stores/userDataStore'
  import {
    exitMiniPlayer,
    isMiniPlayerPinned,
    togglePinMiniPlayer,
    toastMessage,
    showToast
  } from '../stores/settingsStore'
  import { t } from '../stores/i18nStore'

  $: currentTrack = $playerStore.currentTrack
  $: isPlaying = $playerStore.isPlaying
  $: isLoading = $playerStore.isLoading
  $: currentTime = $playerStore.currentTime
  $: duration = $playerStore.duration
  $: isMuted = $playerStore.isMuted
  $: volume = $playerStore.volume
  $: isFav = currentTrack ? $userDataStore.favorites.some(f => f.id === currentTrack?.id) : false
  $: isPinned = $isMiniPlayerPinned

  $: progressPercent = duration > 0 ? Math.min(100, Math.max(0, (currentTime / duration) * 100)) : 0

  function formatTime(sec: number): string {
    if (!sec || isNaN(sec) || sec <= 0) return '00:00'
    const m = Math.floor(sec / 60)
    const s = Math.floor(sec % 60)
    return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  }

  function handleSeek(e: Event) {
    const val = parseFloat((e.target as HTMLInputElement).value)
    playerStore.seek(val)
  }

  function handleToggleFav() {
    if (!currentTrack) return
    const wasFav = isFav
    userDataStore.toggleFavorite(currentTrack)
    showToast(wasFav ? $t('removed_from_favorites') : $t('added_to_favorites'))
  }
</script>

<div
  class="w-full h-full bg-gradient-to-br from-[#121422] via-[#0A0C14] to-[#06070A] text-zinc-100 flex flex-col justify-between p-2.5 sm:p-3 select-none overflow-hidden relative border border-white/[0.12] rounded-2xl shadow-2xl shadow-black/80"
  style="--wails-draggable:drag"
>
  <!-- Ambient subtle glowing accent in background -->
  <div class="absolute -top-10 -left-10 w-32 h-32 bg-crimson-600/10 rounded-full blur-2xl pointer-events-none"></div>

  <!-- Top Bar: Brand, Live Playing Visualizer, Pin & Expand Controls -->
  <header class="flex items-center justify-between gap-2 shrink-0" style="--wails-draggable:drag">
    <!-- Brand & Playing Visualizer -->
    <div class="flex items-center gap-1.5 min-w-0" style="--wails-draggable:drag">
      <div class="w-4 h-4 rounded-md bg-gradient-to-tr from-crimson-600 to-rose-500 flex items-center justify-center shadow-xs shrink-0">
        <svg class="w-2.5 h-2.5 text-white" viewBox="0 0 269.22 267.55" fill="none">
          <path fill="currentColor" d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"/>
        </svg>
      </div>
      <span class="text-xs font-bold text-zinc-200 tracking-wide">Dengar</span>

      {#if isPlaying}
        <div class="flex items-center gap-0.5 ml-1">
          <span class="inline-block w-0.5 h-2.5 bg-crimson-400 rounded-xs animate-pulse"></span>
          <span class="inline-block w-0.5 h-3.5 bg-crimson-500 rounded-xs animate-bounce"></span>
          <span class="inline-block w-0.5 h-2 bg-crimson-400 rounded-xs animate-pulse"></span>
        </div>
      {:else}
        <span class="text-[9px] font-mono px-1.5 py-0.2 rounded bg-white/[0.06] text-zinc-400 border border-white/5">Mini</span>
      {/if}
    </div>

    <!-- Window Controls: Pin (Always on top) & Expand/Restore -->
    <div class="flex items-center gap-1.5 shrink-0" style="--wails-draggable:no-drag">
      <!-- Pin / Always on Top Button -->
      <button
        class="flex items-center gap-1 px-2 py-0.5 sm:py-1 rounded-lg border text-xs transition cursor-pointer {isPinned ? 'bg-crimson-500/20 border-crimson-500/40 text-crimson-300 shadow-[0_0_8px_rgba(230,30,56,0.3)]' : 'bg-white/[0.04] hover:bg-white/[0.09] border-white/10 text-zinc-400 hover:text-zinc-200'}"
        on:click={togglePinMiniPlayer}
        title={isPinned ? 'Lepas Sematan (Always on Top Aktif)' : 'Sematkan di Atas (Always on Top)'}
      >
        <svg class="w-3 h-3 {isPinned ? 'fill-current' : 'fill-none stroke-current'}" viewBox="0 0 24 24" stroke-width="2">
          <path d="M16 12V4h1V2H7v2h1v8l-2 2v2h5.2v6l1 1 1-1v-6H18v-2l-2-2z"/>
        </svg>
        <span class="text-[10px] font-medium hidden sm:inline">{isPinned ? 'Disematkan' : 'Sematkan'}</span>
      </button>

      <!-- Expand / Restore to Full Window -->
      <button
        class="p-1 sm:p-1.5 rounded-lg bg-white/[0.04] hover:bg-white/[0.1] border border-white/10 text-zinc-400 hover:text-white transition cursor-pointer"
        on:click={exitMiniPlayer}
        title="Kembali ke Pemutar Penuh"
      >
        <svg class="w-3 h-3 sm:w-3.5 sm:h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="15 3 21 3 21 9"></polyline>
          <polyline points="9 21 3 21 3 15"></polyline>
          <line x1="21" y1="3" x2="14" y2="10"></line>
          <line x1="3" y1="21" x2="10" y2="14"></line>
        </svg>
      </button>
    </div>
  </header>

  <!-- Middle: Track Info, Favorite & Ergonomic Playback Controls -->
  <div class="flex items-center justify-between gap-2.5 my-1" style="--wails-draggable:drag">
    <!-- Album Art & Track Meta -->
    <div class="flex items-center gap-2.5 min-w-0 flex-1" style="--wails-draggable:no-drag">
      <div class="relative w-10 h-10 sm:w-11 sm:h-11 rounded-xl overflow-hidden shadow-md shadow-black/80 border border-white/10 shrink-0 bg-surface-highlight flex items-center justify-center">
        {#if currentTrack?.thumbnailUrl}
          <img src={currentTrack.thumbnailUrl} alt={currentTrack.title} class="w-full h-full object-cover" />
        {:else}
          <svg class="w-5 h-5 text-zinc-500" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
          </svg>
        {/if}
      </div>

      <div class="min-w-0 flex-1 flex flex-col justify-center">
        <p class="text-xs font-bold text-white truncate hover:text-crimson-400 transition" title={currentTrack?.title || $t('no_track_playing')}>
          {currentTrack?.title || $t('no_track_playing')}
        </p>
        <p class="text-[10px] text-zinc-400 truncate mt-0.5" title={currentTrack?.uploaderName}>
          {currentTrack?.uploaderName || 'Dengar Player'}
        </p>
      </div>

      <!-- Quick Favorite Toggle -->
      {#if currentTrack}
        <button
          class="p-1 text-crimson-500 hover:text-crimson-400 hover:scale-110 active:scale-95 transition shrink-0 cursor-pointer"
          on:click={handleToggleFav}
          title={isFav ? 'Hapus dari favorit' : 'Tambah ke favorit'}
        >
          <svg class="w-3.5 h-3.5 {isFav ? 'fill-current drop-shadow-[0_0_6px_rgba(230,30,56,0.6)]' : 'text-zinc-500'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
          </svg>
        </button>
      {/if}
    </div>

    <!-- Controls -->
    <div class="flex items-center gap-1 sm:gap-1.5 shrink-0" style="--wails-draggable:no-drag">
      <!-- Previous -->
      <button
        class="p-1.5 text-zinc-300 hover:text-white active:scale-90 transition cursor-pointer"
        on:click={playerStore.playPrev}
        title="Lagu Sebelumnya (J)"
      >
        <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24">
          <polygon points="19 20 9 12 19 4 19 20"></polygon>
          <line stroke="currentColor" stroke-width="2.5" x1="5" x2="5" y1="19" y2="5"></line>
        </svg>
      </button>

      <!-- Play / Pause Accent Button -->
      <button
        class="w-7 h-7 sm:w-8 sm:h-8 rounded-full bg-crimson-500 hover:bg-crimson-600 text-white flex items-center justify-center shadow-lg shadow-crimson-500/40 hover:scale-105 active:scale-95 transition cursor-pointer shrink-0"
        on:click={playerStore.togglePlay}
        title="Play/Pause (Spasi)"
      >
        {#if isLoading}
          <div class="w-3.5 h-3.5 border-2 border-white/20 border-t-white rounded-full animate-spin"></div>
        {:else if isPlaying}
          <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24">
            <rect height="12" width="3" x="7" y="6"></rect>
            <rect height="12" width="3" x="14" y="6"></rect>
          </svg>
        {:else}
          <svg class="w-3.5 h-3.5 fill-current ml-0.5" viewBox="0 0 24 24">
            <polygon points="5 3 19 12 5 21 5 3"></polygon>
          </svg>
        {/if}
      </button>

      <!-- Next -->
      <button
        class="p-1.5 text-zinc-300 hover:text-white active:scale-90 transition cursor-pointer"
        on:click={playerStore.playNext}
        title="Lagu Selanjutnya (L)"
      >
        <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24">
          <polygon points="5 4 15 12 5 20 5 4"></polygon>
          <line stroke="currentColor" stroke-width="2.5" x1="19" x2="19" y1="5" y2="19"></line>
        </svg>
      </button>
    </div>
  </div>

  <!-- Bottom: Mini Scrubber & Quick Mute -->
  <div class="flex items-center gap-2 text-[10px] font-mono text-zinc-400" style="--wails-draggable:no-drag">
    <span class="w-7 text-right text-zinc-300 shrink-0">{formatTime(currentTime)}</span>

    <div class="relative flex-1 group py-1 flex items-center">
      <div class="h-1 w-full bg-white/[0.12] rounded-full overflow-hidden group-hover:h-1.5 transition-all relative">
        <div
          class="h-full bg-gradient-to-r from-crimson-500 to-rose-400 rounded-full shadow-[0_0_8px_rgba(230,30,56,0.6)]"
          style="width: {progressPercent}%"
        ></div>
      </div>
      <input
        type="range"
        class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
        min="0"
        max={duration || 100}
        step="0.1"
        value={currentTime}
        on:input={handleSeek}
        disabled={!currentTrack}
      />
    </div>

    <span class="w-7 text-zinc-400 shrink-0">{formatTime(duration)}</span>

    <!-- Quick Mute / Unmute -->
    <button
      class="p-1 text-zinc-400 hover:text-white transition cursor-pointer shrink-0"
      on:click={playerStore.toggleMute}
      title={isMuted || volume === 0 ? 'Bunyikan (M)' : 'Bisukan (M)'}
    >
      {#if isMuted || volume === 0}
        <svg class="w-3 h-3 text-crimson-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
          <line x1="23" y1="9" x2="17" y2="15"></line>
          <line x1="17" y1="9" x2="23" y2="15"></line>
        </svg>
      {:else}
        <svg class="w-3 h-3" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
          <path d="M19.07 4.93a10 10 0 0 1 0 14.14M15.54 8.46a5 5 0 0 1 0 7.07"></path>
        </svg>
      {/if}
    </button>
  </div>

  <!-- Floating Toast Pill inside Mini Player for immediate visual confirmation -->
  {#if $toastMessage}
    <div class="absolute bottom-2 left-1/2 -translate-x-1/2 z-50 pointer-events-none animate-slide-up duration-150">
      <div class="px-2.5 py-0.5 rounded-full bg-[#181A26]/95 border border-white/10 shadow-lg text-[9px] font-medium text-white flex items-center gap-1.5 backdrop-blur-md whitespace-nowrap">
        <span class="w-1.5 h-1.5 rounded-full bg-crimson-500 shadow-[0_0_6px_#E61E38]"></span>
        <span>{$toastMessage}</span>
      </div>
    </div>
  {/if}
</div>
