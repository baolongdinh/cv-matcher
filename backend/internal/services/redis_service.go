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

	err = s.client.HSet(ctx, key, result.ProcessingMetadata.RequestID, data).Err()
	if err != nil {
		return err
	}
	return nil
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
