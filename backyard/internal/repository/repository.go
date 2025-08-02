package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ravikantteq/cupcake/backyard/internal/models"
	"github.com/ravikantteq/cupcake/backyard/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository handles all database operations
type Repository struct {
	db *storage.MongoDB
}

// NewRepository creates a new repository instance
func NewRepository(db *storage.MongoDB) *Repository {
	return &Repository{db: db}
}

// FlowRepository handles test flow operations
type FlowRepository struct {
	collection *mongo.Collection
}

// NewFlowRepository creates a new flow repository
func (r *Repository) NewFlowRepository() *FlowRepository {
	return &FlowRepository{
		collection: r.db.GetCollection("flows"),
	}
}

// CreateFlow creates a new test flow
func (fr *FlowRepository) CreateFlow(ctx context.Context, flow *models.TestFlow) (*models.TestFlow, error) {
	flow.ID = primitive.NewObjectID()
	flow.CreatedAt = time.Now()
	flow.UpdatedAt = time.Now()

	_, err := fr.collection.InsertOne(ctx, flow)
	if err != nil {
		return nil, fmt.Errorf("failed to create flow: %w", err)
	}

	return flow, nil
}

// GetFlowByID retrieves a flow by ID
func (fr *FlowRepository) GetFlowByID(ctx context.Context, id primitive.ObjectID) (*models.TestFlow, error) {
	var flow models.TestFlow
	err := fr.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&flow)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("flow not found")
		}
		return nil, fmt.Errorf("failed to get flow: %w", err)
	}

	return &flow, nil
}

// GetAllFlows retrieves all flows
func (fr *FlowRepository) GetAllFlows(ctx context.Context) ([]models.TestFlow, error) {
	cursor, err := fr.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find flows: %w", err)
	}
	defer cursor.Close(ctx)

	var flows []models.TestFlow
	if err = cursor.All(ctx, &flows); err != nil {
		return nil, fmt.Errorf("failed to decode flows: %w", err)
	}

	return flows, nil
}

// UpdateFlow updates an existing flow
func (fr *FlowRepository) UpdateFlow(ctx context.Context, id primitive.ObjectID, flow *models.TestFlow) (*models.TestFlow, error) {
	flow.UpdatedAt = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": flow}

	result, err := fr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update flow: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("flow not found")
	}

	return fr.GetFlowByID(ctx, id)
}

// DeleteFlow deletes a flow
func (fr *FlowRepository) DeleteFlow(ctx context.Context, id primitive.ObjectID) error {
	result, err := fr.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete flow: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("flow not found")
	}

	return nil
}

// SuiteRepository handles test suite operations
type SuiteRepository struct {
	collection *mongo.Collection
}

// NewSuiteRepository creates a new suite repository
func (r *Repository) NewSuiteRepository() *SuiteRepository {
	return &SuiteRepository{
		collection: r.db.GetCollection("suites"),
	}
}

// CreateSuite creates a new test suite
func (sr *SuiteRepository) CreateSuite(ctx context.Context, suite *models.TestSuite) (*models.TestSuite, error) {
	suite.ID = primitive.NewObjectID()
	suite.CreatedAt = time.Now()
	suite.UpdatedAt = time.Now()

	_, err := sr.collection.InsertOne(ctx, suite)
	if err != nil {
		return nil, fmt.Errorf("failed to create suite: %w", err)
	}

	return suite, nil
}

// GetSuiteByID retrieves a suite by ID
func (sr *SuiteRepository) GetSuiteByID(ctx context.Context, id primitive.ObjectID) (*models.TestSuite, error) {
	var suite models.TestSuite
	err := sr.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&suite)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("suite not found")
		}
		return nil, fmt.Errorf("failed to get suite: %w", err)
	}

	return &suite, nil
}

// GetAllSuites retrieves all suites
func (sr *SuiteRepository) GetAllSuites(ctx context.Context) ([]models.TestSuite, error) {
	cursor, err := sr.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find suites: %w", err)
	}
	defer cursor.Close(ctx)

	var suites []models.TestSuite
	if err = cursor.All(ctx, &suites); err != nil {
		return nil, fmt.Errorf("failed to decode suites: %w", err)
	}

	return suites, nil
}

// ExecutionRepository handles test execution operations
type ExecutionRepository struct {
	collection *mongo.Collection
}

// NewExecutionRepository creates a new execution repository
func (r *Repository) NewExecutionRepository() *ExecutionRepository {
	return &ExecutionRepository{
		collection: r.db.GetCollection("executions"),
	}
}

// CreateExecution creates a new test execution
func (er *ExecutionRepository) CreateExecution(ctx context.Context, execution *models.TestExecution) (*models.TestExecution, error) {
	execution.ID = primitive.NewObjectID()
	execution.StartTime = time.Now()
	execution.Status = models.StatusRunning

	_, err := er.collection.InsertOne(ctx, execution)
	if err != nil {
		return nil, fmt.Errorf("failed to create execution: %w", err)
	}

	return execution, nil
}

// UpdateExecution updates an execution
func (er *ExecutionRepository) UpdateExecution(ctx context.Context, id primitive.ObjectID, execution *models.TestExecution) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": execution}

	_, err := er.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update execution: %w", err)
	}

	return nil
}

// GetExecutionByID retrieves an execution by ID
func (er *ExecutionRepository) GetExecutionByID(ctx context.Context, id primitive.ObjectID) (*models.TestExecution, error) {
	var execution models.TestExecution
	err := er.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&execution)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("execution not found")
		}
		return nil, fmt.Errorf("failed to get execution: %w", err)
	}

	return &execution, nil
}

// GetExecutionsByFlowID retrieves executions for a specific flow
func (er *ExecutionRepository) GetExecutionsByFlowID(ctx context.Context, flowID primitive.ObjectID, limit int) ([]models.TestExecution, error) {
	opts := options.Find().SetSort(bson.D{{Key: "startTime", Value: -1}}).SetLimit(int64(limit))
	cursor, err := er.collection.Find(ctx, bson.M{"flowId": flowID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find executions: %w", err)
	}
	defer cursor.Close(ctx)

	var executions []models.TestExecution
	if err = cursor.All(ctx, &executions); err != nil {
		return nil, fmt.Errorf("failed to decode executions: %w", err)
	}

	return executions, nil
}

// ConsumerRepository handles consumer operations
type ConsumerRepository struct {
	collection *mongo.Collection
}

// NewConsumerRepository creates a new consumer repository
func (r *Repository) NewConsumerRepository() *ConsumerRepository {
	return &ConsumerRepository{
		collection: r.db.GetCollection("consumers"),
	}
}

// CreateConsumer creates a new consumer
func (cr *ConsumerRepository) CreateConsumer(ctx context.Context, consumer *models.Consumer) (*models.Consumer, error) {
	consumer.ID = primitive.NewObjectID()
	consumer.CreatedAt = time.Now()
	consumer.Status = models.ConsumerInactive

	_, err := cr.collection.InsertOne(ctx, consumer)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return consumer, nil
}

// GetConsumerByGroupID retrieves a consumer by group ID
func (cr *ConsumerRepository) GetConsumerByGroupID(ctx context.Context, groupID string) (*models.Consumer, error) {
	var consumer models.Consumer
	err := cr.collection.FindOne(ctx, bson.M{"groupId": groupID}).Decode(&consumer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("consumer not found")
		}
		return nil, fmt.Errorf("failed to get consumer: %w", err)
	}

	return &consumer, nil
}

// GetAllConsumers retrieves all consumers
func (cr *ConsumerRepository) GetAllConsumers(ctx context.Context) ([]models.Consumer, error) {
	cursor, err := cr.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find consumers: %w", err)
	}
	defer cursor.Close(ctx)

	var consumers []models.Consumer
	if err = cursor.All(ctx, &consumers); err != nil {
		return nil, fmt.Errorf("failed to decode consumers: %w", err)
	}

	return consumers, nil
}

// UpdateConsumerStatus updates consumer status
func (cr *ConsumerRepository) UpdateConsumerStatus(ctx context.Context, groupID string, status models.ConsumerStatus, errorMsg string) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":        status,
			"lastHeartbeat": now,
		},
	}

	if errorMsg != "" {
		update["$set"].(bson.M)["errorMessage"] = errorMsg
	}

	_, err := cr.collection.UpdateOne(ctx, bson.M{"groupId": groupID}, update)
	if err != nil {
		return fmt.Errorf("failed to update consumer status: %w", err)
	}

	return nil
}

// MessageRepository handles message storage operations
type MessageRepository struct {
	collection *mongo.Collection
}

// NewMessageRepository creates a new message repository
func (r *Repository) NewMessageRepository() *MessageRepository {
	return &MessageRepository{
		collection: r.db.GetCollection("messages"),
	}
}

// StoreMessage stores a Kafka message
func (mr *MessageRepository) StoreMessage(ctx context.Context, message *models.Message) error {
	message.ID = primitive.NewObjectID()
	message.Timestamp = time.Now()

	_, err := mr.collection.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to store message: %w", err)
	}

	return nil
}

// GetMessagesByTopic retrieves messages for a specific topic
func (mr *MessageRepository) GetMessagesByTopic(ctx context.Context, topic string, limit int) ([]models.Message, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(int64(limit))
	cursor, err := mr.collection.Find(ctx, bson.M{"topic": topic}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages: %w", err)
	}

	return messages, nil
}

// GetMessagesByExecutionID retrieves messages for a specific execution
func (mr *MessageRepository) GetMessagesByExecutionID(ctx context.Context, executionID primitive.ObjectID) ([]models.Message, error) {
	cursor, err := mr.collection.Find(ctx, bson.M{"executionId": executionID})
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages: %w", err)
	}

	return messages, nil
}

// ProducerHistoryRepository handles producer history operations
type ProducerHistoryRepository struct {
	collection *mongo.Collection
}

// NewProducerHistoryRepository creates a new producer history repository
func (r *Repository) NewProducerHistoryRepository() *ProducerHistoryRepository {
	return &ProducerHistoryRepository{
		collection: r.db.GetCollection("producer_history"),
	}
}

// CreateHistoryEntry creates a new producer history entry
func (phr *ProducerHistoryRepository) CreateHistoryEntry(ctx context.Context, history *models.ProducerHistory) (*models.ProducerHistory, error) {
	history.ID = primitive.NewObjectID()
	if history.Timestamp.IsZero() {
		history.Timestamp = time.Now()
	}

	_, err := phr.collection.InsertOne(ctx, history)
	if err != nil {
		return nil, fmt.Errorf("failed to create history entry: %w", err)
	}

	return history, nil
}

// GetRecentHistory retrieves recent producer history entries
func (phr *ProducerHistoryRepository) GetRecentHistory(ctx context.Context, limit int) ([]models.ProducerHistory, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(int64(limit))

	cursor, err := phr.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find history entries: %w", err)
	}
	defer cursor.Close(ctx)

	var history []models.ProducerHistory
	if err = cursor.All(ctx, &history); err != nil {
		return nil, fmt.Errorf("failed to decode history entries: %w", err)
	}

	return history, nil
}

// GetHistoryByUser retrieves producer history for a specific user
func (phr *ProducerHistoryRepository) GetHistoryByUser(ctx context.Context, userID string, limit int, offset int) ([]models.ProducerHistory, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	filter := bson.M{}
	if userID != "" {
		filter["userId"] = userID
	}

	cursor, err := phr.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find history entries: %w", err)
	}
	defer cursor.Close(ctx)

	var history []models.ProducerHistory
	if err = cursor.All(ctx, &history); err != nil {
		return nil, fmt.Errorf("failed to decode history entries: %w", err)
	}

	return history, nil
}

// GetHistoryCount returns the total count of history entries
func (phr *ProducerHistoryRepository) GetHistoryCount(ctx context.Context, userID string) (int64, error) {
	filter := bson.M{}
	if userID != "" {
		filter["userId"] = userID
	}

	count, err := phr.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count history entries: %w", err)
	}

	return count, nil
}

// DeleteOldHistory deletes history entries older than the specified duration
func (phr *ProducerHistoryRepository) DeleteOldHistory(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-olderThan)
	filter := bson.M{"timestamp": bson.M{"$lt": cutoffTime}}

	result, err := phr.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old history entries: %w", err)
	}

	return result.DeletedCount, nil
}
