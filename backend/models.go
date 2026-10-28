package backend

// SearchResultItem represents normalized track/video metadata for the player.
type SearchResultItem struct {
	ID                string `json:"id"`
	Title             string `json:"title"`
	UploaderName      string `json:"uploaderName"`
	Duration          int64  `json:"duration"`
	DurationFormatted string `json:"durationFormatted"`
	ThumbnailURL      string `json:"thumbnailUrl"`
}

// SearchResponse holds search results and an optional continuation token for pagination.
type SearchResponse struct {
	Items             []SearchResultItem `json:"items"`
	ContinuationToken string             `json:"continuationToken,omitempty"`
}

// StreamResult represents resolved media URLs and metadata for playback.
type StreamResult struct {
	VideoID   string `json:"videoId"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Thumbnail string `json:"thumbnail"`
	AudioURL  string `json:"audioUrl"`
	VideoURL  string `json:"videoUrl"`
	Duration  int64  `json:"duration"`
}

// UserData stores bookmarks and playback history.
type UserData struct {
	Favorites []SearchResultItem `json:"favorites"`
	History   []SearchResultItem `json:"history"`
}
