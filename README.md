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
│  ├─ web/                   # Vue / React app
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

- `service-user`
- `service-access-control`

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

The frontend application:

- Communicates only with the API Aggregator
- Is unaware of individual microservices
- Treats the backend as a single logical system

This keeps frontend logic simple and backend evolution flexible.

---

## Local Development

### Prerequisites

- Docker & Docker Compose
- Make (optional)

### Running the Stack

```bash
docker-compose up -d
```

This will start:

- API Aggregator
- All referenced microservices
- PostgreSQL
- Redis
- RabbitMQ
- Frontend

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
