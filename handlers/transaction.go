package handlers

import (
	"net/http"
	"strconv"

	"bbank/middleware"
	"bbank/services"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// Credit money to user account
func (h *TransactionHandler) Credit(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req services.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionService.Credit(userID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Credit successful",
		"transaction": transaction,
	})
}

// Debit money from user account
func (h *TransactionHandler) Debit(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req services.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionService.Debit(userID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Debit successful",
		"transaction": transaction,
	})
}

// Transfer money between users
func (h *TransactionHandler) Transfer(c *gin.Context) {
	fromUserID := getUserIDFromContext(c)

	var req services.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionService.Transfer(fromUserID, req.ToUserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Transfer successful",
		"transaction": transaction,
	})
}

// Get transaction history
func (h *TransactionHandler) GetHistory(c *gin.Context) {
	userID := getUserIDFromContext(c)

	// Get query parameters for pagination
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	transactions, err := h.transactionService.GetUserTransactions(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"count":        len(transactions),
	})
}

// Get single transaction
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromParam(c)
	if !ok {
		return
	}

	transactionIDStr := c.Param("id")
	transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := h.transactionService.GetTransaction(uint(transactionID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
