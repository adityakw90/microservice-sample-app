package client

import (
	"context"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
)

// UserClient defines the interface for calling user gRPC service.
// This is a secondary port - what the application needs from external services.
type UserClient interface {
	// Get retrieves a user by UID.
	Get(ctx context.Context, uid string) (*model.User, error)

	// List returns a paginated list of users based on filters.
	List(ctx context.Context, param *params.ListUsersParam) (*model.Users, error)

	// Create creates a new user and returns its UID.
	Create(ctx context.Context, param *params.CreateUserParam) (string, error)

	// Update updates an existing user.
	Update(ctx context.Context, param *params.UpdateUserParam) error

	// Delete soft-deletes a user by UID.
	Delete(ctx context.Context, uid string) error

	// GetProfile retrieves a user's profile.
	GetProfile(ctx context.Context, userUID string) (*model.Profile, error)

	// UpdateProfile updates a user's profile.
	UpdateProfile(ctx context.Context, param *params.UpdateUserProfileParam) error

	// ListDevice returns a paginated list of user devices.
	ListDevice(ctx context.Context, param *params.ListUserDevicesParam) (*model.Devices, error)

	// RevokeDevice revokes a device for a user.
	RevokeDevice(ctx context.Context, userUID, deviceUID string) error
}
