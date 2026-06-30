package manager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ravikantteq/cupcake/backyard/internal"
	"github.com/ravikantteq/cupcake/backyard/internal/store"
	"github.com/ravikantteq/cupcake/backyard/pkg/netw"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConsumerRuntime represents a running consumer instance
type ConsumerRuntime struct {
	Consumer   *netw.KafkaConsumer
	Context    context.Context
	CancelFunc context.CancelFunc
	StartedAt  time.Time
}

// ConsumerManager manages Kafka consumers
type ConsumerManager struct {
	store           store.Store
	activeConsumers map[primitive.ObjectID]*ConsumerRuntime
	mu              sync.RWMutex
	kafkaBroker     string
}

// NewConsumerManager creates a new consumer manager
func NewConsumerManager(store store.Store, kafkaBroker string) *ConsumerManager {
	return &ConsumerManager{
		store:           store,
		activeConsumers: make(map[primitive.ObjectID]*ConsumerRuntime),
		kafkaBroker:     kafkaBroker,
	}
}

// CreateConsumer creates a new consumer
func (cm *ConsumerManager) CreateConsumer(ctx context.Context, req *internal.CreateConsumerRequest) (*internal.Consumer, error) {
	// Set default config values
	if req.Config.AutoOffsetReset == "" {
		req.Config.AutoOffsetReset = "latest"
	}
	if req.Config.SessionTimeoutMs == 0 {
		req.Config.SessionTimeoutMs = 10000
	}
	if req.Config.HeartbeatIntervalMs == 0 {
		req.Config.HeartbeatIntervalMs = 3000
	}
	if req.Config.MaxPollRecords == 0 {
		req.Config.MaxPollRecords = 100
	}
	if !req.Config.EnableAutoCommit {
		req.Config.EnableAutoCommit = true
	}

	consumer := &internal.Consumer{
		Name:        req.Name,
		Description: req.Description,
		Broker:      req.Broker,
		GroupID:     req.GroupID,
		Topics:      req.Topics,
		Config:      req.Config,
	}

	return cm.store.CreateConsumer(ctx, consumer)
}

// GetConsumer retrieves a consumer by ID
func (cm *ConsumerManager) GetConsumer(ctx context.Context, id primitive.ObjectID) (*internal.Consumer, error) {
	consumer, err := cm.store.GetConsumer(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update status if it's running
	cm.mu.RLock()
	runtime, isRunning := cm.activeConsumers[id]
	cm.mu.RUnlock()

	if isRunning && runtime.Consumer != nil && runtime.Consumer.IsRunning() {
		consumer.Status = internal.StatusActive
		now := time.Now()
		consumer.LastHeartbeat = &now
	}

	return consumer, nil
}

// GetConsumers retrieves all consumers
func (cm *ConsumerManager) GetConsumers(ctx context.Context) ([]*internal.Consumer, error) {
	consumers, err := cm.store.GetConsumers(ctx)
	if err != nil {
		return nil, err
	}

	// Update status for running consumers
	cm.mu.RLock()
	for _, consumer := range consumers {
		if runtime, isRunning := cm.activeConsumers[consumer.ID]; isRunning && runtime.Consumer != nil && runtime.Consumer.IsRunning() {
			consumer.Status = internal.StatusActive
			now := time.Now()
			consumer.LastHeartbeat = &now
		}
	}
	cm.mu.RUnlock()

	return consumers, nil
}

// StartConsumer starts a consumer
func (cm *ConsumerManager) StartConsumer(ctx context.Context, consumerID primitive.ObjectID) error {
	consumer, err := cm.store.GetConsumer(ctx, consumerID)
	if err != nil {
		return fmt.Errorf("consumer not found: %w", err)
	}

	// Check if already running
	cm.mu.RLock()
	if _, exists := cm.activeConsumers[consumerID]; exists {
		cm.mu.RUnlock()
		return fmt.Errorf("consumer is already running")
	}
	cm.mu.RUnlock()

	// Create Kafka consumer
	kafkaConsumer, err := netw.NewKafkaConsumer(
		consumerID,
		consumer.Broker,
		consumer.GroupID,
		consumer.Topics,
		consumer.Config,
	)
	if err != nil {
		return fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	// Create runtime context
	consumerCtx, cancelFunc := context.WithCancel(context.Background())
	runtime := &ConsumerRuntime{
		Consumer:   kafkaConsumer,
		Context:    consumerCtx,
		CancelFunc: cancelFunc,
		StartedAt:  time.Now(),
	}

	// Set message and error handlers
	kafkaConsumer.SetMessageHandler(func(msg *kafka.Message) {
		cm.handleMessage(consumerID, consumer.GroupID, msg)
	})

	kafkaConsumer.SetErrorHandler(func(err error) {
		cm.handleError(consumerID, err)
	})

	// Start the consumer in a goroutine
	go cm.runConsumerAsync(runtime, consumer)

	// Store runtime
	cm.mu.Lock()
	cm.activeConsumers[consumerID] = runtime
	cm.mu.Unlock()

	// Update status in database
	err = cm.store.UpdateConsumerStatus(ctx, consumerID, internal.StatusActive, "")
	if err != nil {
		fmt.Printf("Warning: Failed to update consumer status: %v\n", err)
	}

	return nil
}

// StopConsumer stops a consumer
func (cm *ConsumerManager) StopConsumer(ctx context.Context, consumerID primitive.ObjectID) error {
	cm.mu.RLock()
	runtime, exists := cm.activeConsumers[consumerID]
	cm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("consumer is not running")
	}

	// Cancel the context to stop the goroutine
	runtime.CancelFunc()

	// Stop the Kafka consumer
	if runtime.Consumer != nil {
		err := runtime.Consumer.Stop()
		if err != nil {
			return fmt.Errorf("failed to stop consumer: %w", err)
		}
	}

	// Remove from active consumers
	cm.mu.Lock()
	delete(cm.activeConsumers, consumerID)
	cm.mu.Unlock()

	// Update status in database
	err := cm.store.UpdateConsumerStatus(ctx, consumerID, internal.StatusInactive, "")
	if err != nil {
		fmt.Printf("Warning: Failed to update consumer status: %v\n", err)
	}

	return nil
}

// DeleteConsumer deletes a consumer
func (cm *ConsumerManager) DeleteConsumer(ctx context.Context, consumerID primitive.ObjectID) error {
	// Check if running
	cm.mu.RLock()
	_, isRunning := cm.activeConsumers[consumerID]
	cm.mu.RUnlock()

	if isRunning {
		return fmt.Errorf("cannot delete running consumer - stop it first")
	}

	return cm.store.DeleteConsumer(ctx, consumerID)
}

// runConsumerAsync runs the consumer in a goroutine
func (cm *ConsumerManager) runConsumerAsync(runtime *ConsumerRuntime, consumer *internal.Consumer) {
	defer func() {
		// Clean up
		cm.mu.Lock()
		delete(cm.activeConsumers, consumer.ID)
		cm.mu.Unlock()

		// Update database status
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cm.store.UpdateConsumerStatus(ctx, consumer.ID, internal.StatusInactive, "")

		fmt.Printf("Consumer %s goroutine finished\n", consumer.ID.Hex())
	}()

	fmt.Printf("Starting consumer %s for topics: %v\n", consumer.GroupID, consumer.Topics)

	// Start the consumer
	err := runtime.Consumer.Start(runtime.Context)
	if err != nil {
		fmt.Printf("Failed to start Kafka consumer %s: %v\n", consumer.ID.Hex(), err)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cm.store.UpdateConsumerStatus(ctx, consumer.ID, internal.StatusError, err.Error())
		return
	}

	// Keep the goroutine alive until context is cancelled
	<-runtime.Context.Done()
	fmt.Printf("Consumer %s context cancelled, stopping...\n", consumer.ID.Hex())
}

// handleMessage processes consumed messages
func (cm *ConsumerManager) handleMessage(consumerID primitive.ObjectID, groupID string, msg *kafka.Message) {
	message := &internal.Message{
		Topic:           *msg.TopicPartition.Topic,
		Partition:       msg.TopicPartition.Partition,
		Offset:          int64(msg.TopicPartition.Offset),
		Key:             string(msg.Key),
		Value:           string(msg.Value),
		Headers:         make(map[string]string),
		Timestamp:       time.Now(),
		ConsumerGroupID: groupID,
		ConsumerID:      consumerID,
	}

	// Convert headers
	for _, header := range msg.Headers {
		message.Headers[header.Key] = string(header.Value)
	}

	// Store message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := cm.store.StoreMessage(ctx, message)
	if err != nil {
		fmt.Printf("Failed to store consumed message: %v\n", err)
	}

	// Update message count
	err = cm.store.IncrementMessageCount(ctx, consumerID)
	if err != nil {
		fmt.Printf("Failed to update message count: %v\n", err)
	}
}

// handleError processes consumer errors
func (cm *ConsumerManager) handleError(consumerID primitive.ObjectID, err error) {
	fmt.Printf("Consumer %s encountered error: %v\n", consumerID.Hex(), err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateErr := cm.store.UpdateConsumerStatus(ctx, consumerID, internal.StatusError, err.Error())
	if updateErr != nil {
		fmt.Printf("Failed to update consumer error status: %v\n", updateErr)
	}
}
