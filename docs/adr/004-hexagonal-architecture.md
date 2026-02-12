# ADR 004: Hexagonal Architecture (Ports and Adapters)

## Status
Accepted

## Context
We need to build maintainable, testable microservices where business logic is independent of infrastructure concerns. Traditional layered architecture often leads to:
- Business logic leaking into infrastructure code
- Difficult to test without external dependencies
- Tight coupling to frameworks and databases
- Hard to swap infrastructure components

## Decision
Implement **Hexagonal Architecture** (Ports and Adapters) for all backend services:

### Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│                   Application Core                      │
│  (Business Logic, Use Cases, Domain Models)             │
└────────────────────┬────────────────────────────────────┘
                     │
    ┌────────────────┼────────────────┐
    │                │                │
┌───▼────┐    ┌─────▼─────┐    ┌────▼────┐
│Primary │    │  Primary  │    │Primary  │
│ Ports  │    │  Ports    │    │ Ports   │
│(Driving│    │ (Driving  │    │(Driving │
│Adapters)│    │ Adapters) │    │Adapters)│
├────────┤    ├───────────┤    ├─────────┤
│HTTP    │    │ gRPC      │    │Event    │
│Handler │    │ Server    │    │Consumer │
└────────┘    └───────────┘    └─────────┘
                     │
    ┌────────────────┼────────────────┐
    │                │                │
┌───▼────┐    ┌─────▼─────┐    ┌────▼────┐
│Secondary│    │ Secondary │    │Secondary│
│ Ports  │    │  Ports    │    │ Ports   │
│(Driven │    │ (Driven   │    │(Driven  │
│Adapters)│    │ Adapters) │    │Adapters)│
├────────┤    ├───────────┤    ├─────────┤
│Database│    │ Message   │    │External │
│Client  │    │ Queue     │    │API      │
└────────┘    └───────────┘    └─────────┘
```

### Implementation in Go (API Gateway)

```
api-gateway/
├── internal/
│   ├── handler/      # Primary Adapters (Driving) - HTTP handlers
│   ├── service/      # Application Core - Business logic, use cases
│   ├── client/       # Secondary Adapters (Driven) - gRPC clients
│   ├── model/        # Domain models and DTOs
│   └── config/       # Configuration
```

### Key Principles

1. **Dependency Inversion**: High-level modules don't depend on low-level modules
2. **Port Interface**: Application defines interfaces (ports) for external dependencies
3. **Adapter Implementation**: Infrastructure implements adapters for ports
4. **Business Logic Independence**: Core doesn't know about HTTP, gRPC, databases

### Directory Structure Convention

```
service-name/
├── internal/
│   ├── domain/        # Entities, value objects, domain logic
│   ├── application/   # Use cases, application services
│   ├── infrastructure/ # External dependencies (DB, APIs)
│   └── interfaces/    # Controllers, presenters (primary adapters)
```

## Consequences

### Positive
- **Testability**: Business logic tested without infrastructure
- **Maintainability**: Clear separation of concerns
- **Flexibility**: Easy to swap implementations (e.g., Redis → Memcached)
- **Independence**: Domain logic not coupled to frameworks
- **Parallel Development**: Interfaces allow parallel work

### Negative
- **Complexity**: More files and indirection than simple layered architecture
- **Learning curve**: Team must understand ports/adapters concept
- **Overhead**: More boilerplate for simple CRUD operations
- **Debugging**: More layers to trace through

### Mitigations
- Start simple: Only create ports when multiple implementations exist
- Clear naming conventions: handler → service → client
- Dependency injection in main.go for clarity
- Documentation showing data flow through layers

## Examples

### API Gateway Layer Flow

```go
// Primary Adapter (Driving)
handler (user_handler.go)
    → Converts HTTP to internal request
    → Calls service

// Application Core
service (user_service.go)
    → Business logic, validation
    → Calls client interface

// Secondary Adapter (Driven)
client (user_client.go)
    → Implements gRPC communication
    → Returns domain objects
```

## References

- Alistair Cockburn's Hexagonal Architecture
- https://alistair.cockburn.us/hexagonal-architecture/
- https://herbertograca.com/2017/09/14/ports-adapters-architecture/
