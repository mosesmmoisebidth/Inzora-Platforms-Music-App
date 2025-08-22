package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mosesmmoisebidth/music_backend/internal/config"
	"github.com/mosesmmoisebidth/music_backend/internal/server"
	"github.com/mosesmmoisebidth/music_backend/internal/storage"
	"github.com/mosesmmoisebidth/music_backend/pkg/logger"
)

// @title           Music App Backend API
// @version         1.2.0
// @description     This is the backend API for the Music App, providing endpoints for authentication, music discovery, playlists, and user library management.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name   MIT
// @license.url    https://opensource.org/licenses/MIT

// @host      localhost:8085
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	logger := logger.NewLogger()

	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", "error", err)
	}

	storage, err := storage.New(cfg.Database, cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to initialize storage", "error", err)
	}

	if err := storage.AutoMigrate(); err != nil {
		logger.Fatal("Failed to run migrations", "error", err)
	}

	srv, err := server.New(cfg, storage, logger)
	if err != nil {
		logger.Fatal("Failed to initialize server", "error", err)
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      srv.Router(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.Info("Starting server", "port", cfg.Server.Port, "url", fmt.Sprintf("http://localhost:%s", cfg.Server.Port))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	if err := storage.Close(); err != nil {
		logger.Error("Error closing storage connections", "error", err)
	}

	logger.Info("Server exited")
}