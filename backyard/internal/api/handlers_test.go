package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/systemgenes/cupcake/backyard/internal/models"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	kafkaHandler := NewKafkaHandler("localhost:9092")

	r.GET("/health", kafkaHandler.HealthCheck)

	apiGroup := r.Group("/api")
	{
		kafkaGroup := apiGroup.Group("/kafka")
		{
			kafkaGroup.POST("/publish", kafkaHandler.PublishMessage)
		}
	}

	return r
}

func TestHealthCheck(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
}

func TestPublishMessage_InvalidJSON(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/kafka/publish", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPublishMessage_MissingTopic(t *testing.T) {
	router := setupTestRouter()

	message := models.KafkaMessage{
		Broker: "localhost:9092",
		Key:    "test-key",
		Value:  "test-value",
		// Topic is missing
	}

	jsonData, _ := json.Marshal(message)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/kafka/publish", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPublishMessage_MissingValue(t *testing.T) {
	router := setupTestRouter()

	message := models.KafkaMessage{
		Broker: "localhost:9092",
		Topic:  "test-topic",
		Key:    "test-key",
		// Value is missing
	}

	jsonData, _ := json.Marshal(message)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/kafka/publish", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// Note: This test requires a running Kafka instance
func TestPublishMessage_ValidMessage(t *testing.T) {
	// Skip this test if no Kafka broker is available
	t.Skip("Skipping integration test - requires running Kafka broker")

	router := setupTestRouter()

	message := models.KafkaMessage{
		Broker: "localhost:9092",
		Topic:  "test-topic",
		Key:    "test-key",
		Value:  `{"data": "test message"}`,
	}

	jsonData, _ := json.Marshal(message)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/kafka/publish", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
}
