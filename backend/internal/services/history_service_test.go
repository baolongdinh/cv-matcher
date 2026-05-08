package services

import (
	"context"
	"testing"
	"time"

	"cv-jd-matcher/internal/models"
	"cv-jd-matcher/internal/utils"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func setupHistoryServiceTest(t *testing.T) *HistoryService {
	t.Helper()

	mr := miniredis.RunT(t)
	utils.RedisClient = redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() {
		_ = utils.RedisClient.Close()
		mr.Close()
	})

	return NewHistoryService()
}

func sampleResult(requestID, fileHash, fileName, jobID string, score float64) *models.AnalysisResult {
	return &models.AnalysisResult{
		ExecutiveSummary: "summary",
		FileHash:         fileHash,
		MatchingScore: models.MatchingScore{
			Overall: score,
		},
		QualityScore: models.QualityScore{
			Overall: score - 1,
		},
		TechnicalSkills: []models.SkillScore{{Name: "Go", Score: score}},
		ProcessingMetadata: models.ProcessingMetadata{
			RequestID:  requestID,
			CVFileName: fileName,
			JobID:      jobID,
			Status:     "completed",
			Timestamp:  time.Now().Format(time.RFC3339),
			FileHash:   fileHash,
		},
	}
}

func TestSaveLoadDeleteResultWithSession(t *testing.T) {
	service := setupHistoryServiceTest(t)
	ctx := context.Background()
	result := sampleResult("req_1", "hash_1", "alice.pdf", "fe-role", 8.4)

	require.NoError(t, service.SaveResultWithSession(ctx, result, "session-a"))

	history, err := service.GetRankedHistoryWithSession(ctx, "", 10, "session-a")
	require.NoError(t, err)
	require.Len(t, history, 1)
	require.Equal(t, "req_1", history[0].ProcessingMetadata.RequestID)

	require.NoError(t, service.DeleteResultWithSession(ctx, "req_1", "session-a"))

	history, err = service.GetRankedHistoryWithSession(ctx, "", 10, "session-a")
	require.NoError(t, err)
	require.Len(t, history, 0)
}

func TestCloneResultToSessionCreatesNewScopedRecord(t *testing.T) {
	service := setupHistoryServiceTest(t)
	ctx := context.Background()

	original := sampleResult("req_original", "hash_cache", "bob.pdf", "be-role", 9.1)
	require.NoError(t, service.SaveResultWithSession(ctx, original, "session-a"))

	cached, err := service.GetResultByFileHashForSession(ctx, "hash_cache", "session-b")
	require.NoError(t, err)
	require.Equal(t, "req_original", cached.ProcessingMetadata.RequestID)

	cloned, err := service.CloneResultToSession(ctx, cached, "session-b", "req_cloned", "be-role", "bob-copy.pdf")
	require.NoError(t, err)
	require.Equal(t, "req_cloned", cloned.ProcessingMetadata.RequestID)
	require.Equal(t, "bob-copy.pdf", cloned.ProcessingMetadata.CVFileName)

	sessionAHistory, err := service.GetRankedHistoryWithSession(ctx, "", 10, "session-a")
	require.NoError(t, err)
	require.Len(t, sessionAHistory, 1)

	sessionBHistory, err := service.GetRankedHistoryWithSession(ctx, "", 10, "session-b")
	require.NoError(t, err)
	require.Len(t, sessionBHistory, 1)
	require.Equal(t, "req_cloned", sessionBHistory[0].ProcessingMetadata.RequestID)
}

func TestBatchProgressAndResultsStaySessionScoped(t *testing.T) {
	service := setupHistoryServiceTest(t)
	ctx := context.Background()

	require.NoError(t, service.InitBatch(ctx, "session-a", "batch-1", "jd-1", 2))
	_, err := service.UpdateBatchProgress(ctx, "session-a", "batch-1", models.BatchItemStatus{
		RequestID:  "req_1",
		CVFileName: "alice.pdf",
		Status:     "completed",
	}, true)
	require.NoError(t, err)

	require.NoError(t, service.SaveBatchResult(ctx, "session-a", "batch-1", "req_1", sampleResult("req_1", "hash_1", "alice.pdf", "jd-1", 8.0)))
	require.NoError(t, service.SaveBatchFailure(ctx, "session-a", "batch-1", models.BatchFailure{
		RequestID:    "req_2",
		CVFileName:   "bob.pdf",
		ErrorCode:    "PDF_PARSE_ERROR",
		ErrorMessage: "Failed to extract text from PDF",
	}))
	_, err = service.UpdateBatchProgress(ctx, "session-a", "batch-1", models.BatchItemStatus{
		RequestID:    "req_2",
		CVFileName:   "bob.pdf",
		Status:       "failed",
		ErrorCode:    "PDF_PARSE_ERROR",
		ErrorMessage: "Failed to extract text from PDF",
	}, true)
	require.NoError(t, err)

	status, err := service.GetBatchMetaWithSession(ctx, "session-a", "batch-1")
	require.NoError(t, err)
	require.Equal(t, 2, status.Total)
	require.Equal(t, 2, status.Completed)
	require.Equal(t, 1, status.Successful)
	require.Equal(t, 1, status.Failed)
	require.Equal(t, "completed", status.Status)

	results, err := service.GetBatchResultsWithSession(ctx, "session-a", "batch-1")
	require.NoError(t, err)
	require.Len(t, results.Candidates, 1)
	require.Len(t, results.Failed, 1)
	require.Equal(t, 1, results.FailedCount)

	_, err = service.GetBatchMetaWithSession(ctx, "session-b", "batch-1")
	require.Error(t, err)
}
