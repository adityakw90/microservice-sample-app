package model

import "time"

// User represents a domain user entity with business logic.
type User struct {
	UID       string
	Username  string
	Email     string
	Status    int32
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// UserStatus constants
const (
	UserStatusActive   int32 = 1
	UserStatusInactive int32 = 0
	UserStatusBanned   int32 = 2
)

// IsActive returns true if the user is active.
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsInactive returns true if the user is inactive.
func (u *User) IsInactive() bool {
	return u.Status == UserStatusInactive
}

// IsBanned returns true if the user is banned.
func (u *User) IsBanned() bool {
	return u.Status == UserStatusBanned
}

// CanDisplay returns true if user can be shown in listings (not deleted).
func (u *User) CanDisplay() bool {
	return u.DeletedAt == nil || u.DeletedAt.IsZero()
}

// CanLogin returns true if user can login (active and not deleted).
func (u *User) CanLogin() bool {
	return u.IsActive() && u.CanDisplay()
}

// Ban marks the user as banned.
func (u *User) Ban() {
	u.Status = UserStatusBanned
}

// Activate marks the user as active.
func (u *User) Activate() {
	u.Status = UserStatusActive
}

// Deactivate marks the user as inactive.
func (u *User) Deactivate() {
	u.Status = UserStatusInactive
}

// FormatTime formats the user's timestamps for JSON serialization.
func (u *User) FormatTime() string {
	return u.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
}

// Users represents a collection of users with pagination metadata.
type Users struct {
	Items []User
	Meta  PaginationMeta
}

// PaginationMeta represents pagination metadata.
type PaginationMeta struct {
	Page  int32
	Limit int32
	Total int64
	Pages int32
}

// Tokens represents authentication tokens.
type Tokens struct {
	AccessToken  string
	RefreshToken string
}

// TokenClaims represents claims from a validated token.
type TokenClaims struct {
	UID            string
	Identifier     string
	IdentifierType string
	Claims         map[string]any
}
