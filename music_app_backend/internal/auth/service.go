package auth

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mosesmmoisebidth/music_backend/pkg/logger"
	"gorm.io/gorm"
)

// RefreshTokenRepository defines the refresh token repository interface
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	GetByTokenID(ctx context.Context, tokenID string) (*RefreshToken, error)
	RevokeByTokenID(ctx context.Context, tokenID string) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
	CleanExpiredTokens(ctx context.Context) error
}

// AuthService provides authentication business logic
type AuthService struct {
	jwtService      *JWTService
	googleService   *GoogleService
	refreshTokenRepo RefreshTokenRepository
	logger          logger.Logger
}

// NewAuthService creates a new authentication service
func NewAuthService(
	jwtService *JWTService,
	googleService *GoogleService,
	refreshTokenRepo RefreshTokenRepository,
	logger logger.Logger,
) *AuthService {
	return &AuthService{
		jwtService:      jwtService,
		googleService:   googleService,
		refreshTokenRepo: refreshTokenRepo,
		logger:          logger,
	}
}

// GenerateTokens generates access and refresh tokens for a user
func (s *AuthService) GenerateTokens(ctx context.Context, userID uuid.UUID, email string, roles []string, userAgent, ip string) (*TokenPair, error) {
	// Generate JWT token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(userID, email, roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	// Parse refresh token to get claims
	refreshClaims, err := s.jwtService.VerifyRefreshToken(tokenPair.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated refresh token: %w", err)
	}

	// Store refresh token in database
	refreshToken := &RefreshToken{
		UserID:    userID,
		TokenID:   refreshClaims.ID,
		IssuedAt:  refreshClaims.IssuedAt.Time,
		ExpiresAt: refreshClaims.ExpiresAt.Time,
		UserAgent: &userAgent,
		IP:        &ip,
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return tokenPair, nil
}

// RefreshTokens refreshes access and refresh tokens
func (s *AuthService) RefreshTokens(ctx context.Context, refreshTokenString, userAgent, ip string) (*TokenPair, error) {
	// Verify refresh token
	refreshClaims, err := s.jwtService.VerifyRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Get stored refresh token
	storedToken, err := s.refreshTokenRepo.GetByTokenID(ctx, refreshClaims.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	// Check if token is valid
	if !storedToken.IsValid() {
		return nil, fmt.Errorf("refresh token is revoked or expired")
	}

	// Generate new token pair
	newTokenPair, err := s.jwtService.GenerateTokenPair(refreshClaims.UserID, refreshClaims.Email, refreshClaims.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new token pair: %w", err)
	}

	// Parse new refresh token
	newRefreshClaims, err := s.jwtService.VerifyRefreshToken(newTokenPair.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse new refresh token: %w", err)
	}

	// Revoke old token and create new one
	storedToken.RevokeAndReplace(newRefreshClaims.ID)
	
	// Create new refresh token record
	newRefreshToken := &RefreshToken{
		UserID:    refreshClaims.UserID,
		TokenID:   newRefreshClaims.ID,
		IssuedAt:  newRefreshClaims.IssuedAt.Time,
		ExpiresAt: newRefreshClaims.ExpiresAt.Time,
		UserAgent: &userAgent,
		IP:        &ip,
	}

	if err := s.refreshTokenRepo.Create(ctx, newRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	s.logger.Info("Tokens refreshed successfully", "user_id", refreshClaims.UserID)
	return newTokenPair, nil
}

// RevokeToken revokes a refresh token (logout)
func (s *AuthService) RevokeToken(ctx context.Context, refreshTokenString string) error {
	// Verify refresh token
	refreshClaims, err := s.jwtService.VerifyRefreshToken(refreshTokenString)
	if err != nil {
		return fmt.Errorf("invalid refresh token")
	}

	// Revoke the token
	if err := s.refreshTokenRepo.RevokeByTokenID(ctx, refreshClaims.ID); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	s.logger.Info("Token revoked successfully", "user_id", refreshClaims.UserID, "token_id", refreshClaims.ID)
	return nil
}

// RevokeAllUserTokens revokes all tokens for a user
func (s *AuthService) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	if err := s.refreshTokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke user tokens: %w", err)
	}

	s.logger.Info("All user tokens revoked", "user_id", userID)
	return nil
}

// VerifyGoogleIDToken verifies a Google ID token and returns user info
func (s *AuthService) VerifyGoogleIDToken(ctx context.Context, idToken string) (*GoogleUser, error) {
	return s.googleService.VerifyIDToken(ctx, idToken)
}

// CleanupExpiredTokens removes expired refresh tokens
func (s *AuthService) CleanupExpiredTokens(ctx context.Context) error {
	return s.refreshTokenRepo.CleanExpiredTokens(ctx)
}
