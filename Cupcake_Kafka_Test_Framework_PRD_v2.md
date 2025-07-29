# 🧁 Cupcake Kafka Test Framework - Product Requirements Document (PRD) v2.0

## 📋 Executive Summary

**Project Name:** Cupcake Kafka Test Framework  
**Version:** 2.0  
**Date:** July 26, 2025  
**Type:** Enterprise Kafka Testing Platform  

Cupcake is a comprehensive, dockerized Kafka testing framework designed for enterprise-level applications. It provides end-to-end testing capabilities for Kafka message flows with intelligent consumer setup, flexible test suite management, and advanced message matching capabilities.

## 🎯 Project Vision & Goals

### Primary Goal
Transform Cupcake into a production-ready Kafka testing platform that enables teams to:
- Design, execute, and validate complex Kafka message flows
- Create reusable test suites for regression testing
- Verify message transformations and processing logic
- Monitor and debug Kafka-based microservices

### Key Success Metrics
- Reduce Kafka testing setup time by 80%
- Enable non-technical stakeholders to create test flows
- Support testing of 100+ concurrent message flows
- Provide comprehensive audit trail and reporting

## 🏗️ System Architecture

### Technology Stack
- **Backend:** Go 1.20+ (Gin framework)
- **Frontend:** Angular 19+ (Standalone components)
- **Database:** MongoDB (message storage, test configurations)
- **Message Broker:** Apache Kafka
- **Containerization:** Docker & Docker Compose
- **Documentation:** Swagger/OpenAPI 3.0

### Architecture Principles
1. **Microservices Design:** Modular, loosely coupled components
2. **Event-Driven:** Reactive message processing
3. **Cloud-Native:** 12-factor app compliance
4. **API-First:** RESTful APIs with comprehensive documentation
5. **Test-Driven:** Comprehensive test coverage

## 🎪 Core Features

### 1. Dynamic Consumer Management
- **Auto-Consumer Setup:** Automatically create consumers for subscribed topics
- **Multi-Topic Support:** Single consumer can listen to multiple topics
- **Consumer Groups:** Configurable consumer group management
- **Offset Management:** Manual and automatic offset handling
- **Dead Letter Queue:** Handle failed message processing

### 2. Test Flow Designer
- **Visual Flow Builder:** Drag-and-drop interface for creating test flows
- **Step Configuration:** Input/output message configuration per step
- **Branching Logic:** Support for conditional flow paths
- **Parallel Execution:** Execute multiple flows simultaneously
- **Flow Templates:** Pre-built templates for common patterns

### 3. Intelligent Message Matching
- **Flexible Assertions:** JSON path-based message validation
- **Dynamic Values:** Support for `any()`, `regex()`, `timestamp()`, `uuid()` matchers
- **Schema Validation:** JSON Schema and Avro schema validation
- **Partial Matching:** Match subset of message fields
- **Custom Validators:** Extensible validation framework

### 4. Test Suite Management
- **Suite Organization:** Hierarchical organization of test flows
- **Versioning:** Version control for test suites
- **Environment Configs:** Different configurations per environment
- **Scheduling:** Automated test execution
- **Regression Testing:** Continuous validation of message flows

### 5. Real-time Monitoring
- **Live Dashboard:** Real-time view of test execution
- **Message Tracking:** End-to-end message traceability
- **Performance Metrics:** Latency, throughput, and error rates
- **Alerting:** Configurable alerts for test failures
- **Audit Logs:** Comprehensive logging and audit trail

## 🔧 Technical Specifications

### API Design

#### Core Endpoints
```
# Test Flow Management
POST   /api/v1/flows                    # Create test flow
GET    /api/v1/flows                    # List all flows
GET    /api/v1/flows/{id}               # Get flow details
PUT    /api/v1/flows/{id}               # Update flow
DELETE /api/v1/flows/{id}               # Delete flow

# Test Suite Management
POST   /api/v1/suites                   # Create test suite
GET    /api/v1/suites                   # List all suites
GET    /api/v1/suites/{id}              # Get suite details
POST   /api/v1/suites/{id}/execute      # Execute suite

# Consumer Management
POST   /api/v1/consumers                # Create consumer
GET    /api/v1/consumers                # List consumers
PUT    /api/v1/consumers/{id}/topics    # Subscribe to topics
DELETE /api/v1/consumers/{id}           # Stop consumer

# Message Operations
POST   /api/v1/messages/produce         # Produce message
GET    /api/v1/messages/consume/{topic} # Consume messages
POST   /api/v1/messages/validate        # Validate message
```

#### Message Validation Framework
```json
{
  "messagePattern": {
    "orderId": "uuid()",
    "timestamp": "timestamp()",
    "amount": "number(min=0, max=10000)",
    "status": "enum(pending,completed,failed)",
    "metadata": {
      "source": "string(regex='^[A-Z]+$')",
      "tags": "array(string, minLength=1)"
    }
  }
}
```

### Database Schema (MongoDB)

#### Collections Design
```javascript
// flows collection
{
  _id: ObjectId,
  name: "string",
  description: "string",
  version: "string",
  steps: [
    {
      stepId: "string",
      type: "produce|consume|validate",
      config: {
        topic: "string",
        message: "object",
        expectedResponse: "object",
        timeout: "number",
        retries: "number"
      }
    }
  ],
  createdAt: Date,
  updatedAt: Date,
  createdBy: "string"
}

// suites collection
{
  _id: ObjectId,
  name: "string",
  description: "string",
  flows: ["ObjectId"],
  environment: "string",
  config: {
    kafkaBroker: "string",
    consumerGroups: ["string"],
    timeouts: "object"
  },
  createdAt: Date,
  updatedAt: Date
}

// executions collection
{
  _id: ObjectId,
  suiteId: ObjectId,
  flowId: ObjectId,
  status: "running|completed|failed",
  startTime: Date,
  endTime: Date,
  steps: [
    {
      stepId: "string",
      status: "string",
      input: "object",
      output: "object",
      errors: ["string"],
      duration: "number"
    }
  ],
  metrics: {
    totalDuration: "number",
    messagesProduced: "number",
    messagesConsumed: "number",
    errorsCount: "number"
  }
}

// consumers collection
{
  _id: ObjectId,
  groupId: "string",
  topics: ["string"],
  status: "active|inactive",
  config: {
    autoOffsetReset: "string",
    enableAutoCommit: "boolean",
    maxPollRecords: "number"
  },
  lastHeartbeat: Date,
  createdAt: Date
}
```

### Docker Configuration

#### Enhanced docker-compose.yml
```yaml
version: '3.8'

services:
  # Kafka Infrastructure
  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    volumes:
      - zk-data:/var/lib/zookeeper/data
      - zk-logs:/var/lib/zookeeper/log

  kafka:
    image: confluentinc/cp-kafka:latest
    depends_on: [zookeeper]
    ports: ["9093:9093"]
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9093
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: 'true'
      KAFKA_LOG_RETENTION_HOURS: 24
    volumes:
      - kafka-data:/var/lib/kafka/data

  # Database
  mongodb:
    image: mongo:7.0
    ports: ["27017:27017"]
    environment:
      MONGO_INITDB_ROOT_USERNAME: cupcake
      MONGO_INITDB_ROOT_PASSWORD: cupcake123
      MONGO_INITDB_DATABASE: cupcake
    volumes:
      - mongo-data:/data/db
      - ./scripts/mongo-init.js:/docker-entrypoint-initdb.d/mongo-init.js

  # Backend Service
  cupcake-backend:
    build: 
      context: ./backyard
      dockerfile: Dockerfile
    ports: ["8080:8080"]
    depends_on: [kafka, mongodb]
    environment:
      KAFKA_BROKER: kafka:9093
      MONGO_URI: mongodb://cupcake:cupcake123@mongodb:27017/cupcake
      LOG_LEVEL: info
    volumes:
      - ./logs:/app/logs

  # Frontend Service
  cupcake-frontend:
    build:
      context: ./cupcake_ui
      dockerfile: Dockerfile
    ports: ["4200:80"]
    depends_on: [cupcake-backend]
    environment:
      API_BASE_URL: http://cupcake-backend:8080

  # Monitoring & Management
  kafka-ui:
    image: provectuslabs/kafka-ui:latest
    ports: ["8081:8080"]
    environment:
      KAFKA_CLUSTERS_0_NAME: local
      KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:9093
    depends_on: [kafka]

volumes:
  mongo-data:
  kafka-data:
  zk-data:
  zk-logs:
```

## 🎨 User Interface Design

### Main Dashboard
- **Test Suites Overview:** Grid view of all test suites with status indicators
- **Recent Executions:** Timeline of recent test runs
- **System Health:** Real-time status of Kafka, consumers, and backend services
- **Quick Actions:** Shortcuts to create flows, execute suites, view reports

### Test Flow Designer
- **Canvas Area:** Visual flow designer with drag-and-drop nodes
- **Toolbox:** Library of step types (produce, consume, validate, delay)
- **Properties Panel:** Configuration panel for selected steps
- **Preview Mode:** Test individual steps before full execution

### Message Designer
- **Schema Builder:** Visual JSON schema designer
- **Template Library:** Pre-built message templates for common patterns
- **Validation Rules:** Configure complex validation patterns
- **Mock Data Generator:** Generate test data based on schemas

### Execution Monitor
- **Live Progress:** Real-time execution progress with step-by-step status
- **Message Inspector:** View actual vs expected messages
- **Error Details:** Detailed error information with suggestions
- **Performance Metrics:** Charts showing latency and throughput

## 🔄 User Workflows

### Workflow 1: Creating a New Test Flow
1. **Navigate** to Test Flows section
2. **Click** "Create New Flow"
3. **Configure** flow metadata (name, description, tags)
4. **Design** flow using visual designer:
   - Add "Produce Message" step
   - Configure message template and topic
   - Add "Consume Response" step
   - Configure expected response pattern
   - Add validation rules
5. **Test** individual steps
6. **Save** flow for future use

### Workflow 2: Setting Up Consumer Subscription
1. **Navigate** to Consumer Management
2. **Click** "Create Consumer"
3. **Configure** consumer group and topics
4. **Set** consumption parameters (offset, polling)
5. **Start** consumer
6. **Monitor** consumer health and message processing

### Workflow 3: Executing Test Suite
1. **Select** test suite from dashboard
2. **Choose** environment configuration
3. **Review** flow execution order
4. **Click** "Execute Suite"
5. **Monitor** real-time progress
6. **Review** results and generate report

### Workflow 4: Debugging Failed Test
1. **Identify** failed test from dashboard
2. **Open** execution details
3. **Analyze** error messages and logs
4. **Compare** expected vs actual messages
5. **Edit** test flow to fix issues
6. **Re-execute** test to validate fix

## 📊 Advanced Features

### Message Assertion Library
```javascript
// Built-in assertion functions
const assertions = {
  any: () => true,
  uuid: (value) => /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(value),
  timestamp: (value) => !isNaN(Date.parse(value)),
  regex: (pattern) => (value) => new RegExp(pattern).test(value),
  number: (options) => (value) => {
    const num = Number(value);
    return !isNaN(num) && 
           (options.min === undefined || num >= options.min) &&
           (options.max === undefined || num <= options.max);
  },
  enum: (...values) => (value) => values.includes(value),
  arrayOf: (validator) => (value) => Array.isArray(value) && value.every(validator)
};
```

### Environment Management
- **Multi-Environment Support:** Dev, staging, production configurations
- **Secret Management:** Secure handling of credentials and API keys
- **Configuration Templating:** Parameterized configurations
- **Environment Promotion:** Promote test suites across environments

### Reporting & Analytics
- **Execution Reports:** Detailed reports with pass/fail statistics
- **Trend Analysis:** Historical analysis of test execution trends
- **Performance Dashboards:** System performance and bottleneck identification
- **Export Capabilities:** PDF, Excel, and JSON export formats

## 🚀 Implementation Roadmap

### Phase 1: Core Infrastructure (Weeks 1-4)
- [ ] Enhance Docker setup with MongoDB
- [ ] Implement consumer management system
- [ ] Create basic test flow data models
- [ ] Set up API authentication and authorization

### Phase 2: Test Flow Engine (Weeks 5-8)
- [ ] Build test flow execution engine
- [ ] Implement message assertion framework
- [ ] Create visual flow designer (basic)
- [ ] Add message pattern matching

### Phase 3: Advanced Features (Weeks 9-12)
- [ ] Implement test suite management
- [ ] Add real-time monitoring dashboard
- [ ] Create comprehensive reporting system
- [ ] Build environment management

### Phase 4: Production Readiness (Weeks 13-16)
- [ ] Performance optimization
- [ ] Security hardening
- [ ] Comprehensive testing
- [ ] Documentation and training materials

## 🛡️ Security Considerations

### Authentication & Authorization
- **JWT-based Authentication:** Secure token-based authentication
- **Role-based Access Control:** Different permissions for different user roles
- **API Rate Limiting:** Prevent abuse and ensure fair usage
- **Audit Logging:** Track all user actions and system changes

### Data Security
- **Encryption at Rest:** MongoDB data encryption
- **Encryption in Transit:** HTTPS/TLS for all communications
- **Secret Management:** Secure storage of sensitive configuration
- **Data Sanitization:** Prevent injection attacks

## 📈 Performance Requirements

### Scalability Targets
- **Concurrent Users:** Support 50+ concurrent users
- **Message Throughput:** Handle 10,000+ messages per second
- **Test Executions:** Support 100+ parallel test executions
- **Data Storage:** Efficiently store millions of messages and execution results

### Performance Benchmarks
- **API Response Time:** <200ms for standard operations
- **Test Execution Time:** <30 seconds for typical flows
- **UI Load Time:** <3 seconds for initial page load
- **Real-time Updates:** <1 second latency for live monitoring

## 🎯 Success Criteria

### Technical Success Metrics
- [ ] 100% containerized deployment
- [ ] <1% test execution failure rate due to system issues
- [ ] 99.9% uptime for core services
- [ ] Comprehensive API test coverage (>90%)

### User Experience Metrics
- [ ] <10 minutes for new user onboarding
- [ ] <5 clicks to create and execute a basic test flow
- [ ] Positive user feedback score (>4.5/5)
- [ ] Reduction in manual testing effort (>70%)

### Business Impact Metrics
- [ ] Faster time-to-market for Kafka-based features
- [ ] Reduced production incidents related to message processing
- [ ] Improved team collaboration on Kafka testing
- [ ] Cost reduction in QA resources

## 📚 Documentation Strategy

### Technical Documentation
- **API Documentation:** Complete OpenAPI/Swagger specifications
- **Architecture Guide:** System design and component interactions
- **Deployment Guide:** Step-by-step deployment instructions
- **Troubleshooting Guide:** Common issues and solutions

### User Documentation
- **User Manual:** Comprehensive guide for all features
- **Quick Start Guide:** Get users productive quickly
- **Tutorial Videos:** Screen recordings for complex workflows
- **Best Practices:** Guidelines for effective test design

## 🔮 Future Enhancements

### Short-term (6 months)
- **Kafka Schema Registry Integration**
- **GraphQL API Support**
- **Mobile-responsive UI**
- **Integration with CI/CD pipelines**

### Medium-term (12 months)
- **Multi-cluster Kafka support**
- **Advanced analytics and ML-based insights**
- **Plugin architecture for custom validators**
- **Integration with external monitoring tools**

### Long-term (18+ months)
- **Multi-cloud deployment support**
- **Advanced workflow orchestration**
- **AI-powered test generation**
- **Enterprise SSO integration**

---

## 📄 Appendices

### A. Glossary
- **Test Flow:** A sequence of steps that define a complete message testing scenario
- **Test Suite:** A collection of related test flows that can be executed together
- **Message Pattern:** A template that defines the expected structure and values of a message
- **Consumer Group:** A group of consumers that work together to consume messages from topics

### B. Technical Dependencies
- Go 1.20+
- Angular 19+
- MongoDB 7.0+
- Apache Kafka 3.0+
- Docker 24.0+
- Docker Compose 2.0+

### C. Resource Requirements
- **Development:** 16GB RAM, 8 CPU cores, 100GB storage
- **Production:** 32GB RAM, 16 CPU cores, 500GB storage
- **High Availability:** Load balancer, redundant services, backup storage

---

**Document Status:** Draft v2.0  
**Next Review:** August 15, 2025  
**Stakeholders:** Engineering Team, QA Team, DevOps Team, Product Management
