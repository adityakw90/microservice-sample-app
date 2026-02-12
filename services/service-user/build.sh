#!/bin/bash
set -e

# sementara ambil binary dari folder service-user
# nanti kalo service-user udah release, download dari github

PROJECT_DIR="/media/adit/SSD/project/MTAmedia/repo/microservice-sample-app/services/service-user"
SERVICE_USER_DIR="/media/adit/SSD/project/MTAmedia/repo/service-user/service-user"

# build binary
echo "Building service-user binary..."
cd "$SERVICE_USER_DIR" && make build

# copy binary ke folder service-user
echo "Copying binary to $PROJECT_DIR..."
cp "$SERVICE_USER_DIR/bin/app" "$PROJECT_DIR/server"

# build docker image
echo "Building Docker image..."
cd "$PROJECT_DIR"
docker build -t service-user:latest .

echo "✅ Done! service-user:latest image created"
