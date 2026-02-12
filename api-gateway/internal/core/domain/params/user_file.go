package params

// ListUserFilesParam contains parameters for listing user files.
type ListUserFilesParam struct {
	Page      int
	Limit     int
	UserUIDs  []string
	FileType  *string
	Visibility *int32
}

// UploadUserFileParam contains parameters for uploading a file.
type UploadUserFileParam struct {
	UserUID    string
	FileName   string
	FileData   []byte
	MimeType   string
	Visibility *int32
}

// CreateUserFileParam contains parameters for creating a file.
type CreateUserFileParam struct {
	UserUID    string
	FileName   string
	FileData   []byte
	MimeType   string
	Visibility *int32
}

// UpdateUserFileParam contains parameters for updating a file.
type UpdateUserFileParam struct {
	UID        string
	FileName   *string
	Visibility *int32
}
