package model

import "time"

// Device represents a user's device domain entity.
type Device struct {
	DeviceUID    string
	DeviceName   string
	IPAddress    string
	LastActiveAt time.Time
	CreatedAt    time.Time
	RevokedAt    *time.Time
}

// IsActive returns true if the device is active (not revoked).
func (d *Device) IsActive() bool {
	return d.RevokedAt == nil || d.RevokedAt.IsZero()
}

// IsRecentlyActive returns true if the device was active within the specified duration.
func (d *Device) IsRecentlyActive(duration time.Duration) bool {
	return time.Since(d.LastActiveAt) <= duration
}

// Revoke marks the device as revoked.
func (d *Device) Revoke() {
	now := time.Now()
	d.RevokedAt = &now
}

// Devices represents a collection of devices with pagination metadata.
type Devices struct {
	Items []Device
	Meta  PaginationMeta
}
