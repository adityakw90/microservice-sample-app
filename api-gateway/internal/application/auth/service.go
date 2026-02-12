package auth

import (
	"context"
	"fmt"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/errors"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/port/client"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/port/service"
)

// AuthApplicationService implements service.AuthService use cases.
type AuthApplicationService struct {
	authClient client.AuthClient
}

// NewAuthApplicationService creates a new AuthApplicationService.
func NewAuthApplicationService(authClient client.AuthClient) service.AuthService {
	return &AuthApplicationService{
		authClient: authClient,
	}
}

// Login authenticates a user and returns tokens with business validation.
func (s *AuthApplicationService) Login(ctx context.Context, param *params.LoginParam) (*model.Tokens, error) {
	// Business validation
	if param.Identifier == "" {
		return nil, errors.ErrInvalidInput
	}
	if param.IdentifierType == "" {
		return nil, errors.ErrInvalidInput
	}
	if param.Password == "" {
		return nil, errors.ErrInvalidCredentials
	}

	// Validate identifier type
	if param.IdentifierType != "username" && param.IdentifierType != "email" {
		return nil, errors.ErrInvalidInput
	}

	tokens, err := s.authClient.Authenticate(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return tokens, nil
}

// RefreshToken refreshes an access token using a refresh token.
func (s *AuthApplicationService) RefreshToken(ctx context.Context, param *params.RefreshTokenParam) (*model.Tokens, error) {
	if param.RefreshToken == "" {
		return nil, errors.ErrInvalidToken
	}

	tokens, err := s.authClient.RefreshToken(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return tokens, nil
}

// ValidateToken validates an access token and returns user claims.
func (s *AuthApplicationService) ValidateToken(ctx context.Context, param *params.ValidateTokenParam) (*model.TokenClaims, error) {
	if param.AccessToken == "" {
		return nil, errors.ErrInvalidToken
	}

	claims, err := s.authClient.ValidateToken(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	return claims, nil
}

// VerifyPin verifies a user PIN for sensitive actions.
func (s *AuthApplicationService) VerifyPin(ctx context.Context, param *params.VerifyPinParam) (bool, error) {
	if param.UID == "" {
		return false, errors.ErrInvalidUID
	}
	if param.Code == "" {
		return false, errors.ErrInvalidPIN
	}

	valid, err := s.authClient.VerifyPin(ctx, param)
	if err != nil {
		return false, fmt.Errorf("failed to verify PIN: %w", err)
	}

	return valid, nil
}

// GoogleOAuth initiates Google OAuth flow and returns authorization URL.
func (s *AuthApplicationService) GoogleOAuth(ctx context.Context, redirectURI string) (string, error) {
	return s.authClient.GoogleOAuth(ctx, redirectURI)
}

// HandleGoogleOAuth handles Google OAuth callback and returns tokens.
func (s *AuthApplicationService) HandleGoogleOAuth(ctx context.Context, code, redirectURI string) (*model.Tokens, error) {
	return s.authClient.HandleGoogleOAuth(ctx, code, redirectURI)
}
