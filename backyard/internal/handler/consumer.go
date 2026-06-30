package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ravikantteq/cupcake/backyard/internal"
	"github.com/ravikantteq/cupcake/backyard/internal/manager"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConsumerHandler handles consumer-related HTTP requests
type ConsumerHandler struct {
	mgr *manager.ConsumerManager
}

// NewConsumerHandler creates a new consumer handler
func NewConsumerHandler(mgr *manager.ConsumerManager) *ConsumerHandler {
	return &ConsumerHandler{mgr: mgr}
}

// CreateConsumer handles POST /consumers
func (h *ConsumerHandler) CreateConsumer(c *gin.Context) {
	var req internal.CreateConsumerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, internal.APIError{
			Error:   "Validation Error",
			Message: err.Error(),
		})
		return
	}

	consumer, err := h.mgr.CreateConsumer(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, internal.APIError{
			Error:   "Failed to create consumer",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, internal.APIResponse{
		Success: true,
		Message: "Consumer created successfully",
		Data:    consumer,
	})
}

// GetConsumers handles GET /consumers
func (h *ConsumerHandler) GetConsumers(c *gin.Context) {
	consumers, err := h.mgr.GetConsumers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, internal.APIError{
			Error:   "Failed to get consumers",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, internal.APIResponse{
		Success: true,
		Message: "Consumers retrieved successfully",
		Data:    consumers,
	})
}

// GetConsumer handles GET /consumers/:id
func (h *ConsumerHandler) GetConsumer(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, internal.APIError{
			Error:   "Invalid consumer ID",
			Message: err.Error(),
		})
		return
	}

	consumer, err := h.mgr.GetConsumer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, internal.APIError{
			Error:   "Consumer not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, internal.APIResponse{
		Success: true,
		Message: "Consumer retrieved successfully",
		Data:    consumer,
	})
}

// StartConsumer handles POST /consumers/:id/start
func (h *ConsumerHandler) StartConsumer(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, internal.APIError{
			Error:   "Invalid consumer ID",
			Message: err.Error(),
		})
		return
	}

	err = h.mgr.StartConsumer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, internal.APIError{
			Error:   "Failed to start consumer",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, internal.APIResponse{
		Success: true,
		Message: "Consumer started successfully",
	})
}

// StopConsumer handles POST /consumers/:id/stop
func (h *ConsumerHandler) StopConsumer(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, internal.APIError{
			Error:   "Invalid consumer ID",
			Message: err.Error(),
		})
		return
	}

	err = h.mgr.StopConsumer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, internal.APIError{
			Error:   "Failed to stop consumer",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, internal.APIResponse{
		Success: true,
		Message: "Consumer stopped successfully",
	})
}

// DeleteConsumer handles DELETE /consumers/:id
func (h *ConsumerHandler) DeleteConsumer(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, internal.APIError{
			Error:   "Invalid consumer ID",
			Message: err.Error(),
		})
		return
	}

	err = h.mgr.DeleteConsumer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, internal.APIError{
			Error:   "Failed to delete consumer",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, internal.APIResponse{
		Success: true,
		Message: "Consumer deleted successfully",
	})
}
