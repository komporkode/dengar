<script lang="ts">
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'
  import { exitMiniPlayer } from '../stores/settingsStore'
  import { t } from '../stores/i18nStore'

  $: currentTrack = $playerStore.currentTrack
  $: isPlaying = $playerStore.isPlaying
  $: currentTime = $playerStore.currentTime
  $: duration = $playerStore.duration
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
</script>

<div class="w-full h-full bg-[#090A0E] text-zinc-100 flex flex-col justify-between p-3 select-none overflow-hidden border border-white/10 rounded-xl" style="--wails-draggable:drag">
  <!-- Top bar in Mini Player -->
  <div class="flex items-center justify-between gap-2" style="--wails-draggable:drag">
    <div class="flex items-center gap-1.5 min-w-0" style="--wails-draggable:no-drag">
      <div class="w-3.5 h-3.5 rounded bg-gradient-to-tr from-crimson-600 to-rose-500 flex items-center justify-center shrink-0">
        <svg class="w-2 h-2 text-white" viewBox="0 0 269.22 267.55" fill="none">
          <path fill="currentColor" d="M135.71 0c26.85,0 50.37,5.44 70.57,16.33 20.2,10.89 35.74,26 46.62,45.35 10.88,19.35 16.33,41.36 16.33,66.03 0,24.43 -5.44,46.44 -16.33,66.03 -10.88,19.59 -26.42,35.01 -46.62,46.26 -20.2,11.25 -43.72,16.87 -70.57,16.87l-26.48 0 -8.81 0c-10.22,6.75 -22.47,10.68 -35.64,10.68 -35.78,0 -64.78,-29 -64.78,-64.78 0,-22.41 11.38,-42.16 28.68,-53.79l0 -130.95c0,-9.91 8.11,-18.02 18.02,-18.02l43.87 0 18.66 0 26.48 0zm-107.03 209.99l0 -0.01c5.01,20.39 27.57,32.83 51.03,27.89 24.62,-5.18 38.29,-31.72 34.99,-52.84l13.75 0c18.14,0 32.53,-4.96 43.18,-14.88 10.64,-9.92 15.96,-24.07 15.96,-42.45 0,-18.38 -5.32,-32.53 -15.96,-42.45 -10.64,-9.92 -25.03,-14.87 -43.18,-14.87l-19.23 0 0 97.01c-16.53,-25.13 -72.69,-15.5 -80.54,25.92l0 -0c-0.07,0.3 -0.14,0.59 -0.2,0.89l-0 0.01 -0.01 0.04c-0.06,0.29 -0.12,0.59 -0.18,0.88l-0.01 0.03 -0.01 0.04c-0.17,0.93 -0.31,1.87 -0.41,2.82l-0 0.03 -0 0.03c-0.03,0.31 -0.06,0.63 -0.08,0.94l-0 0.03c-0.02,0.31 -0.04,0.62 -0.06,0.94l-0 0.01 -0 0.02c-0.02,0.32 -0.03,0.63 -0.03,0.95l-0 0.02c-0.01,0.32 -0.01,0.64 -0.01,0.96l0 0.02c0,0.32 0.01,0.64 0.02,0.96l0 0 0 0.02c0.01,0.32 0.03,0.64 0.05,0.97l0 0.01c0.02,0.32 0.05,0.64 0.08,0.95l0 0.02 0 0.01c0.03,0.32 0.07,0.65 0.11,0.98l0 0.01c0.04,0.32 0.08,0.63 0.13,0.95l0 0.03 0 0.01 0 0.02c0.05,0.31 0.1,0.63 0.16,0.94l0 0.01 0 0 0 0.01c0.06,0.32 0.12,0.65 0.19,0.97l0 0.02 0 0.02c0.06,0.3 0.13,0.6 0.2,0.9l0 0.02 0 0.02 0 0.01 0 0.01 0 0.02 0 0 0.03 0.14z"/>
        </svg>
      </div>
      <span class="text-[11px] font-bold text-zinc-300">Dengar</span>
      <span class="text-[9px] font-mono px-1 py-0.2 rounded bg-crimson-500/20 text-crimson-400">Mini</span>
    </div>

    <button
      class="p-1 rounded text-zinc-400 hover:text-white hover:bg-white/10 transition cursor-pointer"
      style="--wails-draggable:no-drag"
      on:click={exitMiniPlayer}
      title={$t('restore_player')}
    >
      <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="15 3 21 3 21 9"></polyline>
        <polyline points="9 21 3 21 3 15"></polyline>
        <line x1="21" y1="3" x2="14" y2="10"></line>
        <line x1="3" y1="21" x2="10" y2="14"></line>
      </svg>
    </button>
  </div>

  <!-- Middle: Track info & playback -->
  <div class="flex items-center justify-between gap-3 my-1" style="--wails-draggable:no-drag">
    <div class="flex items-center gap-2.5 min-w-0 flex-1">
      <div class="w-10 h-10 rounded-lg overflow-hidden bg-surface-highlight border border-white/10 shrink-0 flex items-center justify-center">
        {#if currentTrack?.thumbnailUrl}
          <img src={currentTrack.thumbnailUrl} alt={currentTrack.title} class="w-full h-full object-cover" />
        {:else}
          <svg class="w-4 h-4 text-zinc-500" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
          </svg>
        {/if}
      </div>

      <div class="min-w-0 flex-1">
        <p class="text-xs font-semibold text-white truncate" title={currentTrack?.title || $t('no_track_playing')}>
          {currentTrack?.title || $t('no_track_playing')}
        </p>
        <p class="text-[10px] text-zinc-400 truncate">
          {currentTrack?.uploaderName || 'Dengar Player'}
        </p>
      </div>
    </div>

    <!-- Controls -->
    <div class="flex items-center gap-2 shrink-0">
      <button class="p-1 text-zinc-300 hover:text-white transition cursor-pointer" on:click={playerStore.playPrev} title="Lagu Sebelumnya">
        <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24">
          <polygon points="19 20 9 12 19 4 19 20"></polygon>
          <line stroke="currentColor" stroke-width="2.5" x1="5" x2="5" y1="19" y2="5"></line>
        </svg>
      </button>

      <button
        class="w-7 h-7 rounded-full bg-crimson-500 hover:bg-crimson-600 text-white flex items-center justify-center transition cursor-pointer"
        on:click={playerStore.togglePlay}
        title="Play/Pause"
      >
        {#if isPlaying}
          <svg class="w-3 h-3 fill-current" viewBox="0 0 24 24">
            <rect height="12" width="3" x="7" y="6"></rect>
            <rect height="12" width="3" x="14" y="6"></rect>
          </svg>
        {:else}
          <svg class="w-3 h-3 fill-current ml-0.5" viewBox="0 0 24 24">
            <polygon points="5 3 19 12 5 21 5 3"></polygon>
          </svg>
        {/if}
      </button>

      <button class="p-1 text-zinc-300 hover:text-white transition cursor-pointer" on:click={playerStore.playNext} title="Lagu Selanjutnya">
        <svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24">
          <polygon points="5 4 15 12 5 20 5 4"></polygon>
          <line stroke="currentColor" stroke-width="2.5" x1="19" x2="19" y1="5" y2="19"></line>
        </svg>
      </button>
    </div>
  </div>

  <!-- Bottom: Mini Scrubber -->
  <div class="flex items-center gap-2 text-[10px] font-mono text-zinc-400" style="--wails-draggable:no-drag">
    <span class="w-7 text-right">{formatTime(currentTime)}</span>
    <div class="relative flex-1 py-1 flex items-center">
      <div class="h-1 w-full bg-white/[0.12] rounded-full overflow-hidden relative">
        <div class="h-full bg-crimson-500 rounded-full" style="width: {progressPercent}%"></div>
      </div>
      <input
        type="range"
        class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
        min="0"
        max={duration || 100}
        value={currentTime}
        on:input={handleSeek}
        disabled={!currentTrack}
      />
    </div>
    <span class="w-7">{formatTime(duration)}</span>
  </div>
</div>
