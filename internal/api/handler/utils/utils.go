package utils

import (
	"employee-app/internal/api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context, userType string) (uint, int, string) {
	userIDVal, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		return 0, http.StatusUnauthorized, "user ID not found in token"
	}

	role, exists := c.Get(middleware.ContextRoleKey)
	if !exists {
		return 0, http.StatusUnauthorized, "role not found in token"
	}

	if role != userType {
		return 0, http.StatusForbidden, "forbidden"
	}

	userID := userIDVal.(uint)

	return userID, http.StatusOK, ""
}