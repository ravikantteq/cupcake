package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/ravikantteq/cupcake/backyard/docs"
	"github.com/ravikantteq/cupcake/backyard/internal/handler"
	"github.com/ravikantteq/cupcake/backyard/internal/manager"
	"github.com/ravikantteq/cupcake/backyard/internal/store"
	"github.com/ravikantteq/cupcake/backyard/pkg/storage"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Cupcake Kafka Test Framework API
// @version 2.0
// @description Enterprise-ready Kafka testing platform with clean Go architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/ravikantteq/cupcake
// @contact.email support@cupcake.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
func main() {
	// Configuration
	config := &Config{
		MongoURI:    getEnv("MONGO_URI", "mongodb://cupcake:cupcake123@localhost:27017/cupcake?authSource=admin"),
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9093"),
		Port:        getEnv("PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),
	}

	// Set Gin mode
	gin.SetMode(config.GinMode)

	// Initialize database
	db, err := storage.NewMongoDB(config.MongoURI, "cupcake")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}()

	// Initialize layers
	dataStore := store.NewMongoDB(db)
	mgrs := manager.NewManagers(dataStore, config.KafkaBroker)
	handlers := handler.NewHandlers(mgrs)

	// Setup router
	router := setupRouter(handlers)

	// Start server
	log.Println("🧁 Starting Cupcake Kafka Test Framework API server")
	log.Printf("📊 MongoDB connected: %s", config.MongoURI)
	log.Printf("📨 Kafka broker: %s", config.KafkaBroker)
	log.Printf("🌐 Server listening on port: %s", config.Port)
	log.Printf("📖 Swagger UI: http://localhost:%s/swagger/index.html", config.Port)
	log.Printf("🔧 Health Check: http://localhost:%s/health", config.Port)

	if err := router.Run(":" + config.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Config holds application configuration
type Config struct {
	MongoURI    string
	KafkaBroker string
	Port        string
	GinMode     string
}

// setupRouter configures the HTTP router
func setupRouter(h *handler.Handlers) *gin.Engine {
	r := gin.Default()

	// CORS middleware
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

	// Health check
	r.GET("/health", h.Health.HealthCheck)

	// Legacy API routes (backward compatibility)
	apiGroup := r.Group("/api")
	{
		kafkaGroup := apiGroup.Group("/kafka")
		{
			kafkaGroup.POST("/publish", h.Producer.PublishMessage)
		}
	}

	// API v1 routes (new clean structure)
	v1Group := r.Group("/api/v1")
	{
		// Consumer management
		consumersGroup := v1Group.Group("/consumers")
		{
			consumersGroup.POST("", h.Consumer.CreateConsumer)
			consumersGroup.GET("", h.Consumer.GetConsumers)
			consumersGroup.GET("/:id", h.Consumer.GetConsumer)
			consumersGroup.POST("/:id/start", h.Consumer.StartConsumer)
			consumersGroup.POST("/:id/stop", h.Consumer.StopConsumer)
			consumersGroup.DELETE("/:id", h.Consumer.DeleteConsumer)
		}

		// Flow management
		flowsGroup := v1Group.Group("/flows")
		{
			flowsGroup.POST("", h.Flow.CreateFlow)
			flowsGroup.GET("", h.Flow.GetFlows)
			flowsGroup.GET("/:id", h.Flow.GetFlow)
			flowsGroup.PUT("/:id", h.Flow.UpdateFlow)
			flowsGroup.POST("/:id/execute", h.Flow.ExecuteFlow)
			flowsGroup.GET("/:id/executions", h.Flow.GetExecutions)
			flowsGroup.DELETE("/:id", h.Flow.DeleteFlow)
		}

		// Execution management
		executionsGroup := v1Group.Group("/executions")
		{
			executionsGroup.GET("/:id", h.Flow.GetExecution)
			executionsGroup.GET("", h.Flow.GetExecutions)
		}

		// Producer history
		historyGroup := v1Group.Group("/history")
		{
			historyGroup.GET("", h.Producer.GetProducerHistory)
			historyGroup.GET("/recent", h.Producer.GetRecentProducerHistory)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
