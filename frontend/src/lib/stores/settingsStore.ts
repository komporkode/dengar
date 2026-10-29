import { writable, get } from 'svelte/store'
import {
  WindowSetSize,
  WindowGetSize,
  WindowSetPosition,
  WindowGetPosition,
  WindowIsMaximised,
  WindowMaximise,
  WindowUnmaximise,
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
        if (parsed && typeof parsed.width === 'number' && parsed.width >= 600) {
          return parsed
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
  if (get(isMiniPlayer) || get(isTransitioning)) return
  isTransitioning.set(true)

  try {
    // 1. Capture current window state
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
      if (p) {
        posX = p.x
        posY = p.y
      }
    } catch (_) {}

    persistWindowState({
      wasMaximised: isMax,
      width: curW,
      height: curH,
      x: posX,
      y: posY
    })

    // 2. Unmaximise if maximized
    if (isMax) {
      try {
        WindowUnmaximise()
      } catch (_) {}
    }

    // 3. Wait for full desktop view to fade out smoothly
    await new Promise(r => setTimeout(r, 120))

    // 4. Set mini player constraints & window size
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

    // 6. Give OS window time to settle at mini size before mounting MiniPlayer
    await new Promise(r => setTimeout(r, 90))

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
  if (!get(isMiniPlayer) || get(isTransitioning)) return
  isTransitioning.set(true)

  try {
    // 1. Disable Always On Top immediately
    try {
      WindowSetAlwaysOnTop(false)
    } catch (_) {}

    // 2. Mini player fades out smoothly
    await new Promise(r => setTimeout(r, 120))

    // 3. Reset min size for desktop
    try {
      WindowSetMinSize(900, 600)
    } catch (_) {}

    // 4. Restore window size & position or maximised state
    try {
      if (savedWindowState.wasMaximised) {
        WindowMaximise()
      } else {
        WindowSetSize(savedWindowState.width, savedWindowState.height)
        if (typeof savedWindowState.x === 'number' && typeof savedWindowState.y === 'number') {
          WindowSetPosition(savedWindowState.x, savedWindowState.y)
        }
      }
    } catch (e) {
      console.warn('Failed to restore window size:', e)
      try {
        WindowSetSize(1200, 800)
      } catch (_) {}
    }

    // 5. Allow OS window time to expand to full size before showing desktop UI
    await new Promise(r => setTimeout(r, 140))

    // 6. Unmount mini player and show Desktop UI
    isMiniPlayer.set(false)
  } catch (err) {
    console.error('Failed to exit mini player:', err)
    isMiniPlayer.set(false)
  } finally {
    setTimeout(() => {
      isTransitioning.set(false)
    }, 120)
  }
}
