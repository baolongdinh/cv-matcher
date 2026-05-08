package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cv-jd-matcher/internal/models"
	"cv-jd-matcher/internal/services"
	"cv-jd-matcher/internal/utils"
	"cv-jd-matcher/internal/workers"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const maxBulkFiles = 20

type BulkHandler struct {
	pdfService     *services.PDFService
	historyService *services.HistoryService
}

func NewBulkHandler(ps *services.PDFService, hs *services.HistoryService) *BulkHandler {
	return &BulkHandler{
		pdfService:     ps,
		historyService: hs,
	}
}

func (h *BulkHandler) AnalyzeBulk(c *gin.Context) {
	sessionID, ok := requireSessionID(c)
	if !ok {
		return
	}

	jobDescription := strings.TrimSpace(c.PostForm("job_description"))
	apiKey := strings.TrimSpace(c.PostForm("api_key"))
	jobID := strings.TrimSpace(c.PostForm("job_id"))

	if jobDescription == "" || apiKey == "" {
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Job description and API key are required", nil)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		writeError(c, http.StatusBadRequest, "MULTIPART_PARSE_ERROR", "Failed to parse multipart form", err.Error())
		return
	}

	files := form.File["cv_files"]
	if len(files) == 0 {
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "At least one CV file is required", nil)
		return
	}
	if len(files) > maxBulkFiles {
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf("A maximum of %d CV files are allowed per batch", maxBulkFiles), nil)
		return
	}

	batchID := fmt.Sprintf("batch_%d", time.Now().UnixNano())
	ctx := c.Request.Context()

	if err := h.historyService.InitBatch(ctx, sessionID, batchID, jobID, len(files)); err != nil {
		writeError(c, http.StatusInternalServerError, "BATCH_INIT_ERROR", "Failed to initialize batch", err.Error())
		return
	}

	for _, fileHeader := range files {
		requestID := fmt.Sprintf("req_%d", time.Now().UnixNano())
		itemStatus := models.BatchItemStatus{
			RequestID:  requestID,
			CVFileName: fileHeader.Filename,
			Status:     "queued",
			Cached:     false,
		}
		if err := h.historyService.SaveBatchItemStatus(ctx, sessionID, batchID, itemStatus); err != nil {
			writeError(c, http.StatusInternalServerError, "BATCH_ITEM_INIT_ERROR", "Failed to initialize batch item", err.Error())
			return
		}

		if strings.ToLower(fileHeader.Header.Get("Content-Type")) != "application/pdf" && !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".pdf") {
			itemStatus.Status = "failed"
			itemStatus.ErrorCode = "INVALID_FILE_TYPE"
			itemStatus.ErrorMessage = "Only PDF files are supported"
			if _, err := h.historyService.UpdateBatchProgress(ctx, sessionID, batchID, itemStatus, true); err != nil {
				writeError(c, http.StatusInternalServerError, "BATCH_PROGRESS_ERROR", "Failed to update batch progress", err.Error())
				return
			}
			_ = h.historyService.SaveBatchFailure(ctx, sessionID, batchID, models.BatchFailure{
				RequestID:    requestID,
				CVFileName:   fileHeader.Filename,
				ErrorCode:    "INVALID_FILE_TYPE",
				ErrorMessage: "Only PDF files are supported",
			})
			continue
		}

		openedFile, err := fileHeader.Open()
		if err != nil {
			itemStatus.Status = "failed"
			itemStatus.ErrorCode = "FILE_READ_ERROR"
			itemStatus.ErrorMessage = "Failed to open uploaded file"
			_, _ = h.historyService.UpdateBatchProgress(ctx, sessionID, batchID, itemStatus, true)
			_ = h.historyService.SaveBatchFailure(ctx, sessionID, batchID, models.BatchFailure{
				RequestID:    requestID,
				CVFileName:   fileHeader.Filename,
				ErrorCode:    "FILE_READ_ERROR",
				ErrorMessage: "Failed to open uploaded file",
				Details:      err.Error(),
			})
			continue
		}

		fileBytes, readErr := io.ReadAll(openedFile)
		openedFile.Close()
		if readErr != nil {
			itemStatus.Status = "failed"
			itemStatus.ErrorCode = "FILE_READ_ERROR"
			itemStatus.ErrorMessage = "Failed to read uploaded file"
			_, _ = h.historyService.UpdateBatchProgress(ctx, sessionID, batchID, itemStatus, true)
			_ = h.historyService.SaveBatchFailure(ctx, sessionID, batchID, models.BatchFailure{
				RequestID:    requestID,
				CVFileName:   fileHeader.Filename,
				ErrorCode:    "FILE_READ_ERROR",
				ErrorMessage: "Failed to read uploaded file",
				Details:      readErr.Error(),
			})
			continue
		}

		fileHash := h.pdfService.GenerateFileHash(fileBytes)
		cachedResult, cacheErr := h.historyService.GetResultByFileHashForSession(ctx, fileHash, sessionID)
		if cacheErr == nil && cachedResult != nil {
			cloned, cloneErr := h.historyService.CloneResultToSession(ctx, cachedResult, sessionID, requestID, jobID, fileHeader.Filename)
			if cloneErr != nil {
				itemStatus.Status = "failed"
				itemStatus.ErrorCode = "CACHE_CLONE_ERROR"
				itemStatus.ErrorMessage = "Failed to copy cached result into the current session"
				_, _ = h.historyService.UpdateBatchProgress(ctx, sessionID, batchID, itemStatus, true)
				_ = h.historyService.SaveBatchFailure(ctx, sessionID, batchID, models.BatchFailure{
					RequestID:    requestID,
					CVFileName:   fileHeader.Filename,
					ErrorCode:    "CACHE_CLONE_ERROR",
					ErrorMessage: "Failed to copy cached result into the current session",
					Details:      cloneErr.Error(),
				})
				continue
			}
			itemStatus.Status = "completed"
			itemStatus.Cached = true
			if err := h.historyService.SaveBatchResult(ctx, sessionID, batchID, requestID, cloned); err != nil {
				writeError(c, http.StatusInternalServerError, "BATCH_RESULT_SAVE_ERROR", "Failed to save cached batch result", err.Error())
				return
			}
			if _, err := h.historyService.UpdateBatchProgress(ctx, sessionID, batchID, itemStatus, true); err != nil {
				writeError(c, http.StatusInternalServerError, "BATCH_PROGRESS_ERROR", "Failed to update batch progress", err.Error())
				return
			}
			continue
		}
		if cacheErr != nil && cacheErr != redis.Nil {
			writeError(c, http.StatusInternalServerError, "CACHE_LOOKUP_ERROR", "Failed to check cached result", cacheErr.Error())
			return
		}

		text, _, err := h.pdfService.ExtractText(fileBytes)
		if err != nil {
			itemStatus.Status = "failed"
			itemStatus.ErrorCode = "PDF_PARSE_ERROR"
			itemStatus.ErrorMessage = "Failed to extract text from PDF"
			_, _ = h.historyService.UpdateBatchProgress(ctx, sessionID, batchID, itemStatus, true)
			_ = h.historyService.SaveBatchFailure(ctx, sessionID, batchID, models.BatchFailure{
				RequestID:    requestID,
				CVFileName:   fileHeader.Filename,
				ErrorCode:    "PDF_PARSE_ERROR",
				ErrorMessage: "Failed to extract text from PDF",
				Details:      err.Error(),
			})
			continue
		}

		task := workers.AnalysisTask{
			BatchID:        batchID,
			RequestID:      requestID,
			SessionID:      sessionID,
			JobID:          jobID,
			JobDescription: jobDescription,
			CVText:         text,
			CVFileName:     fileHeader.Filename,
			APIKey:         apiKey,
			FileHash:       fileHash,
		}

		taskJSON, err := json.Marshal(task)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "TASK_SERIALIZATION_ERROR", "Failed to enqueue analysis task", err.Error())
			return
		}
		if err := utils.RedisClient.RPush(ctx, "cv:queue", taskJSON).Err(); err != nil {
			writeError(c, http.StatusInternalServerError, "QUEUE_PUSH_ERROR", "Failed to enqueue analysis task", err.Error())
			return
		}
	}

	writeSuccess(c, http.StatusAccepted, models.BatchAnalyzeData{
		BatchID:       batchID,
		JobID:         jobID,
		TotalFiles:    len(files),
		Message:       "Analysis batch queued successfully",
		SessionScoped: true,
	})
}

func (h *BulkHandler) GetBatchStatus(c *gin.Context) {
	sessionID, ok := requireSessionID(c)
	if !ok {
		return
	}

	status, err := h.historyService.GetBatchMetaWithSession(c.Request.Context(), sessionID, c.Param("batch_id"))
	if err == redis.Nil {
		writeError(c, http.StatusNotFound, "BATCH_NOT_FOUND", "Batch not found for this session", nil)
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "BATCH_STATUS_ERROR", "Failed to fetch batch status", err.Error())
		return
	}

	writeSuccess(c, http.StatusOK, status)
}

func (h *BulkHandler) GetBatchResults(c *gin.Context) {
	sessionID, ok := requireSessionID(c)
	if !ok {
		return
	}

	results, err := h.historyService.GetBatchResultsWithSession(c.Request.Context(), sessionID, c.Param("batch_id"))
	if err == redis.Nil {
		writeError(c, http.StatusNotFound, "BATCH_NOT_FOUND", "Batch not found for this session", nil)
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "BATCH_RESULTS_ERROR", "Failed to fetch batch results", err.Error())
		return
	}

	writeSuccess(c, http.StatusOK, results)
}

func (h *BulkHandler) GetBatchNotification(c *gin.Context) {
	sessionID, ok := requireSessionID(c)
	if !ok {
		return
	}

	status, err := h.historyService.GetBatchMetaWithSession(c.Request.Context(), sessionID, c.Param("batch_id"))
	if err == redis.Nil {
		writeError(c, http.StatusNotFound, "BATCH_NOT_FOUND", "Batch not found for this session", nil)
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "BATCH_NOTIFICATION_ERROR", "Failed to fetch batch notification", err.Error())
		return
	}

	complete := status.Total > 0 && status.Completed >= status.Total
	notificationType := "info"
	message := fmt.Sprintf("Processing... %d of %d completed", status.Completed, status.Total)
	if complete {
		notificationType = "success"
		message = fmt.Sprintf("Analysis completed successfully! %d candidates processed", status.Successful)
		if status.Failed > 0 {
			notificationType = "warning"
			message = fmt.Sprintf("Analysis completed with %d successful and %d failed results", status.Successful, status.Failed)
		}
	}

	writeSuccess(c, http.StatusOK, models.BatchNotificationData{
		BatchID:    status.BatchID,
		Status:     status.Status,
		Complete:   complete,
		Total:      status.Total,
		Completed:  status.Completed,
		Successful: status.Successful,
		Failed:     status.Failed,
		Message:    message,
		Type:       notificationType,
	})
}
