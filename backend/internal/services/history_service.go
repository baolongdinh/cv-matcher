package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"cv-jd-matcher/internal/models"
	"cv-jd-matcher/internal/utils"

	"github.com/redis/go-redis/v9"
)

const batchTTL = 24 * time.Hour

type HistoryService struct {
	redis *redis.Client
}

func NewHistoryService() *HistoryService {
	return &HistoryService{
		redis: utils.RedisClient,
	}
}

func (s *HistoryService) SaveResult(ctx context.Context, result *models.AnalysisResult) error {
	return s.SaveResultWithSession(ctx, result, "")
}

func (s *HistoryService) SaveResultWithSession(ctx context.Context, result *models.AnalysisResult, sessionID string) error {
	if result == nil {
		return fmt.Errorf("analysis result is required")
	}

	id := result.ProcessingMetadata.RequestID
	if id == "" {
		id = fmt.Sprintf("req_%d", time.Now().UnixNano())
		result.ProcessingMetadata.RequestID = id
	}
	if result.FileHash != "" && result.ProcessingMetadata.FileHash == "" {
		result.ProcessingMetadata.FileHash = result.FileHash
	}
	if result.ProcessingMetadata.Timestamp == "" {
		result.ProcessingMetadata.Timestamp = time.Now().Format(time.RFC3339)
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resultsKey := s.resultsKey(sessionID)
	globalRankingKey := s.rankingKey(sessionID, "")
	jobRankingKey := s.rankingKey(sessionID, result.ProcessingMetadata.JobID)

	if err := s.redis.HSet(ctx, resultsKey, id, jsonData).Err(); err != nil {
		return err
	}

	if err := s.redis.ZAdd(ctx, globalRankingKey, redis.Z{
		Score:  result.MatchingScore.Overall,
		Member: id,
	}).Err(); err != nil {
		return err
	}

	if jobRankingKey != globalRankingKey {
		if err := s.redis.ZAdd(ctx, jobRankingKey, redis.Z{
			Score:  result.MatchingScore.Overall,
			Member: id,
		}).Err(); err != nil {
			return err
		}
	}

	if result.FileHash != "" {
		if err := s.redis.HSet(ctx, "cv:hash_to_result", result.FileHash, jsonData).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (s *HistoryService) DeleteResult(ctx context.Context, id string) error {
	return s.DeleteResultWithSession(ctx, id, "")
}

func (s *HistoryService) DeleteResultWithSession(ctx context.Context, id string, sessionID string) error {
	result, err := s.GetResultByIDWithSession(ctx, id, sessionID)
	if err != nil {
		return err
	}

	globalRankingKey := s.rankingKey(sessionID, "")
	if err := s.redis.ZRem(ctx, globalRankingKey, id).Err(); err != nil {
		return err
	}

	jobRankingKey := s.rankingKey(sessionID, result.ProcessingMetadata.JobID)
	if jobRankingKey != globalRankingKey {
		if err := s.redis.ZRem(ctx, jobRankingKey, id).Err(); err != nil {
			return err
		}
	}

	return s.redis.HDel(ctx, s.resultsKey(sessionID), id).Err()
}

func (s *HistoryService) GetResultByFileHash(ctx context.Context, fileHash string) (*models.AnalysisResult, error) {
	return s.GetResultByFileHashForSession(ctx, fileHash, "")
}

func (s *HistoryService) GetResultByFileHashForSession(ctx context.Context, fileHash string, sessionID string) (*models.AnalysisResult, error) {
	if fileHash == "" {
		return nil, redis.Nil
	}

	if sessionID != "" {
		rawResults, err := s.redis.HGetAll(ctx, s.resultsKey(sessionID)).Result()
		if err == nil {
			for _, value := range rawResults {
				var result models.AnalysisResult
				if err := json.Unmarshal([]byte(value), &result); err != nil {
					continue
				}
				if result.FileHash == fileHash || result.ProcessingMetadata.FileHash == fileHash {
					return &result, nil
				}
			}
		}
	}

	jsonData, err := s.redis.HGet(ctx, "cv:hash_to_result", fileHash).Result()
	if err != nil {
		return nil, err
	}

	var result models.AnalysisResult
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *HistoryService) CloneResultToSession(ctx context.Context, source *models.AnalysisResult, sessionID, requestID, jobID, fileName string) (*models.AnalysisResult, error) {
	if source == nil {
		return nil, fmt.Errorf("source result is required")
	}

	clonedJSON, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	var cloned models.AnalysisResult
	if err := json.Unmarshal(clonedJSON, &cloned); err != nil {
		return nil, err
	}

	if requestID == "" {
		requestID = fmt.Sprintf("req_%d", time.Now().UnixNano())
	}

	cloned.ProcessingMetadata.RequestID = requestID
	cloned.ProcessingMetadata.Timestamp = time.Now().Format(time.RFC3339)
	cloned.ProcessingMetadata.Status = "completed"
	if jobID != "" {
		cloned.ProcessingMetadata.JobID = jobID
	}
	if fileName != "" {
		cloned.ProcessingMetadata.CVFileName = fileName
	}
	if cloned.FileHash != "" {
		cloned.ProcessingMetadata.FileHash = cloned.FileHash
	}

	if err := s.SaveResultWithSession(ctx, &cloned, sessionID); err != nil {
		return nil, err
	}

	return &cloned, nil
}

func (s *HistoryService) GetResultByID(ctx context.Context, requestID string) (*models.AnalysisResult, error) {
	return s.GetResultByIDWithSession(ctx, requestID, "")
}

func (s *HistoryService) GetResultByIDWithSession(ctx context.Context, requestID string, sessionID string) (*models.AnalysisResult, error) {
	jsonData, err := s.redis.HGet(ctx, s.resultsKey(sessionID), requestID).Result()
	if err != nil {
		return nil, err
	}

	var result models.AnalysisResult
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *HistoryService) GetRankedHistory(ctx context.Context, jobID string, limit int64) ([]models.AnalysisResult, error) {
	return s.GetRankedHistoryWithSession(ctx, jobID, limit, "")
}

func (s *HistoryService) GetRankedHistoryWithSession(ctx context.Context, jobID string, limit int64, sessionID string) ([]models.AnalysisResult, error) {
	if limit <= 0 {
		limit = 50
	}

	ids, err := s.redis.ZRevRange(ctx, s.rankingKey(sessionID, jobID), 0, limit-1).Result()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []models.AnalysisResult{}, nil
	}

	data, err := s.redis.HMGet(ctx, s.resultsKey(sessionID), ids...).Result()
	if err != nil {
		return nil, err
	}

	results := make([]models.AnalysisResult, 0, len(data))
	for _, item := range data {
		itemStr, ok := item.(string)
		if !ok || itemStr == "" {
			continue
		}
		var res models.AnalysisResult
		if err := json.Unmarshal([]byte(itemStr), &res); err != nil {
			continue
		}
		results = append(results, res)
	}

	return results, nil
}

func (s *HistoryService) InitBatch(ctx context.Context, sessionID, batchID, jobID string, total int) error {
	metaKey := s.batchMetaKey(sessionID, batchID)
	if err := s.redis.HSet(ctx, metaKey, map[string]interface{}{
		"batch_id":   batchID,
		"job_id":     jobID,
		"status":     "queued",
		"total":      total,
		"completed":  0,
		"successful": 0,
		"failed":     0,
		"created_at": time.Now().Format(time.RFC3339),
	}).Err(); err != nil {
		return err
	}

	return s.redis.Expire(ctx, metaKey, batchTTL).Err()
}

func (s *HistoryService) SaveBatchItemStatus(ctx context.Context, sessionID, batchID string, item models.BatchItemStatus) error {
	key := s.batchDetailsKey(sessionID, batchID)
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return err
	}
	if err := s.redis.HSet(ctx, key, item.RequestID, itemJSON).Err(); err != nil {
		return err
	}
	return s.redis.Expire(ctx, key, batchTTL).Err()
}

func (s *HistoryService) SaveBatchResult(ctx context.Context, sessionID, batchID, requestID string, result *models.AnalysisResult) error {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}
	key := s.batchResultsKey(sessionID, batchID)
	if err := s.redis.HSet(ctx, key, requestID, resultJSON).Err(); err != nil {
		return err
	}
	return s.redis.Expire(ctx, key, batchTTL).Err()
}

func (s *HistoryService) SaveBatchFailure(ctx context.Context, sessionID, batchID string, failure models.BatchFailure) error {
	failureJSON, err := json.Marshal(failure)
	if err != nil {
		return err
	}
	key := s.batchFailuresKey(sessionID, batchID)
	if err := s.redis.HSet(ctx, key, failure.RequestID, failureJSON).Err(); err != nil {
		return err
	}
	return s.redis.Expire(ctx, key, batchTTL).Err()
}

func (s *HistoryService) UpdateBatchProgress(ctx context.Context, sessionID, batchID string, item models.BatchItemStatus, terminal bool) (*models.BatchStatusData, error) {
	if err := s.SaveBatchItemStatus(ctx, sessionID, batchID, item); err != nil {
		return nil, err
	}

	metaKey := s.batchMetaKey(sessionID, batchID)
	if err := s.redis.HSet(ctx, metaKey, "status", "processing").Err(); err != nil {
		return nil, err
	}

	if terminal {
		if _, err := s.redis.HIncrBy(ctx, metaKey, "completed", 1).Result(); err != nil {
			return nil, err
		}
		if item.Status == "completed" {
			if _, err := s.redis.HIncrBy(ctx, metaKey, "successful", 1).Result(); err != nil {
				return nil, err
			}
		}
		if item.Status == "failed" {
			if _, err := s.redis.HIncrBy(ctx, metaKey, "failed", 1).Result(); err != nil {
				return nil, err
			}
		}
	}

	status, err := s.GetBatchMetaWithSession(ctx, sessionID, batchID)
	if err != nil {
		return nil, err
	}

	if status.Completed >= status.Total && status.Total > 0 {
		status.Status = "completed"
		if err := s.redis.HSet(ctx, metaKey, "status", "completed").Err(); err != nil {
			return nil, err
		}
	}

	return s.GetBatchMetaWithSession(ctx, sessionID, batchID)
}

func (s *HistoryService) GetBatchMetaWithSession(ctx context.Context, sessionID, batchID string) (*models.BatchStatusData, error) {
	meta, err := s.redis.HGetAll(ctx, s.batchMetaKey(sessionID, batchID)).Result()
	if err != nil {
		return nil, err
	}
	if len(meta) == 0 {
		return nil, redis.Nil
	}

	items, err := s.GetBatchItemsWithSession(ctx, sessionID, batchID)
	if err != nil {
		return nil, err
	}

	return &models.BatchStatusData{
		BatchID:    batchID,
		Status:     meta["status"],
		JobID:      meta["job_id"],
		Total:      atoi(meta["total"]),
		Completed:  atoi(meta["completed"]),
		Successful: atoi(meta["successful"]),
		Failed:     atoi(meta["failed"]),
		Items:      items,
	}, nil
}

func (s *HistoryService) GetBatchItemsWithSession(ctx context.Context, sessionID, batchID string) ([]models.BatchItemStatus, error) {
	rawItems, err := s.redis.HGetAll(ctx, s.batchDetailsKey(sessionID, batchID)).Result()
	if err != nil {
		return nil, err
	}

	items := make([]models.BatchItemStatus, 0, len(rawItems))
	for _, value := range rawItems {
		var item models.BatchItemStatus
		if err := json.Unmarshal([]byte(value), &item); err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *HistoryService) GetBatchResultsWithSession(ctx context.Context, sessionID, batchID string) (*models.BatchResultsData, error) {
	status, err := s.GetBatchMetaWithSession(ctx, sessionID, batchID)
	if err != nil {
		return nil, err
	}

	rawResults, err := s.redis.HGetAll(ctx, s.batchResultsKey(sessionID, batchID)).Result()
	if err != nil {
		return nil, err
	}

	rawFailures, err := s.redis.HGetAll(ctx, s.batchFailuresKey(sessionID, batchID)).Result()
	if err != nil {
		return nil, err
	}

	candidates := make([]models.AnalysisResult, 0, len(rawResults))
	for _, value := range rawResults {
		var result models.AnalysisResult
		if err := json.Unmarshal([]byte(value), &result); err != nil {
			continue
		}
		candidates = append(candidates, result)
	}

	failures := make([]models.BatchFailure, 0, len(rawFailures))
	for _, value := range rawFailures {
		var failure models.BatchFailure
		if err := json.Unmarshal([]byte(value), &failure); err != nil {
			continue
		}
		failures = append(failures, failure)
	}

	return &models.BatchResultsData{
		BatchID:     batchID,
		Candidates:  candidates,
		Failed:      failures,
		Total:       status.Total,
		Successful:  status.Successful,
		FailedCount: status.Failed,
	}, nil
}

func (s *HistoryService) resultsKey(sessionID string) string {
	if sessionID == "" {
		return "cv:results"
	}
	return fmt.Sprintf("cv:results:session:%s", sessionID)
}

func (s *HistoryService) rankingKey(sessionID, jobID string) string {
	if sessionID != "" {
		if jobID != "" {
			return fmt.Sprintf("cv:ranking:session:%s:%s", sessionID, jobID)
		}
		return fmt.Sprintf("cv:ranking:session:%s", sessionID)
	}
	if jobID != "" {
		return fmt.Sprintf("cv:ranking:%s", jobID)
	}
	return "cv:ranking:global"
}

func (s *HistoryService) batchMetaKey(sessionID, batchID string) string {
	return fmt.Sprintf("cv:batch:meta:session:%s:%s", sessionID, batchID)
}

func (s *HistoryService) batchDetailsKey(sessionID, batchID string) string {
	return fmt.Sprintf("cv:batch:details:session:%s:%s", sessionID, batchID)
}

func (s *HistoryService) batchResultsKey(sessionID, batchID string) string {
	return fmt.Sprintf("cv:batch:results:session:%s:%s", sessionID, batchID)
}

func (s *HistoryService) batchFailuresKey(sessionID, batchID string) string {
	return fmt.Sprintf("cv:batch:failures:session:%s:%s", sessionID, batchID)
}

func atoi(raw string) int {
	if raw == "" {
		return 0
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return value
}
