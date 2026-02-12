# OAuth and File Management Integration Design

**Date:** 2026-02-12
**Status:** Approved
**Author:** Claude

## Overview

Integrate Google OAuth and UserFileService from service-user into microservice-sample-app with complete backend and frontend implementation.

## Scope

1. **Google OAuth** - Full social login via Google
2. **File Management** - Complete UserFileService integration (upload, list, update, delete)
3. **Frontend UI** - OAuth login button, file management interface, UserDetailView

## Architecture Principles

- **Clean microservice design:** service-user owns all OAuth config (client_id, client_secret, scopes, state, PKCE, token exchange)
- **Gateway minimal responsibility:** Only HTTP routing and redirects
- **Frontend simplicity:** Initiate flow, receive callback via redirect

---

## Section 1: Google OAuth Integration

### Responsibilities

| Component | Owns |
|-----------|------|
| **service-user** | client_id, client_secret, scopes, PKCE, state generation/validation, token exchange |
| **Gateway** | HTTP endpoints, redirect routing only |
| **Frontend** | Initiate flow, receive callback via URL params |

### gRPC Methods (Existing in service-user)

```protobuf
rpc GoogleOAuth(GoogleOAuthRequest) returns (GoogleOAuthResponse);
rpc HandleGoogleOAuth(HandleGoogleOAuthRequest) returns (AuthResponse);
```

### API Gateway Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/auth/google` | Initiate OAuth flow |
| GET | `/api/v1/auth/google/callback` | Handle Google callback |

### OAuth Flow Sequence

```
1. Frontend → Gateway: GET /api/v1/auth/google
2. Gateway → service-user (gRPC): GoogleOAuth()
3. service-user: generates state (32-byte), stores in Redis (10min TTL), returns auth_url
4. Gateway → Frontend: 302 Redirect to auth_url (Google)
5. User: approves on Google
6. Google → Gateway: GET /api/v1/auth/google/callback?code=xxx&state=yyy
7. Gateway → service-user (gRPC): HandleGoogleOAuth(code, state)
8. service-user: validates state, exchanges code for tokens, creates/returns JWT
9. service-user → Gateway: AuthResponse (access_token, refresh_token, user)
10. Gateway → Frontend: sets HTTP-only cookies, redirects with access_token
11. Frontend: stores token in Pinia, redirects to dashboard
```

### Gateway Configuration Required

```env
USER_SERVICE_ADDRESS=service-user:50051  # existing
OAUTH_REDIRECT_URI=http://gateway:8080/api/v1/auth/google/callback
FRONTEND_REDIRECT_URI=http://localhost:3000/login
```

### Error Handling

| Scenario | Handling |
|----------|----------|
| State mismatch | Gateway returns 400, redirects to login with error param |
| Token exchange fails | Gateway returns 502, redirects to login with error param |
| User cancels OAuth | Google redirects with error=access_denied, Gateway shows message |
| Email already exists | service-user links to existing account |
| New user | service-user creates account, returns new user flag |

---

## Section 2: File Management Integration

### API Gateway Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/files` | Upload file (multipart/form-data) |
| GET | `/api/v1/files` | List user files with pagination |
| GET | `/api/v1/files/{fileUid}` | Get file metadata/content |
| PUT | `/api/v1/files/{fileUid}` | Update file (name, visibility) |
| DELETE | `/api/v1/files/{fileUid}` | Delete file |

### gRPC Methods (Existing in service-user)

```protobuf
rpc Get(GetUserFileRequest) returns (UserFile);
rpc List(ListUserFilesRequest) returns (ListUserFilesResponse);
rpc Add(AddUserFileRequest) returns (UserFile);
rpc Update(UpdateUserFileRequest) returns (UserFile);
rpc Delete(DeleteUserFileRequest) returns (google.protobuf.Empty);
```

### Key Design Decisions

1. **File Upload:** Gateway receives `multipart/form-data`, streams to service-user via gRPC
2. **Storage:** service-user manages file storage path (configured in `config.yaml`)
3. **Download:** Gateway proxies file through service-user for consistency
4. **File Size Limits:** Gateway validates max size (default 10MB, configurable)
5. **Thumbnails:** service-user generates; Gateway returns URLs in metadata

### Error Handling

| Scenario | Handling |
|----------|----------|
| File too large | Gateway validates before forwarding, returns 413 |
| Invalid file type | Gateway validates allowlist, returns 400 |
| Storage quota exceeded | service-user returns 403 |
| File not found | service-user returns 404 |
| Upload interrupted | Gateway handles timeout, returns 408 |

### Security

- Validate MIME-type on gateway (don't trust extension)
- Sanitize filenames (remove path traversal)
- Rate limiting on upload endpoint

---

## Section 3: Frontend UI Components

### OAuth Frontend Changes

**Modified Components:**

| Component | Changes |
|-----------|---------|
| `LoginView.vue` | Add "Continue with Google" button, handle URL params for callback |
| `api.ts` | Add `googleOAuth()` method |

**OAuth Button Implementation:**
```typescript
const googleOAuth = () => {
  window.location.href = `${API_BASE_URL}/api/v1/auth/google`
}
```

**Callback Handling:**
- Gateway redirects to frontend with `?token=xxx&refresh=yyy`
- LoginView checks URL params on mount
- Store tokens in Pinia auth store and redirect to dashboard

### File Management UI

**New Components:**

| Component | Route | Purpose |
|-----------|-------|---------|
| `FilesView.vue` | `/files` | Main file management page |
| `FileUpload.vue` | - | Upload component (drag & drop) |
| `FileCard.vue` | - | Display file with actions |
| `FilePreview.vue` | - | Modal for file preview |

**FilesView Features:**
- Drag & drop upload area
- Grid/list toggle
- Filter by file type, visibility
- Search files by name
- Pagination
- Context menu: download, rename, delete, change visibility

### UserDetailView Implementation

**New Component:** `UserDetailView.vue` (route exists, not implemented)

**Features:**
- User profile card with avatar
- Account information
- Tabs: Files, Devices, Activity
- Edit profile modal
- Reset password button (admin)

### Pinia Stores

**New/Updated Stores:**

```typescript
// auth.store.ts - Add OAuth state
oauthInProgress: boolean
oauthError: string | null

// files.store.ts - New store
files: UserFile[]
pagination: { page, limit, total }
uploadProgress: Map<fileUid, number>
```

---

## Section 4: Implementation Order

### Phase 1: OAuth Backend
1. Gateway: Add OAuth handlers (`/google`, `/google/callback`)
2. Implement state management and gRPC calls to service-user
3. Test flow with Google or mock

### Phase 2: OAuth Frontend
4. Add Google button to LoginView
5. Implement callback redirect handling

### Phase 3: File Backend
6. Gateway: Add file handlers (upload, list, get, update, delete)
7. Implement multipart/form-data handling
8. Test file operations

### Phase 4: File Frontend
9. Create FilesView, FileUpload, FileCard components
10. Implement files.store.ts
11. Add routing and navigation

### Phase 5: UserDetailView
12. Implement UserDetailView with files/devices tabs

---

## Testing Strategy

### Backend Testing
- **Unit:** Handler OAuth callback parsing, file upload validation
- **Integration:** Gateway → service-user gRPC calls
- **E2E:** Full OAuth flow, file upload/download cycle

### Frontend Testing
- **Component:** FileUpload dropzone, FileCard actions
- **Integration:** File store mutations, OAuth redirect handling
- **E2E:** Login with Google, upload and verify file

---

## Configuration Notes

### service-user Config
- OAuth credentials already configured in service-user
- File storage path in `config.yaml`

### Gateway Config (New)
```env
# OAuth
OAUTH_REDIRECT_URI=http://gateway:8080/api/v1/auth/google/callback
FRONTEND_REDIRECT_URI=http://localhost:3000/login

# File Upload
MAX_FILE_SIZE=10485760  # 10MB
ALLOWED_FILE_TYPES=image/*,application/pdf,text/*
```

---

## Open Questions / Decisions Needed

None - design approved.
