package netw

import (
	"testing"
)

func TestNewKafkaProducer(t *testing.T) {
	broker := "localhost:9092"
	topic := "test-topic"

	producer := NewKafkaProducer(broker, topic)

	if producer == nil {
		t.Fatal("Expected producer to be created, got nil")
	}

	if producer.broker != broker {
		t.Errorf("Expected broker %s, got %s", broker, producer.broker)
	}

	if producer.topic != topic {
		t.Errorf("Expected topic %s, got %s", topic, producer.topic)
	}
}

// Note: This test requires a running Kafka instance
// For integration testing, you would need to set up a test Kafka broker
func TestProduceJSON(t *testing.T) {
	// Skip this test if no Kafka broker is available
	t.Skip("Skipping integration test - requires running Kafka broker")

	broker := "localhost:9092"
	topic := "test-topic"

	producer := NewKafkaProducer(broker, topic)

	err := producer.ProduceJSON("test-key", `{"message": "test-value"}`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
