package model

// Profile represents a user's profile domain entity.
type Profile struct {
	UID      string
	FirstName string
	LastName  string
	Bio      string
	Avatar   *Avatar
}

// Avatar represents a user's avatar.
type Avatar struct {
	UID string
	URL string
}

// HasAvatar returns true if the profile has an avatar set.
func (p *Profile) HasAvatar() bool {
	return p.Avatar != nil
}

// GetFullName returns the user's full name.
func (p *Profile) GetFullName() string {
	return p.FirstName + " " + p.LastName
}

// IsEmpty returns true if the profile is essentially empty.
func (p *Profile) IsEmpty() bool {
	return p.FirstName == "" && p.LastName == "" && p.Bio == "" && p.Avatar == nil
}
