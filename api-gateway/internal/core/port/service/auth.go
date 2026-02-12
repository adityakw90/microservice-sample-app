package service

import (
	"context"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
)

// AuthService defines the interface for authentication-related use cases.
// This is a primary port - what the application provides to the outside world.
type AuthService interface {
	// Login authenticates a user and returns tokens.
	Login(ctx context.Context, param *params.LoginParam) (*model.Tokens, error)

	// RefreshToken refreshes an access token using a refresh token.
	RefreshToken(ctx context.Context, param *params.RefreshTokenParam) (*model.Tokens, error)

	// ValidateToken validates an access token and returns user claims.
	ValidateToken(ctx context.Context, param *params.ValidateTokenParam) (*model.TokenClaims, error)

	// VerifyPin verifies a user PIN for sensitive actions.
	VerifyPin(ctx context.Context, param *params.VerifyPinParam) (bool, error)

	// GoogleOAuth initiates Google OAuth flow and returns authorization URL.
	GoogleOAuth(ctx context.Context, redirectURI string) (authorizationURL string, err error)

	// HandleGoogleOAuth handles Google OAuth callback and returns tokens.
	HandleGoogleOAuth(ctx context.Context, code, redirectURI string) (*model.Tokens, error)
}
