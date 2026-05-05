package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"cv-jd-matcher/internal/handlers"
	"cv-jd-matcher/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	// Set Gin mode
	if os.Getenv("DEBUG") == "false" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Initialize Redis
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	// Test Redis connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis at %s: %v", redisURL, err)
	} else {
		log.Printf("Connected to Redis at %s", redisURL)
	}

	// Initialize services
	pdfService := services.NewPDFService()
	historyService := services.NewHistoryService(rdb)
	analyzeHandler := handlers.NewAnalyzeHandler(pdfService, historyService)
	historyHandler := handlers.NewHistoryHandler(historyService)

	// CORS middleware
	config := cors.DefaultConfig()
	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = strings.Split(corsOrigins, ",")
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Session-ID", "X-File-Hash"}
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

		api.POST("/analyze", analyzeHandler.Analyze)
		api.GET("/history", historyHandler.GetHistory)
		api.DELETE("/history/:id", historyHandler.DeleteHistory)
	}

	// Serve Static Frontend (Placeholder)
	// After building the frontend, the files will be in web/dist
	// r.StaticFS("/", http.Dir("./web/dist"))
	// r.NoRoute(func(c *gin.Context) {
	// 	c.File("./web/dist/index.html")
	// })

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
