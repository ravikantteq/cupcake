package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemgenes/cupcake/backyard/internal/models"
	"github.com/systemgenes/cupcake/backyard/pkg/netw"
)

type KafkaHandler struct {
	producer *netw.KafkaProducer
}

func NewKafkaHandler(broker string) *KafkaHandler {
	return &KafkaHandler{}
}

// PublishMessage publishes a message to Kafka
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
func (h *KafkaHandler) PublishMessage(c *gin.Context) {
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

// HealthCheck checks the health of the service
// @Summary Health check
// @Description Get the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} models.Response
// @Router /health [get]
func (h *KafkaHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Service is healthy",
		Data:    map[string]interface{}{"status": "ok"},
	})
}
