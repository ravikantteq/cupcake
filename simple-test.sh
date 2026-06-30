#!/bin/bash

echo "Testing Cupcake API..."

echo "1. Health Check:"
curl -s http://localhost:8080/health && echo

echo -e "\n2. Simple Kafka Publish Test:"
curl -s -X POST http://localhost:8080/api/kafka/publish \
  -H "Content-Type: application/json" \
  -d '{
    "broker": "localhost:9093",
    "topic": "test-topic",
    "key": "test-key",
    "value": "{\"message\": \"test\"}"
  }' && echo

echo -e "\nDone!"
