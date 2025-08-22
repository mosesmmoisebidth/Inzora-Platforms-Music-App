package music

import (
	"context"
	"time"
)

// Track represents a music track from any provider
type Track struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Artist        string  `json:"artist"`
	Album         string  `json:"album"`
	Duration      int64   `json:"duration_ms"`
	ArtworkURL    string  `json:"artwork_url"`
	PreviewURL    string  `json:"preview_url,omitempty"`
	TrackNumber   int     `json:"track_number,omitempty"`
	ReleaseDate   string  `json:"release_date,omitempty"`
	Genre         string  `json:"genre,omitempty"`
	Provider      string  `json:"provider"`
	ExternalURL   string  `json:"external_url,omitempty"`
	Explicit      bool    `json:"explicit"`
	Popularity    int     `json:"popularity,omitempty"`
}

// PlaylistSummary represents a playlist summary from any provider
type PlaylistSummary struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	TrackCount  int    `json:"track_count"`
	Provider    string `json:"provider"`
	ExternalURL string `json:"external_url,omitempty"`
	Creator     string `json:"creator,omitempty"`
}

// Category represents a music category
type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url,omitempty"`
}

// PageInfo represents pagination information
type PageInfo struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	Total      int64 `json:"total"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
	TotalPages int   `json:"total_pages"`
}

// SearchFilters represents search filter options
type SearchFilters struct {
	Genre       string `json:"genre,omitempty"`
	Year        string `json:"year,omitempty"`
	Explicit    *bool  `json:"explicit,omitempty"`
	Duration    string `json:"duration,omitempty"` // short, medium, long
	SortBy      string `json:"sort_by,omitempty"`  // relevance, popularity, date
	SortOrder   string `json:"sort_order,omitempty"` // asc, desc
}

// MusicProvider defines the interface that all music providers must implement
type MusicProvider interface {
	// GetName returns the provider name
	GetName() string
	
	// SearchTracks searches for tracks with optional filters
	SearchTracks(ctx context.Context, query string, page, size int, filters *SearchFilters) ([]Track, *PageInfo, error)
	
	// GetTrack gets a specific track by ID
	GetTrack(ctx context.Context, trackID string) (*Track, error)
	
	// GetTopCharts gets top charts for a country
	GetTopCharts(ctx context.Context, country string, page, size int) ([]Track, *PageInfo, error)
	
	// GetCategories gets available music categories
	GetCategories(ctx context.Context) ([]Category, error)
	
	// GetPlaylistsByCategory gets playlists for a specific category
	GetPlaylistsByCategory(ctx context.Context, categoryID string, page, size int) ([]PlaylistSummary, *PageInfo, error)
	
	// IsHealthy checks if the provider is healthy and accessible
	IsHealthy(ctx context.Context) error
}

// ProviderConfig holds common configuration for music providers
type ProviderConfig struct {
	Timeout   time.Duration
	RateLimit int
	CacheTTL  time.Duration
	UserAgent string
}

// ProviderError represents an error from a music provider
type ProviderError struct {
	Provider string
	Message  string
	Code     string
	Err      error
}

func (e *ProviderError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *ProviderError) Unwrap() error {
	return e.Err
}

// NewProviderError creates a new provider error
func NewProviderError(provider, message, code string, err error) *ProviderError {
	return &ProviderError{
		Provider: provider,
		Message:  message,
		Code:     code,
		Err:      err,
	}
}
