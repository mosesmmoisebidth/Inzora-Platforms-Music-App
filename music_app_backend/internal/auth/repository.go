package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// refreshTokenRepository implements RefreshTokenRepository interface
type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *refreshTokenRepository) Create(ctx context.Context, token *RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// GetByTokenID retrieves a refresh token by token ID
func (r *refreshTokenRepository) GetByTokenID(ctx context.Context, tokenID string) (*RefreshToken, error) {
	var token RefreshToken
	err := r.db.WithContext(ctx).Where("token_id = ?", tokenID).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// RevokeByTokenID revokes a refresh token by token ID
func (r *refreshTokenRepository) RevokeByTokenID(ctx context.Context, tokenID string) error {
	return r.db.WithContext(ctx).
		Model(&RefreshToken{}).
		Where("token_id = ?", tokenID).
		Update("revoked", true).Error
}

// RevokeAllUserTokens revokes all tokens for a user
func (r *refreshTokenRepository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}

// CleanExpiredTokens removes expired refresh tokens
func (r *refreshTokenRepository) CleanExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&RefreshToken{}).Error
}
