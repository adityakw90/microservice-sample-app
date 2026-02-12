package request

// CreateUserRequest represents the HTTP request for creating a user.
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// UpdateUserRequest represents the HTTP request for updating a user.
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=8"`
	Status   *int32  `json:"status,omitempty"`
}

// UpdateUserProfileRequest represents the HTTP request for updating a user profile.
type UpdateUserProfileRequest struct {
	FirstName  string              `json:"first_name" validate:"required,max=100"`
	LastName   string              `json:"last_name" validate:"required,max=100"`
	Bio        string              `json:"bio" validate:"max=500"`
	Attributes map[string]any `json:"attributes,omitempty"`
}
