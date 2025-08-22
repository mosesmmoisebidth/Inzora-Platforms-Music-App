package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mosesmmoisebidth/music_backend/internal/playlist"
	"github.com/mosesmmoisebidth/music_backend/pkg/logger"
	"github.com/mosesmmoisebidth/music_backend/pkg/response"
)

// PlaylistHandlers contains playlist HTTP handlers
type PlaylistHandlers struct {
	service *playlist.Service
	logger  logger.Logger
}

// NewPlaylistHandlers creates new playlist handlers
func NewPlaylistHandlers(service *playlist.Service, logger logger.Logger) *PlaylistHandlers {
	return &PlaylistHandlers{service: service, logger: logger}
}

// GetPlaylists retrieves the user's playlists.
// @Summary      Get user playlists
// @Description  Retrieves a paginated list of the authenticated user's playlists.
// @Tags         Playlists
// @Produce      json
// @Security     Bearer
// @Param        page query int false "Page number" default(1)
// @Param        size query int false "Page size" default(20)
// @Success      200 {object} response.APIResponse{data=response.PaginatedData{playlists=[]PlaylistResponse}}
// @Failure      401 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /playlists [get]
func (h *PlaylistHandlers) GetPlaylists(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_FOUND", "User not authenticated")
		return
	}

	var req GetPlaylistsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	playlists, total, err := h.service.GetUserPlaylists(c.Request.Context(), userID.(string), req.Page, req.Size)
	if err != nil {
		h.logger.Error("failed to get playlists", "error", err, "user_id", userID)
		response.InternalError(c, "PLAYLISTS_FETCH_FAILED", "Failed to fetch playlists")
		return
	}

	var playlistResponses []PlaylistResponse
	for _, p := range playlists {
		playlistResponses = append(playlistResponses, mapPlaylistToResponse(&p))
	}

	response.Success(c, response.NewPaginatedData(playlistResponses, req.Page, req.Size, total))
}

// CreatePlaylist creates a new playlist.
// @Summary      Create a new playlist
// @Description  Creates a new playlist for the authenticated user.
// @Tags         Playlists
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body CreatePlaylistRequest true "Playlist details"
// @Success      201 {object} response.APIResponse{data=PlaylistResponse}
// @Failure      400 {object} response.APIResponse{error=response.APIError}
// @Failure      401 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /playlists [post]
func (h *PlaylistHandlers) CreatePlaylist(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_FOUND", "User not authenticated")
		return
	}

	var req CreatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	desc := ""
	if req.Description != nil {
		desc = *req.Description
	}

	newPlaylist, err := h.service.CreatePlaylist(c.Request.Context(), userID.(string), req.Title, desc)
	if err != nil {
		h.logger.Error("failed to create playlist", "error", err, "user_id", userID)
		response.InternalError(c, "PLAYLIST_CREATE_FAILED", "Failed to create playlist")
		return
	}

	response.Created(c, mapPlaylistToResponse(newPlaylist))
}

// GetPlaylist retrieves a single playlist by its ID.
// @Summary      Get a single playlist
// @Description  Retrieves details for a single playlist. Can be accessed without authentication if the playlist is public.
// @Tags         Playlists
// @Produce      json
// @Security     Bearer
// @Param        playlistId path string true "Playlist ID"
// @Success      200 {object} response.APIResponse{data=PlaylistResponse}
// @Failure      403 {object} response.APIResponse{error=response.APIError}
// @Failure      404 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /playlists/{playlistId} [get]
func (h *PlaylistHandlers) GetPlaylist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	playlistID := c.Param("playlistId")

	playlistData, err := h.service.GetPlaylist(c.Request.Context(), playlistID, userIDStr)
	if err != nil {
		switch err {
		case playlist.ErrPlaylistNotFound:
			response.NotFound(c, "PLAYLIST_NOT_FOUND", err.Error())
		case playlist.ErrNotPlaylistOwner:
			response.Forbidden(c, "FORBIDDEN", err.Error())
		default:
			h.logger.Error("failed to get playlist", "error", err, "playlist_id", playlistID)
			response.InternalError(c, "PLAYLIST_FETCH_FAILED", "Failed to fetch playlist")
		}
		return
	}

	response.Success(c, mapPlaylistToResponse(playlistData))
}

// UpdatePlaylist updates a playlist's details.
// @Summary      Update a playlist
// @Description  Updates a playlist's title, description, or public status.
// @Tags         Playlists
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        playlistId path string true "Playlist ID"
// @Param        request body UpdatePlaylistRequest true "Fields to update"
// @Success      200 {object} response.APIResponse{data=PlaylistResponse}
// @Failure      400 {object} response.APIResponse{error=response.APIError}
// @Failure      401 {object} response.APIResponse{error=response.APIError}
// @Failure      403 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /playlists/{playlistId} [patch]
func (h *PlaylistHandlers) UpdatePlaylist(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_FOUND", "User not authenticated")
		return
	}

	playlistID := c.Param("playlistId")

	var req UpdatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	updatedPlaylist, err := h.service.UpdatePlaylist(c.Request.Context(), playlistID, userID.(string), req.Title, req.Description, req.IsPublic)
	if err != nil {
		switch err {
		case playlist.ErrPlaylistNotFound, playlist.ErrNotPlaylistOwner:
			response.Forbidden(c, "FORBIDDEN", "You do not have permission to update this playlist")
		default:
			h.logger.Error("failed to update playlist", "error", err, "playlist_id", playlistID)
			response.InternalError(c, "PLAYLIST_UPDATE_FAILED", "Failed to update playlist")
		}
		return
	}

	response.Success(c, mapPlaylistToResponse(updatedPlaylist))
}

// DeletePlaylist deletes a playlist.
// @Summary      Delete a playlist
// @Description  Permanently deletes a user's playlist.
// @Tags         Playlists
// @Produce      json
// @Security     Bearer
// @Param        playlistId path string true "Playlist ID"
// @Success      200 {object} response.APIResponse{data=response.SuccessMessage}
// @Failure      401 {object} response.APIResponse{error=response.APIError}
// @Failure      403 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /playlists/{playlistId} [delete]
func (h *PlaylistHandlers) DeletePlaylist(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_FOUND", "User not authenticated")
		return
	}

	playlistID := c.Param("playlistId")

	if err := h.service.DeletePlaylist(c.Request.Context(), playlistID, userID.(string)); err != nil {
		switch err {
		case playlist.ErrPlaylistNotFound, playlist.ErrNotPlaylistOwner:
			response.Forbidden(c, "FORBIDDEN", "You do not have permission to delete this playlist")
		default:
			h.logger.Error("failed to delete playlist", "error", err, "playlist_id", playlistID)
			response.InternalError(c, "PLAYLIST_DELETE_FAILED", "Failed to delete playlist")
		}
		return
	}

	response.Success(c, &response.SuccessMessage{Message: "Playlist deleted successfully"})
}

// AddTrackToPlaylist adds a track to a playlist.
// @Summary      Add track to playlist
// @Description  Adds a single track to the end of a specified playlist.
// @Tags         Playlists
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        playlistId path string true "Playlist ID"
// @Param        request body AddTrackToPlaylistRequest true "Track information"
// @Success      200 {object} response.APIResponse{data=PlaylistResponse}
// @Failure      400 {object} response.APIResponse{error=response.APIError}
// @Failure      401 {object} response.APIResponse{error=response.APIError}
// @Failure      403 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /playlists/{playlistId}/tracks [post]
func (h *PlaylistHandlers) AddTrackToPlaylist(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_FOUND", "User not authenticated")
		return
	}

	playlistID := c.Param("playlistId")

	var req AddTrackToPlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	trackData := playlist.TrackData{
		Provider:        playlist.MusicProvider(req.Provider),
		ProviderTrackID: req.ProviderTrackID,
		Title:           req.Title,
		Artist:          req.Artist,
		Album:           req.Album,
		DurationMs:      req.DurationMs,
		ArtworkURL:      req.ArtworkURL,
	}

	updatedPlaylist, err := h.service.AddTrackToPlaylist(c.Request.Context(), playlistID, userID.(string), trackData)
	if err != nil {
		switch err {
		case playlist.ErrPlaylistNotFound, playlist.ErrNotPlaylistOwner:
			response.Forbidden(c, "FORBIDDEN", "You do not have permission to modify this playlist")
		default:
			h.logger.Error("failed to add track to playlist", "error", err, "playlist_id", playlistID)
			response.InternalError(c, "TRACK_ADD_FAILED", "Failed to add track to playlist")
		}
		return
	}

	response.Success(c, mapPlaylistToResponse(updatedPlaylist))
}

// RemoveTrackFromPlaylist removes a track from a playlist.
// @Summary      Remove track from playlist
// @Description  Removes a single track from a specified playlist.
// @Tags         Playlists
// @Produce      json
// @Security     Bearer
// @Param        playlistId path string true "Playlist ID"
// @Param        trackId path string true "Playlist Track ID"
// @Success      200 {object} response.APIResponse{data=PlaylistResponse}
// @Failure      401 {object} response.APIResponse{error=response.APIError}
// @Failure      403 {object} response.APIResponse{error=response.APIError}
// @Failure      404 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /playlists/{playlistId}/tracks/{trackId} [delete]
func (h *PlaylistHandlers) RemoveTrackFromPlaylist(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_FOUND", "User not authenticated")
		return
	}

	playlistID := c.Param("playlistId")
	trackID := c.Param("trackId")

	updatedPlaylist, err := h.service.RemoveTrackFromPlaylist(c.Request.Context(), playlistID, userID.(string), trackID)
	if err != nil {
		switch err {
		case playlist.ErrPlaylistNotFound, playlist.ErrNotPlaylistOwner:
			response.Forbidden(c, "FORBIDDEN", "You do not have permission to modify this playlist")
		case playlist.ErrTrackNotFound:
			response.NotFound(c, "TRACK_NOT_FOUND", err.Error())
		default:
			h.logger.Error("failed to remove track from playlist", "error", err, "playlist_id", playlistID)
			response.InternalError(c, "TRACK_REMOVE_FAILED", "Failed to remove track from playlist")
		}
		return
	}

	response.Success(c, mapPlaylistToResponse(updatedPlaylist))
}