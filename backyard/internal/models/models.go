package models

// KafkaMessage represents a message to be published to Kafka
type KafkaMessage struct {
	Broker string `json:"broker" binding:"required" example:"localhost:9093"`
	Topic  string `json:"topic" binding:"required" example:"test-topic"`
	Key    string `json:"key" example:"message-key"`
	Value  string `json:"value" binding:"required" example:"{'data': 'test message'}"`
}

// Response represents a successful API response
type Response struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Error   string `json:"error" example:"Validation Error"`
	Message string `json:"message" example:"Field is required"`
}
