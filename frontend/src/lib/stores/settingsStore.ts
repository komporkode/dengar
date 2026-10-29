import { writable, get } from 'svelte/store'
import {
  WindowSetSize,
  WindowGetSize,
  WindowSetPosition,
  WindowGetPosition,
  WindowIsMaximised,
  WindowMaximise,
  WindowUnmaximise,
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
  x?: number
  y?: number
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
            height: Math.min(Math.max(parsed.height || 800, 600), 2160),
            x: parsed.x,
            y: parsed.y
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
  isTransitioning.set(true)

  try {
    // 1. Capture current window state safely
    let isMax = false
    let curW = 1200
    let curH = 800
    let posX: number | undefined
    let posY: number | undefined

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

    try {
      const p = await WindowGetPosition()
      if (p && p.x >= 0 && p.y >= 0) {
        posX = p.x
        posY = p.y
      }
    } catch (_) {}

    savedWindowState = {
      wasMaximised: isMax,
      width: curW,
      height: curH,
      x: posX,
      y: posY
    }
    persistWindowState(savedWindowState)

    // 2. If maximised, unmaximise first before shrinking
    if (isMax) {
      try {
        WindowUnmaximise()
      } catch (_) {}
    }

    // 3. Smooth fade out of full player (80ms)
    await new Promise(r => setTimeout(r, 80))

    // 4. Set minimum bounds & target mini size
    try {
      WindowSetMinSize(360, 140)
      WindowSetSize(420, 175)
    } catch (e) {
      console.warn('Wails WindowSetSize error:', e)
    }

    // 5. Pin / Always on top
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
    setTimeout(() => {
      isTransitioning.set(false)
    }, 80)
  }
}

export async function exitMiniPlayer() {
  if (!get(isMiniPlayer)) return
  isTransitioning.set(true)

  try {
    // 1. Immediately disable Always On Top
    try {
      WindowSetAlwaysOnTop(false)
    } catch (_) {}

    // 2. Allow MiniPlayer to fade out smoothly (100ms)
    await new Promise(r => setTimeout(r, 100))

    // 3. Reset min size for desktop
    try {
      WindowSetMinSize(900, 600)
    } catch (_) {}

    // 4. Restore window size and position smoothly near current mini position
    try {
      if (savedWindowState.wasMaximised) {
        WindowMaximise()
      } else {
        const targetW = Math.min(Math.max(savedWindowState.width || 1200, 900), 2560)
        const targetH = Math.min(Math.max(savedWindowState.height || 800, 600), 1440)

        // Find position near current mini player so window doesn't teleport
        const currentMiniPos = await WindowGetPosition().catch(() => null)
        let targetX = savedWindowState.x ?? 100
        let targetY = savedWindowState.y ?? 100

        if (currentMiniPos && currentMiniPos.x >= 0 && currentMiniPos.y >= 0) {
          targetX = Math.max(30, currentMiniPos.x - Math.floor((targetW - 420) * 0.75))
          targetY = Math.max(30, currentMiniPos.y - Math.floor((targetH - 175) * 0.75))
        }

        WindowSetSize(targetW, targetH)
        if (targetX >= 0 && targetY >= 0) {
          WindowSetPosition(targetX, targetY)
        }
      }
      WindowShow()
    } catch (e) {
      console.warn('Failed to restore window size:', e)
      try {
        WindowSetSize(1200, 800)
        WindowShow()
      } catch (_) {}
    }

    // Small delay for OS window resize to settle before unmounting mini player
    await new Promise(r => setTimeout(r, 80))
    isMiniPlayer.set(false)
  } catch (err) {
    console.error('Failed to exit mini player:', err)
    isMiniPlayer.set(false)
  } finally {
    setTimeout(() => {
      isTransitioning.set(false)
    }, 100)
  }
}
