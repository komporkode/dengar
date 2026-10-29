import { writable, get } from 'svelte/store'
import { WindowSetSize, WindowSetAlwaysOnTop, WindowSetMinSize } from '../../../wailsjs/runtime/runtime.js'

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

export function enterMiniPlayer() {
  try {
    WindowSetMinSize(380, 160)
    WindowSetSize(420, 175)
    const pinned = get(isMiniPlayerPinned)
    WindowSetAlwaysOnTop(pinned)
    isMiniPlayer.set(true)
  } catch (e) {
    console.warn('Wails WindowSetSize error:', e)
    isMiniPlayer.set(true)
  }
}

export function exitMiniPlayer() {
  try {
    WindowSetMinSize(900, 600)
    WindowSetSize(1200, 800)
    WindowSetAlwaysOnTop(false)
    isMiniPlayer.set(false)
  } catch (e) {
    console.warn('Wails WindowSetSize error:', e)
    isMiniPlayer.set(false)
  }
}
