package service

import (
	"context"
	"fmt"

	authpb "github.com/adityakw90/service-user-proto/gen/go/auth"
)

// AuthService handles business logic for authentication operations
type AuthService struct {
	client authpb.AuthServiceClient
}

// NewAuthService creates a new auth service
func NewAuthService(client authpb.AuthServiceClient) *AuthService {
	return &AuthService{
		client: client,
	}
}

// LoginRequest represents the request to login
type LoginRequest struct {
	Identifier        string `json:"identifier"`
	IdentifierType    string `json:"identifier_type"` // "username" or "email"
	Password          string `json:"password"`
	DeviceFingerprint string `json:"device_fingerprint,omitempty"`
	DeviceName        string `json:"device_name,omitempty"`
}

// TokenResponse represents the token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenRequest represents the request to refresh a token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// ValidateTokenRequest represents the request to validate a token
type ValidateTokenRequest struct {
	AccessToken string `json:"access_token"`
}

// ValidateTokenResponse represents the response from validating a token
type ValidateTokenResponse struct {
	UID            string                 `json:"uid"`
	Identifier     string                 `json:"identifier"`
	IdentifierType string                 `json:"identifier_type"`
	Claims         map[string]any `json:"claims,omitempty"`
}

// VerifyPinRequest represents the request to verify a PIN
type VerifyPinRequest struct {
	UID  string `json:"uid"`
	Code string `json:"code"`
}

// VerifyPinResponse represents the response from verifying a PIN
type VerifyPinResponse struct {
	Valid bool `json:"valid"`
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	pbResp, err := s.client.Auth(ctx, &authpb.AuthRequest{
		Identifier:        req.Identifier,
		IdentifierType:    req.IdentifierType,
		Password:          req.Password,
		DeviceFingerprint: req.DeviceFingerprint,
		DeviceName:        req.DeviceName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return &TokenResponse{
		AccessToken:  pbResp.AccessToken,
		RefreshToken: pbResp.RefreshToken,
	}, nil
}

// RefreshToken refreshes the access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*TokenResponse, error) {
	pbResp, err := s.client.RefreshToken(ctx, &authpb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  pbResp.AccessToken,
		RefreshToken: pbResp.RefreshToken,
	}, nil
}

// ValidateToken validates the access token
func (s *AuthService) ValidateToken(ctx context.Context, req ValidateTokenRequest) (*ValidateTokenResponse, error) {
	pbResp, err := s.client.ValidateToken(ctx, &authpb.ValidateTokenRequest{
		AccessToken: req.AccessToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	resp := &ValidateTokenResponse{
		UID:            pbResp.Uid,
		Identifier:     pbResp.Identifier,
		IdentifierType: pbResp.IdentifierType,
	}

	if pbResp.Claims != nil {
		resp.Claims = pbResp.Claims.AsMap()
	}

	return resp, nil
}

// VerifyPin verifies the user PIN for sensitive actions
func (s *AuthService) VerifyPin(ctx context.Context, req VerifyPinRequest) (*VerifyPinResponse, error) {
	pbResp, err := s.client.VerifyPin(ctx, &authpb.VerifyPinRequest{
		Uid:  req.UID,
		Code: req.Code,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to verify PIN: %w", err)
	}

	return &VerifyPinResponse{
		Valid: pbResp.Valid,
	}, nil
}
