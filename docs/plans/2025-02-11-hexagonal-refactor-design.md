# Hexagonal Architecture Refactor Design

**Date:** 2025-02-11
**Component:** API Gateway
**Status:** Design Complete
**Reference:** ADR 004

## Overview

Refactor the API Gateway to follow **hexagonal architecture** (ports and adapters), aligning with the service-user microservice structure. This refactoring improves testability, maintainability, and consistency across the codebase.

## Current State

The API Gateway follows a basic layered architecture:

```
cmd/main.go
    ↓
handler/ (HTTP handlers)
    ↓
service/ (business logic)
    ↓
client/ (gRPC clients)
```

**Issues:**
- Services depend directly on concrete gRPC types (`userpb.UserServiceClient`)
- No explicit port interfaces - tight coupling between layers
- Domain concepts mixed with application logic
- Difficult to test without real gRPC connections

## Target State

Hexagonal architecture aligned with service-user:

```
cmd/main.go (wiring)
    ↓
adapter/primary/http/ (Primary Adapters - Driving)
    ↓
application/ (Use Cases)
    ↓
core/port/client/ (Ports - Interfaces)
    ↓
adapter/secondary/grpc/ (Secondary Adapters - Driven)
```

## Architecture Structure

```
api-gateway/
├── cmd/
│   └── main.go                        # Dependency wiring
│
├── internal/
│   ├── core/                          # Pure core (no external dependencies)
│   │   ├── domain/
│   │   │   ├── model/                 # Domain entities with business logic
│   │   │   │   ├── user.go
│   │   │   │   ├── profile.go
│   │   │   │   └── device.go
│   │   │   ├── params/                # Use case parameters
│   │   │   │   ├── user.go
│   │   │   │   ├── auth.go
│   │   │   │   └── common.go
│   │   │   └── errors/                # Domain errors
│   │   │       └── errors.go
│   │   │
│   │   └── port/
│   │       ├── service/               # Service ports (use case interfaces)
│   │       │   ├── user.go
│   │       │   └── auth.go
│   │       └── client/                # Client ports (gRPC interfaces)
│   │           ├── user_client.go
│   │           └── auth_client.go
│   │
│   ├── application/                   # Use case implementations
│   │   ├── user/
│   │   │   └── service.go
│   │   └── auth/
│   │       └── service.go
│   │
│   └── config/
│       └── config.go
│
├── adapter/                           # External implementations
│   ├── primary/                       # Driving adapters
│   │   └── http/
│   │       ├── handler/               # HTTP request handlers
│   │       │   ├── user_handler.go
│   │       │   └── auth_handler.go
│   │       ├── request/               # Request DTOs with validation tags
│   │       │   └── user_request.go
│   │       ├── response/              # Response DTOs
│   │       │   └── response_mapper.go
│   │       └── validator/             # go-playground/validator
│   │           └── validator.go
│   │
│   └── secondary/                     # Driven adapters
│       └── grpc/
│           ├── user_client.go         # Implements client.UserClient
│           ├── auth_client.go         # Implements client.AuthClient
│           └── mapper/                # gRPC ↔ Domain converters
│               ├── user_mapper.go
│               └── auth_mapper.go
│
├── pkg/
│   └── http/
│       └── middleware.go              # Shared middleware
│
└── Dockerfile
```

## Layer Responsibilities

### Domain Layer (`core/domain/`)

**Purpose:** Pure business entities with no external dependencies.

- **model/**: Domain entities with business logic methods
  - Example: `User.CanLogin()`, `User.Activate()`, `User.Ban()`
- **params/**: Plain structs for use case inputs
- **errors/**: Domain-specific errors

**Rules:**
- No imports from external packages (gRPC, HTTP)
- Business logic lives in domain methods
- Pure Go structs and interfaces

### Port Layer (`core/port/`)

**Purpose:** Interface definitions (contracts) without implementation.

- **service/**: Primary ports - what the application provides
  - Example: `UserService` interface with use case methods
- **client/**: Secondary ports - what the application needs
  - Example: `UserClient` interface for gRPC operations

**Rules:**
- Ports are pure interfaces
- All types referenced are domain types
- No imports of gRPC, HTTP, or other external packages

### Application Layer (`application/`)

**Purpose:** Use case implementations that coordinate between ports.

- Implements service port interfaces
- Contains business rules that don't belong in domain entities
- Orchestrates calls to secondary adapters
- No external input validation (handled in adapter layer)

**Rules:**
- Depends only on ports (interfaces), not concrete implementations
- Business rules only - no HTTP, no gRPC
- Domain errors mapped to application errors

### Adapter Layer (`adapter/`)

**Purpose:** Implementations that connect the application to external concerns.

#### Primary Adapters (`adapter/primary/http/`)

**Driving adapters** that control the application:

- **handler/**: HTTP request handlers
- **request/**: Request DTOs with validation tags
- **validator/**: go-playground/validator setup

**Responsibilities:**
- Parse HTTP requests
- Validate input using go-playground/validator
- Convert request DTOs to domain params
- Call application services
- Convert responses to HTTP format

#### Secondary Adapters (`adapter/secondary/grpc/`)

**Driven adapters** that implement client ports:

- Implement `client.UserClient`, `client.AuthClient`
- Convert domain params to gRPC requests
- Convert gRPC responses to domain models
- Handle gRPC-specific concerns

**Responsibilities:**
- Implement client port interfaces
- Handle gRPC communication
- Map between gRPC types and domain types
- Handle connection lifecycle

## Data Flow

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Request                            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              Primary Adapter (handler)                       │
│  - Parse HTTP request                                        │
│  - Validate using go-playground/validator                    │
│  - Convert to domain params                                  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│           Application Service (use case)                     │
│  - Apply business rules                                      │
│  - Orchestrate operations                                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│          Secondary Adapter (gRPC client)                     │
│  - Convert domain to gRPC types                              │
│  - Call external service                                     │
│  - Convert gRPC to domain types                              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   External gRPC Service                      │
└─────────────────────────────────────────────────────────────┘
```

## Key Patterns

### 1. Port Interface Example

```go
// core/port/client/user_client.go

package client

import (
    "context"
    "github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
    "github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
)

// UserClient defines the interface for calling user gRPC service.
type UserClient interface {
    Get(ctx context.Context, uid string) (*model.User, error)
    List(ctx context.Context, param *params.ListUsersParam) (*model.Users, error)
    Create(ctx context.Context, param *params.CreateUserParam) (string, error)
    Update(ctx context.Context, param *params.UpdateUserParam) error
    Delete(ctx context.Context, uid string) error
    GetProfile(ctx context.Context, userUID string) (*model.Profile, error)
    UpdateProfile(ctx context.Context, param *params.UpdateUserProfileParam) error
    ListDevice(ctx context.Context, param *params.ListUserDevicesParam) (*model.Devices, error)
    RevokeDevice(ctx context.Context, userUID, deviceUID string) error
}
```

### 2. Domain Entity Example

```go
// core/domain/model/user.go

package model

import "time"

type User struct {
    UID       string
    Username  string
    Email     string
    Status    int32
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}

// IsActive returns true if the user is active.
func (u *User) IsActive() bool {
    return u.Status == 1
}

// CanDisplay returns true if user can be shown in listings.
func (u *User) CanDisplay() bool {
    return !u.DeletedAt.IsZero()
}
```

### 3. Application Service Example

```go
// application/user/service.go

package user

import (
    "context"
    "fmt"
    "github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/errors"
    "github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
    "github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
    "github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/port/client"
    "github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/port/service"
)

type UserApplicationService struct {
    userClient client.UserClient
}

func NewUserApplicationService(userClient client.UserClient) service.UserService {
    return &UserApplicationService{userClient: userClient}
}

func (s *UserApplicationService) GetUser(ctx context.Context, uid string) (*model.User, error) {
    if uid == "" {
        return nil, errors.ErrInvalidUID
    }

    user, err := s.userClient.Get(ctx, uid)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    if !user.CanDisplay() {
        return nil, errors.ErrUserNotFound
    }

    return user, nil
}
```

### 4. Validation in Adapter

```go
// adapter/primary/http/request/user_request.go

package request

// CreateUserRequest represents the HTTP request for creating a user.
type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

### 5. Mapper Example

```go
// adapter/secondary/grpc/mapper/user_mapper.go

package mapper

import (
    "time"
    "google.golang.org/protobuf/types/known/timestamppb"
    userpb "github.com/adityakw90/service-user-proto/gen/go/user"
    "github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
)

// UserFromProto converts a gRPC User protobuf to a domain User.
func UserFromProto(u *userpb.User) *model.User {
    return &model.User{
        UID:       u.Uid,
        Username:  u.Username,
        Email:     u.Email,
        Status:    u.Status,
        CreatedAt: TimestampToTime(u.CreatedAt),
        UpdatedAt: TimestampToTime(u.UpdatedAt),
        DeletedAt: TimestampToTimePtr(u.DeletedAt),
    }
}
```

## Testing Strategy

### Test Categories

| Type | Location | What It Tests | Dependencies |
|------|----------|---------------|--------------|
| **Unit Tests** | `application/*/service_test.go` | Use case logic | Mock clients |
| **Domain Tests** | `core/domain/model/*_test.go` | Business rules | None |
| **Handler Tests** | `adapter/primary/http/handler/*_test.go` | HTTP handling | Mock services |
| **Adapter Tests** | `adapter/secondary/grpc/*_test.go` | gRPC integration | Real service |

### Mock Generation

Using `mockgen` for generating port interface mocks:

```bash
# Generate mocks for client ports
mockgen -source=internal/core/port/client/user_client.go \
    -destination=internal/core/port/client/user_client_mock.go

# Generate mocks for service ports
mockgen -source=internal/core/port/service/user.go \
    -destination=internal/core/port/service/user_mock.go
```

### Unit Test Example

```go
func TestUserApplicationService_GetUser_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mockclient.NewMockUserClient(ctrl)
    service := user.NewUserApplicationService(mockClient)

    expectedUser := &model.User{
        UID:      "user-123",
        Username: "testuser",
        Status:   1,
    }

    mockClient.EXPECT().
        Get(gomock.Any(), "user-123").
        Return(expectedUser, nil)

    result, err := service.GetUser(context.Background(), "user-123")

    require.NoError(t, err)
    assert.Equal(t, "user-123", result.UID)
}
```

## Dependency Injection

All dependencies wired in `cmd/main.go`:

```go
func main() {
    cfg := config.Load()

    // Secondary Adapters (Driven)
    grpcUserClient := grpcAdapter.NewUserClientAdapter(cfg.UserServiceAddress)
    defer grpcUserClient.Close()

    // Application Services
    userAppService := user.NewUserApplicationService(grpcUserClient)

    // Primary Adapters (Driving)
    v := validator.NewValidator()
    userHandler := handler.NewUserHandler(userAppService, v)

    // HTTP Server
    router := mux.NewRouter()
    api := router.PathPrefix("/api/v1").Subrouter()
    userHandler.RegisterRoutes(api)
    // ...
}
```

## Migration Strategy

### Phase 1: Set up new structure (non-breaking)

Create new directories alongside existing code. Keep existing handlers/services working.

### Phase 2: Implement core and ports

1. Create domain models (`core/domain/model/`)
2. Create parameter objects (`core/domain/params/`)
3. Define port interfaces (`core/port/`)

### Phase 3: Implement adapters

1. Create gRPC adapter implementing `client.UserClient`
2. Create application service implementing `service.UserService`
3. Create HTTP handler using the application service

### Phase 4: Migrate routes incrementally

Update `main.go` to use new handlers for specific routes while keeping old handlers for others.

### Phase 5: Remove old code

After all routes migrated, remove old `handler/`, `service/`, `client/` directories.

### Migration Checklist

| Step | Task | Status |
|------|------|--------|
| 1 | Create `core/domain/model/` | ⬜ |
| 2 | Create `core/domain/params/` | ⬜ |
| 3 | Create `core/domain/errors/` | ⬜ |
| 4 | Create `core/port/service/` interfaces | ⬜ |
| 5 | Create `core/port/client/` interfaces | ⬜ |
| 6 | Create `adapter/secondary/grpc/` | ⬜ |
| 7 | Create `application/user/` | ⬜ |
| 8 | Create `adapter/primary/http/handler/` | ⬜ |
| 9 | Update `main.go` wiring | ⬜ |
| 10 | Add tests | ⬜ |
| 11 | Migrate all routes | ⬜ |
| 12 | Remove old code | ⬜ |

## Dependencies

### Go Packages

```go
// Validation
"github.com/go-playground/validator/v10"

// Testing
"github.com/golang/mock/gomock"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"

// Existing
"github.com/gorilla/mux"
"google.golang.org/grpc"
```

### Tooling

```bash
# Install mockgen
go install github.com/golang/mock/mockgen@latest

# Generate mocks
make mocks
```

## References

- ADR 004: Hexagonal Architecture
- service-user repository: `/media/adit/SSD/project/MTAmedia/repo/service-user/service-user`
- Alistair Cockburn's Hexagonal Architecture: https://alistair.cockburn.us/hexagonal-architecture/
