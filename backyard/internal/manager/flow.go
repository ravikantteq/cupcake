package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/ravikantteq/cupcake/backyard/internal"
	"github.com/ravikantteq/cupcake/backyard/internal/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FlowManager manages test flows and their execution
type FlowManager struct {
	store store.Store
}

// NewFlowManager creates a new flow manager
func NewFlowManager(store store.Store) *FlowManager {
	return &FlowManager{
		store: store,
	}
}

// CreateFlow creates a new test flow
func (fm *FlowManager) CreateFlow(ctx context.Context, req *internal.CreateFlowRequest) (*internal.Flow, error) {
	// Validate steps
	if err := fm.validateFlowSteps(req.Steps); err != nil {
		return nil, fmt.Errorf("invalid flow steps: %w", err)
	}

	flow := &internal.Flow{
		Name:        req.Name,
		Description: req.Description,
		Version:     req.Version,
		Steps:       req.Steps,
		CreatedBy:   "system", // TODO: Get from auth context
	}

	return fm.store.CreateFlow(ctx, flow)
}

// GetFlow retrieves a flow by ID
func (fm *FlowManager) GetFlow(ctx context.Context, id primitive.ObjectID) (*internal.Flow, error) {
	return fm.store.GetFlow(ctx, id)
}

// GetFlows retrieves all flows
func (fm *FlowManager) GetFlows(ctx context.Context) ([]*internal.Flow, error) {
	return fm.store.GetFlows(ctx)
}

// UpdateFlow updates an existing flow
func (fm *FlowManager) UpdateFlow(ctx context.Context, id primitive.ObjectID, req *internal.CreateFlowRequest) (*internal.Flow, error) {
	// Get existing flow
	existingFlow, err := fm.store.GetFlow(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("flow not found: %w", err)
	}

	// Validate steps
	if err := fm.validateFlowSteps(req.Steps); err != nil {
		return nil, fmt.Errorf("invalid flow steps: %w", err)
	}

	// Update flow
	existingFlow.Name = req.Name
	existingFlow.Description = req.Description
	existingFlow.Version = req.Version
	existingFlow.Steps = req.Steps
	existingFlow.UpdatedAt = time.Now()

	err = fm.store.UpdateFlow(ctx, existingFlow)
	if err != nil {
		return nil, err
	}

	return existingFlow, nil
}

// DeleteFlow deletes a flow
func (fm *FlowManager) DeleteFlow(ctx context.Context, id primitive.ObjectID) error {
	return fm.store.DeleteFlow(ctx, id)
}

// ExecuteFlow executes a test flow
func (fm *FlowManager) ExecuteFlow(ctx context.Context, flowID primitive.ObjectID) (*internal.Execution, error) {
	flow, err := fm.store.GetFlow(ctx, flowID)
	if err != nil {
		return nil, fmt.Errorf("flow not found: %w", err)
	}

	// Create execution record
	execution := &internal.Execution{
		SuiteID: primitive.NewObjectID(), // Create a temporary suite ID for standalone flow execution
		FlowID:  flowID,
		Status:  internal.StatusRunning,
		Steps:   make([]internal.ExecutionStep, 0, len(flow.Steps)),
	}

	createdExecution, err := fm.store.CreateExecution(ctx, execution)
	if err != nil {
		return nil, fmt.Errorf("failed to create execution: %w", err)
	}

	// Execute steps asynchronously
	go fm.executeFlowStepsAsync(context.Background(), createdExecution, flow)

	return createdExecution, nil
}

// GetExecution retrieves an execution by ID
func (fm *FlowManager) GetExecution(ctx context.Context, id primitive.ObjectID) (*internal.Execution, error) {
	return fm.store.GetExecution(ctx, id)
}

// GetExecutions retrieves executions for a flow
func (fm *FlowManager) GetExecutions(ctx context.Context, flowID primitive.ObjectID) ([]*internal.Execution, error) {
	return fm.store.GetExecutions(ctx, flowID)
}

// executeFlowStepsAsync executes flow steps in a goroutine
func (fm *FlowManager) executeFlowStepsAsync(ctx context.Context, execution *internal.Execution, flow *internal.Flow) {
	fmt.Printf("Starting execution of flow %s (ID: %s)\n", flow.Name, execution.ID.Hex())

	stepResults := make(map[string]interface{})

	for _, step := range flow.Steps {
		stepStartTime := time.Now()

		executionStep := internal.ExecutionStep{
			StepID: step.StepID,
			Status: internal.StatusRunning,
		}

		var stepErr error
		switch step.Type {
		case internal.StepTypeProduce:
			stepErr = fm.executeProduceStep(ctx, &step, &executionStep, stepResults)
		case internal.StepTypeConsume:
			stepErr = fm.executeConsumeStep(ctx, &step, &executionStep, stepResults)
		case internal.StepTypeValidate:
			stepErr = fm.executeValidateStep(ctx, &step, &executionStep, stepResults)
		case internal.StepTypeDelay:
			stepErr = fm.executeDelayStep(ctx, &step, &executionStep)
		default:
			stepErr = fmt.Errorf("unknown step type: %s", step.Type)
		}

		// Update step execution result
		stepDuration := time.Since(stepStartTime).Milliseconds()
		executionStep.Duration = stepDuration

		if stepErr != nil {
			executionStep.Status = internal.StatusFailed
			executionStep.Errors = []string{stepErr.Error()}
			execution.Metrics.ErrorsCount++
			fmt.Printf("Step %s failed: %v\n", step.StepID, stepErr)
		} else {
			executionStep.Status = internal.StatusSuccess
			execution.Metrics.StepsCompleted++
		}

		execution.Steps = append(execution.Steps, executionStep)

		// Store step output for reference by subsequent steps
		if executionStep.Output != nil {
			for key, value := range executionStep.Output {
				stepResults[fmt.Sprintf("%s.%s", step.StepID, key)] = value
			}
		}

		// Stop execution if step failed
		if stepErr != nil {
			execution.Status = internal.StatusFailed
			break
		}
	}

	// Update final execution status
	if execution.Status == internal.StatusRunning {
		execution.Status = internal.StatusSuccess
	}

	endTime := time.Now()
	execution.EndTime = &endTime
	execution.Metrics.TotalDuration = endTime.Sub(execution.StartTime).Milliseconds()

	// Update execution in database
	if err := fm.store.UpdateExecution(ctx, execution); err != nil {
		fmt.Printf("Failed to update execution: %v\n", err)
	}

	fmt.Printf("Completed execution of flow %s with status %s\n", flow.Name, execution.Status)
}

// executeProduceStep executes a produce step
func (fm *FlowManager) executeProduceStep(ctx context.Context, step *internal.FlowStep, execStep *internal.ExecutionStep, stepResults map[string]interface{}) error {
	// TODO: Implement message production
	// This would integrate with your Kafka producer
	fmt.Printf("Executing produce step %s to topic %s\n", step.StepID, step.Config.Topic)

	execStep.Input = step.Config.Message
	execStep.Output = map[string]interface{}{
		"topic":   step.Config.Topic,
		"message": "produced",
	}

	return nil
}

// executeConsumeStep executes a consume step
func (fm *FlowManager) executeConsumeStep(ctx context.Context, step *internal.FlowStep, execStep *internal.ExecutionStep, stepResults map[string]interface{}) error {
	// TODO: Implement message consumption
	// This would integrate with your Kafka consumer
	fmt.Printf("Executing consume step %s from topic %s\n", step.StepID, step.Config.Topic)

	timeout := time.Duration(step.Config.Timeout) * time.Millisecond
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	// Simulate message consumption
	select {
	case <-time.After(1 * time.Second):
		execStep.Output = map[string]interface{}{
			"topic":   step.Config.Topic,
			"message": "consumed",
		}
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for message from topic %s", step.Config.Topic)
	}
}

// executeValidateStep executes a validation step
func (fm *FlowManager) executeValidateStep(ctx context.Context, step *internal.FlowStep, execStep *internal.ExecutionStep, stepResults map[string]interface{}) error {
	// TODO: Implement message validation
	fmt.Printf("Executing validate step %s\n", step.StepID)

	execStep.Input = map[string]interface{}{
		"expected": step.Config.ExpectedMessage,
	}
	execStep.Output = map[string]interface{}{
		"valid": true,
	}

	return nil
}

// executeDelayStep executes a delay step
func (fm *FlowManager) executeDelayStep(ctx context.Context, step *internal.FlowStep, execStep *internal.ExecutionStep) error {
	delayMs := step.Config.DelayMs
	if delayMs <= 0 {
		delayMs = 1000 // Default 1 second
	}

	delay := time.Duration(delayMs) * time.Millisecond
	fmt.Printf("Delaying for %v\n", delay)

	execStep.Input = map[string]interface{}{
		"delayMs": delayMs,
	}

	time.Sleep(delay)

	execStep.Output = map[string]interface{}{
		"delayed": delayMs,
	}

	return nil
}

// validateFlowSteps validates the structure of flow steps
func (fm *FlowManager) validateFlowSteps(steps []internal.FlowStep) error {
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
		case internal.StepTypeProduce:
			if step.Config.Topic == "" {
				return fmt.Errorf("produce step %s must have a topic", step.StepID)
			}
			if step.Config.Message == nil {
				return fmt.Errorf("produce step %s must have a message", step.StepID)
			}
		case internal.StepTypeConsume:
			if step.Config.Topic == "" {
				return fmt.Errorf("consume step %s must have a topic", step.StepID)
			}
		case internal.StepTypeValidate:
			if step.Config.ExpectedMessage == nil {
				return fmt.Errorf("validate step %s must have an expected message", step.StepID)
			}
		case internal.StepTypeDelay:
			// Delay steps are always valid
		default:
			return fmt.Errorf("unknown step type: %s", step.Type)
		}
	}

	return nil
}
