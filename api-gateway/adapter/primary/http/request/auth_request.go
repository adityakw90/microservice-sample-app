package request

// LoginRequest represents the HTTP request for user login.
type LoginRequest struct {
	Identifier        string `json:"identifier" validate:"required"`
	IdentifierType    string `json:"identifier_type" validate:"required,oneof=username email"`
	Password          string `json:"password" validate:"required"`
	DeviceFingerprint string `json:"device_fingerprint,omitempty"`
	DeviceName        string `json:"device_name,omitempty" validate:"max=100"`
}

// RefreshTokenRequest represents the HTTP request for refreshing a token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ValidateTokenRequest represents the HTTP request for validating a token.
type ValidateTokenRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}

// VerifyPinRequest represents the HTTP request for verifying a PIN.
type VerifyPinRequest struct {
	UID  string `json:"uid" validate:"required"`
	Code string `json:"code" validate:"required,min=4,max=8"`
}
