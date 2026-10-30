import { writable, get } from 'svelte/store'

export interface EqualizerBand {
  frequency: number
  type: BiquadFilterType
  label: string
  q?: number
}

export const EQ_BANDS: EqualizerBand[] = [
  { frequency: 60, type: 'lowshelf', label: '60 Hz' },
  { frequency: 230, type: 'peaking', label: '230 Hz', q: 1.0 },
  { frequency: 910, type: 'peaking', label: '910 Hz', q: 1.0 },
  { frequency: 3600, type: 'peaking', label: '3.6 kHz', q: 1.0 },
  { frequency: 14000, type: 'highshelf', label: '14 kHz' },
]

export const EQ_PRESETS: Record<string, number[]> = {
  'Flat': [0, 0, 0, 0, 0],
  'Bass Boost': [6, 4, 1, 0, 0],
  'Bass Reducer': [-6, -4, 0, 0, 0],
  'Vocal Booster': [-2, 2, 5, 3, -1],
  'Treble Boost': [0, 0, 1, 4, 6],
  'Rock': [4, 2, -1, 3, 4],
  'Pop': [-1, 2, 4, 2, -1],
  'Electronic': [5, 3, 0, 2, 4],
}

export interface EqualizerState {
  isEnabled: boolean
  preset: string
  gains: number[]
  isSupported: boolean
}

const STORAGE_KEY = 'dengerin_eq_settings'

function loadSavedSettings(): { isEnabled: boolean; preset: string; gains: number[] } {
  if (typeof window === 'undefined') {
    return { isEnabled: false, preset: 'Flat', gains: [0, 0, 0, 0, 0] }
  }
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) {
      const parsed = JSON.parse(raw)
      return {
        isEnabled: Boolean(parsed.isEnabled),
        preset: typeof parsed.preset === 'string' ? parsed.preset : 'Flat',
        gains: Array.isArray(parsed.gains) && parsed.gains.length === 5 ? parsed.gains : [0, 0, 0, 0, 0],
      }
    }
  } catch (e) {
    console.warn('Failed to load equalizer settings from localStorage:', e)
  }
  return { isEnabled: false, preset: 'Flat', gains: [0, 0, 0, 0, 0] }
}

function saveSettings(state: EqualizerState) {
  if (typeof window === 'undefined') return
  try {
    localStorage.setItem(
      STORAGE_KEY,
      JSON.stringify({
        isEnabled: state.isEnabled,
        preset: state.preset,
        gains: state.gains,
      })
    )
  } catch (e) {
    console.warn('Failed to save equalizer settings:', e)
  }
}

function createEqualizerStore() {
  const initial = loadSavedSettings()
  const { subscribe, set, update } = writable<EqualizerState>({
    isEnabled: initial.isEnabled,
    preset: initial.preset,
    gains: initial.gains,
    isSupported: true,
  })

  let audioCtx: AudioContext | null = null
  let sourceNode: MediaElementAudioSourceNode | null = null
  let filterNodes: BiquadFilterNode[] = []
  let isConnected = false

  function applyFiltersToNodes() {
    if (!filterNodes.length) return
    const s = get({ subscribe })
    filterNodes.forEach((node, idx) => {
      // If disabled, bypass by setting gain to 0
      const targetGain = s.isEnabled ? (s.gains[idx] ?? 0) : 0
      if (node.gain) {
        node.gain.setTargetAtTime(targetGain, audioCtx ? audioCtx.currentTime : 0, 0.05)
      }
    })
  }

  function connectAudio(audioElement: HTMLAudioElement) {
    if (typeof window === 'undefined') return
    if (isConnected || !audioElement) return

    try {
      const AudioContextClass = window.AudioContext || (window as any).webkitAudioContext
      if (!AudioContextClass) {
        update(st => ({ ...st, isSupported: false }))
        return
      }

      if (!audioCtx) {
        audioCtx = new AudioContextClass()
      }

      if (!sourceNode) {
        audioElement.crossOrigin = 'anonymous'
        sourceNode = audioCtx.createMediaElementSource(audioElement)
      }

      filterNodes = EQ_BANDS.map(band => {
        const filter = audioCtx!.createBiquadFilter()
        filter.type = band.type
        filter.frequency.value = band.frequency
        if (band.q) filter.Q.value = band.q
        filter.gain.value = 0
        return filter
      })

      // Chain: source -> filter0 -> ... -> filter4 -> destination
      let current: AudioNode = sourceNode
      for (const filter of filterNodes) {
        current.connect(filter)
        current = filter
      }
      current.connect(audioCtx.destination)

      isConnected = true
      applyFiltersToNodes()
    } catch (e) {
      console.warn('Web Audio API Equalizer initialization warning:', e)
      // Normal fallback playback continues unaffected
    }
  }

  function resumeContext() {
    if (audioCtx && audioCtx.state === 'suspended') {
      audioCtx.resume().catch(e => console.warn('AudioContext resume error:', e))
    }
  }

  function setEnabled(enabled: boolean) {
    resumeContext()
    update(st => {
      const next = { ...st, isEnabled: enabled }
      saveSettings(next)
      return next
    })
    applyFiltersToNodes()
  }

  function setPreset(name: string) {
    resumeContext()
    const gains = EQ_PRESETS[name]
    if (!gains) return
    update(st => {
      const next = { ...st, preset: name, gains: [...gains], isEnabled: true }
      saveSettings(next)
      return next
    })
    applyFiltersToNodes()
  }

  function setBandGain(index: number, gain: number) {
    resumeContext()
    update(st => {
      const newGains = [...st.gains]
      newGains[index] = Math.max(-12, Math.min(12, Math.round(gain * 10) / 10))
      const next = { ...st, gains: newGains, preset: 'Kustom', isEnabled: true }
      saveSettings(next)
      return next
    })
    applyFiltersToNodes()
  }

  function resetToFlat() {
    setPreset('Flat')
  }

  function getProxiedStreamUrl(rawUrl: string, videoId?: string): string {
    if (!rawUrl) return ''
    // If it's already a relative / local API URL, don't double proxy
    if (rawUrl.startsWith('/api/') || rawUrl.startsWith('http://localhost') || rawUrl.startsWith('http://127.0.0.1')) {
      return rawUrl
    }
    let proxied = `/api/audio-proxy?url=${encodeURIComponent(rawUrl)}`
    if (videoId) {
      proxied += `&video_id=${encodeURIComponent(videoId)}`
    }
    return proxied
  }

  return {
    subscribe,
    connectAudio,
    resumeContext,
    setEnabled,
    setPreset,
    setBandGain,
    resetToFlat,
    getProxiedStreamUrl,
  }
}

export const equalizerStore = createEqualizerStore()

