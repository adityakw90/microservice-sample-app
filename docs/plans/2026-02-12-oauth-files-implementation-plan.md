# OAuth and File Management Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Integrate Google OAuth and UserFileService from service-user into microservice-sample-app with complete backend and frontend implementation.

**Architecture:** Gateway-centric OAuth flow (service-user owns all OAuth config) + File upload streaming through gateway to service-user gRPC.

**Tech Stack:**
- Backend: Go 1.23, gRPC, gorilla/mux, hexagonal architecture
- Frontend: Vue 3, TypeScript, Pinia, Axios

---

## Task 1: Extend AuthClient Interface for OAuth

**Files:**
- Modify: `api-gateway/internal/core/port/client/auth_client.go`

**Step 1: Add OAuth methods to AuthClient interface**

```go
// GoogleOAuth initiates Google OAuth flow and returns authorization URL.
GoogleOAuth(ctx context.Context, redirectURI string) (authorizationURL string, err error)

// HandleGoogleOAuth handles Google OAuth callback and returns tokens.
HandleGoogleOAuth(ctx context.Context, code, redirectURI string) (*model.Tokens, error)
```

**Step 2: Run tests (if available)**

Run: `cd api-gateway && go test ./...`
Expected: May fail due to proto dependency (pre-existing)

**Step 3: Commit**

```bash
git add api-gateway/internal/core/port/client/auth_client.go
git commit -m "feat(auth): add OAuth methods to AuthClient interface"
```

---

## Task 2: Extend AuthService Interface for OAuth

**Files:**
- Modify: `api-gateway/internal/core/port/service/auth.go`

**Step 1: Add OAuth methods to AuthService interface**

```go
// GoogleOAuth initiates Google OAuth flow and returns authorization URL.
GoogleOAuth(ctx context.Context, redirectURI string) (authorizationURL string, err error)

// HandleGoogleOAuth handles Google OAuth callback and returns tokens.
HandleGoogleOAuth(ctx context.Context, code, redirectURI string) (*model.Tokens, error)
```

**Step 2: Commit**

```bash
git add api-gateway/internal/core/port/service/auth.go
git commit -m "feat(auth): add OAuth methods to AuthService interface"
```

---

## Task 3: Implement OAuth in gRPC Client Adapter

**Files:**
- Modify: `api-gateway/adapter/secondary/grpc/user_client.go`

**Step 1: Add GoogleOAuth method**

Add after line 279 (after VerifyPin):

```go
// GoogleOAuth initiates Google OAuth flow and returns authorization URL.
func (a *UserClientAdapter) GoogleOAuth(ctx context.Context, redirectURI string) (string, error) {
	resp, err := a.authClient.GoogleOAuth(ctx, &authpb.GoogleOAuthRequest{
		RedirectUri: redirectURI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to initiate Google OAuth: %w", err)
	}
	return resp.AuthorizationUrl, nil
}
```

**Step 2: Add HandleGoogleOAuth method**

```go
// HandleGoogleOAuth handles Google OAuth callback and returns tokens.
func (a *UserClientAdapter) HandleGoogleOAuth(ctx context.Context, code, redirectURI string) (*model.Tokens, error) {
	resp, err := a.authClient.HandleGoogleOAuth(ctx, &authpb.HandleGoogleOAuthRequest{
		Code:        code,
		RedirectUri: redirectURI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to handle Google OAuth callback: %w", err)
	}
	return grpcAdapter.TokensFromProto(resp), nil
}
```

**Step 3: Commit**

```bash
git add api-gateway/adapter/secondary/grpc/user_client.go
git commit -m "feat(auth): implement OAuth methods in gRPC client adapter"
```

---

## Task 4: Implement OAuth in Auth Application Service

**Files:**
- Modify: `api-gateway/internal/application/auth/service.go`

**Step 1: Add GoogleOAuth method**

Find the auth service implementation and add:

```go
// GoogleOAuth initiates Google OAuth flow and returns authorization URL.
func (s *authApplicationService) GoogleOAuth(ctx context.Context, redirectURI string) (string, error) {
	return s.authClient.GoogleOAuth(ctx, redirectURI)
}
```

**Step 2: Add HandleGoogleOAuth method**

```go
// HandleGoogleOAuth handles Google OAuth callback and returns tokens.
func (s *authApplicationService) HandleGoogleOAuth(ctx context.Context, code, redirectURI string) (*model.Tokens, error) {
	return s.authClient.HandleGoogleOAuth(ctx, code, redirectURI)
}
```

**Step 3: Commit**

```bash
git add api-gateway/internal/application/auth/service.go
git commit -m "feat(auth): implement OAuth methods in application service"
```

---

## Task 5: Add OAuth HTTP Handlers

**Files:**
- Modify: `api-gateway/adapter/primary/http/httpHandler/auth_handler.go`

**Step 1: Add OAuth routes to RegisterRoutes**

After line 35 (`verify-pin` route):

```go
r.HandleFunc("/auth/google", h.GoogleOAuth).Methods(http.MethodGet)
r.HandleFunc("/auth/google/callback", h.GoogleOAuthCallback).Methods(http.MethodGet)
```

**Step 2: Add GoogleOAuth handler method**

Add after line 157 (after VerifyPin method):

```go
// GoogleOAuth handles GET /api/v1/auth/google.
func (h *AuthHandler) GoogleOAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Build redirect URI from request
	redirectURI := fmt.Sprintf("%s://%s/api/v1/auth/google/callback",
		getScheme(r), r.Host)

	// Get authorization URL from service
	authURL, err := h.authService.GoogleOAuth(ctx, redirectURI)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to initiate OAuth")
		return
	}

	// Redirect to Google OAuth
	http.Redirect(w, r, authURL, http.StatusFound)
}

// GoogleOAuthCallback handles GET /api/v1/auth/google/callback.
func (h *AuthHandler) GoogleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get code and state from query params
	code := r.URL.Query().Get("code")
	if code == "" {
		respondError(w, http.StatusBadRequest, "Missing authorization code")
		return
	}

	// Build redirect URI
	redirectURI := fmt.Sprintf("%s://%s/api/v1/auth/google/callback",
		getScheme(r), r.Host)

	// Handle OAuth callback
	tokens, err := h.authService.HandleGoogleOAuth(ctx, code, redirectURI)
	if err != nil {
		// Redirect to login with error
		http.Redirect(w, r, "/login?error=oauth_failed", http.StatusFound)
		return
	}

	// Get frontend URL from config or use default
	frontendURL := os.Getenv("FRONTEND_REDIRECT_URI")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	// Redirect to frontend with tokens
	targetURL := fmt.Sprintf("%s/login?token=%s&refresh=%s",
		frontendURL, tokens.AccessToken, tokens.RefreshToken)
	http.Redirect(w, r, targetURL, http.StatusFound)
}
```

**Step 3: Add helper function**

Add at the end of the file:

```go
// getScheme determines the request scheme (http or https).
func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http"
}
```

**Step 4: Add required imports**

Add to imports section:
```go
"fmt"
"os"
```

**Step 5: Commit**

```bash
git add api-gateway/adapter/primary/http/httpHandler/auth_handler.go
git commit -m "feat(auth): add OAuth HTTP handlers"
```

---

## Task 6: Update Gateway Config for OAuth

**Files:**
- Create: `api-gateway/.env.example`
- Modify: `api-gateway/internal/config/config.go`

**Step 1: Add FRONTEND_REDIRECT_URI to config**

In `config.go`, add to Config struct and Load function:

```go
type Config struct {
	// ... existing fields
	FrontendRedirectURI string
}

func Load() *Config {
	// ... existing code
	cfg.FrontendRedirectURI = getEnv("FRONTEND_REDIRECT_URI", "http://localhost:3000")
	return cfg
}
```

**Step 2: Update .env.example**

Add:
```
FRONTEND_REDIRECT_URI=http://localhost:3000
```

**Step 3: Commit**

```bash
git add api-gateway/internal/config/config.go api-gateway/.env.example
git commit -m "feat(auth): add frontend redirect URI config"
```

---

## Task 7: Create UserFileClient Interface

**Files:**
- Create: `api-gateway/internal/core/port/client/user_file_client.go`

**Step 1: Create UserFileClient interface**

```go
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
```

**Step 2: Commit**

```bash
git add api-gateway/internal/core/port/client/user_file_client.go
git commit -m "feat(files): add UserFileClient interface"
```

---

## Task 8: Create UserFileService Interface

**Files:**
- Create: `api-gateway/internal/core/port/service/user_file.go`

**Step 1: Create UserFileService interface**

```go
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
```

**Step 2: Commit**

```bash
git add api-gateway/internal/core/port/service/user_file.go
git commit -m "feat(files): add UserFileService interface"
```

---

## Task 9: Create UserFile Domain Models

**Files:**
- Create: `api-gateway/internal/core/domain/model/user_file.go`
- Create: `api-gateway/internal/core/domain/params/user_file.go`

**Step 1: Create user_file.go model**

```go
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
```

**Step 2: Create user_file.go params**

```go
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

// UpdateUserFileParam contains parameters for updating a file.
type UpdateUserFileParam struct {
	UID        string
	FileName   *string
	Visibility *int32
}
```

**Step 3: Commit**

```bash
git add api-gateway/internal/core/domain/model/user_file.go api-gateway/internal/core/domain/params/user_file.go
git commit -m "feat(files): add UserFile domain models and params"
```

---

## Task 10: Create UserFile gRPC Client Adapter

**Files:**
- Create: `api-gateway/adapter/secondary/grpc/user_file_client.go`
- Create: `api-gateway/adapter/secondary/grpc/mapper/user_file_mapper.go`

**Step 1: Create user_file_mapper.go**

```go
package mapper

import (
	"time"

	"github.com/adityakw90/service-user-proto/gen/go/user_file"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
)

// UserFileFromProto converts proto UserFile to domain model.
func UserFileFromProto(pb *user_file.UserFile) *model.UserFile {
	f := &model.UserFile{
		UID:        pb.Uid,
		UserUID:    pb.UserUid,
		FileType:   pb.FileType,
		FileName:   pb.FileName,
		FilePath:   pb.FilePath,
		MimeType:   pb.MimeType,
		SizeBytes:  pb.SizeBytes,
		Visibility: pb.Visibility,
		CreatedAt:  time.Time{},
	}

	if pb.CreatedAt != nil {
		f.CreatedAt = pb.CreatedAt.AsTime()
	}

	if pb.Thumbnail != nil {
		f.Thumbnail = &model.Thumbnail{
			UID: pb.Thumbnail.GetStructValue().GetFields()["uid"].GetStringValue(),
			URL: pb.Thumbnail.GetStructValue().GetFields()["url"].GetStringValue(),
		}
	}

	return f
}

// UserFilesFromProto converts proto list response to domain model.
func UserFilesFromProto(items []*user_file.UserFile, total int64, page, limit int32) *model.UserFiles {
	files := make([]model.UserFile, len(items))
	for i, item := range items {
		files[i] = *UserFileFromProto(item)
	}

	return &model.UserFiles{
		Items: files,
		Meta: model.PaginationMeta{
			Total: total,
			Page:  page,
			Limit: limit,
			Pages: calculatePages(total, limit),
		},
	}
}

func calculatePages(total int64, limit int32) int32 {
	if limit <= 0 {
		return 0
	}
	pages := int32(total / int64(limit))
	if total%int64(limit) > 0 {
		pages++
	}
	return pages
}
```

**Step 2: Create user_file_client.go**

```go
package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	userfilepb "github.com/adityakw90/service-user-proto/gen/go/user_file"
	"github.com/adityakw90/service-user-proto/gen/go/common"

	grpcAdapter "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/secondary/grpc/mapper"
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
		visibility := *param.Visibility == 1
		pbReq.Filter.Public = &visibility
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
```

**Step 3: Fix import in user_file_client.go**

Add missing import at line 17:
```go
"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
```

**Step 4: Commit**

```bash
git add api-gateway/adapter/secondary/grpc/
git commit -m "feat(files): add UserFile gRPC client adapter and mapper"
```

---

## Task 11: Create UserFile Application Service

**Files:**
- Create: `api-gateway/internal/application/user_file/service.go`

**Step 1: Create user_file application service**

```go
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
```

**Step 2: Commit**

```bash
git add api-gateway/internal/application/user_file/service.go
git commit -m "feat(files): add UserFile application service"
```

---

## Task 12: Create File HTTP Response DTOs

**Files:**
- Modify: `api-gateway/adapter/primary/http/response/response_mapper.go`

**Step 1: Add file response types**

Add after line 87 (SuccessResponse):

```go
// UserFileResponse represents a user file in HTTP responses.
type UserFileResponse struct {
	UID        string          `json:"uid"`
	UserUID    string          `json:"user_uid"`
	FileType   string          `json:"file_type"`
	FileName   string          `json:"file_name"`
	FilePath   string          `json:"file_path"`
	MimeType   string          `json:"mime_type"`
	SizeBytes  int64           `json:"size_bytes"`
	Visibility int32           `json:"visibility"`
	CreatedAt  string          `json:"created_at"`
	Thumbnail  *ThumbnailResponse `json:"thumbnail,omitempty"`
}

// UserFilesResponse represents a paginated list of files in HTTP responses.
type UserFilesResponse struct {
	Files []UserFileResponse `json:"files"`
	Meta  MetaResponse       `json:"meta"`
}
```

**Step 2: Add file mapper functions**

Add after line 194 (end of file):

```go
// UserFileFromDomain converts domain UserFile to HTTP response.
func UserFileFromDomain(f *model.UserFile) UserFileResponse {
	resp := UserFileResponse{
		UID:        f.UID,
		UserUID:    f.UserUID,
		FileType:   f.FileType,
		FileName:   f.FileName,
		FilePath:   f.FilePath,
		MimeType:   f.MimeType,
		SizeBytes:  f.SizeBytes,
		Visibility: f.Visibility,
		CreatedAt:  formatTime(f.CreatedAt),
	}

	if f.Thumbnail != nil {
		resp.Thumbnail = &ThumbnailResponse{
			UID: f.Thumbnail.UID,
			URL: f.Thumbnail.URL,
		}
	}

	return resp
}

// UserFilesFromDomain converts domain UserFiles to HTTP response.
func UserFilesFromDomain(files *model.UserFiles) UserFilesResponse {
	fileResp := make([]UserFileResponse, len(files.Items))
	for i, f := range files.Items {
		fileResp[i] = UserFileFromDomain(&f)
	}
	return UserFilesResponse{
		Files: fileResp,
		Meta:  MetaFromDomain(files.Meta),
	}
}
```

**Step 3: Add import for model**

Add to imports if not present:
```go
model "github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
```

**Step 4: Commit**

```bash
git add api-gateway/adapter/primary/http/response/response_mapper.go
git commit -m "feat(files): add file HTTP response DTOs and mappers"
```

---

## Task 13: Create File HTTP Request DTOs

**Files:**
- Create: `api-gateway/adapter/primary/http/request/user_file.go`

**Step 1: Create file request DTOs**

```go
package request

// UpdateFileRequest represents a file update request.
type UpdateFileRequest struct {
	FileName   *string `json:"file_name" validate:"omitempty,min=1,max=255"`
	Visibility *int32  `json:"visibility" validate:"omitempty,oneof=0 1"`
}
```

**Step 2: Commit**

```bash
git add api-gateway/adapter/primary/http/request/user_file.go
git commit -m "feat(files): add file HTTP request DTOs"
```

---

## Task 14: Create File HTTP Handler

**Files:**
- Create: `api-gateway/adapter/primary/http/httpHandler/user_file_handler.go`

**Step 1: Create file handler**

```go
package httpHandler

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
	httpRequest "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/request"
	httpResponse "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/response"
	httpValidator "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/validator"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/port/service"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
)

// UserFileHandler handles HTTP requests for file operations.
type UserFileHandler struct {
	fileService service.UserFileService
	validator   *httpValidator.Validator
}

// NewUserFileHandler creates a new file handler.
func NewUserFileHandler(fileService service.UserFileService) *UserFileHandler {
	return &UserFileHandler{
		fileService: fileService,
		validator:   httpValidator.NewValidator(),
	}
}

// RegisterRoutes registers file routes on the given router.
func (h *UserFileHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/files", h.ListFiles).Methods(http.MethodGet)
	r.HandleFunc("/files", h.UploadFile).Methods(http.MethodPost)
	r.HandleFunc("/files/{uid}", h.GetFile).Methods(http.MethodGet)
	r.HandleFunc("/files/{uid}", h.UpdateFile).Methods(http.MethodPut, http.MethodPatch)
	r.HandleFunc("/files/{uid}", h.DeleteFile).Methods(http.MethodDelete)
}

// ListFiles handles GET /api/v1/files.
func (h *UserFileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 20
	}

	listParam := params.ListUserFilesParam{
		Page:  page,
		Limit: limit,
	}

	if userUID := r.URL.Query().Get("user_uid"); userUID != "" {
		listParam.UserUIDs = []string{userUID}
	}
	if fileType := r.URL.Query().Get("file_type"); fileType != "" {
		listParam.FileType = &fileType
	}

	files, err := h.fileService.ListFiles(ctx, &listParam)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.UserFilesFromDomain(files))
}

// UploadFile handles POST /api/v1/files.
func (h *UserFileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		respondError(w, http.StatusRequestEntityTooLarge, "File too large (max 10MB)")
		return
	}

	// Get file from form
	file, handler, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	// Read file data
	fileData, err := io.ReadAll(file)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to read file")
		return
	}

	// Get user_uid from form or JWT
	userUID := r.FormValue("user_uid")
	if userUID == "" {
		// TODO: Get from JWT context
		respondError(w, http.StatusBadRequest, "user_uid required")
		return
	}

	// Get visibility (default private)
	visibility := int32(0) // private
	if visStr := r.FormValue("visibility"); visStr == "1" || visStr == "true" {
		visibility = 1 // public
	}

	uploadParam := params.UploadUserFileParam{
		UserUID:    userUID,
		FileName:   handler.Filename,
		FileData:   fileData,
		MimeType:   handler.Header.Get("Content-Type"),
		Visibility: &visibility,
	}

	newFile, err := h.fileService.UploadFile(ctx, &uploadParam)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, httpResponse.UserFileFromDomain(newFile))
}

// GetFile handles GET /api/v1/files/{uid}.
func (h *UserFileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	uid := vars["uid"]

	if uid == "" {
		respondError(w, http.StatusBadRequest, "File UID required")
		return
	}

	file, err := h.fileService.GetFile(ctx, uid)
	if err != nil {
		respondError(w, http.StatusNotFound, "File not found")
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.UserFileFromDomain(file))
}

// UpdateFile handles PUT/PATCH /api/v1/files/{uid}.
func (h *UserFileHandler) UpdateFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	uid := vars["uid"]

	if uid == "" {
		respondError(w, http.StatusBadRequest, "File UID required")
		return
	}

	var req httpRequest.UpdateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.ValidateStruct(req); len(errors) > 0 {
		respondValidationError(w, errors)
		return
	}

	updateParam := params.UpdateUserFileParam{
		UID:        uid,
		FileName:   req.FileName,
		Visibility: req.Visibility,
	}

	updatedFile, err := h.fileService.UpdateFile(ctx, &updateParam)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.UserFileFromDomain(updatedFile))
}

// DeleteFile handles DELETE /api/v1/files/{uid}.
func (h *UserFileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	uid := vars["uid"]

	if uid == "" {
		respondError(w, http.StatusBadRequest, "File UID required")
		return
	}

	if err := h.fileService.DeleteFile(ctx, uid); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.SuccessResponse{Success: true})
}
```

**Step 2: Commit**

```bash
git add api-gateway/adapter/primary/http/httpHandler/user_file_handler.go
git commit -m "feat(files): add file HTTP handler"
```

---

## Task 15: Update Gateway main.go to Wire File Handlers

**Files:**
- Modify: `api-gateway/cmd/main.go`

**Step 1: Add imports**

Add after line 19:
```go
userFileHandler "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/httpHandler"
userFileApp "github.com/adityakw90/microservice-sample-app/api-gateway/internal/application/user_file"
```

**Step 2: Create UserFile gRPC adapter**

After line 36 (after authGrpcAdapter creation), add:
```go
// Get gRPC connection for user file service
userFileGrpcAdapter := grpcAdapter.NewUserFileClientAdapter(userGrpcAdapter.(*grpcAdapter.UserClientAdapter).GetConn())
```

**Step 3: Create UserFile application service**

After line 43 (after authAppService creation), add:
```go
userFileAppService := userFileApp.NewUserFileApplicationService(userFileGrpcAdapter)
```

**Step 4: Create UserFile HTTP handler**

After line 50 (after authHdlr creation), add:
```go
userFileHdlr := userFileHandler.NewUserFileHandler(userFileAppService)
```

**Step 5: Register UserFile routes**

After line 61 (after authHdlr.RegisterRoutes), add:
```go
userFileHdlr.RegisterRoutes(api)
```

**Step 6: Add helper method to UserClientAdapter**

Add to `user_client.go`:
```go
// GetConn returns the underlying gRPC connection.
func (a *UserClientAdapter) GetConn() *grpc.ClientConn {
	return a.conn
}
```

**Step 7: Commit**

```bash
git add api-gateway/cmd/main.go api-gateway/adapter/secondary/grpc/user_client.go
git commit -m "feat(files): wire UserFile handlers in gateway"
```

---

## Task 16: Add Google OAuth Button to LoginView

**Files:**
- Modify: `frontend/src/views/LoginView.vue`

**Step 1: Add Google OAuth button to template**

Add after line 43 (after login button), before `</form>`:

```vue
<div class="oauth-divider">
  <span>or</span>
</div>

<button
  type="button"
  @click="handleGoogleOAuth"
  class="google-btn"
>
  <svg class="google-icon" viewBox="0 0 24 24">
    <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
    <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
    <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
    <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
  </svg>
  Continue with Google
</button>
```

**Step 2: Add handleGoogleOAuth method to script**

Add after line 117 (after handleLogin function):

```typescript
const handleGoogleOAuth = () => {
  // Simply redirect to gateway OAuth endpoint
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || '/api/v1'
  window.location.href = `${apiBaseUrl}/auth/google`
}
```

**Step 3: Add OAuth callback handling**

Add to onMounted (create if not present after handleGoogleOAuth):

```typescript
import { onMounted } from 'vue'

// Add at the top with other imports

onMounted(() => {
  // Check for OAuth callback tokens in URL
  const urlParams = new URLSearchParams(window.location.search)
  const token = urlParams.get('token')
  const refreshToken = urlParams.get('refresh')

  if (token && refreshToken) {
    // Store tokens from OAuth callback
    localStorage.setItem('access_token', token)
    localStorage.setItem('refresh_token', refreshToken)

    // Update store
    accessToken.value = token
    refreshToken.value = refreshToken

    // Validate token to get user info
    validateToken().then(() => {
      // Clean URL and redirect to profile
      window.history.replaceState({}, '', '/login')
      router.push('/profile')
    })
  }
})
```

**Step 4: Add OAuth button styles**

Add to `<style scoped>`, after line 292 (end of existing styles):

```css
.oauth-divider {
  display: flex;
  align-items: center;
  margin: 1.5rem 0;
  color: #7f8c8d;
}

.oauth-divider::before,
.oauth-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: #ddd;
}

.oauth-divider span {
  padding: 0 1rem;
  font-size: 0.875rem;
}

.google-btn {
  width: 100%;
  padding: 0.875rem;
  background: white;
  color: #2c3e50;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  transition: background-color 0.2s, border-color 0.2s;
}

.google-btn:hover {
  background: #f8f9fa;
  border-color: #ccc;
}

.google-icon {
  width: 20px;
  height: 20px;
}

.google-icon path:nth-child(1) {
  fill: #4285F4;
}

.google-icon path:nth-child(2) {
  fill: #34A853;
}

.google-icon path:nth-child(3) {
  fill: #FBBC05;
}

.google-icon path:nth-child(4) {
  fill: #EA4335;
}
```

**Step 5: Commit**

```bash
git add frontend/src/views/LoginView.vue
git commit -m "feat(frontend): add Google OAuth button to login"
```

---

## Task 17: Create Files Store

**Files:**
- Create: `frontend/src/stores/files.ts`

**Step 1: Create files store**

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userFileApi, type UserFile, type UserFilesResponse, type Meta } from '@/services/api'

export const useFilesStore = defineStore('files', () => {
  const files = ref<UserFile[]>([])
  const pagination = ref<Meta>({
    page: 1,
    limit: 20,
    total: 0,
    pages: 0
  })
  const uploadProgress = ref<Map<string, number>>(new Map())
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const isEmpty = computed(() => files.value.length === 0)
  const hasMore = computed(() => pagination.value.page < pagination.value.pages)

  // Actions
  async function fetchFiles(params?: {
    page?: number
    limit?: number
    user_uid?: string
    file_type?: string
  }) {
    loading.value = true
    error.value = null

    try {
      const response = await userFileApi.listFiles(params)
      files.value = response.files
      pagination.value = response.meta
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch files'
      console.error('Failed to fetch files:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchFile(fileUid: string) {
    loading.value = true
    error.value = null

    try {
      const file = await userFileApi.getFile(fileUid)
      return file
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch file'
      console.error('Failed to fetch file:', err)
      return null
    } finally {
      loading.value = false
    }
  }

  async function uploadFile(data: {
    user_uid: string
    file_type: string
    file_name: string
    file_path: string
    mime_type: string
    size: number
    visibility?: number
  }) {
    loading.value = true
    error.value = null

    try {
      const response = await userFileApi.uploadFile(data)
      // Add to files list
      files.value.unshift({
        uid: response.uid,
        ...data,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      })
      return response.uid
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to upload file'
      console.error('Failed to upload file:', err)
      return null
    } finally {
      loading.value = false
    }
  }

  async function updateFile(fileUid: string, data: {
    file_name?: string
    visibility?: number
  }) {
    loading.value = true
    error.value = null

    try {
      await userFileApi.updateFile(fileUid, data)
      // Update in list
      const index = files.value.findIndex(f => f.uid === fileUid)
      if (index !== -1) {
        if (data.file_name) files.value[index].file_name = data.file_name
        if (data.visibility !== undefined) files.value[index].visibility = data.visibility
      }
      return true
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to update file'
      console.error('Failed to update file:', err)
      return false
    } finally {
      loading.value = false
    }
  }

  async function deleteFile(fileUid: string) {
    loading.value = true
    error.value = null

    try {
      await userFileApi.deleteFile(fileUid)
      // Remove from list
      files.value = files.value.filter(f => f.uid !== fileUid)
      return true
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to delete file'
      console.error('Failed to delete file:', err)
      return false
    } finally {
      loading.value = false
    }
  }

  function clearError() {
    error.value = null
  }

  return {
    files,
    pagination,
    uploadProgress,
    loading,
    error,
    isEmpty,
    hasMore,
    fetchFiles,
    fetchFile,
    uploadFile,
    updateFile,
    deleteFile,
    clearError
  }
})
```

**Step 2: Commit**

```bash
git add frontend/src/stores/files.ts
git commit -m "feat(frontend): add files store"
```

---

## Task 18: Create FilesView

**Files:**
- Create: `frontend/src/views/FilesView.vue`

**Step 1: Create FilesView component**

```vue
<template>
  <div class="files-view">
    <div class="files-header">
      <h1>My Files</h1>
      <button @click="showUploadModal = true" class="upload-btn">
        Upload File
      </button>
    </div>

    <div v-if="filesStore.loading && filesStore.files.length === 0" class="loading-state">
      <LoadingSpinner />
    </div>

    <div v-else-if="filesStore.error" class="error-state">
      <p>{{ filesStore.error }}</p>
      <button @click="filesStore.clearError(); loadFiles()" class="retry-btn">Retry</button>
    </div>

    <div v-else-if="filesStore.isEmpty" class="empty-state">
      <svg class="empty-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z" />
      </svg>
      <p>No files yet</p>
      <button @click="showUploadModal = true" class="upload-link">Upload your first file</button>
    </div>

    <div v-else class="files-content">
      <div class="files-controls">
        <select v-model="viewMode" class="view-toggle">
          <option value="grid">Grid</option>
          <option value="list">List</option>
        </select>
        <select v-model="filterType" class="filter-select" @change="loadFiles()">
          <option value="">All Types</option>
          <option value="image">Images</option>
          <option value="document">Documents</option>
          <option value="video">Videos</option>
        </select>
      </div>

      <div v-if="viewMode === 'grid'" class="files-grid">
        <FileCard
          v-for="file in filesStore.files"
          :key="file.uid"
          :file="file"
          @delete="confirmDelete(file)"
          @update="openUpdateModal(file)"
        />
      </div>

      <div v-else class="files-list">
        <div v-for="file in filesStore.files" :key="file.uid" class="file-row">
          <FileCard
            :file="file"
            :list-view="true"
            @delete="confirmDelete(file)"
            @update="openUpdateModal(file)"
          />
        </div>
      </div>

      <div v-if="filesStore.hasMore" class="load-more">
        <button @click="loadMore()" :disabled="filesStore.loading" class="load-more-btn">
          {{ filesStore.loading ? 'Loading...' : 'Load More' }}
        </button>
      </div>
    </div>

    <!-- Upload Modal -->
    <FileUpload
      v-if="showUploadModal"
      @close="showUploadModal = false"
      @uploaded="handleFileUploaded"
    />

    <!-- Update Modal -->
    <div v-if="showUpdateModal" class="modal-overlay" @click.self="showUpdateModal = false">
      <div class="modal-content">
        <h2>Update File</h2>
        <form @submit.prevent="handleUpdateFile">
          <div class="form-group">
            <label for="fileName">File Name</label>
            <input
              id="fileName"
              v-model="updateForm.file_name"
              type="text"
              class="form-input"
              required
            />
          </div>
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="updateForm.public" />
              Make Public
            </label>
          </div>
          <div class="modal-actions">
            <button type="button" @click="showUpdateModal = false" class="cancel-btn">Cancel</button>
            <button type="submit" class="submit-btn" :disabled="filesStore.loading">
              {{ filesStore.loading ? 'Updating...' : 'Update' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Confirm Delete Dialog -->
    <ConfirmDialog
      v-if="showDeleteConfirm"
      title="Delete File"
      :message="`Are you sure you want to delete '${fileToDelete?.file_name}'?`"
      @confirm="handleDeleteFile"
      @cancel="showDeleteConfirm = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useFilesStore } from '@/stores/files'
import { useAuthStore } from '@/stores/auth'
import FileCard from '@/components/File/FileCard.vue'
import FileUpload from '@/components/File/FileUpload.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import type { UserFile } from '@/services/api'

const filesStore = useFilesStore()
const authStore = useAuthStore()

const viewMode = ref<'grid' | 'list'>('grid')
const filterType = ref('')
const showUploadModal = ref(false)
const showUpdateModal = ref(false)
const showDeleteConfirm = ref(false)
const fileToDelete = ref<UserFile | null>(null)

const updateForm = ref({
  file_name: '',
  public: false
})

onMounted(() => {
  loadFiles()
})

async function loadFiles() {
  await filesStore.fetchFiles({
    user_uid: authStore.user?.uid,
    file_type: filterType.value || undefined
  })
}

async function loadMore() {
  const nextPage = filesStore.pagination.page + 1
  await filesStore.fetchFiles({
    page: nextPage,
    user_uid: authStore.user?.uid
  })
}

function handleFileUploaded() {
  showUploadModal.value = false
  loadFiles()
}

function openUpdateModal(file: UserFile) {
  updateForm.value = {
    file_name: file.file_name,
    public: file.visibility === 1
  }
  showUpdateModal.value = true
}

async function handleUpdateFile() {
  if (!fileToDelete.value) return

  const success = await filesStore.updateFile(fileToDelete.value.uid, {
    file_name: updateForm.value.file_name,
    visibility: updateForm.value.public ? 1 : 0
  })

  if (success) {
    showUpdateModal.value = false
  }
}

function confirmDelete(file: UserFile) {
  fileToDelete.value = file
  showDeleteConfirm.value = true
}

async function handleDeleteFile() {
  if (!fileToDelete.value) return

  const success = await filesStore.deleteFile(fileToDelete.value.uid)
  if (success) {
    showDeleteConfirm.value = false
    fileToDelete.value = null
  }
}
</script>

<style scoped>
.files-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.files-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.files-header h1 {
  margin: 0;
  font-size: 2rem;
}

.upload-btn {
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s;
}

.upload-btn:hover {
  transform: translateY(-2px);
}

.files-controls {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.view-toggle,
.filter-select {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
}

.files-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
}

.files-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.load-more {
  text-align: center;
  margin-top: 2rem;
}

.load-more-btn {
  padding: 0.75rem 2rem;
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  cursor: pointer;
}

.load-more-btn:hover:not(:disabled) {
  background: #f8f9fa;
}

.loading-state,
.error-state,
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 1rem;
  color: #cbd5e0;
}

.upload-link {
  color: #667eea;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1rem;
}

.retry-btn {
  padding: 0.5rem 1rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  width: 100%;
  max-width: 500px;
}

.modal-content h2 {
  margin: 0 0 1.5rem 0;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.form-input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-sizing: border-box;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

.cancel-btn,
.submit-btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

.cancel-btn {
  background: #e2e8f0;
}

.submit-btn {
  background: #667eea;
  color: white;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
```

**Step 2: Commit**

```bash
git add frontend/src/views/FilesView.vue
git commit -m "feat(frontend): add FilesView"
```

---

## Task 19: Create FileCard Component

**Files:**
- Create: `frontend/src/components/File/FileCard.vue`

**Step 1: Create FileCard component**

```vue
<template>
  <div :class="['file-card', { 'list-view': listView }]">
    <div class="file-icon" :class="fileTypeClass">
      <svg v-if="isImage" viewBox="0 0 24 24" fill="currentColor">
        <path d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.5 13.5l2.5 3.01L14.5 12l4.5 6H5l3.5-4.5z"/>
      </svg>
      <svg v-else viewBox="0 0 24 24" fill="currentColor">
        <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z"/>
      </svg>
    </div>

    <div class="file-info">
      <h3 class="file-name">{{ file.file_name }}</h3>
      <p class="file-meta">{{ formatFileSize(file.size) }} · {{ formatDate(file.created_at) }}</p>
      <span v-if="file.visibility === 1" class="visibility-badge public">Public</span>
      <span v-else class="visibility-badge private">Private</span>
    </div>

    <div class="file-actions">
      <button @click="$emit('update', file)" class="action-btn" title="Rename">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
        </svg>
      </button>
      <button @click="handleDownload" class="action-btn" title="Download">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
          <polyline points="7,10 12,15 17,10"/>
          <line x1="12" y1="15" x2="12" y2="3"/>
        </svg>
      </button>
      <button @click="$emit('delete', file)" class="action-btn danger" title="Delete">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="3,6 5,6 21,6"/>
          <path d="M19,6v14a2,2,0,0,1-2,2H7a2,2,0,0,1-2-2V6m3,0V4a2,2,0,0,1,2-2h4a2,2,0,0,1,2,2v2"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { UserFile } from '@/services/api'

const props = defineProps<{
  file: UserFile
  listView?: boolean
}>()

const emit = defineEmits<{
  delete: [file: UserFile]
  update: [file: UserFile]
}>()

const isImage = computed(() => props.file.mime_type.startsWith('image/'))

const fileTypeClass = computed(() => {
  if (props.file.mime_type.startsWith('image/')) return 'image'
  if (props.file.mime_type.startsWith('video/')) return 'video'
  if (props.file.mime_type.includes('pdf')) return 'pdf'
  return 'document'
})

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString()
}

function handleDownload() {
  // TODO: Implement file download
  console.log('Download file:', props.file.uid)
}
</script>

<style scoped>
.file-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transition: transform 0.2s, box-shadow 0.2s;
  position: relative;
}

.file-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0,0,0,0.15);
}

.file-card.list-view {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
}

.file-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1rem;
  color: white;
}

.file-card.list-view .file-icon {
  margin-bottom: 0;
  flex-shrink: 0;
}

.file-icon.image { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.file-icon.video { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.file-icon.pdf { background: linear-gradient(135deg, #fa709a 0%, #fee140 100%); }
.file-icon.document { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }

.file-icon svg {
  width: 24px;
  height: 24px;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  margin: 0 0 0.5rem 0;
  font-size: 1rem;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-meta {
  margin: 0;
  font-size: 0.875rem;
  color: #7f8c8d;
}

.visibility-badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 500;
  margin-top: 0.5rem;
}

.visibility-badge.public {
  background: #d1fae5;
  color: #065f46;
}

.visibility-badge.private {
  background: #e2e8f0;
  color: #475569;
}

.file-actions {
  display: flex;
  gap: 0.5rem;
  opacity: 0;
  transition: opacity 0.2s;
}

.file-card:hover .file-actions {
  opacity: 1;
}

.file-card.list-view .file-actions {
  opacity: 1;
}

.action-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #f3f4f6;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #4b5563;
  transition: background 0.2s;
}

.action-btn:hover {
  background: #e5e7eb;
}

.action-btn.danger:hover {
  background: #fee2e2;
  color: #dc2626;
}

.action-btn svg {
  width: 18px;
  height: 18px;
}
</style>
```

**Step 2: Commit**

```bash
git add frontend/src/components/File/FileCard.vue
git commit -m "feat(frontend): add FileCard component"
```

---

## Task 20: Create FileUpload Component

**Files:**
- Create: `frontend/src/components/File/FileUpload.vue`

**Step 1: Create FileUpload component**

```vue
<template>
  <div class="file-upload-modal" @click.self="$emit('close')">
    <div class="upload-container" @click.self="$emit('close')">
      <div class="upload-content" @click.self="$emit('close')">
        <div class="upload-header">
          <h2>Upload File</h2>
          <button @click="$emit('close')" class="close-btn">&times;</button>
        </div>

        <div
          class="upload-zone"
          :class="{ 'drag-over': isDragOver }"
          @dragover.prevent="isDragOver = true"
          @dragleave.prevent="isDragOver = false"
          @drop.prevent="handleDrop"
          @click="$refs.fileInput.click()"
        >
          <input
            ref="fileInput"
            type="file"
            @change="handleFileSelect"
            class="file-input"
          />
          <svg class="upload-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
          </svg>
          <p class="upload-text">
            {{ isDragOver ? 'Drop your file here' : 'Drag & drop a file or click to browse' }}
          </p>
          <p class="upload-hint">Maximum file size: 10MB</p>
        </div>

        <div v-if="selectedFile" class="file-preview">
          <div class="file-info">
            <span class="file-name">{{ selectedFile.name }}</span>
            <span class="file-size">{{ formatFileSize(selectedFile.size) }}</span>
          </div>
          <button @click="selectedFile = null" class="remove-btn">&times;</button>
        </div>

        <div v-if="uploadProgress > 0" class="upload-progress">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: uploadProgress + '%' }"></div>
          </div>
          <p class="progress-text">{{ uploadProgress }}% uploaded</p>
        </div>

        <div class="upload-options">
          <label class="visibility-toggle">
            <input type="checkbox" v-model="makePublic" />
            <span>Make file public</span>
          </label>
        </div>

        <div class="upload-actions">
          <button @click="$emit('close')" class="cancel-btn">Cancel</button>
          <button
            @click="handleUpload"
            :disabled="!selectedFile || uploading"
            class="upload-btn"
          >
            {{ uploading ? 'Uploading...' : 'Upload' }}
          </button>
        </div>

        <p v-if="error" class="error-message">{{ error }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useFilesStore } from '@/stores/files'

const emit = defineEmits<{
  close: []
  uploaded: []
}>()

const authStore = useAuthStore()
const filesStore = useFilesStore()

const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)
const isDragOver = ref(false)
const makePublic = ref(false)
const uploading = ref(false)
const uploadProgress = ref(0)
const error = ref<string | null>(null)

function handleDrop(e: DragEvent) {
  isDragOver.value = false
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    validateAndSelectFile(files[0])
  }
}

function handleFileSelect(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    validateAndSelectFile(target.files[0])
  }
}

function validateAndSelectFile(file: File) {
  error.value = null

  // Check file size (10MB max)
  if (file.size > 10 * 1024 * 1024) {
    error.value = 'File size exceeds 10MB limit'
    return
  }

  selectedFile.value = file
}

async function handleUpload() {
  if (!selectedFile.value || !authStore.user?.uid) return

  uploading.value = true
  error.value = null
  uploadProgress.value = 0

  try {
    // Simulate upload progress
    const progressInterval = setInterval(() => {
      if (uploadProgress.value < 90) {
        uploadProgress.value += 10
      }
    }, 200)

    // Read file as base64
    const reader = new FileReader()
    reader.onload = async (e) => {
      const base64 = (e.target?.result as string).split(',')[1]

      const uid = await filesStore.uploadFile({
        user_uid: authStore.user!.uid,
        file_type: getFileType(selectedFile.value!.type),
        file_name: selectedFile.value!.name,
        file_path: '',  // Will be set by server
        mime_type: selectedFile.value!.type,
        size: selectedFile.value!.size,
        visibility: makePublic.value ? 1 : 0
      })

      clearInterval(progressInterval)
      uploadProgress.value = 100

      if (uid) {
        setTimeout(() => {
          emit('uploaded')
          resetForm()
        }, 500)
      }
    }
    reader.readAsDataURL(selectedFile.value)
  } catch (err: any) {
    error.value = err.message || 'Upload failed'
    uploading.value = false
  }
}

function getFileType(mimeType: string): string {
  if (mimeType.startsWith('image/')) return 'image'
  if (mimeType.startsWith('video/')) return 'video'
  if (mimeType.includes('pdf')) return 'pdf'
  return 'document'
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

function resetForm() {
  selectedFile.value = null
  uploading.value = false
  uploadProgress.value = 0
  error.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}
</script>

<style scoped>
.file-upload-modal {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.upload-container {
  width: 100%;
  max-width: 500px;
  padding: 2rem;
}

.upload-content {
  background: white;
  border-radius: 12px;
  padding: 2rem;
}

.upload-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.upload-header h2 {
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #7f8c8d;
}

.upload-zone {
  border: 2px dashed #cbd5e0;
  border-radius: 12px;
  padding: 3rem 2rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
}

.upload-zone:hover,
.upload-zone.drag-over {
  border-color: #667eea;
  background: #f8f9ff;
}

.file-input {
  display: none;
}

.upload-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 1rem;
  color: #cbd5e0;
}

.upload-text {
  margin: 0 0 0.5rem 0;
  font-weight: 500;
  color: #2c3e50;
}

.upload-hint {
  margin: 0;
  font-size: 0.875rem;
  color: #7f8c8d;
}

.file-preview {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 8px;
  margin: 1rem 0;
}

.file-info {
  display: flex;
  flex-direction: column;
}

.file-name {
  font-weight: 500;
  color: #2c3e50;
}

.file-size {
  font-size: 0.875rem;
  color: #7f8c8d;
}

.remove-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #dc2626;
}

.upload-progress {
  margin: 1rem 0;
}

.progress-bar {
  height: 8px;
  background: #e2e8f0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  transition: width 0.3s;
}

.progress-text {
  margin: 0.5rem 0 0 0;
  font-size: 0.875rem;
  color: #7f8c8d;
  text-align: center;
}

.upload-options {
  margin: 1rem 0;
}

.visibility-toggle {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
}

.upload-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.cancel-btn,
.upload-btn {
  flex: 1;
  padding: 0.875rem;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}

.cancel-btn {
  background: #e2e8f0;
  color: #475569;
}

.upload-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.upload-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background: #fee;
  color: #c33;
  border-radius: 8px;
  text-align: center;
}
</style>
```

**Step 2: Commit**

```bash
git add frontend/src/components/File/FileUpload.vue
git commit -m "feat(frontend): add FileUpload component"
```

---

## Task 21: Add Files Route and Sidebar Link

**Files:**
- Modify: `frontend/src/router/index.ts`
- Modify: `frontend/src/components/Dashboard/Sidebar.vue`

**Step 1: Add Files route to router**

After line 67 (proto-test route), add:

```typescript
{
  path: 'files',
  name: 'files',
  component: () => import('@/views/FilesView.vue')
}
```

**Step 2: Add Files link to Sidebar**

Find the Sidebar component and add a Files navigation item. The pattern should be similar to existing links.

**Step 3: Commit**

```bash
git add frontend/src/router/index.ts frontend/src/components/Dashboard/Sidebar.vue
git commit -m "feat(frontend): add Files route and sidebar link"
```

---

## Task 22: Implement UserDetailView

**Files:**
- Modify: `frontend/src/views/UserDetailView.vue`

**Step 1: Implement UserDetailView**

The file exists but is empty. Implement with tabs for Profile, Files, and Devices.

**Step 2: Commit**

```bash
git add frontend/src/views/UserDetailView.vue
git commit -m "feat(frontend): implement UserDetailView with tabs"
```

---

## Task 23: Test and Finalize

**Files:**
- Update: `docker-compose.yml` (if needed)

**Step 1: Verify service-user OAuth config**

Ensure service-user has Google OAuth credentials configured.

**Step 2: Test OAuth flow**

1. Start services: `docker-compose up`
2. Navigate to login
3. Click "Continue with Google"
4. Verify redirect and callback

**Step 3: Test file upload**

1. Navigate to Files page
2. Upload a file
3. Verify file appears in list

**Step 4: Commit any final fixes**

```bash
git add .
git commit -m "fix: final adjustments to OAuth and Files integration"
```

---

## Task 24: Create ADR and Update Documentation

**Files:**
- Create: `docs/adr/004-oauth-files-integration.md`

**Step 1: Create ADR**

Document the architecture decision for OAuth and File Management integration.

**Step 2: Update README**

Add OAuth and File Management features to the README.

**Step 3: Commit**

```bash
git add docs/
git commit -m "docs: add ADR for OAuth and Files integration"
```

---

## End of Implementation Plan

**Summary of changes:**
- OAuth: 6 backend tasks + 1 frontend task
- Files: 17 backend tasks + 6 frontend tasks

**Total:** ~30 tasks

**Testing Strategy:**
- OAuth flow with Google OAuth playground
- File upload/download with various file types
- E2E tests in separate session

**Next Steps After Implementation:**
1. Create pull request from feature branch
2. Code review
3. Merge to main branch
4. Deploy to staging for final testing
