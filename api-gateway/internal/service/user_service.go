package service

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	userpb "github.com/adityakw90/service-user-proto/gen/go/user"
	"github.com/adityakw90/service-user-proto/gen/go/common"
)

// UserService handles business logic for user operations
type UserService struct {
	client userpb.UserServiceClient
}

// NewUserService creates a new user service
func NewUserService(client userpb.UserServiceClient) *UserService {
	return &UserService{
		client: client,
	}
}

// ListUsersRequest represents the request to list users
type ListUsersRequest struct {
	Page     int
	Limit    int
	UIDs     []string
	Username *string
	Email    *string
	Query    *string
	Active   *bool
}

// ListUsersResponse represents the response from listing users
type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
	Meta  MetaResponse   `json:"meta"`
}

// UserResponse represents a user in HTTP responses
type UserResponse struct {
	UID       string `json:"uid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Status    int32  `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at,omitempty"`
}

// MetaResponse represents pagination metadata
type MetaResponse struct {
	Page  int32 `json:"page"`
	Limit int32 `json:"limit"`
	Total int64 `json:"total"`
	Pages int32 `json:"pages"`
}

// ProfileResponse represents a user profile in HTTP responses
type ProfileResponse struct {
	UID      string        `json:"uid"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Bio       string       `json:"bio"`
	Avatar    *AvatarResponse `json:"avatar,omitempty"`
}

// AvatarResponse represents an avatar in HTTP responses
type AvatarResponse struct {
	UID string `json:"uid"`
	URL string `json:"url"`
}

// DeviceResponse represents a device in HTTP responses
type DeviceResponse struct {
	DeviceUID    string `json:"device_uid"`
	DeviceName   string `json:"device_name"`
	IPAddress    string `json:"ip_address"`
	LastActiveAt string `json:"last_active_at"`
	CreatedAt    string `json:"created_at"`
	RevokedAt    string `json:"revoked_at,omitempty"`
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	UID      string  `json:"uid"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	Status   *int32  `json:"status,omitempty"`
}

// UpdateUserProfileRequest represents the request to update a user profile
type UpdateUserProfileRequest struct {
	UserUID    string  `json:"user_uid"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Bio        string  `json:"bio"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

// ListUserDevicesRequest represents the request to list user devices
type ListUserDevicesRequest struct {
	UserUID string
	Page    int
	Limit   int
}

// ListUserDevicesResponse represents the response from listing user devices
type ListUserDevicesResponse struct {
	Devices []DeviceResponse `json:"devices"`
	Meta    MetaResponse     `json:"meta"`
}

// ListUsers returns a paginated list of users
func (s *UserService) ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error) {
	pbReq := &userpb.ListRequest{
		Pagination: &common.Pagination{
			Page:  int32(req.Page),
			Limit: int32(req.Limit),
		},
		Filter: &userpb.FilterRequest{},
	}

	if len(req.UIDs) > 0 {
		pbReq.Filter.Uids = req.UIDs
	}
	if req.Username != nil {
		pbReq.Filter.Username = req.Username
	}
	if req.Email != nil {
		pbReq.Filter.Email = req.Email
	}
	if req.Query != nil {
		pbReq.Filter.Query = req.Query
	}
	if req.Active != nil {
		pbReq.Filter.Active = req.Active
	}

	pbResp, err := s.client.List(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]UserResponse, len(pbResp.Items))
	for i, u := range pbResp.Items {
		users[i] = userToResponse(u)
	}

	return &ListUsersResponse{
		Users: users,
		Meta:  metaToResponse(pbResp.Meta),
	}, nil
}

// GetUser returns a user by UID
func (s *UserService) GetUser(ctx context.Context, uid string) (*UserResponse, error) {
	pbResp, err := s.client.Get(ctx, &userpb.GetRequest{Uid: uid})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	resp := userToResponse(pbResp)
	return &resp, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (string, error) {
	pbResp, err := s.client.Add(ctx, &userpb.AddRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return pbResp.Uid, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, req UpdateUserRequest) error {
	pbReq := &userpb.UpdateRequest{
		Uid: req.UID,
	}

	if req.Username != nil {
		pbReq.Username = req.Username
	}
	if req.Email != nil {
		pbReq.Email = req.Email
	}
	if req.Password != nil {
		pbReq.Password = req.Password
	}
	if req.Status != nil {
		pbReq.Status = req.Status
	}

	_, err := s.client.Update(ctx, pbReq)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user by UID
func (s *UserService) DeleteUser(ctx context.Context, uid string) error {
	_, err := s.client.Delete(ctx, &userpb.DeleteRequest{Uid: uid})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// GetUserProfile returns a user's profile
func (s *UserService) GetUserProfile(ctx context.Context, userUID string) (*ProfileResponse, error) {
	pbResp, err := s.client.GetProfile(ctx, &userpb.GetProfileRequest{UserUid: userUID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	resp := &ProfileResponse{
		UID:      pbResp.Uid,
		FirstName: pbResp.FirstName,
		LastName:  pbResp.LastName,
		Bio:       pbResp.Bio,
	}

	if pbResp.Avatar != nil {
		resp.Avatar = &AvatarResponse{
			UID: pbResp.Avatar.Uid,
			URL: pbResp.Avatar.Url,
		}
	}

	return resp, nil
}

// UpdateUserProfile updates a user's profile
func (s *UserService) UpdateUserProfile(ctx context.Context, req UpdateUserProfileRequest) error {
	pbReq := &userpb.UpdateProfileRequest{
		UserUid:   req.UserUID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Bio:       req.Bio,
		Attributes: &structpb.Struct{
			Fields: make(map[string]*structpb.Value),
		},
	}

	// Convert attributes map to protobuf struct
	for k, v := range req.Attributes {
		var val *structpb.Value
		switch converted := v.(type) {
		case string:
			val = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: converted}}
		case bool:
			val = &structpb.Value{Kind: &structpb.Value_BoolValue{BoolValue: converted}}
		case float64:
			val = &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: converted}}
		case int:
			val = &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(converted)}}
		case int64:
			val = &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(converted)}}
		case map[string]interface{}:
			// Convert nested map to struct
			nestedFields := make(map[string]*structpb.Value)
			for nk, nv := range converted {
				var nestedVal *structpb.Value
				switch nconverted := nv.(type) {
				case string:
					nestedVal = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: nconverted}}
				case bool:
					nestedVal = &structpb.Value{Kind: &structpb.Value_BoolValue{BoolValue: nconverted}}
				case float64:
					nestedVal = &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: nconverted}}
				default:
					// For unsupported types, convert to string
					nestedVal = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: fmt.Sprintf("%v", nv)}}
				}
				nestedFields[nk] = nestedVal
			}
			val = &structpb.Value{Kind: &structpb.Value_StructValue{StructValue: &structpb.Struct{Fields: nestedFields}}}
		default:
			// For unsupported types, convert to string
			val = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: fmt.Sprintf("%v", v)}}
		}
		pbReq.Attributes.Fields[k] = val
	}

	_, err := s.client.UpdateProfile(ctx, pbReq)
	if err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	return nil
}

// ListUserDevices returns a list of devices for a user
func (s *UserService) ListUserDevices(ctx context.Context, req ListUserDevicesRequest) (*ListUserDevicesResponse, error) {
	pbReq := &userpb.ListDevicesRequest{
		UserUid: req.UserUID,
		Pagination: &common.Pagination{
			Page:  int32(req.Page),
			Limit: int32(req.Limit),
		},
		Filter: &userpb.FilterDeviceRequest{},
	}

	pbResp, err := s.client.ListDevice(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list user devices: %w", err)
	}

	devices := make([]DeviceResponse, len(pbResp.Items))
	for i, d := range pbResp.Items {
		devices[i] = DeviceResponse{
			DeviceUID:    d.DeviceUid,
			DeviceName:   d.DeviceName,
			IPAddress:    d.IpAddress,
			LastActiveAt: timestampToString(d.LastActiveAt),
			CreatedAt:    timestampToString(d.CreatedAt),
			RevokedAt:    timestampToString(d.RevokedAt),
		}
	}

	return &ListUserDevicesResponse{
		Devices: devices,
		Meta:    metaToResponse(pbResp.Meta),
	}, nil
}

// RevokeUserDevice revokes a device for a user
func (s *UserService) RevokeUserDevice(ctx context.Context, userUID, deviceUID string) error {
	_, err := s.client.RevokeDevice(ctx, &userpb.RevokeDeviceRequest{
		UserUid:   userUID,
		DeviceUid: deviceUID,
	})
	if err != nil {
		return fmt.Errorf("failed to revoke device: %w", err)
	}

	return nil
}

// Helper functions

func userToResponse(u *userpb.User) UserResponse {
	return UserResponse{
		UID:       u.Uid,
		Username:  u.Username,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: timestampToString(u.CreatedAt),
		UpdatedAt: timestampToString(u.UpdatedAt),
		DeletedAt: timestampToString(u.DeletedAt),
	}
}

func metaToResponse(m *common.Meta) MetaResponse {
	return MetaResponse{
		Page:  m.Page,
		Limit: m.Limit,
		Total: m.Total,
		Pages: m.Pages,
	}
}

func timestampToString(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().Format("2006-01-02T15:04:05Z07:00")
}
