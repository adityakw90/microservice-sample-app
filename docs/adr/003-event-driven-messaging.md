# ADR 003: Event-Driven Messaging with RabbitMQ and Kafka

## Status
Accepted

## Context
Microservices need to communicate asynchronously for:
- Background job processing
- Event notification between services
- Decoupling of long-running operations

## Decision
Implement dual messaging infrastructure:
- **RabbitMQ** for traditional message queuing (work queues, pub/sub)
- **Kafka** for event streaming and log aggregation

## Consequences
### Positive
- **Loose coupling** - services don't need to be online simultaneously
- **Scalability** - workers can scale independently
- **Resilience** - messages persist during outages
- **Flexibility** - RabbitMQ for work queues, Kafka for event streams

### Negative
- **Infrastructure complexity** - two messaging systems to operate
- **Operational overhead** - monitoring, scaling, backup
- **Message ordering challenges** - need proper design for ordering guarantees

### Mitigations
- Docker Compose for local development parity
- Health checks and proper service dependencies
- Clear documentation on when to use RabbitMQ vs Kafka

## References
- RabbitMQ: https://www.rabbitmq.com/
- Kafka: https://kafka.apache.org/
