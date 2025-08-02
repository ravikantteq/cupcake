package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ravikantteq/cupcake/backyard/internal/models"
	"github.com/ravikantteq/cupcake/backyard/internal/repository"
	"github.com/ravikantteq/cupcake/backyard/internal/services"
	"github.com/ravikantteq/cupcake/backyard/pkg/netw"
	"github.com/ravikantteq/cupcake/backyard/pkg/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handlers struct {
	flowService     *services.FlowService
	consumerService *services.ConsumerService
	db              *storage.MongoDB
	repository      *repository.Repository
	historyRepo     *repository.ProducerHistoryRepository
}

func NewHandlers(flowService *services.FlowService, consumerService *services.ConsumerService, db *storage.MongoDB) *Handlers {
	repo := repository.NewRepository(db)
	return &Handlers{
		flowService:     flowService,
		consumerService: consumerService,
		db:              db,
		repository:      repo,
		historyRepo:     repo.NewProducerHistoryRepository(),
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

	// Create history entry
	history := &models.ProducerHistory{
		Broker:    message.Broker,
		Topic:     message.Topic,
		Key:       message.Key,
		Value:     message.Value,
		Success:   false, // Will be updated after publish attempt
		Timestamp: time.Now(),
	}

	// Publish message
	err := producer.ProduceJSON(message.Key, message.Value)

	// Prepare response data
	responseData := map[string]interface{}{
		"topic": message.Topic,
		"key":   message.Key,
	}

	if err != nil {
		// Store failed attempt in history
		history.Success = false
		history.Error = err.Error()

		// Save to database (don't fail if history save fails)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, histErr := h.historyRepo.CreateHistoryEntry(ctx, history); histErr != nil {
			// Log error but don't fail the request
			// Could use a proper logger here
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Kafka Error",
			Message: err.Error(),
		})
		return
	}

	// Store successful attempt in history
	history.Success = true
	history.Response = responseData

	// Save to database (don't fail if history save fails)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, histErr := h.historyRepo.CreateHistoryEntry(ctx, history); histErr != nil {
		// Log error but don't fail the request
		// Could use a proper logger here
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Message published successfully",
		Data:    responseData,
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

// GetProducerHistory gets recent producer history
// @Summary Get producer history
// @Description Get recent producer message history with pagination
// @Tags history
// @Produce json
// @Param limit query int false "Number of records to return (default 10, max 50)"
// @Param offset query int false "Number of records to skip (default 0)"
// @Success 200 {object} models.Response{data=[]models.ProducerHistory}
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/history [get]
func (h *Handlers) GetProducerHistory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	history, err := h.historyRepo.GetHistoryByUser(ctx, "", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database Error",
			Message: err.Error(),
		})
		return
	}

	// Get total count for pagination info
	totalCount, err := h.historyRepo.GetHistoryCount(ctx, "")
	if err != nil {
		// If count fails, just return the results without count
		totalCount = 0
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Producer history retrieved successfully",
		Data: map[string]interface{}{
			"history":    history,
			"totalCount": totalCount,
			"limit":      limit,
			"offset":     offset,
		},
	})
}

// GetRecentProducerHistory gets the most recent producer history (for UI caching)
// @Summary Get recent producer history
// @Description Get the most recent 10 producer messages for UI caching
// @Tags history
// @Produce json
// @Success 200 {object} models.Response{data=[]models.ProducerHistory}
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/history/recent [get]
func (h *Handlers) GetRecentProducerHistory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	history, err := h.historyRepo.GetRecentHistory(ctx, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Database Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Recent producer history retrieved successfully",
		Data:    history,
	})
}

// CreateConsumer creates a new consumer configuration
// @Summary Create a new consumer
// @Description Create a new consumer configuration that can be started later
// @Tags consumers
// @Accept json
// @Produce json
// @Param consumer body models.CreateConsumerRequest true "Consumer to create"
// @Success 201 {object} models.Response{data=models.Consumer}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/consumers [post]
func (h *Handlers) CreateConsumer(c *gin.Context) {
	var req models.CreateConsumerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid JSON",
			Message: err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	consumer, err := h.consumerService.CreateConsumer(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create consumer",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Consumer created successfully",
		Data:    consumer,
	})
}

// GetConsumers retrieves all consumers
// @Summary Get all consumers
// @Description Get a list of all consumer configurations
// @Tags consumers
// @Produce json
// @Success 200 {object} models.Response{data=[]models.Consumer}
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/consumers [get]
func (h *Handlers) GetConsumers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	consumers, err := h.consumerService.GetAllConsumers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to retrieve consumers",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Consumers retrieved successfully",
		Data:    consumers,
	})
}

// GetConsumerByID retrieves a consumer by ID
// @Summary Get consumer by ID
// @Description Get a specific consumer configuration by ID
// @Tags consumers
// @Produce json
// @Param id path string true "Consumer ID"
// @Success 200 {object} models.Response{data=models.Consumer}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/consumers/{id} [get]
func (h *Handlers) GetConsumerByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Invalid consumer ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	consumer, err := h.consumerService.GetConsumerStatus(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Consumer not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Consumer retrieved successfully",
		Data:    consumer,
	})
}

// StartConsumer starts a consumer
// @Summary Start a consumer
// @Description Start a consumer to begin listening for messages
// @Tags consumers
// @Accept json
// @Produce json
// @Param id path string true "Consumer ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/consumers/{id}/start [post]
func (h *Handlers) StartConsumer(c *gin.Context) {
	idStr := c.Param("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Invalid consumer ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = h.consumerService.StartConsumer(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to start consumer",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Consumer started successfully",
	})
}

// StopConsumer stops a consumer
// @Summary Stop a consumer
// @Description Stop a running consumer
// @Tags consumers
// @Accept json
// @Produce json
// @Param id path string true "Consumer ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/consumers/{id}/stop [post]
func (h *Handlers) StopConsumer(c *gin.Context) {
	idStr := c.Param("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Invalid consumer ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = h.consumerService.StopConsumer(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to stop consumer",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Consumer stopped successfully",
	})
}

// DeleteConsumer deletes a consumer
// @Summary Delete a consumer
// @Description Delete a consumer configuration (only if not running)
// @Tags consumers
// @Accept json
// @Produce json
// @Param id path string true "Consumer ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/consumers/{id} [delete]
func (h *Handlers) DeleteConsumer(c *gin.Context) {
	idStr := c.Param("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid ID",
			Message: "Invalid consumer ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = h.consumerService.DeleteConsumer(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete consumer",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Consumer deleted successfully",
	})
}
