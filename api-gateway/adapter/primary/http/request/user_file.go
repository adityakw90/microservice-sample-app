package request

// UpdateFileRequest represents a file update request.
type UpdateFileRequest struct {
	FileName   *string `json:"file_name" validate:"omitempty,min=1,max=255"`
	Visibility *int32  `json:"visibility" validate:"omitempty,oneof=0 1"`
}
