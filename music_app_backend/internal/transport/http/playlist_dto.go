
package http

import (
	"time"

	"github.com/google/uuid"
	"github.com/mosesmmoisebidth/music_backend/internal/playlist"
)

// --- Playlist Requests ---

type GetPlaylistsRequest struct {
	Page int `form:"page,default=1"`
	Size int `form:"size,default=20"`
}

type CreatePlaylistRequest struct {
	Title       string  `json:"title" binding:"required,min=3,max=100"`
	Description *string `json:"description,omitempty" binding:"max=500"`
}

type UpdatePlaylistRequest struct {
	Title       *string `json:"title,omitempty" binding:"omitempty,min=3,max=100"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=500"`
	IsPublic    *bool   `json:"is_public,omitempty"`
}

type AddTrackToPlaylistRequest struct {
	Provider        string `json:"provider" binding:"required"`
	ProviderTrackID string `json:"provider_track_id" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Artist          string `json:"artist" binding:"required"`
	Album           string `json:"album"`
	DurationMs      int    `json:"duration_ms"`
	ArtworkURL      string `json:"artwork_url"`
	TrackNumber     int    `json:"track_number"`
}

type ReorderPlaylistRequest struct {
	TrackID  string `json:"track_id" binding:"required"`
	Position int    `json:"position" binding:"required,gte=0"`
}

// --- Playlist Responses ---

type PlaylistTrackResponse struct {
	ID              uuid.UUID `json:"id"`
	Provider        string    `json:"provider"`
	ProviderTrackID string    `json:"provider_track_id"`
	Title           string    `json:"title"`
	Artist          string    `json:"artist"`
	Album           string    `json:"album"`
	DurationMs      int       `json:"duration_ms"`
	ArtworkURL      string    `json:"artwork_url"`
	TrackNumber     int       `json:"track_number"`
	Position        int       `json:"position"`
	AddedAt         time.Time `json:"added_at"`
}

type PlaylistResponse struct {
	ID          uuid.UUID               `json:"id"`
	UserID      uuid.UUID               `json:"user_id"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	CoverURL    *string                 `json:"cover_url,omitempty"`
	IsPublic    bool                    `json:"is_public"`
	ShareCode   *string                 `json:"share_code,omitempty"`
	TrackCount  int                     `json:"track_count"`
	Tracks      []PlaylistTrackResponse `json:"tracks"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

type SharePlaylistResponse struct {
	ShareCode string `json:"share_code"`
	ShareURL  string `json:"share_url"`
}

func mapPlaylistToResponse(p *playlist.Playlist) PlaylistResponse {
	var tracks []PlaylistTrackResponse
	for _, track := range p.Tracks {
		tracks = append(tracks, PlaylistTrackResponse{
			ID:              track.ID,
			Provider:        string(track.Provider),
			ProviderTrackID: track.ProviderTrackID,
			Title:           track.Title,
			Artist:          track.Artist,
			Album:           track.Album,
			DurationMs:      track.DurationMs,
			ArtworkURL:      track.ArtworkURL,
			TrackNumber:     track.TrackNumber,
			Position:        track.Position,
			AddedAt:         track.AddedAt,
		})
	}

	return PlaylistResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		Title:       p.Title,
		Description: p.Description,
		CoverURL:    p.CoverURL,
		IsPublic:    p.IsPublic,
		ShareCode:   p.ShareCode,
		TrackCount:  len(tracks),
		Tracks:      tracks,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
