import { writable } from 'svelte/store'
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
    WindowSetMinSize(360, 140)
    WindowSetSize(380, 150)
    WindowSetAlwaysOnTop(true)
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
