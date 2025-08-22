package library

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository provides access to the library storage.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new library repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// --- Favorites ---

// AddFavorite adds a track to a user's favorites.
func (r *Repository) AddFavorite(ctx context.Context, favorite *Favorite) error {
	return r.db.WithContext(ctx).Create(favorite).Error
}

// GetFavorites retrieves a paginated list of a user's favorite tracks.
func (r *Repository) GetFavorites(ctx context.Context, userID uuid.UUID, page, size int) ([]Favorite, int64, error) {
	var favorites []Favorite
	var total int64

	db := r.db.WithContext(ctx).Model(&Favorite{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := db.Order("added_at DESC").Limit(size).Offset(offset).Find(&favorites).Error

	return favorites, total, err
}

// RemoveFavorite removes a track from a user's favorites.
func (r *Repository) RemoveFavorite(ctx context.Context, userID, favoriteID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, favoriteID).Delete(&Favorite{}).Error
}

// FindFavoriteByTrackID checks if a track is already in the user's favorites.
func (r *Repository) FindFavoriteByTrackID(ctx context.Context, userID uuid.UUID, provider string, providerTrackID string) (*Favorite, error) {
	var favorite Favorite
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND provider = ? AND provider_track_id = ?", userID, provider, providerTrackID).
		First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

// --- History ---

// AddHistory adds a track to a user's listening history.
func (r *Repository) AddHistory(ctx context.Context, history *History) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetHistory retrieves a paginated list of a user's listening history.
func (r *Repository) GetHistory(ctx context.Context, userID uuid.UUID, page, size int) ([]History, int64, error) {
	var history []History
	var total int64

	db := r.db.WithContext(ctx).Model(&History{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := db.Order("played_at DESC").Limit(size).Offset(offset).Find(&history).Error

	return history, total, err
}

// --- Downloads ---

// AddDownload adds a track to the user's download list.
func (r *Repository) AddDownload(ctx context.Context, download *Download) error {
	return r.db.WithContext(ctx).Create(download).Error
}

// GetDownloads retrieves a paginated list of a user's downloads, optionally filtered by state.
func (r *Repository) GetDownloads(ctx context.Context, userID uuid.UUID, page, size int, state *DownloadState) ([]Download, int64, error) {
	var downloads []Download
	var total int64

	db := r.db.WithContext(ctx).Model(&Download{}).Where("user_id = ?", userID)
	if state != nil {
		db = db.Where("state = ?", *state)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := db.Order("created_at DESC").Limit(size).Offset(offset).Find(&downloads).Error

	return downloads, total, err
}

// GetDownloadByID retrieves a specific download by its ID.
func (r *Repository) GetDownloadByID(ctx context.Context, userID, downloadID uuid.UUID) (*Download, error) {
	var download Download
	err := r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, downloadID).First(&download).Error
	return &download, err
}

// UpdateDownload updates a download's information (e.g., state, path, size).
func (r *Repository) UpdateDownload(ctx context.Context, download *Download) error {
	return r.db.WithContext(ctx).Save(download).Error
}

// RemoveDownload removes a download entry.
func (r *Repository) RemoveDownload(ctx context.Context, userID, downloadID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, downloadID).Delete(&Download{}).Error
}