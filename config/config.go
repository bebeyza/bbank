package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	JWTSecret  string
	ServerPort string
}

func LoadConfig() *Config {
	// Set config file
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		JWTSecret:  getEnv("JWT_SECRET", "change-this-secret"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	return config
}

// Helper function to get env with default
func getEnv(key, defaultValue string) string {
	if value := viper.GetString(key); value != "" {
		return value
	}
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
