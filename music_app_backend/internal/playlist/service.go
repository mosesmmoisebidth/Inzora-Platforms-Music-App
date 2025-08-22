
package playlist

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mosesmmoisebidth/music_backend/pkg/logger"
	"gorm.io/gorm"
)

var (
	ErrPlaylistNotFound = errors.New("playlist not found")
	ErrTrackNotFound    = errors.New("track not found in playlist")
	ErrNotPlaylistOwner = errors.New("user is not the owner of the playlist")
)

// Service provides playlist business logic.
type Service struct {
	repo   *Repository
	logger logger.Logger
}

// NewService creates a new playlist service.
func NewService(repo *Repository, logger logger.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// CreatePlaylist creates a new playlist for a user.
func (s *Service) CreatePlaylist(ctx context.Context, userIDStr, title, description string) (*Playlist, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	playlist := &Playlist{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       title,
		Description: description,
		IsPublic:    false, // Default to private
	}

	if err := s.repo.Create(ctx, playlist); err != nil {
		s.logger.Error("failed to create playlist", "error", err, "userID", userID)
		return nil, err
	}

	return playlist, nil
}

// GetPlaylist retrieves a single playlist, checking for ownership or public status.
func (s *Service) GetPlaylist(ctx context.Context, playlistIDStr, userIDStr string) (*Playlist, error) {
	playlistID, err := uuid.Parse(playlistIDStr)
	if err != nil {
		return nil, errors.New("invalid playlist ID format")
	}

	playlist, err := s.repo.GetByID(ctx, playlistID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlaylistNotFound
		}
		s.logger.Error("failed to get playlist by id", "error", err, "playlistID", playlistID)
		return nil, err
	}

	if !playlist.IsPublic && userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil || playlist.UserID != userID {
			return nil, ErrNotPlaylistOwner
		}
	} else if !playlist.IsPublic && userIDStr == "" {
		return nil, ErrNotPlaylistOwner
	}

	return playlist, nil
}

// GetUserPlaylists retrieves all playlists for a given user.
func (s *Service) GetUserPlaylists(ctx context.Context, userIDStr string, page, size int) ([]Playlist, int64, error) {
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

	playlists, total, err := s.repo.GetUserPlaylists(ctx, userID, page, size)
	if err != nil {
		s.logger.Error("failed to get user playlists", "error", err, "userID", userID)
		return nil, 0, err
	}

	return playlists, total, nil
}

// UpdatePlaylist updates a playlist's details, checking for ownership.
func (s *Service) UpdatePlaylist(ctx context.Context, playlistIDStr, userIDStr string, title, description *string, isPublic *bool) (*Playlist, error) {
	playlist, err := s.getAndVerifyOwner(ctx, playlistIDStr, userIDStr)
	if err != nil {
		return nil, err
	}

	if title != nil {
		playlist.Title = *title
	}
	if description != nil {
		playlist.Description = *description
	}
	if isPublic != nil {
		playlist.IsPublic = *isPublic
	}

	if err := s.repo.Update(ctx, playlist); err != nil {
		s.logger.Error("failed to update playlist", "error", err, "playlistID", playlist.ID)
		return nil, err
	}

	return playlist, nil
}

// DeletePlaylist deletes a playlist, checking for ownership.
func (s *Service) DeletePlaylist(ctx context.Context, playlistIDStr, userIDStr string) error {
	playlist, err := s.getAndVerifyOwner(ctx, playlistIDStr, userIDStr)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, playlist.ID); err != nil {
		s.logger.Error("failed to delete playlist", "error", err, "playlistID", playlist.ID)
		return err
	}

	return nil
}

// AddTrackToPlaylist adds a track to a playlist, checking for ownership.
func (s *Service) AddTrackToPlaylist(ctx context.Context, playlistIDStr, userIDStr string, data TrackData) (*Playlist, error) {
	playlist, err := s.getAndVerifyOwner(ctx, playlistIDStr, userIDStr)
	if err != nil {
		return nil, err
	}

	maxPos, err := s.repo.GetMaxPosition(ctx, playlist.ID)
	if err != nil {
		s.logger.Error("failed to get max position for track", "error", err, "playlistID", playlist.ID)
		return nil, err
	}

	newTrack := &PlaylistTrack{
		ID:              uuid.New(),
		PlaylistID:      playlist.ID,
		Provider:        data.Provider,
		ProviderTrackID: data.ProviderTrackID,
		Title:           data.Title,
		Artist:          data.Artist,
		Album:           data.Album,
		DurationMs:      data.DurationMs,
		ArtworkURL:      data.ArtworkURL,
		Position:        maxPos + 1,
		AddedAt:         time.Now(),
	}

	if err := s.repo.AddTrack(ctx, newTrack); err != nil {
		s.logger.Error("failed to add track to playlist", "error", err, "playlistID", playlist.ID)
		return nil, err
	}

	return s.repo.GetByID(ctx, playlist.ID)
}

// RemoveTrackFromPlaylist removes a track from a playlist, checking for ownership.
func (s *Service) RemoveTrackFromPlaylist(ctx context.Context, playlistIDStr, userIDStr, trackIDStr string) (*Playlist, error) {
	playlist, err := s.getAndVerifyOwner(ctx, playlistIDStr, userIDStr)
	if err != nil {
		return nil, err
	}

	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		return nil, errors.New("invalid track ID format")
	}

	if err := s.repo.RemoveTrack(ctx, playlist.ID, trackID); err != nil {
		s.logger.Error("failed to remove track from playlist", "error", err, "playlistID", playlist.ID, "trackID", trackID)
		return nil, err
	}

	return s.repo.GetByID(ctx, playlist.ID)
}

// ReorderPlaylistTracks reorders tracks within a playlist.
func (s *Service) ReorderPlaylistTracks(ctx context.Context, playlistIDStr, userIDStr string, trackPositions []TrackPosition) error {
	_, err := s.getAndVerifyOwner(ctx, playlistIDStr, userIDStr)
	if err != nil {
		return err
	}

	// This is a complex operation and would require a transaction and careful updates.
	// The simplified version is removed to avoid introducing buggy behavior.
	// A full implementation would fetch all tracks, reorder them in memory, and then update all positions in a single transaction.
	s.logger.Warn("ReorderPlaylistTracks is not fully implemented and is a complex operation.")
	return errors.New("reordering tracks is not implemented yet")
}

// getAndVerifyOwner is a helper function to get a playlist and check if the user is the owner.
func (s *Service) getAndVerifyOwner(ctx context.Context, playlistIDStr, userIDStr string) (*Playlist, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	playlistID, err := uuid.Parse(playlistIDStr)
	if err != nil {
		return nil, errors.New("invalid playlist ID format")
	}

	playlist, err := s.repo.GetByID(ctx, playlistID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlaylistNotFound
		}
		s.logger.Error("failed to get playlist by id for ownership check", "error", err, "playlistID", playlistID)
		return nil, err
	}

	if playlist.UserID != userID {
		return nil, ErrNotPlaylistOwner
	}

	return playlist, nil
}
