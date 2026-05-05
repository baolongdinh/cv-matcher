package handlers

import (
	"net/http"

	"cv-jd-matcher/internal/services"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	historyService *services.HistoryService
}

func NewHistoryHandler(historyService *services.HistoryService) *HistoryHandler {
	return &HistoryHandler{
		historyService: historyService,
	}
}

func (h *HistoryHandler) GetHistory(c *gin.Context) {
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "MISSING_SESSION_ID",
				"message": "X-Session-ID header is required",
			},
		})
		return
	}

	history, err := h.historyService.GetHistory(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "REDIS_ERROR",
				"message": "Failed to fetch history",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   history,
	})
}

func (h *HistoryHandler) DeleteHistory(c *gin.Context) {
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "MISSING_SESSION_ID",
				"message": "X-Session-ID header is required",
			},
		})
		return
	}

	requestID := c.Param("id")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "MISSING_REQUEST_ID",
				"message": "request ID param is required",
			},
		})
		return
	}

	err := h.historyService.DeleteAnalysis(c.Request.Context(), sessionID, requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": gin.H{
				"code":    "REDIS_ERROR",
				"message": "Failed to delete history item",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Item deleted",
	})
}
