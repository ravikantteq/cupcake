package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Status represents the status of various entities
type Status string

const (
	StatusInactive  Status = "inactive"
	StatusActive    Status = "active"
	StatusError     Status = "error"
	StatusRunning   Status = "running"
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed" // Changed from "success" to match MongoDB schema
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// StepType defines different types of test steps
type StepType string

const (
	StepTypeProduce  StepType = "produce"
	StepTypeConsume  StepType = "consume"
	StepTypeValidate StepType = "validate"
	StepTypeDelay    StepType = "delay"
)

// Consumer represents a Kafka consumer
type Consumer struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name"`
	Description   string             `json:"description,omitempty" bson:"description,omitempty"`
	Broker        string             `json:"broker" bson:"broker"`
	GroupID       string             `json:"groupId" bson:"groupId"`
	Topics        []string           `json:"topics" bson:"topics"`
	Status        Status             `json:"status" bson:"status"`
	Config        ConsumerConfig     `json:"config" bson:"config"`
	MessageCount  int64              `json:"messageCount" bson:"messageCount"`
	LastHeartbeat *time.Time         `json:"lastHeartbeat,omitempty" bson:"lastHeartbeat,omitempty"`
	StartedAt     *time.Time         `json:"startedAt,omitempty" bson:"startedAt,omitempty"`
	StoppedAt     *time.Time         `json:"stoppedAt,omitempty" bson:"stoppedAt,omitempty"`
	ErrorMessage  string             `json:"errorMessage,omitempty" bson:"errorMessage,omitempty"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// ConsumerConfig holds Kafka consumer configuration
type ConsumerConfig struct {
	AutoOffsetReset     string `json:"autoOffsetReset" bson:"autoOffsetReset"`
	EnableAutoCommit    bool   `json:"enableAutoCommit" bson:"enableAutoCommit"`
	MaxPollRecords      int    `json:"maxPollRecords" bson:"maxPollRecords"`
	SessionTimeoutMs    int    `json:"sessionTimeoutMs,omitempty" bson:"sessionTimeoutMs,omitempty"`
	HeartbeatIntervalMs int    `json:"heartbeatIntervalMs,omitempty" bson:"heartbeatIntervalMs,omitempty"`
}

// Message represents a Kafka message
type Message struct {
	ID              primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Topic           string              `json:"topic" bson:"topic"`
	Partition       int32               `json:"partition" bson:"partition"`
	Offset          int64               `json:"offset" bson:"offset"`
	Key             string              `json:"key,omitempty" bson:"key,omitempty"`
	Value           interface{}         `json:"value" bson:"value"`
	Headers         map[string]string   `json:"headers,omitempty" bson:"headers,omitempty"`
	Timestamp       time.Time           `json:"timestamp" bson:"timestamp"`
	ConsumerGroupID string              `json:"consumerGroupId" bson:"consumerGroupId"`
	ConsumerID      primitive.ObjectID  `json:"consumerId,omitempty" bson:"consumerId,omitempty"`
	ExecutionID     *primitive.ObjectID `json:"executionId,omitempty" bson:"executionId,omitempty"`
}

// Flow represents a test flow
type Flow struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Version     string             `json:"version" bson:"version"`
	Steps       []FlowStep         `json:"steps" bson:"steps"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	CreatedBy   string             `json:"createdBy" bson:"createdBy"`
}

// FlowStep represents a single step in a test flow
type FlowStep struct {
	StepID string     `json:"stepId" bson:"stepId"`
	Type   StepType   `json:"type" bson:"type"`
	Config StepConfig `json:"config" bson:"config"`
}

// StepConfig holds configuration for different step types
type StepConfig struct {
	Topic           string                 `json:"topic,omitempty" bson:"topic,omitempty"`
	Message         map[string]interface{} `json:"message,omitempty" bson:"message,omitempty"`
	ExpectedMessage map[string]interface{} `json:"expectedMessage,omitempty" bson:"expectedMessage,omitempty"`
	Timeout         int                    `json:"timeout,omitempty" bson:"timeout,omitempty"`
	Retries         int                    `json:"retries,omitempty" bson:"retries,omitempty"`
	ExpectedCount   int                    `json:"expectedCount,omitempty" bson:"expectedCount,omitempty"`
	DelayMs         int                    `json:"delayMs,omitempty" bson:"delayMs,omitempty"`
}

// Execution represents the execution of a test flow
type Execution struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SuiteID   primitive.ObjectID `json:"suiteId" bson:"suiteId"`
	FlowID    primitive.ObjectID `json:"flowId" bson:"flowId"`
	Status    Status             `json:"status" bson:"status"`
	StartTime time.Time          `json:"startTime" bson:"startTime"`
	EndTime   *time.Time         `json:"endTime,omitempty" bson:"endTime,omitempty"`
	Steps     []ExecutionStep    `json:"steps" bson:"steps"`
	Metrics   ExecutionMetrics   `json:"metrics" bson:"metrics"`
	Logs      []string           `json:"logs,omitempty" bson:"logs,omitempty"`
}

// ExecutionStep represents the execution state of a step
type ExecutionStep struct {
	StepID   string                 `json:"stepId" bson:"stepId"`
	Status   Status                 `json:"status" bson:"status"`
	Input    map[string]interface{} `json:"input,omitempty" bson:"input,omitempty"`
	Output   map[string]interface{} `json:"output,omitempty" bson:"output,omitempty"`
	Errors   []string               `json:"errors,omitempty" bson:"errors,omitempty"`
	Duration int64                  `json:"duration" bson:"duration"` // milliseconds
}

// ExecutionMetrics represents metrics for a test execution
type ExecutionMetrics struct {
	TotalDuration     int64 `json:"totalDuration" bson:"totalDuration"` // milliseconds
	MessagesProduced  int   `json:"messagesProduced" bson:"messagesProduced"`
	MessagesConsumed  int   `json:"messagesConsumed" bson:"messagesConsumed"`
	ErrorsCount       int   `json:"errorsCount" bson:"errorsCount"`
	StepsCompleted    int   `json:"stepsCompleted" bson:"stepsCompleted"`
	ValidationsPassed int   `json:"validationsPassed" bson:"validationsPassed"`
	ValidationsFailed int   `json:"validationsFailed" bson:"validationsFailed"`
}

// ProducerHistory represents a record of producer message attempts
type ProducerHistory struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Broker    string             `json:"broker" bson:"broker"`
	Topic     string             `json:"topic" bson:"topic"`
	Key       string             `json:"key,omitempty" bson:"key,omitempty"`
	Value     string             `json:"value" bson:"value"`
	Success   bool               `json:"success" bson:"success"`
	Response  interface{}        `json:"response,omitempty" bson:"response,omitempty"`
	Error     string             `json:"error,omitempty" bson:"error,omitempty"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	UserID    string             `json:"userId,omitempty" bson:"userId,omitempty"`
}

// Request/Response types for API
type (
	// CreateConsumerRequest represents a request to create a consumer
	CreateConsumerRequest struct {
		Name        string         `json:"name" binding:"required"`
		Description string         `json:"description,omitempty"`
		Broker      string         `json:"broker" binding:"required"`
		GroupID     string         `json:"groupId" binding:"required"`
		Topics      []string       `json:"topics" binding:"required"`
		Config      ConsumerConfig `json:"config"`
	}

	// CreateFlowRequest represents a request to create a flow
	CreateFlowRequest struct {
		Name        string     `json:"name" binding:"required"`
		Description string     `json:"description"`
		Version     string     `json:"version"`
		Steps       []FlowStep `json:"steps" binding:"required"`
	}

	// KafkaMessageRequest represents a message to be published to Kafka
	KafkaMessageRequest struct {
		Broker  string            `json:"broker" binding:"required"`
		Topic   string            `json:"topic" binding:"required"`
		Key     string            `json:"key"`
		Value   interface{}       `json:"value" binding:"required"`
		Headers map[string]string `json:"headers,omitempty"`
	}

	// APIResponse represents a successful API response
	APIResponse struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	// APIError represents an error API response
	APIError struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}

	// HealthResponse represents health check response
	HealthResponse struct {
		Status    string            `json:"status"`
		Timestamp time.Time         `json:"timestamp"`
		Services  map[string]string `json:"services"`
		Version   string            `json:"version"`
	}
)
