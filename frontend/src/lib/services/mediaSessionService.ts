import { playerStore } from '../stores/playerStore'
import { isDiscordRpcEnabled } from '../stores/settingsStore'
import { UpdateDiscordActivity, ClearDiscordActivity, UpdateTaskbarPlayState } from '../../../wailsjs/go/main/App'
import { get } from 'svelte/store'
import type { SearchResultItem } from '../types'

/**
 * Service to integrate Hardware Media Keys, Windows SMTC, and Discord Rich Presence
 */
class MediaSessionService {
  private isInitialized = false
  private lastTrackId: string | null = null
  private lastPlaying: boolean | null = null

  public init() {
    if (this.isInitialized) return

    this.setupActionHandlers()
    this.subscribeToPlayer()
    this.subscribeToDiscordToggle()
    this.isInitialized = true
  }

  private subscribeToDiscordToggle() {
    isDiscordRpcEnabled.subscribe((enabled) => {
      if (!enabled) {
        ClearDiscordActivity().catch(() => {})
      } else {
        const state = this.getCurrentState()
        if (state && state.currentTrack) {
          UpdateDiscordActivity(
            state.currentTrack.title,
            state.currentTrack.uploaderName || 'Dengerin',
            state.duration || state.currentTrack.duration || 0,
            Math.round(state.currentTime),
            state.isPlaying
          ).catch(() => {})
        }
      }
    })
  }

  private setupActionHandlers() {
    if (typeof navigator === 'undefined' || !('mediaSession' in navigator)) return
    const ms = navigator.mediaSession

    const safeSetHandler = (action: MediaSessionAction, handler: MediaSessionActionHandler) => {
      try {
        ms.setActionHandler(action, handler)
      } catch (err) {
        console.warn(`[MediaSession] Action handler '${action}' not supported:`, err)
      }
    }

    safeSetHandler('play', () => playerStore.togglePlay())
    safeSetHandler('pause', () => playerStore.togglePlay())
    safeSetHandler('nexttrack', () => playerStore.playNext())
    safeSetHandler('previoustrack', () => playerStore.playPrev())
    safeSetHandler('seekbackward', (details) => {
      const skipTime = details.seekOffset || 10
      const state = this.getCurrentState()
      if (state) playerStore.seek(Math.max(0, state.currentTime - skipTime))
    })
    safeSetHandler('seekforward', (details) => {
      const skipTime = details.seekOffset || 10
      const state = this.getCurrentState()
      if (state) playerStore.seek(Math.min(state.duration, state.currentTime + skipTime))
    })
    safeSetHandler('seekto', (details) => {
      if (details.seekTime !== undefined && details.seekTime !== null) {
        playerStore.seek(details.seekTime)
      }
    })
    safeSetHandler('stop', () => playerStore.togglePlay())
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
      const hasMediaSession = typeof navigator !== 'undefined' && 'mediaSession' in navigator

      // 1. MediaSession state sync
      if (hasMediaSession) {
        const ms = navigator.mediaSession
        ms.playbackState = state.isPlaying ? 'playing' : (state.currentTrack ? 'paused' : 'none')

        if (state.currentTrack && state.currentTrack.id !== this.lastTrackId) {
          this.updateMetadata(state.currentTrack)
        } else if (!state.currentTrack && this.lastTrackId !== null) {
          ms.metadata = null
        }

        if (state.duration > 0 && typeof ms.setPositionState === 'function') {
          try {
            const pos = Math.max(0, Math.min(state.currentTime, state.duration))
            ms.setPositionState({
              duration: state.duration,
              playbackRate: 1.0,
              position: pos,
            })
          } catch (_) {}
        }
      }

      // 2. Discord RPC sync on track or play state change
      const trackChanged = state.currentTrack?.id !== this.lastTrackId
      const playChanged = state.isPlaying !== this.lastPlaying

      if (trackChanged || playChanged) {
        this.lastTrackId = state.currentTrack?.id || null
        this.lastPlaying = state.isPlaying

        // Synchronize Windows Taskbar Thumbnail Toolbar button (Play / Pause toggle)
        UpdateTaskbarPlayState(state.isPlaying).catch(() => {})

        if (get(isDiscordRpcEnabled)) {
          if (state.currentTrack) {
            UpdateDiscordActivity(
              state.currentTrack.title,
              state.currentTrack.uploaderName || 'Dengerin',
              state.duration || state.currentTrack.duration || 0,
              Math.round(state.currentTime),
              state.isPlaying
            ).catch(() => {})
          } else {
            ClearDiscordActivity().catch(() => {})
          }
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
