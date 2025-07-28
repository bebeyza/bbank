package models

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID         uint           `gorm:"primaryKey"`
	EntityType string         `gorm:"not null"` // e.g. "user", "transaction", "api_request"
	EntityID   uint           `gorm:"not null"` // user_id
	Action     string         `gorm:"not null"` // e.g. "updated", "created", "called"
	Details    datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt  time.Time
}
