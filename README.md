# microservice-sample-app

## Overview

`microservice-sample-app` is an **end-to-end reference application** that demonstrates how independently developed microservices can be composed into a cohesive system.

This repository focuses on **architecture, boundaries, and integration**, not on reimplementing individual services. Each microservice is treated as an independent deployable unit, mirroring a real production environment.

The system includes:

- Multiple backend **gRPC microservices**
- An **API Aggregator (BFF)** exposed to the frontend
- **Async background workers**
- Shared infrastructure (PostgreSQL, Redis, RabbitMQ)
- A **frontend application** consuming the system as a whole

---

## Goals of This Repository

- Showcase a **production-style microservice architecture**
- Demonstrate clean **service boundaries and contracts**
- Show how an **API Aggregator / BFF** simplifies frontend consumption
- Provide a realistic **local development setup** using Docker
- Act as a **reference implementation**, not a tutorial toy project

---

## High-Level Architecture

```
Frontend
   ↓ (HTTP)
API Aggregator (BFF)
   ↓ (gRPC)
---------------------------------
|  User Service                  |
|  Access Control Service (IAM)  |
|  Exchange Service              |
|  Other Domain Services         |
---------------------------------
   ↓
Shared Infrastructure
(PostgreSQL, Redis, RabbitMQ)
```

### Key Architectural Decisions

- **gRPC for service-to-service communication**
- **API Aggregator (BFF)** as the single entry point for frontend
- **No direct frontend → microservice calls**
- **Services own their data and schema**
- **Async workers** for background and long-running tasks

---

## Repository Structure

```
microservice-sample-app/
├─ services/                 # References to existing microservices
│  └─ README.md              # Service list and image references
│
├─ api-gateway/              # API Aggregator / BFF
│  ├─ internal/
│  │  ├─ handler/            # HTTP handlers
│  │  ├─ service/            # Aggregation logic
│  │  └─ client/             # gRPC clients
│  └─ main.go
│
├─ frontend/                 # Frontend application
│  └─ README.md
│
├─ proto/                    # Shared gRPC contracts
│  ├─ user.proto
│  └─ access_control.proto
│
├─ infra/                    # Local infrastructure setup
│  ├─ docker-compose.yml
│  ├─ postgres/
│  ├─ redis/
│  └─ rabbitmq/
│
├─ docs/                     # Architecture documentation
│  ├─ architecture.md
│  ├─ service-boundaries.md
│  └─ data-flow.md
│
├─ env.example
└─ README.md
```

---

## Services

This repository **does not reimplement** microservices.

Instead, it consumes existing services such as:

- `service-user` - User management and authentication
- `service-access-control` - IAM and permissions

### Service User Showcase

The current showcase demonstrates integration with the **service-user** microservice, which provides:

**User Management:**

- CRUD operations for users
- User profile management
- User device tracking and revocation

**Authentication:**

- Login with username or email
- JWT token issuance (access + refresh tokens)
- Token validation
- PIN verification for sensitive actions

**gRPC Endpoints:**

- `UserService.List` - List users with pagination and filtering
- `UserService.Get` - Get user by UID
- `UserService.Add` - Create new user
- `UserService.Update` - Update user details
- `UserService.Delete` - Delete user
- `UserService.GetProfile` - Get user profile
- `UserService.UpdateProfile` - Update user profile
- `UserService.ListDevice` - List user devices
- `UserService.RevokeDevice` - Revoke a device

- `AuthService.Auth` - Authenticate user
- `AuthService.RefreshToken` - Refresh access token
- `AuthService.ValidateToken` - Validate access token
- `AuthService.VerifyPin` - Verify user PIN

Each service is expected to be:

- Independently versioned
- Independently deployable
- Accessed only through **gRPC contracts**

Services are integrated via **Docker images** (or optionally git submodules for development).

---

## API Aggregator (BFF)

The API Aggregator acts as a **Backend For Frontend**:

- Exposes HTTP/REST APIs tailored for frontend needs
- Aggregates data from multiple gRPC services
- Handles orchestration, caching, and response shaping
- Shields the frontend from service topology changes

The frontend never communicates directly with gRPC services.

---

## Frontend

The frontend application is a **Vue 3 SPA** built with TypeScript:

- **Communicates only with the API Aggregator** via HTTP/REST
- **Is unaware of individual microservices** - treats backend as a single system
- **JWT-based authentication** with automatic token refresh
- **Pinia state management** for auth and data
- **Modern component architecture** with Vue 3 Composition API

### Frontend Features

**Authentication:**
- Login with username or email
- JWT token management (access + refresh)
- Automatic token refresh on expiration
- Protected routes with auth guards

**User Management:**
- Browse users with pagination and search
- Create new users
- View user profiles with full details
- Edit user information
- Delete users

**User Profile:**
- View and edit personal profile
- Device management (view and revoke)
- Real-time profile updates

### Frontend Tech Stack

- **Vue 3** with Composition API
- **TypeScript** for type safety
- **Vite** for fast development and building
- **Vue Router** for navigation
- **Pinia** for state management
- **Axios** for HTTP requests with interceptors

This keeps frontend logic simple and backend evolution flexible.

---

## Local Development

### Prerequisites

- Docker & Docker Compose
- Go 1.23+ (for building service-user)
- Node.js 20+ (for frontend development)
- Make (optional)

### Building service-user

The service-user microservice requires proto files from the service-user-proto repository. Before running the stack, build the service-user Docker image:

```bash
# Option 1: Use the provided build script (recommended)
./scripts/build-service-user.sh

# Option 2: Manual build
# 1. Copy proto files
cp -r ../../../service-user/service-user-proto/proto ./services/service-user/
# 2. Build Docker image
docker build -t service-user:latest ./services/service-user
```

See `services/service-user/README.md` for more details.

### Running the Stack

```bash
docker-compose up -d
```

This will start:

- **Frontend** (port 80) - Vue 3 application
- **API Gateway** (port 8080) - HTTP/REST entry point
- **Service User** (port 50051) - User management gRPC service
- **PostgreSQL** (port 5432) - Database
- **Redis** (port 6379) - Cache
- **RabbitMQ** (ports 5672, 15672) - Message broker

### Accessing the Application

Once the stack is running:

- **Frontend**: http://localhost
- **API Gateway**: http://localhost:8080
- **API Docs**: See `api-gateway/README.md`
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)

### Testing the API

Once the stack is running, you can test the API:

```bash
# Health check
curl http://localhost:8080/health

# Create a user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"testuser","identifier_type":"username","password":"password123"}'
```

### Frontend Development

For local development of the frontend:

```bash
cd frontend
npm install
npm run dev
```

The frontend dev server will run on `http://localhost:3000` and proxy API requests to the API Gateway.

---

## What This Repo Is (and Is Not)

### This repo **is**:

- A reference architecture
- A system integration example
- A realistic microservice composition

### This repo **is not**:

- A single-service codebase
- A step-by-step beginner tutorial
- A place to implement domain logic

---

## Intended Audience

- Backend engineers designing microservice systems
- Engineers evaluating gRPC-based architectures
- Teams looking for a realistic reference implementation

---

## License

MIT (or your preferred license)
