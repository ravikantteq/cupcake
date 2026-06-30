package manager

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ravikantteq/cupcake/backyard/internal"
	"github.com/ravikantteq/cupcake/backyard/internal/store"
	"github.com/ravikantteq/cupcake/backyard/pkg/netw"
)

// ProducerManager manages Kafka message production and history
type ProducerManager struct {
	store store.Store
}

// NewProducerManager creates a new producer manager
func NewProducerManager(store store.Store) *ProducerManager {
	return &ProducerManager{
		store: store,
	}
}

// PublishMessage publishes a message to Kafka and stores history
func (pm *ProducerManager) PublishMessage(ctx context.Context, req *internal.KafkaMessageRequest) error {
	// Convert value to string - handle both string and JSON object inputs
	var valueStr string
	switch v := req.Value.(type) {
	case string:
		valueStr = v
	case map[string]interface{}, []interface{}, map[string]string:
		// If it's a JSON object/array, marshal it to JSON string
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON value: %w", err)
		}
		valueStr = string(jsonBytes)
	default:
		// For other types, convert to JSON
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		valueStr = string(jsonBytes)
	}

	// Create producer
	producer := netw.NewKafkaProducer(req.Broker, req.Topic)

	// Publish message
	err := producer.ProduceJSON(req.Key, valueStr)

	// Create history record
	history := &internal.ProducerHistory{
		Broker:  req.Broker,
		Topic:   req.Topic,
		Key:     req.Key,
		Value:   valueStr,
		Success: err == nil,
	}

	if err != nil {
		history.Error = err.Error()
	} else {
		history.Response = "Message published successfully"
	}

	// Store history
	storeErr := pm.store.StoreProducerHistory(ctx, history)
	if storeErr != nil {
		fmt.Printf("Failed to store producer history: %v\n", storeErr)
	}

	return err
}

// GetProducerHistory retrieves producer history
func (pm *ProducerManager) GetProducerHistory(ctx context.Context, limit int) ([]*internal.ProducerHistory, error) {
	return pm.store.GetProducerHistory(ctx, limit)
}

// GetRecentProducerHistory retrieves recent producer history (last 50)
func (pm *ProducerManager) GetRecentProducerHistory(ctx context.Context) ([]*internal.ProducerHistory, error) {
	return pm.store.GetProducerHistory(ctx, 50)
}
