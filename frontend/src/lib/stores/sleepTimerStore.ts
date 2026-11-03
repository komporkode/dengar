import { writable, get } from 'svelte/store'
import { playerStore } from './playerStore'
import { showToast } from './settingsStore'

export const isSleepTimerModalOpen = writable<boolean>(false)
export const sleepTimerRemaining = writable<number | null>(null) // in seconds
export const sleepTimerMode = writable<'minutes' | 'end_of_track' | null>(null)

let tickInterval: any = null
let unsubscribePlayer: (() => void) | null = null
let isFadingOut = false

export function openSleepTimerModal() {
  isSleepTimerModalOpen.set(true)
}

export function closeSleepTimerModal() {
  isSleepTimerModalOpen.set(false)
}

export function setTimerMinutes(minutes: number) {
  clearActiveTimer()
  if (minutes <= 0) return

  const totalSeconds = Math.round(minutes * 60)
  sleepTimerMode.set('minutes')
  sleepTimerRemaining.set(totalSeconds)
  isFadingOut = false

  showToast(`Sleep timer disetel untuk ${minutes} menit`)

  tickInterval = setInterval(() => {
    sleepTimerRemaining.update(remaining => {
      if (remaining === null) return null
      const next = remaining - 1

      // 5-second smooth fade-out before pause
      if (next === 5 && !isFadingOut) {
        isFadingOut = true
        playerStore.fadeVolumeToZero(5000, () => {
          playerStore.pause()
          showToast('Sleep timer: Musik dihentikan.')
          clearActiveTimer()
        })
      }

      if (next <= 0) {
        if (!isFadingOut) {
          playerStore.pause()
          showToast('Sleep timer: Musik dihentikan.')
          clearActiveTimer()
        }
        return null
      }
      return next
    })
  }, 1000)
}

export function setEndOfTrackTimer() {
  clearActiveTimer()
  sleepTimerMode.set('end_of_track')
  isFadingOut = false

  const playerState = get(playerStore)
  if (playerState.duration > 0 && playerState.currentTime < playerState.duration) {
    const left = Math.max(0, Math.round(playerState.duration - playerState.currentTime))
    sleepTimerRemaining.set(left)
  }

  showToast('Sleep timer aktif: Berhenti di akhir lagu ini')

  unsubscribePlayer = playerStore.subscribe(state => {
    if (!state.currentTrack || !state.isPlaying) return

    const remainingSec = state.duration - state.currentTime
    if (remainingSec > 0) {
      sleepTimerRemaining.set(Math.round(remainingSec))
    }

    if (remainingSec <= 5 && remainingSec > 0 && !isFadingOut) {
      isFadingOut = true
      playerStore.fadeVolumeToZero(5000, () => {
        playerStore.pause()
        showToast('Sleep timer: Musik dihentikan.')
        clearActiveTimer()
      })
    }
  })
}

export function cancelSleepTimer() {
  clearActiveTimer()
  playerStore.cancelVolumeFade()
  showToast('Sleep timer dimatikan')
}

function clearActiveTimer() {
  if (tickInterval) {
    clearInterval(tickInterval)
    tickInterval = null
  }
  if (unsubscribePlayer) {
    unsubscribePlayer()
    unsubscribePlayer = null
  }
  isFadingOut = false
  sleepTimerRemaining.set(null)
  sleepTimerMode.set(null)
}
