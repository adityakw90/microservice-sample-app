package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userpb "github.com/adityakw90/service-user-proto/gen/go/user"
	authpb "github.com/adityakw90/service-user-proto/gen/go/auth"
	"github.com/adityakw90/service-user-proto/gen/go/common"

	grpcAdapter "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/secondary/grpc/mapper"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/port/client"
)

// UserClientAdapter implements client.UserClient using gRPC.
// This is a secondary adapter - driven by external services.
type UserClientAdapter struct {
	userClient userpb.UserServiceClient
	authClient authpb.AuthServiceClient
	conn       *grpc.ClientConn
}

// NewUserClientAdapter creates a new UserClientAdapter with gRPC connections.
func NewUserClientAdapter(serviceAddress string) (client.UserClient, client.AuthClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		serviceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(4*1024*1024),
			grpc.MaxCallSendMsgSize(4*1024*1024),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial user service: %w", err)
	}

	return &UserClientAdapter{
		userClient: userpb.NewUserServiceClient(conn),
		authClient: authpb.NewAuthServiceClient(conn),
		conn:       conn,
	}, nil, nil
}

// Close closes the gRPC connection.
func (a *UserClientAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

// === UserClient Implementation ===

// Get retrieves a user by UID.
func (a *UserClientAdapter) Get(ctx context.Context, uid string) (*model.User, error) {
	resp, err := a.userClient.Get(ctx, &userpb.GetRequest{Uid: uid})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return grpcAdapter.UserFromProto(resp), nil
}

// List returns a paginated list of users based on filters.
func (a *UserClientAdapter) List(ctx context.Context, param *params.ListUsersParam) (*model.Users, error) {
	pbReq := &userpb.ListRequest{
		Pagination: &common.Pagination{
			Page:  int32(param.Page),
			Limit: int32(param.Limit),
		},
		Filter: &userpb.FilterRequest{},
	}

	if len(param.UIDs) > 0 {
		pbReq.Filter.Uids = param.UIDs
	}
	if param.Username != nil {
		pbReq.Filter.Username = param.Username
	}
	if param.Email != nil {
		pbReq.Filter.Email = param.Email
	}
	if param.Query != nil {
		pbReq.Filter.Query = param.Query
	}
	if param.Active != nil {
		// Convert Active bool to Status int32
		var status int32
		if *param.Active {
			status = model.UserStatusActive
		} else {
			status = model.UserStatusInactive
		}
		pbReq.Filter.Status = &status
	}

	resp, err := a.userClient.List(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return grpcAdapter.UsersFromProto(resp.Items, resp.Meta.Total, resp.Meta.Page, resp.Meta.Limit), nil
}

// Create creates a new user and returns its UID.
func (a *UserClientAdapter) Create(ctx context.Context, param *params.CreateUserParam) (string, error) {
	resp, err := a.userClient.Add(ctx, &userpb.AddRequest{
		Username: param.Username,
		Email:    param.Email,
		Password: param.Password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	return resp.Uid, nil
}

// Update updates an existing user.
func (a *UserClientAdapter) Update(ctx context.Context, param *params.UpdateUserParam) error {
	pbReq := &userpb.UpdateRequest{
		Uid: param.UID,
	}

	if param.Username != nil {
		pbReq.Username = param.Username
	}
	if param.Email != nil {
		pbReq.Email = param.Email
	}
	if param.Password != nil {
		pbReq.Password = param.Password
	}
	if param.Status != nil {
		pbReq.Status = param.Status
	}

	_, err := a.userClient.Update(ctx, pbReq)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete soft-deletes a user by UID.
func (a *UserClientAdapter) Delete(ctx context.Context, uid string) error {
	_, err := a.userClient.Delete(ctx, &userpb.DeleteRequest{Uid: uid})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetProfile retrieves a user's profile.
func (a *UserClientAdapter) GetProfile(ctx context.Context, userUID string) (*model.Profile, error) {
	resp, err := a.userClient.GetProfile(ctx, &userpb.GetProfileRequest{UserUid: userUID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return grpcAdapter.ProfileFromProto(resp), nil
}

// UpdateProfile updates a user's profile.
func (a *UserClientAdapter) UpdateProfile(ctx context.Context, param *params.UpdateUserProfileParam) error {
	attrsProto, err := grpcAdapter.AttributesToProto(param.Attributes)
	if err != nil {
		return fmt.Errorf("failed to convert attributes: %w", err)
	}

	pbReq := &userpb.UpdateProfileRequest{
		UserUid:   param.UserUID,
		FirstName: param.FirstName,
		LastName:  param.LastName,
		Bio:       param.Bio,
		Attributes: attrsProto,
	}

	_, err = a.userClient.UpdateProfile(ctx, pbReq)
	if err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}
	return nil
}

// ListDevice returns a paginated list of user devices.
func (a *UserClientAdapter) ListDevice(ctx context.Context, param *params.ListUserDevicesParam) (*model.Devices, error) {
	pbReq := &userpb.ListDevicesRequest{
		UserUid: param.UserUID,
		Pagination: &common.Pagination{
			Page:  int32(param.Page),
			Limit: int32(param.Limit),
		},
		Filter: &userpb.FilterDeviceRequest{},
	}

	resp, err := a.userClient.ListDevice(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list user devices: %w", err)
	}

	return grpcAdapter.DevicesFromProto(resp.Items, resp.Meta.Total, resp.Meta.Page, resp.Meta.Limit), nil
}

// RevokeDevice revokes a device for a user.
func (a *UserClientAdapter) RevokeDevice(ctx context.Context, userUID, deviceUID string) error {
	_, err := a.userClient.RevokeDevice(ctx, &userpb.RevokeDeviceRequest{
		UserUid:   userUID,
		DeviceUid: deviceUID,
	})
	if err != nil {
		return fmt.Errorf("failed to revoke device: %w", err)
	}
	return nil
}

// === AuthClient Implementation ===

// Authenticate authenticates a user and returns tokens.
func (a *UserClientAdapter) Authenticate(ctx context.Context, param *params.LoginParam) (*model.Tokens, error) {
	pbReq := &authpb.AuthRequest{
		Identifier:     param.Identifier,
		IdentifierType: param.IdentifierType,
		Password:       param.Password,
	}
	if param.DeviceFingerprint != "" {
		pbReq.DeviceFingerprint = &param.DeviceFingerprint
	}
	if param.DeviceName != "" {
		pbReq.DeviceName = &param.DeviceName
	}

	resp, err := a.authClient.Auth(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}
	return grpcAdapter.TokensFromProto(resp), nil
}

// RefreshToken refreshes an access token using a refresh token.
func (a *UserClientAdapter) RefreshToken(ctx context.Context, param *params.RefreshTokenParam) (*model.Tokens, error) {
	resp, err := a.authClient.RefreshToken(ctx, &authpb.RefreshTokenRequest{
		RefreshToken: param.RefreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	return grpcAdapter.TokensFromProto(resp), nil
}

// ValidateToken validates an access token and returns user claims.
func (a *UserClientAdapter) ValidateToken(ctx context.Context, param *params.ValidateTokenParam) (*model.TokenClaims, error) {
	resp, err := a.authClient.ValidateToken(ctx, &authpb.ValidateTokenRequest{
		AccessToken: param.AccessToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}
	return grpcAdapter.TokenClaimsFromProto(resp), nil
}

// VerifyPin verifies a user PIN for sensitive actions.
func (a *UserClientAdapter) VerifyPin(ctx context.Context, param *params.VerifyPinParam) (bool, error) {
	resp, err := a.authClient.VerifyPin(ctx, &authpb.VerifyPinRequest{
		Uid:  param.UID,
		Code: param.Code,
	})
	if err != nil {
		return false, fmt.Errorf("failed to verify PIN: %w", err)
	}
	return resp.Valid, nil
}

// GoogleOAuth initiates Google OAuth flow and returns authorization URL.
func (a *UserClientAdapter) GoogleOAuth(ctx context.Context, redirectURI string) (string, error) {
	resp, err := a.authClient.GoogleOAuth(ctx, &authpb.GoogleOAuthRequest{
		RedirectUri: redirectURI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to initiate Google OAuth: %w", err)
	}
	return resp.AuthorizationUrl, nil
}

// HandleGoogleOAuth handles Google OAuth callback and returns tokens.
func (a *UserClientAdapter) HandleGoogleOAuth(ctx context.Context, code, redirectURI string) (*model.Tokens, error) {
	resp, err := a.authClient.HandleGoogleOAuth(ctx, &authpb.HandleGoogleOAuthRequest{
		Code:        code,
		RedirectUri: redirectURI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to handle Google OAuth callback: %w", err)
	}
	return grpcAdapter.TokensFromProto(resp), nil
}

// GetConn returns the underlying gRPC connection.
func (a *UserClientAdapter) GetConn() *grpc.ClientConn {
	return a.conn
}
