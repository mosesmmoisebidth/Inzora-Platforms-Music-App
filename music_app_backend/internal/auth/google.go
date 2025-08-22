package auth

import (
	"context"
	"fmt"
	"google.golang.org/api/idtoken"
)

// GoogleUser represents user information from Google
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// GoogleService provides Google Sign-In verification
type GoogleService struct {
	clientID string
}

// NewGoogleService creates a new Google service
func NewGoogleService(clientID string) *GoogleService {
	return &GoogleService{
		clientID: clientID,
	}
}

// VerifyIDToken verifies a Google ID token and returns user information
func (g *GoogleService) VerifyIDToken(ctx context.Context, idToken string) (*GoogleUser, error) {
	// Validate the ID token
	payload, err := idtoken.Validate(ctx, idToken, g.clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	// Validate issuer
	if payload.Issuer != "accounts.google.com" && payload.Issuer != "https://accounts.google.com" {
		return nil, fmt.Errorf("invalid issuer: %s", payload.Issuer)
	}

	// The idtoken.Validate function already checks the audience (clientID).
	// It also checks if the email is verified if the "email_verified" claim exists.

	// Extract claims
	claims := payload.Claims
	email, _ := claims["email"].(string)
	emailVerified, _ := claims["email_verified"].(bool)
	name, _ := claims["name"].(string)
	givenName, _ := claims["given_name"].(string)
	familyName, _ := claims["family_name"].(string)
	picture, _ := claims["picture"].(string)
	locale, _ := claims["locale"].(string)

	if !emailVerified {
		return nil, fmt.Errorf("email not verified")
	}

	return &GoogleUser{
		ID:            payload.Subject, // Use Subject for the user ID
		Email:         email,
		VerifiedEmail: emailVerified,
		Name:          name,
		GivenName:     givenName,
		FamilyName:    familyName,
		Picture:       picture,
		Locale:        locale,
	}, nil
}
