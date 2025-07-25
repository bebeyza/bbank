package services

import (
	"errors"
	"time"

	"bbank/models"

	"gorm.io/gorm"
)

type TransactionService struct {
	db             *gorm.DB
	balanceService *BalanceService
}

type TransactionRequest struct {
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	ToUserID *uint   `json:"to_user_id,omitempty"`
}

type TransferRequest struct {
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	ToUserID uint    `json:"to_user_id" binding:"required"`
}

func NewTransactionService(db *gorm.DB, balanceService *BalanceService) *TransactionService {
	return &TransactionService{
		db:             db,
		balanceService: balanceService,
	}
}

// Credit money to user account
func (s *TransactionService) Credit(userID uint, amount float64) (*models.Transaction, error) {
	var transaction models.Transaction

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Create transaction record
		transaction = models.Transaction{
			ToUserID: userID,
			Amount:   amount,
			Type:     "credit",
			Status:   "completed",
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// Update balance directly in this transaction (avoid nested transaction)
		var balance models.Balance
		if err := tx.Where("user_id = ?", userID).First(&balance).Error; err != nil {
			return err
		}

		balance.Amount += amount
		balance.LastUpdatedAt = time.Now()

		if err := tx.Save(&balance).Error; err != nil {
			return err
		}

		return nil
	})

	return &transaction, err
}

// Debit money from user account
func (s *TransactionService) Debit(userID uint, amount float64) (*models.Transaction, error) {
	var transaction models.Transaction

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Check balance and update in same transaction
		var balance models.Balance
		if err := tx.Where("user_id = ?", userID).First(&balance).Error; err != nil {
			return err
		}

		if balance.Amount < amount {
			return errors.New("insufficient funds")
		}

		// Create transaction record
		transaction = models.Transaction{
			FromUserID: &userID,
			ToUserID:   userID,
			Amount:     amount,
			Type:       "debit",
			Status:     "completed",
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// Update balance
		balance.Amount -= amount
		balance.LastUpdatedAt = time.Now()

		if err := tx.Save(&balance).Error; err != nil {
			return err
		}

		return nil
	})

	return &transaction, err
}

// Transfer money between users
func (s *TransactionService) Transfer(fromUserID, toUserID uint, amount float64) (*models.Transaction, error) {
	if fromUserID == toUserID {
		return nil, errors.New("cannot transfer to same account")
	}

	var transaction models.Transaction

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Get both balances in single transaction
		var fromBalance, toBalance models.Balance

		if err := tx.Where("user_id = ?", fromUserID).First(&fromBalance).Error; err != nil {
			return errors.New("sender account not found")
		}

		if err := tx.Where("user_id = ?", toUserID).First(&toBalance).Error; err != nil {
			return errors.New("recipient account not found")
		}

		if fromBalance.Amount < amount {
			return errors.New("insufficient funds")
		}

		// Create transaction record
		transaction = models.Transaction{
			FromUserID: &fromUserID,
			ToUserID:   toUserID,
			Amount:     amount,
			Type:       "transfer",
			Status:     "completed",
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// Update both balances
		fromBalance.Amount -= amount
		fromBalance.LastUpdatedAt = time.Now()

		toBalance.Amount += amount
		toBalance.LastUpdatedAt = time.Now()

		if err := tx.Save(&fromBalance).Error; err != nil {
			return err
		}

		if err := tx.Save(&toBalance).Error; err != nil {
			return err
		}

		return nil
	})

	return &transaction, err
}

// Get transaction history for user
func (s *TransactionService) GetUserTransactions(userID uint, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := s.db.Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Preload("FromUser").
		Preload("ToUser").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

// Get single transaction by ID
func (s *TransactionService) GetTransaction(transactionID, userID uint) (*models.Transaction, error) {
	var transaction models.Transaction

	err := s.db.Where("id = ? AND (from_user_id = ? OR to_user_id = ?)",
		transactionID, userID, userID).
		Preload("FromUser").
		Preload("ToUser").
		First(&transaction).Error

	if err != nil {
		return nil, errors.New("transaction not found")
	}

	return &transaction, nil
}
