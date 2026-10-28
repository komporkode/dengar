import type { backend, main } from '../../wailsjs/go/models'

export type SearchResultItem = backend.SearchResultItem
export type SearchResponse = backend.SearchResponse
export type StreamResult = backend.StreamResult
export type UserData = backend.UserData
export interface Playlist {
  id: string
  name: string
  description: string
  tracks: SearchResultItem[]
  createdAt: number
  updatedAt: number
}
export type HomeData = main.HomeData

export type PlayMode = 'audio' | 'video'
export type RepeatMode = 'off' | 'all' | 'one'
export type NavigationTab = 'home' | 'search' | 'favorites' | 'history' | 'playlist'
export type SearchFilter = 'music_songs' | 'videos' | 'all'

