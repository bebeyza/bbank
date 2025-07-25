package handlers

import (
	"net/http"
	"time"

	"bbank/services"

	"github.com/gin-gonic/gin"
)

type BalanceHandler struct {
	balanceService *services.BalanceService
}

func NewBalanceHandler(balanceService *services.BalanceService) *BalanceHandler {
	return &BalanceHandler{
		balanceService: balanceService,
	}
}

// Get current balance
func (h *BalanceHandler) GetCurrentBalance(c *gin.Context) {
	userID := getUserIDFromContext(c)

	balance, err := h.balanceService.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balance)
}

// Get historical balance (simplified - returns current balance)
func (h *BalanceHandler) GetHistoricalBalance(c *gin.Context) {
	userID := getUserIDFromContext(c)

	balance, err := h.balanceService.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"current_balance": balance.Amount,
		"message":         "Historical balance tracking coming soon",
	})
}

// Get balance at specific time
func (h *BalanceHandler) GetBalanceAtTime(c *gin.Context) {
	userID := getUserIDFromContext(c)

	timeStr := c.Query("time")
	if timeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "time parameter required"})
		return
	}

	timestamp, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format. Use RFC3339"})
		return
	}

	balance, err := h.balanceService.GetBalanceAtTime(userID, timestamp)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance":   balance,
		"timestamp": timestamp,
	})
}
