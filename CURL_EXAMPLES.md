# Cupcake Kafka API - curl Examples

## Health Check
```bash
curl -X GET http://localhost:8080/health
```

## Publish Simple Message
```bash
curl -X POST http://localhost:8080/api/kafka/publish \
  -H "Content-Type: application/json" \
  -d '{
    "broker": "localhost:9092",
    "topic": "test-topic",
    "key": "simple-key",
    "value": "{\"message\": \"Hello World\", \"timestamp\": \"2025-07-26\"}"
  }'
```

## Publish Complex JSON Message
```bash
curl -X POST http://localhost:8080/api/kafka/publish \
  -H "Content-Type: application/json" \
  -d '{
    "broker": "localhost:9092",
    "topic": "user-events",
    "key": "user-123",
    "value": "{\"eventType\": \"login\", \"userId\": 123, \"timestamp\": \"2025-07-26T09:00:00Z\", \"metadata\": {\"ip\": \"192.168.1.1\", \"userAgent\": \"PostmanRuntime/7.29.0\", \"sessionId\": \"sess_abc123\", \"country\": \"US\"}}"
  }'
```

## Using JSON File
Create a file `message.json`:
```json
{
  "broker": "localhost:9092",
  "topic": "orders",
  "key": "order-456",
  "value": "{\"orderId\": 456, \"customerId\": 789, \"items\": [{\"productId\": 101, \"quantity\": 2, \"price\": 29.99}, {\"productId\": 102, \"quantity\": 1, \"price\": 15.50}], \"total\": 75.48, \"timestamp\": \"2025-07-26T09:15:00Z\"}"
}
```

Then run:
```bash
curl -X POST http://localhost:8080/api/kafka/publish \
  -H "Content-Type: application/json" \
  -d @message.json
```

## Expected Responses

### Successful Response:
```json
{
  "success": true,
  "message": "Message published successfully",
  "data": {
    "topic": "test-topic",
    "key": "simple-key"
  }
}
```

### Error Response (missing required field):
```json
{
  "error": "Validation Error",
  "message": "Topic is required"
}
```

### Error Response (invalid JSON):
```json
{
  "error": "Invalid JSON",
  "message": "invalid character 'i' looking for beginning of object key string"
}
```

## Postman Collection

You can also import this as a Postman collection:

### 1. Health Check
- Method: GET
- URL: `http://localhost:8080/health`

### 2. Publish Message
- Method: POST
- URL: `http://localhost:8080/api/kafka/publish`
- Headers: `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "broker": "localhost:9092",
  "topic": "test-topic",
  "key": "postman-key",
  "value": "{\"source\": \"postman\", \"message\": \"Hello from Postman!\", \"timestamp\": \"2025-07-26T09:30:00Z\"}"
}
```

## Notes
- The `broker` field should point to your Kafka broker (default: localhost:9092)
- The `topic` field is required and will be created automatically if it doesn't exist (depending on Kafka configuration)
- The `key` field is optional but recommended for message ordering
- The `value` field should contain your JSON message as a string (escaped JSON)
- For testing without a real Kafka broker, the API will return an error but the validation will still work
