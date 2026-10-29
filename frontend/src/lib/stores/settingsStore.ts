import { writable, get } from 'svelte/store'
import {
  WindowSetSize,
  WindowGetSize,
  WindowIsMaximised,
  WindowMaximise,
  WindowUnmaximise,
  WindowCenter,
  WindowShow,
  WindowSetAlwaysOnTop,
  WindowSetMinSize
} from '../../../wailsjs/runtime/runtime.js'

export type AudioQuality = 'high' | 'saver'

const savedQuality = (typeof localStorage !== 'undefined' ? localStorage.getItem('dengar_quality') : null) as AudioQuality | null

export const audioQualityStore = writable<AudioQuality>(savedQuality === 'saver' ? 'saver' : 'high')

audioQualityStore.subscribe((val) => {
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('dengar_quality', val)
  }
})

export const isSettingsOpen = writable<boolean>(false)
export const isMiniPlayer = writable<boolean>(false)
export const isTransitioning = writable<boolean>(false)
export const toastMessage = writable<string | null>(null)

// Mini Player Pin (Always on top) store
const savedMiniPlayerPinned = typeof localStorage !== 'undefined' ? localStorage.getItem('dengar_mini_pinned') !== 'false' : true
export const isMiniPlayerPinned = writable<boolean>(savedMiniPlayerPinned)

isMiniPlayerPinned.subscribe((pinned) => {
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('dengar_mini_pinned', String(pinned))
  }
})

export function togglePinMiniPlayer(): boolean {
  const next = !get(isMiniPlayerPinned)
  isMiniPlayerPinned.set(next)
  try {
    WindowSetAlwaysOnTop(next)
  } catch (e) {
    console.warn('Wails WindowSetAlwaysOnTop error:', e)
  }
  showToast(next ? 'Disematkan di atas (Always on Top)' : 'Sematan dilepas')
  return next
}

const savedSidebarCollapsed = typeof localStorage !== 'undefined' ? localStorage.getItem('dengar_sidebar_collapsed') === 'true' : false
export const isSidebarCollapsed = writable<boolean>(savedSidebarCollapsed)

isSidebarCollapsed.subscribe((collapsed) => {
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('dengar_sidebar_collapsed', String(collapsed))
  }
})

export function toggleSidebar() {
  isSidebarCollapsed.update(v => !v)
}

let toastTimer: any

export function showToast(message: string) {
  toastMessage.set(message)
  clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toastMessage.set(null)
  }, 2500)
}

interface SavedWindowState {
  wasMaximised: boolean
  width: number
  height: number
}

const defaultWindowState: SavedWindowState = {
  wasMaximised: false,
  width: 1200,
  height: 800
}

function loadSavedWindowState(): SavedWindowState {
  try {
    if (typeof localStorage !== 'undefined') {
      const raw = localStorage.getItem('dengar_window_state')
      if (raw) {
        const parsed = JSON.parse(raw)
        if (parsed && typeof parsed.width === 'number' && parsed.width >= 600 && parsed.width <= 3840) {
          return {
            wasMaximised: !!parsed.wasMaximised,
            width: parsed.width,
            height: Math.min(Math.max(parsed.height || 800, 600), 2160)
          }
        }
      }
    }
  } catch (e) {
    // fallback
  }
  return defaultWindowState
}

let savedWindowState: SavedWindowState = loadSavedWindowState()

function persistWindowState(state: SavedWindowState) {
  savedWindowState = state
  try {
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('dengar_window_state', JSON.stringify(state))
    }
  } catch (e) {
    // ignore
  }
}

export async function enterMiniPlayer() {
  if (get(isMiniPlayer)) return

  try {
    // 1. Capture current window state safely
    let isMax = false
    let curW = 1200
    let curH = 800

    try {
      isMax = await WindowIsMaximised()
    } catch (_) {}

    try {
      const s = await WindowGetSize()
      if (s && s.w >= 600 && s.h >= 400) {
        curW = s.w
        curH = s.h
      }
    } catch (_) {}

    savedWindowState = {
      wasMaximised: isMax,
      width: curW,
      height: curH
    }
    persistWindowState(savedWindowState)

    // 2. If maximised, unmaximise first before shrinking
    if (isMax) {
      try {
        WindowUnmaximise()
      } catch (_) {}
    }

    // 3. Set minimum bounds & target mini size
    try {
      WindowSetMinSize(360, 140)
      WindowSetSize(420, 175)
    } catch (e) {
      console.warn('Wails WindowSetSize error:', e)
    }

    // 4. Pin / Always on top
    const pinned = get(isMiniPlayerPinned)
    try {
      WindowSetAlwaysOnTop(pinned)
    } catch (_) {}

    try {
      WindowShow()
    } catch (_) {}

    isMiniPlayer.set(true)
  } catch (err) {
    console.error('Failed to enter mini player:', err)
    isMiniPlayer.set(true)
  } finally {
    isTransitioning.set(false)
  }
}

export async function exitMiniPlayer() {
  if (!get(isMiniPlayer)) return

  try {
    // 1. Immediately disable Always On Top
    try {
      WindowSetAlwaysOnTop(false)
    } catch (_) {}

    // 2. Reset min size for desktop
    try {
      WindowSetMinSize(900, 600)
    } catch (_) {}

    // 3. Restore window size and position SAFELY (Always visible and centered on screen!)
    try {
      if (savedWindowState.wasMaximised) {
        WindowMaximise()
      } else {
        const targetW = Math.min(Math.max(savedWindowState.width || 1200, 900), 2560)
        const targetH = Math.min(Math.max(savedWindowState.height || 800, 600), 1440)
        WindowSetSize(targetW, targetH)
        WindowCenter() // GUARANTEES window is dead-center on the monitor!
      }
      WindowShow() // GUARANTEES window is displayed and brought to foreground!
    } catch (e) {
      console.warn('Failed to restore window size:', e)
      try {
        WindowSetSize(1200, 800)
        WindowCenter()
        WindowShow()
      } catch (_) {}
    }

    isMiniPlayer.set(false)
  } catch (err) {
    console.error('Failed to exit mini player:', err)
    isMiniPlayer.set(false)
  } finally {
    isTransitioning.set(false)
  }
}
