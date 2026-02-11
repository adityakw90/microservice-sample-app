package model_test

import (
	"testing"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
)

// TestProfileHasAvatar tests the HasAvatar method.
func TestProfileHasAvatar(t *testing.T) {
	tests := []struct {
		name     string
		profile  *model.Profile
		expected bool
	}{
		{
			name: "profile with avatar",
			profile: &model.Profile{
				UID:      "user-123",
				FirstName: "John",
				LastName:  "Doe",
				Avatar: &model.Avatar{
					UID: "avatar-123",
					URL: "https://example.com/avatar.jpg",
				},
			},
			expected: true,
		},
		{
			name: "profile without avatar",
			profile: &model.Profile{
				UID:      "user-456",
				FirstName: "Jane",
				LastName:  "Smith",
				Avatar:   nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.profile.HasAvatar(); got != tt.expected {
				t.Errorf("HasAvatar() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestProfileGetFullName tests the GetFullName method.
func TestProfileGetFullName(t *testing.T) {
	tests := []struct {
		name     string
		profile  *model.Profile
		expected string
	}{
		{
			name: "full profile",
			profile: &model.Profile{
				UID:      "user-123",
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "John Doe",
		},
		{
			name: "only first name",
			profile: &model.Profile{
				UID:      "user-456",
				FirstName: "Jane",
				LastName:  "",
			},
			expected: "Jane ",
		},
		{
			name: "only last name",
			profile: &model.Profile{
				UID:      "user-789",
				FirstName: "",
				LastName:  "Smith",
			},
			expected: " Smith",
		},
		{
			name: "no name",
			profile: &model.Profile{
				UID:       "user-000",
				FirstName: "",
				LastName:  "",
			},
			expected: " ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.profile.GetFullName(); got != tt.expected {
				t.Errorf("GetFullName() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestProfileIsEmpty tests the IsEmpty method.
func TestProfileIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		profile  *model.Profile
		expected bool
	}{
		{
			name: "empty profile",
			profile: &model.Profile{
				UID:       "user-123",
				FirstName: "",
				LastName:  "",
				Bio:       "",
				Avatar:    nil,
			},
			expected: true,
		},
		{
			name: "profile with first name",
			profile: &model.Profile{
				UID:       "user-456",
				FirstName: "John",
				LastName:  "",
				Bio:       "",
				Avatar:    nil,
			},
			expected: false,
		},
		{
			name: "profile with last name",
			profile: &model.Profile{
				UID:       "user-789",
				FirstName: "",
				LastName:  "Doe",
				Bio:       "",
				Avatar:    nil,
			},
			expected: false,
		},
		{
			name: "profile with bio",
			profile: &model.Profile{
				UID:       "user-000",
				FirstName: "",
				LastName:  "",
				Bio:       "Some bio",
				Avatar:    nil,
			},
			expected: false,
		},
		{
			name: "profile with avatar",
			profile: &model.Profile{
				UID:       "user-111",
				FirstName: "",
				LastName:  "",
				Bio:       "",
				Avatar:    &model.Avatar{UID: "avatar-1", URL: "http://example.com"},
			},
			expected: false,
		},
		{
			name: "complete profile",
			profile: &model.Profile{
				UID:       "user-222",
				FirstName: "Jane",
				LastName:  "Smith",
				Bio:       "Developer",
				Avatar:    &model.Avatar{UID: "avatar-2", URL: "http://example.com/avatar.jpg"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.profile.IsEmpty(); got != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestAvatarEntity tests the Avatar entity.
func TestAvatarEntity(t *testing.T) {
	avatar := &model.Avatar{
		UID: "avatar-123",
		URL: "https://example.com/avatar.jpg",
	}

	if avatar.UID == "" {
		t.Error("Avatar.UID should not be empty")
	}

	if avatar.URL == "" {
		t.Error("Avatar.URL should not be empty")
	}
}

// TestProfileUID tests that Profile UID is set.
func TestProfileUID(t *testing.T) {
	profile := &model.Profile{
		UID:       "user-123",
		FirstName: "John",
		LastName:  "Doe",
	}

	if profile.UID == "" {
		t.Error("Profile.UID should not be empty")
	}
}
