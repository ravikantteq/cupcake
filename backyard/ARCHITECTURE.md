# Cupcake Kafka Test Framework - Idiomatic Go Architecture

## Overview

The project has been restructured to follow idiomatic Go patterns and best practices. The new architecture is clean, simple, and follows Go conventions while maintaining all existing functionality.

## Project Structure

```
backyard/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── types.go                # Core domain types and models
│   ├── store/                  # Data persistence layer
│   │   ├── store.go           # Store interfaces
│   │   └── mongodb.go         # MongoDB implementation
│   ├── manager/                # Business logic layer
│   │   ├── manager.go         # Manager aggregator
│   │   ├── consumer.go        # Consumer management
│   │   ├── flow.go            # Flow management
│   │   └── producer.go        # Producer management
│   └── handler/                # HTTP API layer
│       ├── handler.go         # Handler aggregator
│       ├── consumer.go        # Consumer HTTP handlers
│       ├── flow.go            # Flow HTTP handlers
│       └── producer.go        # Producer HTTP handlers
└── pkg/                        # Shared packages
    ├── netw/                   # Kafka networking
    └── storage/                # MongoDB client
```

## Architecture Principles

### 1. **Layered Architecture**
- **Handler Layer**: HTTP request/response handling
- **Manager Layer**: Business logic and orchestration
- **Store Layer**: Data persistence abstraction
- **Types**: Core domain models

### 2. **Go Conventions**
- Using "manager" instead of "service" (Go idiom)
- Using "store" instead of "repository" (Go idiom)
- Clear separation of concerns
- Dependency injection through constructors

### 3. **Clean Design**
- No circular dependencies
- Interfaces for testability
- Simple and straightforward structure
- Minimal abstractions

## Key Components

### Types (`internal/types.go`)
Central definition of all domain models:
- `Consumer` - Kafka consumer configuration
- `Message` - Kafka message representation
- `Flow` - Test flow definition
- `Execution` - Flow execution tracking
- Request/Response types for API

### Store Layer (`internal/store/`)
Data persistence abstraction:
- `Store` interface defines all data operations
- `MongoDB` struct implements the interface
- Clean separation between business logic and data access

### Manager Layer (`internal/manager/`)
Business logic and orchestration:
- `ConsumerManager` - Manages Kafka consumers and their lifecycle
- `FlowManager` - Handles test flow creation and execution
- `ProducerManager` - Manages message production and history
- Each manager handles specific domain logic

### Handler Layer (`internal/handler/`)
HTTP API endpoints:
- `ConsumerHandler` - Consumer-related endpoints
- `FlowHandler` - Flow-related endpoints
- `ProducerHandler` - Producer and health endpoints
- Clean request/response handling with proper error responses

## API Structure

### New V1 API Routes
```
GET    /health                      # System health check
POST   /api/kafka/publish           # Legacy compatibility

# Consumer Management
POST   /api/v1/consumers            # Create consumer
GET    /api/v1/consumers            # List consumers
GET    /api/v1/consumers/:id        # Get consumer
POST   /api/v1/consumers/:id/start  # Start consumer
POST   /api/v1/consumers/:id/stop   # Stop consumer
DELETE /api/v1/consumers/:id        # Delete consumer

# Flow Management
POST   /api/v1/flows                # Create flow
GET    /api/v1/flows                # List flows
GET    /api/v1/flows/:id            # Get flow
PUT    /api/v1/flows/:id            # Update flow
POST   /api/v1/flows/:id/execute    # Execute flow
DELETE /api/v1/flows/:id            # Delete flow

# Producer History
GET    /api/v1/history              # Get producer history
GET    /api/v1/history/recent       # Get recent history
```

## Improvements Made

### 1. **Idiomatic Go Structure**
- Removed complex abstractions
- Used Go naming conventions (manager vs service)
- Clear package organization
- Proper dependency injection

### 2. **Simplified Architecture**
- Removed unnecessary layers
- Clean separation of concerns
- Single responsibility principle
- Easy to understand and maintain

### 3. **Better Error Handling**
- Consistent error responses
- Proper HTTP status codes
- Centralized error types

### 4. **Improved Testability**
- Interface-based design
- Dependency injection
- Clear separation of layers
- Mockable components

### 5. **Enhanced Maintainability**
- Single source of truth for types
- Consistent patterns across components
- Clear naming conventions
- Minimal coupling

## Running the Application

```bash
# Build the application
go build -o bin/cupcake ./cmd

# Run the application
./bin/cupcake

# Environment variables
export MONGO_URI="mongodb://cupcake:cupcake123@localhost:27017/cupcake?authSource=admin"
export KAFKA_BROKER="localhost:9093"
export PORT="8080"
export GIN_MODE="release"  # for production
```

## Benefits of the New Structure

1. **Go Idiomatic**: Follows Go community standards and conventions
2. **Simple**: Easy to understand and navigate
3. **Testable**: Clean interfaces and dependency injection
4. **Maintainable**: Clear separation of concerns
5. **Scalable**: Easy to add new features and components
6. **Performance**: Minimal abstractions and overhead

## Migration Notes

- All existing functionality is preserved
- Legacy API endpoints remain for backward compatibility
- New V1 API provides cleaner interface
- Database schema unchanged
- Kafka integration unchanged

This restructuring provides a solid foundation for future development while maintaining all existing capabilities with improved code quality and maintainability.
