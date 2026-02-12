package service

import (
	"context"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
)

// UserFileService defines the interface for file-related use cases.
// This is a primary port - what the application provides to the outside world.
type UserFileService interface {
	// GetFile retrieves a file by UID.
	GetFile(ctx context.Context, uid string) (*model.UserFile, error)

	// ListFiles returns a paginated list of files based on filters.
	ListFiles(ctx context.Context, param *params.ListUserFilesParam) (*model.UserFiles, error)

	// UploadFile uploads a new file and returns its UID.
	UploadFile(ctx context.Context, param *params.UploadUserFileParam) (*model.UserFile, error)

	// UpdateFile updates an existing file.
	UpdateFile(ctx context.Context, param *params.UpdateUserFileParam) (*model.UserFile, error)

	// DeleteFile deletes a file by UID.
	DeleteFile(ctx context.Context, uid string) error
}
