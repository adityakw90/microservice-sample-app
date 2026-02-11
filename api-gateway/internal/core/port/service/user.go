package service

import (
	"context"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
)

// UserService defines the interface for user-related use cases.
// This is a primary port - what the application provides to the outside world.
type UserService interface {
	// GetUser retrieves a user by UID.
	GetUser(ctx context.Context, uid string) (*model.User, error)

	// ListUsers returns a paginated list of users based on filters.
	ListUsers(ctx context.Context, param *params.ListUsersParam) (*model.Users, error)

	// CreateUser creates a new user and returns its UID.
	CreateUser(ctx context.Context, param *params.CreateUserParam) (string, error)

	// UpdateUser updates an existing user.
	UpdateUser(ctx context.Context, param *params.UpdateUserParam) error

	// DeleteUser soft-deletes a user by UID.
	DeleteUser(ctx context.Context, uid string) error

	// GetUserProfile retrieves a user's profile.
	GetUserProfile(ctx context.Context, userUID string) (*model.Profile, error)

	// UpdateUserProfile updates a user's profile.
	UpdateUserProfile(ctx context.Context, param *params.UpdateUserProfileParam) error

	// ListUserDevices returns a paginated list of user devices.
	ListUserDevices(ctx context.Context, param *params.ListUserDevicesParam) (*model.Devices, error)

	// RevokeUserDevice revokes a device for a user.
	RevokeUserDevice(ctx context.Context, userUID, deviceUID string) error
}
