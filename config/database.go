package config

import (
	"fmt"
	"log"
	"os"

	"bbank/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Use different path for Docker
	dbPath := "banking.db"
	if os.Getenv("DOCKER_ENV") == "true" {
		dbPath = "/root/data/banking.db"
	}

	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}

	// Configure connection pool for SQLite
	sqlDB.SetMaxOpenConns(1) // SQLite works best with single connection
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(0)

	// Enable WAL mode for better concurrency
	sqlDB.Exec("PRAGMA journal_mode=WAL;")
	sqlDB.Exec("PRAGMA busy_timeout=5000;") // 5 second timeout
	sqlDB.Exec("PRAGMA synchronous=NORMAL;")
	sqlDB.Exec("PRAGMA cache_size=1000;")
	sqlDB.Exec("PRAGMA foreign_keys=ON;")
	sqlDB.Exec("PRAGMA temp_store=memory;")

	DB = database

	// Auto-migrate all models
	err = DB.AutoMigrate(models.GetAllModels()...)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("SQLite database connected and migrated successfully!")
}

func GetDB() *gorm.DB {
	return DB
}
