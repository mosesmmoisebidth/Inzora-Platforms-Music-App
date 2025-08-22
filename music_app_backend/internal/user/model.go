
package user

import (
	"time"
	"github.com/lib/pq"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// User represents a user in the system.
type User struct {
	ID             uuid.UUID              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Email          *string                `gorm:"uniqueIndex;size:255" json:"email"`
	Password       *string                `gorm:"size:255" json:"-"`
	DisplayName    *string                `gorm:"size:255" json:"display_name"`
	PhotoURL       *string                `gorm:"size:1024" json:"photo_url"`
	Roles          pq.StringArray  `gorm:"type:text[]" json:"roles"`
	IsActive       bool                   `gorm:"default:true" json:"is_active"`
	GoogleID       *string                `gorm:"uniqueIndex;size:255" json:"google_id,omitempty"`
	LastLoginAt    *time.Time             `json:"last_login_at,omitempty"`
	Preferences    datatypes.JSON         `gorm:"type:jsonb" json:"preferences,omitempty"`
	FavoriteGenres pq.StringArray  `gorm:"type:text[]" json:"favorite_genres,omitempty"`
	CreatedAt      time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
