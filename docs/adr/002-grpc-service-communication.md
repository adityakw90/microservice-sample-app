# ADR 002: gRPC for Service-to-Service Communication

## Status
Accepted

## Context
Microservices need to communicate with each other. Options include REST/HTTP, gRPC, message queues, or GraphQL. We need a protocol that is type-safe, performant, and supports contract-first development.

## Decision
Use gRPC with Protocol Buffers for all service-to-service communication:
- Services expose gRPC endpoints defined in `.proto` files
- API Gateway translates HTTP/REST to gRPC
- Strong typing through protobuf contracts

## Consequences
### Positive
- **Type safety** - compile-time contract validation
- **Performance** - binary serialization, HTTP/2 multiplexing
- **Contract-first** - proto files are the single source of truth
- **Polyglot support** - code generation for any language
- **Built-in streaming** - support for unary and streaming RPCs

### Negative
- **Browser support** - gRPC-Web required for direct browser access (we avoid this via BFF)
- **Debugging complexity** - binary format less human-readable than JSON
- **Tooling overhead** - need protoc compiler and code generation

### Mitigations
- API Gateway handles HTTP/gRPC translation - frontend uses REST
- grpcurl or similar tools for debugging
- Proto files in shared repository for versioning

## References
- https://grpc.io/docs/what-is-grpc/
