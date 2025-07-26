#!/bin/bash

echo "🧁 Cupcake Kafka API Testing Script"
echo "===================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

API_BASE="http://localhost:8080"

echo -e "${BLUE}1. Testing Health Check...${NC}"
curl -X GET "${API_BASE}/health" \
  -H "Accept: application/json" \
  -w "\n\nStatus: %{http_code}\nTime: %{time_total}s\n\n" \
  -s

echo -e "${BLUE}2. Testing Kafka Publish with inline JSON...${NC}"
curl -X POST "${API_BASE}/api/kafka/publish" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d '{
    "broker": "localhost:9093",
    "topic": "test-topic",
    "key": "test-key-1",
    "value": "{\"message\": \"Hello from curl inline!\", \"timestamp\": \"2025-07-26\", \"user\": \"test-user\"}"
  }' \
  -w "\n\nStatus: %{http_code}\nTime: %{time_total}s\n\n" \
  -s

echo -e "${BLUE}3. Testing Kafka Publish with complex JSON...${NC}"
curl -X POST "${API_BASE}/api/kafka/publish" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d '{
    "broker": "localhost:9093",
    "topic": "user-events",
    "key": "user-123",
    "value": "{\"eventType\": \"login\", \"userId\": 123, \"timestamp\": \"2025-07-26T08:45:00Z\", \"metadata\": {\"ip\": \"192.168.1.1\", \"userAgent\": \"curl/7.68.0\", \"sessionId\": \"sess_abc123\"}}"
  }' \
  -w "\n\nStatus: %{http_code}\nTime: %{time_total}s\n\n" \
  -s

echo -e "${BLUE}4. Testing with missing required fields (should fail)...${NC}"
curl -X POST "${API_BASE}/api/kafka/publish" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d '{
    "broker": "localhost:9093",
    "key": "test-key"
  }' \
  -w "\n\nStatus: %{http_code}\nTime: %{time_total}s\n\n" \
  -s

echo -e "${BLUE}5. Testing with invalid JSON (should fail)...${NC}"
curl -X POST "${API_BASE}/api/kafka/publish" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d '{
    "broker": "localhost:9093",
    "topic": "test-topic",
    "value": "invalid json here
  }' \
  -w "\n\nStatus: %{http_code}\nTime: %{time_total}s\n\n" \
  -s

echo -e "${GREEN}Testing completed!${NC}"
echo ""
echo "📚 For interactive testing, visit:"
echo "   Swagger UI: ${API_BASE}/swagger/index.html"
echo ""
echo "💡 Example curl commands:"
echo ""
echo "# Simple message:"
echo "curl -X POST ${API_BASE}/api/kafka/publish \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{"
echo "    \"broker\": \"localhost:9093\","
echo "    \"topic\": \"my-topic\","
echo "    \"key\": \"my-key\","
echo "    \"value\": \"{\\\"data\\\": \\\"my message\\\"}\""
echo "  }'"
echo ""
echo "# Using a file:"
echo "curl -X POST ${API_BASE}/api/kafka/publish \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d @my-message.json"
