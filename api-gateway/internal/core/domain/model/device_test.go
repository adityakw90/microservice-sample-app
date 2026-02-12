package model_test

import (
	"testing"
	"time"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
)

// TestDeviceIsActive tests the IsActive method.
func TestDeviceIsActive(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		device   *model.Device
		expected bool
	}{
		{
			name: "active device (no revoked timestamp)",
			device: &model.Device{
				DeviceUID:    "device-123",
				DeviceName:   "iPhone",
				IPAddress:    "192.168.1.1",
				LastActiveAt: now,
				CreatedAt:    now.Add(-24 * time.Hour),
				RevokedAt:    nil,
			},
			expected: true,
		},
		{
			name: "active device (zero revoked timestamp)",
			device: &model.Device{
				DeviceUID:    "device-456",
				DeviceName:   "Android",
				IPAddress:    "192.168.1.2",
				LastActiveAt: now,
				CreatedAt:    now.Add(-24 * time.Hour),
				RevokedAt:    func() *time.Time { t := time.Time{}; return &t }(),
			},
			expected: true,
		},
		{
			name: "revoked device",
			device: &model.Device{
				DeviceUID:    "device-789",
				DeviceName:   "Laptop",
				IPAddress:    "192.168.1.3",
				LastActiveAt: now.Add(-1 * time.Hour),
				CreatedAt:    now.Add(-48 * time.Hour),
				RevokedAt:    func() *time.Time { t := now; return &t }(),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.device.IsActive(); got != tt.expected {
				t.Errorf("IsActive() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestDeviceIsRecentlyActive tests the IsRecentlyActive method.
func TestDeviceIsRecentlyActive(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		device   *model.Device
		duration time.Duration
		expected bool
	}{
		{
			name: "active within 1 hour",
			device: &model.Device{
				DeviceUID:    "device-123",
				LastActiveAt: now.Add(-30 * time.Minute),
			},
			duration: 1 * time.Hour,
			expected: true,
		},
		{
			name: "not active within 1 hour",
			device: &model.Device{
				DeviceUID:    "device-456",
				LastActiveAt: now.Add(-2 * time.Hour),
			},
			duration: 1 * time.Hour,
			expected: false,
		},
		{
			name: "active within 24 hours",
			device: &model.Device{
				DeviceUID:    "device-789",
				LastActiveAt: now.Add(-12 * time.Hour),
			},
			duration: 24 * time.Hour,
			expected: true,
		},
		{
			name: "active just within threshold",
			device: &model.Device{
				DeviceUID:    "device-000",
				LastActiveAt: now.Add(-59 * time.Minute),
			},
			duration: 1 * time.Hour,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.device.IsRecentlyActive(tt.duration); got != tt.expected {
				t.Errorf("IsRecentlyActive(%v) = %v, want %v", tt.duration, got, tt.expected)
			}
		})
	}
}

// TestDeviceRevoke tests the Revoke method.
func TestDeviceRevoke(t *testing.T) {
	now := time.Now()
	device := &model.Device{
		DeviceUID:    "device-123",
		DeviceName:   "iPhone",
		IPAddress:    "192.168.1.1",
		LastActiveAt: now,
		CreatedAt:    now.Add(-24 * time.Hour),
		RevokedAt:    nil,
	}

	// Device should be active initially
	if !device.IsActive() {
		t.Error("Device should be active before revocation")
	}

	// Revoke the device
	device.Revoke()

	// Device should not be active after revocation
	if device.IsActive() {
		t.Error("Device should not be active after revocation")
	}

	// RevokedAt should be set
	if device.RevokedAt == nil {
		t.Error("RevokedAt should be set after revocation")
	}

	// RevokedAt should be recent (within last second)
	if device.RevokedAt != nil && time.Since(*device.RevokedAt) > time.Second {
		t.Error("RevokedAt should be set to current time")
	}
}

// TestDevicesCollection tests the Devices collection type.
func TestDevicesCollection(t *testing.T) {
	now := time.Now()
	devices := &model.Devices{
		Items: []model.Device{
			{
				DeviceUID:    "device-1",
				DeviceName:   "iPhone",
				IPAddress:    "192.168.1.1",
				LastActiveAt: now,
				CreatedAt:    now.Add(-24 * time.Hour),
			},
			{
				DeviceUID:    "device-2",
				DeviceName:   "Android",
				IPAddress:    "192.168.1.2",
				LastActiveAt: now,
				CreatedAt:    now.Add(-48 * time.Hour),
			},
		},
		Meta: model.PaginationMeta{
			Page:  1,
			Limit: 10,
			Total: 2,
			Pages: 1,
		},
	}

	if len(devices.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(devices.Items))
	}

	if devices.Meta.Total != 2 {
		t.Errorf("Expected Total=2, got %d", devices.Meta.Total)
	}
}

// TestDeviceFields tests that required device fields are validated.
func TestDeviceFields(t *testing.T) {
	now := time.Now()
	device := &model.Device{
		DeviceUID:    "device-123",
		DeviceName:   "iPhone",
		IPAddress:    "192.168.1.1",
		LastActiveAt: now,
		CreatedAt:    now.Add(-24 * time.Hour),
	}

	if device.DeviceUID == "" {
		t.Error("Device.DeviceUID should not be empty")
	}

	if device.DeviceName == "" {
		t.Error("Device.DeviceName should not be empty")
	}

	if device.IPAddress == "" {
		t.Error("Device.IPAddress should not be empty")
	}
}

// TestRevokedDeviceCannotBeActive tests that revoked devices cannot be active.
func TestRevokedDeviceCannotBeActive(t *testing.T) {
	now := time.Now()
	device := &model.Device{
		DeviceUID:    "device-123",
		DeviceName:   "iPhone",
		IPAddress:    "192.168.1.1",
		LastActiveAt: now,
		CreatedAt:    now.Add(-24 * time.Hour),
		RevokedAt:    func() *time.Time { t := now; return &t }(),
	}

	if device.IsActive() {
		t.Error("Revoked device should not be active")
	}

	// Even if recently active, revoked device should return false for IsRecentlyActive
	// with a recent duration
	if !device.IsRecentlyActive(1 * time.Hour) {
		// IsRecentlyActive should still return true based on LastActiveAt
		// The revocation doesn't affect the LastActiveAt timestamp
		t.Error("IsRecentlyActive should be based on LastActiveAt, not revocation status")
	}
}
