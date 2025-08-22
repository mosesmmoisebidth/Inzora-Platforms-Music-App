package playlist

import (
	"time"

	"github.com/google/uuid"
)

// MusicProvider represents the source of the music track (e.g., "itunes", "spotify").
type MusicProvider string

const (
	ITunes MusicProvider = "itunes"
	Spotify MusicProvider = "spotify"
)

// Playlist represents a user-created playlist.
type Playlist struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Title       string    `gorm:"not null;size:100"`
	Description string    `gorm:"size:500"`
	CoverURL    *string   `gorm:"size:1024"`
	IsPublic    bool      `gorm:"default:false"`
	ShareCode   *string   `gorm:"uniqueIndex;size:10"`
	Tracks      []PlaylistTrack `gorm:"foreignKey:PlaylistID;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// PlaylistTrack represents a track within a playlist.
type PlaylistTrack struct {
	ID              uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	PlaylistID      uuid.UUID     `gorm:"type:uuid;not null;index"`
	Provider        MusicProvider `gorm:"not null;size:20"`
	ProviderTrackID string        `gorm:"not null;size:100"`
	Title           string        `gorm:"not null;size:255"`
	Artist          string        `gorm:"not null;size:255"`
	Album           string        `gorm:"size:255"`
	DurationMs      int
	ArtworkURL      string `gorm:"size:1024"`
	TrackNumber     int
	Position        int    `gorm:"not null"`
	AddedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// TrackData is used to pass track information to the service layer.
type TrackData struct {
	Provider        MusicProvider
	ProviderTrackID string
	Title           string
	Artist          string
	Album           string
	DurationMs      int
	ArtworkURL      string
}

// TrackPosition is used for reordering tracks in a playlist.
type TrackPosition struct {
	TrackID  uuid.UUID
	Position int
}