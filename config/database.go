package config

import (
	"fmt"
	"log"

	"bbank/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load config
	cfg := LoadConfig()

	// Build PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	// Connect to PostgreSQL
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	// Configure connection pool
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}

	// PostgreSQL connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(0)

	DB = database

	// Auto-migrate all models
	err = DB.AutoMigrate(models.GetAllModels()...)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("PostgreSQL database connected and migrated successfully!")
}

func GetDB() *gorm.DB {
	return DB
}
