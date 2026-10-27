import { writable, get } from 'svelte/store'
import type { SearchResultItem, StreamResult, PlayMode } from '../types'
import { GetStream } from '../../../wailsjs/go/main/App'
import { userDataStore } from './userDataStore'
import { queueStore } from './queueStore'

interface PlayerState {
  currentTrack: SearchResultItem | null
  streamResult: StreamResult | null
  isPlaying: boolean
  currentTime: number
  duration: number
  volume: number
  isMuted: boolean
  mode: PlayMode
  isLoading: boolean
  error: string | null
  retryCount: number
}

function createPlayerStore() {
  const { subscribe, set, update } = writable<PlayerState>({
    currentTrack: null,
    streamResult: null,
    isPlaying: false,
    currentTime: 0,
    duration: 0,
    volume: 0.8,
    isMuted: false,
    mode: 'audio',
    isLoading: false,
    error: null,
    retryCount: 0,
  })

  let audioElement: HTMLAudioElement | null = null
  let videoElement: HTMLVideoElement | null = null
  let currentRequestId = 0

  function bindMedia(audio: HTMLAudioElement | null, video: HTMLVideoElement | null) {
    audioElement = audio
    videoElement = video
  }

  function getActiveMedia(): HTMLMediaElement | null {
    const s = get({ subscribe })
    return s.mode === 'video' ? videoElement : audioElement
  }

  async function playTrack(track: SearchResultItem, modeOverride?: PlayMode) {
    const requestId = ++currentRequestId
    const s = get({ subscribe })
    const activeMode = modeOverride || s.mode

    update(st => ({
      ...st,
      currentTrack: track,
      isLoading: true,
      error: null,
      mode: activeMode,
      retryCount: 0,
    }))

    // Record to history
    userDataStore.recordHistory(track)

    try {
      const stream = await GetStream(track.id, false)
      if (requestId !== currentRequestId) {
        // Obsolete request, ignore
        return
      }
      if (!stream.audioUrl && !stream.videoUrl) {
        throw new Error('Tidak ada stream audio/video yang valid')
      }

      update(st => ({
        ...st,
        streamResult: stream,
        duration: stream.duration || track.duration || 0,
        isLoading: false,
        isPlaying: true,
      }))
    } catch (e: any) {
      if (requestId !== currentRequestId) return
      console.error('Failed to resolve stream:', e)
      update(st => ({
        ...st,
        isLoading: false,
        error: e?.message || 'Gagal memutar stream media',
      }))
    }
  }

  async function handleMediaError() {
    const s = get({ subscribe })
    if (!s.currentTrack || s.retryCount >= 3) {
      update(st => ({ ...st, error: 'Stream error setelah beberapa percobaan', isLoading: false }))
      return
    }

    console.warn(`[Auto-Recovery] Stream error on track ${s.currentTrack.id}, retrying with fresh instance...`)
    const savedTime = s.currentTime

    update(st => ({
      ...st,
      isLoading: true,
      retryCount: st.retryCount + 1,
    }))

    try {
      const freshStream = await GetStream(s.currentTrack.id, true)
      update(st => ({
        ...st,
        streamResult: freshStream,
        isLoading: false,
        isPlaying: true,
      }))

      const media = getActiveMedia()
      if (media) {
        media.currentTime = savedTime
        media.play().catch(console.error)
      }
    } catch (err: any) {
      console.error('Auto re-resolve failed:', err)
      update(st => ({
        ...st,
        isLoading: false,
        error: 'Gagal memperbarui stream cadangan',
      }))
    }
  }

  function togglePlay() {
    const media = getActiveMedia()
    if (!media) return

    if (media.paused) {
      media.play().then(() => {
        update(st => ({ ...st, isPlaying: true }))
      }).catch(console.error)
    } else {
      media.pause()
      update(st => ({ ...st, isPlaying: false }))
    }
  }

  function seek(time: number) {
    const media = getActiveMedia()
    if (media) {
      media.currentTime = time
    }
    update(st => ({ ...st, currentTime: time }))
  }

  function setVolume(v: number) {
    const vol = Math.max(0, Math.min(1, v))
    const media = getActiveMedia()
    if (media) {
      media.volume = vol
    }
    update(st => ({ ...st, volume: vol, isMuted: vol === 0 }))
  }

  function toggleMute() {
    update(st => {
      const newMuted = !st.isMuted
      const media = getActiveMedia()
      if (media) {
        media.muted = newMuted
      }
      return { ...st, isMuted: newMuted }
    })
  }

  function toggleMode() {
    const s = get({ subscribe })
    const newMode: PlayMode = s.mode === 'audio' ? 'video' : 'audio'
    const currentTime = s.currentTime

    update(st => ({ ...st, mode: newMode }))

    // Resume time in newly active element
    setTimeout(() => {
      const media = getActiveMedia()
      if (media) {
        media.currentTime = currentTime
        if (s.isPlaying) {
          media.play().catch(console.error)
        }
      }
    }, 50)
  }

  function onTimeUpdate(currentTime: number, duration: number) {
    update(st => ({
      ...st,
      currentTime,
      duration: duration || st.duration,
    }))
  }

  function onTrackEnded() {
    const nextTrack = queueStore.next()
    if (nextTrack) {
      playTrack(nextTrack)
    } else {
      update(st => ({ ...st, isPlaying: false, currentTime: 0 }))
    }
  }

  function playNext() {
    const nextTrack = queueStore.next()
    if (nextTrack) {
      playTrack(nextTrack)
    }
  }

  function playPrev() {
    const prevTrack = queueStore.prev()
    if (prevTrack) {
      playTrack(prevTrack)
    }
  }

  return {
    subscribe,
    bindMedia,
    playTrack,
    togglePlay,
    seek,
    setVolume,
    toggleMute,
    toggleMode,
    onTimeUpdate,
    onTrackEnded,
    handleMediaError,
    playNext,
    playPrev,
  }
}

export const playerStore = createPlayerStore()
