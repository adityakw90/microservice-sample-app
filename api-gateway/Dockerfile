# Runtime image for service-user (Debian for glibc compatibility)
FROM debian:bookworm-slim

# Install ca-certificates for HTTPS and grpc_health_probe for health checks
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl wget && rm -rf /var/lib/apt/lists/*


WORKDIR /app

# Copy the pre-built binary from build.sh
COPY server .

# Expose rest port
EXPOSE 8080

# Run the application
CMD ["/app/server"]
