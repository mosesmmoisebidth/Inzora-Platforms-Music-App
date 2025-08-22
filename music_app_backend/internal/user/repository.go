package user

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// repository implements the Repository interface for user data.
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository instance.
func NewRepository(db *gorm.DB) Repository { // Returns interface
	return &repository{db: db}
}

// CreateUser creates a new user record in the database.
func (r *repository) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetUserByEmail retrieves a user by their email address.
func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

// GetUserByID retrieves a user by their ID.
func (r *repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

// GetUserByGoogleID retrieves a user by their Google ID.
func (r *repository) GetUserByGoogleID(ctx context.Context, googleID string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("google_id = ?", googleID).First(&user).Error
	return &user, err
}

// UpdateUser updates an existing user record.
func (r *repository) UpdateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}