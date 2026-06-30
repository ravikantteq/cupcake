package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ravikantteq/cupcake/backyard/internal"
	"github.com/ravikantteq/cupcake/backyard/internal/manager"
)

// ProducerHandler handles producer-related HTTP requests
type ProducerHandler struct {
	mgr *manager.ProducerManager
}

// NewProducerHandler creates a new producer handler
func NewProducerHandler(mgr *manager.ProducerManager) *ProducerHandler {
	return &ProducerHandler{mgr: mgr}
}

// PublishMessage handles POST /kafka/publish
func (h *ProducerHandler) PublishMessage(c *gin.Context) {
	var req internal.KafkaMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, internal.APIError{
			Error:   "Validation Error",
			Message: err.Error(),
		})
		return
	}

	err := h.mgr.PublishMessage(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, internal.APIError{
			Error:   "Failed to publish message",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, internal.APIResponse{
		Success: true,
		Message: "Message published successfully",
		Data: map[string]interface{}{
			"topic": req.Topic,
			"key":   req.Key,
		},
	})
}

// GetProducerHistory handles GET /history
func (h *ProducerHandler) GetProducerHistory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	history, err := h.mgr.GetProducerHistory(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, internal.APIError{
			Error:   "Failed to get producer history",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, internal.APIResponse{
		Success: true,
		Message: "Producer history retrieved successfully",
		Data:    history,
	})
}

// GetRecentProducerHistory handles GET /history/recent
func (h *ProducerHandler) GetRecentProducerHistory(c *gin.Context) {
	history, err := h.mgr.GetRecentProducerHistory(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, internal.APIError{
			Error:   "Failed to get recent producer history",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, internal.APIResponse{
		Success: true,
		Message: "Recent producer history retrieved successfully",
		Data:    history,
	})
}

// HealthHandler handles system health checks
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck handles GET /health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, internal.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services: map[string]string{
			"api":      "running",
			"database": "connected",
			"kafka":    "available",
		},
		Version: "2.0.0",
	})
}
