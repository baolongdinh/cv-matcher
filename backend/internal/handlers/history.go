package handlers

import (
	"net/http"
	"strconv"

	"cv-jd-matcher/internal/services"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	historyService *services.HistoryService
}

func NewHistoryHandler(hs *services.HistoryService) *HistoryHandler {
	return &HistoryHandler{historyService: hs}
}

func (h *HistoryHandler) GetHistory(c *gin.Context) {
	sessionID, ok := requireSessionID(c)
	if !ok {
		return
	}

	jobID := c.Query("job_id")
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	results, err := h.historyService.GetRankedHistoryWithSession(c.Request.Context(), jobID, limit, sessionID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "HISTORY_FETCH_ERROR", "Failed to fetch history", err.Error())
		return
	}

	writeSuccess(c, http.StatusOK, results)
}

func (h *HistoryHandler) DeleteHistory(c *gin.Context) {
	sessionID, ok := requireSessionID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	if err := h.historyService.DeleteResultWithSession(c.Request.Context(), id, sessionID); err != nil {
		writeError(c, http.StatusInternalServerError, "HISTORY_DELETE_ERROR", "Failed to delete record", err.Error())
		return
	}

	writeSuccess(c, http.StatusOK, gin.H{"message": "Record deleted"})
}
