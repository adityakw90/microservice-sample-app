package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	userfilepb "github.com/adityakw90/service-user-proto/gen/go/user_file"
	"github.com/adityakw90/service-user-proto/gen/go/common"

	grpcAdapter "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/secondary/grpc/mapper"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/port/client"
)

// UserFileClientAdapter implements client.UserFileClient using gRPC.
type UserFileClientAdapter struct {
	client userfilepb.UserFileServiceClient
}

// NewUserFileClientAdapter creates a new UserFileClientAdapter.
func NewUserFileClientAdapter(conn *grpc.ClientConn) client.UserFileClient {
	return &UserFileClientAdapter{
		client: userfilepb.NewUserFileServiceClient(conn),
	}
}

// Get retrieves a file by UID.
func (a *UserFileClientAdapter) Get(ctx context.Context, uid string) (*model.UserFile, error) {
	resp, err := a.client.Get(ctx, &userfilepb.GetRequest{Uid: uid})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	return grpcAdapter.UserFileFromProto(resp), nil
}

// List returns a paginated list of files based on filters.
func (a *UserFileClientAdapter) List(ctx context.Context, param *params.ListUserFilesParam) (*model.UserFiles, error) {
	pbReq := &userfilepb.ListRequest{
		Pagination: &common.Pagination{
			Page:  int32(param.Page),
			Limit: int32(param.Limit),
		},
		Filter: &userfilepb.FilterRequest{},
	}

	if len(param.UserUIDs) > 0 {
		pbReq.Filter.UserUid = param.UserUIDs
	}
	if param.FileType != nil {
		pbReq.Filter.Filetype = param.FileType
	}
	if param.Visibility != nil {
		visibility := "public"
		if *param.Visibility == 0 {
			visibility = "private"
		}
		pbReq.Filter.Visibility = &visibility
	}

	resp, err := a.client.List(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return grpcAdapter.UserFilesFromProto(resp.Items, resp.Meta.Total, resp.Meta.Page, resp.Meta.Limit), nil
}

// Create uploads a new file and returns its UID.
func (a *UserFileClientAdapter) Create(ctx context.Context, param *params.CreateUserFileParam) (string, error) {
	pbReq := &userfilepb.AddRequest{
		UserUid:  param.UserUID,
		Name:     param.FileName,
		Filename: param.FileName,
		Filedata: param.FileData,
	}

	if param.Visibility != nil {
		public := *param.Visibility == 1
		pbReq.Public = &public
	}

	resp, err := a.client.Add(ctx, pbReq)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	return resp.Uid, nil
}

// Update updates an existing file.
func (a *UserFileClientAdapter) Update(ctx context.Context, param *params.UpdateUserFileParam) error {
	pbReq := &userfilepb.UpdateRequest{Uid: param.UID}

	if param.FileName != nil {
		pbReq.Name = param.FileName
	}
	if param.Visibility != nil {
		public := *param.Visibility == 1
		pbReq.Public = &public
	}

	_, err := a.client.Update(ctx, pbReq)
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}
	return nil
}

// Delete deletes a file by UID.
func (a *UserFileClientAdapter) Delete(ctx context.Context, uid string) error {
	_, err := a.client.Delete(ctx, &userfilepb.DeleteRequest{Uid: uid})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
