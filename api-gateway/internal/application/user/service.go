package user

import (
	"context"
	"fmt"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/errors"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/port/client"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/port/service"
)

// UserApplicationService implements service.UserService use cases.
type UserApplicationService struct {
	userClient client.UserClient
}

// NewUserApplicationService creates a new UserApplicationService.
func NewUserApplicationService(userClient client.UserClient) service.UserService {
	return &UserApplicationService{
		userClient: userClient,
	}
}

// GetUser retrieves a user by UID with business validation.
func (s *UserApplicationService) GetUser(ctx context.Context, uid string) (*model.User, error) {
	if uid == "" {
		return nil, errors.ErrInvalidUID
	}

	user, err := s.userClient.Get(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.CanDisplay() {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

// ListUsers returns a paginated list of users based on filters.
func (s *UserApplicationService) ListUsers(ctx context.Context, param *params.ListUsersParam) (*model.Users, error) {
	// Set default pagination if not provided
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 20
	}

	users, err := s.userClient.List(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// CreateUser creates a new user and returns its UID.
func (s *UserApplicationService) CreateUser(ctx context.Context, param *params.CreateUserParam) (string, error) {
	// Business validation
	if param.Username == "" {
		return "", errors.ErrInvalidInput
	}
	if param.Email == "" {
		return "", errors.ErrInvalidInput
	}
	if param.Password == "" {
		return "", errors.ErrInvalidInput
	}

	uid, err := s.userClient.Create(ctx, param)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return uid, nil
}

// UpdateUser updates an existing user.
func (s *UserApplicationService) UpdateUser(ctx context.Context, param *params.UpdateUserParam) error {
	if param.UID == "" {
		return errors.ErrInvalidUID
	}

	err := s.userClient.Update(ctx, param)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser soft-deletes a user by UID.
func (s *UserApplicationService) DeleteUser(ctx context.Context, uid string) error {
	if uid == "" {
		return errors.ErrInvalidUID
	}

	err := s.userClient.Delete(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// GetUserProfile retrieves a user's profile.
func (s *UserApplicationService) GetUserProfile(ctx context.Context, userUID string) (*model.Profile, error) {
	if userUID == "" {
		return nil, errors.ErrInvalidUID
	}

	profile, err := s.userClient.GetProfile(ctx, userUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	if profile.IsEmpty() {
		return nil, errors.ErrProfileNotFound
	}

	return profile, nil
}

// UpdateUserProfile updates a user's profile.
func (s *UserApplicationService) UpdateUserProfile(ctx context.Context, param *params.UpdateUserProfileParam) error {
	if param.UserUID == "" {
		return errors.ErrInvalidUID
	}

	err := s.userClient.UpdateProfile(ctx, param)
	if err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	return nil
}

// ListUserDevices returns a paginated list of user devices.
func (s *UserApplicationService) ListUserDevices(ctx context.Context, param *params.ListUserDevicesParam) (*model.Devices, error) {
	if param.UserUID == "" {
		return nil, errors.ErrInvalidUID
	}

	// Set default pagination if not provided
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 20
	}

	devices, err := s.userClient.ListDevice(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("failed to list user devices: %w", err)
	}

	return devices, nil
}

// RevokeUserDevice revokes a device for a user.
func (s *UserApplicationService) RevokeUserDevice(ctx context.Context, userUID, deviceUID string) error {
	if userUID == "" {
		return errors.ErrInvalidUID
	}
	if deviceUID == "" {
		return errors.ErrInvalidInput
	}

	err := s.userClient.RevokeDevice(ctx, userUID, deviceUID)
	if err != nil {
		return fmt.Errorf("failed to revoke device: %w", err)
	}

	return nil
}
