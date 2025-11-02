package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	HTTPPort string
	Env      string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Pagination
	DefaultPageSize int
	MaxPageSize     int
}

func Load() (*Config, error) {
	// Load .env file if exists (ignore error in production)
	_ = godotenv.Load()

	defaultPageSize, _ := strconv.Atoi(getEnv("DEFAULT_PAGE_SIZE", "20"))
	maxPageSize, _ := strconv.Atoi(getEnv("MAX_PAGE_SIZE", "100"))

	return &Config{
		HTTPPort: getEnv("HTTP_PORT", "8084"),
		Env:      getEnv("ENV", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "order_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		DefaultPageSize: defaultPageSize,
		MaxPageSize:     maxPageSize,
	}, nil
}

func (c *Config) Print() {
	log.Println("=== Order Service Configuration ===")
	log.Printf("HTTP Port: %s", c.HTTPPort)
	log.Printf("Environment: %s", c.Env)
	log.Printf("Database: %s@%s:%s/%s", c.DBUser, c.DBHost, c.DBPort, c.DBName)
	log.Printf("Pagination: default=%d, max=%d", c.DefaultPageSize, c.MaxPageSize)
	log.Println("====================================")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
