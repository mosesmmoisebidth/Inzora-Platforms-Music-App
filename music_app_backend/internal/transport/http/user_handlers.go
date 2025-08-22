
package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mosesmmoisebidth/music_backend/internal/user"
	"github.com/mosesmmoisebidth/music_backend/pkg/logger"
	"github.com/mosesmmoisebidth/music_backend/pkg/response"
)

// UserHandlers contains user HTTP handlers
type UserHandlers struct {
	service *user.Service
	logger  logger.Logger
}

// NewUserHandlers creates new user handlers
func NewUserHandlers(service *user.Service, logger logger.Logger) *UserHandlers {
	return &UserHandlers{service: service, logger: logger}
}

// GetCurrentUser retrieves the profile of the currently authenticated user.
// @Summary      Get current user profile
// @Description  Retrieves the full profile for the user associated with the JWT token.
// @Tags         Users
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  response.APIResponse{data=UserResponse}
// @Failure      401  {object}  response.APIResponse{error=response.APIError}
// @Failure      500  {object}  response.APIResponse{error=response.APIError}
// @Router       /users/me [get]
func (h *UserHandlers) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_FOUND", "User not authenticated")
		return
	}

	userData, err := h.service.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error("failed to get current user", "error", err, "user_id", userID)
		response.InternalError(c, "USER_FETCH_FAILED", "Failed to fetch user information")
		return
	}

	response.Success(c, mapUserToResponse(userData))
}

// UpdateCurrentUser updates the profile of the currently authenticated user.
// @Summary      Update current user profile
// @Description  Updates the display name, photo URL, or preferences for the current user.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body UpdateUserRequest true "Fields to update"
// @Success      200  {object}  response.APIResponse{data=UserResponse}
// @Failure      400  {object}  response.APIResponse{error=response.APIError}
// @Failure      401  {object}  response.APIResponse{error=response.APIError}
// @Failure      500  {object}  response.APIResponse{error=response.APIError}
// @Router       /users/me [patch]
func (h *UserHandlers) UpdateCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "USER_NOT_FOUND", "User not authenticated")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	updatedUser, err := h.service.UpdateUser(c.Request.Context(), userID.(string), req.DisplayName, req.PhotoURL, req.Preferences)
	if err != nil {
		h.logger.Error("failed to update current user", "error", err, "user_id", userID)
		response.InternalError(c, "USER_UPDATE_FAILED", "Failed to update user information")
		return
	}

	response.Success(c, mapUserToResponse(updatedUser))
}
