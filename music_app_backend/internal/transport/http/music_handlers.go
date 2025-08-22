package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mosesmmoisebidth/music_backend/internal/music"
	"github.com/mosesmmoisebidth/music_backend/pkg/logger"
	"github.com/mosesmmoisebidth/music_backend/pkg/response"
)

// MusicHandlers contains music discovery HTTP handlers
type MusicHandlers struct {
	service *music.MusicService
	logger  logger.Logger
}

// NewMusicHandlers creates new music handlers
func NewMusicHandlers(service *music.MusicService, logger logger.Logger) *MusicHandlers {
	return &MusicHandlers{service: service, logger: logger}
}

// SearchTracks searches for tracks across music providers.
// @Summary      Search for tracks
// @Description  Searches for tracks by a query string, with optional filters.
// @Tags         Music
// @Produce      json
// @Param        q query string true "Search query"
// @Param        provider query string false "Provider to search (e.g., itunes, spotify)"
// @Param        page query int false "Page number" default(1)
// @Param        size query int false "Page size" default(20)
// @Param        genre query string false "Filter by genre"
// @Param        year query int false "Filter by release year"
// @Param        explicit query bool false "Filter by explicit content"
// @Success      200 {object} response.APIResponse{data=response.PaginatedData{tracks=[]TrackResponse}}
// @Failure      400 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /music/search [get]
func (h *MusicHandlers) SearchTracks(c *gin.Context) {
	var req SearchTracksRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	filters := &music.SearchFilters{}
	if req.Genre != nil {
		filters.Genre = *req.Genre
	}
	if req.Year != nil {
		filters.Year = strconv.Itoa(*req.Year)
	}
	if req.Explicit != nil {
		filters.Explicit = req.Explicit
	}

	provider := ""
	if req.Provider != nil {
		provider = *req.Provider
	}

	tracks, pageInfo, err := h.service.SearchTracks(c.Request.Context(), provider, req.Query, req.Page, req.Size, filters)
	if err != nil {
		h.logger.Error("failed to search tracks", "error", err, "query", req.Query)
		response.InternalError(c, "SEARCH_FAILED", "Failed to search for tracks")
		return
	}

	var trackResponses []TrackResponse
	for _, t := range tracks {
		trackResponses = append(trackResponses, mapTrackToResponse(&t))
	}

	response.Success(c, response.NewPaginatedData(trackResponses, pageInfo.Page, pageInfo.Size, pageInfo.Total))
}

// GetTrack retrieves a single track by its ID from a specific provider.
// @Summary      Get a single track
// @Description  Retrieves full details for a single track by its ID and provider.
// @Tags         Music
// @Produce      json
// @Param        trackId path string true "Track ID"
// @Param        provider query string true "Provider name (e.g., itunes, spotify)"
// @Success      200 {object} response.APIResponse{data=TrackResponse}
// @Failure      400 {object} response.APIResponse{error=response.APIError}
// @Failure      404 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /music/tracks/{trackId} [get]
func (h *MusicHandlers) GetTrack(c *gin.Context) {
	var req GetTrackRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	track, err := h.service.GetTrack(c.Request.Context(), req.Provider, req.TrackID)
	if err != nil {
		h.logger.Error("failed to get track", "error", err, "track_id", req.TrackID)
		response.NotFound(c, "TRACK_NOT_FOUND", "Track not found")
		return
	}

	response.Success(c, mapTrackToResponse(track))
}

// GetTopCharts retrieves top charts for a given country.
// @Summary      Get top charts
// @Description  Retrieves a paginated list of top tracks for a specific country.
// @Tags         Music
// @Produce      json
// @Param        country query string false "Country code (ISO 3166-1 alpha-2)" default(US)
// @Param        provider query string false "Provider name (e.g., itunes, spotify)" default(itunes)
// @Param        page query int false "Page number" default(1)
// @Param        size query int false "Page size" default(20)
// @Success      200 {object} response.APIResponse{data=response.PaginatedData{tracks=[]TrackResponse}}
// @Failure      400 {object} response.APIResponse{error=response.APIError}
// @Failure      500 {object} response.APIResponse{error=response.APIError}
// @Router       /music/top-charts [get]
func (h *MusicHandlers) GetTopCharts(c *gin.Context) {
	var req GetTopChartsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	provider := "itunes"
	if req.Provider != nil {
		provider = *req.Provider
	}

	tracks, pageInfo, err := h.service.GetTopCharts(c.Request.Context(), provider, req.Country, req.Page, req.Size)
	if err != nil {
		h.logger.Error("failed to get top charts", "error", err, "provider", provider)
		response.InternalError(c, "CHARTS_FETCH_FAILED", "Failed to fetch top charts")
		return
	}

	var trackResponses []TrackResponse
	for _, t := range tracks {
		trackResponses = append(trackResponses, mapTrackToResponse(&t))
	}

	response.Success(c, response.NewPaginatedData(trackResponses, pageInfo.Page, pageInfo.Size, pageInfo.Total))
}