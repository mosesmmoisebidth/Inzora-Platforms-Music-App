package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Providers ProvidersConfig `mapstructure:"providers"`
	Google    GoogleConfig    `mapstructure:"google"`
	Spotify   SpotifyConfig   `mapstructure:"spotify"`
}

// AppConfig contains general application configuration
type AppConfig struct {
	Name        string `mapstructure:"name" default:"music-app-backend"`
	Environment string `mapstructure:"environment" default:"development"`
	LogLevel    string `mapstructure:"log_level" default:"info"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Port         string        `mapstructure:"port" default:"8080"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" default:"30s"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" default:"30s"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout" default:"120s"`
	CORS         CORSConfig    `mapstructure:"cors"`
}

// CORSConfig contains CORS configuration
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins" default:"*"`
	AllowMethods     []string `mapstructure:"allow_methods" default:"GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS"`
	AllowHeaders     []string `mapstructure:"allow_headers" default:"Origin,Content-Length,Content-Type,Authorization"`
	AllowCredentials bool     `mapstructure:"allow_credentials" default:"true"`
	MaxAge           int      `mapstructure:"max_age" default:"86400"`
}

// DatabaseConfig contains database connection configuration
type DatabaseConfig struct {
	Host         string `mapstructure:"host" default:"localhost"`
	Port         int    `mapstructure:"port" default:"5432"`
	User         string `mapstructure:"user" default:"postgres"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name" default:"music_app"`
	SSLMode      string `mapstructure:"ssl_mode" default:"disable"`
	MaxOpenConns int    `mapstructure:"max_open_conns" default:"25"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" default:"25"`
	MaxLifetime  string `mapstructure:"max_lifetime" default:"5m"`
}

// DSN returns the database connection string
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

// RedisConfig contains Redis connection configuration
type RedisConfig struct {
	Addr         string `mapstructure:"addr" default:"localhost:6379"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db" default:"0"`
	PoolSize     int    `mapstructure:"pool_size" default:"10"`
	MinIdleConns int    `mapstructure:"min_idle_conns" default:"5"`
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	JWTSigningMethod    string        `mapstructure:"jwt_signing_method" default:"HS256"`
	JWTAccessSecret     string        `mapstructure:"jwt_access_secret"`
	JWTRefreshSecret    string        `mapstructure:"jwt_refresh_secret"`
	AccessTokenTTL      time.Duration `mapstructure:"access_token_ttl" default:"15m"`
	RefreshTokenTTL     time.Duration `mapstructure:"refresh_token_ttl" default:"720h"`
	PasswordHashMemory  uint32        `mapstructure:"password_hash_memory" default:"65536"`
	PasswordHashTime    uint32        `mapstructure:"password_hash_time" default:"3"`
	PasswordHashThreads uint8         `mapstructure:"password_hash_threads" default:"2"`
}

// GoogleConfig contains Google OAuth configuration
type GoogleConfig struct {
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	RedirectURLs []string `mapstructure:"redirect_urls"`
}

// SpotifyConfig contains Spotify API configuration
type SpotifyConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

// ProvidersConfig contains music providers configuration
type ProvidersConfig struct {
	Enabled        []string `mapstructure:"enabled" default:"itunes"`
	DefaultTimeout string   `mapstructure:"default_timeout" default:"30s"`
	CacheTTL       string   `mapstructure:"cache_ttl" default:"300s"`
	RateLimit      int      `mapstructure:"rate_limit" default:"100"`
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	// Set up environment variable handling
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("MUSIC_APP")

	viper.BindEnv("auth.jwt_access_secret")
	viper.BindEnv("auth.jwt_refresh_secret")
	viper.BindEnv("database.password")
	// Also bind secrets for providers if you use them
	viper.BindEnv("google.client_id")
	viper.BindEnv("google.client_secret")
	viper.BindEnv("spotify.client_id")
	viper.BindEnv("spotify.client_secret")

	// Set defaults
	setDefaults()

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found is okay, we can use env vars and defaults
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate required configuration
	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "music-app-backend")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.log_level", "info")

	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")

	// CORS defaults
	viper.SetDefault("server.cors.allow_origins", []string{"*"})
	viper.SetDefault("server.cors.allow_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"})
	viper.SetDefault("server.cors.allow_headers", []string{"Origin", "Content-Length", "Content-Type", "Authorization"})
	viper.SetDefault("server.cors.allow_credentials", true)
	viper.SetDefault("server.cors.max_age", 86400)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.name", "music_app")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 25)
	viper.SetDefault("database.max_lifetime", "5m")

	// Redis defaults
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)

	// Auth defaults
	viper.SetDefault("auth.jwt_signing_method", "HS256")
	viper.SetDefault("auth.access_token_ttl", "15m")
	viper.SetDefault("auth.refresh_token_ttl", "720h")
	viper.SetDefault("auth.password_hash_memory", 65536)
	viper.SetDefault("auth.password_hash_time", 3)
	viper.SetDefault("auth.password_hash_threads", 2)

	// Providers defaults
	viper.SetDefault("providers.enabled", []string{"itunes"})
	viper.SetDefault("providers.default_timeout", "30s")
	viper.SetDefault("providers.cache_ttl", "300s")
	viper.SetDefault("providers.rate_limit", 100)
}

func validate(config *Config) error {
	// Validate required auth secrets
	if config.Auth.JWTAccessSecret == "" {
		return fmt.Errorf("JWT access secret is required")
	}
	if config.Auth.JWTRefreshSecret == "" {
		return fmt.Errorf("JWT refresh secret is required")
	}

	// Validate database password
	if config.Database.Password == "" {
		return fmt.Errorf("database password is required")
	}

	// Validate Google config if Google is enabled
	for _, provider := range config.Providers.Enabled {
		if provider == "google" && config.Google.ClientID == "" {
			return fmt.Errorf("Google client ID is required when Google provider is enabled")
		}
	}

	// Validate Spotify config if Spotify is enabled
	for _, provider := range config.Providers.Enabled {
		if provider == "spotify" && config.Spotify.ClientID == "" {
			return fmt.Errorf("Spotify client ID is required when Spotify provider is enabled")
		}
	}

	return nil
}
