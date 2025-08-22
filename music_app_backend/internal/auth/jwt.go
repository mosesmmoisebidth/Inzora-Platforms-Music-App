package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenType represents the type of JWT token
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// Claims represents JWT claims with custom fields
type Claims struct {
	UserID uuid.UUID `json:"uid"`
	Email  string    `json:"email"`
	Roles  []string  `json:"roles"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// JWTService provides JWT token operations
type JWTService struct {
	accessSecret     []byte
	refreshSecret    []byte
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	signingMethod    jwt.SigningMethod
}

// NewJWTService creates a new JWT service
func NewJWTService(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration, method string) *JWTService {
	var signingMethod jwt.SigningMethod
	switch method {
	case "HS256":
		signingMethod = jwt.SigningMethodHS256
	case "HS384":
		signingMethod = jwt.SigningMethodHS384
	case "HS512":
		signingMethod = jwt.SigningMethodHS512
	case "RS256":
		signingMethod = jwt.SigningMethodRS256
	case "RS384":
		signingMethod = jwt.SigningMethodRS384
	case "RS512":
		signingMethod = jwt.SigningMethodRS512
	default:
		signingMethod = jwt.SigningMethodHS256
	}

	return &JWTService{
		accessSecret:     []byte(accessSecret),
		refreshSecret:    []byte(refreshSecret),
		accessTokenTTL:   accessTTL,
		refreshTokenTTL:  refreshTTL,
		signingMethod:    signingMethod,
	}
}

// GenerateTokenPair generates a new access and refresh token pair
func (j *JWTService) GenerateTokenPair(userID uuid.UUID, email string, roles []string) (*TokenPair, error) {
	now := time.Now()
	
	// Generate access token
	accessTokenID := uuid.New().String()
	accessClaims := &Claims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		Type:   TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        accessTokenID,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "music-app-backend",
		},
	}

	accessToken := jwt.NewWithClaims(j.signingMethod, accessClaims)
	accessTokenString, err := accessToken.SignedString(j.accessSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	refreshTokenID := uuid.New().String()
	refreshClaims := &Claims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		Type:   TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        refreshTokenID,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "music-app-backend",
		},
	}

	refreshToken := jwt.NewWithClaims(j.signingMethod, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(j.refreshSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(j.accessTokenTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// VerifyAccessToken validates and parses an access token
func (j *JWTService) VerifyAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != j.signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.accessSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims.Type != TokenTypeAccess {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

// VerifyRefreshToken validates and parses a refresh token
func (j *JWTService) VerifyRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != j.signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.refreshSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims.Type != TokenTypeRefresh {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

// ExtractTokenFromHeader extracts JWT token from Authorization header
func ExtractTokenFromHeader(authHeader string) string {
	// Try to parse as "Bearer <token>"
	parts := strings.Fields(authHeader)
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}

	// If the above fails, check if the header itself might be the token.
	// A simple check is to see if it contains dots, as JWTs do.
	if !strings.Contains(authHeader, " ") && strings.Count(authHeader, ".") == 2 {
		return authHeader
	}

	return ""
}
