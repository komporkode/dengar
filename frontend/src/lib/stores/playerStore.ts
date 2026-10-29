import { writable, get } from 'svelte/store'
import type { SearchResultItem, StreamResult, PlayMode } from '../types'
import { GetStream, RecordPlay, Search } from '../../../wailsjs/go/main/App'
import { userDataStore } from './userDataStore'
import { queueStore, isAutoplayEnabled } from './queueStore'
import { audioQualityStore } from './settingsStore'

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

  let audio: HTMLAudioElement | null = null
  let currentRequestId = 0
  let hasRecordedPlay = false
  let recordedTrackId = ''

  function initAudio() {
    if (typeof window === 'undefined' || typeof Audio === 'undefined') return
    if (audio) return

    audio = new Audio()
    audio.preload = 'auto'
    audio.volume = initialVolume
    audio.muted = initialMuted

    audio.addEventListener('timeupdate', () => {
      if (audio) {
        onTimeUpdate(audio.currentTime, audio.duration)
      }
    })

    audio.addEventListener('durationchange', () => {
      if (audio && audio.duration && !isNaN(audio.duration) && audio.duration > 0) {
        update(st => ({ ...st, duration: audio!.duration }))
      }
    })

    audio.addEventListener('ended', () => {
      onTrackEnded()
    })

    audio.addEventListener('error', () => {
      handleMediaError()
    })

    let stallTimer: any = null
    audio.addEventListener('stalled', () => {
      clearTimeout(stallTimer)
      stallTimer = setTimeout(() => {
        const s = get({ subscribe })
        if (s.isPlaying && audio && audio.paused === false) {
          console.warn('[Auto-Recovery] Audio playback stalled for > 6s, refreshing stream...')
          handleMediaError()
        }
      }, 6000)
    })

    audio.addEventListener('waiting', () => {
      update(st => ({ ...st, isLoading: true }))
    })

    audio.addEventListener('playing', () => {
      clearTimeout(stallTimer)
      update(st => ({ ...st, isPlaying: true, isLoading: false }))
    })

    audio.addEventListener('pause', () => {
      update(st => ({ ...st, isPlaying: false }))
    })

    audio.addEventListener('canplay', () => {
      update(st => ({ ...st, isLoading: false }))
    })
  }

  // Initialize singleton audio instance in browser environment
  initAudio()

  function applyVolumeToAudio() {
    if (!audio) return
    const s = get({ subscribe })
    audio.volume = s.volume
    audio.muted = s.isMuted
  }

  function bindAudio(_audio: HTMLAudioElement | null) {
    // Deprecated: Audio lifecycle is managed directly by playerStore singleton
  }

  function bindVideo(_video: HTMLVideoElement | null) {
    // Pure audio mode - video no-op for backward compatibility
  }

  function bindMedia(_audio: HTMLAudioElement | null, _video: HTMLVideoElement | null) {
    // Deprecated: Audio lifecycle is managed directly by playerStore singleton
  }

  function getActiveMedia(): HTMLMediaElement | null {
    return audio
  }

  async function playTrack(track: SearchResultItem, modeOverride?: PlayMode, forceRestart = false) {
    initAudio()
    const s = get({ subscribe })
    const activeMode = modeOverride || s.mode

    // Replay check for the exact same track (repeat-one or re-clicking current song)
    if (!forceRestart && s.currentTrack?.id === track.id && s.streamResult) {
      if (audio) {
        audio.currentTime = 0
        audio.play().catch(console.error)
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
      let stream: StreamResult
      const currentQuality = get(audioQualityStore) || 'high'
      try {
        stream = await GetStream(track.id, false, currentQuality)
      } catch (firstErr) {
        console.warn('Initial stream resolution failed, attempting with rotated instance...', firstErr)
        stream = await GetStream(track.id, true, currentQuality)
      }

      if (requestId !== currentRequestId) {
        // Obsolete request, ignore
        return
      }
      const streamUrl = stream.audioUrl || stream.videoUrl
      if (!streamUrl) {
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

      if (audio) {
        if (audio.src !== streamUrl) {
          audio.src = streamUrl
        }
        audio.currentTime = 0
        audio.volume = s.volume
        audio.muted = s.isMuted
        audio.play().catch(e => {
          console.warn('Audio play error:', e)
        })
      }
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
    initAudio()
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
      const currentQuality = get(audioQualityStore) || 'high'
      const freshStream = await GetStream(s.currentTrack.id, true, currentQuality)
      update(st => ({
        ...st,
        streamResult: freshStream,
        isLoading: false,
        isPlaying: true,
      }))

      if (audio) {
        const streamUrl = freshStream.audioUrl || freshStream.videoUrl
        if (streamUrl) {
          audio.src = streamUrl
          audio.currentTime = savedTime
          audio.play().catch(console.error)
        }
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
    initAudio()
    if (!audio) return

    if (audio.paused) {
      audio.play().then(() => {
        update(st => ({ ...st, isPlaying: true }))
      }).catch(console.error)
    } else {
      audio.pause()
      update(st => ({ ...st, isPlaying: false }))
    }
  }

  function seek(time: number) {
    initAudio()
    if (audio) {
      try {
        audio.currentTime = time
      } catch (e) {
        console.warn('Seek error:', e)
      }
    }
    update(st => ({ ...st, currentTime: time }))
  }

  function setVolume(v: number) {
    initAudio()
    const vol = Math.max(0, Math.min(1, Math.round(v * 100) / 100))
    const isMuted = vol === 0
    if (vol > 0) {
      lastNonZeroVolume = vol
    }
    update(st => {
      saveVolumeState(vol, isMuted)
      return { ...st, volume: vol, isMuted }
    })
    if (audio) {
      audio.volume = vol
      audio.muted = isMuted
    }
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
    initAudio()
    update(st => {
      const nextMuted = !st.isMuted
      let nextVol = st.volume
      if (!nextMuted && nextVol === 0) {
        nextVol = lastNonZeroVolume > 0 ? lastNonZeroVolume : 0.5
      }
      saveVolumeState(nextVol, nextMuted)
      if (audio) {
        audio.muted = nextMuted
        audio.volume = nextVol
      }
      return { ...st, isMuted: nextMuted, volume: nextVol }
    })
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
      if (get(isAutoplayEnabled) && s.currentTrack) {
        triggerAutoplayNext(s.currentTrack)
      } else {
        update(st => ({ ...st, isPlaying: false, currentTime: 0 }))
      }
    }
  }

  async function triggerAutoplayNext(current: SearchResultItem) {
    update(st => ({ ...st, isLoading: true }))
    try {
      // YouTube-style dynamic Up Next recommendations based on artist / title
      const query = current.uploaderName ? `${current.uploaderName} music` : current.title
      const res = await Search(query, 'music_songs')
      if (res && res.items && res.items.length > 0) {
        // Exclude current track from newly queued items
        const candidates = res.items.filter(it => it.id !== current.id)
        if (candidates.length > 0) {
          // Push top 6 related tracks into the queue
          for (let i = 0; i < Math.min(6, candidates.length); i++) {
            queueStore.addToQueue(candidates[i])
          }
          const nextTrack = queueStore.next(false)
          if (nextTrack) {
            playTrack(nextTrack, undefined, true)
            return
          }
        }
      }
    } catch (e) {
      console.warn('Autoplay failed to fetch next recommendations:', e)
    }
    update(st => ({ ...st, isPlaying: false, currentTime: 0, isLoading: false }))
  }

  function playNext() {
    // When manually clicking next, advance forward even if repeat one is set
    const nextTrack = queueStore.next(true)
    if (nextTrack) {
      playTrack(nextTrack, undefined, true)
    } else {
      const s = get({ subscribe })
      if (get(isAutoplayEnabled) && s.currentTrack) {
        triggerAutoplayNext(s.currentTrack)
      }
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

  // Hot-swap stream when user switches audio quality in settings
  let activeQuality = get(audioQualityStore)
  audioQualityStore.subscribe(async (newQuality) => {
    if (newQuality === activeQuality) return
    activeQuality = newQuality

    const s = get({ subscribe })
    if (s.currentTrack && audio && s.streamResult) {
      console.log(`[AudioQuality] Hot-swapping stream quality to "${newQuality}" for track ${s.currentTrack.id}...`)
      const savedTime = audio.currentTime
      const wasPlaying = !audio.paused
      try {
        const freshStream = await GetStream(s.currentTrack.id, true, newQuality)
        const streamUrl = freshStream.audioUrl || freshStream.videoUrl
        if (streamUrl && audio) {
          audio.src = streamUrl
          audio.currentTime = savedTime
          if (wasPlaying) {
            audio.play().catch(e => console.warn('Hot-swap playback resume error:', e))
          }
          update(st => ({
            ...st,
            streamResult: freshStream,
          }))
        }
      } catch (e) {
        console.warn('[AudioQuality] Failed to hot-swap quality on active track:', e)
      }
    }
  })

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
