import { writable } from 'svelte/store'
import type { SearchResultItem } from '../types'
import { GetUserData, ToggleFavorite as wailsToggleFav, AddToHistory as wailsAddHist } from '../../../wailsjs/go/main/App'

interface UserDataState {
  favorites: SearchResultItem[]
  history: SearchResultItem[]
  isLoading: boolean
}

function createUserDataStore() {
  const { subscribe, set, update } = writable<UserDataState>({
    favorites: [],
    history: [],
    isLoading: false,
  })

  async function load() {
    update(s => ({ ...s, isLoading: true }))
    try {
      const data = await GetUserData()
      set({
        favorites: data.favorites || [],
        history: data.history || [],
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

  return {
    subscribe,
    load,
    toggleFavorite,
    recordHistory,
  }
}

export const userDataStore = createUserDataStore()
