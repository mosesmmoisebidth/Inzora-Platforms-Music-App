package playlist

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository provides access to the playlist storage.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new playlist repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new playlist.
func (r *Repository) Create(ctx context.Context, playlist *Playlist) error {
	return r.db.WithContext(ctx).Create(playlist).Error
}

// GetByID retrieves a playlist by its ID, preloading its tracks.
func (r *Repository) GetByID(ctx context.Context, playlistID uuid.UUID) (*Playlist, error) {
	var playlist Playlist
	err := r.db.WithContext(ctx).Preload("Tracks").First(&playlist, playlistID).Error
	return &playlist, err
}

// GetUserPlaylists retrieves a paginated list of playlists for a specific user.
func (r *Repository) GetUserPlaylists(ctx context.Context, userID uuid.UUID, page, size int) ([]Playlist, int64, error) {
	var playlists []Playlist
	var total int64

	db := r.db.WithContext(ctx).Model(&Playlist{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := db.Order("created_at DESC").Limit(size).Offset(offset).Find(&playlists).Error

	return playlists, total, err
}

// Update updates a playlist's details.
func (r *Repository) Update(ctx context.Context, playlist *Playlist) error {
	return r.db.WithContext(ctx).Save(playlist).Error
}

// Delete removes a playlist.
func (r *Repository) Delete(ctx context.Context, playlistID uuid.UUID) error {
	// GORM will automatically handle deleting associated tracks if the relationship is configured with cascading deletes.
	return r.db.WithContext(ctx).Select("Tracks").Delete(&Playlist{ID: playlistID}).Error
}

// AddTrack adds a track to a playlist.
func (r *Repository) AddTrack(ctx context.Context, track *PlaylistTrack) error {
	return r.db.WithContext(ctx).Create(track).Error
}

// GetTrack retrieves a specific track from a playlist.
func (r *Repository) GetTrack(ctx context.Context, playlistID, trackID uuid.UUID) (*PlaylistTrack, error) {
	var track PlaylistTrack
	err := r.db.WithContext(ctx).Where("playlist_id = ? AND id = ?", playlistID, trackID).First(&track).Error
	return &track, err
}

// RemoveTrack removes a track from a playlist.
func (r *Repository) RemoveTrack(ctx context.Context, playlistID, trackID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("playlist_id = ? AND id = ?", playlistID, trackID).Delete(&PlaylistTrack{}).Error
}

// GetMaxPosition returns the highest position number for tracks in a playlist.
func (r *Repository) GetMaxPosition(ctx context.Context, playlistID uuid.UUID) (int, error) {
	var maxPos int
	err := r.db.WithContext(ctx).Model(&PlaylistTrack{}).Where("playlist_id = ?", playlistID).Select("COALESCE(MAX(position), 0)").Row().Scan(&maxPos)
	return maxPos, err
}

// UpdateTrackPositions updates the positions of multiple tracks in a playlist.
func (r *Repository) UpdateTrackPositions(ctx context.Context, tracks []PlaylistTrack) error {
	return r.db.WithContext(ctx).Save(&tracks).Error
}