package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mosesmmoisebidth/music_backend/internal/library"
	"github.com/mosesmmoisebidth/music_backend/pkg/logger"
	"github.com/mosesmmoisebidth/music_backend/pkg/response"
)

// LibraryHandlers contains library HTTP handlers
type LibraryHandlers struct {
	service *library.Service
	logger  logger.Logger
}

// NewLibraryHandlers creates new library handlers
func NewLibraryHandlers(service *library.Service, logger logger.Logger) *LibraryHandlers {
	return &LibraryHandlers{service: service, logger: logger}
}

// GetFavorites retrieves the user's favorite tracks.
// @Summary      Get favorite tracks
// @Description  Retrieves a paginated list of the authenticated user's favorite tracks.
// @Tags         Library
// @Produce      json
// @Security     Bearer
// @Param        page query int false "Page number" default(1)
// @Param        size query int false "Page size" default(20)
// @Success      200 {object} response.APIResponse{data=response.PaginatedData{favorites=[]FavoriteResponse}}
// @Failure      401 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /favorites [get]
func (h *LibraryHandlers) GetFavorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_FOUND", "User not authenticated")
		return
	}

	var req GetFavoritesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	favorites, total, err := h.service.GetFavorites(c.Request.Context(), userID.(string), req.Page, req.Size)
	if err != nil {
		h.logger.Error("failed to get favorites", "error", err, "user_id", userID)
		response.InternalError(c, "FAVORITES_FETCH_FAILED", "Failed to fetch favorites")
		return
	}

	var favResponses []FavoriteResponse
	for _, f := range favorites {
		favResponses = append(favResponses, mapFavoriteToResponse(&f))
	}

	response.Success(c, response.NewPaginatedData(favResponses, req.Page, req.Size, total))
}

// AddFavorite adds a track to the user's favorites.
// @Summary      Add a favorite track
// @Description  Adds a track to the authenticated user's list of favorite tracks.
// @Tags         Library
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body AddFavoriteRequest true "Favorite track information"
// @Success      201 {object} response.APIResponse{data=FavoriteResponse}
// @Failure      400 {object} response.APIResponse{error=response.APIError}
// @Failure      401 {object} response.APIResponse{error=response.APIError}
// @Failure      409 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /favorites [post]
func (h *LibraryHandlers) AddFavorite(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_FOUND", "User not authenticated")
		return
	}

	var req AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	trackData := library.TrackData{
		Provider:        req.Provider,
		ProviderTrackID: req.ProviderTrackID,
		Title:           req.Title,
		Artist:          req.Artist,
		Album:           req.Album,
		DurationMs:      req.DurationMs,
		ArtworkURL:      req.ArtworkURL,
	}

	favorite, err := h.service.AddFavorite(c.Request.Context(), userID.(string), trackData)
	if err != nil {
		if err == library.ErrFavoriteExists {
			response.Conflict(c, "FAVORITE_EXISTS", err.Error())
			return
		}
		h.logger.Error("failed to add favorite", "error", err, "user_id", userID)
		response.InternalError(c, "FAVORITE_ADD_FAILED", "Failed to add favorite")
		return
	}

	response.Created(c, mapFavoriteToResponse(favorite))
}

// RemoveFavorite removes a track from the user's favorites.
// @Summary      Remove a favorite track
// @Description  Removes a track from the authenticated user's list of favorite tracks.
// @Tags         Library
// @Produce      json
// @Security     Bearer
// @Param        favoriteId path string true "Favorite ID"
// @Success      200 {object} response.APIResponse{data=response.SuccessMessage}
// @Failure      401 {object} response.APIResponse{error=response.APIError}
// @Failure      404 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /favorites/{favoriteId} [delete]
func (h *LibraryHandlers) RemoveFavorite(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_FOUND", "User not authenticated")
		return
	}

	favoriteID := c.Param("favoriteId")

	if err := h.service.RemoveFavorite(c.Request.Context(), userID.(string), favoriteID); err != nil {
		if err == library.ErrFavoriteNotFound {
			response.NotFound(c, "FAVORITE_NOT_FOUND", err.Error())
			return
		}
		h.logger.Error("failed to remove favorite", "error", err, "favorite_id", favoriteID)
		response.InternalError(c, "FAVORITE_REMOVE_FAILED", "Failed to remove favorite")
		return
	}

	response.Success(c, &response.SuccessMessage{Message: "Removed from favorites successfully"})
}

// GetHistory retrieves the user's listening history.
// @Summary      Get listening history
// @Description  Retrieves a paginated list of the authenticated user's recently played tracks.
// @Tags         Library
// @Produce      json
// @Security     Bearer
// @Param        page query int false "Page number" default(1)
// @Param        size query int false "Page size" default(50)
// @Success      200 {object} response.APIResponse{data=response.PaginatedData{history=[]HistoryResponse}}
// @Failure      401 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /history [get]
func (h *LibraryHandlers) GetHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_FOUND", "User not authenticated")
		return
	}

	var req GetHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	history, total, err := h.service.GetUserHistory(c.Request.Context(), userID.(string), req.Page, req.Size)
	if err != nil {
		h.logger.Error("failed to get history", "error", err, "user_id", userID)
		response.InternalError(c, "HISTORY_FETCH_FAILED", "Failed to fetch history")
		return
	}

	var histResponses []HistoryResponse
	for _, h := range history {
		histResponses = append(histResponses, mapHistoryToResponse(&h))
	}

	response.Success(c, response.NewPaginatedData(histResponses, req.Page, req.Size, total))
}

// AddHistory adds a track to the user's listening history.
// @Summary      Add to listening history
// @Description  Adds a track to the authenticated user's listening history.
// @Tags         Library
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body AddHistoryRequest true "History track information"
// @Success      201 {object} response.APIResponse{data=HistoryResponse}
// @Failure      400 {object} response.APIResponse{error=response.APIError}
// @Failure      401 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /history [post]
func (h *LibraryHandlers) AddHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_FOUND", "User not authenticated")
		return
	}

	var req AddHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	trackData := library.TrackData{
		Provider:        req.Provider,
		ProviderTrackID: req.ProviderTrackID,
		Title:           req.Title,
		Artist:          req.Artist,
		Album:           req.Album,
		DurationMs:      req.DurationMs,
		ArtworkURL:      req.ArtworkURL,
	}

	history, err := h.service.AddHistory(c.Request.Context(), userID.(string), trackData)
	if err != nil {
		h.logger.Error("failed to add history", "error", err, "user_id", userID)
		response.InternalError(c, "HISTORY_ADD_FAILED", "Failed to add to history")
		return
	}

	response.Created(c, mapHistoryToResponse(history))
}