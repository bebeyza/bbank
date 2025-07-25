package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	FromUserID *uint          `json:"from_user_id"` // pointer for nullable
	ToUserID   uint           `json:"to_user_id" gorm:"not null"`
	Amount     float64        `json:"amount" gorm:"not null"`
	Type       string         `json:"type" gorm:"not null"` // credit, debit, transfer
	Status     string         `json:"status" gorm:"default:pending"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	FromUser *User `json:"from_user,omitempty" gorm:"foreignKey:FromUserID"`
	ToUser   User  `json:"to_user" gorm:"foreignKey:ToUserID"`
}
