package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/systemgenes/cupcake/backyard/docs" // This line is required for swagger
	"github.com/systemgenes/cupcake/backyard/internal/api"
)

// @title Cupcake Kafka API
// @version 1.0
// @description A simple Kafka producer API for testing and development
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	// Set up Gin router
	r := gin.Default()

	// Enable CORS for the Angular frontend
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	kafkaHandler := api.NewKafkaHandler("localhost:9092")

	// Health check endpoint
	r.GET("/health", kafkaHandler.HealthCheck)

	// API routes
	apiGroup := r.Group("/api")
	{
		kafkaGroup := apiGroup.Group("/kafka")
		{
			kafkaGroup.POST("/publish", kafkaHandler.PublishMessage)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Starting Cupcake Kafka API server on :8080")
	log.Println("Swagger UI available at: http://localhost:8080/swagger/index.html")

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
