package workers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"cv-jd-matcher/internal/models"
	"cv-jd-matcher/internal/services"
	"cv-jd-matcher/internal/utils"

	"github.com/redis/go-redis/v9"
)

type AnalysisTask struct {
	BatchID        string `json:"batch_id"`
	RequestID      string `json:"request_id"`
	SessionID      string `json:"session_id"`
	JobID          string `json:"job_id"`
	JobDescription string `json:"job_description"`
	CVText         string `json:"cv_text"`
	CVFileName     string `json:"cv_file_name"`
	APIKey         string `json:"api_key"`
	FileHash       string `json:"file_hash"`
}

type AnalysisWorker struct {
	historyService *services.HistoryService
	pdfService     *services.PDFService
	redis          *redis.Client
}

func NewAnalysisWorker(hs *services.HistoryService, ps *services.PDFService) *AnalysisWorker {
	return &AnalysisWorker{
		historyService: hs,
		pdfService:     ps,
		redis:          utils.RedisClient,
	}
}

func (w *AnalysisWorker) Start(ctx context.Context) {
	log.Println("[Worker] Analysis worker started...")
	for {
		select {
		case <-ctx.Done():
			log.Println("[Worker] Analysis worker stopping...")
			return
		default:
			result, err := w.redis.BLPop(ctx, 5*time.Second, "cv:queue").Result()
			if err != nil {
				if err != redis.Nil {
					log.Printf("[Worker] Error popping task: %v", err)
				}
				continue
			}

			var task AnalysisTask
			if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
				log.Printf("[Worker] Error unmarshaling task: %v", err)
				continue
			}

			w.processTask(ctx, task)
		}
	}
}

func (w *AnalysisWorker) processTask(ctx context.Context, task AnalysisTask) {
	log.Printf("[Worker] Processing CV: %s (Batch: %s)", task.CVFileName, task.BatchID)

	_, err := w.historyService.UpdateBatchProgress(ctx, task.SessionID, task.BatchID, models.BatchItemStatus{
		RequestID:  task.RequestID,
		CVFileName: task.CVFileName,
		Status:     "processing",
	}, false)
	if err != nil {
		log.Printf("[Worker] Error setting processing status for %s: %v", task.RequestID, err)
		return
	}

	geminiSvc, err := services.NewGeminiService(task.APIKey)
	if err != nil {
		w.failTask(ctx, task, "INVALID_API_KEY", "Failed to initialize Gemini service", err.Error())
		return
	}

	startTime := time.Now()
	analysisResult, err := geminiSvc.AnalyzeCVJD(ctx, task.JobDescription, task.CVText)
	if err != nil {
		w.failTask(ctx, task, "ANALYSIS_ERROR", "Failed to analyze CV with AI", err.Error())
		return
	}

	analysisResult.FileHash = task.FileHash
	analysisResult.ProcessingMetadata = models.ProcessingMetadata{
		ProcessingTime:   time.Since(startTime).Seconds(),
		CVPagesProcessed: 0,
		JDWordCount:      len(task.JobDescription),
		ModelUsed:        "gemini-2.5-flash-lite",
		Timestamp:        time.Now().Format(time.RFC3339),
		RequestID:        task.RequestID,
		CVFileName:       task.CVFileName,
		JobID:            task.JobID,
		Status:           "completed",
		FileHash:         task.FileHash,
	}

	if err := w.historyService.SaveResultWithSession(ctx, analysisResult, task.SessionID); err != nil {
		w.failTask(ctx, task, "HISTORY_SAVE_ERROR", "Failed to save analysis result", err.Error())
		return
	}

	if err := w.historyService.SaveBatchResult(ctx, task.SessionID, task.BatchID, task.RequestID, analysisResult); err != nil {
		w.failTask(ctx, task, "BATCH_RESULT_SAVE_ERROR", "Failed to save batch result", err.Error())
		return
	}

	if _, err := w.historyService.UpdateBatchProgress(ctx, task.SessionID, task.BatchID, models.BatchItemStatus{
		RequestID:  task.RequestID,
		CVFileName: task.CVFileName,
		Status:     "completed",
	}, true); err != nil {
		log.Printf("[Worker] Error marking task completed for %s: %v", task.RequestID, err)
	}

	log.Printf("[Worker] Finished CV: %s", task.CVFileName)
}

func (w *AnalysisWorker) failTask(ctx context.Context, task AnalysisTask, code, message, details string) {
	log.Printf("[Worker] %s for %s: %s", code, task.CVFileName, details)

	failure := models.BatchFailure{
		RequestID:    task.RequestID,
		CVFileName:   task.CVFileName,
		ErrorCode:    code,
		ErrorMessage: message,
		Details:      details,
	}
	if err := w.historyService.SaveBatchFailure(ctx, task.SessionID, task.BatchID, failure); err != nil {
		log.Printf("[Worker] Error saving batch failure for %s: %v", task.RequestID, err)
	}

	if _, err := w.historyService.UpdateBatchProgress(ctx, task.SessionID, task.BatchID, models.BatchItemStatus{
		RequestID:    task.RequestID,
		CVFileName:   task.CVFileName,
		Status:       "failed",
		ErrorCode:    code,
		ErrorMessage: message,
	}, true); err != nil {
		log.Printf("[Worker] Error updating failed task %s: %v", task.RequestID, err)
	}
}
