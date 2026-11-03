import { writable } from 'svelte/store'
import type { SearchResultItem } from '../types'
import {
  GetOfflineTracks,
  SaveTrackOffline as wailsSaveOffline,
  RemoveOfflineTrack as wailsRemoveOffline,
} from '../../../wailsjs/go/main/App'

export interface OfflineState {
  tracks: SearchResultItem[]
  trackIds: Set<string>
  isLoading: boolean
}

function createOfflineStore() {
  const { subscribe, set, update } = writable<OfflineState>({
    tracks: [],
    trackIds: new Set<string>(),
    isLoading: false,
  })

  async function load() {
    update(s => ({ ...s, isLoading: true }))
    try {
      const items = (await GetOfflineTracks()) || []
      const ids = new Set(items.map(t => t.id))
      set({
        tracks: items,
        trackIds: ids,
        isLoading: false,
      })
    } catch (e) {
      console.warn('Failed to load offline tracks:', e)
      update(s => ({ ...s, isLoading: false }))
    }
  }

  async function saveTrack(track: SearchResultItem): Promise<boolean> {
    try {
      await wailsSaveOffline(track)
      update(s => {
        const nextIds = new Set(s.trackIds)
        nextIds.add(track.id)
        const exists = s.tracks.some(t => t.id === track.id)
        const nextTracks = exists ? s.tracks : [track, ...s.tracks]
        return {
          ...s,
          trackIds: nextIds,
          tracks: nextTracks,
        }
      })
      return true
    } catch (err) {
      console.error('Failed to save offline track:', err)
      return false
    }
  }

  async function removeTrack(trackId: string): Promise<boolean> {
    try {
      await wailsRemoveOffline(trackId)
      update(s => {
        const nextIds = new Set(s.trackIds)
        nextIds.delete(trackId)
        return {
          ...s,
          trackIds: nextIds,
          tracks: s.tracks.filter(t => t.id !== trackId),
        }
      })
      return true
    } catch (err) {
      console.error('Failed to remove offline track:', err)
      return false
    }
  }

  return {
    subscribe,
    load,
    saveTrack,
    removeTrack,
  }
}

export const offlineStore = createOfflineStore()

/**
 * Fallback handler for cover images: if external URL fails (e.g. offline),
 * switch to local cover proxy which serves from disk cache.
 */
export function handleCoverError(e: Event, videoId?: string) {
  const img = e.currentTarget as HTMLImageElement | null
  if (img && !img.dataset.fallback && videoId) {
    img.dataset.fallback = 'true'
    img.src = `/api/cover-proxy?cache_id=${videoId}`
  }
}

