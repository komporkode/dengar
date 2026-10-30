import { writable, get } from 'svelte/store'
import type { SearchResultItem, StreamResult, PlayMode } from '../types'
import { GetStream, RecordPlay, Search, PurgeTrackCache } from '../../../wailsjs/go/main/App'
import { userDataStore } from './userDataStore'
import { queueStore, isAutoplayEnabled } from './queueStore'
import { audioQualityStore, showToast } from './settingsStore'
import { equalizerStore } from './equalizerStore'
import { networkStore } from './networkStore'

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
  autoRetryCountdown: number | null
  autoRetryAction: 'retry' | 'skip' | null
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
    autoRetryCountdown: null,
    autoRetryAction: null,
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

    // Connect Web Audio Equalizer graph
    try {
      equalizerStore.connectAudio(audio)
    } catch (e) {
      console.warn('Equalizer connect audio warning:', e)
    }

    let lastTokenRefresh = 0
    audio.addEventListener('timeupdate', () => {
      if (audio) {
        onTimeUpdate(audio.currentTime, audio.duration)

        // For long tracks (>30 mins = 1800s), refresh stream token every 25 minutes
        const s = get({ subscribe })
        if (s.duration > 1800 && audio.currentTime - lastTokenRefresh > 1500) {
          lastTokenRefresh = audio.currentTime
          checkAndRefreshToken(s.currentTrack, audio.currentTime)
        }
      }
    })

    audio.addEventListener('durationchange', () => {
      if (audio && audio.duration && !isNaN(audio.duration) && audio.duration > 0) {
        update(st => ({ ...st, duration: audio!.duration }))
      }
    })

    audio.addEventListener('ended', () => {
      clearLoadWatchdog()
      onTrackEnded()
    })

    audio.addEventListener('error', () => {
      clearLoadWatchdog()
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
      startLoadWatchdog()
    })

    audio.addEventListener('playing', () => {
      clearTimeout(stallTimer)
      clearLoadWatchdog()
      update(st => ({ ...st, isPlaying: true, isLoading: false }))
    })

    audio.addEventListener('pause', () => {
      clearLoadWatchdog()
      update(st => ({ ...st, isPlaying: false }))
    })

    audio.addEventListener('canplay', () => {
      clearLoadWatchdog()
      update(st => ({ ...st, isLoading: false }))
    })
  }

  let loadWatchdogTimer: any = null
  function startLoadWatchdog() {
    clearTimeout(loadWatchdogTimer)
    loadWatchdogTimer = setTimeout(() => {
      const s = get({ subscribe })
      if (s.isLoading && !s.isPlaying && s.currentTrack) {
        console.warn('[StallWatchdog] Playback failed to start within 7s, triggering automatic instance rotation...')
        handleMediaError()
      }
    }, 7000)
  }

  function clearLoadWatchdog() {
    clearTimeout(loadWatchdogTimer)
  }

  function setAudioSourceAndPlay(streamUrl: string, startTime = 0, autoPlay = true) {
    if (!audio) return
    const proxied = equalizerStore.getProxiedStreamUrl(streamUrl)
    if (audio.src !== proxied) {
      audio.src = proxied
    }
    audio.currentTime = startTime
    applyVolumeToAudio()
    equalizerStore.resumeContext()
    if (autoPlay) {
      startLoadWatchdog()
      audio.play().catch(e => {
        console.warn('Audio play error, attempting recovery:', e)
        handleMediaError()
      })
    }
  }

  // Initialize singleton audio instance in browser environment
  initAudio()

  // Auto-recovery when device reconnects to internet
  networkStore.onReconnect(() => {
    const s = get({ subscribe })
    if (s.currentTrack && s.error && s.error.includes('Koneksi internet')) {
      console.log('[Auto-Recovery] Network restored! Resuming track:', s.currentTrack.title)
      retryWithDifferentSource()
    }
  })

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
    const isOnline = get(networkStore).isOnline

    // Replay check for the exact same track (repeat-one or re-clicking current song)
    if (!forceRestart && s.currentTrack?.id === track.id && s.streamResult) {
      if (audio) {
        audio.currentTime = 0
        audio.play().catch(e => {
          console.warn('Replay error:', e)
          handleMediaError()
        })
      }
      update(st => ({
        ...st,
        currentTime: 0,
        isLoading: true,
        isPlaying: false,
        error: null,
        mode: activeMode,
      }))
      return
    }

    cancelAutoRetry()
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
      autoRetryCountdown: null,
      autoRetryAction: null,
    }))

    // Record to history
    userDataStore.recordHistory(track)

    try {
      let stream: StreamResult
      const currentQuality = get(audioQualityStore) || 'high'
      try {
        stream = await GetStream(track.id, false, currentQuality)
      } catch (firstErr) {
        if (!isOnline) {
          update(st => ({
            ...st,
            isLoading: false,
            isPlaying: false,
            error: 'Lagu belum tersedia di cache offline. Sambungkan internet untuk memutar.',
          }))
          showToast('Lagu belum tersedia di cache offline. Sambungkan internet untuk memutar.')
          return
        }
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
        isLoading: true,
        isPlaying: false,
        error: null,
      }))

      if (audio) {
        setAudioSourceAndPlay(streamUrl, 0, true)
      }
    } catch (e: any) {
      if (requestId !== currentRequestId) return
      console.error('Failed to resolve stream:', e)
      handleMediaError()
    }
  }

  let autoRetryTimer: any = null
  let countdownInterval: any = null

  function cancelAutoRetry() {
    if (autoRetryTimer) {
      clearTimeout(autoRetryTimer)
      autoRetryTimer = null
    }
    if (countdownInterval) {
      clearInterval(countdownInterval)
      countdownInterval = null
    }
    update(st => ({
      ...st,
      autoRetryCountdown: null,
      autoRetryAction: null,
    }))
  }

  function scheduleAutoRetry(action: 'retry' | 'skip', seconds = 3) {
    cancelAutoRetry()
    const s = get({ subscribe })
    if (!s.currentTrack) return

    let remaining = seconds
    update(st => ({
      ...st,
      isLoading: true,
      isPlaying: false,
      autoRetryCountdown: remaining,
      autoRetryAction: action,
      error: action === 'retry'
        ? `Gagal memutar audio dari server saat ini (${st.retryCount + 1}/3)`
        : 'Semua server gagal memutar lagu ini',
    }))

    countdownInterval = setInterval(() => {
      remaining -= 1
      if (remaining > 0) {
        update(st => ({ ...st, autoRetryCountdown: remaining }))
      } else {
        if (countdownInterval) {
          clearInterval(countdownInterval)
          countdownInterval = null
        }
      }
    }, 1000)

    autoRetryTimer = setTimeout(async () => {
      cancelAutoRetry()
      const currentState = get({ subscribe })
      if (!currentState.currentTrack) return

      if (action === 'retry') {
        await executeRetry()
      } else if (action === 'skip') {
        showToast('Lagu tidak tersedia di server, memutar lagu berikutnya...')
        playNext()
      }
    }, seconds * 1000)
  }

  async function executeRetry() {
    const s = get({ subscribe })
    if (!s.currentTrack) return

    const nextRetry = s.retryCount + 1
    const savedTime = s.currentTime

    update(st => ({
      ...st,
      isLoading: true,
      isPlaying: false,
      retryCount: nextRetry,
      error: null,
    }))

    try {
      const currentQuality = get(audioQualityStore) || 'high'
      console.log(`[Auto-Recovery] Rotating instance for track ${s.currentTrack.id} (attempt ${nextRetry}/3)...`)
      const freshStream = await GetStream(s.currentTrack.id, true, currentQuality)
      const streamUrl = freshStream.audioUrl || freshStream.videoUrl
      if (!streamUrl) throw new Error('Tidak ada stream audio yang valid')

      update(st => ({
        ...st,
        streamResult: freshStream,
        duration: freshStream.duration || s.currentTrack?.duration || 0,
        isLoading: true,
        isPlaying: false,
        error: null,
      }))

      if (audio) {
        setAudioSourceAndPlay(streamUrl, savedTime, true)
      }
    } catch (err: any) {
      console.error('[Auto-Recovery] Retry failed:', err)
      if (nextRetry >= 3) {
        console.warn('[Auto-Recovery] Max retries reached, scheduling auto-skip...')
        scheduleAutoRetry('skip', 3)
      } else {
        scheduleAutoRetry('retry', 3)
      }
    }
  }

  async function handleMediaError() {
    initAudio()
    const s = get({ subscribe })

    // If an error occurred on the current track, immediately purge any corrupted disk cache
    if (s.currentTrack) {
      PurgeTrackCache(s.currentTrack.id).catch(() => {})
    }

    if (!get(networkStore).isOnline) {
      console.warn('[Player] Media error while offline, waiting for connection...')
      cancelAutoRetry()
      update(st => ({
        ...st,
        isLoading: false,
        isPlaying: false,
        error: 'Koneksi internet terputus. Pemutaran akan dilanjutkan otomatis saat kembali online.',
      }))
      return
    }

    if (!s.currentTrack) {
      cancelAutoRetry()
      return
    }

    if (s.retryCount >= 3) {
      console.warn(`[Auto-Recovery] Retries exhausted (${s.retryCount}/3) for track ${s.currentTrack.id}. Scheduling auto-skip...`)
      scheduleAutoRetry('skip', 3)
      return
    }

    console.warn(`[Auto-Recovery] Stream error on track ${s.currentTrack.id}, scheduling auto-retry (${s.retryCount + 1}/3) in 3s...`)
    scheduleAutoRetry('retry', 3)
  }

  async function retryWithDifferentSource() {
    cancelAutoRetry()
    const s = get({ subscribe })
    if (!s.currentTrack) return
    console.log(`[Manual-Retry] User requested new streaming source for: ${s.currentTrack.title}`)
    update(st => ({
      ...st,
      retryCount: 0,
      isLoading: true,
      isPlaying: false,
      error: null,
    }))
    try {
      const currentQuality = get(audioQualityStore) || 'high'
      const freshStream = await GetStream(s.currentTrack.id, true, currentQuality)
      const streamUrl = freshStream.audioUrl || freshStream.videoUrl
      if (!streamUrl) throw new Error('Tidak ada stream audio yang valid')

      update(st => ({
        ...st,
        streamResult: freshStream,
        duration: freshStream.duration || s.currentTrack?.duration || 0,
        isLoading: true,
        isPlaying: false,
        error: null,
      }))

      if (audio) {
        setAudioSourceAndPlay(streamUrl, s.currentTime || 0, true)
      }
    } catch (err: any) {
      console.error('Manual retry failed:', err)
      scheduleAutoRetry('skip', 3)
    }
  }

  function clearError() {
    cancelAutoRetry()
    update(st => ({ ...st, error: null }))
  }

  async function checkAndRefreshToken(track: SearchResultItem | null, currentTime: number) {
    if (!track) return
    try {
      console.log(`[TokenRenewal] Refreshing streaming token for long track (${Math.round(currentTime / 60)}m)...`)
      const currentQuality = get(audioQualityStore) || 'high'
      const fresh = await GetStream(track.id, false, currentQuality)
      const streamUrl = fresh.audioUrl || fresh.videoUrl
      if (streamUrl && audio && audio.src !== streamUrl) {
        const wasPaused = audio.paused
        audio.src = streamUrl
        audio.currentTime = currentTime
        if (!wasPaused) {
          audio.play().catch(console.warn)
        }
      }
    } catch (e) {
      console.warn('[TokenRenewal] Failed to refresh token in background:', e)
    }
  }

  function togglePlay() {
    initAudio()
    if (!audio) return

    if (audio.paused) {
      equalizerStore.resumeContext()
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
    cancelAutoRetry()
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
    cancelAutoRetry()
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
    retryWithDifferentSource,
    cancelAutoRetry,
    clearError,
    playNext,
    playPrev,
  }
}

export const playerStore = createPlayerStore()
