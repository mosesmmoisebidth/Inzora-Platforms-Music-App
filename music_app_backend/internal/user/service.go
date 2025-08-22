package user

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"github.com/lib/pq"
	"github.com/google/uuid"
	"github.com/mosesmmoisebidth/music_backend/internal/auth"
	"github.com/mosesmmoisebidth/music_backend/pkg/logger"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrEmailExists          = errors.New("user with this email already exists")
	ErrAuthenticationFailed = errors.New("authentication failed")
)

// Repository defines the interface for user data storage.
type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByGoogleID(ctx context.Context, googleID string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
}

// Service provides user business logic.
type Service struct {
	repo           Repository // Use the interface
	passwordHasher *auth.PasswordHasher
	logger         logger.Logger
}

// NewService creates a new user service.
func NewService(repo Repository, passwordHasher *auth.PasswordHasher, logger logger.Logger) *Service {
	return &Service{repo: repo, passwordHasher: passwordHasher, logger: logger}
}

// CreateUser creates a new user with a hashed password.
func (s *Service) CreateUser(ctx context.Context, email, password, displayName string) (*User, error) {
	_, err := s.repo.GetUserByEmail(ctx, email)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == nil {
			return nil, ErrEmailExists
		}
		return nil, err
	}

	hashedPassword, err := s.passwordHasher.Hash(password)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return nil, err
	}

	user := &User{
		ID:          uuid.New(),
		Email:       &email,
		Password:    &hashedPassword,
		DisplayName: &displayName,
		Roles:       pq.StringArray{"user"},
		IsActive:    true,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		s.logger.Error("failed to create user in repo", "error", err)
		return nil, err
	}

	return user, nil
}

// CreateGoogleUser creates or retrieves a user from a Google Sign-In.
func (s *Service) CreateGoogleUser(ctx context.Context, googleID, email, displayName, photoURL string) (*User, error) {
	user, err := s.repo.GetUserByGoogleID(ctx, googleID)
	if err == nil {
		now := time.Now()
		user.LastLoginAt = &now
		return user, s.repo.UpdateUser(ctx, user)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("failed to get user by google id", "error", err)
		return nil, err
	}

	newUser := &User{
		ID:          uuid.New(),
		GoogleID:    &googleID,
		Email:       &email,
		DisplayName: &displayName,
		PhotoURL:    &photoURL,
		Roles:       pq.StringArray{"user"},
		IsActive:    true,
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		s.logger.Error("failed to create google user in repo", "error", err)
		return nil, err
	}

	return newUser, nil
}

// AuthenticateUser checks a user's credentials and returns the user if valid.
func (s *Service) AuthenticateUser(ctx context.Context, email, password string) (*User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuthenticationFailed
		}
		return nil, err
	}

	if user.Password == nil {
		return nil, ErrAuthenticationFailed
	}

	match, err := s.passwordHasher.Matches(password, *user.Password)
	if err != nil || !match {
		return nil, ErrAuthenticationFailed
	}

	now := time.Now()
	user.LastLoginAt = &now
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		s.logger.Error("failed to update last login time", "error", err, "userID", user.ID)
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (s *Service) GetUserByID(ctx context.Context, userIDStr string) (*User, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// UpdateUser updates a user's profile information.
func (s *Service) UpdateUser(ctx context.Context, userIDStr string, displayName, photoURL *string, preferences map[string]interface{}) (*User, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if displayName != nil {
		user.DisplayName = displayName
	}
	if photoURL != nil {
		user.PhotoURL = photoURL
	}

	if preferences != nil {
		currentPrefs := make(map[string]interface{})
		if user.Preferences != nil {
			if err := json.Unmarshal(user.Preferences, &currentPrefs); err != nil {
				return nil, errors.New("failed to parse existing user preferences")
			}
		}

		for key, value := range preferences {
			currentPrefs[key] = value
		}

		newPrefs, err := json.Marshal(currentPrefs)
		if err != nil {
			return nil, errors.New("failed to serialize new user preferences")
		}
		user.Preferences = newPrefs
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		s.logger.Error("failed to update user in repo", "error", err, "userID", user.ID)
		return nil, err
	}

	return user, nil
}