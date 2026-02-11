# ADR 001: Backend For Frontend (BFF) Pattern

## Status
Accepted

## Context
We need to expose multiple gRPC microservices to a frontend application. Direct gRPC from the browser is complex, and calling multiple services from the frontend creates tight coupling and complexity.

## Decision
Implement an API Gateway using the Backend For Frontend (BFF) pattern:
- Frontend communicates only with the API Gateway via HTTP/REST
- API Gateway aggregates calls to multiple gRPC microservices
- Gateway handles orchestration, caching, and response shaping

## Consequences
### Positive
- **Single API surface** for frontend - simpler integration
- **Reduced frontend complexity** - no gRPC handling, no multiple service URLs
- **Flexible backend evolution** - can change service topology without affecting frontend
- **Centralized cross-cutting concerns** - auth, logging, caching at gateway layer

### Negative
- **Additional network hop** - frontend → gateway → service
- **Gateway becomes critical path** - single point of failure
- **Potential gateway bloat** - business logic can creep into gateway

### Mitigations
- Gateway is stateless and horizontally scalable
- Health checks and graceful degradation
- Clear separation: handlers → services → clients (business logic in service layer)

## References
- https://samnewman.io/patterns/architectural/bff/
