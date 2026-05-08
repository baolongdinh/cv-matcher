package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cv-jd-matcher/internal/handlers"
	"cv-jd-matcher/internal/services"
	"cv-jd-matcher/internal/utils"
	"cv-jd-matcher/internal/workers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	// 1. Initialize Redis
	if err := utils.InitRedis(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	} else {
		log.Printf("Connected to Redis successfully")
	}

	// Set Gin mode
	if os.Getenv("DEBUG") == "false" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 2. Initialize Services
	pdfService := services.NewPDFService()
	historyService := services.NewHistoryService()

	// 3. Initialize Handlers
	analyzeHandler := handlers.NewAnalyzeHandler(pdfService, historyService)
	bulkHandler := handlers.NewBulkHandler(pdfService, historyService)
	historyHandler := handlers.NewHistoryHandler(historyService)

	// 4. Start Background Worker
	worker := workers.NewAnalysisWorker(historyService, pdfService)
	go worker.Start(context.Background())

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Session-ID", "X-File-Hash", "x-goog-api-key"}
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	// API routes group
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": "CV-JD Matcher API is running",
			})
		})

		// Single analysis
		api.POST("/analyze", analyzeHandler.Analyze)

		// Bulk analysis (Queue)
		api.POST("/analyze/bulk", bulkHandler.AnalyzeBulk)
		api.GET("/jobs/:batch_id", bulkHandler.GetBatchStatus)
		api.GET("/jobs/:batch_id/results", bulkHandler.GetBatchResults)
		api.GET("/jobs/:batch_id/notification", bulkHandler.GetBatchNotification)

		// History & Ranking
		api.GET("/history", historyHandler.GetHistory)
		api.DELETE("/history/:id", historyHandler.DeleteHistory)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
