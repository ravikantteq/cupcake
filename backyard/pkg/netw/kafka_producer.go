package netw

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	broker    string
	topic     string
	partition int32
}

func NewKafkaProducer(broker string, topic string) *KafkaProducer {
	return &KafkaProducer{
		broker:    broker,
		topic:     topic,
		partition: kafka.PartitionAny,
	}
}

// Produce publishes a message to Kafka (legacy method)
func (kp *KafkaProducer) Produce(key string, value string) {
	err := kp.ProduceJSON(key, value)
	if err != nil {
		log.Printf("Failed to produce message: %s", err)
	}
}

// ProduceJSON publishes a message to Kafka and returns error for API handling
func (kp *KafkaProducer) ProduceJSON(key string, value string) error {
	config := kafka.ConfigMap{
		"bootstrap.servers": kp.broker,
		"client.id":         "cupcake-producer",
		"acks":              "all",
	}

	producer, err := kafka.NewProducer(&config)
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	defer producer.Close()

	// Create the message
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kp.topic,
			Partition: kp.partition,
		},
		Key:   []byte(key),
		Value: []byte(value),
		Headers: []kafka.Header{
			{
				Key:   "backyard-header",
				Value: []byte("Message produced through backyard service"),
			},
		},
	}

	// Delivery report handler for produced messages
	deliveryChan := make(chan kafka.Event, 1000)
	defer close(deliveryChan)

	// Produce the message
	err = producer.Produce(message, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	// Wait for delivery report with timeout
	select {
	case e := <-deliveryChan:
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				return fmt.Errorf("delivery failed: %v", ev.TopicPartition.Error)
			}
			log.Printf("Delivered message to topic %s [%d] at offset %v",
				*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
		case kafka.Error:
			return fmt.Errorf("delivery error: %v", ev)
		}
	case <-time.After(10 * time.Second):
		return fmt.Errorf("delivery timeout")
	}

	return nil
}
