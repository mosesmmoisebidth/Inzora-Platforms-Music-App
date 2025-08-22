
package http

import (
	"github.com/mosesmmoisebidth/music_backend/internal/music"
)

// --- Music Requests ---

type SearchTracksRequest struct {
	Query    string  `form:"q" binding:"required"`
	Page     int     `form:"page,default=1"`
	Size     int     `form:"size,default=20"`
	Provider *string `form:"provider,omitempty"`
	Genre    *string `form:"genre,omitempty"`
	Year     *int    `form:"year,omitempty"`
	Explicit *bool   `form:"explicit,omitempty"`
}

type GetTrackRequest struct {
	TrackID  string `uri:"trackId" binding:"required"`
	Provider string `form:"provider" binding:"required"`
}

type GetTopChartsRequest struct {
	Country  string  `form:"country,default=US"`
	Page     int     `form:"page,default=1"`
	Size     int     `form:"size,default=20"`
	Provider *string `form:"provider,omitempty"`
}

// --- Music Responses ---

type TrackResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Duration    int64  `json:"duration"` // Changed to int64
	ArtworkURL  string `json:"artwork_url"`
	PreviewURL  string `json:"preview_url"`
	TrackNumber int    `json:"track_number"`
	ReleaseDate string `json:"release_date"` // Changed to string
	Genre       string `json:"genre"`
	Provider    string `json:"provider"`
	ExternalURL string `json:"external_url"`
	Explicit    bool   `json:"explicit"`
	Popularity  int    `json:"popularity"`
}

func mapTrackToResponse(t *music.Track) TrackResponse {
	return TrackResponse{
		ID:          t.ID,
		Title:       t.Title,
		Artist:      t.Artist,
		Album:       t.Album,
		Duration:    t.Duration, // Now matches int64
		ArtworkURL:  t.ArtworkURL,
		PreviewURL:  t.PreviewURL,
		TrackNumber: t.TrackNumber,
		ReleaseDate: t.ReleaseDate, // Now matches string
		Genre:       t.Genre,
		Provider:    t.Provider,
		ExternalURL: t.ExternalURL,
		Explicit:    t.Explicit,
		Popularity:  t.Popularity,
	}
}
