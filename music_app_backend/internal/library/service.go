
package library

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mosesmmoisebidth/music_backend/pkg/logger"
	"gorm.io/gorm"
)

var (
	ErrFavoriteExists   = errors.New("track is already in favorites")
	ErrFavoriteNotFound = errors.New("favorite not found")
	ErrDownloadNotFound = errors.New("download not found")
)

// Service provides library business logic.
type Service struct {
	repo   *Repository
	logger logger.Logger
}

// NewService creates a new library service.
func NewService(repo *Repository, logger logger.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// --- Favorites ---

// AddFavorite adds a track to a user's favorites, preventing duplicates.
func (s *Service) AddFavorite(ctx context.Context, userIDStr string, data TrackData) (*Favorite, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	_, err = s.repo.FindFavoriteByTrackID(ctx, userID, data.Provider, data.ProviderTrackID)
	if err == nil {
		return nil, ErrFavoriteExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("failed to check for existing favorite", "error", err, "userID", userID)
		return nil, err
	}

	favorite := &Favorite{
		ID:              uuid.New(),
		UserID:          userID,
		Provider:        data.Provider,
		ProviderTrackID: data.ProviderTrackID,
		Title:           data.Title,
		Artist:          data.Artist,
		Album:           data.Album,
		DurationMs:      data.DurationMs,
		ArtworkURL:      data.ArtworkURL,
		AddedAt:         time.Now(),
	}

	if err := s.repo.AddFavorite(ctx, favorite); err != nil {
		s.logger.Error("failed to add favorite", "error", err, "userID", userID)
		return nil, err
	}

	return favorite, nil
}

// GetFavorites retrieves a paginated list of a user's favorite tracks.
func (s *Service) GetFavorites(ctx context.Context, userIDStr string, page, size int) ([]Favorite, int64, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 0, errors.New("invalid user ID format")
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	favorites, total, err := s.repo.GetFavorites(ctx, userID, page, size)
	if err != nil {
		s.logger.Error("failed to get favorites", "error", err, "userID", userID)
		return nil, 0, err
	}

	return favorites, total, nil
}

// RemoveFavorite removes a track from a user's favorites.
func (s *Service) RemoveFavorite(ctx context.Context, userIDStr, favoriteIDStr string) error {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.New("invalid user ID format")
	}
	favoriteID, err := uuid.Parse(favoriteIDStr)
	if err != nil {
		return errors.New("invalid favorite ID format")
	}

	if err := s.repo.RemoveFavorite(ctx, userID, favoriteID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFavoriteNotFound
		}
		s.logger.Error("failed to remove favorite", "error", err, "userID", userID, "favoriteID", favoriteID)
		return err
	}

	return nil
}

// --- History ---

// AddHistory adds a track to a user's listening history.
func (s *Service) AddHistory(ctx context.Context, userIDStr string, data TrackData) (*History, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	history := &History{
		ID:              uuid.New(),
		UserID:          userID,
		Provider:        data.Provider,
		ProviderTrackID: data.ProviderTrackID,
		Title:           data.Title,
		Artist:          data.Artist,
		Album:           data.Album,
		DurationMs:      data.DurationMs,
		ArtworkURL:      data.ArtworkURL,
		PlayedAt:        time.Now(),
	}

	if err := s.repo.AddHistory(ctx, history); err != nil {
		s.logger.Error("failed to add history", "error", err, "userID", userID)
		return nil, err
	}

	return history, nil
}

// GetUserHistory retrieves a paginated list of a user's listening history.
func (s *Service) GetUserHistory(ctx context.Context, userIDStr string, page, size int) ([]History, int64, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 0, errors.New("invalid user ID format")
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 50
	}

	history, total, err := s.repo.GetHistory(ctx, userID, page, size)
	if err != nil {
		s.logger.Error("failed to get user history", "error", err, "userID", userID)
		return nil, 0, err
	}

	return history, total, nil
}

// --- Downloads ---

// AddDownload adds a track to the user's download list with a 'Pending' state.
func (s *Service) AddDownload(ctx context.Context, userIDStr string, data TrackData, quality string) (*Download, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	download := &Download{
		ID:              uuid.New(),
		UserID:          userID,
		Provider:        data.Provider,
		ProviderTrackID: data.ProviderTrackID,
		Title:           data.Title,
		Artist:          data.Artist,
		Album:           data.Album,
		DurationMs:      data.DurationMs,
		ArtworkURL:      data.ArtworkURL,
		State:           StatePending,
		Quality:         quality,
	}

	if err := s.repo.AddDownload(ctx, download); err != nil {
		s.logger.Error("failed to add download", "error", err, "userID", userID)
		return nil, err
	}

	return download, nil
}

// GetUserDownloads retrieves a user's downloads, optionally filtered by state.
func (s *Service) GetUserDownloads(ctx context.Context, userIDStr string, page, size int, state *DownloadState) ([]Download, int64, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 0, errors.New("invalid user ID format")
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	downloads, total, err := s.repo.GetDownloads(ctx, userID, page, size, state)
	if err != nil {
		s.logger.Error("failed to get user downloads", "error", err, "userID", userID)
		return nil, 0, err
	}

	return downloads, total, nil
}

// RemoveDownload removes a download entry.
func (s *Service) RemoveDownload(ctx context.Context, userIDStr, downloadIDStr string) error {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.New("invalid user ID format")
	}
	downloadID, err := uuid.Parse(downloadIDStr)
	if err != nil {
		return errors.New("invalid download ID format")
	}

	if err := s.repo.RemoveDownload(ctx, userID, downloadID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDownloadNotFound
		}
		s.logger.Error("failed to remove download", "error", err, "userID", userID, "downloadID", downloadID)
		return err
	}

	return nil
}
