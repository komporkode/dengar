import { playerStore } from '../stores/playerStore'
import type { SearchResultItem } from '../types'

/**
 * Service to integrate Hardware Media Keys and Windows System Media Transport Controls (SMTC)
 * via standard Web MediaSession API in WebView2.
 */
class MediaSessionService {
  private isInitialized = false
  private lastTrackId: string | null = null

  public init() {
    if (this.isInitialized || typeof navigator === 'undefined' || !('mediaSession' in navigator)) {
      return
    }

    this.setupActionHandlers()
    this.subscribeToPlayer()
    this.isInitialized = true
  }

  private setupActionHandlers() {
    const ms = navigator.mediaSession

    const safeSetHandler = (action: MediaSessionAction, handler: MediaSessionActionHandler) => {
      try {
        ms.setActionHandler(action, handler)
      } catch (err) {
        // Some browser engines / webview versions might not support specific optional actions
        console.warn(`[MediaSession] Action handler '${action}' not supported:`, err)
      }
    }

    safeSetHandler('play', () => {
      playerStore.togglePlay()
    })

    safeSetHandler('pause', () => {
      playerStore.togglePlay()
    })

    safeSetHandler('nexttrack', () => {
      playerStore.playNext()
    })

    safeSetHandler('previoustrack', () => {
      playerStore.playPrev()
    })

    safeSetHandler('seekbackward', (details) => {
      const skipTime = details.seekOffset || 10
      // We will adjust via store
      const state = this.getCurrentState()
      if (state) {
        playerStore.seek(Math.max(0, state.currentTime - skipTime))
      }
    })

    safeSetHandler('seekforward', (details) => {
      const skipTime = details.seekOffset || 10
      const state = this.getCurrentState()
      if (state) {
        playerStore.seek(Math.min(state.duration, state.currentTime + skipTime))
      }
    })

    safeSetHandler('seekto', (details) => {
      if (details.seekTime !== undefined && details.seekTime !== null) {
        playerStore.seek(details.seekTime)
      }
    })

    safeSetHandler('stop', () => {
      playerStore.togglePlay()
    })
  }

  private getCurrentState() {
    let state: any = null
    const unsub = playerStore.subscribe((s) => {
      state = s
    })
    unsub()
    return state
  }

  private subscribeToPlayer() {
    playerStore.subscribe((state) => {
      if (typeof navigator === 'undefined' || !('mediaSession' in navigator)) return

      const ms = navigator.mediaSession

      // 1. Update Playback State
      ms.playbackState = state.isPlaying ? 'playing' : (state.currentTrack ? 'paused' : 'none')

      // 2. Update Track Metadata on Track Change
      if (state.currentTrack && state.currentTrack.id !== this.lastTrackId) {
        this.lastTrackId = state.currentTrack.id
        this.updateMetadata(state.currentTrack)
      } else if (!state.currentTrack && this.lastTrackId !== null) {
        this.lastTrackId = null
        ms.metadata = null
      }

      // 3. Update Position State if duration is valid and positive
      if (state.duration > 0 && typeof ms.setPositionState === 'function') {
        try {
          const pos = Math.max(0, Math.min(state.currentTime, state.duration))
          ms.setPositionState({
            duration: state.duration,
            playbackRate: 1.0,
            position: pos,
          })
        } catch (e) {
          // Ignored if out of range or precision issues
        }
      }
    })
  }

  private updateMetadata(track: SearchResultItem) {
    if (typeof MediaMetadata === 'undefined') return

    const artworkList: MediaImage[] = []
    if (track.thumbnailUrl) {
      artworkList.push({
        src: track.thumbnailUrl,
        sizes: '512x512',
        type: 'image/jpeg',
      })
    }

    navigator.mediaSession.metadata = new MediaMetadata({
      title: track.title,
      artist: track.uploaderName || 'Dengerin',
      album: 'Dengerin Desktop',
      artwork: artworkList,
    })
  }
}

export const mediaSessionService = new MediaSessionService()
