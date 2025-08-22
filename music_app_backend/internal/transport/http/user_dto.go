package http

import (
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/mosesmmoisebidth/music_backend/internal/user"
)

// --- User Requests ---

type UpdateUserRequest struct {
	DisplayName *string                `json:"display_name,omitempty" binding:"omitempty,min=2,max=50"`
	PhotoURL    *string                `json:"photo_url,omitempty" binding:"omitempty,url"`
	Preferences map[string]interface{} `json:"preferences,omitempty"`
}

// --- User Responses ---

type UserResponse struct {
	ID          uuid.UUID              `json:"id"`
	Email       *string                `json:"email,omitempty"`
	DisplayName *string                `json:"display_name,omitempty"`
	PhotoURL    *string                `json:"photo_url,omitempty"`
	Roles       []string               `json:"roles"`
	IsActive    bool                   `json:"is_active"`
	Preferences map[string]interface{} `json:"preferences,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

func mapUserToResponse(u *user.User) UserResponse {
	prefs := make(map[string]interface{})
	if u.Preferences != nil {
		_ = json.Unmarshal(u.Preferences, &prefs) // Ignore error for DTO mapping
	}

	return UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		PhotoURL:    u.PhotoURL,
		Roles:       u.Roles,
		IsActive:    u.IsActive,
		Preferences: prefs,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}