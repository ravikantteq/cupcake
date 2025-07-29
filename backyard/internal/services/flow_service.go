package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ravikantteq/cupcake/backyard/internal/models"
	"github.com/ravikantteq/cupcake/backyard/internal/repository"
	"github.com/ravikantteq/cupcake/backyard/pkg/netw"
	"github.com/ravikantteq/cupcake/backyard/pkg/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FlowService handles test flow operations
type FlowService struct {
	flowRepo      *repository.FlowRepository
	executionRepo *repository.ExecutionRepository
	messageRepo   *repository.MessageRepository
	validator     *validation.MessageValidator
	kafkaBroker   string
}

// NewFlowService creates a new flow service
func NewFlowService(repo *repository.Repository, kafkaBroker string) *FlowService {
	return &FlowService{
		flowRepo:      repo.NewFlowRepository(),
		executionRepo: repo.NewExecutionRepository(),
		messageRepo:   repo.NewMessageRepository(),
		validator:     validation.NewMessageValidator(),
		kafkaBroker:   kafkaBroker,
	}
}

// CreateFlow creates a new test flow
func (fs *FlowService) CreateFlow(ctx context.Context, req *models.CreateFlowRequest) (*models.TestFlow, error) {
	flow := &models.TestFlow{
		Name:        req.Name,
		Description: req.Description,
		Version:     req.Version,
		Steps:       req.Steps,
		CreatedBy:   "system", // TODO: Get from auth context
	}

	// Validate flow steps
	if err := fs.validateFlowSteps(flow.Steps); err != nil {
		return nil, fmt.Errorf("invalid flow steps: %w", err)
	}

	return fs.flowRepo.CreateFlow(ctx, flow)
}

// UpdateFlow updates an existing test flow
func (fs *FlowService) UpdateFlow(ctx context.Context, id primitive.ObjectID, req *models.CreateFlowRequest) (*models.TestFlow, error) {
	// First, get the existing flow to preserve some fields
	existingFlow, err := fs.flowRepo.GetFlowByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("flow not found: %w", err)
	}

	// Update the flow with new data
	flow := &models.TestFlow{
		ID:          existingFlow.ID,
		Name:        req.Name,
		Description: req.Description,
		Version:     req.Version,
		Steps:       req.Steps,
		CreatedBy:   existingFlow.CreatedBy,
		CreatedAt:   existingFlow.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	// Validate flow steps
	if err := fs.validateFlowSteps(flow.Steps); err != nil {
		return nil, fmt.Errorf("invalid flow steps: %w", err)
	}

	return fs.flowRepo.UpdateFlow(ctx, id, flow)
}

// GetFlowByID retrieves a flow by ID
func (fs *FlowService) GetFlowByID(ctx context.Context, id primitive.ObjectID) (*models.TestFlow, error) {
	return fs.flowRepo.GetFlowByID(ctx, id)
}

// GetAllFlows retrieves all flows
func (fs *FlowService) GetAllFlows(ctx context.Context) ([]models.TestFlow, error) {
	return fs.flowRepo.GetAllFlows(ctx)
}

// ExecuteFlow executes a test flow
func (fs *FlowService) ExecuteFlow(ctx context.Context, flowID primitive.ObjectID, suiteID primitive.ObjectID) (*models.TestExecution, error) {
	flow, err := fs.flowRepo.GetFlowByID(ctx, flowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get flow: %w", err)
	}

	// Create execution record
	execution := &models.TestExecution{
		FlowID:  flowID,
		SuiteID: suiteID,
		Status:  models.StatusRunning,
		Steps:   make([]models.ExecutionStep, 0, len(flow.Steps)),
		Metrics: models.ExecutionMetrics{},
	}

	createdExecution, err := fs.executionRepo.CreateExecution(ctx, execution)
	if err != nil {
		return nil, fmt.Errorf("failed to create execution: %w", err)
	}

	// Execute steps asynchronously
	go fs.executeFlowSteps(context.Background(), createdExecution, flow)

	return createdExecution, nil
}

// executeFlowSteps executes all steps in a flow
func (fs *FlowService) executeFlowSteps(ctx context.Context, execution *models.TestExecution, flow *models.TestFlow) {
	log.Printf("Starting execution of flow %s (ID: %s)", flow.Name, execution.ID.Hex())

	stepResults := make(map[string]interface{})

	for _, step := range flow.Steps {
		stepStartTime := time.Now()

		executionStep := models.ExecutionStep{
			StepID: step.StepID,
			Status: models.StatusRunning,
		}

		var stepErr error
		switch step.Type {
		case models.StepTypeProduce:
			stepErr = fs.executeProduceStep(ctx, &step, &executionStep, stepResults, execution.ID)
		case models.StepTypeConsume:
			stepErr = fs.executeConsumeStep(ctx, &step, &executionStep, stepResults, execution.ID)
		case models.StepTypeValidate:
			stepErr = fs.executeValidateStep(ctx, &step, &executionStep, stepResults)
		case models.StepTypeDelay:
			stepErr = fs.executeDelayStep(ctx, &step, &executionStep)
		default:
			stepErr = fmt.Errorf("unknown step type: %s", step.Type)
		}

		// Update step execution result
		stepDuration := time.Since(stepStartTime).Milliseconds()
		executionStep.Duration = stepDuration

		if stepErr != nil {
			executionStep.Status = models.StatusFailed
			executionStep.Errors = []string{stepErr.Error()}
			execution.Metrics.ErrorsCount++
			log.Printf("Step %s failed: %v", step.StepID, stepErr)
		} else {
			executionStep.Status = models.StatusCompleted
			execution.Metrics.StepsCompleted++
		}

		execution.Steps = append(execution.Steps, executionStep)

		// Store step output for reference by subsequent steps
		if executionStep.Output != nil {
			for key, value := range executionStep.Output {
				fs.validator.SetStepValue(step.StepID, key, value)
				stepResults[fmt.Sprintf("%s.%s", step.StepID, key)] = value
			}
		}

		// Stop execution if step failed and not configured to continue
		if stepErr != nil {
			execution.Status = models.StatusFailed
			break
		}
	}

	// Update final execution status
	if execution.Status == models.StatusRunning {
		execution.Status = models.StatusCompleted
	}

	endTime := time.Now()
	execution.EndTime = &endTime
	execution.Metrics.TotalDuration = endTime.Sub(execution.StartTime).Milliseconds()

	// Update execution in database
	if err := fs.executionRepo.UpdateExecution(ctx, execution.ID, execution); err != nil {
		log.Printf("Failed to update execution: %v", err)
	}

	log.Printf("Completed execution of flow %s with status %s", flow.Name, execution.Status)
}

// executeProduceStep executes a produce step
func (fs *FlowService) executeProduceStep(ctx context.Context, step *models.FlowStep, execStep *models.ExecutionStep, stepResults map[string]interface{}, executionID primitive.ObjectID) error {
	config := step.Config

	// Generate message from template
	messageData, err := fs.generateMessageFromTemplate(config.Message, stepResults)
	if err != nil {
		return fmt.Errorf("failed to generate message: %w", err)
	}

	// Extract key and value like the legacy producer does
	key := ""
	if keyVal, exists := messageData["key"]; exists {
		key = fmt.Sprintf("%v", keyVal)
	}

	value := ""
	if valueVal, exists := messageData["value"]; exists {
		value = fmt.Sprintf("%v", valueVal)
	} else {
		return fmt.Errorf("message must have a 'value' field")
	}

	// Create producer
	producer := netw.NewKafkaProducer(fs.kafkaBroker, config.Topic)

	// Produce message using the same method as legacy producer
	err = producer.ProduceJSON(key, value)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	// Store message in database - store the actual key/value that was sent
	message := &models.Message{
		Topic:       config.Topic,
		Key:         key,
		Value:       map[string]interface{}{"key": key, "value": value},
		ExecutionID: &executionID,
	}

	if err := fs.messageRepo.StoreMessage(ctx, message); err != nil {
		log.Printf("Failed to store produced message: %v", err)
	}

	// Update step execution
	execStep.Input = config.Message
	execStep.Output = map[string]interface{}{"key": key, "value": value}

	log.Printf("Produced message to topic %s: %s", config.Topic, value)
	return nil
}

// executeConsumeStep executes a consume step
func (fs *FlowService) executeConsumeStep(ctx context.Context, step *models.FlowStep, execStep *models.ExecutionStep, _ map[string]interface{}, _ primitive.ObjectID) error {
	config := step.Config

	// TODO: Implement Kafka consumer to read messages
	// For now, simulate consuming a message
	timeout := time.Duration(config.Timeout) * time.Millisecond
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	// Create a context with timeout
	consumeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Simulate message consumption (in real implementation, this would use Kafka consumer)
	select {
	case <-time.After(1 * time.Second): // Simulate processing time
		// Simulate consumed message
		consumedMessage := map[string]interface{}{
			"topic":     config.Topic,
			"partition": 0,
			"offset":    1,
			"key":       "test-key",
			"value": map[string]interface{}{
				"status":      "processed",
				"processedAt": time.Now().Format(time.RFC3339),
			},
		}

		execStep.Output = consumedMessage
		log.Printf("Consumed message from topic %s", config.Topic)
		return nil

	case <-consumeCtx.Done():
		return fmt.Errorf("timeout waiting for message from topic %s", config.Topic)
	}
}

// executeValidateStep executes a validation step
func (fs *FlowService) executeValidateStep(_ context.Context, step *models.FlowStep, execStep *models.ExecutionStep, stepResults map[string]interface{}) error {
	config := step.Config

	// Get the message to validate (from previous consume step)
	var actualMessage map[string]interface{}

	// Look for the most recent consume step output
	for _, value := range stepResults {
		if valueMap, ok := value.(map[string]interface{}); ok {
			if _, hasValue := valueMap["value"]; hasValue {
				if actualValue, ok := valueMap["value"].(map[string]interface{}); ok {
					actualMessage = actualValue
					break
				}
			}
		}
	}

	if actualMessage == nil {
		return fmt.Errorf("no message found to validate")
	}

	// Validate message against expected pattern
	result := fs.validator.ValidateMessage(actualMessage, config.ExpectedMessage)

	execStep.Input = map[string]interface{}{
		"actual":   actualMessage,
		"expected": config.ExpectedMessage,
	}
	execStep.Output = map[string]interface{}{
		"validationResult": result,
	}

	if !result.Valid {
		return fmt.Errorf("validation failed: %s", result.Message)
	}

	log.Printf("Validation passed for step %s", step.StepID)
	return nil
}

// executeDelayStep executes a delay step
func (fs *FlowService) executeDelayStep(_ context.Context, step *models.FlowStep, execStep *models.ExecutionStep) error {
	config := step.Config
	delayMs := config.DelayMs
	if delayMs <= 0 {
		delayMs = 1000 // Default 1 second
	}

	delay := time.Duration(delayMs) * time.Millisecond

	execStep.Input = map[string]interface{}{
		"delayMs": delayMs,
	}

	log.Printf("Delaying for %v", delay)
	time.Sleep(delay)

	execStep.Output = map[string]interface{}{
		"delayed": delayMs,
	}

	return nil
}

// generateMessageFromTemplate generates a message from a template with dynamic values
func (fs *FlowService) generateMessageFromTemplate(template map[string]interface{}, stepResults map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for key, value := range template {
		generatedValue, err := fs.generateValue(value, stepResults)
		if err != nil {
			return nil, fmt.Errorf("failed to generate value for field %s: %w", key, err)
		}
		result[key] = generatedValue
	}

	return result, nil
}

// generateValue generates a value based on the template specification
func (fs *FlowService) generateValue(valueSpec interface{}, stepResults map[string]interface{}) (interface{}, error) {
	// Handle nested objects
	if objSpec, ok := valueSpec.(map[string]interface{}); ok {
		return fs.generateMessageFromTemplate(objSpec, stepResults)
	}

	// Handle arrays
	if arrSpec, ok := valueSpec.([]interface{}); ok {
		result := make([]interface{}, len(arrSpec))
		for i, itemSpec := range arrSpec {
			generatedItem, err := fs.generateValue(itemSpec, stepResults)
			if err != nil {
				return nil, err
			}
			result[i] = generatedItem
		}
		return result, nil
	}

	// Handle string templates
	if strSpec, ok := valueSpec.(string); ok {
		return fs.generateStringValue(strSpec, stepResults)
	}

	// Return value as-is for other types
	return valueSpec, nil
}

// generateStringValue generates a value from a string specification
func (fs *FlowService) generateStringValue(spec string, stepResults map[string]interface{}) (interface{}, error) {
	switch spec {
	case "uuid()":
		return primitive.NewObjectID().Hex(), nil // Using ObjectID as UUID substitute
	case "timestamp()":
		return time.Now().Format(time.RFC3339), nil
	default:
		// Check for step references
		if value, exists := stepResults[spec]; exists {
			return value, nil
		}
		// Return as literal string
		return spec, nil
	}
}

// validateFlowSteps validates the structure of flow steps
func (fs *FlowService) validateFlowSteps(steps []models.FlowStep) error {
	if len(steps) == 0 {
		return fmt.Errorf("flow must have at least one step")
	}

	stepIDs := make(map[string]bool)

	for _, step := range steps {
		if step.StepID == "" {
			return fmt.Errorf("step ID cannot be empty")
		}

		if stepIDs[step.StepID] {
			return fmt.Errorf("duplicate step ID: %s", step.StepID)
		}
		stepIDs[step.StepID] = true

		switch step.Type {
		case models.StepTypeProduce:
			if step.Config.Topic == "" {
				return fmt.Errorf("produce step %s must have a topic", step.StepID)
			}
			if step.Config.Message == nil {
				return fmt.Errorf("produce step %s must have a message", step.StepID)
			}
		case models.StepTypeConsume:
			if step.Config.Topic == "" {
				return fmt.Errorf("consume step %s must have a topic", step.StepID)
			}
		case models.StepTypeValidate:
			if step.Config.ExpectedMessage == nil {
				return fmt.Errorf("validate step %s must have an expected message", step.StepID)
			}
		case models.StepTypeDelay:
			// Delay steps are always valid
		default:
			return fmt.Errorf("unknown step type: %s", step.Type)
		}
	}

	return nil
}
