# API Gateway

The API Gateway acts as a Backend For Frontend (BFF) in the microservice architecture. It provides HTTP/REST endpoints for the frontend while communicating with backend microservices via gRPC.

## Architecture

```
Frontend
   ↓ (HTTP/REST)
API Gateway (this service)
   ↓ (gRPC)
User Service
```

## Features

- **HTTP to gRPC translation**: Converts REST API calls to gRPC calls
- **Aggregation**: Combines data from multiple microservices
- **Single entry point**: Frontend only needs to know about the API Gateway
- **Service discovery**: Connects to microservices via configurable addresses

## API Endpoints

### Authentication

- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/validate` - Validate access token
- `POST /api/v1/auth/verify-pin` - Verify user PIN

### Users

- `GET /api/v1/users` - List users (with pagination and filtering)
- `GET /api/v1/users/{uid}` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/{uid}` - Update user
- `DELETE /api/v1/users/{uid}` - Delete user

### User Profiles

- `GET /api/v1/users/{uid}/profile` - Get user profile
- `PUT /api/v1/users/{uid}/profile` - Update user profile

### User Devices

- `GET /api/v1/users/{uid}/devices` - List user devices
- `DELETE /api/v1/users/{uid}/devices/{deviceUid}` - Revoke device

### Health Check

- `GET /health` - Health check endpoint

## Configuration

The API Gateway is configured via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `USER_SERVICE_ADDRESS` | Address of the user gRPC service | `localhost:50051` |
| `SERVER_PORT` | HTTP server port | `8080` |

## Running Locally

### Prerequisites

- Go 1.22+
- Access to the service-user gRPC service

### Run directly

```bash
cd api-gateway
export USER_SERVICE_ADDRESS=localhost:50051
go run main.go
```

### Run with Docker

```bash
# From the project root
docker-compose up api-gateway
```

## Development

### Project Structure

```
api-gateway/
├── internal/
│   ├── client/      # gRPC clients
│   ├── handler/     # HTTP request handlers
│   ├── service/     # Business logic layer
│   └── config/      # Configuration
├── main.go          # Application entry point
├── Dockerfile       # Docker build file
└── go.mod           # Go module definition
```

### Adding New Endpoints

1. Define the request/response types in `service/`
2. Add the business logic in `service/*_service.go`
3. Create the HTTP handler in `handler/*_handler.go`
4. Register the route in the handler's `RegisterRoutes` method
5. Wire up the handler in `main.go`

## Example Requests

### Create a user

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "johndoe",
    "identifier_type": "username",
    "password": "securepassword"
  }'
```

### List users

```bash
curl http://localhost:8080/api/v1/users?page=1&limit=10
```

### Get user profile

```bash
curl http://localhost:8080/api/v1/users/{uid}/profile
```
