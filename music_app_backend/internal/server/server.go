
package server

import (
	"time"
	"github.com/gin-gonic/gin"
	_ "github.com/mosesmmoisebidth/music_backend/docs" // This is required for swag to find docs
	"github.com/mosesmmoisebidth/music_backend/internal/auth"
	"github.com/mosesmmoisebidth/music_backend/internal/config"
	"github.com/mosesmmoisebidth/music_backend/internal/library"
	"github.com/mosesmmoisebidth/music_backend/internal/middleware"
	"github.com/mosesmmoisebidth/music_backend/internal/music"
	"github.com/mosesmmoisebidth/music_backend/internal/playlist"
	"github.com/mosesmmoisebidth/music_backend/internal/storage"
	httpTransport "github.com/mosesmmoisebidth/music_backend/internal/transport/http"
	"github.com/mosesmmoisebidth/music_backend/internal/user"
	"github.com/mosesmmoisebidth/music_backend/pkg/logger"
	"github.com/mosesmmoisebidth/music_backend/pkg/response"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the HTTP server
type Server struct {
	router  *gin.Engine
	config  *config.Config
	storage *storage.Storage
	logger  logger.Logger
}

// New creates a new server instance
func New(cfg *config.Config, storage *storage.Storage, logger logger.Logger) (*Server, error) {
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	server := &Server{
		config:  cfg,
		storage: storage,
		logger:  logger,
	}

	server.setupRouter()
	return server, nil
}

// setupRouter configures the Gin router
func (s *Server) setupRouter() {
	router := gin.New()

	// Global middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(s.logger))
	router.Use(middleware.Recovery(s.logger))
	router.Use(middleware.CORS(s.config.Server.CORS))
	router.Use(middleware.SecurityHeaders())

	// Health check and version endpoints
	router.GET("/healthz", s.healthCheck)
	router.GET("/version", s.versionInfo)

	// Swagger documentation route
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- Initialize Services and Repositories ---
	passwordHasher := auth.NewPasswordHasher(
		s.config.Auth.PasswordHashTime,
		s.config.Auth.PasswordHashMemory,
		s.config.Auth.PasswordHashThreads,
	)

	// Repositories
	userRepo := user.NewRepository(s.storage.DB)
	refreshTokenRepo := auth.NewRefreshTokenRepository(s.storage.DB)
	playlistRepo := playlist.NewRepository(s.storage.DB)
	libraryRepo := library.NewRepository(s.storage.DB)

	// Services
	jwtService := auth.NewJWTService(
		s.config.Auth.JWTAccessSecret,
		s.config.Auth.JWTRefreshSecret,
		s.config.Auth.AccessTokenTTL,
		s.config.Auth.RefreshTokenTTL,
		s.config.Auth.JWTSigningMethod,
	)
	googleService := auth.NewGoogleService(s.config.Google.ClientID)

	userService := user.NewService(userRepo, passwordHasher, s.logger)
	authService := auth.NewAuthService(jwtService, googleService, refreshTokenRepo, s.logger)
	playlistService := playlist.NewService(playlistRepo, s.logger)
	libraryService := library.NewService(libraryRepo, s.logger)
	musicService := music.NewMusicService(
		s.config.Providers.Enabled,
		30*time.Second, // timeout
		5*time.Minute,  // cache TTL
		s.config.Spotify.ClientID,
		s.config.Spotify.ClientSecret,
	)

	// --- Initialize Handlers ---
	authHandlers := httpTransport.NewAuthHandlers(userService, authService, s.logger)
	userHandlers := httpTransport.NewUserHandlers(userService, s.logger)
	playlistHandlers := httpTransport.NewPlaylistHandlers(playlistService, s.logger)
	libraryHandlers := httpTransport.NewLibraryHandlers(libraryService, s.logger)
	musicHandlers := httpTransport.NewMusicHandlers(musicService, s.logger)

	// --- API Routes ---
	api := router.Group("/api/v1")
	api.Use(middleware.ContentType())

	// Authentication routes (public)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandlers.Register)
		authGroup.POST("/login", authHandlers.Login)
		authGroup.POST("/google", authHandlers.GoogleSignIn)
		authGroup.POST("/refresh", authHandlers.RefreshToken)
		authGroup.POST("/logout", authHandlers.Logout)
	}

	jwtAuth := middleware.JWTAuth(jwtService)

	// User routes
	userGroup := api.Group("/users", jwtAuth)
	{
		userGroup.GET("/me", userHandlers.GetCurrentUser)
		userGroup.PATCH("/me", userHandlers.UpdateCurrentUser)
	}

	// Music routes
	musicGroup := api.Group("/music", middleware.OptionalAuth(jwtService))
	{
		musicGroup.GET("/search", musicHandlers.SearchTracks)
		musicGroup.GET("/tracks/:trackId", musicHandlers.GetTrack)
		musicGroup.GET("/top-charts", musicHandlers.GetTopCharts)
	}

	// Playlist routes
	playlistGroup := api.Group("/playlists", jwtAuth)
	{
		playlistGroup.GET("", playlistHandlers.GetPlaylists)
		playlistGroup.POST("", playlistHandlers.CreatePlaylist)
		playlistGroup.GET("/:playlistId", playlistHandlers.GetPlaylist)
		playlistGroup.PATCH("/:playlistId", playlistHandlers.UpdatePlaylist)
		playlistGroup.DELETE("/:playlistId", playlistHandlers.DeletePlaylist)
		playlistGroup.POST("/:playlistId/tracks", playlistHandlers.AddTrackToPlaylist)
		playlistGroup.DELETE("/:playlistId/tracks/:trackId", playlistHandlers.RemoveTrackFromPlaylist)
	}

	// Library routes
	libraryGroup := api.Group("", jwtAuth)
	{
		libraryGroup.GET("/favorites", libraryHandlers.GetFavorites)
		libraryGroup.POST("/favorites", libraryHandlers.AddFavorite)
		libraryGroup.DELETE("/favorites/:favoriteId", libraryHandlers.RemoveFavorite)
		libraryGroup.GET("/history", libraryHandlers.GetHistory)
		libraryGroup.POST("/history", libraryHandlers.AddHistory)
	}

	s.router = router
}

// Router returns the configured Gin router
func (s *Server) Router() *gin.Engine {
	return s.router
}

// healthCheck handles health check requests
func (s *Server) healthCheck(c *gin.Context) {
	if err := s.storage.Health(); err != nil {
		s.logger.Error("Health check failed", "error", err)
		response.InternalError(c, "UNHEALTHY", "Service is unhealthy")
		return
	}

	healthData := gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"services": gin.H{
			"database": "healthy",
			"redis":    "healthy",
		},
	}

	response.Success(c, healthData)
}

// versionInfo handles version information requests
func (s *Server) versionInfo(c *gin.Context) {
	versionData := gin.H{
		"version":     "1.2.0",
		"environment": s.config.App.Environment,
	}
	response.Success(c, versionData)
}
