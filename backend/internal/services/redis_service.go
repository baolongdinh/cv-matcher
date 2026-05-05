package services

import (
	"context"
	"encoding/json"
	"fmt"

	"cv-jd-matcher/internal/models"

	"github.com/redis/go-redis/v9"
)

type HistoryService struct {
	client *redis.Client
}

func NewHistoryService(rdb *redis.Client) *HistoryService {
	return &HistoryService{
		client: rdb,
	}
}

func getHistoryKey(sessionID string) string {
	return fmt.Sprintf("cv_matcher:history:%s", sessionID)
}

func (s *HistoryService) SaveAnalysis(ctx context.Context, sessionID string, result *models.AnalysisResult) error {
	key := getHistoryKey(sessionID)
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	// Save to history hash
	err = s.client.HSet(ctx, key, result.ProcessingMetadata.RequestID, data).Err()
	if err != nil {
		return err
	}

	// Save semantic cache mapping if hash exists
	if result.FileHash != "" {
		cacheKey := fmt.Sprintf("cv_matcher:session_cache:%s:%s", sessionID, result.FileHash)
		s.client.Set(ctx, cacheKey, result.ProcessingMetadata.RequestID, 0)
	}

	return nil
}

func (s *HistoryService) GetByHash(ctx context.Context, sessionID string, hash string) (*models.AnalysisResult, error) {
	if hash == "" {
		return nil, fmt.Errorf("empty hash")
	}

	cacheKey := fmt.Sprintf("cv_matcher:session_cache:%s:%s", sessionID, hash)
	requestID, err := s.client.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	key := getHistoryKey(sessionID)
	data, err := s.client.HGet(ctx, key, requestID).Result()
	if err != nil {
		return nil, err
	}

	var result models.AnalysisResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *HistoryService) GetHistory(ctx context.Context, sessionID string) ([]models.AnalysisResult, error) {
	key := getHistoryKey(sessionID)

	resultsMap, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var history []models.AnalysisResult
	for _, jsonStr := range resultsMap {
		var res models.AnalysisResult
		if err := json.Unmarshal([]byte(jsonStr), &res); err != nil {
			continue // skip invalid entries
		}
		history = append(history, res)
	}

	return history, nil
}

func (s *HistoryService) DeleteAnalysis(ctx context.Context, sessionID string, requestID string) error {
	key := getHistoryKey(sessionID)
	return s.client.HDel(ctx, key, requestID).Err()
}
