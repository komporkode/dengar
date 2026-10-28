import { writable } from 'svelte/store'
import type { Playlist, SearchResultItem } from '../types'

export const isCreatePlaylistOpen = writable<boolean>(false)
export const editingPlaylist = writable<Playlist | null>(null)
export const trackToAddToPlaylist = writable<SearchResultItem | null>(null)
export const sharingPlaylist = writable<Playlist | null>(null)
export const isImportPlaylistOpen = writable<boolean>(false)
export const activePlaylistId = writable<string | null>(null)

export function openCreatePlaylistModal() {
  editingPlaylist.set(null)
  isCreatePlaylistOpen.set(true)
}

export function openEditPlaylistModal(playlist: Playlist) {
  editingPlaylist.set(playlist)
  isCreatePlaylistOpen.set(true)
}

export function openAddToPlaylistModal(track: SearchResultItem) {
  trackToAddToPlaylist.set(track)
}

export function openSharePlaylistModal(playlist: Playlist) {
  sharingPlaylist.set(playlist)
}

export function openImportPlaylistModal() {
  isImportPlaylistOpen.set(true)
}
