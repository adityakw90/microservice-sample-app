# Documentation Index

This directory contains architecture and technical documentation for the microservice-sample-app project.

## Architecture Decision Records (ADR/)

Key architectural decisions and their rationale:

| ADR | Title | Status |
|-----|-------|--------|
| [000](adr/000-use-adr-template.md) | ADR Template | Template |
| [001](adr/001-bff-pattern.md) | Backend For Frontend Pattern | Accepted |
| [002](adr/002-grpc-service-communication.md) | gRPC for Service-to-Service Communication | Accepted |
| [003](adr/003-event-driven-messaging.md) | Event-Driven Messaging (RabbitMQ & Kafka) | Accepted |
| [004](adr/004-hexagonal-architecture.md) | Hexagonal Architecture (Ports and Adapters) | Accepted |

## Technical Documentation (technical/)

Detailed technical documentation:

| Document | Description |
|----------|-------------|
| [api-gateway.md](technical/api-gateway.md) | API Gateway architecture, endpoints, and configuration |
| [frontend.md](technical/frontend.md) | Frontend architecture, tech stack, and development |
| [deployment.md](technical/deployment.md) | Local development setup and deployment guide |

## Adding New Documentation

### New ADR
1. Copy `docs/adr/000-use-adr-template.md`
2. Name with sequential number: `NNN-title.md`
3. Fill out all sections
4. Update this index

### New Technical Doc
1. Create in `docs/technical/`
2. Use descriptive `.md` filename
3. Update this index
