package main

import (
	"bbank/config"
	"bbank/handlers"
	"bbank/middleware"
	"bbank/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	config.ConnectDatabase()

	// Load config
	cfg := config.LoadConfig()

	// Initialize services
	authService := services.NewAuthService(config.GetDB(), cfg.JWTSecret)
	balanceService := services.NewBalanceService(config.GetDB())
	transactionService := services.NewTransactionService(config.GetDB(), balanceService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	balanceHandler := handlers.NewBalanceHandler(balanceService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Setup Gin router
	r := gin.Default()

	// Middleware for logging
	r.Use(middleware.AuditLogger(config.GetDB()))
	// Middleware for CORS
	r.Use(middleware.CORSMiddleware())

	// Public routes
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		//auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(authService))
	{
		// Transaction routes
		transactions := api.Group("/transactions")
		{
			transactions.POST("/credit", transactionHandler.Credit)
			transactions.POST("/debit", transactionHandler.Debit)
			transactions.POST("/transfer", transactionHandler.Transfer)
			transactions.GET("/history", transactionHandler.GetHistory)
			transactions.GET("/:id", transactionHandler.GetTransaction)
		}

		// Balance routes
		balances := api.Group("/balances")
		{
			balances.GET("/current", balanceHandler.GetCurrentBalance)
			balances.GET("/historical", balanceHandler.GetHistoricalBalance)
			balances.GET("/at-time", balanceHandler.GetBalanceAtTime)
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("", authHandler.GetAllUsers)
			users.GET("/:id", authHandler.GetUser)
			users.PUT("/:id", authHandler.UpdateUser)
			users.DELETE("/:id", authHandler.DeleteUser)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// Start server
	r.Run(":8080")
}
