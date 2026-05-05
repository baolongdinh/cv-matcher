package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cv-jd-matcher/internal/models"
	"cv-jd-matcher/internal/services"

	"github.com/gin-gonic/gin"
)

type AnalyzeHandler struct {
	pdfService     *services.PDFService
	geminiService  *services.GeminiService
	historyService *services.HistoryService
}

func NewAnalyzeHandler(pdfService *services.PDFService, historyService *services.HistoryService) *AnalyzeHandler {
	return &AnalyzeHandler{
		pdfService:     pdfService,
		historyService: historyService,
	}
}

func (h *AnalyzeHandler) Analyze(c *gin.Context) {
	startTime := time.Now()

	// 1. Parse multipart form
	jobDescription := c.PostForm("job_description")
	apiKey := c.PostForm("api_key")
	file, err := c.FormFile("cv_file")
	sessionID := c.GetHeader("X-Session-ID")

	if jobDescription == "" || file == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Job description and CV file are required",
			},
		})
		return
	}

	// 2. Read CV file
	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "FILE_READ_ERROR",
				"message": "Failed to read uploaded file",
			},
		})
		return
	}
	defer openedFile.Close()

	fileBytes := make([]byte, file.Size)
	_, err = openedFile.Read(fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "FILE_READ_ERROR",
				"message": "Failed to read file content",
			},
		})
		return
	}

	// 3. Extract text from PDF
	cvText, numPages, err := h.pdfService.ExtractText(fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "PDF_PARSE_ERROR",
				"message": "Failed to extract text from PDF",
				"details": err.Error(),
			},
		})
		return
	}

	// 4. Initialize Gemini service (per request or global)
	geminiSvc, err := services.NewGeminiService(apiKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "INVALID_API_KEY",
				"message": "Failed to initialize Gemini service",
				"details": err.Error(),
			},
		})
		return
	}

	// 5. Analyze with Gemini
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := geminiSvc.AnalyzeCVJD(ctx, jobDescription, cvText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "ANALYSIS_ERROR",
				"message": "Failed to analyze CV-JD",
				"details": err.Error(),
			},
		})
		return
	}

	// 6. Add metadata
	result.ProcessingMetadata = models.ProcessingMetadata{
		ProcessingTime:   time.Since(startTime).Seconds(),
		CVPagesProcessed: numPages,
		JDWordCount:      len(jobDescription), // Simple word count
		ModelUsed:        "gemini-2.5-flash-lite",
		Timestamp:        time.Now().Format(time.RFC3339),
		RequestID:        fmt.Sprintf("req_%d", time.Now().UnixNano()),
		CVFileName:       file.Filename,
	}

	// 7. Save to Redis History asynchronously if Session ID is present
	if sessionID != "" {
		go func(sid string, res *models.AnalysisResult) {
			_ = h.historyService.SaveAnalysis(context.Background(), sid, res)
		}(sessionID, result)
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}
