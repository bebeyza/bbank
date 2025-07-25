package models

import (
	"time"
)

type Balance struct {
	UserID        uint      `json:"user_id" gorm:"primaryKey"`
	Amount        float64   `json:"amount" gorm:"default:0"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	User          User      `json:"user" gorm:"foreignKey:UserID"`
}
