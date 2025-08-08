package store

import (
	"context"

	"github.com/ravikantteq/cupcake/backyard/internal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Store defines the interface for data persistence
type Store interface {
	ConsumerStore
	MessageStore
	FlowStore
	ExecutionStore
	ProducerHistoryStore
}

// ConsumerStore defines operations for consumer persistence
type ConsumerStore interface {
	CreateConsumer(ctx context.Context, consumer *internal.Consumer) (*internal.Consumer, error)
	GetConsumer(ctx context.Context, id primitive.ObjectID) (*internal.Consumer, error)
	GetConsumers(ctx context.Context) ([]*internal.Consumer, error)
	UpdateConsumer(ctx context.Context, consumer *internal.Consumer) error
	DeleteConsumer(ctx context.Context, id primitive.ObjectID) error
	UpdateConsumerStatus(ctx context.Context, id primitive.ObjectID, status internal.Status, errorMsg string) error
	IncrementMessageCount(ctx context.Context, id primitive.ObjectID) error
}

// MessageStore defines operations for message persistence
type MessageStore interface {
	StoreMessage(ctx context.Context, message *internal.Message) error
	GetMessages(ctx context.Context, filters MessageFilters) ([]*internal.Message, error)
}

// FlowStore defines operations for flow persistence
type FlowStore interface {
	CreateFlow(ctx context.Context, flow *internal.Flow) (*internal.Flow, error)
	GetFlow(ctx context.Context, id primitive.ObjectID) (*internal.Flow, error)
	GetFlows(ctx context.Context) ([]*internal.Flow, error)
	UpdateFlow(ctx context.Context, flow *internal.Flow) error
	DeleteFlow(ctx context.Context, id primitive.ObjectID) error
}

// ExecutionStore defines operations for execution persistence
type ExecutionStore interface {
	CreateExecution(ctx context.Context, execution *internal.Execution) (*internal.Execution, error)
	GetExecution(ctx context.Context, id primitive.ObjectID) (*internal.Execution, error)
	GetExecutions(ctx context.Context, flowID primitive.ObjectID) ([]*internal.Execution, error)
	GetAllExecutions(ctx context.Context) ([]*internal.Execution, error)
	UpdateExecution(ctx context.Context, execution *internal.Execution) error
}

// ProducerHistoryStore defines operations for producer history persistence
type ProducerHistoryStore interface {
	StoreProducerHistory(ctx context.Context, history *internal.ProducerHistory) error
	GetProducerHistory(ctx context.Context, limit int) ([]*internal.ProducerHistory, error)
}

// MessageFilters defines filters for message queries
type MessageFilters struct {
	Topic      string
	ConsumerID *primitive.ObjectID
	Limit      int
}
