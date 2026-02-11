# Deployment Documentation

## Local Development

### Prerequisites
- Docker & Docker Compose
- Go 1.23+ (for building service-user)
- Node.js 20+ (for frontend development)

### Building service-user

service-user requires proto files from external repository before building:

```bash
# Option 1: Use build script
./scripts/build-service-user.sh

# Option 2: Manual build
cp -r ../../../service-user/service-user-proto/proto ./services/service-user/
docker build -t service-user:latest ./services/service-user
```

### Starting the Stack

```bash
docker-compose up -d
```

This starts:
- Frontend (port 8000)
- API Gateway (port 8080)
- service-user (port 50051)
- PostgreSQL (port 5432)
- Redis (port 6379)
- RabbitMQ (ports 5672, 15672)
- Kafka (port 9092)
- pgAdmin (port 5050)

### Access Points

- **Frontend**: http://localhost:8000
- **API Gateway**: http://localhost:8080
- **API Health**: http://localhost:8080/health
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)
- **pgAdmin**: http://localhost:5050 (admin@localhost.id / password)

### Stopping the Stack

```bash
docker-compose down      # Stop containers
docker-compose down -v   # Stop + remove volumes
```

## Infrastructure Services

### PostgreSQL
- **Image**: postgres:18-alpine
- **Database**: microservice_db
- **Credentials**: postgres/postgres
- **Volume**: `postgres_data:/var/lib/postgresql` (note: not `/var/lib/postgresql/data` due to permission issues)

### Redis
- **Image**: redis:8.4-alpine
- **No authentication** (development only)

### RabbitMQ
- **Image**: rabbitmq:4.2.3-management-alpine
- **Credentials**: rabbit/password
- **Management UI**: port 15672

### Kafka
- **Image**: confluentinc/cp-kafka:latest
- **Mode**: KRaft (no Zookeeper)
- **Port**: 9092 (PLAINTEXT)

## Production Considerations

### Configuration Changes Required
1. Enable Redis authentication
2. Use strong passwords for all services
3. Enable TLS for inter-service communication
4. Configure proper volume mounts for persistence
5. Set up proper monitoring and logging
6. Use external managed services where appropriate

### Health Checks
All services include health checks:
- PostgreSQL: `pg_isready`
- Redis: `redis-cli ping`
- RabbitMQ: `rabbitmq-diagnostics ping`
- service-user: `nc -z localhost 50051`
- API Gateway: HTTP `/health` endpoint
- Frontend: HTTP spider check

## Service Dependencies

```
Frontend
  ↓ (health check)
API Gateway
  ↓ (service started)
service-user
  ↓ (health checks)
├── PostgreSQL
├── Redis
└── RabbitMQ
```
