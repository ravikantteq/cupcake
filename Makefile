# Cupcake Project Makefile

.PHONY: help backend frontend test-backend test-frontend build-backend build-frontend clean docs install-deps

# Default target
help:
	@echo "Cupcake Project - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  backend         - Start the Go backend server"
	@echo "  frontend        - Start the Angular frontend"
	@echo "  dev             - Start both backend and frontend"
	@echo ""
	@echo "Testing:"
	@echo "  test-backend    - Run Go tests"
	@echo "  test-frontend   - Run Angular tests"
	@echo "  test            - Run all tests"
	@echo ""
	@echo "Build:"
	@echo "  build-backend   - Build the Go backend"
	@echo "  build-frontend  - Build the Angular frontend"
	@echo "  build           - Build both backend and frontend"
	@echo ""
	@echo "Documentation:"
	@echo "  docs            - Generate Swagger documentation"
	@echo ""
	@echo "Utilities:"
	@echo "  install-deps    - Install all dependencies"
	@echo "  clean           - Clean build artifacts"

# Install dependencies
install-deps:
	@echo "Installing Go dependencies..."
	cd backyard && go mod tidy
	@echo "Installing Angular dependencies..."
	cd cupcake_ui && npm install

# Backend commands
backend:
	@echo "Starting backend server..."
	cd backyard && go run cmd/main.go

build-backend:
	@echo "Building backend..."
	cd backyard && go build -o bin/backyard cmd/main.go

test-backend:
	@echo "Running backend tests..."
	cd backyard && go test ./...

# Frontend commands
frontend:
	@echo "Starting frontend development server..."
	cd cupcake_ui && npm start

build-frontend:
	@echo "Building frontend..."
	cd cupcake_ui && npm run build

test-frontend:
	@echo "Running frontend tests..."
	cd cupcake_ui && npm test

# Combined commands
dev:
	@echo "Starting both backend and frontend..."
	@echo "Backend will start on :8080, Frontend on :4200"
	@make backend &
	@make frontend

build: build-backend build-frontend

test: test-backend test-frontend

# Documentation
docs:
	@echo "Generating Swagger documentation..."
	cd backyard && ~/go/bin/swag init -g cmd/main.go

# Clean up
clean:
	@echo "Cleaning build artifacts..."
	cd backyard && rm -rf bin/
	cd cupcake_ui && rm -rf dist/
	cd cupcake_ui && rm -rf node_modules/

# Docker commands (for future use)
docker-backend:
	@echo "Building backend Docker image..."
	cd backyard && docker build -t cupcake-backend .

docker-frontend:
	@echo "Building frontend Docker image..."
	cd cupcake_ui && docker build -t cupcake-frontend .
