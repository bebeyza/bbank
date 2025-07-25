package services

import (
	"errors"
	"time"

	"bbank/models"

	"gorm.io/gorm"
)

type BalanceService struct {
	db *gorm.DB
}

func NewBalanceService(db *gorm.DB) *BalanceService {
	return &BalanceService{db: db}
}

// Get current balance for user
func (s *BalanceService) GetBalance(userID uint) (*models.Balance, error) {
	var balance models.Balance
	if err := s.db.Where("user_id = ?", userID).First(&balance).Error; err != nil {
		return nil, errors.New("balance not found")
	}
	return &balance, nil
}

// Update balance (thread-safe with database transaction)
func (s *BalanceService) UpdateBalance(userID uint, amount float64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var balance models.Balance
		if err := tx.Where("user_id = ?", userID).First(&balance).Error; err != nil {
			return err
		}

		// Check if sufficient funds for debit
		if balance.Amount+amount < 0 {
			return errors.New("insufficient funds")
		}

		// Update balance
		balance.Amount += amount
		balance.LastUpdatedAt = time.Now()

		return tx.Save(&balance).Error
	})
}

// Get balance at specific time (simplified version)
func (s *BalanceService) GetBalanceAtTime(userID uint, timestamp time.Time) (float64, error) {
	// For now, return current balance
	// In production, you'd calculate from transaction history
	balance, err := s.GetBalance(userID)
	if err != nil {
		return 0, err
	}
	return balance.Amount, nil
}
