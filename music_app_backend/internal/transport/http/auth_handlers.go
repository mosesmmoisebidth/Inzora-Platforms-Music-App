
package http

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mosesmmoisebidth/music_backend/internal/auth"
	"github.com/mosesmmoisebidth/music_backend/internal/user"
	"github.com/mosesmmoisebidth/music_backend/pkg/logger"
	"github.com/mosesmmoisebidth/music_backend/pkg/response"
)

// AuthHandlers contains authentication HTTP handlers
type AuthHandlers struct {
	userService *user.Service
	authService *auth.AuthService
	logger      logger.Logger
}

// NewAuthHandlers creates new authentication handlers
func NewAuthHandlers(userService *user.Service, authService *auth.AuthService, logger logger.Logger) *AuthHandlers {
	return &AuthHandlers{
		userService: userService,
		authService: authService,
		logger:      logger,
	}
}

// Register handles user registration
// @Summary      Register a new user
// @Description  Register a new user with email, password, and display name.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Registration Request"
// @Success      201  {object}  response.APIResponse{data=AuthResponse}
// @Failure      400  {object}  response.APIResponse{error=response.APIError}
// @Failure      409  {object}  response.APIResponse{error=response.APIError}
// @Failure      500  {object}  response.APIResponse{error=response.APIError}
// @Router       /auth/register [post]
func (h *AuthHandlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	newUser, err := h.userService.CreateUser(ctx, req.Email, req.Password, req.DisplayName)
	if err != nil {
		if err == user.ErrEmailExists {
			response.Conflict(c, "USER_EXISTS", err.Error())
			return
		}
		h.logger.Error("Failed to create user", "error", err, "email", req.Email)
		response.InternalError(c, "REGISTRATION_FAILED", "Failed to create user account")
		return
	}

	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()
	email := ""
	if newUser.Email != nil {
		email = *newUser.Email
	}

	tokens, err := h.authService.GenerateTokens(ctx, newUser.ID, email, newUser.Roles, userAgent, clientIP)
	if err != nil {
		h.logger.Error("Failed to generate tokens", "error", err, "user_id", newUser.ID)
		response.InternalError(c, "TOKEN_GENERATION_FAILED", "Failed to generate authentication tokens")
		return
	}

	authResponse := &AuthResponse{
		User:         mapUserToResponse(newUser),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
	}

	h.logger.Info("User registered successfully", "user_id", newUser.ID)
	response.Created(c, authResponse)
}

// Login handles user login
// @Summary      User login
// @Description  Authenticate a user with email and password to receive JWT tokens.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login Request"
// @Success      200  {object}  response.APIResponse{data=AuthResponse}
// @Failure      400  {object}  response.APIResponse{error=response.APIError}
// @Failure      401  {object}  response.APIResponse{error=response.APIError}
// @Failure      500  {object}  response.APIResponse{error=response.APIError}
// @Router       /auth/login [post]
func (h *AuthHandlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	authenticatedUser, err := h.userService.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		h.logger.Warn("Authentication failed", "error", err, "email", req.Email)
		response.Unauthorized(c, "AUTHENTICATION_FAILED", "Invalid email or password")
		return
	}

	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()
	email := ""
	if authenticatedUser.Email != nil {
		email = *authenticatedUser.Email
	}

	tokens, err := h.authService.GenerateTokens(ctx, authenticatedUser.ID, email, authenticatedUser.Roles, userAgent, clientIP)
	if err != nil {
		h.logger.Error("Failed to generate tokens", "error", err, "user_id", authenticatedUser.ID)
		response.InternalError(c, "TOKEN_GENERATION_FAILED", "Failed to generate authentication tokens")
		return
	}

	authResponse := &AuthResponse{
		User:         mapUserToResponse(authenticatedUser),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
	}

	h.logger.Info("User logged in successfully", "user_id", authenticatedUser.ID)
	response.Success(c, authResponse)
}

// GoogleSignIn handles Google Sign-In
// @Summary      Google Sign-In
// @Description  Authenticate or register a user using a Google ID token.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body GoogleSignInRequest true "Google Sign-In Request"
// @Success      200  {object}  response.APIResponse{data=AuthResponse}
// @Failure      400  {object}  response.APIResponse{error=response.APIError}
// @Failure      401  {object}  response.APIResponse{error=response.APIError}
// @Failure      500  {object}  response.APIResponse{error=response.APIError}
// @Router       /auth/google [post]
func (h *AuthHandlers) GoogleSignIn(c *gin.Context) {
	var req GoogleSignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	googleUser, err := h.authService.VerifyGoogleIDToken(ctx, req.IDToken)
	if err != nil {
		h.logger.Warn("Google ID token verification failed", "error", err)
		response.Unauthorized(c, "INVALID_GOOGLE_TOKEN", "Invalid Google ID token")
		return
	}

	appUser, err := h.userService.CreateGoogleUser(ctx, googleUser.ID, googleUser.Email, googleUser.Name, googleUser.Picture)
	if err != nil {
		h.logger.Error("Failed to create/get Google user", "error", err, "google_sub", googleUser.ID)
		response.InternalError(c, "GOOGLE_USER_CREATION_FAILED", "Failed to process Google user")
		return
	}

	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()
	email := ""
	if appUser.Email != nil {
		email = *appUser.Email
	}

	tokens, err := h.authService.GenerateTokens(ctx, appUser.ID, email, appUser.Roles, userAgent, clientIP)
	if err != nil {
		h.logger.Error("Failed to generate tokens", "error", err, "user_id", appUser.ID)
		response.InternalError(c, "TOKEN_GENERATION_FAILED", "Failed to generate authentication tokens")
		return
	}

	authResponse := &AuthResponse{
		User:         mapUserToResponse(appUser),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
	}

	h.logger.Info("Google user signed in successfully", "user_id", appUser.ID)
	response.Success(c, authResponse)
}

// RefreshToken handles token refresh
// @Summary      Refresh access token
// @Description  Obtain a new access token using a valid refresh token.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body RefreshTokenRequest true "Refresh Token Request"
// @Success      200  {object}  response.APIResponse{data=TokenPair}
// @Failure      400  {object}  response.APIResponse{error=response.APIError}
// @Failure      401  {object}  response.APIResponse{error=response.APIError}
// @Failure      500  {object}  response.APIResponse{error=response.APIError}
// @Router       /auth/refresh [post]
func (h *AuthHandlers) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()

	newTokens, err := h.authService.RefreshTokens(ctx, req.RefreshToken, userAgent, clientIP)
	if err != nil {
		h.logger.Warn("Token refresh failed", "error", err)
		response.Unauthorized(c, "TOKEN_REFRESH_FAILED", "Invalid or expired refresh token")
		return
	}

	h.logger.Info("Tokens refreshed successfully")
	response.Success(c, newTokens)
}

// Logout handles user logout
// @Summary      User logout
// @Description  Revoke the user's refresh token.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body RefreshTokenRequest true "Logout Request"
// @Success      200  {object}  response.APIResponse{data=LogoutResponse}
// @Failure      400  {object}  response.APIResponse{error=response.APIError}
// @Failure      500  {object}  response.APIResponse{error=response.APIError}
// @Router       /auth/logout [post]
func (h *AuthHandlers) Logout(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := h.authService.RevokeToken(ctx, req.RefreshToken); err != nil {
		h.logger.Warn("Token revocation failed", "error", err)
	}

	h.logger.Info("User logged out successfully")
	response.Success(c, &LogoutResponse{
		Message: "Logged out successfully",
	})
}
