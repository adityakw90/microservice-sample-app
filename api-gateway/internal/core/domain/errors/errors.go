package errors

import "errors"

// Domain error definitions
var (
	// User errors
	ErrInvalidUID      = errors.New("invalid user UID")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserDeleted     = errors.New("user has been deleted")
	ErrUserInactive    = errors.New("user is inactive")

	// Profile errors
	ErrProfileNotFound = errors.New("profile not found")

	// Device errors
	ErrDeviceNotFound = errors.New("device not found")
	ErrDeviceRevoked  = errors.New("device has been revoked")

	// Auth errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken      = errors.New("invalid token")
	ErrExpiredToken      = errors.New("token has expired")
	ErrInvalidPIN       = errors.New("invalid PIN code")

	// Validation errors
	ErrInvalidInput = errors.New("invalid input")
)
