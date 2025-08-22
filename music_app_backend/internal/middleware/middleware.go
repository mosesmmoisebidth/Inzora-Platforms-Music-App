package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mosesmmoisebidth/music_backend/internal/auth"
	"github.com/mosesmmoisebidth/music_backend/internal/config"
	"github.com/mosesmmoisebidth/music_backend/pkg/logger"
)

// RequestID middleware adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Logger middleware provides structured logging
func Logger(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := c.GetString("request_id")

		logger.With("request_id", requestID, "method", c.Request.Method, "path", c.Request.URL.Path).
			Info("Request started")

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		logLevel := "info"
		if status >= 400 && status < 500 {
			logLevel = "warn"
		} else if status >= 500 {
			logLevel = "error"
		}

		logFields := map[string]interface{}{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     status,
			"duration":   duration.String(),
			"user_agent": c.GetHeader("User-Agent"),
		}

		if userID := c.GetString("user_id"); userID != "" {
			logFields["user_id"] = userID
		}

		switch logLevel {
		case "warn":
			logger.With(logFields).Warn("Request completed with warning")
		case "error":
			logger.With(logFields).Error("Request completed with error")
		default:
			logger.With(logFields).Info("Request completed")
		}
	}
}

// Recovery middleware handles panics
func Recovery(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString("request_id")
				logger.With("request_id", requestID, "error", err).
					Error("Panic recovered")

				c.JSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{
						"code":    "INTERNAL_ERROR",
						"message": "Internal server error occurred",
					},
					"data": nil,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORS middleware configures CORS settings
func CORS(config config.CORSConfig) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     config.AllowOrigins,
		AllowMethods:     config.AllowMethods,
		AllowHeaders:     config.AllowHeaders,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           time.Duration(config.MaxAge) * time.Second,
	})
}

// JWTAuth middleware validates JWT tokens
func JWTAuth(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "MISSING_TOKEN",
					"message": "Authorization header is required",
				},
				"data": nil,
			})
			c.Abort()
			return
		}

		token := auth.ExtractTokenFromHeader(authHeader)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Invalid authorization header format",
				},
				"data": nil,
			})
			c.Abort()
			return
		}

		claims, err := jwtService.VerifyAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Invalid or expired token",
				},
				"data": nil,
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID.String())
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuth middleware validates JWT tokens but doesn't require them
func OptionalAuth(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		token := auth.ExtractTokenFromHeader(authHeader)
		if token == "" {
			c.Next()
			return
		}

		claims, err := jwtService.VerifyAccessToken(token)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID.String())
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("user_roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "INSUFFICIENT_PERMISSIONS",
					"message": "Insufficient permissions",
				},
				"data": nil,
			})
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "INSUFFICIENT_PERMISSIONS",
					"message": "Insufficient permissions",
				},
				"data": nil,
			})
			c.Abort()
			return
		}

		hasRole := false
		for _, userRole := range userRoles {
			if userRole == role || userRole == "admin" {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "INSUFFICIENT_PERMISSIONS",
					"message": "Insufficient permissions",
				},
				"data": nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecurityHeaders middleware adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Relax CSP for Swagger UI
		if strings.HasPrefix(c.Request.URL.Path, "/docs/") {
			c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';")
		} else {
			c.Header("Content-Security-Policy", "default-src 'self';")
		}

		c.Next()
	}
}

// ContentType middleware ensures JSON content type for API endpoints
func ContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": gin.H{
						"code":    "INVALID_CONTENT_TYPE",
						"message": "Content-Type must be application/json",
					},
					"data": nil,
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}