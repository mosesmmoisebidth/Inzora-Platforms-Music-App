package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshToken represents a refresh token in the database
type RefreshToken struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	TokenID    string     `json:"token_id" gorm:"uniqueIndex;not null;size:255"` // JWT ID (jti claim)
	IssuedAt   time.Time  `json:"issued_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	Revoked    bool       `json:"revoked" gorm:"default:false"`
	ReplacedBy *string    `json:"replaced_by,omitempty" gorm:"size:255"`
	UserAgent  *string    `json:"user_agent,omitempty" gorm:"size:500"`
	IP         *string    `json:"ip,omitempty" gorm:"size:45"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}
	if rt.IssuedAt.IsZero() {
		rt.IssuedAt = time.Now()
	}
	return nil
}

// TableName returns the table name for the RefreshToken model
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsValid checks if the refresh token is valid (not revoked and not expired)
func (rt *RefreshToken) IsValid() bool {
	return !rt.Revoked && rt.ExpiresAt.After(time.Now())
}

// IsExpired checks if the refresh token has expired
func (rt *RefreshToken) IsExpired() bool {
	return rt.ExpiresAt.Before(time.Now())
}

// Revoke marks the refresh token as revoked
func (rt *RefreshToken) Revoke() {
	rt.Revoked = true
	rt.UpdatedAt = time.Now()
}

// RevokeAndReplace marks the refresh token as revoked and sets the replacement
func (rt *RefreshToken) RevokeAndReplace(replacementTokenID string) {
	rt.Revoked = true
	rt.ReplacedBy = &replacementTokenID
	rt.UpdatedAt = time.Now()
}
