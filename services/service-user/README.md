# Service User Docker Build

This directory contains the Dockerfile for building the service-user microservice.

## Prerequisites

The service-user microservice requires proto files from the service-user-proto repository. Ensure the proto files are available before building.

## Building the Docker Image

### Option 1: Using the provided build script

```bash
# From the microservice-sample-app root
./scripts/build-service-user.sh
```

### Option 2: Manual build with proto files

```bash
# Copy proto files first
cp -r ../../../service-user/service-user-proto/proto ./services/service-user/

# Build the image
docker build -t service-user:latest ./services/service-user
```

### Option 3: Using docker-compose with build context

The docker-compose.yml is configured to build the service-user with the necessary proto files included.

```bash
docker-compose build service-user
```

## Docker Image Details

- **Base Image**: golang:1.23-alpine (builder), alpine:latest (runtime)
- **Exposed Port**: 50051 (gRPC)
- **Health Check**: grpc_health_probe on port 50051
- **Binary**: `/root/service-user`

## Environment Variables

The service requires the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | - |
| `REDIS_URL` | Redis connection string | - |
| `RABBITMQ_URL` | RabbitMQ connection string (optional) | - |
| `GRPC_PORT` | gRPC server port | 50051 |
| `JWT_SECRET` | JWT signing secret | - |
| `PASSWORD_HASHER` | Password hashing algorithm | argon2id |
| `PIN_HASHER` | PIN hashing algorithm | argon2id |
| `LOCKOUT_MAX_ATTEMPTS` | Max login attempts before lockout | 5 |
| `LOCKOUT_DURATION` | Account lockout duration | 30m |

## Running the Container

```bash
docker run -d \
  --name service-user \
  -p 50051:50051 \
  -e DATABASE_URL=postgres://user:pass@localhost:5432/db \
  -e REDIS_URL=redis://localhost:6379 \
  -e JWT_SECRET=your-secret-key \
  service-user:latest
```

## Development

For local development without Docker:

```bash
cd ../../../service-user/service-user
go run ./cmd/server/main.go
```

## Troubleshooting

### Build fails with "cannot find package"

Ensure the proto files are copied to the correct location before building:
```bash
ls ./services/service-user/proto/
# Should contain: auth.proto, common.proto, user.proto, etc.
```

### Container exits immediately

Check the logs for configuration errors:
```bash
docker logs service-user
```

Ensure all required environment variables are set.
