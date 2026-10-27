import { writable, get } from 'svelte/store'
import type { SearchResultItem, RepeatMode } from '../types'

interface QueueState {
  items: SearchResultItem[]
  currentIndex: number
  repeatMode: RepeatMode
  isShuffle: boolean
  isDrawerOpen: boolean
}

function createQueueStore() {
  const { subscribe, set, update } = writable<QueueState>({
    items: [],
    currentIndex: -1,
    repeatMode: 'off',
    isShuffle: false,
    isDrawerOpen: false,
  })

  function setQueue(items: SearchResultItem[], startIndex = 0) {
    update(s => ({
      ...s,
      items: [...items],
      currentIndex: startIndex >= 0 && startIndex < items.length ? startIndex : 0,
    }))
  }

  function addToQueue(item: SearchResultItem) {
    update(s => ({
      ...s,
      items: [...s.items, item],
      currentIndex: s.currentIndex === -1 ? 0 : s.currentIndex,
    }))
  }

  function removeFromQueue(index: number) {
    update(s => {
      const newItems = s.items.filter((_, i) => i !== index)
      let newIdx = s.currentIndex
      if (index < s.currentIndex) {
        newIdx--
      } else if (index === s.currentIndex && newIdx >= newItems.length) {
        newIdx = newItems.length - 1
      }
      return {
        ...s,
        items: newItems,
        currentIndex: newIdx,
      }
    })
  }

  function clearQueue() {
    update(s => ({ ...s, items: [], currentIndex: -1 }))
  }

  function next(): SearchResultItem | null {
    const s = get({ subscribe })
    if (s.items.length === 0) return null

    if (s.repeatMode === 'one') {
      return s.items[s.currentIndex] || null
    }

    if (s.isShuffle && s.items.length > 1) {
      let rand = Math.floor(Math.random() * s.items.length)
      if (rand === s.currentIndex) {
        rand = (rand + 1) % s.items.length
      }
      update(st => ({ ...st, currentIndex: rand }))
      return s.items[rand]
    }

    let nextIdx = s.currentIndex + 1
    if (nextIdx >= s.items.length) {
      if (s.repeatMode === 'all') {
        nextIdx = 0
      } else {
        return null // Reached end of queue
      }
    }

    update(st => ({ ...st, currentIndex: nextIdx }))
    return s.items[nextIdx]
  }

  function prev(): SearchResultItem | null {
    const s = get({ subscribe })
    if (s.items.length === 0) return null

    let prevIdx = s.currentIndex - 1
    if (prevIdx < 0) {
      prevIdx = s.repeatMode === 'all' ? s.items.length - 1 : 0
    }

    update(st => ({ ...st, currentIndex: prevIdx }))
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
    update(s => ({ ...s, isShuffle: !s.isShuffle }))
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
    next,
    prev,
    cycleRepeat,
    toggleShuffle,
    toggleDrawer,
  }
}

export const queueStore = createQueueStore()
