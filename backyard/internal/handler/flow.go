package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ravikantteq/cupcake/backyard/internal"
	"github.com/ravikantteq/cupcake/backyard/internal/manager"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FlowHandler handles flow-related HTTP requests
type FlowHandler struct {
	mgr *manager.FlowManager
}

// NewFlowHandler creates a new flow handler
func NewFlowHandler(mgr *manager.FlowManager) *FlowHandler {
	return &FlowHandler{mgr: mgr}
}

// CreateFlow handles POST /flows
func (h *FlowHandler) CreateFlow(c *gin.Context) {
	var req internal.CreateFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, internal.APIError{
			Error:   "Validation Error",
			Message: err.Error(),
		})
		return
	}

	flow, err := h.mgr.CreateFlow(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, internal.APIError{
			Error:   "Failed to create flow",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, internal.APIResponse{
		Success: true,
		Message: "Flow created successfully",
		Data:    flow,
	})
}

// GetFlows handles GET /flows
func (h *FlowHandler) GetFlows(c *gin.Context) {
	flows, err := h.mgr.GetFlows(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, internal.APIError{
			Error:   "Failed to get flows",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, internal.APIResponse{
		Success: true,
		Message: "Flows retrieved successfully",
		Data:    flows,
	})
}

// GetFlow handles GET /flows/:id
func (h *FlowHandler) GetFlow(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, internal.APIError{
			Error:   "Invalid flow ID",
			Message: err.Error(),
		})
		return
	}

	flow, err := h.mgr.GetFlow(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, internal.APIError{
			Error:   "Flow not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, internal.APIResponse{
		Success: true,
		Message: "Flow retrieved successfully",
		Data:    flow,
	})
}

// UpdateFlow handles PUT /flows/:id
func (h *FlowHandler) UpdateFlow(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, internal.APIError{
			Error:   "Invalid flow ID",
			Message: err.Error(),
		})
		return
	}

	var req internal.CreateFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, internal.APIError{
			Error:   "Validation Error",
			Message: err.Error(),
		})
		return
	}

	flow, err := h.mgr.UpdateFlow(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, internal.APIError{
			Error:   "Failed to update flow",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, internal.APIResponse{
		Success: true,
		Message: "Flow updated successfully",
		Data:    flow,
	})
}

// ExecuteFlow handles POST /flows/:id/execute
func (h *FlowHandler) ExecuteFlow(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, internal.APIError{
			Error:   "Invalid flow ID",
			Message: err.Error(),
		})
		return
	}

	execution, err := h.mgr.ExecuteFlow(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, internal.APIError{
			Error:   "Failed to execute flow",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, internal.APIResponse{
		Success: true,
		Message: "Flow execution started",
		Data:    execution,
	})
}

// DeleteFlow handles DELETE /flows/:id
func (h *FlowHandler) DeleteFlow(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, internal.APIError{
			Error:   "Invalid flow ID",
			Message: err.Error(),
		})
		return
	}

	err = h.mgr.DeleteFlow(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, internal.APIError{
			Error:   "Failed to delete flow",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, internal.APIResponse{
		Success: true,
		Message: "Flow deleted successfully",
	})
}
