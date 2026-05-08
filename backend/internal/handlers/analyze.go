package handlers

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"cv-jd-matcher/internal/models"
	"cv-jd-matcher/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AnalyzeHandler struct {
	pdfService     *services.PDFService
	historyService *services.HistoryService
}

func NewAnalyzeHandler(pdfService *services.PDFService, historyService *services.HistoryService) *AnalyzeHandler {
	return &AnalyzeHandler{
		pdfService:     pdfService,
		historyService: historyService,
	}
}

func (h *AnalyzeHandler) Analyze(c *gin.Context) {
	sessionID, ok := requireSessionID(c)
	if !ok {
		return
	}

	startTime := time.Now()
	jobDescription := c.PostForm("job_description")
	apiKey := c.PostForm("api_key")
	jobID := c.PostForm("job_id")
	file, err := c.FormFile("cv_file")

	if jobDescription == "" || file == nil || apiKey == "" {
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Job description, API key and CV file are required", nil)
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		writeError(c, http.StatusInternalServerError, "FILE_READ_ERROR", "Failed to read uploaded file", err.Error())
		return
	}
	defer openedFile.Close()

	fileBytes, err := io.ReadAll(openedFile)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "FILE_READ_ERROR", "Failed to read file content", err.Error())
		return
	}

	fileHash := h.pdfService.GenerateFileHash(fileBytes)
	ctx := c.Request.Context()

	cachedResult, err := h.historyService.GetResultByFileHashForSession(ctx, fileHash, sessionID)
	if err == nil && cachedResult != nil {
		requestID := fmt.Sprintf("req_%d", time.Now().UnixNano())
		cloned, cloneErr := h.historyService.CloneResultToSession(ctx, cachedResult, sessionID, requestID, jobID, file.Filename)
		if cloneErr != nil {
			writeError(c, http.StatusInternalServerError, "CACHE_CLONE_ERROR", "Failed to store cached result for this session", cloneErr.Error())
			return
		}
		writeSuccess(c, http.StatusOK, cloned)
		return
	}
	if err != nil && err != redis.Nil {
		writeError(c, http.StatusInternalServerError, "CACHE_LOOKUP_ERROR", "Failed to check cached result", err.Error())
		return
	}

	cvText, numPages, err := h.pdfService.ExtractText(fileBytes)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "PDF_PARSE_ERROR", "Failed to extract text from PDF", err.Error())
		return
	}

	geminiSvc, err := services.NewGeminiService(apiKey)
	if err != nil {
		writeError(c, http.StatusUnauthorized, "INVALID_API_KEY", "Failed to initialize Gemini service", err.Error())
		return
	}

	result, err := geminiSvc.AnalyzeCVJD(ctx, jobDescription, cvText)
	if err != nil {
		writeError(c, http.StatusServiceUnavailable, "AI_SERVICE_ERROR", "Failed to analyze CV with AI", err.Error())
		return
	}

	requestID := fmt.Sprintf("req_%d", time.Now().UnixNano())
	result.ProcessingMetadata = models.ProcessingMetadata{
		ProcessingTime:   time.Since(startTime).Seconds(),
		CVPagesProcessed: numPages,
		JDWordCount:      len(jobDescription),
		ModelUsed:        "gemini-2.5-flash-lite",
		Timestamp:        time.Now().Format(time.RFC3339),
		RequestID:        requestID,
		CVFileName:       file.Filename,
		JobID:            jobID,
		Status:           "completed",
		FileHash:         fileHash,
	}
	result.FileHash = fileHash

	if err := h.historyService.SaveResultWithSession(ctx, result, sessionID); err != nil {
		writeError(c, http.StatusInternalServerError, "HISTORY_SAVE_ERROR", "Failed to save analysis result", err.Error())
		return
	}

	writeSuccess(c, http.StatusOK, result)
}
