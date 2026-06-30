#!/bin/bash

# Cupcake Local Test Script
# This script runs the backend locally for testing (without Docker)

echo "🧁 Starting Cupcake Backend Locally..."

# Set environment variables for local testing
export MONGO_URI="mongodb://cupcake:cupcake123@localhost:27017/cupcake?authSource=admin"
export KAFKA_BROKER="localhost:9093"
export PORT="8080"
export GIN_MODE="debug"

echo "📊 MongoDB URI: $MONGO_URI"
echo "📨 Kafka broker: $KAFKA_BROKER"
echo "🌐 Server port: $PORT"
echo ""

# Navigate to backend directory
cd backyard

# Build the binary
echo "🔨 Building backend..."
go build -o bin/backyard cmd/main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "🚀 Starting server..."
    echo "📖 Swagger UI will be available at: http://localhost:8080/swagger/index.html"
    echo "🔧 Health Check: http://localhost:8080/health"
    echo ""
    echo "Press Ctrl+C to stop the server"
    echo ""
    
    # Run the backend
    ./bin/backyard
else
    echo "❌ Build failed!"
    exit 1
fi
