package handlers

import (
	"net/http"

	"cv-jd-matcher/internal/models"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, statusCode int, code, message string, details interface{}) {
	c.JSON(statusCode, models.ErrorResponse{
		Status: "error",
		Error: models.APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func writeSuccess[T any](c *gin.Context, statusCode int, data T) {
	c.JSON(statusCode, models.SuccessResponse[T]{
		Status: "success",
		Data:   data,
	})
}

func requireSessionID(c *gin.Context) (string, bool) {
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		writeError(c, http.StatusBadRequest, "SESSION_REQUIRED", "X-Session-ID header is required", nil)
		return "", false
	}
	return sessionID, true
}
