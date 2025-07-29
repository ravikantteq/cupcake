# 🧁 Cupcake Kafka Test Framework - Enhanced Makefile
# Version: 2.0
# Supports: Docker, MongoDB, Test Flows, Consumer Management

.DEFAULT_GOAL := help
SHELL := /bin/bash

# Colors for pretty output
COLOR_RESET = \033[0m
COLOR_BOLD = \033[1m
COLOR_GREEN = \033[32m
COLOR_YELLOW = \033[33m
COLOR_BLUE = \033[34m
COLOR_CYAN = \033[36m

# Project configuration
PROJECT_NAME := cupcake
BACKEND_DIR := backyard
FRONTEND_DIR := cupcake_ui
DOCKER_COMPOSE := docker-compose
LOGS_DIR := logs

# Go configuration
GO_MODULE := github.com/systemgenes/cupcake/$(BACKEND_DIR)
GO_VERSION := 1.20

# Node/Angular configuration
NODE_VERSION := 18
ANGULAR_VERSION := 19

.PHONY: help install-deps dev build test clean docker docs

## 📋 Help and Information
help: ## Display this help message
	@echo "$(COLOR_BOLD)🧁 Cupcake Kafka Test Framework$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)Enterprise-ready Kafka testing platform$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)📚 Available Commands:$(COLOR_RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(COLOR_BOLD)🌐 Service URLs:$(COLOR_RESET)"
	@echo "  $(COLOR_BLUE)Frontend:$(COLOR_RESET)     http://localhost:4200"
	@echo "  $(COLOR_BLUE)Backend API:$(COLOR_RESET)  http://localhost:8080"
	@echo "  $(COLOR_BLUE)API Docs:$(COLOR_RESET)     http://localhost:8080/swagger/index.html"
	@echo "  $(COLOR_BLUE)Kafka UI:$(COLOR_RESET)     http://localhost:8081"
	@echo "  $(COLOR_BLUE)Mongo UI:$(COLOR_RESET)     http://localhost:8082"

status: ## Show status of all services
	@echo "$(COLOR_BOLD)📊 Service Status:$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) ps

## 🚀 Development Environment
dev: ## Start complete development environment (recommended)
	@echo "$(COLOR_BOLD)🚀 Starting Cupcake development environment...$(COLOR_RESET)"
	@$(MAKE) ensure-logs-dir
	@$(DOCKER_COMPOSE) up -d
	@echo "$(COLOR_GREEN)✅ All services started!$(COLOR_RESET)"
	@echo ""
	@$(MAKE) dev-info

dev-info: ## Display development environment information
	@echo "$(COLOR_BOLD)🎯 Development Environment Ready!$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)📱 Frontend:$(COLOR_RESET)     http://localhost:4200"
	@echo "$(COLOR_BOLD)🔧 Backend API:$(COLOR_RESET)  http://localhost:8080"
	@echo "$(COLOR_BOLD)📖 API Docs:$(COLOR_RESET)     http://localhost:8080/swagger/index.html"
	@echo "$(COLOR_BOLD)📊 Kafka UI:$(COLOR_RESET)     http://localhost:8081"
	@echo "$(COLOR_BOLD)🗄️ Database UI:$(COLOR_RESET) http://localhost:8082 (admin/admin123)"
	@echo ""
	@echo "$(COLOR_YELLOW)💡 Use 'make logs' to follow service logs$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)💡 Use 'make stop' to stop all services$(COLOR_RESET)"

dev-local: ## Start backend and frontend locally (without Docker)
	@echo "$(COLOR_BOLD)🔧 Starting local development...$(COLOR_RESET)"
	@$(MAKE) backend &
	@$(MAKE) frontend &
	@echo "$(COLOR_GREEN)✅ Local services started!$(COLOR_RESET)"

## 🛠️ Individual Service Management
backend: ## Start backend service only
	@echo "$(COLOR_BOLD)🔧 Starting backend service...$(COLOR_RESET)"
	@cd $(BACKEND_DIR) && go run cmd/main.go

frontend: ## Start frontend service only
	@echo "$(COLOR_BOLD)🎨 Starting frontend service...$(COLOR_RESET)"
	@cd $(FRONTEND_DIR) && npm start

kafka: ## Start only Kafka infrastructure
	@echo "$(COLOR_BOLD)🔄 Starting Kafka infrastructure...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) up -d zookeeper kafka kafka-ui

database: ## Start only MongoDB
	@echo "$(COLOR_BOLD)🗄️ Starting MongoDB...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) up -d mongodb mongo-express

## 📦 Installation and Dependencies
install-deps: ## Install all dependencies
	@echo "$(COLOR_BOLD)📦 Installing dependencies...$(COLOR_RESET)"
	@$(MAKE) install-backend-deps
	@$(MAKE) install-frontend-deps
	@echo "$(COLOR_GREEN)✅ All dependencies installed!$(COLOR_RESET)"

install-backend-deps: ## Install Go dependencies
	@echo "$(COLOR_BLUE)Installing Go dependencies...$(COLOR_RESET)"
	@cd $(BACKEND_DIR) && go mod tidy && go mod download

install-frontend-deps: ## Install Node.js dependencies
	@echo "$(COLOR_BLUE)Installing Node.js dependencies...$(COLOR_RESET)"
	@cd $(FRONTEND_DIR) && npm install

## 🐳 Docker Operations
docker-build: ## Build all Docker images
	@echo "$(COLOR_BOLD)🐳 Building Docker images...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) build
	@echo "$(COLOR_GREEN)✅ Docker images built!$(COLOR_RESET)"

docker-rebuild: ## Rebuild Docker images without cache
	@echo "$(COLOR_BOLD)🐳 Rebuilding Docker images...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) build --no-cache
	@echo "$(COLOR_GREEN)✅ Docker images rebuilt!$(COLOR_RESET)"

up: ## Start all services with Docker Compose
	@$(MAKE) dev

stop: ## Stop all Docker services
	@echo "$(COLOR_BOLD)🛑 Stopping all services...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) stop
	@echo "$(COLOR_GREEN)✅ All services stopped!$(COLOR_RESET)"

down: ## Stop and remove all containers, networks
	@echo "$(COLOR_BOLD)🧹 Stopping and removing containers...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) down
	@echo "$(COLOR_GREEN)✅ Containers removed!$(COLOR_RESET)"

restart: ## Restart all services
	@echo "$(COLOR_BOLD)🔄 Restarting all services...$(COLOR_RESET)"
	@$(MAKE) down
	@$(MAKE) dev

## 🗄️ Database Operations
db-init: ## Initialize MongoDB with sample data
	@echo "$(COLOR_BOLD)🗄️ Initializing database...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) exec mongodb mongosh cupcake --eval "load('/docker-entrypoint-initdb.d/mongo-init.js')"
	@echo "$(COLOR_GREEN)✅ Database initialized!$(COLOR_RESET)"

db-backup: ## Backup MongoDB data
	@echo "$(COLOR_BOLD)💾 Backing up database...$(COLOR_RESET)"
	@mkdir -p backups
	@$(DOCKER_COMPOSE) exec mongodb mongodump --db cupcake --out /tmp/backup
	@$(DOCKER_COMPOSE) cp mongodb:/tmp/backup ./backups/$(shell date +%Y%m%d_%H%M%S)
	@echo "$(COLOR_GREEN)✅ Database backed up!$(COLOR_RESET)"

db-restore: ## Restore MongoDB data (usage: make db-restore BACKUP=20240726_143000)
	@echo "$(COLOR_BOLD)📥 Restoring database...$(COLOR_RESET)"
	@if [ -z "$(BACKUP)" ]; then echo "$(COLOR_YELLOW)Usage: make db-restore BACKUP=20240726_143000$(COLOR_RESET)"; exit 1; fi
	@$(DOCKER_COMPOSE) cp ./backups/$(BACKUP) mongodb:/tmp/restore
	@$(DOCKER_COMPOSE) exec mongodb mongorestore --db cupcake --drop /tmp/restore/cupcake
	@echo "$(COLOR_GREEN)✅ Database restored!$(COLOR_RESET)"

db-shell: ## Connect to MongoDB shell
	@echo "$(COLOR_BOLD)🗄️ Connecting to MongoDB...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) exec mongodb mongosh cupcake

## 🧪 Testing
test: ## Run all tests
	@echo "$(COLOR_BOLD)🧪 Running all tests...$(COLOR_RESET)"
	@$(MAKE) test-backend
	@$(MAKE) test-frontend
	@echo "$(COLOR_GREEN)✅ All tests completed!$(COLOR_RESET)"

test-backend: ## Run backend tests
	@echo "$(COLOR_BLUE)Testing backend...$(COLOR_RESET)"
	@cd $(BACKEND_DIR) && go test -v ./...

test-frontend: ## Run frontend tests
	@echo "$(COLOR_BLUE)Testing frontend...$(COLOR_RESET)"
	@cd $(FRONTEND_DIR) && npm test -- --watch=false --browsers=ChromeHeadless

test-integration: ## Run integration tests
	@echo "$(COLOR_BLUE)Running integration tests...$(COLOR_RESET)"
	@cd $(BACKEND_DIR) && go test -v -tags=integration ./...

test-coverage: ## Generate test coverage reports
	@echo "$(COLOR_BLUE)Generating coverage reports...$(COLOR_RESET)"
	@cd $(BACKEND_DIR) && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html
	@cd $(FRONTEND_DIR) && npm run test:coverage

## 🏗️ Build Operations
build: ## Build all components
	@echo "$(COLOR_BOLD)🏗️ Building all components...$(COLOR_RESET)"
	@$(MAKE) build-backend
	@$(MAKE) build-frontend
	@echo "$(COLOR_GREEN)✅ Build completed!$(COLOR_RESET)"

build-backend: ## Build backend binary
	@echo "$(COLOR_BLUE)Building backend...$(COLOR_RESET)"
	@cd $(BACKEND_DIR) && go build -o bin/cupcake cmd/main.go

build-frontend: ## Build frontend for production
	@echo "$(COLOR_BLUE)Building frontend...$(COLOR_RESET)"
	@cd $(FRONTEND_DIR) && npm run build

## 📖 Documentation
docs: ## Generate API documentation
	@echo "$(COLOR_BOLD)📖 Generating documentation...$(COLOR_RESET)"
	@cd $(BACKEND_DIR) && swag init -g cmd/main.go
	@echo "$(COLOR_GREEN)✅ Documentation generated!$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)View at: http://localhost:8080/swagger/index.html$(COLOR_RESET)"

docs-serve: ## Serve documentation locally
	@echo "$(COLOR_BOLD)📖 Serving documentation...$(COLOR_RESET)"
	@cd $(BACKEND_DIR)/docs && python3 -m http.server 8090
	@echo "$(COLOR_BLUE)Documentation available at: http://localhost:8090$(COLOR_RESET)"

## 📊 Monitoring and Logs
logs: ## Follow all service logs
	@echo "$(COLOR_BOLD)📊 Following service logs...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) logs -f

logs-backend: ## Follow backend logs only
	@$(DOCKER_COMPOSE) logs -f cupcake-backend

logs-frontend: ## Follow frontend logs only
	@$(DOCKER_COMPOSE) logs -f cupcake-frontend

logs-kafka: ## Follow Kafka logs
	@$(DOCKER_COMPOSE) logs -f kafka

logs-mongodb: ## Follow MongoDB logs
	@$(DOCKER_COMPOSE) logs -f mongodb

health: ## Check health of all services
	@echo "$(COLOR_BOLD)🏥 Checking service health...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Backend Health:$(COLOR_RESET)"
	@curl -s http://localhost:8080/health | jq . || echo "Backend not responding"
	@echo ""
	@echo "$(COLOR_BLUE)Service Status:$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) ps

## 🧹 Cleanup Operations
clean: ## Clean up generated files and caches
	@echo "$(COLOR_BOLD)🧹 Cleaning up...$(COLOR_RESET)"
	@$(MAKE) clean-backend
	@$(MAKE) clean-frontend
	@$(MAKE) clean-docker
	@echo "$(COLOR_GREEN)✅ Cleanup completed!$(COLOR_RESET)"

clean-backend: ## Clean backend build artifacts
	@echo "$(COLOR_BLUE)Cleaning backend...$(COLOR_RESET)"
	@cd $(BACKEND_DIR) && rm -rf bin/ coverage.out coverage.html

clean-frontend: ## Clean frontend build artifacts
	@echo "$(COLOR_BLUE)Cleaning frontend...$(COLOR_RESET)"
	@cd $(FRONTEND_DIR) && rm -rf dist/ node_modules/.cache coverage/

clean-docker: ## Clean Docker resources
	@echo "$(COLOR_BLUE)Cleaning Docker resources...$(COLOR_RESET)"
	@docker system prune -f
	@docker volume prune -f

clean-all: ## Complete cleanup including Docker volumes
	@echo "$(COLOR_BOLD)🧹 Complete cleanup...$(COLOR_RESET)"
	@$(MAKE) down
	@$(MAKE) clean
	@docker-compose down -v
	@echo "$(COLOR_GREEN)✅ Complete cleanup done!$(COLOR_RESET)"

## 🔧 Utility Functions
ensure-logs-dir: ## Ensure logs directory exists
	@mkdir -p $(LOGS_DIR)

check-docker: ## Verify Docker is running
	@docker --version > /dev/null 2>&1 || (echo "$(COLOR_YELLOW)⚠️ Docker not found or not running$(COLOR_RESET)" && exit 1)

check-node: ## Verify Node.js is installed
	@node --version > /dev/null 2>&1 || (echo "$(COLOR_YELLOW)⚠️ Node.js not found$(COLOR_RESET)" && exit 1)

check-go: ## Verify Go is installed
	@go version > /dev/null 2>&1 || (echo "$(COLOR_YELLOW)⚠️ Go not found$(COLOR_RESET)" && exit 1)

## 🚀 Production Operations
prod-deploy: ## Deploy to production environment
	@echo "$(COLOR_BOLD)🚀 Deploying to production...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) -f docker-compose.prod.yml up -d
	@echo "$(COLOR_GREEN)✅ Production deployment completed!$(COLOR_RESET)"

prod-stop: ## Stop production environment
	@echo "$(COLOR_BOLD)🛑 Stopping production environment...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) -f docker-compose.prod.yml down
	@echo "$(COLOR_GREEN)✅ Production environment stopped!$(COLOR_RESET)"

## 🎯 Quick Actions for Development
quick-test: ## Quick test of core functionality
	@echo "$(COLOR_BOLD)⚡ Quick functionality test...$(COLOR_RESET)"
	@curl -s http://localhost:8080/health || echo "Backend not ready"
	@curl -s http://localhost:4200 > /dev/null && echo "Frontend OK" || echo "Frontend not ready"

setup: ## Initial project setup
	@echo "$(COLOR_BOLD)🎯 Setting up Cupcake for first time...$(COLOR_RESET)"
	@$(MAKE) check-docker
	@$(MAKE) check-node  
	@$(MAKE) check-go
	@$(MAKE) install-deps
	@$(MAKE) docker-build
	@echo "$(COLOR_GREEN)✅ Setup completed! Run 'make dev' to start.$(COLOR_RESET)"

## 📈 Development Shortcuts
lint: ## Run linters for all code
	@echo "$(COLOR_BOLD)🔍 Running linters...$(COLOR_RESET)"
	@cd $(BACKEND_DIR) && golint ./...
	@cd $(FRONTEND_DIR) && npm run lint

format: ## Format all code
	@echo "$(COLOR_BOLD)✨ Formatting code...$(COLOR_RESET)"
	@cd $(BACKEND_DIR) && go fmt ./...
	@cd $(FRONTEND_DIR) && npm run format

version: ## Show version information
	@echo "$(COLOR_BOLD)📋 Version Information:$(COLOR_RESET)"
	@echo "Project: $(PROJECT_NAME) v2.0"
	@echo "Go: $(shell go version 2>/dev/null || echo 'Not installed')"
	@echo "Node: $(shell node --version 2>/dev/null || echo 'Not installed')"
	@echo "Docker: $(shell docker --version 2>/dev/null || echo 'Not installed')"

.PHONY: help status dev dev-info dev-local backend frontend kafka database
.PHONY: install-deps install-backend-deps install-frontend-deps
.PHONY: docker-build docker-rebuild up stop down restart
.PHONY: db-init db-backup db-restore db-shell
.PHONY: test test-backend test-frontend test-integration test-coverage
.PHONY: build build-backend build-frontend
.PHONY: docs docs-serve
.PHONY: logs logs-backend logs-frontend logs-kafka logs-mongodb health
.PHONY: clean clean-backend clean-frontend clean-docker clean-all
.PHONY: ensure-logs-dir check-docker check-node check-go
.PHONY: prod-deploy prod-stop quick-test setup lint format version
