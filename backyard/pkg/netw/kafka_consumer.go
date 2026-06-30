package netw

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ravikantteq/cupcake/backyard/internal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// KafkaConsumer represents a Kafka consumer instance
type KafkaConsumer struct {
	ID        primitive.ObjectID
	consumer  *kafka.Consumer
	ctx       context.Context
	cancel    context.CancelFunc
	topics    []string
	groupID   string
	broker    string
	isRunning bool
	mu        sync.RWMutex
	onMessage func(*kafka.Message) // Callback for message handling
	onError   func(error)          // Callback for error handling
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(id primitive.ObjectID, broker, groupID string, topics []string, config internal.ConsumerConfig) (*KafkaConsumer, error) {
	// Configure consumer
	configMap := kafka.ConfigMap{
		"bootstrap.servers":  broker,
		"group.id":           groupID,
		"auto.offset.reset":  config.AutoOffsetReset,
		"enable.auto.commit": config.EnableAutoCommit,
	}

	// Add optional configurations
	if config.SessionTimeoutMs > 0 {
		configMap["session.timeout.ms"] = config.SessionTimeoutMs
	}
	if config.HeartbeatIntervalMs > 0 {
		configMap["heartbeat.interval.ms"] = config.HeartbeatIntervalMs
	}
	// Note: max.poll.records is not available in confluent-kafka-go, it's handled differently

	consumer, err := kafka.NewConsumer(&configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	kc := &KafkaConsumer{
		ID:        id,
		consumer:  consumer,
		topics:    topics,
		groupID:   groupID,
		broker:    broker,
		isRunning: false,
	}

	return kc, nil
}

// SetMessageHandler sets the callback for handling messages
func (kc *KafkaConsumer) SetMessageHandler(handler func(*kafka.Message)) {
	kc.mu.Lock()
	defer kc.mu.Unlock()
	kc.onMessage = handler
}

// SetErrorHandler sets the callback for handling errors
func (kc *KafkaConsumer) SetErrorHandler(handler func(error)) {
	kc.mu.Lock()
	defer kc.mu.Unlock()
	kc.onError = handler
}

// Start begins consuming messages
func (kc *KafkaConsumer) Start(ctx context.Context) error {
	kc.mu.Lock()
	defer kc.mu.Unlock()

	if kc.isRunning {
		return fmt.Errorf("consumer is already running")
	}

	// Subscribe to topics
	err := kc.consumer.SubscribeTopics(kc.topics, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	kc.ctx, kc.cancel = context.WithCancel(ctx)
	kc.isRunning = true

	// Start consuming in a goroutine
	go kc.consumeLoop()

	log.Printf("Consumer %s started for topics: %v", kc.groupID, kc.topics)
	return nil
}

// Stop stops the consumer
func (kc *KafkaConsumer) Stop() error {
	kc.mu.Lock()
	defer kc.mu.Unlock()

	if !kc.isRunning {
		return fmt.Errorf("consumer is not running")
	}

	kc.cancel()
	kc.isRunning = false

	// Close the consumer
	err := kc.consumer.Close()
	if err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}

	log.Printf("Consumer %s stopped", kc.groupID)
	return nil
}

// IsRunning returns whether the consumer is currently running
func (kc *KafkaConsumer) IsRunning() bool {
	kc.mu.RLock()
	defer kc.mu.RUnlock()
	return kc.isRunning
}

// GetTopics returns the topics this consumer is subscribed to
func (kc *KafkaConsumer) GetTopics() []string {
	return kc.topics
}

// GetGroupID returns the consumer group ID
func (kc *KafkaConsumer) GetGroupID() string {
	return kc.groupID
}

// consumeLoop is the main consumption loop
func (kc *KafkaConsumer) consumeLoop() {
	defer func() {
		kc.mu.Lock()
		kc.isRunning = false
		kc.mu.Unlock()
	}()

	connectionRetries := 0
	maxConnectionRetries := 5

	for {
		select {
		case <-kc.ctx.Done():
			log.Printf("Consumer %s: context cancelled, stopping consumption", kc.groupID)
			return
		default:
			// Poll for messages with timeout
			msg, err := kc.consumer.ReadMessage(1 * time.Second)
			if err != nil {
				// Check if it's a timeout (not a real error)
				if kafkaErr, ok := err.(kafka.Error); ok {
					if kafkaErr.Code() == kafka.ErrTimedOut {
						continue // Timeout is normal, just continue polling
					}

					// Check for connection errors
					if kafkaErr.Code() == kafka.ErrTransport ||
						kafkaErr.Code() == kafka.ErrAllBrokersDown ||
						kafkaErr.Code() == kafka.ErrNetworkException {
						connectionRetries++
						log.Printf("Consumer %s connection error (retry %d/%d): %v",
							kc.groupID, connectionRetries, maxConnectionRetries, err)

						if connectionRetries >= maxConnectionRetries {
							log.Printf("Consumer %s: max connection retries reached, stopping", kc.groupID)

							// Call error handler and stop
							kc.mu.RLock()
							errorHandler := kc.onError
							kc.mu.RUnlock()

							if errorHandler != nil {
								errorHandler(fmt.Errorf("max connection retries reached: %w", err))
							}
							return
						}

						// Wait before retrying
						time.Sleep(2 * time.Second)
						continue
					}
				}

				log.Printf("Consumer %s error: %v", kc.groupID, err)

				// Call error handler if set
				kc.mu.RLock()
				errorHandler := kc.onError
				kc.mu.RUnlock()

				if errorHandler != nil {
					errorHandler(err)
				}
				continue
			}

			// Reset connection retry counter on successful message
			connectionRetries = 0

			// Process the message
			log.Printf("Consumer %s received message on topic %s: %s",
				kc.groupID, *msg.TopicPartition.Topic, string(msg.Value))

			// Call message handler if set
			kc.mu.RLock()
			messageHandler := kc.onMessage
			kc.mu.RUnlock()

			if messageHandler != nil {
				messageHandler(msg)
			}
		}
	}
}
