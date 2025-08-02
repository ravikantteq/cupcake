package store

import (
	"context"
	"time"

	"github.com/ravikantteq/cupcake/backyard/internal"
	"github.com/ravikantteq/cupcake/backyard/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB implements Store interface using MongoDB
type MongoDB struct {
	db *storage.MongoDB
}

// NewMongoDB creates a new MongoDB store
func NewMongoDB(db *storage.MongoDB) *MongoDB {
	return &MongoDB{db: db}
}

// Consumer operations
func (m *MongoDB) CreateConsumer(ctx context.Context, consumer *internal.Consumer) (*internal.Consumer, error) {
	if consumer.ID.IsZero() {
		consumer.ID = primitive.NewObjectID()
	}

	now := time.Now()
	consumer.CreatedAt = now
	consumer.UpdatedAt = now
	consumer.Status = internal.StatusInactive

	result, err := m.db.Database.Collection("consumers").InsertOne(ctx, consumer)
	if err != nil {
		return nil, err
	}

	consumer.ID = result.InsertedID.(primitive.ObjectID)
	return consumer, nil
}

func (m *MongoDB) GetConsumer(ctx context.Context, id primitive.ObjectID) (*internal.Consumer, error) {
	var consumer internal.Consumer
	err := m.db.Database.Collection("consumers").FindOne(ctx, bson.M{"_id": id}).Decode(&consumer)
	if err != nil {
		return nil, err
	}
	return &consumer, nil
}

func (m *MongoDB) GetConsumers(ctx context.Context) ([]*internal.Consumer, error) {
	cursor, err := m.db.Database.Collection("consumers").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var consumers []*internal.Consumer
	for cursor.Next(ctx) {
		var consumer internal.Consumer
		if err := cursor.Decode(&consumer); err != nil {
			return nil, err
		}
		consumers = append(consumers, &consumer)
	}

	return consumers, cursor.Err()
}

func (m *MongoDB) UpdateConsumer(ctx context.Context, consumer *internal.Consumer) error {
	consumer.UpdatedAt = time.Now()

	filter := bson.M{"_id": consumer.ID}
	update := bson.M{"$set": consumer}

	_, err := m.db.Database.Collection("consumers").UpdateOne(ctx, filter, update)
	return err
}

func (m *MongoDB) DeleteConsumer(ctx context.Context, id primitive.ObjectID) error {
	_, err := m.db.Database.Collection("consumers").DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (m *MongoDB) UpdateConsumerStatus(ctx context.Context, id primitive.ObjectID, status internal.Status, errorMsg string) error {
	update := bson.M{
		"$set": bson.M{
			"status":       status,
			"errorMessage": errorMsg,
			"updatedAt":    time.Now(),
		},
	}

	_, err := m.db.Database.Collection("consumers").UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (m *MongoDB) IncrementMessageCount(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$inc": bson.M{"messageCount": 1},
		"$set": bson.M{"updatedAt": time.Now()},
	}

	_, err := m.db.Database.Collection("consumers").UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Message operations
func (m *MongoDB) StoreMessage(ctx context.Context, message *internal.Message) error {
	if message.ID.IsZero() {
		message.ID = primitive.NewObjectID()
	}

	_, err := m.db.Database.Collection("messages").InsertOne(ctx, message)
	return err
}

func (m *MongoDB) GetMessages(ctx context.Context, filters MessageFilters) ([]*internal.Message, error) {
	query := bson.M{}

	if filters.Topic != "" {
		query["topic"] = filters.Topic
	}

	if filters.ConsumerID != nil {
		query["consumerId"] = *filters.ConsumerID
	}

	opts := options.Find()
	if filters.Limit > 0 {
		opts.SetLimit(int64(filters.Limit))
	}
	opts.SetSort(bson.M{"timestamp": -1}) // Most recent first

	cursor, err := m.db.Database.Collection("messages").Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*internal.Message
	for cursor.Next(ctx) {
		var message internal.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, cursor.Err()
}

// Flow operations
func (m *MongoDB) CreateFlow(ctx context.Context, flow *internal.Flow) (*internal.Flow, error) {
	if flow.ID.IsZero() {
		flow.ID = primitive.NewObjectID()
	}

	now := time.Now()
	flow.CreatedAt = now
	flow.UpdatedAt = now

	result, err := m.db.Database.Collection("flows").InsertOne(ctx, flow)
	if err != nil {
		return nil, err
	}

	flow.ID = result.InsertedID.(primitive.ObjectID)
	return flow, nil
}

func (m *MongoDB) GetFlow(ctx context.Context, id primitive.ObjectID) (*internal.Flow, error) {
	var flow internal.Flow
	err := m.db.Database.Collection("flows").FindOne(ctx, bson.M{"_id": id}).Decode(&flow)
	if err != nil {
		return nil, err
	}
	return &flow, nil
}

func (m *MongoDB) GetFlows(ctx context.Context) ([]*internal.Flow, error) {
	cursor, err := m.db.Database.Collection("flows").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var flows []*internal.Flow
	for cursor.Next(ctx) {
		var flow internal.Flow
		if err := cursor.Decode(&flow); err != nil {
			return nil, err
		}
		flows = append(flows, &flow)
	}

	return flows, cursor.Err()
}

func (m *MongoDB) UpdateFlow(ctx context.Context, flow *internal.Flow) error {
	flow.UpdatedAt = time.Now()

	filter := bson.M{"_id": flow.ID}
	update := bson.M{"$set": flow}

	_, err := m.db.Database.Collection("flows").UpdateOne(ctx, filter, update)
	return err
}

func (m *MongoDB) DeleteFlow(ctx context.Context, id primitive.ObjectID) error {
	_, err := m.db.Database.Collection("flows").DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Execution operations
func (m *MongoDB) CreateExecution(ctx context.Context, execution *internal.Execution) (*internal.Execution, error) {
	if execution.ID.IsZero() {
		execution.ID = primitive.NewObjectID()
	}

	execution.StartTime = time.Now()
	execution.Status = internal.StatusRunning

	result, err := m.db.Database.Collection("executions").InsertOne(ctx, execution)
	if err != nil {
		return nil, err
	}

	execution.ID = result.InsertedID.(primitive.ObjectID)
	return execution, nil
}

func (m *MongoDB) GetExecution(ctx context.Context, id primitive.ObjectID) (*internal.Execution, error) {
	var execution internal.Execution
	err := m.db.Database.Collection("executions").FindOne(ctx, bson.M{"_id": id}).Decode(&execution)
	if err != nil {
		return nil, err
	}
	return &execution, nil
}

func (m *MongoDB) GetExecutions(ctx context.Context, flowID primitive.ObjectID) ([]*internal.Execution, error) {
	cursor, err := m.db.Database.Collection("executions").Find(ctx, bson.M{"flowId": flowID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var executions []*internal.Execution
	for cursor.Next(ctx) {
		var execution internal.Execution
		if err := cursor.Decode(&execution); err != nil {
			return nil, err
		}
		executions = append(executions, &execution)
	}

	return executions, cursor.Err()
}

func (m *MongoDB) UpdateExecution(ctx context.Context, execution *internal.Execution) error {
	filter := bson.M{"_id": execution.ID}
	update := bson.M{"$set": execution}

	_, err := m.db.Database.Collection("executions").UpdateOne(ctx, filter, update)
	return err
}

// Producer history operations
func (m *MongoDB) StoreProducerHistory(ctx context.Context, history *internal.ProducerHistory) error {
	if history.ID.IsZero() {
		history.ID = primitive.NewObjectID()
	}

	history.Timestamp = time.Now()

	_, err := m.db.Database.Collection("producer_history").InsertOne(ctx, history)
	return err
}

func (m *MongoDB) GetProducerHistory(ctx context.Context, limit int) ([]*internal.ProducerHistory, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	opts.SetSort(bson.M{"timestamp": -1}) // Most recent first

	cursor, err := m.db.Database.Collection("producer_history").Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var histories []*internal.ProducerHistory
	for cursor.Next(ctx) {
		var history internal.ProducerHistory
		if err := cursor.Decode(&history); err != nil {
			return nil, err
		}
		histories = append(histories, &history)
	}

	return histories, cursor.Err()
}
