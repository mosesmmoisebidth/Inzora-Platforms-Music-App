package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/mosesmmoisebidth/music_backend/internal/auth"
	"github.com/mosesmmoisebidth/music_backend/internal/config"
	"github.com/mosesmmoisebidth/music_backend/internal/library"
	"github.com/mosesmmoisebidth/music_backend/internal/playlist"
	"github.com/mosesmmoisebidth/music_backend/internal/user"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Storage holds database and cache connections
type Storage struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// New creates a new Storage instance with database and Redis connections
func New(dbConfig config.DatabaseConfig, redisConfig config.RedisConfig) (*Storage, error) {
	// Initialize PostgreSQL connection
	db, err := initDB(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize Redis connection
	redisClient, err := initRedis(redisConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	return &Storage{
		DB:    db,
		Redis: redisClient,
	}, nil
}

// initDB initializes the PostgreSQL database connection
func initDB(config config.DatabaseConfig) (*gorm.DB, error) {
	dsn := config.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	
	if duration, err := time.ParseDuration(config.MaxLifetime); err == nil {
		sqlDB.SetConnMaxLifetime(duration)
	}

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// initRedis initializes the Redis connection
func initRedis(config config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	})

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

// AutoMigrate runs database migrations for all models
func (s *Storage) AutoMigrate() error {
	return s.DB.AutoMigrate(
		&user.User{},
		&playlist.Playlist{},
		&playlist.PlaylistTrack{},
		&library.Favorite{},
		&library.History{},
		&library.Download{},
		&auth.RefreshToken{},
	)
}

// Close closes all database connections
func (s *Storage) Close() error {
	// Close PostgreSQL connection
	if sqlDB, err := s.DB.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
	}

	// Close Redis connection
	if err := s.Redis.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	return nil
}

// Health checks the health of all storage connections
func (s *Storage) Health() error {
	// Check PostgreSQL health
	if sqlDB, err := s.DB.DB(); err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	} else if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	// Check Redis health
	ctx := context.Background()
	if err := s.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	return nil
}
