package client

import (
	"context"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
)

// UserFileClient defines the interface for calling user file gRPC service.
// This is a secondary port - what the application needs from external services.
type UserFileClient interface {
	// Get retrieves a file by UID.
	Get(ctx context.Context, uid string) (*model.UserFile, error)

	// List returns a paginated list of files based on filters.
	List(ctx context.Context, param *params.ListUserFilesParam) (*model.UserFiles, error)

	// Create uploads a new file and returns its UID.
	Create(ctx context.Context, param *params.CreateUserFileParam) (string, error)

	// Update updates an existing file.
	Update(ctx context.Context, param *params.UpdateUserFileParam) error

	// Delete deletes a file by UID.
	Delete(ctx context.Context, uid string) error
}
