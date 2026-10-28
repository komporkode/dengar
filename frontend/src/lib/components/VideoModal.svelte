<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { playerStore } from '../stores/playerStore'

  let videoElement: HTMLVideoElement

  $: currentTrack = $playerStore.currentTrack
  $: stream = $playerStore.streamResult
  $: isPlaying = $playerStore.isPlaying
  $: mode = $playerStore.mode
  $: videoSrc = stream?.videoUrl || stream?.audioUrl || ''

  onMount(() => {
    if (videoElement) {
      playerStore.bindVideo(videoElement)
    }
  })

  onDestroy(() => {
    playerStore.bindVideo(null)
  })

  // Synchronize video element playback when stream or isPlaying changes in video mode
  $: if (videoElement && videoSrc && mode === 'video') {
    if (videoElement.src !== videoSrc) {
      videoElement.src = videoSrc
      videoElement.currentTime = $playerStore.currentTime
      if (isPlaying) {
        videoElement.play().catch(e => console.warn('Video play error:', e))
      }
    } else if (isPlaying && videoElement.paused) {
      videoElement.play().catch(e => console.warn('Video play error:', e))
    }
  }

  function handleTimeUpdate() {
    if (videoElement && mode === 'video') {
      playerStore.onTimeUpdate(videoElement.currentTime, videoElement.duration)
    }
  }

  function closeVideo() {
    playerStore.toggleMode()
  }
</script>

{#if mode === 'video' && currentTrack}
  <div class="video-overlay" role="dialog" aria-modal="true">
    <div class="video-header">
      <div class="video-meta">
        <span class="video-title text-truncate">{currentTrack.title}</span>
        <span class="video-artist text-truncate text-secondary">{currentTrack.uploaderName}</span>
      </div>

      <button class="close-btn" on:click={closeVideo} title="Tutup Tampilan Video (V)">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>

    <div class="video-player-box">
      <!-- svelte-ignore a11y-media-has-caption -->
      <video
        bind:this={videoElement}
        on:timeupdate={handleTimeUpdate}
        on:ended={playerStore.onTrackEnded}
        on:error={playerStore.handleMediaError}
        playsinline
      ></video>
    </div>
  </div>
{/if}

<style>
  .video-overlay {
    position: fixed;
    top: 24px;
    right: 24px;
    width: 480px;
    background: #000000;
    border: 1px solid var(--border-subtle);
    border-radius: 12px;
    box-shadow: 0 12px 36px rgba(0, 0, 0, 0.7);
    overflow: hidden;
    z-index: 90;
    display: flex;
    flex-direction: column;
    animation: fadeIn 0.2s ease;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: scale(0.96); }
    to { opacity: 1; transform: scale(1); }
  }

  .video-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 14px;
    background: rgba(18, 18, 18, 0.95);
    border-bottom: 1px solid var(--border-subtle);
  }

  .video-meta {
    display: flex;
    flex-direction: column;
    min-width: 0;
    flex: 1;
    margin-right: 12px;
  }

  .video-title {
    font-size: 13px;
    font-weight: 600;
  }

  .video-artist {
    font-size: 11px;
  }

  .close-btn {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    color: var(--text-muted);
  }

  .close-btn:hover {
    color: var(--text-primary);
    background: rgba(255, 255, 255, 0.1);
  }

  .video-player-box {
    width: 100%;
    aspect-ratio: 16 / 9;
    background: #000000;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  video {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }
</style>
