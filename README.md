# 🧁 Cupcake | Backyard

A simple, well-designed Kafka producer testing tool with a clean Angular frontend and robust Go backend.

## 🎯 Features

- **Kafka Producer API**: Publish messages to any Kafka topic via REST API
- **Interactive Web UI**: Angular-based frontend for easy message publishing
- **Swagger Documentation**: Auto-generated API documentation
- **Comprehensive Testing**: Unit and integration tests for reliability
- **CORS Support**: Frontend and backend can run on different ports
- **Health Checks**: Monitor backend service health
- **Clean Architecture**: Modular, maintainable code structure

## 🏗️ Project Structure

```
cupcake/
├── backyard/                 # Go backend
│   ├── cmd/
│   │   └── main.go          # API server entry point
│   ├── internal/
│   │   ├── api/             # HTTP handlers
│   │   └── models/          # Data models
│   ├── pkg/
│   │   └── netw/            # Kafka producer
│   ├── docs/                # Swagger generated docs
│   └── go.mod
├── cupcake_ui/              # Angular frontend
│   ├── src/
│   │   └── app/
│   │       ├── components/  # UI components
│   │       └── services/    # HTTP services
│   └── package.json
└── Makefile                 # Build automation
```

## 🚀 Quick Start

### Prerequisites

- Go 1.20+
- Node.js 18+
- Angular CLI 19+
- Kafka (for actual message publishing)

### Installation

1. **Install dependencies:**
   ```bash
   make install-deps
   ```

2. **Start the backend:**
   ```bash
   make backend
   ```
   Backend runs on `http://localhost:8080`

3. **Start the frontend (in a new terminal):**
   ```bash
   make frontend
   ```
   Frontend runs on `http://localhost:4200`

### Using the Application

1. **Access the UI**: Open `http://localhost:4200`
2. **Fill in the form:**
   - **Broker**: Your Kafka broker address (e.g., `localhost:9093`)
   - **Topic**: Kafka topic name
   - **Key**: Message key (optional)
   - **Value**: JSON message content

3. **Publish**: Click "Publish Message" to send to Kafka

## 📚 API Documentation

Access the Swagger UI at: `http://localhost:8080/swagger/index.html`

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/kafka/publish` | Publish message to Kafka |

### Example API Usage

```bash
curl -X POST http://localhost:8080/api/kafka/publish \
  -H "Content-Type: application/json" \
  -d '{
    "broker": "localhost:9093",
    "topic": "test-topic",
    "key": "test-key",
    "value": "{\"data\": \"hello world\"}"
  }'
```

## 🧪 Testing

### Run All Tests
```bash
make test
```

### Backend Tests Only
```bash
make test-backend
```

### Frontend Tests Only
```bash
make test-frontend
```

### Test Coverage

- **Unit Tests**: Core logic validation
- **Integration Tests**: API endpoint testing
- **Validation Tests**: Input validation and error handling

*Note: Some tests are skipped by default if Kafka is not running. To run full integration tests, ensure Kafka is available.*

## 🔧 Development

### Available Make Commands

```bash
make help                 # Show all available commands
make backend             # Start backend server
make frontend            # Start frontend server
make dev                 # Start both backend and frontend
make build               # Build both components
make docs                # Generate Swagger docs
make clean               # Clean build artifacts
```

### Project Design Principles

1. **Clean Architecture**: Separation of concerns with clear layers
2. **Testability**: Comprehensive test coverage with mocking
3. **Maintainability**: Clear code structure and documentation
4. **Extensibility**: Easy to add new features (consumers, topics, etc.)
5. **Simplicity**: No overengineering, focused on core functionality

## 🔄 Future Enhancements

The architecture supports easy extension for:

- **Kafka Consumer**: Add consumer functionality
- **Topic Management**: Create/list/delete topics
- **Message History**: Store and view published messages
- **Batch Publishing**: Send multiple messages
- **Authentication**: Add security layers
- **Monitoring**: Message delivery tracking
- **Configuration Management**: Environment-based configs

## 🐳 Docker Support (Coming Soon)

The project is designed to support containerization:

```bash
make docker-backend     # Build backend image
make docker-frontend    # Build frontend image
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make your changes with tests
4. Run tests: `make test`
5. Commit changes: `git commit -am 'Add new feature'`
6. Push: `git push origin feature/new-feature`
7. Create a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Troubleshooting

### Common Issues

**Backend won't start:**
- Check if port 8080 is available
- Ensure Go dependencies are installed: `cd backyard && go mod tidy`

**Frontend won't start:**
- Check if port 4200 is available
- Install Angular dependencies: `cd cupcake_ui && npm install`

**Can't publish to Kafka:**
- Ensure Kafka is running on the specified broker address
- Check topic exists or Kafka is configured to auto-create topics
- Verify network connectivity to Kafka broker

**CORS errors:**
- Backend includes CORS headers for development
- For production, configure proper CORS settings

### Health Check

Always check backend health at: `http://localhost:8080/health`

---

Built with ❤️ using Go and Angular
