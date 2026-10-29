import { writable, get } from 'svelte/store'
import type { SearchResultItem, RepeatMode } from '../types'

interface QueueState {
  items: SearchResultItem[]
  currentIndex: number
  repeatMode: RepeatMode
  isShuffle: boolean
  shuffleOrder: number[]
  shufflePointer: number
  isDrawerOpen: boolean
}

/**
 * Fisher-Yates shuffle algorithm that creates a randomized order of indices.
 * If fixedFirstIndex is provided, it guarantees that fixedFirstIndex appears at position 0.
 */
function generateShuffleOrder(length: number, fixedFirstIndex = -1): number[] {
  if (length <= 0) return []
  if (length === 1) return [0]

  const indices: number[] = []
  for (let i = 0; i < length; i++) {
    if (i !== fixedFirstIndex) {
      indices.push(i)
    }
  }

  // Shuffle remaining indices
  for (let i = indices.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    const temp = indices[i]
    indices[i] = indices[j]
    indices[j] = temp
  }

  if (fixedFirstIndex >= 0 && fixedFirstIndex < length) {
    return [fixedFirstIndex, ...indices]
  }

  return indices
}

function createQueueStore() {
  const { subscribe, set, update } = writable<QueueState>({
    items: [],
    currentIndex: -1,
    repeatMode: 'off',
    isShuffle: false,
    shuffleOrder: [],
    shufflePointer: -1,
    isDrawerOpen: false,
  })

  function setQueue(items: SearchResultItem[], startIndex = 0) {
    const validIndex = startIndex >= 0 && startIndex < items.length ? startIndex : 0
    update(s => {
      let shuffleOrder: number[] = []
      let shufflePointer = validIndex

      if (s.isShuffle && items.length > 0) {
        shuffleOrder = generateShuffleOrder(items.length, validIndex)
        shufflePointer = 0
      } else {
        shuffleOrder = items.map((_, i) => i)
        shufflePointer = validIndex
      }

      return {
        ...s,
        items: [...items],
        currentIndex: items.length > 0 ? validIndex : -1,
        shuffleOrder,
        shufflePointer: items.length > 0 ? shufflePointer : -1,
      }
    })
  }

  function addToQueue(item: SearchResultItem) {
    update(s => {
      const newItems = [...s.items, item]
      const newIndex = newItems.length - 1
      let newShuffleOrder = [...s.shuffleOrder]

      if (s.isShuffle && newShuffleOrder.length > 0) {
        // Insert new track randomly after current shuffle pointer so it will be played
        const insertAt = s.shufflePointer + 1 + Math.floor(Math.random() * Math.max(1, newShuffleOrder.length - s.shufflePointer))
        newShuffleOrder.splice(insertAt, 0, newIndex)
      } else {
        newShuffleOrder.push(newIndex)
      }

      return {
        ...s,
        items: newItems,
        currentIndex: s.currentIndex === -1 ? 0 : s.currentIndex,
        shuffleOrder: newShuffleOrder,
        shufflePointer: s.shufflePointer === -1 ? 0 : s.shufflePointer,
      }
    })
  }

  function removeFromQueue(index: number) {
    update(s => {
      if (index < 0 || index >= s.items.length) return s

      const newItems = s.items.filter((_, i) => i !== index)
      let newIdx = s.currentIndex

      if (index < s.currentIndex) {
        newIdx--
      } else if (index === s.currentIndex) {
        if (newIdx >= newItems.length) {
          newIdx = newItems.length - 1
        }
      }

      // Rebuild or adjust shuffle order
      const newShuffleOrder = s.shuffleOrder
        .filter(i => i !== index)
        .map(i => (i > index ? i - 1 : i))

      let newShufflePointer = s.shufflePointer
      if (newShufflePointer >= newShuffleOrder.length) {
        newShufflePointer = Math.max(0, newShuffleOrder.length - 1)
      }

      return {
        ...s,
        items: newItems,
        currentIndex: newItems.length > 0 ? newIdx : -1,
        shuffleOrder: newShuffleOrder,
        shufflePointer: newItems.length > 0 ? newShufflePointer : -1,
      }
    })
  }

  function clearQueue() {
    update(s => ({
      ...s,
      items: [],
      currentIndex: -1,
      shuffleOrder: [],
      shufflePointer: -1,
    }))
  }

  function jumpToIndex(index: number): SearchResultItem | null {
    const s = get({ subscribe })
    if (index < 0 || index >= s.items.length) return null

    let pointer = index
    if (s.isShuffle) {
      const found = s.shuffleOrder.indexOf(index)
      pointer = found !== -1 ? found : 0
    }

    update(st => ({
      ...st,
      currentIndex: index,
      shufflePointer: pointer,
    }))

    return s.items[index]
  }

  function getCurrent(): SearchResultItem | null {
    const s = get({ subscribe })
    if (s.currentIndex >= 0 && s.currentIndex < s.items.length) {
      return s.items[s.currentIndex]
    }
    return null
  }

  /**
   * Advance to the next track.
   * @param manualAdvance If true (e.g. user clicked Next button), overrides repeatMode === 'one' to move forward.
   */
  function next(manualAdvance = false): SearchResultItem | null {
    const s = get({ subscribe })
    if (s.items.length === 0) return null

    // If repeatMode is 'one' and this is an automatic track completion (not user pressing Next)
    if (s.repeatMode === 'one' && !manualAdvance) {
      return s.items[s.currentIndex] || null
    }

    // Shuffle playback
    if (s.isShuffle && s.items.length > 1) {
      const nextPointer = s.shufflePointer + 1

      if (nextPointer < s.shuffleOrder.length) {
        const nextIdx = s.shuffleOrder[nextPointer]
        update(st => ({
          ...st,
          currentIndex: nextIdx,
          shufflePointer: nextPointer,
        }))
        return s.items[nextIdx]
      }

      // Reached the end of the shuffled bag
      if (s.repeatMode === 'all') {
        // Reshuffle all items, ensuring first track is not the same as the last track
        const lastTrackIdx = s.currentIndex
        let newOrder = generateShuffleOrder(s.items.length)
        if (newOrder.length > 1 && newOrder[0] === lastTrackIdx) {
          const swap = newOrder[0]
          newOrder[0] = newOrder[1]
          newOrder[1] = swap
        }
        const firstIdx = newOrder[0]
        update(st => ({
          ...st,
          shuffleOrder: newOrder,
          shufflePointer: 0,
          currentIndex: firstIdx,
        }))
        return s.items[firstIdx]
      }

      // Reached end and repeat is off
      return null
    }

    // Linear sequential playback
    let nextIdx = s.currentIndex + 1
    if (nextIdx >= s.items.length) {
      if (s.repeatMode === 'all') {
        nextIdx = 0
      } else {
        return null // Reached end of queue
      }
    }

    update(st => ({
      ...st,
      currentIndex: nextIdx,
      shufflePointer: nextIdx,
    }))
    return s.items[nextIdx]
  }

  function prev(): SearchResultItem | null {
    const s = get({ subscribe })
    if (s.items.length === 0) return null

    // Shuffle playback
    if (s.isShuffle && s.items.length > 1) {
      let prevPointer = s.shufflePointer - 1
      if (prevPointer < 0) {
        if (s.repeatMode === 'all') {
          prevPointer = s.shuffleOrder.length - 1
        } else {
          prevPointer = 0
        }
      }
      const prevIdx = s.shuffleOrder[prevPointer]
      update(st => ({
        ...st,
        currentIndex: prevIdx,
        shufflePointer: prevPointer,
      }))
      return s.items[prevIdx]
    }

    // Linear sequential playback
    let prevIdx = s.currentIndex - 1
    if (prevIdx < 0) {
      prevIdx = s.repeatMode === 'all' ? s.items.length - 1 : 0
    }

    update(st => ({
      ...st,
      currentIndex: prevIdx,
      shufflePointer: prevIdx,
    }))
    return s.items[prevIdx]
  }

  function cycleRepeat() {
    update(s => {
      const nextMode: Record<RepeatMode, RepeatMode> = {
        off: 'all',
        all: 'one',
        one: 'off',
      }
      return { ...s, repeatMode: nextMode[s.repeatMode] }
    })
  }

  function toggleShuffle() {
    update(s => {
      const nextShuffle = !s.isShuffle
      let newOrder: number[] = []
      let newPointer = s.currentIndex

      if (nextShuffle && s.items.length > 0) {
        // Build shuffle order with current track as the active pointer 0
        newOrder = generateShuffleOrder(s.items.length, s.currentIndex)
        newPointer = 0
      } else {
        newOrder = s.items.map((_, i) => i)
        newPointer = s.currentIndex
      }

      return {
        ...s,
        isShuffle: nextShuffle,
        shuffleOrder: newOrder,
        shufflePointer: newPointer,
      }
    })
  }

  function enableShuffle() {
    update(s => {
      let newOrder: number[] = []
      let newPointer = s.currentIndex >= 0 ? s.currentIndex : 0
      if (s.items.length > 0) {
        newOrder = generateShuffleOrder(s.items.length, newPointer)
        newPointer = 0
      }
      return {
        ...s,
        isShuffle: true,
        shuffleOrder: newOrder,
        shufflePointer: newPointer,
      }
    })
  }

  function toggleDrawer(open?: boolean) {
    update(s => ({ ...s, isDrawerOpen: open !== undefined ? open : !s.isDrawerOpen }))
  }

  return {
    subscribe,
    setQueue,
    addToQueue,
    removeFromQueue,
    clearQueue,
    jumpToIndex,
    getCurrent,
    next,
    prev,
    cycleRepeat,
    toggleShuffle,
    enableShuffle,
    toggleDrawer,
  }
}

export const queueStore = createQueueStore()

const savedAutoplay = typeof localStorage !== 'undefined' ? localStorage.getItem('dengar_autoplay') !== 'false' : true
export const isAutoplayEnabled = writable<boolean>(savedAutoplay)

isAutoplayEnabled.subscribe(val => {
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('dengar_autoplay', String(val))
  }
})

export function toggleAutoplay(): boolean {
  const next = !get(isAutoplayEnabled)
  isAutoplayEnabled.set(next)
  return next
}
