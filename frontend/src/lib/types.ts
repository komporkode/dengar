import type { backend } from '../../wailsjs/go/models'

export type SearchResultItem = backend.SearchResultItem
export type StreamResult = backend.StreamResult
export type UserData = backend.UserData

export type PlayMode = 'audio' | 'video'
export type RepeatMode = 'off' | 'all' | 'one'
export type NavigationTab = 'home' | 'search' | 'favorites' | 'history'
export type SearchFilter = 'music_songs' | 'videos' | 'all'
