package library

import (
	"time"

	"github.com/google/uuid"
)

// DownloadState represents the state of a download.
type DownloadState string

const (
	StatePending   DownloadState = "pending"
	StateDownloading DownloadState = "downloading"
	StateCompleted  DownloadState = "completed"
	StateFailed      DownloadState = "failed"
	StatePaused      DownloadState = "paused"
)

// Favorite represents a user's favorite track.
type Favorite struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"`
	Provider        string    `gorm:"not null;size:20"`
	ProviderTrackID string    `gorm:"not null;size:100"`
	Title           string    `gorm:"not null;size:255"`
	Artist          string    `gorm:"not null;size:255"`
	Album           string    `gorm:"size:255"`
	DurationMs      int
	ArtworkURL      string `gorm:"size:1024"`
	AddedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// History represents a user's listening history.
type History struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"`
	Provider        string    `gorm:"not null;size:20"`
	ProviderTrackID string    `gorm:"not null;size:100"`
	Title           string    `gorm:"not null;size:255"`
	Artist          string    `gorm:"not null;size:255"`
	Album           string    `gorm:"size:255"`
	DurationMs      int
	ArtworkURL      string `gorm:"size:1024"`
	PlayedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// Download represents a user's downloaded track.
type Download struct {
	ID              uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID          uuid.UUID     `gorm:"type:uuid;not null;index"`
	Provider        string        `gorm:"not null;size:20"`
	ProviderTrackID string        `gorm:"not null;size:100"`
	Title           string        `gorm:"not null;size:255"`
	Artist          string        `gorm:"not null;size:255"`
	Album           string        `gorm:"size:255"`
	DurationMs      int
	ArtworkURL      string        `gorm:"size:1024"`
	State           DownloadState `gorm:"not null;size:20"`
	LocalPath       *string       `gorm:"size:1024"`
	FileSize        *int64
	Quality         string `gorm:"not null;size:50"` // Changed to string
	CreatedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP"`
}

// TrackData is used to pass track information to the service layer.
type TrackData struct {
	Provider        string
	ProviderTrackID string
	Title           string
	Artist          string
	Album           string
	DurationMs      int
	ArtworkURL      string
}