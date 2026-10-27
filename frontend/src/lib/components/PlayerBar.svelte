<script lang="ts">
  import { onMount } from 'svelte'
  import { playerStore } from '../stores/playerStore'
  import { queueStore } from '../stores/queueStore'
  import { userDataStore } from '../stores/userDataStore'

  let audioElement: HTMLAudioElement

  $: currentTrack = $playerStore.currentTrack
  $: stream = $playerStore.streamResult
  $: isPlaying = $playerStore.isPlaying
  $: currentTime = $playerStore.currentTime
  $: duration = $playerStore.duration
  $: volume = $playerStore.volume
  $: isMuted = $playerStore.isMuted
  $: mode = $playerStore.mode
  $: isLoading = $playerStore.isLoading
  $: isFav = currentTrack ? $userDataStore.favorites.some(f => f.id === currentTrack?.id) : false

  $: repeatMode = $queueStore.repeatMode
  $: isShuffle = $queueStore.isShuffle
  $: queueCount = $queueStore.items.length

  $: progressPercent = duration > 0 ? Math.min(100, Math.max(0, (currentTime / duration) * 100)) : 0
  $: volumePercent = isMuted ? 0 : Math.round(volume * 100)

  onMount(() => {
    if (audioElement) {
      playerStore.bindMedia(audioElement, null)
    }
  })

  // Synchronize audio element playback when stream or isPlaying changes in audio mode
  $: if (audioElement && stream?.audioUrl && mode === 'audio') {
    if (audioElement.src !== stream.audioUrl) {
      audioElement.src = stream.audioUrl
      if (isPlaying) {
        audioElement.play().catch(e => console.warn('Audio play error:', e))
      }
    }
  }

  function handleTimeUpdate() {
    if (audioElement && mode === 'audio') {
      playerStore.onTimeUpdate(audioElement.currentTime, audioElement.duration)
    }
  }

  function handleSeek(e: Event) {
    const val = parseFloat((e.target as HTMLInputElement).value)
    playerStore.seek(val)
  }

  function handleVolumeChange(e: Event) {
    const val = parseFloat((e.target as HTMLInputElement).value)
    playerStore.setVolume(val)
  }

  function formatTime(sec: number): string {
    if (!sec || isNaN(sec) || sec <= 0) return '00:00'
    const m = Math.floor(sec / 60)
    const s = Math.floor(sec % 60)
    return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  }
</script>

<!-- Hidden Audio Element -->
<audio
  bind:this={audioElement}
  on:timeupdate={handleTimeUpdate}
  on:ended={playerStore.onTrackEnded}
  on:error={playerStore.handleMediaError}
  preload="auto"
></audio>

<!-- Floating Modern Glassmorphic Island Player -->
<footer class="fixed bottom-4 left-6 right-6 md:left-72 md:right-8 z-40 pointer-events-none select-none">
  <div class="pointer-events-auto floating-island rounded-2xl border border-white/[0.12] shadow-island-glow px-4 py-3 flex items-center justify-between gap-4 transition-all duration-300">
    
    <!-- Sisi Kiri: Album art melengkung, interaktif favorit, detail lagu & bit-perfect info -->
    <div class="flex items-center gap-3.5 min-w-[200px] max-w-[280px] shrink-0">
      <div class="relative w-12 h-12 rounded-xl overflow-hidden shadow-md shadow-black/80 border border-white/10 shrink-0 bg-surface-highlight flex items-center justify-center">
        {#if currentTrack?.thumbnailUrl}
          <img alt={currentTrack.title} class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105" src={currentTrack.thumbnailUrl} />
        {:else}
          <svg class="w-6 h-6 text-zinc-500" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
          </svg>
        {/if}
      </div>

      <div class="min-w-0 flex flex-col justify-center">
        {#if currentTrack}
          <div class="flex items-center gap-1.5">
            <h4 class="text-sm font-bold text-white truncate hover:text-crimson-400 transition" title={currentTrack.title}>
              {currentTrack.title}
            </h4>
            <span class="text-[9px] font-mono px-1 py-0.2 rounded bg-amber-500/20 text-amber-300 font-semibold border border-amber-500/30 shrink-0">24-BIT</span>
          </div>
          <p class="text-xs text-zinc-400 truncate hover:text-zinc-200" title={currentTrack.uploaderName}>
            {currentTrack.uploaderName}
          </p>
          <div class="flex items-center gap-1 text-[10px] text-zinc-400 font-mono">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
            <span>Bit-Perfect • Direct Stream</span>
          </div>
        {:else}
          <span class="text-xs text-zinc-400">Belum ada lagu yang diputar</span>
          <span class="text-[10px] text-zinc-500">Pilih lagu untuk memutar</span>
        {/if}
      </div>

      {#if currentTrack}
        <button
          class="p-1.5 text-crimson-500 hover:text-crimson-400 hover:scale-110 active:scale-95 transition shrink-0 ml-0.5 cursor-pointer"
          on:click={() => currentTrack && userDataStore.toggleFavorite(currentTrack)}
          title={isFav ? 'Hapus dari favorit' : 'Tambah ke favorit'}
        >
          <svg class="w-4 h-4 {isFav ? 'fill-current drop-shadow-[0_0_6px_rgba(230,30,56,0.6)]' : 'text-zinc-400'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
          </svg>
        </button>
      {/if}
    </div>

    <!-- Bagian Tengah: Kontrol playback ergonomis, dual time indicator, & smooth scrub bar -->
    <div class="flex flex-col items-center justify-center flex-1 max-w-xl px-2">
      <!-- Playback Control Buttons -->
      <div class="flex items-center gap-4 mb-1">
        <!-- Shuffle -->
        <button
          class="p-1.5 transition active:scale-95 cursor-pointer {isShuffle ? 'text-crimson-400' : 'text-zinc-400 hover:text-zinc-200'}"
          on:click={queueStore.toggleShuffle}
          title="Acak Lagu"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <polyline points="16 3 21 3 21 8"></polyline>
            <line x1="4" x2="21" y1="20" y2="3"></line>
            <polyline points="21 16 21 21 16 21"></polyline>
            <line x1="15" x2="21" y1="15" y2="21"></line>
            <line x1="4" x2="9" y1="4" y2="9"></line>
          </svg>
        </button>

        <!-- Previous -->
        <button
          class="p-1.5 text-zinc-300 hover:text-white transition active:scale-90 cursor-pointer"
          on:click={playerStore.playPrev}
          title="Lagu Sebelumnya (J)"
        >
          <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
            <polygon points="19 20 9 12 19 4 19 20"></polygon>
            <line stroke="currentColor" stroke-width="2.5" x1="5" x2="5" y1="19" y2="5"></line>
          </svg>
        </button>

        <!-- Play / Pause Circle Accent Button -->
        <button
          class="w-9 h-9 rounded-full bg-crimson-500 hover:bg-crimson-600 text-white flex items-center justify-center shadow-lg shadow-crimson-500/40 hover:scale-105 active:scale-95 transition duration-150 cursor-pointer"
          on:click={playerStore.togglePlay}
          title="Play/Pause (Spasi)"
        >
          {#if isLoading}
            <div class="w-4 h-4 border-2 border-white/20 border-t-white rounded-full animate-spin"></div>
          {:else if isPlaying}
            <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
              <rect height="14" width="3.5" x="6" y="5"></rect>
              <rect height="14" width="3.5" x="14.5" y="5"></rect>
            </svg>
          {:else}
            <svg class="w-4 h-4 fill-current ml-0.5" viewBox="0 0 24 24">
              <polygon points="5 3 19 12 5 21 5 3"></polygon>
            </svg>
          {/if}
        </button>

        <!-- Next -->
        <button
          class="p-1.5 text-zinc-300 hover:text-white transition active:scale-90 cursor-pointer"
          on:click={playerStore.playNext}
          title="Lagu Selanjutnya (L)"
        >
          <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
            <polygon points="5 4 15 12 5 20 5 4"></polygon>
            <line stroke="currentColor" stroke-width="2.5" x1="19" x2="19" y1="5" y2="19"></line>
          </svg>
        </button>

        <!-- Repeat Mode -->
        <button
          class="p-1.5 transition relative active:scale-95 cursor-pointer {repeatMode !== 'off' ? 'text-crimson-400' : 'text-zinc-400 hover:text-zinc-200'}"
          on:click={queueStore.cycleRepeat}
          title="Ulangi: {repeatMode.toUpperCase()}"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <polyline points="17 1 21 5 17 9"></polyline>
            <path d="M3 11V9a4 4 0 0 1 4-4h14"></path>
            <polyline points="7 23 3 19 7 15"></polyline>
            <path d="M21 13v2a4 4 0 0 1-4 4H3"></path>
          </svg>
          {#if repeatMode === 'one'}
            <span class="absolute -top-1 -right-1 text-[9px] font-bold font-mono text-crimson-400">1</span>
          {:else if repeatMode === 'all'}
            <span class="absolute -bottom-0.5 right-1 w-1 h-1 rounded-full bg-crimson-400 shadow-[0_0_4px_#FF4D62]"></span>
          {/if}
        </button>
      </div>

      <!-- Dual time indicator & scrubber bar -->
      <div class="w-full flex items-center gap-2.5 text-[11px] font-mono font-medium text-zinc-400">
        <span class="w-9 text-right text-zinc-300 shrink-0">{formatTime(currentTime)}</span>
        
        <div class="relative flex-1 group py-1 flex items-center">
          <div class="h-1 w-full bg-white/[0.12] rounded-full overflow-hidden group-hover:h-1.5 transition-all relative">
            <div
              class="h-full bg-gradient-to-r from-crimson-500 to-rose-400 rounded-full shadow-[0_0_12px_rgba(230,30,56,0.6)]"
              style="width: {progressPercent}%"
            ></div>
          </div>
          <!-- Real seek input overlay -->
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

        <span class="w-9 text-zinc-400 shrink-0">{formatTime(duration)}</span>
      </div>
    </div>

    <!-- Sisi Kanan: Audio Device Switcher, Mode Toggle, Queue, dan Volume -->
    <div class="flex items-center justify-end gap-2.5 min-w-[200px] shrink-0">
      <!-- Mode Switcher Pill (Audio vs Video) -->
      <button
        class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl border text-xs font-semibold transition cursor-pointer {mode === 'video' ? 'bg-crimson-500/20 border-crimson-500/40 text-crimson-300' : 'bg-white/[0.06] hover:bg-white/[0.1] border-white/10 text-zinc-200'}"
        on:click={playerStore.toggleMode}
        title="Ganti Mode Audio / Video (V)"
      >
        <svg class="w-3.5 h-3.5 {mode === 'video' ? 'text-crimson-400' : 'text-zinc-400'}" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path d="M3 18v-6a9 9 0 0 1 18 0v6"></path>
          <path d="M21 19a2 2 0 0 1-2 2h-1a2 2 0 0 1-2-2v-3a2 2 0 0 1 2-2h3zM3 19a2 2 0 0 0 2 2h1a2 2 0 0 0 2-2v-3a2 2 0 0 0-2-2H3z"></path>
        </svg>
        <span class="text-[11px] hidden sm:inline">{mode === 'video' ? 'Video MV' : 'Hi-Fi Audio'}</span>
      </button>

      <!-- Queue Button with badge -->
      <button
        class="p-2 text-zinc-400 hover:text-white transition rounded-xl hover:bg-white/[0.06] relative cursor-pointer {$queueStore.isDrawerOpen ? 'bg-white/[0.1] text-white' : ''}"
        on:click={() => queueStore.toggleDrawer()}
        title="Antrean Pemutaran"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <line x1="8" x2="21" y1="6" y2="6"></line>
          <line x1="8" x2="21" y1="12" y2="12"></line>
          <line x1="8" x2="21" y1="18" y2="18"></line>
          <line x1="3" x2="3.01" y1="6" y2="6"></line>
          <line x1="3" x2="3.01" y1="12" y2="12"></line>
          <line x1="3" x2="3.01" y1="18" y2="18"></line>
        </svg>
        {#if queueCount > 0}
          <span class="absolute top-1 right-1 w-3.5 h-3.5 rounded-full bg-crimson-500 text-[9px] font-bold text-white flex items-center justify-center font-mono shadow-sm">
            {queueCount}
          </span>
        {/if}
      </button>

      <!-- Volume Slider Ramping -->
      <div class="flex items-center gap-1.5 group w-24">
        <button
          class="p-1 text-zinc-400 hover:text-white transition cursor-pointer"
          on:click={playerStore.toggleMute}
          title="Mute/Unmute (M)"
        >
          {#if isMuted || volume === 0}
            <svg class="w-3.5 h-3.5 text-crimson-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
              <line x1="23" y1="9" x2="17" y2="15"></line>
              <line x1="17" y1="9" x2="23" y2="15"></line>
            </svg>
          {:else}
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
              <path d="M19.07 4.93a10 10 0 0 1 0 14.14M15.54 8.46a5 5 0 0 1 0 7.07"></path>
            </svg>
          {/if}
        </button>

        <div class="relative flex-1 py-1 flex items-center">
          <div class="h-1 w-full bg-white/[0.15] rounded-full group-hover:h-1.5 transition-all relative overflow-hidden">
            <div
              class="h-full bg-zinc-200 group-hover:bg-crimson-400 rounded-full transition-colors"
              style="width: {volumePercent}%"
            ></div>
          </div>
          <input
            type="range"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            min="0"
            max="1"
            step="0.02"
            value={isMuted ? 0 : volume}
            on:input={handleVolumeChange}
          />
        </div>
      </div>
    </div>
  </div>
</footer>
