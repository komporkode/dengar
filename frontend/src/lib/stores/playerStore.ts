import { writable, get } from 'svelte/store'
import type { SearchResultItem, StreamResult, PlayMode } from '../types'
import { GetStream, RecordPlay } from '../../../wailsjs/go/main/App'
import { userDataStore } from './userDataStore'
import { queueStore } from './queueStore'

const STORAGE_VOLUME_KEY = 'dengar_player_volume'
const STORAGE_MUTED_KEY = 'dengar_player_muted'

function loadSavedVolume(): number {
  if (typeof localStorage === 'undefined') return 0.8
  try {
    const raw = localStorage.getItem(STORAGE_VOLUME_KEY)
    if (raw !== null) {
      const parsed = parseFloat(raw)
      if (!isNaN(parsed) && parsed >= 0 && parsed <= 1) {
        return parsed
      }
    }
  } catch (e) {
    console.warn('Failed to load volume from localStorage:', e)
  }
  return 0.8
}

function loadSavedMuted(): boolean {
  if (typeof localStorage === 'undefined') return false
  try {
    return localStorage.getItem(STORAGE_MUTED_KEY) === 'true'
  } catch (e) {
    return false
  }
}

function saveVolumeState(volume: number, isMuted: boolean) {
  if (typeof localStorage === 'undefined') return
  try {
    localStorage.setItem(STORAGE_VOLUME_KEY, volume.toString())
    localStorage.setItem(STORAGE_MUTED_KEY, isMuted.toString())
  } catch (e) {
    console.warn('Failed to save volume to localStorage:', e)
  }
}

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
  const initialVolume = loadSavedVolume()
  const initialMuted = loadSavedMuted()
  let lastNonZeroVolume = initialVolume > 0 ? initialVolume : 0.8

  const { subscribe, set, update } = writable<PlayerState>({
    currentTrack: null,
    streamResult: null,
    isPlaying: false,
    currentTime: 0,
    duration: 0,
    volume: initialVolume,
    isMuted: initialMuted,
    mode: 'audio',
    isLoading: false,
    error: null,
    retryCount: 0,
  })

  let audioElement: HTMLAudioElement | null = null
  let currentRequestId = 0
  let hasRecordedPlay = false
  let recordedTrackId = ''

  function applyVolumeToAudio() {
    if (!audioElement) return
    const s = get({ subscribe })
    audioElement.volume = s.volume
    audioElement.muted = s.isMuted
  }

  function bindAudio(audio: HTMLAudioElement | null) {
    audioElement = audio
    applyVolumeToAudio()
  }

  function bindVideo(_video: HTMLVideoElement | null) {
    // Pure audio mode - video no-op for backward compatibility
  }

  function bindMedia(audio: HTMLAudioElement | null, _video: HTMLVideoElement | null) {
    if (audio !== null) {
      audioElement = audio
      applyVolumeToAudio()
    }
  }

  function getActiveMedia(): HTMLMediaElement | null {
    return audioElement
  }

  async function playTrack(track: SearchResultItem, modeOverride?: PlayMode, forceRestart = false) {
    const s = get({ subscribe })
    const activeMode = modeOverride || s.mode

    // Replay check for the exact same track (repeat-one or re-clicking current song)
    if (s.currentTrack?.id === track.id && s.streamResult) {
      const media = getActiveMedia()
      if (media) {
        media.currentTime = 0
        media.play().catch(console.error)
      }
      update(st => ({
        ...st,
        isPlaying: true,
        currentTime: 0,
        isLoading: false,
        error: null,
        mode: activeMode,
      }))
      return
    }

    const requestId = ++currentRequestId
    hasRecordedPlay = false
    recordedTrackId = track.id
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
        currentTime: 0,
        isLoading: false,
        isPlaying: true,
      }))

      // Ensure active media starts playing
      setTimeout(() => {
        const media = getActiveMedia()
        if (media) {
          media.currentTime = 0
          media.play().catch(console.error)
        }
      }, 20)
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
    const vol = Math.max(0, Math.min(1, Math.round(v * 100) / 100))
    const isMuted = vol === 0
    if (vol > 0) {
      lastNonZeroVolume = vol
    }
    update(st => {
      saveVolumeState(vol, isMuted)
      return { ...st, volume: vol, isMuted }
    })
    applyVolumeToAudio()
  }

  function adjustVolume(delta: number) {
    const s = get({ subscribe })
    if (s.isMuted && delta > 0) {
      setVolume(Math.min(1, (lastNonZeroVolume || 0.5) + delta))
      return
    }
    const target = Math.max(0, Math.min(1, Math.round((s.volume + delta) * 100) / 100))
    setVolume(target)
  }

  function toggleMute() {
    update(st => {
      const nextMuted = !st.isMuted
      let nextVol = st.volume
      if (!nextMuted && nextVol === 0) {
        nextVol = lastNonZeroVolume > 0 ? lastNonZeroVolume : 0.5
      }
      saveVolumeState(nextVol, nextMuted)
      return { ...st, isMuted: nextMuted, volume: nextVol }
    })
    applyVolumeToAudio()
  }

  function toggleMode() {
    // Pure audio mode - kept for backward compatibility
  }

  function onTimeUpdate(currentTime: number, duration: number) {
    if (!hasRecordedPlay && currentTime >= 20) {
      const s = get({ subscribe })
      if (s.currentTrack && s.currentTrack.id === recordedTrackId) {
        hasRecordedPlay = true
        RecordPlay(s.currentTrack).catch(console.error)
      }
    }

    update(st => ({
      ...st,
      currentTime,
      duration: duration || st.duration,
    }))
  }

  function onTrackEnded() {
    const s = get({ subscribe })
    const q = get(queueStore)

    if (!hasRecordedPlay && s.currentTrack && s.currentTrack.id === recordedTrackId) {
      hasRecordedPlay = true
      RecordPlay(s.currentTrack).catch(console.error)
    }

    // Repeat one mode: loop same track immediately
    if (q.repeatMode === 'one' && s.currentTrack) {
      const media = getActiveMedia()
      if (media) {
        media.currentTime = 0
        media.play().catch(console.error)
      }
      update(st => ({ ...st, isPlaying: true, currentTime: 0 }))
      return
    }

    const nextTrack = queueStore.next(false)
    if (nextTrack) {
      playTrack(nextTrack, undefined, true)
    } else {
      update(st => ({ ...st, isPlaying: false, currentTime: 0 }))
    }
  }

  function playNext() {
    // When manually clicking next, advance forward even if repeat one is set
    const nextTrack = queueStore.next(true)
    if (nextTrack) {
      playTrack(nextTrack, undefined, true)
    }
  }

  function playPrev() {
    const prevTrack = queueStore.prev()
    if (prevTrack) {
      playTrack(prevTrack, undefined, true)
    } else {
      seek(0)
    }
  }

  return {
    subscribe,
    bindMedia,
    bindAudio,
    bindVideo,
    playTrack,
    togglePlay,
    seek,
    setVolume,
    adjustVolume,
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
