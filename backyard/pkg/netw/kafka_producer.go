package netw

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kp.broker})
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	defer p.Close()

	// Optional delivery channel
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.topic, Partition: kp.partition},
		Key:            []byte(key),
		Value:          []byte(value),
		Headers:        []kafka.Header{{Key: "backyard-header", Value: []byte("Message produced through backyard service")}},
	}, deliveryChan)

	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	// Wait for delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %v", m.TopicPartition.Error)
	}

	log.Printf("Delivered message to topic %s [%d] at offset %v",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)

	return nil
}
