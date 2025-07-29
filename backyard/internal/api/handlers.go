package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ravikantteq/cupcake/backyard/internal/models"
	"github.com/ravikantteq/cupcake/backyard/internal/services"
	"github.com/ravikantteq/cupcake/backyard/pkg/netw"
	"github.com/ravikantteq/cupcake/backyard/pkg/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handlers struct {
	flowService *services.FlowService
	db          *storage.MongoDB
}

func NewHandlers(flowService *services.FlowService, db *storage.MongoDB) *Handlers {
	return &Handlers{
		flowService: flowService,
		db:          db,
	}
}

// PublishMessage publishes a message to Kafka (legacy handler)
// @Summary Publish message to Kafka topic
// @Description Publish a JSON message to specified Kafka topic
// @Tags kafka
// @Accept json
// @Produce json
// @Param message body models.KafkaMessage true "Message to publish"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/kafka/publish [post]
func (h *Handlers) PublishMessage(c *gin.Context) {
	var message models.KafkaMessage

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid JSON",
			Message: err.Error(),
		})
		return
	}

	// Validate required fields
	if message.Topic == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation Error",
			Message: "Topic is required",
		})
		return
	}

	if message.Value == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation Error",
			Message: "Value is required",
		})
		return
	}

	// Create producer for the specific topic
	producer := netw.NewKafkaProducer(message.Broker, message.Topic)

	// Publish message
	err := producer.ProduceJSON(message.Key, message.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Kafka Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Message published successfully",
		Data: map[string]interface{}{
			"topic": message.Topic,
			"key":   message.Key,
		},
	})
}

// CreateFlow creates a new test flow
// @Summary Create a new test flow
// @Description Create a new test flow with steps
// @Tags flows
// @Accept json
// @Produce json
// @Param flow body models.CreateFlowRequest true "Flow to create"
// @Success 201 {object} models.Response{data=models.TestFlow}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/flows [post]
func (h *Handlers) CreateFlow(c *gin.Context) {
	var req models.CreateFlowRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid JSON",
			Message: err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	flow, err := h.flowService.CreateFlow(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create flow",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Flow created successfully",
		Data:    flow,
	})
}

// GetFlows retrieves all test flows
// @Summary Get all test flows
// @Description Retrieve all test flows
// @Tags flows
// @Produce json
// @Success 200 {object} models.Response{data=[]models.TestFlow}
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/flows [get]
func (h *Handlers) GetFlows(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	flows, err := h.flowService.GetAllFlows(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get flows",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Flows retrieved successfully",
		Data:    flows,
	})
}

// GetFlowByID retrieves a specific test flow
// @Summary Get test flow by ID
// @Description Retrieve a specific test flow by ID
// @Tags flows
// @Produce json
// @Param id path string true "Flow ID"
// @Success 200 {object} models.Response{data=models.TestFlow}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/flows/{id} [get]
func (h *Handlers) GetFlowByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Invalid flow ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	flow, err := h.flowService.GetFlowByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Flow not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Flow retrieved successfully",
		Data:    flow,
	})
}

// UpdateFlow updates an existing test flow
// @Summary Update a test flow
// @Description Update an existing test flow by ID
// @Tags flows
// @Accept json
// @Produce json
// @Param id path string true "Flow ID"
// @Param flow body models.CreateFlowRequest true "Flow data"
// @Success 200 {object} models.Response{data=models.TestFlow}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/flows/{id} [put]
func (h *Handlers) UpdateFlow(c *gin.Context) {
	idStr := c.Param("id")

	flowID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Invalid flow ID format",
		})
		return
	}

	var req models.CreateFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid JSON",
			Message: err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	flow, err := h.flowService.UpdateFlow(ctx, flowID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update flow",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Flow updated successfully",
		Data:    flow,
	})
}

// ExecuteFlow executes a test flow
// @Summary Execute a test flow
// @Description Execute a test flow by ID
// @Tags flows
// @Accept json
// @Produce json
// @Param id path string true "Flow ID"
// @Success 202 {object} models.Response{data=models.TestExecution}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/flows/{id}/execute [post]
func (h *Handlers) ExecuteFlow(c *gin.Context) {
	idStr := c.Param("id")

	flowID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Invalid flow ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use a dummy suite ID for now
	suiteID := primitive.NewObjectID()

	execution, err := h.flowService.ExecuteFlow(ctx, flowID, suiteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to execute flow",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, models.Response{
		Success: true,
		Message: "Flow execution started",
		Data:    execution,
	})
}

// HealthCheck checks the health of the service
// @Summary Health check
// @Description Get the health status of the service and its dependencies
// @Tags health
// @Produce json
// @Success 200 {object} models.Response{data=models.HealthStatus}
// @Failure 500 {object} models.ErrorResponse
// @Router /health [get]
func (h *Handlers) HealthCheck(c *gin.Context) {
	status := models.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
		Version:   "2.0.0",
	}

	// Check MongoDB connection
	if err := h.db.Health(); err != nil {
		status.Services["mongodb"] = "unhealthy: " + err.Error()
		status.Status = "unhealthy"
	} else {
		status.Services["mongodb"] = "healthy"
	}

	// Check Kafka connection (basic test)
	status.Services["kafka"] = "healthy" // TODO: Implement actual Kafka health check

	if status.Status == "unhealthy" {
		c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
			Error:   "Service Unhealthy",
			Message: "One or more dependencies are unhealthy",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Service is healthy",
		Data:    status,
	})
}
