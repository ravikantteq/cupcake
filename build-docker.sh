#!/bin/bash

# Cupcake Docker Build Script
# This script builds the Docker images for the Cupcake Kafka Test Framework

echo "🧁 Building Cupcake Docker Images..."

# Set environment to disable Docker BuildKit if needed
export DOCKER_BUILDKIT=0

# Build backend image
echo "📦 Building cupcake-backend..."
docker build -t cupcake-backend -f backyard/Dockerfile backyard/

if [ $? -eq 0 ]; then
    echo "✅ Backend build successful!"
else
    echo "❌ Backend build failed!"
    exit 1
fi

# Build frontend image
echo "🎨 Building cupcake-ui..."
docker build -t cupcake-ui -f cupcake_ui/Dockerfile cupcake_ui/

if [ $? -eq 0 ]; then
    echo "✅ Frontend build successful!"
else
    echo "❌ Frontend build failed!"
    exit 1
fi

echo "🎉 All Docker images built successfully!"
echo ""
echo "To start the services:"
echo "  docker-compose up -d"
echo ""
echo "To start individual services:"
echo "  docker run -p 8080:8080 --network cupcake_default cupcake-backend"
echo "  docker run -p 4200:80 cupcake-ui"
