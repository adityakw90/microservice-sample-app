package client

import (
	"context"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
)

// AuthClient defines the interface for calling auth gRPC service.
// This is a secondary port - what the application needs from external services.
type AuthClient interface {
	// Authenticate authenticates a user and returns tokens.
	Authenticate(ctx context.Context, param *params.LoginParam) (*model.Tokens, error)

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
