package http

import (
	"time"

	"github.com/google/uuid"
	"github.com/mosesmmoisebidth/music_backend/internal/library"
)

// --- Library Requests ---

type GetFavoritesRequest struct {
	Page int `form:"page,default=1"`
	Size int `form:"size,default=20"`
}

type AddFavoriteRequest struct {
	Provider        string `json:"provider" binding:"required"`
	ProviderTrackID string `json:"provider_track_id" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Artist          string `json:"artist" binding:"required"`
	Album           string `json:"album"`
	DurationMs      int    `json:"duration_ms"`
	ArtworkURL      string `json:"artwork_url"`
}

type GetHistoryRequest struct {
	Page int `form:"page,default=1"`
	Size int `form:"size,default=50"`
}

type AddHistoryRequest struct {
	Provider        string `json:"provider" binding:"required"`
	ProviderTrackID string `json:"provider_track_id" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Artist          string `json:"artist" binding:"required"`
	Album           string `json:"album"`
	DurationMs      int    `json:"duration_ms"`
	ArtworkURL      string `json:"artwork_url"`
}

// --- Library Responses ---

type FavoriteResponse struct {
	ID              uuid.UUID `json:"id"`
	Provider        string    `json:"provider"`
	ProviderTrackID string    `json:"provider_track_id"`
	Title           string    `json:"title"`
	Artist          string    `json:"artist"`
	Album           string    `json:"album"`
	DurationMs      int       `json:"duration_ms"`
	ArtworkURL      string    `json:"artwork_url"`
	AddedAt         time.Time `json:"added_at"`
}

type HistoryResponse struct {
	ID              uuid.UUID `json:"id"`
	Provider        string    `json:"provider"`
	ProviderTrackID string    `json:"provider_track_id"`
	Title           string    `json:"title"`
	Artist          string    `json:"artist"`
	Album           string    `json:"album"`
	DurationMs      int       `json:"duration_ms"`
	ArtworkURL      string    `json:"artwork_url"`
	PlayedAt        time.Time `json:"played_at"`
}

func mapFavoriteToResponse(f *library.Favorite) FavoriteResponse {
	return FavoriteResponse{
		ID:              f.ID,
		Provider:        f.Provider,
		ProviderTrackID: f.ProviderTrackID,
		Title:           f.Title,
		Artist:          f.Artist,
		Album:           f.Album,
		DurationMs:      f.DurationMs,
		ArtworkURL:      f.ArtworkURL,
		AddedAt:         f.AddedAt,
	}
}

func mapHistoryToResponse(h *library.History) HistoryResponse {
	return HistoryResponse{
		ID:              h.ID,
		Provider:        h.Provider,
		ProviderTrackID: h.ProviderTrackID,
		Title:           h.Title,
		Artist:          h.Artist,
		Album:           h.Album,
		DurationMs:      h.DurationMs,
		ArtworkURL:      h.ArtworkURL,
		PlayedAt:        h.PlayedAt,
	}
}