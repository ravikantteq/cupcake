package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

// StepType defines the type of step in a test flow
type StepType string

const (
	StepTypeProduce  StepType = "produce"
	StepTypeConsume  StepType = "consume"
	StepTypeValidate StepType = "validate"
	StepTypeDelay    StepType = "delay"
)

// ExecutionStatus defines the status of test execution
type ExecutionStatus string

const (
	StatusRunning   ExecutionStatus = "running"
	StatusCompleted ExecutionStatus = "completed"
	StatusFailed    ExecutionStatus = "failed"
	StatusCancelled ExecutionStatus = "cancelled"
)

// ConsumerStatus defines the status of a consumer
type ConsumerStatus string

const (
	ConsumerActive   ConsumerStatus = "active"
	ConsumerInactive ConsumerStatus = "inactive"
	ConsumerError    ConsumerStatus = "error"
)

// StepConfig represents the configuration for a test step
type StepConfig struct {
	Topic           string                 `json:"topic,omitempty" bson:"topic,omitempty"`
	Message         map[string]interface{} `json:"message,omitempty" bson:"message,omitempty"`
	ExpectedMessage map[string]interface{} `json:"expectedMessage,omitempty" bson:"expectedMessage,omitempty"`
	Timeout         int                    `json:"timeout,omitempty" bson:"timeout,omitempty"`
	Retries         int                    `json:"retries,omitempty" bson:"retries,omitempty"`
	ExpectedCount   int                    `json:"expectedCount,omitempty" bson:"expectedCount,omitempty"`
	DelayMs         int                    `json:"delayMs,omitempty" bson:"delayMs,omitempty"`
}

// FlowStep represents a single step in a test flow
type FlowStep struct {
	StepID string     `json:"stepId" bson:"stepId" binding:"required"`
	Type   StepType   `json:"type" bson:"type" binding:"required"`
	Config StepConfig `json:"config" bson:"config"`
}

// TestFlow represents a complete test flow definition
type TestFlow struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" binding:"required"`
	Description string             `json:"description" bson:"description"`
	Version     string             `json:"version" bson:"version"`
	Steps       []FlowStep         `json:"steps" bson:"steps" binding:"required"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	CreatedBy   string             `json:"createdBy" bson:"createdBy"`
}

// TestSuite represents a collection of test flows
type TestSuite struct {
	ID          primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string               `json:"name" bson:"name" binding:"required"`
	Description string               `json:"description" bson:"description"`
	Flows       []primitive.ObjectID `json:"flows" bson:"flows"`
	Environment string               `json:"environment" bson:"environment"`
	Config      SuiteConfig          `json:"config" bson:"config"`
	CreatedAt   time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt" bson:"updatedAt"`
}

// SuiteConfig represents configuration for a test suite
type SuiteConfig struct {
	KafkaBroker     string            `json:"kafkaBroker" bson:"kafkaBroker"`
	ConsumerGroups  []string          `json:"consumerGroups" bson:"consumerGroups"`
	Timeouts        map[string]int    `json:"timeouts" bson:"timeouts"`
	CustomVariables map[string]string `json:"customVariables,omitempty" bson:"customVariables,omitempty"`
}

// ExecutionStep represents the execution state of a step
type ExecutionStep struct {
	StepID   string                 `json:"stepId" bson:"stepId"`
	Status   ExecutionStatus        `json:"status" bson:"status"`
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

// TestExecution represents the execution of a test flow or suite
type TestExecution struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SuiteID   primitive.ObjectID `json:"suiteId" bson:"suiteId"`
	FlowID    primitive.ObjectID `json:"flowId" bson:"flowId"`
	Status    ExecutionStatus    `json:"status" bson:"status"`
	StartTime time.Time          `json:"startTime" bson:"startTime"`
	EndTime   *time.Time         `json:"endTime,omitempty" bson:"endTime,omitempty"`
	Steps     []ExecutionStep    `json:"steps" bson:"steps"`
	Metrics   ExecutionMetrics   `json:"metrics" bson:"metrics"`
	Logs      []string           `json:"logs,omitempty" bson:"logs,omitempty"`
}

// ConsumerConfig represents configuration for a Kafka consumer
type ConsumerConfig struct {
	AutoOffsetReset     string `json:"autoOffsetReset" bson:"autoOffsetReset"`
	EnableAutoCommit    bool   `json:"enableAutoCommit" bson:"enableAutoCommit"`
	MaxPollRecords      int    `json:"maxPollRecords" bson:"maxPollRecords"`
	SessionTimeoutMs    int    `json:"sessionTimeoutMs,omitempty" bson:"sessionTimeoutMs,omitempty"`
	HeartbeatIntervalMs int    `json:"heartbeatIntervalMs,omitempty" bson:"heartbeatIntervalMs,omitempty"`
}

// Consumer represents a Kafka consumer instance
type Consumer struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	GroupID       string             `json:"groupId" bson:"groupId" binding:"required"`
	Topics        []string           `json:"topics" bson:"topics" binding:"required"`
	Status        ConsumerStatus     `json:"status" bson:"status"`
	Config        ConsumerConfig     `json:"config" bson:"config"`
	LastHeartbeat *time.Time         `json:"lastHeartbeat,omitempty" bson:"lastHeartbeat,omitempty"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	ErrorMessage  string             `json:"errorMessage,omitempty" bson:"errorMessage,omitempty"`
}

// Message represents a Kafka message stored in the database
type Message struct {
	ID              primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Topic           string              `json:"topic" bson:"topic"`
	Partition       int32               `json:"partition" bson:"partition"`
	Offset          int64               `json:"offset" bson:"offset"`
	Key             string              `json:"key,omitempty" bson:"key,omitempty"`
	Value           interface{}         `json:"value" bson:"value"`
	Headers         map[string]string   `json:"headers,omitempty" bson:"headers,omitempty"`
	Timestamp       time.Time           `json:"timestamp" bson:"timestamp"`
	ConsumerGroupID string              `json:"consumerGroupId,omitempty" bson:"consumerGroupId,omitempty"`
	ExecutionID     *primitive.ObjectID `json:"executionId,omitempty" bson:"executionId,omitempty"`
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
	UserID    string             `json:"userId,omitempty" bson:"userId,omitempty"` // For future user tracking
}

// CreateFlowRequest represents a request to create a new test flow
type CreateFlowRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	Version     string     `json:"version"`
	Steps       []FlowStep `json:"steps" binding:"required"`
}

// CreateSuiteRequest represents a request to create a new test suite
type CreateSuiteRequest struct {
	Name        string               `json:"name" binding:"required"`
	Description string               `json:"description"`
	Flows       []primitive.ObjectID `json:"flows"`
	Environment string               `json:"environment"`
	Config      SuiteConfig          `json:"config"`
}

// ExecuteSuiteRequest represents a request to execute a test suite
type ExecuteSuiteRequest struct {
	SuiteID     primitive.ObjectID `json:"suiteId" binding:"required"`
	Environment string             `json:"environment"`
	Variables   map[string]string  `json:"variables,omitempty"`
}

// CreateConsumerRequest represents a request to create a consumer
type CreateConsumerRequest struct {
	GroupID string         `json:"groupId" binding:"required"`
	Topics  []string       `json:"topics" binding:"required"`
	Config  ConsumerConfig `json:"config"`
}

// HealthStatus represents the health status of the system
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}
