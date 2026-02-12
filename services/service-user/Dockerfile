# Runtime image for service-user (Debian for glibc compatibility)
FROM debian:bookworm-slim

# Install ca-certificates for HTTPS and netcat for health checks
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl wget netcat-openbsd && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the pre-built binary from build.sh
COPY server .
COPY config.yaml .

# Expose gRPC port
EXPOSE 50051

# Health check - TCP connection check
HEALTHCHECK --interval=10s --timeout=5s --retries=5 --start-period=30s \
  CMD nc -z localhost 50051 || exit 1

# Run the application
CMD ["/app/server"]
