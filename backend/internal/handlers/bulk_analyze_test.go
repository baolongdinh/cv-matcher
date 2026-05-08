package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"cv-jd-matcher/internal/services"
	"cv-jd-matcher/internal/utils"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func setupBulkHandlerTest(t *testing.T) *services.HistoryService {
	t.Helper()

	gin.SetMode(gin.TestMode)
	mr := miniredis.RunT(t)
	utils.RedisClient = redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() {
		_ = utils.RedisClient.Close()
		mr.Close()
	})

	return services.NewHistoryService()
}

func TestGetBatchStatusIsSessionScoped(t *testing.T) {
	historyService := setupBulkHandlerTest(t)
	require.NoError(t, historyService.InitBatch(context.Background(), "session-a", "batch-1", "jd-1", 1))

	handler := NewBulkHandler(services.NewPDFService(), historyService)

	router := gin.New()
	router.GET("/api/jobs/:batch_id", handler.GetBatchStatus)

	req := httptest.NewRequest(http.MethodGet, "/api/jobs/batch-1", nil)
	req.Header.Set("X-Session-ID", "session-b")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotFound, resp.Code)
}
