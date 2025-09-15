package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"local.dev/doc-analyzer/internal/pkg/gateway/clients"
	_ "local.dev/doc-analyzer/internal/pkg/gateway/docs"
	"local.dev/doc-analyzer/internal/pkg/gateway/handlers"
)

// @title File Processing API
// @version 1.0
// @description API Gateway for file storing and analysis services
// @BasePath /

func main() {
	log.Println("Starting API Gateway...")

	// Initialize File Storing Service client
	fileStoringAddress := getEnvOrDefault("FILE_STORING_SERVICE_ADDRESS", "file-storing-service:50051")
	fileStoringClient, err := clients.NewFileStoringClient(fileStoringAddress)
	if err != nil {
		log.Fatalf("Failed to initialize File Storing Service client: %v", err)
	}
	defer fileStoringClient.Close()

	// Initialize File Analysis Service client
	fileAnalysisAddress := getEnvOrDefault("FILE_ANALYSIS_SERVICE_ADDRESS", "file-analysis-service:50052")
	fileAnalysisClient, err := clients.NewFileAnalysisClient(fileAnalysisAddress)
	if err != nil {
		log.Fatalf("Failed to initialize File Analysis Service client: %v", err)
	}
	defer fileAnalysisClient.Close()

	// Create Gin router
	router := gin.Default()

	// Initialize handlers
	fileHandler := handlers.NewFileHandler(fileStoringClient)
	analysisHandler := handlers.NewAnalysisHandler(fileAnalysisClient)

	// Setup API routes
	v1 := router.Group("/api/v1")
	{
		// File routes
		v1.POST("/files", fileHandler.UploadFile)
		v1.GET("/files/:file_id", fileHandler.GetFile)

		// Analysis routes
		v1.POST("/analysis", analysisHandler.AnalyzeFile)
		v1.GET("/wordcloud/:location", analysisHandler.GetWordCloud)
	}

	// Setup Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start HTTP server
	httpPort := getEnvOrDefault("HTTP_PORT", "8080")
	log.Printf("Starting HTTP server on port %s...", httpPort)
	if err := router.Run(":" + httpPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnvOrDefault returns the value of the environment variable or the default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("%s not set, using default: %s", key, defaultValue)
		return defaultValue
	}
	return value
}
