package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server
	HTTPPort string
	GRPCPort string
	Env      string

	// Database
	DatabaseURL      string
	DBMaxIdleConns   int
	DBMaxOpenConns   int

	// JWT
	JWTSecret          string
	JWTAccessDuration  time.Duration
	JWTRefreshDuration time.Duration

	// Password
	BcryptCost int

	// CORS
	CORSAllowOrigins     []string
	CORSAllowCredentials bool

	// Logging
	LogLevel string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		HTTPPort: getEnv("HTTP_PORT", "8082"),
		GRPCPort: getEnv("GRPC_PORT", "50052"),
		Env:      getEnv("ENV", "development"),

		DatabaseURL:    getEnv("DATABASE_URL", "postgresql://bookstore:dev_password@localhost:5432/users_db?sslmode=disable"),
		DBMaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		DBMaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 100),

		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key-change-in-production-min-32-chars"),
		JWTAccessDuration:  getEnvAsDuration("JWT_ACCESS_DURATION", "24h"),
		JWTRefreshDuration: getEnvAsDuration("JWT_REFRESH_DURATION", "168h"),

		BcryptCost: getEnvAsInt("BCRYPT_COST", 12),

		CORSAllowOrigins:     getEnvAsSlice("CORS_ALLOW_ORIGINS", []string{"http://localhost:3000", "http://localhost:8080"}, ","),
		CORSAllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),

		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	// Validate required fields
	if config.JWTSecret == "your-secret-key-change-in-production-min-32-chars" && config.Env == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production")
	}

	return config, nil
}

// Helper functions

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		valueStr = defaultValue
	}
	duration, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Printf("Invalid duration for %s, using default: %s", key, defaultValue)
		duration, _ = time.ParseDuration(defaultValue)
	}
	return duration
}

func getEnvAsSlice(key string, defaultValue []string, separator string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	
	values := []string{}
	for _, v := range splitAndTrim(valueStr, separator) {
		values = append(values, v)
	}
	return values
}

func splitAndTrim(s, sep string) []string {
	parts := []string{}
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	parts := []string{}
	current := ""
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			parts = append(parts, current)
			current = ""
			i += len(sep) - 1
		} else {
			current += string(s[i])
		}
	}
	parts = append(parts, current)
	return parts
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	
	return s[start:end]
}
