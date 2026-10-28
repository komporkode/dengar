import { writable } from 'svelte/store'
import type { SearchResultItem, Playlist } from '../types'
import {
  GetUserData,
  ToggleFavorite as wailsToggleFav,
  AddToHistory as wailsAddHist,
  CreatePlaylist as wailsCreatePl,
  UpdatePlaylist as wailsUpdatePl,
  DeletePlaylist as wailsDeletePl,
  AddTrackToPlaylist as wailsAddTrackPl,
  RemoveTrackFromPlaylist as wailsRemoveTrackPl,
  ExportPlaylist as wailsExportPl,
  ImportPlaylist as wailsImportPl,
} from '../../../wailsjs/go/main/App'

interface UserDataState {
  favorites: SearchResultItem[]
  history: SearchResultItem[]
  playlists: Playlist[]
  isLoading: boolean
}

function createUserDataStore() {
  const { subscribe, set, update } = writable<UserDataState>({
    favorites: [],
    history: [],
    playlists: [],
    isLoading: false,
  })

  async function load() {
    update(s => ({ ...s, isLoading: true }))
    try {
      const data = await GetUserData()
      set({
        favorites: data.favorites || [],
        history: data.history || [],
        playlists: data.playlists || [],
        isLoading: false,
      })
    } catch (e) {
      console.error('Failed to load user data:', e)
      update(s => ({ ...s, isLoading: false }))
    }
  }

  async function toggleFavorite(item: SearchResultItem): Promise<boolean> {
    try {
      const isFav = await wailsToggleFav(item)
      update(s => {
        let newFavs: SearchResultItem[]
        if (isFav) {
          newFavs = [item, ...s.favorites.filter(f => f.id !== item.id)]
        } else {
          newFavs = s.favorites.filter(f => f.id !== item.id)
        }
        return { ...s, favorites: newFavs }
      })
      return isFav
    } catch (e) {
      console.error('Failed to toggle favorite:', e)
      return false
    }
  }

  async function recordHistory(item: SearchResultItem) {
    try {
      await wailsAddHist(item)
      update(s => {
        const filtered = s.history.filter(h => h.id !== item.id)
        const newHist = [item, ...filtered].slice(0, 100)
        return { ...s, history: newHist }
      })
    } catch (e) {
      console.error('Failed to record history:', e)
    }
  }

  async function createPlaylist(name: string, description: string = ''): Promise<Playlist | null> {
    try {
      const pl = await wailsCreatePl(name, description)
      update(s => ({ ...s, playlists: [...s.playlists, pl] }))
      return pl
    } catch (e) {
      console.error('Failed to create playlist:', e)
      return null
    }
  }

  async function updatePlaylist(id: string, name: string, description: string): Promise<Playlist | null> {
    try {
      const updated = await wailsUpdatePl(id, name, description)
      update(s => ({
        ...s,
        playlists: s.playlists.map(p => (p.id === id ? updated : p)),
      }))
      return updated
    } catch (e) {
      console.error('Failed to update playlist:', e)
      return null
    }
  }

  async function deletePlaylist(id: string): Promise<void> {
    try {
      await wailsDeletePl(id)
      update(s => ({
        ...s,
        playlists: s.playlists.filter(p => p.id !== id),
      }))
    } catch (e) {
      console.error('Failed to delete playlist:', e)
    }
  }

  async function addTrackToPlaylist(playlistId: string, track: SearchResultItem): Promise<void> {
    try {
      await wailsAddTrackPl(playlistId, track)
      update(s => ({
        ...s,
        playlists: s.playlists.map(p => {
          if (p.id === playlistId) {
            const exists = (p.tracks || []).some(t => t.id === track.id)
            const tracks = exists ? p.tracks : [...(p.tracks || []), track]
            return { ...p, tracks, updatedAt: Date.now() }
          }
          return p
        }),
      }))
    } catch (e) {
      console.error('Failed to add track to playlist:', e)
    }
  }

  async function removeTrackFromPlaylist(playlistId: string, trackId: string): Promise<void> {
    try {
      await wailsRemoveTrackPl(playlistId, trackId)
      update(s => ({
        ...s,
        playlists: s.playlists.map(p => {
          if (p.id === playlistId) {
            return {
              ...p,
              tracks: (p.tracks || []).filter(t => t.id !== trackId),
              updatedAt: Date.now(),
            }
          }
          return p
        }),
      }))
    } catch (e) {
      console.error('Failed to remove track from playlist:', e)
    }
  }

  async function exportPlaylist(playlistId: string): Promise<string> {
    return await wailsExportPl(playlistId)
  }

  async function importPlaylist(payload: string): Promise<Playlist | null> {
    try {
      const pl = await wailsImportPl(payload)
      update(s => ({
        ...s,
        playlists: [...s.playlists.filter(p => p.id !== pl.id), pl],
      }))
      return pl
    } catch (e) {
      console.error('Failed to import playlist:', e)
      throw e
    }
  }

  return {
    subscribe,
    load,
    toggleFavorite,
    recordHistory,
    createPlaylist,
    updatePlaylist,
    deletePlaylist,
    addTrackToPlaylist,
    removeTrackFromPlaylist,
    exportPlaylist,
    importPlaylist,
  }
}

export const userDataStore = createUserDataStore()

