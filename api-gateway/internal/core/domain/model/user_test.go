package model_test

import (
	"testing"
	"time"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
)

// TestUserCreationTime validates that a User is properly initialized with creation timestamp.
func TestUserCreationTime(t *testing.T) {
	user := &model.User{
		UID:       "user-123",
		Username:  "testuser",
		Email:     "test@example.com",
		Status:    model.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if user.CreatedAt.IsZero() {
		t.Error("User.CreatedAt should be set during creation")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("User.UpdatedAt should be set during creation")
	}

	if user.DeletedAt != nil {
		t.Error("User.DeletedAt should be nil for active user")
	}

	if user.Status == model.UserStatusInactive {
		t.Error("User.Status should be UserStatusActive for active user")
	}

	if user.UID == "" {
		t.Error("User.UID should not be empty")
	}

	if user.Username == "" {
		t.Error("User.Username should not be empty")
	}

	if user.Email == "" {
		t.Error("User.Email should not be empty")
	}

	if !user.IsActive() {
		t.Error("Active user should be true")
	}

	if !user.CanDisplay() {
		t.Error("Active user should be displayable")
	}
}

// TestInactiveUser tests inactive user behavior.
func TestInactiveUser(t *testing.T) {
	now := time.Now()
	inactiveUser := &model.User{
		UID:       "user-456",
		Username:  "inactive_user",
		Email:     "inactive@example.com",
		Status:    model.UserStatusInactive,
		CreatedAt: now.Add(-48 * time.Hour),
		UpdatedAt: now.Add(-48 * time.Hour),
		DeletedAt: &now,
	}

	if inactiveUser.IsActive() {
		t.Error("Inactive user should not be active")
	}

	if inactiveUser.UID == "" {
		t.Error("Inactive user.UID should not be empty")
	}

	if inactiveUser.Username == "" {
		t.Error("Inactive user.Username should not be empty")
	}

	if inactiveUser.Email == "" {
		t.Error("Inactive user.Email should not be empty")
	}

	if inactiveUser.CanDisplay() {
		t.Error("Inactive user with DeletedAt set should not be displayable")
	}
}

// TestBannedUser tests banned user behavior.
func TestBannedUser(t *testing.T) {
	bannedUser := &model.User{
		UID:       "user-banned",
		Username:  "banned_user",
		Email:     "banned@example.com",
		Status:    model.UserStatusBanned,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if bannedUser.Status != model.UserStatusBanned {
		t.Error("Banned user should have UserStatusBanned")
	}

	if bannedUser.UID == "" {
		t.Error("Banned user.UID should not be empty")
	}

	if bannedUser.Username == "" {
		t.Error("Banned user.Username should not be empty")
	}

	if bannedUser.Email == "" {
		t.Error("Banned user.Email should not be empty")
	}

	// Banned users can be displayed (not deleted), but cannot login
	if !bannedUser.CanDisplay() {
		t.Error("Banned user should be displayable (not deleted)")
	}

	if bannedUser.CanLogin() {
		t.Error("Banned user should not be able to login")
	}
}

// TestUserMethods tests user business logic methods.
func TestUserMethods(t *testing.T) {
	tests := []struct {
		name     string
		user     *model.User
		active   bool
		inactive bool
		banned   bool
		display  bool
		login    bool
	}{
		{
			name: "active user",
			user: &model.User{
				UID:       "user-1",
				Username:  "active",
				Email:     "active@example.com",
				Status:    model.UserStatusActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			active:   true,
			inactive: false,
			banned:   false,
			display:  true,
			login:    true,
		},
		{
			name: "inactive user",
			user: &model.User{
				UID:       "user-2",
				Username:  "inactive",
				Email:     "inactive@example.com",
				Status:    model.UserStatusInactive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			active:   false,
			inactive: true,
			banned:   false,
			display:  true,
			login:    false,
		},
		{
			name: "banned user",
			user: &model.User{
				UID:       "user-3",
				Username:  "banned",
				Email:     "banned@example.com",
				Status:    model.UserStatusBanned,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			active:   false,
			inactive: false,
			banned:   true,
			display:  true,
			login:    false,
		},
		{
			name: "deleted user (soft-deleted but status remains active)",
			user: &model.User{
				UID:       "user-4",
				Username:  "deleted",
				Email:     "deleted@example.com",
				Status:    model.UserStatusActive, // Status stays active even when soft-deleted
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: func() *time.Time { t := time.Now(); return &t }(),
			},
			// Note: IsActive() returns true because it only checks Status field
			// The DeletedAt field is checked separately by CanDisplay() and CanLogin()
			active:   true,
			inactive: false,
			banned:   false,
			display:  false, // CanDisplay() checks DeletedAt
			login:    false, // CanLogin() checks both IsActive() and CanDisplay()
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.IsActive(); got != tt.active {
				t.Errorf("IsActive() = %v, want %v", got, tt.active)
			}
			if got := tt.user.IsInactive(); got != tt.inactive {
				t.Errorf("IsInactive() = %v, want %v", got, tt.inactive)
			}
			if got := tt.user.IsBanned(); got != tt.banned {
				t.Errorf("IsBanned() = %v, want %v", got, tt.banned)
			}
			if got := tt.user.CanDisplay(); got != tt.display {
				t.Errorf("CanDisplay() = %v, want %v", got, tt.display)
			}
			if got := tt.user.CanLogin(); got != tt.login {
				t.Errorf("CanLogin() = %v, want %v", got, tt.login)
			}
		})
	}
}

// TestUserStatusChanges tests user status change methods.
func TestUserStatusChanges(t *testing.T) {
	user := &model.User{
		UID:       "user-123",
		Username:  "testuser",
		Email:     "test@example.com",
		Status:    model.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test Ban
	user.Ban()
	if user.Status != model.UserStatusBanned {
		t.Error("Ban() should set status to UserStatusBanned")
	}

	// Test Activate
	user.Activate()
	if user.Status != model.UserStatusActive {
		t.Error("Activate() should set status to UserStatusActive")
	}

	// Test Deactivate
	user.Deactivate()
	if user.Status != model.UserStatusInactive {
		t.Error("Deactivate() should set status to UserStatusInactive")
	}
}

// TestUsersCollection tests the Users collection type.
func TestUsersCollection(t *testing.T) {
	users := &model.Users{
		Items: []model.User{
			{UID: "user-1", Username: "user1", Email: "user1@example.com", Status: model.UserStatusActive},
			{UID: "user-2", Username: "user2", Email: "user2@example.com", Status: model.UserStatusActive},
		},
		Meta: model.PaginationMeta{
			Page:  1,
			Limit: 10,
			Total: 2,
			Pages: 1,
		},
	}

	if len(users.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(users.Items))
	}

	if users.Meta.Total != 2 {
		t.Errorf("Expected Total=2, got %d", users.Meta.Total)
	}
}
