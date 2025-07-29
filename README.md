# 🧁 Cupcake Kafka Test Framework

A comprehensive, enterprise-ready Kafka testing platform with advanced flow design, intelligent message matching, and real-time monitoring capabilities.

## 🎯 Features

- **🔄 Test Flow Designer**: Visual interface for creating complex Kafka message flows
- **🎭 Intelligent Message Matching**: Advanced assertion framework with dynamic value matching
- **📊 Real-time Monitoring**: Live dashboard for test execution and system health
- **🗃️ Persistent Storage**: MongoDB integration for test configurations and execution history
- **🐳 Fully Dockerized**: Complete containerized environment with one-command setup
- **🎪 Test Suite Management**: Organize and execute multiple test flows as suites
- **🕵️ Consumer Management**: Auto-setup consumers with flexible topic subscription
- **📈 Advanced Reporting**: Comprehensive execution reports and analytics
- **🔐 Enterprise Ready**: Built for scale with security and performance in mind

## 🏗️ Architecture

```
cupcake/
├── 📄 Cupcake_Kafka_Test_Framework_PRD_v2.md  # Comprehensive PRD
├── 🐳 docker-compose.yml                      # Complete infrastructure setup
├── 📁 backyard/                               # Go backend service
│   ├── cmd/                                   # Entry points
│   ├── internal/                              # Core business logic
│   │   ├── api/                              # HTTP handlers & routes
│   │   ├── models/                           # Data models
│   │   ├── services/                         # Business services
│   │   └── consumers/                        # Kafka consumer management
│   ├── pkg/                                  # Shared packages
│   │   ├── netw/                            # Kafka client
│   │   ├── validation/                       # Message validation
│   │   └── storage/                          # MongoDB operations
│   └── docs/                                 # API documentation
├── 📁 cupcake_ui/                            # Angular frontend
│   ├── src/app/                              
│   │   ├── components/                       # UI components
│   │   ├── services/                         # API services
│   │   ├── models/                           # TypeScript models
│   │   └── pages/                            # Application pages
│   └── public/                               # Static assets
└── 📁 scripts/                               # Database & deployment scripts
```

## 🚀 Quick Start

### Prerequisites

- **Docker** 24.0+ and **Docker Compose** 2.0+
- **Git** for version control
- 8GB+ RAM recommended for full stack

### Installation & Startup

1. **Clone and navigate:**
   ```bash
   git clone <repository-url>
   cd cupcake
   ```

2. **Start the complete stack:**
   ```bash
   docker-compose up -d
   ```
   
   This will start:
   - 🧩 **Kafka & Zookeeper** (message infrastructure)
   - 🗄️ **MongoDB** (data persistence)
   - 🔧 **Backend Service** (Go API server)
   - 🎨 **Frontend** (Angular UI)
   - 📊 **Kafka UI** (broker monitoring)
   - 🗃️ **Mongo Express** (database management)

3. **Access the application:**
   - **Main UI**: http://localhost:4200
   - **API Docs**: http://localhost:8080/swagger/index.html
   - **Kafka UI**: http://localhost:8081
   - **Database UI**: http://localhost:8082 (admin/admin123)

## 🎪 Key Capabilities

### 1. Test Flow Creation
Create sophisticated test flows with multiple steps:

```json
{
  "name": "Order Processing Flow",
  "steps": [
    {
      "type": "produce",
      "config": {
        "topic": "orders-input",
        "message": {
          "orderId": "uuid()",
          "amount": "number(min=10, max=1000)",
          "timestamp": "timestamp()"
        }
      }
    },
    {
      "type": "consume",
      "config": {
        "topic": "orders-processed",
        "timeout": 10000
      }
    },
    {
      "type": "validate",
      "config": {
        "expectedMessage": {
          "orderId": "match(step-1.orderId)",
          "status": "enum(processed,completed)",
          "processedAt": "timestamp()"
        }
      }
    }
  ]
}
```

### 2. Dynamic Message Matching
Support for intelligent assertions:
- `uuid()` - Match any valid UUID
- `timestamp()` - Match any valid timestamp
- `number(min=X, max=Y)` - Match numbers within range
- `enum(val1,val2,val3)` - Match specific values
- `regex(pattern)` - Match regex patterns
- `any()` - Match any value
- `match(step-X.field)` - Reference previous step values

### 3. Consumer Management
- **Auto-Consumer Setup**: Automatically create consumers for test topics
- **Group Management**: Organize consumers into logical groups
- **Offset Control**: Manual and automatic offset management
- **Health Monitoring**: Real-time consumer health tracking

### 4. Test Suite Execution
- **Batch Execution**: Run multiple flows as a coordinated suite
- **Environment Configs**: Different settings for dev/staging/prod
- **Parallel Processing**: Execute independent flows simultaneously
- **Detailed Reporting**: Comprehensive execution reports

## 📚 API Reference

### Core Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | System health check |
| `POST` | `/api/v1/flows` | Create test flow |
| `GET` | `/api/v1/flows` | List all flows |
| `POST` | `/api/v1/suites` | Create test suite |
| `POST` | `/api/v1/suites/{id}/execute` | Execute suite |
| `POST` | `/api/v1/consumers` | Create consumer |
| `GET` | `/api/v1/consumers` | List active consumers |

### Example: Create and Execute Flow

```bash
# Create a test flow
curl -X POST http://localhost:8080/api/v1/flows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Payment Processing Test",
    "description": "Test payment flow end-to-end",
    "steps": [
      {
        "stepId": "produce-payment",
        "type": "produce",
        "config": {
          "topic": "payments-input",
          "message": {
            "paymentId": "uuid()",
            "amount": 100.50,
            "currency": "USD"
          }
        }
      },
      {
        "stepId": "consume-result",
        "type": "consume",
        "config": {
          "topic": "payments-processed",
          "timeout": 5000
        }
      }
    ]
  }'

# Execute the flow as part of a suite
curl -X POST http://localhost:8080/api/v1/suites/execute \
  -H "Content-Type: application/json" \
  -d '{
    "suiteId": "suite-id-here",
    "environment": "development"
  }'
```

## 🧪 Testing Framework Features

### Message Validation Patterns
```javascript
// Dynamic timestamp validation
"timestamp": "timestamp()"

// Range-based number validation  
"amount": "number(min=0, max=10000)"

// Enum validation
"status": "enum(pending,completed,failed)"

// UUID validation
"id": "uuid()"

// Reference previous step values
"correlationId": "match(step-1.paymentId)"

// Custom regex patterns
"code": "regex(^[A-Z]{3}\\d{3}$)"
```

### Consumer Configuration
```json
{
  "groupId": "test-consumer-group",
  "topics": ["orders", "payments", "notifications"],
  "config": {
    "autoOffsetReset": "earliest",
    "enableAutoCommit": true,
    "maxPollRecords": 100
  }
}
```

## 🔧 Development Workflow

### Available Commands

```bash
# Infrastructure management
docker-compose up -d          # Start all services
docker-compose down           # Stop all services
docker-compose logs -f        # Follow logs

# Development shortcuts
make dev                      # Start development environment
make test                     # Run all tests
make docs                     # Generate API documentation
make clean                    # Clean up resources

# Database operations
make db-init                  # Initialize database
make db-backup                # Backup database
make db-restore               # Restore database
```

### Service Health Checks

Monitor service health:
```bash
# Backend health
curl http://localhost:8080/health

# Check all services
docker-compose ps
```

## 📊 Monitoring & Observability

### Built-in Dashboards
- **Test Execution Dashboard**: Real-time view of running tests
- **System Health**: Monitor Kafka, MongoDB, and service status
- **Message Flow Tracking**: End-to-end message traceability
- **Performance Metrics**: Latency, throughput, and error rates

### External Tools Integration
- **Kafka UI**: Topic management and message browsing
- **Mongo Express**: Database inspection and management
- **Health Endpoints**: Integration with monitoring systems

## 🔒 Security Features

- **Input Validation**: Comprehensive request validation
- **Error Handling**: Secure error responses
- **Resource Limits**: Prevent resource exhaustion
- **Audit Logging**: Complete action audit trail

## 🚀 Production Deployment

### Environment Variables
```bash
# Backend configuration
KAFKA_BROKER=your-kafka-broker:9093
MONGO_URI=mongodb://username:password@your-mongo:27017/cupcake
LOG_LEVEL=info
GIN_MODE=release

# Frontend configuration  
API_BASE_URL=https://your-api-domain.com
```

### Docker Production Setup
```bash
# Use production docker-compose
docker-compose -f docker-compose.prod.yml up -d

# Scale services
docker-compose up -d --scale cupcake-backend=3
```

## 🛠️ Troubleshooting

### Common Issues

**Services won't start:**
```bash
# Check Docker resources
docker system df
docker system prune  # Clean up if needed

# Check logs
docker-compose logs cupcake-backend
docker-compose logs cupcake-frontend
```

**Kafka connection issues:**
```bash
# Verify Kafka is running
docker-compose exec kafka kafka-topics --bootstrap-server localhost:9093 --list

# Check network connectivity
docker-compose exec cupcake-backend ping kafka
```

**Database connection problems:**
```bash
# Test MongoDB connection
docker-compose exec mongodb mongosh cupcake

# Check initialization
docker-compose logs mongodb | grep "initialization"
```

### Performance Tuning

For high-throughput scenarios:
```yaml
# Kafka configuration
environment:
  KAFKA_NUM_NETWORK_THREADS: 8
  KAFKA_NUM_IO_THREADS: 16
  KAFKA_SOCKET_SEND_BUFFER_BYTES: 102400
  KAFKA_SOCKET_RECEIVE_BUFFER_BYTES: 102400
```

## 🤝 Contributing

1. **Fork** the repository
2. **Create** feature branch: `git checkout -b feature/amazing-feature`
3. **Follow** the coding standards and add tests
4. **Test** your changes: `make test`
5. **Commit** changes: `git commit -am 'Add amazing feature'`
6. **Push** to branch: `git push origin feature/amazing-feature`
7. **Create** Pull Request

### Development Guidelines
- Follow Go best practices and Angular style guide
- Write comprehensive tests for new features
- Update documentation for API changes
- Use conventional commit messages

## 📄 Documentation

- **📋 [Complete PRD](./Cupcake_Kafka_Test_Framework_PRD_v2.md)**: Comprehensive product requirements
- **🔧 [API Documentation](http://localhost:8080/swagger/index.html)**: Interactive API docs
- **🏗️ [Architecture Guide](./docs/architecture.md)**: System design details
- **🚀 [Deployment Guide](./docs/deployment.md)**: Production deployment

## 🔮 Roadmap

### Current Version (v2.0)
- ✅ Visual test flow designer
- ✅ Intelligent message matching
- ✅ MongoDB integration
- ✅ Docker containerization

### Upcoming Features
- 🔄 **CI/CD Integration**: GitHub Actions, Jenkins plugins
- 📊 **Advanced Analytics**: ML-based insights and predictions
- 🔌 **Plugin Architecture**: Custom validators and integrations
- ☁️ **Cloud Support**: AWS, GCP, Azure deployment templates

---

## 📞 Support

- **Issues**: [GitHub Issues](../../issues)
- **Discussions**: [GitHub Discussions](../../discussions)
- **Documentation**: [Wiki](../../wiki)

Built with ❤️ for the Kafka testing community

---

**License**: MIT | **Version**: 2.0.0 | **Maintainer**: Ravikant P
