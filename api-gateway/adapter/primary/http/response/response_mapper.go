package response

import (
	"time"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
)

// UserResponse represents a user in HTTP responses.
type UserResponse struct {
	UID       string `json:"uid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Status    int32  `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at,omitempty"`
}

// UsersResponse represents a paginated list of users in HTTP responses.
type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Meta  MetaResponse    `json:"meta"`
}

// MetaResponse represents pagination metadata in HTTP responses.
type MetaResponse struct {
	Page  int32 `json:"page"`
	Limit int32 `json:"limit"`
	Total int64 `json:"total"`
	Pages int32 `json:"pages"`
}

// ProfileResponse represents a user profile in HTTP responses.
type ProfileResponse struct {
	UID      string        `json:"uid"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Bio       string       `json:"bio"`
	Avatar    *AvatarResponse `json:"avatar,omitempty"`
}

// AvatarResponse represents an avatar in HTTP responses.
type AvatarResponse struct {
	UID string `json:"uid"`
	URL string `json:"url"`
}

// DeviceResponse represents a device in HTTP responses.
type DeviceResponse struct {
	DeviceUID    string `json:"device_uid"`
	DeviceName   string `json:"device_name"`
	IPAddress    string `json:"ip_address"`
	LastActiveAt string `json:"last_active_at"`
	CreatedAt    string `json:"created_at"`
	RevokedAt    string `json:"revoked_at,omitempty"`
}

// DevicesResponse represents a paginated list of devices in HTTP responses.
type DevicesResponse struct {
	Devices []DeviceResponse `json:"devices"`
	Meta    MetaResponse     `json:"meta"`
}

// TokensResponse represents authentication tokens in HTTP responses.
type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// TokenClaimsResponse represents validated token claims in HTTP responses.
type TokenClaimsResponse struct {
	UID            string                 `json:"uid"`
	Identifier     string                 `json:"identifier"`
	IdentifierType string                 `json:"identifier_type"`
	Claims         map[string]any `json:"claims,omitempty"`
}

// ErrorResponse represents an error in HTTP responses.
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a simple success response.
type SuccessResponse struct {
	Success bool `json:"success"`
}

// ThumbnailResponse represents a file thumbnail in HTTP responses.
type ThumbnailResponse struct {
	UID string `json:"uid"`
	URL string `json:"url"`
}

// UserFileResponse represents a user file in HTTP responses.
type UserFileResponse struct {
	UID        string             `json:"uid"`
	UserUID    string             `json:"user_uid"`
	FileType   string             `json:"file_type"`
	FileName   string             `json:"file_name"`
	FilePath   string             `json:"file_path"`
	MimeType   string             `json:"mime_type"`
	SizeBytes  int64              `json:"size_bytes"`
	Visibility int32              `json:"visibility"`
	CreatedAt  string             `json:"created_at"`
	Thumbnail  *ThumbnailResponse  `json:"thumbnail,omitempty"`
}

// UserFilesResponse represents a paginated list of files in HTTP responses.
type UserFilesResponse struct {
	Files []UserFileResponse `json:"files"`
	Meta  MetaResponse       `json:"meta"`
}

// UserFromDomain converts a domain User to HTTP response.
func UserFromDomain(u *model.User) UserResponse {
	return UserResponse{
		UID:       u.UID,
		Username:  u.Username,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: formatTime(u.CreatedAt),
		UpdatedAt: formatTime(u.UpdatedAt),
		DeletedAt: formatTimePtr(u.DeletedAt),
	}
}

// UsersFromDomain converts domain Users to HTTP response.
func UsersFromDomain(users *model.Users) UsersResponse {
	userResp := make([]UserResponse, len(users.Items))
	for i, u := range users.Items {
		userResp[i] = UserFromDomain(&u)
	}
	return UsersResponse{
		Users: userResp,
		Meta:  MetaFromDomain(users.Meta),
	}
}

// ProfileFromDomain converts a domain Profile to HTTP response.
func ProfileFromDomain(p *model.Profile) ProfileResponse {
	resp := ProfileResponse{
		UID:      p.UID,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Bio:       p.Bio,
	}
	if p.Avatar != nil {
		resp.Avatar = &AvatarResponse{
			UID: p.Avatar.UID,
			URL: p.Avatar.URL,
		}
	}
	return resp
}

// DeviceFromDomain converts a domain Device to HTTP response.
func DeviceFromDomain(d *model.Device) DeviceResponse {
	return DeviceResponse{
		DeviceUID:    d.DeviceUID,
		DeviceName:   d.DeviceName,
		IPAddress:    d.IPAddress,
		LastActiveAt: formatTime(d.LastActiveAt),
		CreatedAt:    formatTime(d.CreatedAt),
		RevokedAt:    formatTimePtr(d.RevokedAt),
	}
}

// DevicesFromDomain converts domain Devices to HTTP response.
func DevicesFromDomain(devices *model.Devices) DevicesResponse {
	deviceResp := make([]DeviceResponse, len(devices.Items))
	for i, d := range devices.Items {
		deviceResp[i] = DeviceFromDomain(&d)
	}
	return DevicesResponse{
		Devices: deviceResp,
		Meta:    MetaFromDomain(devices.Meta),
	}
}

// TokensFromDomain converts domain Tokens to HTTP response.
func TokensFromDomain(tokens *model.Tokens) TokensResponse {
	return TokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

// TokenClaimsFromDomain converts domain TokenClaims to HTTP response.
func TokenClaimsFromDomain(claims *model.TokenClaims) TokenClaimsResponse {
	return TokenClaimsResponse{
		UID:            claims.UID,
		Identifier:     claims.Identifier,
		IdentifierType: claims.IdentifierType,
		Claims:         claims.Claims,
	}
}

// MetaFromDomain converts domain PaginationMeta to HTTP response.
func MetaFromDomain(meta model.PaginationMeta) MetaResponse {
	return MetaResponse{
		Page:  meta.Page,
		Limit: meta.Limit,
		Total: meta.Total,
		Pages: meta.Pages,
	}
}

// Helper functions for time formatting

func formatTime(t time.Time) string {
	return t.Format("2006-01-02T15:04:05Z07:00")
}

func formatTimePtr(tp *time.Time) string {
	if tp == nil || tp.IsZero() {
		return ""
	}
	return formatTime(*tp)
}

// UserFileFromDomain converts domain UserFile to HTTP response.
func UserFileFromDomain(f *model.UserFile) UserFileResponse {
	resp := UserFileResponse{
		UID:        f.UID,
		UserUID:    f.UserUID,
		FileType:   f.FileType,
		FileName:   f.FileName,
		FilePath:   f.FilePath,
		MimeType:   f.MimeType,
		SizeBytes:  f.SizeBytes,
		Visibility: f.Visibility,
		CreatedAt:  formatTime(f.CreatedAt),
	}

	if f.Thumbnail != nil {
		resp.Thumbnail = &ThumbnailResponse{
			UID: f.Thumbnail.UID,
			URL: f.Thumbnail.URL,
		}
	}

	return resp
}

// UserFilesFromDomain converts domain UserFiles to HTTP response.
func UserFilesFromDomain(files *model.UserFiles) UserFilesResponse {
	fileResp := make([]UserFileResponse, len(files.Items))
	for i, f := range files.Items {
		fileResp[i] = UserFileFromDomain(&f)
	}
	return UserFilesResponse{
		Files: fileResp,
		Meta:  MetaFromDomain(files.Meta),
	}
}
