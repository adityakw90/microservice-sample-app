package user_file

import (
	"context"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/port/client"
)

// UserFileApplicationService implements service.UserFileService.
type UserFileApplicationService struct {
	fileClient client.UserFileClient
}

// NewUserFileApplicationService creates a new UserFileApplicationService.
func NewUserFileApplicationService(fileClient client.UserFileClient) *UserFileApplicationService {
	return &UserFileApplicationService{
		fileClient: fileClient,
	}
}

// GetFile retrieves a file by UID.
func (s *UserFileApplicationService) GetFile(ctx context.Context, uid string) (*model.UserFile, error) {
	return s.fileClient.Get(ctx, uid)
}

// ListFiles returns a paginated list of files based on filters.
func (s *UserFileApplicationService) ListFiles(ctx context.Context, param *params.ListUserFilesParam) (*model.UserFiles, error) {
	return s.fileClient.List(ctx, param)
}

// UploadFile uploads a new file and returns its UID.
func (s *UserFileApplicationService) UploadFile(ctx context.Context, param *params.UploadUserFileParam) (*model.UserFile, error) {
	// Convert UploadUserFileParam to CreateUserFileParam for client
	createParam := &params.CreateUserFileParam{
		UserUID:    param.UserUID,
		FileName:   param.FileName,
		FileData:   param.FileData,
		MimeType:   param.MimeType,
		Visibility: param.Visibility,
	}

	uid, err := s.fileClient.Create(ctx, createParam)
	if err != nil {
		return nil, err
	}

	// Fetch the created file to return complete data
	return s.fileClient.Get(ctx, uid)
}

// UpdateFile updates an existing file.
func (s *UserFileApplicationService) UpdateFile(ctx context.Context, param *params.UpdateUserFileParam) (*model.UserFile, error) {
	if err := s.fileClient.Update(ctx, param); err != nil {
		return nil, err
	}

	// Fetch updated file
	return s.fileClient.Get(ctx, param.UID)
}

// DeleteFile deletes a file by UID.
func (s *UserFileApplicationService) DeleteFile(ctx context.Context, uid string) error {
	return s.fileClient.Delete(ctx, uid)
}
