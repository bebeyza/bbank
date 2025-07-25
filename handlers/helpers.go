package handlers

import (
	"github.com/gin-gonic/gin"
)

// Get user ID from gin context (set by auth middleware)
func getUserIDFromContext(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}
