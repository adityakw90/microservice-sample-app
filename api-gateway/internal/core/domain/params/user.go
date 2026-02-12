package params

// ListUsersParam represents parameters for listing users.
type ListUsersParam struct {
	ListParam
	UIDs     []string
	Username *string
	Email    *string
	Query    *string
	Active   *bool
}

// GetUserParam represents parameters for getting a user.
type GetUserParam struct {
	UID string
}

// CreateUserParam represents parameters for creating a user.
type CreateUserParam struct {
	Username string
	Email    string
	Password string
}

// UpdateUserParam represents parameters for updating a user.
type UpdateUserParam struct {
	UID      string
	Username *string
	Email    *string
	Password *string
	Status   *int32
}

// DeleteUserParam represents parameters for deleting a user.
type DeleteUserParam struct {
	UID string
}

// GetUserProfileParam represents parameters for getting a user profile.
type GetUserProfileParam struct {
	UserUID string
}

// UpdateUserProfileParam represents parameters for updating a user profile.
type UpdateUserProfileParam struct {
	UserUID    string
	FirstName  string
	LastName   string
	Bio        string
	Attributes map[string]any
}

// ListUserDevicesParam represents parameters for listing user devices.
type ListUserDevicesParam struct {
	ListParam
	UserUID string
}

// RevokeUserDeviceParam represents parameters for revoking a user device.
type RevokeUserDeviceParam struct {
	UserUID   string
	DeviceUID string
}
