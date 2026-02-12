# OAuth and File Management Integration

## Status
Proposed: 2025-02-12

## Context
The microservice-sample-app needs to integrate Google OAuth authentication and User File Management functionality from the service-user microservice. This integration will enable users to:
1. Sign in using their Google account
2. Upload, manage, and share files with visibility controls

## Decision
We will integrate Google OAuth and File Management features into the API Gateway and Frontend, following the existing BFF (Backend For Frontend) pattern and hexagonal architecture.

### Key Design Decisions:

1. **OAuth Flow**: Service-user owns all OAuth configuration (Google Client ID, secrets). API Gateway acts as a proxy, redirecting to service-user for authorization URL generation and token exchange.

2. **File Upload**: Files are uploaded via multipart/form-data through API Gateway to service-user gRPC endpoint. Maximum file size is 10MB.

3. **File Storage**: service-user handles all file storage (path generation, S3/cloud configuration). Gateway is stateless regarding file data.

4. **Visibility Model**: Files have visibility attribute (0=private, 1=public) to control access control.

5. **Frontend Integration**: Vue 3 frontend with Pinia stores for state management, Axios interceptors for JWT refresh, and component-based architecture.

## Architecture Diagram

```
┌─────────────┐           ┌──────────────┐
│   Frontend  │           │   gRPC        │
│  (Vue 3)    │◄──────►│ service-user  │
│             │           │              │  (OAuth +     │
└─────────────┘           │              │  Files API)  │
                     └──────────────┘           └──────────────┘
                               │
                         ┌──────────────┐
                         │   API Gateway  │
                         │ (BFF - Go)    │
                         └───────────────┘
```

### Component Structure (Backend - API Gateway):

```
api-gateway/
├── cmd/main.go                          # Application entry point
├── internal/
│   ├── application/                    # Application services (use cases)
│   │   ├── auth/                   # Auth application service
│   │   └── user_file/             # UserFile application service
│   ├── core/
│   │   ├── domain/                  # Domain models and params
│   │   │   ├── model/              # User, Tokens, UserFile, etc.
│   │   │   └── params/             # Input parameters
│   │   └── port/                   # Primary/Secondary ports
│   │       ├── service/              # Service interfaces (primary)
│   │       └── client/              # Client interfaces (secondary)
│   └── adapter/                    # Adapters
│       ├── primary/                 # Driving adapters (HTTP handlers)
│       │   └── http/
│       │       ├── httpHandler/   # HTTP request handlers
│       │       ├── response/       # HTTP response DTOs
│       │       └── request/        # HTTP request DTOs
│       └── secondary/               # Driven adapters
│           └── grpc/           # gRPC clients
│               └── mapper/         # Proto-to-domain mappers
└── pkg/                              # Shared utilities
    ├── http/                       # HTTP middleware (CORS, recovery)
    └── config/                     # Configuration management
```

### Component Structure (Frontend):

```
frontend/
├── src/
│   ├── components/
│   │   ├── common/            # LoadingSpinner, ConfirmDialog, NotificationToast
│   │   ├── File/             # FileCard, FileUpload
│   │   └── Dashboard/        # Sidebar, Header
│   ├── stores/
│   │   ├── auth.ts           # Authentication state (Pinia)
│   │   └── files.ts          # File management state (Pinia)
│   ├── views/
│   │   ├── LoginView.vue     # Login with Google OAuth button
│   │   ├── FilesView.vue     # File listing and management
│   │   └── UserDetailView.vue # User details with tabs (Profile, Files, Devices)
│   ├── services/
│   │   └── api.ts            # API client with Axios interceptors
│   └── router/
│       └── index.ts        # Vue Router configuration
└── package.json               # Dependencies and scripts
```

### API Endpoints:

**OAuth:**
- `GET /api/v1/auth/google` - Initiate OAuth flow, redirects to Google
- `GET /api/v1/auth/google/callback` - Handle OAuth callback, returns tokens via redirect

**Files:**
- `GET /api/v1/files?page={n}&limit={n}` - List user files with pagination
- `GET /api/v1/files/{uid}` - Get single file details
- `POST /api/v1/files` - Upload file (multipart/form-data, max 10MB)
- `PUT/PATCH /api/v1/files/{uid}` - Update file (name, visibility)
- `DELETE /api/v1/files/{uid}` - Delete file

## Alternatives Considered

1. **Direct OAuth Integration**: Frontend could call Google OAuth directly. Rejected because service-user already owns OAuth config and secrets management is centralized.

2. **Frontend File Upload to Cloud Storage**: Rejected because we need service-level metadata (thumbnails, visibility) and access control managed by service-user.

3. **Separate File Service**: Rejected to avoid adding another microservice. File management is core to service-user domain.

4. **WebSocket for File Upload Progress**: Considered but not implemented to keep HTTP simplicity. Multipart upload with progress feedback via polling is sufficient.

## Consequences

### Positive:
- Single entry point for authentication (service-user via gateway)
- Consistent file access patterns across the application
- Reusable file management components for future features
- Gateway-centric architecture keeps frontend lightweight and focused on UI

### Negative:
- Additional dependency on service-user availability for file operations
- Service-user becomes single point of failure for file-related features
- OAuth tokens issued by service-user have same lifecycle as regular auth tokens

### Mitigation:
- Implement circuit breaker pattern for gRPC calls to service-user
- Add health check endpoint for service-user availability monitoring
- Implement retry logic with exponential backoff for failed requests
- Graceful degradation when file service is unavailable (show cached/stub data)

## Implementation Checklist

- [x] Task 1: Extend AuthClient interface for OAuth
- [x] Task 2: Extend AuthService interface for OAuth
- [x] Task 3: Implement OAuth in gRPC client adapter
- [x] Task 4: Implement OAuth in auth application service
- [x] Task 5: Add OAuth HTTP handlers
- [x] Task 6: Update gateway config for OAuth
- [x] Task 7: Create UserFileClient interface
- [x] Task 8: Create UserFileService interface
- [x] Task 9: Create UserFile domain models
- [x] Task 10: Create UserFile gRPC client adapter
- [x] Task 11: Create UserFile application service
- [x] Task 12: Create File HTTP response DTOs
- [x] Task 13: Create File HTTP request DTOs
- [x] Task 14: Create File HTTP handler
- [x] Task 15: Wire UserFile handlers in gateway
- [x] Task 16: Add Google OAuth button to LoginView
- [x] Task 17: Create Files store
- [x] Task 18: Create FilesView
- [x] Task 19: Create FileCard component
- [x] Task 20: Create FileUpload component
- [x] Task 21: Add Files route and sidebar link
- [x] Task 22: Implement UserDetailView with tabs
- [x] Task 23: Verify service-user OAuth config and test flows
- [x] Task 24: Create ADR and update documentation

## References

- [RFC 6749 - OAuth 2.0](https://datatracker.ietf.org/doc/html/rfc6749)
- [service-user proto definitions](../../../service-user/service-user-proto/proto/)
- [Hexagonal Architecture pattern](https://alistair.cwi.cc/20/hexagonal-architecture/)
- [BFF pattern by Martin Fowler](https://martinfowler.com/articles/blikiPatterns/BusinessDelegatePattern)
- [Vue 3 Composition API](https://vuejs.org/guide/introduction.html#composition-api)
- [Pinia state management](https://pinia.vuejs.org/)
