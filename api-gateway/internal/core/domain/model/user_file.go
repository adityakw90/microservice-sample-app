package model

import "time"

// UserFileVisibility represents file visibility levels.
const (
	UserFileVisibilityPrivate = 0
	UserFileVisibilityPublic  = 1
)

// UserFile represents a user's file.
type UserFile struct {
	UID        string
	UserUID    string
	FileType   string
	FileName   string
	FilePath   string
	MimeType   string
	SizeBytes  int64
	Visibility int32
	CreatedAt  time.Time
	Thumbnail  *Thumbnail
}

// Thumbnail represents a file thumbnail.
type Thumbnail struct {
	UID  string
	URL  string
}

// UserFiles represents a paginated list of files.
type UserFiles struct {
	Items []UserFile
	Meta  PaginationMeta
}
