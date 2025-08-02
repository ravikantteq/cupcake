package services

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ravikantteq/cupcake/backyard/internal/models"
	"github.com/ravikantteq/cupcake/backyard/internal/repository"
	"github.com/ravikantteq/cupcake/backyard/pkg/netw"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConsumerRuntime represents a running consumer with its management context
type ConsumerRuntime struct {
	Consumer   *netw.KafkaConsumer
	Context    context.Context
	CancelFunc context.CancelFunc
	StartedAt  time.Time
	ConsumerID primitive.ObjectID
}

// ConsumerService manages Kafka consumers with goroutine-based lifecycle management
type ConsumerService struct {
	consumerRepo    *repository.ConsumerRepository
	messageRepo     *repository.MessageRepository
	activeConsumers map[primitive.ObjectID]*ConsumerRuntime // Changed to hold runtime info
	mu              sync.RWMutex
	kafkaBroker     string
}

// NewConsumerService creates a new consumer service
func NewConsumerService(repo *repository.Repository, kafkaBroker string) *ConsumerService {
	return &ConsumerService{
		consumerRepo:    repo.NewConsumerRepository(),
		messageRepo:     repo.NewMessageRepository(),
		activeConsumers: make(map[primitive.ObjectID]*ConsumerRuntime),
		kafkaBroker:     kafkaBroker,
	}
}

// CreateConsumer creates a new consumer configuration
func (cs *ConsumerService) CreateConsumer(ctx context.Context, req *models.CreateConsumerRequest) (*models.Consumer, error) {
	consumer := &models.Consumer{
		Name:        req.Name,
		Description: req.Description,
		Broker:      req.Broker,
		GroupID:     req.GroupID,
		Topics:      req.Topics,
		Status:      models.ConsumerInactive,
		Config:      req.Config,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set default config values if not provided
	if consumer.Config.AutoOffsetReset == "" {
		consumer.Config.AutoOffsetReset = "latest" // Only consume new messages
	}
	if consumer.Config.SessionTimeoutMs == 0 {
		consumer.Config.SessionTimeoutMs = 10000
	}
	if consumer.Config.HeartbeatIntervalMs == 0 {
		consumer.Config.HeartbeatIntervalMs = 3000
	}
	if consumer.Config.MaxPollRecords == 0 {
		consumer.Config.MaxPollRecords = 100
	}

	return cs.consumerRepo.CreateConsumer(ctx, consumer)
}

// GetConsumerByID retrieves a consumer by ID
func (cs *ConsumerService) GetConsumerByID(ctx context.Context, id primitive.ObjectID) (*models.Consumer, error) {
	return cs.consumerRepo.GetConsumerByID(ctx, id)
}

// GetAllConsumers retrieves all consumers
func (cs *ConsumerService) GetAllConsumers(ctx context.Context) ([]models.Consumer, error) {
	return cs.consumerRepo.GetAllConsumers(ctx)
}

// StartConsumer starts a consumer in a separate goroutine for scalability
func (cs *ConsumerService) StartConsumer(ctx context.Context, consumerID primitive.ObjectID) error {
	// Get consumer details from database
	consumer, err := cs.consumerRepo.GetConsumerByID(ctx, consumerID)
	if err != nil {
		return fmt.Errorf("consumer not found: %w", err)
	}

	// Check if consumer is already running
	cs.mu.RLock()
	if _, exists := cs.activeConsumers[consumerID]; exists {
		cs.mu.RUnlock()
		return fmt.Errorf("consumer is already running")
	}
	cs.mu.RUnlock()

	// Update status to starting
	err = cs.consumerRepo.UpdateConsumerStatus(ctx, consumer.ID, models.ConsumerActive, "")
	if err != nil {
		return fmt.Errorf("failed to update consumer status: %w", err)
	}

	// Create consumer runtime with independent context
	consumerCtx, cancelFunc := context.WithCancel(context.Background())
	runtime := &ConsumerRuntime{
		Context:    consumerCtx,
		CancelFunc: cancelFunc,
		StartedAt:  time.Now(),
		ConsumerID: consumerID,
	}

	// Start consumer in a goroutine for async execution
	go cs.runConsumerAsync(runtime, consumer)

	// Store runtime in active consumers (before async start to prevent race conditions)
	cs.mu.Lock()
	cs.activeConsumers[consumerID] = runtime
	cs.mu.Unlock()

	return nil
}

// runConsumerAsync runs a consumer asynchronously in a goroutine
func (cs *ConsumerService) runConsumerAsync(runtime *ConsumerRuntime, consumer *models.Consumer) {
	defer func() {
		// Clean up: remove from active consumers when done
		cs.mu.Lock()
		delete(cs.activeConsumers, runtime.ConsumerID)
		cs.mu.Unlock()

		// Update database status
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		now := time.Now()
		cs.consumerRepo.UpdateConsumerWithDetails(ctx, runtime.ConsumerID, models.ConsumerInactive, "", nil, &now)

		fmt.Printf("Consumer %s goroutine finished\n", runtime.ConsumerID.Hex())
	}()

	fmt.Printf("Starting consumer %s in goroutine for topics: %v\n", consumer.GroupID, consumer.Topics)

	// Create Kafka consumer
	kafkaConsumer, err := netw.NewKafkaConsumer(
		runtime.ConsumerID,
		consumer.Broker,
		consumer.GroupID,
		consumer.Topics,
		consumer.Config,
	)
	if err != nil {
		// Update status to error
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cs.consumerRepo.UpdateConsumerStatus(ctx, runtime.ConsumerID, models.ConsumerError, err.Error())
		fmt.Printf("Failed to create Kafka consumer for %s: %v\n", runtime.ConsumerID.Hex(), err)
		return
	}

	// Store the kafka consumer in runtime
	runtime.Consumer = kafkaConsumer

	// Set message handler
	kafkaConsumer.SetMessageHandler(func(msg *kafka.Message) {
		cs.handleMessage(runtime.ConsumerID, msg)
	})

	// Set error handler
	kafkaConsumer.SetErrorHandler(func(err error) {
		cs.handleError(runtime.ConsumerID, err)
	})

	// Start the consumer with the runtime context
	err = kafkaConsumer.Start(runtime.Context)
	if err != nil {
		// Update status to error
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cs.consumerRepo.UpdateConsumerStatus(ctx, runtime.ConsumerID, models.ConsumerError, err.Error())
		fmt.Printf("Failed to start Kafka consumer for %s: %v\n", runtime.ConsumerID.Hex(), err)
		return
	}

	// Update status to active with started timestamp
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	now := time.Now()
	err = cs.consumerRepo.UpdateConsumerWithDetails(ctx, runtime.ConsumerID, models.ConsumerActive, "", &now, nil)
	if err != nil {
		fmt.Printf("Warning: Failed to update consumer status in DB: %v\n", err)
	}

	fmt.Printf("Consumer %s started successfully in goroutine\n", runtime.ConsumerID.Hex())

	// Keep the goroutine alive until context is cancelled
	<-runtime.Context.Done()
	fmt.Printf("Consumer %s context cancelled, stopping...\n", runtime.ConsumerID.Hex())
}

// StopConsumer stops a consumer by cancelling its goroutine
func (cs *ConsumerService) StopConsumer(ctx context.Context, consumerID primitive.ObjectID) error {
	// Get consumer runtime from active map
	cs.mu.RLock()
	runtime, exists := cs.activeConsumers[consumerID]
	cs.mu.RUnlock()

	if !exists {
		return fmt.Errorf("consumer is not running")
	}

	// Update status to stopping
	err := cs.consumerRepo.UpdateConsumerStatus(ctx, consumerID, models.ConsumerInactive, "")
	if err != nil {
		return fmt.Errorf("failed to update consumer status: %w", err)
	}

	// Cancel the context to stop the goroutine
	runtime.CancelFunc()

	// Stop the Kafka consumer if it exists
	if runtime.Consumer != nil {
		err = runtime.Consumer.Stop()
		if err != nil {
			// Update status to error
			cs.consumerRepo.UpdateConsumerStatus(ctx, consumerID, models.ConsumerError, err.Error())
			return fmt.Errorf("failed to stop consumer: %w", err)
		}
	}

	// Remove from active consumers map (will also be cleaned up by goroutine defer)
	cs.mu.Lock()
	delete(cs.activeConsumers, consumerID)
	cs.mu.Unlock()

	// Update status to stopped
	now := time.Now()
	err = cs.consumerRepo.UpdateConsumerWithDetails(ctx, consumerID, models.ConsumerInactive, "", nil, &now)
	if err != nil {
		// Consumer is stopped but we couldn't update DB - log error but don't fail
		fmt.Printf("Warning: Failed to update consumer status in DB: %v\n", err)
	}

	fmt.Printf("Consumer %s stopped successfully\n", consumerID.Hex())
	return nil
}

// GetConsumerStatus returns the status of a consumer
func (cs *ConsumerService) GetConsumerStatus(ctx context.Context, consumerID primitive.ObjectID) (*models.Consumer, error) {
	consumer, err := cs.consumerRepo.GetConsumerByID(ctx, consumerID)
	if err != nil {
		return nil, err
	}

	// Check if it's in active consumers map
	cs.mu.RLock()
	runtime, isActive := cs.activeConsumers[consumerID]
	cs.mu.RUnlock()

	if isActive && runtime.Consumer != nil && runtime.Consumer.IsRunning() {
		// Update heartbeat for active consumer
		now := time.Now()
		consumer.LastHeartbeat = &now
		consumer.Status = models.ConsumerActive
	}

	return consumer, nil
}

// DeleteConsumer deletes a consumer (only if not running)
func (cs *ConsumerService) DeleteConsumer(ctx context.Context, consumerID primitive.ObjectID) error {
	// Check if consumer is running
	cs.mu.RLock()
	_, isRunning := cs.activeConsumers[consumerID]
	cs.mu.RUnlock()

	if isRunning {
		return fmt.Errorf("cannot delete running consumer - stop it first")
	}

	return cs.consumerRepo.DeleteConsumer(ctx, consumerID)
}

// handleMessage processes consumed messages
func (cs *ConsumerService) handleMessage(consumerID primitive.ObjectID, msg *kafka.Message) {
	// Get consumer details to access the real group ID
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	consumer, err := cs.consumerRepo.GetConsumerByID(ctx, consumerID)
	if err != nil {
		fmt.Printf("Failed to get consumer details for message handling: %v\n", err)
		return
	}

	// Convert Kafka message to our Message model
	message := &models.Message{
		Topic:           *msg.TopicPartition.Topic,
		Partition:       msg.TopicPartition.Partition,
		Offset:          int64(msg.TopicPartition.Offset),
		Key:             string(msg.Key),
		Value:           string(msg.Value),
		Headers:         make(map[string]string),
		Timestamp:       time.Now(),
		ConsumerGroupID: consumer.GroupID, // Use the actual Kafka consumer group ID
		ConsumerID:      consumerID,       // Track which consumer instance processed this
	}

	// Convert headers
	for _, header := range msg.Headers {
		message.Headers[header.Key] = string(header.Value)
	}

	// Store message in database
	storeCtx, storeCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer storeCancel()

	storeErr := cs.messageRepo.StoreMessage(storeCtx, message)
	if storeErr != nil {
		fmt.Printf("Failed to store consumed message: %v\n", storeErr)
	}

	// Update message count for consumer
	countErr := cs.consumerRepo.IncrementMessageCount(storeCtx, consumerID)
	if countErr != nil {
		fmt.Printf("Failed to update message count: %v\n", countErr)
	}
}

// handleError processes consumer errors
func (cs *ConsumerService) handleError(consumerID primitive.ObjectID, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("Consumer %s encountered error: %v\n", consumerID.Hex(), err)

	// Check if this is a fatal connection error
	errStr := err.Error()
	isFatalError := strings.Contains(errStr, "max connection retries reached") ||
		strings.Contains(errStr, "Connection refused") ||
		strings.Contains(errStr, "All broker connections are down") ||
		strings.Contains(errStr, "connect to ipv4") ||
		strings.Contains(errStr, "failed: Connection refused")

	if isFatalError {
		fmt.Printf("Fatal connection error detected for consumer %s, stopping completely\n", consumerID.Hex())

		// Stop and remove the consumer from active map
		cs.mu.Lock()
		if runtime, exists := cs.activeConsumers[consumerID]; exists {
			// Cancel the context to stop the goroutine
			runtime.CancelFunc()
			// Stop the Kafka consumer if it exists
			if runtime.Consumer != nil {
				runtime.Consumer.Stop()
			}
			delete(cs.activeConsumers, consumerID)
		}
		cs.mu.Unlock()

		// Update status to error with stopped timestamp
		now := time.Now()
		updateErr := cs.consumerRepo.UpdateConsumerWithDetails(ctx, consumerID, models.ConsumerError, err.Error(), nil, &now)
		if updateErr != nil {
			fmt.Printf("Failed to update consumer error status: %v\n", updateErr)
		}
	} else {
		// For non-fatal errors, just update the status but keep the consumer running
		updateErr := cs.consumerRepo.UpdateConsumerStatus(ctx, consumerID, models.ConsumerError, err.Error())
		if updateErr != nil {
			fmt.Printf("Failed to update consumer error status: %v\n", updateErr)
		}
	}
}
