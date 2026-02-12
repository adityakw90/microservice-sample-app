package params

// LoginParam represents parameters for user login.
type LoginParam struct {
	Identifier        string
	IdentifierType    string // "username" or "email"
	Password          string
	DeviceFingerprint string
	DeviceName        string
}

// RefreshTokenParam represents parameters for refreshing a token.
type RefreshTokenParam struct {
	RefreshToken string
}

// ValidateTokenParam represents parameters for validating a token.
type ValidateTokenParam struct {
	AccessToken string
}

// VerifyPinParam represents parameters for verifying a PIN.
type VerifyPinParam struct {
	UID  string
	Code string
}
