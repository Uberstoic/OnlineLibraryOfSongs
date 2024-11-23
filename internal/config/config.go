package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	MusicAPIURL string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USERNAME", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "music_library"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		MusicAPIURL: getEnv("MUSIC_API_URL", "http://localhost:8081"),
	}
}

// GetDBURL returns URL for database migrations
func (c *Config) GetDBURL() string {
	return fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// GetGormDBURL returns URL for GORM database connection
func (c *Config) GetGormDBURL() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
