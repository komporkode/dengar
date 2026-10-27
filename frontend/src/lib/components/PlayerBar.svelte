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

<footer class="player-bar">
  <!-- Left: Track Info -->
  <div class="left-section">
    {#if currentTrack}
      <div class="thumb-container">
        {#if currentTrack.thumbnailUrl}
          <img src={currentTrack.thumbnailUrl} alt={currentTrack.title} />
        {:else}
          <div class="no-thumb">🎵</div>
        {/if}
      </div>
      <div class="track-info">
        <span class="track-title text-truncate" title={currentTrack.title}>{currentTrack.title}</span>
        <span class="track-artist text-truncate text-secondary" title={currentTrack.uploaderName}>{currentTrack.uploaderName}</span>
      </div>
      <button
        class="fav-btn"
        class:is-fav={isFav}
        on:click={() => currentTrack && userDataStore.toggleFavorite(currentTrack)}
        title={isFav ? 'Hapus dari favorit' : 'Tambah ke favorit'}
      >
        <svg width="18" height="18" viewBox="0 0 24 24" fill={isFav ? '#ff0000' : 'none'} stroke="currentColor" stroke-width="2">
          <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
        </svg>
      </button>
    {:else}
      <div class="empty-track text-muted">
        <span>Belum ada lagu yang diputar</span>
      </div>
    {/if}
  </div>

  <!-- Center: Controls & Timeline -->
  <div class="center-section">
    <div class="playback-controls">
      <!-- Shuffle -->
      <button
        class="ctrl-btn"
        class:active={isShuffle}
        on:click={queueStore.toggleShuffle}
        title="Acak Antrean"
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="16 3 21 3 21 8"></polyline>
          <line x1="4" y1="20" x2="21" y2="3"></line>
          <polyline points="21 16 21 21 16 21"></polyline>
          <line x1="15" y1="15" x2="21" y2="21"></line>
          <line x1="4" y1="4" x2="9" y2="9"></line>
        </svg>
      </button>

      <!-- Prev -->
      <button class="ctrl-btn" on:click={playerStore.playPrev} title="Lagu Sebelumnya (J)">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
          <polygon points="19 20 9 12 19 4 19 20"></polygon>
          <line x1="5" y1="19" x2="5" y2="5" stroke="currentColor" stroke-width="2"></line>
        </svg>
      </button>

      <!-- Play / Pause -->
      <button class="play-btn" on:click={playerStore.togglePlay} title="Play/Pause (Spasi)">
        {#if isLoading}
          <div class="mini-spinner"></div>
        {:else if isPlaying}
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
            <rect x="6" y="4" width="4" height="16"></rect>
            <rect x="14" y="4" width="4" height="16"></rect>
          </svg>
        {:else}
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
            <polygon points="5 3 19 12 5 21 5 3"></polygon>
          </svg>
        {/if}
      </button>

      <!-- Next -->
      <button class="ctrl-btn" on:click={playerStore.playNext} title="Lagu Berikutnya (L)">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
          <polygon points="5 4 15 12 5 20 5 4"></polygon>
          <line x1="19" y1="5" x2="19" y2="19" stroke="currentColor" stroke-width="2"></line>
        </svg>
      </button>

      <!-- Repeat -->
      <button
        class="ctrl-btn"
        class:active={repeatMode !== 'off'}
        on:click={queueStore.cycleRepeat}
        title="Ulangi: {repeatMode.toUpperCase()}"
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="17 1 21 5 17 9"></polyline>
          <path d="M3 11V9a4 4 0 0 1 4-4h14"></path>
          <polyline points="7 23 3 19 7 15"></polyline>
          <path d="M21 13v2a4 4 0 0 1-4 4H3"></path>
        </svg>
        {#if repeatMode === 'one'}
          <span class="repeat-one">1</span>
        {/if}
      </button>
    </div>

    <!-- Timeline Progress -->
    <div class="timeline-row">
      <span class="time-label">{formatTime(currentTime)}</span>
      <input
        type="range"
        class="seek-bar"
        min="0"
        max={duration || 100}
        value={currentTime}
        on:input={handleSeek}
        disabled={!currentTrack}
      />
      <span class="time-label text-muted">{formatTime(duration)}</span>
    </div>
  </div>

  <!-- Right: Video toggle, Volume, Queue -->
  <div class="right-section">
    <!-- Mode Switcher (Audio / Video) -->
    <button
      class="mode-badge-btn"
      class:video-mode={mode === 'video'}
      on:click={playerStore.toggleMode}
      title="Ganti Mode Audio/Video (V)"
    >
      {#if mode === 'video'}
        🎬 <span>Video</span>
      {:else}
        🎵 <span>Audio</span>
      {/if}
    </button>

    <!-- Volume -->
    <div class="volume-container">
      <button class="vol-btn" on:click={playerStore.toggleMute} title="Mute/Unmute (M)">
        {#if isMuted || volume === 0}
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
            <line x1="23" y1="9" x2="17" y2="15"></line>
            <line x1="17" y1="9" x2="23" y2="15"></line>
          </svg>
        {:else if volume < 0.5}
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
            <path d="M15.54 8.46a5 5 0 0 1 0 7.07"></path>
          </svg>
        {:else}
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
            <path d="M19.07 4.93a10 10 0 0 1 0 14.14M15.54 8.46a5 5 0 0 1 0 7.07"></path>
          </svg>
        {/if}
      </button>

      <input
        type="range"
        class="vol-bar"
        min="0"
        max="1"
        step="0.02"
        value={isMuted ? 0 : volume}
        on:input={handleVolumeChange}
      />
    </div>

    <!-- Queue Toggle Button -->
    <button
      class="queue-toggle-btn"
      class:active={$queueStore.isDrawerOpen}
      on:click={() => queueStore.toggleDrawer()}
      title="Antrean Lagu"
    >
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <line x1="8" y1="6" x2="21" y2="6"></line>
        <line x1="8" y1="12" x2="21" y2="12"></line>
        <line x1="8" y1="18" x2="21" y2="18"></line>
        <line x1="3" y1="6" x2="3.01" y2="6"></line>
        <line x1="3" y1="12" x2="3.01" y2="12"></line>
        <line x1="3" y1="18" x2="3.01" y2="18"></line>
      </svg>
      {#if queueCount > 0}
        <span class="queue-badge">{queueCount}</span>
      {/if}
    </button>
  </div>
</footer>

<style>
  .player-bar {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: var(--player-height);
    background: rgba(18, 18, 18, 0.88);
    backdrop-filter: blur(24px);
    -webkit-backdrop-filter: blur(24px);
    border-top: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 24px;
    z-index: 100;
  }

  /* Left Section */
  .left-section {
    display: flex;
    align-items: center;
    gap: 14px;
    width: 280px;
    min-width: 220px;
  }

  .thumb-container {
    width: 50px;
    height: 50px;
    border-radius: 8px;
    overflow: hidden;
    background: var(--bg-card);
    flex-shrink: 0;
  }

  .thumb-container img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .no-thumb {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
  }

  .track-info {
    display: flex;
    flex-direction: column;
    min-width: 0;
    flex: 1;
  }

  .track-title {
    font-size: 14px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .track-artist {
    font-size: 12px;
  }

  .fav-btn {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    color: var(--text-muted);
    transition: color 0.15s ease;
  }

  .fav-btn:hover {
    color: var(--text-primary);
  }

  .fav-btn.is-fav {
    color: var(--accent-red);
  }

  .empty-track {
    font-size: 13px;
  }

  /* Center Section */
  .center-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    flex: 1;
    max-width: 580px;
    padding: 0 16px;
  }

  .playback-controls {
    display: flex;
    align-items: center;
    gap: 18px;
  }

  .ctrl-btn {
    width: 34px;
    height: 34px;
    border-radius: 50%;
    color: var(--text-secondary);
    position: relative;
    transition: all 0.15s ease;
  }

  .ctrl-btn:hover {
    color: var(--text-primary);
    transform: scale(1.08);
  }

  .ctrl-btn.active {
    color: var(--accent-red);
  }

  .repeat-one {
    position: absolute;
    top: 6px;
    right: 4px;
    font-size: 9px;
    font-weight: 700;
  }

  .play-btn {
    width: 44px;
    height: 44px;
    border-radius: 50%;
    background: var(--text-primary);
    color: var(--bg-primary);
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.4);
    transition: transform 0.15s ease, background-color 0.15s ease;
  }

  .play-btn:hover {
    transform: scale(1.06);
    background: #ffffff;
  }

  .mini-spinner {
    width: 18px;
    height: 18px;
    border: 2px solid rgba(0, 0, 0, 0.2);
    border-top-color: var(--accent-red);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .timeline-row {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
  }

  .time-label {
    font-size: 11px;
    width: 36px;
    text-align: center;
    font-variant-numeric: tabular-nums;
  }

  .seek-bar {
    flex: 1;
  }

  /* Right Section */
  .right-section {
    display: flex;
    align-items: center;
    gap: 16px;
    width: 280px;
    min-width: 220px;
    justify-content: flex-end;
  }

  .mode-badge-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 5px 12px;
    border-radius: 14px;
    font-size: 12px;
    font-weight: 600;
    background: var(--bg-card);
    border: 1px solid var(--border-subtle);
    color: var(--text-secondary);
    transition: all 0.15s ease;
  }

  .mode-badge-btn:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  .mode-badge-btn.video-mode {
    background: rgba(255, 0, 0, 0.15);
    border-color: var(--accent-red);
    color: var(--accent-red);
  }

  .volume-container {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 120px;
  }

  .vol-btn {
    width: 28px;
    height: 28px;
    color: var(--text-secondary);
  }

  .vol-btn:hover {
    color: var(--text-primary);
  }

  .vol-bar {
    width: 80px;
  }

  .queue-toggle-btn {
    position: relative;
    width: 34px;
    height: 34px;
    border-radius: 8px;
    color: var(--text-secondary);
    background: var(--bg-card);
    border: 1px solid var(--border-subtle);
  }

  .queue-toggle-btn:hover, .queue-toggle-btn.active {
    color: var(--text-primary);
    background: var(--bg-hover);
    border-color: var(--border-focus);
  }

  .queue-badge {
    position: absolute;
    top: -4px;
    right: -4px;
    background: var(--accent-red);
    color: white;
    font-size: 10px;
    font-weight: 700;
    min-width: 16px;
    height: 16px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 4px;
  }
</style>
