package middleware

import (
	"encoding/json"
	"time"

	"bbank/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuditLogger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // Process request

		userIDVal, exists := c.Get("user_id")
		if !exists {
			return // skip unauthenticated requests
		}

		userID, ok := userIDVal.(uint)
		if !ok {
			return
		}

		details := map[string]interface{}{
			"method":      c.Request.Method,
			"path":        c.FullPath(),
			"status_code": c.Writer.Status(),
			"duration_ms": time.Since(start).Milliseconds(),
		}

		jsonDetails, _ := json.Marshal(details)

		log := models.AuditLog{
			EntityType: "api_request",
			EntityID:   userID,
			Action:     "called",
			Details:    jsonDetails,
			CreatedAt:  time.Now(),
		}

		// Store asynchronously (non-blocking)
		go db.Create(&log)
	}
}
