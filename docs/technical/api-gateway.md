# API Gateway Technical Documentation

## Overview

The API Gateway serves as the Backend For Frontend (BFF), exposing HTTP/REST endpoints to the frontend and translating requests to gRPC calls to backend microservices.

## Architecture

```
HTTP Request
    ↓
Middleware (CORS, Recovery, Logging)
    ↓
Handler (HTTP Layer)
    ↓
Service (Business Logic Layer)
    ↓
Client (gRPC Layer)
    ↓
Backend Service
```

## Directory Structure

```
api-gateway/
├── cmd/
│   └── gateway/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── handler/              # HTTP request handlers
│   ├── service/              # Business logic & orchestration
│   ├── client/               # gRPC client connections
│   ├── middleware/           # Custom middleware
│   └── model/                # Data transfer objects
├── Dockerfile
└── README.md
```

## Configuration

Environment variables:
- `USER_SERVICE_ADDRESS` - gRPC address of service-user (default: service-user:50051)
- `SERVER_PORT` - HTTP port for gateway (default: 8080)

## API Endpoints

### Authentication (`/api/v1/auth`)
- `POST /login` - Authenticate user, return JWT tokens
- `POST /refresh` - Refresh access token using refresh token
- `POST /validate` - Validate access token
- `POST /verify-pin` - Verify user PIN for sensitive actions

### Users (`/api/v1/users`)
- `GET /` - List users with pagination and filtering
- `GET /{uid}` - Get user by ID
- `POST /` - Create new user
- `PUT /{uid}` - Update user details
- `DELETE /{uid}` - Delete user
- `GET /{uid}/profile` - Get user profile
- `PUT /{uid}/profile` - Update user profile
- `GET /{uid}/devices` - List user devices
- `DELETE /{uid}/devices/{deviceUid}` - Revoke a device

## Middleware Chain

1. **CORS** - Cross-origin resource sharing
2. **Recovery** - Panic recovery
3. **Logging** - Request/response logging
4. **Authentication** (future) - JWT validation

## Health Check

`GET /health` - Returns service health status

## Error Handling

All errors follow consistent JSON format:
```json
{
  "error": "error message",
  "code": "ERROR_CODE"
}
```
