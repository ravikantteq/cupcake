package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/ravikantteq/cupcake/backyard/docs" // This line is required for swagger
	"github.com/ravikantteq/cupcake/backyard/internal/api"
	"github.com/ravikantteq/cupcake/backyard/internal/repository"
	"github.com/ravikantteq/cupcake/backyard/internal/services"
	"github.com/ravikantteq/cupcake/backyard/pkg/storage"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Cupcake Kafka Test Framework API
// @version 2.0
// @description Enterprise-ready Kafka testing platform with advanced flow design and intelligent message matching
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/ravikantteq/cupcake
// @contact.email support@cupcake.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
func main() {
	// Get environment variables
	mongoURI := getEnv("MONGO_URI", "mongodb://cupcake:cupcake123@localhost:27017/cupcake?authSource=admin")
	kafkaBroker := getEnv("KAFKA_BROKER", "localhost:9093")
	port := getEnv("PORT", "8080")
	ginMode := getEnv("GIN_MODE", "debug")

	// Set Gin mode
	gin.SetMode(ginMode)

	// Initialize MongoDB
	db, err := storage.NewMongoDB(mongoURI, "cupcake")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}()

	// Initialize repositories and services
	repo := repository.NewRepository(db)
	flowService := services.NewFlowService(repo, kafkaBroker)

	// Initialize handlers
	handlers := api.NewHandlers(flowService, db)

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

	// Health check endpoint
	r.GET("/health", handlers.HealthCheck)

	// Legacy API routes (for backward compatibility)
	apiGroup := r.Group("/api")
	{
		kafkaGroup := apiGroup.Group("/kafka")
		{
			kafkaGroup.POST("/publish", handlers.PublishMessage)
		}
	}

	// New API v1 routes
	v1Group := r.Group("/api/v1")
	{
		// Flow management
		flowsGroup := v1Group.Group("/flows")
		{
			flowsGroup.POST("", handlers.CreateFlow)
			flowsGroup.GET("", handlers.GetFlows)
			flowsGroup.GET("/:id", handlers.GetFlowByID)
			flowsGroup.PUT("/:id", handlers.UpdateFlow)
			flowsGroup.POST("/:id/execute", handlers.ExecuteFlow)
		}

		// Producer history management
		historyGroup := v1Group.Group("/history")
		{
			historyGroup.GET("", handlers.GetProducerHistory)
			historyGroup.GET("/recent", handlers.GetRecentProducerHistory)
		}

		// TODO: Add more endpoints for suites, consumers, executions, etc.
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("🧁 Starting Cupcake Kafka Test Framework API server")
	log.Printf("📊 MongoDB connected: %s", mongoURI)
	log.Printf("📨 Kafka broker: %s", kafkaBroker)
	log.Printf("🌐 Server listening on port: %s", port)
	log.Printf("📖 Swagger UI: http://localhost:%s/swagger/index.html", port)
	log.Printf("🔧 Health Check: http://localhost:%s/health", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
